package service

import (
	"context"
	"errors"
	"math"
	"regexp"
	"strings"

	"github.com/denden-dr/openbench/apps/backend/internal/dto"
	"github.com/denden-dr/openbench/apps/backend/internal/model"
	"github.com/denden-dr/openbench/apps/backend/internal/repository"
	"github.com/go-playground/validator/v10"
	"golang.org/x/sync/errgroup"
)

type TicketService interface {
	CreateTicket(ctx context.Context, req *dto.CreateTicketRequest) (*dto.TicketResponse, error)
	GetTicket(ctx context.Context, id string) (*dto.TicketResponse, error)
	UpdateTicket(ctx context.Context, id string, req *dto.UpdateTicketRequest) (*dto.TicketResponse, error)
	ListTickets(ctx context.Context, page int, limit int, search string, status string) (*dto.PaginatedTicketsResult, error)
	DeleteTicket(ctx context.Context, id string) error
	GetPublicTicket(ctx context.Context, id string) (*dto.PublicTicketResponse, error)
	TrackPublicTicket(ctx context.Context, req *dto.PublicTrackRequest) (string, error)
}

const (
	defaultTicketLimit = 20
	maxTicketLimit     = 100
	maxTicketPage      = 10000
)

type ticketService struct {
	repo     repository.TicketRepository
	validate *validator.Validate
}

func NewTicketService(repo repository.TicketRepository) TicketService {
	return &ticketService{
		repo:     repo,
		validate: validator.New(),
	}
}

func (s *ticketService) CreateTicket(ctx context.Context, req *dto.CreateTicketRequest) (*dto.TicketResponse, error) {
	if err := s.validate.Struct(req); err != nil {
		return nil, err
	}

	ticket := &model.Ticket{
		CustomerName:   req.CustomerName,
		CustomerPhone:  req.CustomerPhone,
		CustomerGender: req.CustomerGender,
		Brand:          req.Brand,
		Model:          req.Model,
		Issue:          req.Issue,
		Price:          req.Price,
	}
	if req.WarrantyDays != nil {
		ticket.WarrantyDays = *req.WarrantyDays
	} else {
		ticket.WarrantyDays = model.DefaultWarrantyDays
	}
	if req.AdditionalDescription != "" {
		ticket.AdditionalDescription = &req.AdditionalDescription
	}
	if req.Accessories != "" {
		ticket.Accessories = &req.Accessories
	}

	if err := ticket.PrepareForCreate(); err != nil {
		return nil, MapModelError(err)
	}

	if err := s.repo.Create(ctx, ticket); err != nil {
		return nil, MapRepositoryError(err)
	}

	return s.mapToResponse(ticket), nil
}

func (s *ticketService) GetTicket(ctx context.Context, id string) (*dto.TicketResponse, error) {
	ticket, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, MapRepositoryError(err)
	}
	return s.mapToResponse(ticket), nil
}

func (s *ticketService) UpdateTicket(ctx context.Context, id string, req *dto.UpdateTicketRequest) (*dto.TicketResponse, error) {
	if err := s.validate.Struct(req); err != nil {
		return nil, err
	}

	update := model.TicketUpdate{
		CustomerName:          req.CustomerName,
		CustomerPhone:         req.CustomerPhone,
		CustomerGender:        req.CustomerGender,
		Brand:                 req.Brand,
		Model:                 req.Model,
		Issue:                 req.Issue,
		AdditionalDescription: req.AdditionalDescription,
		Accessories:           req.Accessories,
		Price:                 req.Price,
		Status:                req.Status,
		PaymentStatus:         req.PaymentStatus,
		WarrantyDays:          req.WarrantyDays,
		ExitDate:              req.ExitDate,
	}
	if err := model.ValidateTicketUpdate(update); err != nil {
		return nil, MapModelError(err)
	}

	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return nil, MapRepositoryError(err)
	}
	defer rollbackTx(tx)

	ticket, err := s.repo.GetByIDForUpdateTx(ctx, tx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrTicketNotFound
		}
		return nil, MapRepositoryError(err)
	}

	if err := ticket.ApplyUpdate(update); err != nil {
		return nil, MapModelError(err)
	}

	if err := s.repo.UpdateTx(ctx, tx, ticket); err != nil {
		return nil, MapRepositoryError(err)
	}

	if err := tx.Commit(); err != nil {
		return nil, MapRepositoryError(err)
	}

	return s.mapToResponse(ticket), nil
}

func (s *ticketService) ListTickets(ctx context.Context, page int, limit int, search string, status string) (*dto.PaginatedTicketsResult, error) {
	status = strings.TrimSpace(status)
	if status != "" && status != "all" {
		validStatuses := map[string]bool{
			"service_in":           true,
			"on_process":           true,
			"waiting_confirmation": true,
			"cancelled":            true,
			"fixed":                true,
			"picked_up":            true,
		}
		if !validStatuses[status] {
			return nil, ErrInvalidStatus
		}
	}

	if page < 1 {
		page = 1
	} else if page > maxTicketPage {
		page = maxTicketPage
	}
	if limit <= 0 {
		limit = defaultTicketLimit
	} else if limit > maxTicketLimit {
		limit = maxTicketLimit
	}

	offset := (page - 1) * limit

	var total int64
	var tickets []model.Ticket
	var countsMap map[string]int64

	g, ctxGroup := errgroup.WithContext(ctx)

	g.Go(func() error {
		var err error
		total, err = s.repo.CountPaginated(ctxGroup, search, status)
		return err
	})

	g.Go(func() error {
		var err error
		tickets, err = s.repo.ListPaginated(ctxGroup, search, status, limit, offset)
		return err
	})

	g.Go(func() error {
		var err error
		countsMap, err = s.repo.GetStatusCounts(ctxGroup, search)
		return err
	})

	if err := g.Wait(); err != nil {
		return nil, MapRepositoryError(err)
	}

	responses := make([]dto.TicketResponse, len(tickets))
	for i, t := range tickets {
		responses[i] = *s.mapToResponse(&t)
	}

	allStatuses := []string{"service_in", "on_process", "waiting_confirmation", "cancelled", "fixed", "picked_up"}
	resCounts := make(map[string]int64)
	var activeSum int64
	for _, st := range allStatuses {
		cnt := countsMap[st]
		resCounts[st] = cnt
		if st != "picked_up" {
			activeSum += cnt
		}
	}
	resCounts["all"] = activeSum

	totalPages := int64(math.Ceil(float64(total) / float64(limit)))
	if totalPages < 1 {
		totalPages = 1
	}

	return &dto.PaginatedTicketsResult{
		Data:         responses,
		Total:        total,
		TotalPages:   totalPages,
		Page:         page,
		Limit:        limit,
		StatusCounts: resCounts,
	}, nil
}

func (s *ticketService) DeleteTicket(ctx context.Context, id string) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		return MapRepositoryError(err)
	}
	return nil
}

func (s *ticketService) GetPublicTicket(ctx context.Context, id string) (*dto.PublicTicketResponse, error) {
	ticket, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrTicketNotFound
		}
		return nil, MapRepositoryError(err)
	}

	res := &dto.PublicTicketResponse{
		ID:                    ticket.ID,
		CustomerNameMasked:    maskName(ticket.CustomerName),
		CustomerPhoneMasked:   maskPhone(ticket.CustomerPhone),
		Brand:                 ticket.Brand,
		Model:                 ticket.Model,
		Issue:                 ticket.Issue,
		Status:                string(ticket.Status),
		EntryDate:             ticket.EntryDate,
		ExitDate:              ticket.ExitDate,
		WarrantyDays:          ticket.WarrantyDays,
		PaymentStatus:         string(ticket.PaymentStatus),
		Price:                 &ticket.Price,
		AdditionalDescription: ticket.AdditionalDescription,
		Accessories:           ticket.Accessories,
		WarrantyExpiryDate:    ticket.WarrantyExpiryDate(),
	}
	return res, nil
}

func (s *ticketService) TrackPublicTicket(ctx context.Context, req *dto.PublicTrackRequest) (string, error) {
	if err := s.validate.Struct(req); err != nil {
		return "", err
	}

	tickets, err := s.repo.GetByShortID(ctx, req.ShortID)
	if err != nil {
		return "", MapRepositoryError(err)
	}
	if len(tickets) == 0 {
		return "", ErrTicketNotFound
	}

	normalizedInputPhone := normalizePhone(req.Phone)
	if normalizedInputPhone == "" {
		return "", ErrTicketNotFound
	}

	var match *model.Ticket
	for i := range tickets {
		t := &tickets[i]
		if t.CustomerPhone == "" {
			continue
		}
		if normalizePhone(t.CustomerPhone) == normalizedInputPhone {
			match = t
			break
		}
	}

	if match == nil {
		return "", ErrTicketNotFound
	}

	return match.ID, nil
}

func (s *ticketService) mapToResponse(ticket *model.Ticket) *dto.TicketResponse {
	return MapTicketToResponse(ticket)
}

func MapTicketToResponse(ticket *model.Ticket) *dto.TicketResponse {
	if ticket == nil {
		return nil
	}
	return &dto.TicketResponse{
		ID:                    ticket.ID,
		CustomerName:          ticket.CustomerName,
		CustomerPhone:         ticket.CustomerPhone,
		CustomerGender:        ticket.CustomerGender,
		Brand:                 ticket.Brand,
		Model:                 ticket.Model,
		Issue:                 ticket.Issue,
		AdditionalDescription: ticket.AdditionalDescription,
		Accessories:           ticket.Accessories,
		Price:                 ticket.Price,
		Status:                string(ticket.Status),
		PaymentStatus:         string(ticket.PaymentStatus),
		WarrantyDays:          ticket.WarrantyDays,
		EntryDate:             ticket.EntryDate,
		ExitDate:              ticket.ExitDate,
		WarrantyExpiryDate:    ticket.WarrantyExpiryDate(),
		IsWarranty:            ticket.IsWarranty,
		ParentTicketID:        ticket.ParentTicketID,
	}
}

// Masking and normalization helpers

func maskName(name string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		return ""
	}
	parts := strings.Fields(name)
	maskedParts := make([]string, len(parts))
	for i, part := range parts {
		runes := []rune(part)
		if len(runes) == 0 {
			continue
		}
		masked := string(runes[0])
		if len(runes) > 1 {
			masked += strings.Repeat("*", len(runes)-1)
		}
		maskedParts[i] = masked
	}
	return strings.Join(maskedParts, " ")
}

func maskPhone(phone string) string {
	phone = normalizePhone(phone)
	runes := []rune(phone)
	l := len(runes)
	if l == 0 {
		return ""
	}
	if l <= 4 {
		return strings.Repeat("*", l)
	}
	if l <= 6 {
		return string(runes[:2]) + strings.Repeat("*", l-3) + string(runes[l-1:])
	}
	return string(runes[:4]) + strings.Repeat("*", l-6) + string(runes[l-2:])
}

var nonDigitsRegexp = regexp.MustCompile(`\D`)

func normalizePhone(phone string) string {
	cleaned := nonDigitsRegexp.ReplaceAllString(phone, "")
	if strings.HasPrefix(cleaned, "62") {
		cleaned = "0" + cleaned[2:]
	}
	return cleaned
}
