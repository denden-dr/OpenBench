package ticket

import (
	"context"
	"crypto/rand"
	"database/sql"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/denden-dr/openbench/apps/backend/internal/database"
	"github.com/denden-dr/openbench/apps/backend/internal/pkg/api"
	"github.com/jmoiron/sqlx"
)

var (
	ErrTicketNotFound = errors.New("ticket not found")
	ErrInvalidInput   = errors.New("invalid input data")
)

type TicketService interface {
	GetTicketByNumber(ctx context.Context, ticketNumber string) (*Ticket, error)
}

type AdminTicketService interface {
	TicketService
	GetTicket(ctx context.Context, id string) (*Ticket, error)
	CreateTicket(ctx context.Context, req *api.TicketCreate) (*Ticket, error)
	ListTickets(ctx context.Context) ([]*Ticket, error)
	UpdateTicket(ctx context.Context, id string, req *api.TicketUpdate) (*Ticket, error)
	ListWarranties(ctx context.Context) ([]*Warranty, error)
	EmergencyUpdateTicket(ctx context.Context, id string, req *api.TicketUpdate) (*Ticket, error)
}

type ticketService struct {
	repo TicketRepository
	db   *database.Database
}

// NewService creates a new ticket service for public access
func NewService(repo TicketRepository, db *database.Database) TicketService {
	return &ticketService{
		repo: repo,
		db:   db,
	}
}

// NewAdminService creates a new ticket service for admin access
func NewAdminService(repo TicketRepository, db *database.Database) AdminTicketService {
	return &ticketService{
		repo: repo,
		db:   db,
	}
}

func (s *ticketService) CreateTicket(ctx context.Context, req *api.TicketCreate) (*Ticket, error) {
	// Validate invariants using NewTicket before starting a transaction
	if _, err := NewTicket("", req.CustomerName, req.CustomerPhone, req.BrandPhone, req.ModelPhone, "", req.DamageDescription, "", 0, req.WarrantyDurationDays); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidInput, err)
	}

	tx, err := s.db.DB.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	ticketNumber, err := s.generateTicketNumber(ctx, tx)
	if err != nil {
		return nil, err
	}

	serialNum := ""
	if req.SerialNumber != nil {
		serialNum = *req.SerialNumber
	}

	var costVal float64 = 0.0
	if req.Cost != nil {
		costVal = float64(*req.Cost)
	}

	repairAct := ""
	if req.RepairAction != nil {
		repairAct = *req.RepairAction
	}

	t, err := NewTicket(
		ticketNumber,
		req.CustomerName,
		req.CustomerPhone,
		req.BrandPhone,
		req.ModelPhone,
		serialNum,
		req.DamageDescription,
		repairAct,
		costVal,
		req.WarrantyDurationDays,
	)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidInput, err)
	}

	if err := s.repo.Create(ctx, tx, t); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return t, nil
}

func (s *ticketService) GetTicket(ctx context.Context, id string) (*Ticket, error) {
	t, err := s.repo.GetByID(ctx, nil, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrTicketNotFound
		}
		return nil, err
	}
	return t, nil
}

func (s *ticketService) ListTickets(ctx context.Context) ([]*Ticket, error) {
	return s.repo.List(ctx, nil)
}

func (s *ticketService) UpdateTicket(ctx context.Context, id string, req *api.TicketUpdate) (*Ticket, error) {
	tx, err := s.db.DB.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	t, err := s.repo.GetByIDWithLock(ctx, tx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrTicketNotFound
		}
		return nil, err
	}

	// Normal updates: disallow reversal if device has already been picked up
	if t.DevicePosition == PositionPickedUp {
		if req.DevicePosition != nil && *req.DevicePosition != PositionPickedUp {
			return nil, fmt.Errorf("%w: %v", ErrInvalidInput, ErrStatusReversalNotAllowed)
		}
		if req.Status != nil && *req.Status != StatusCompleted && *req.Status != StatusCancelled {
			return nil, fmt.Errorf("%w: %v", ErrInvalidInput, ErrStatusReversalNotAllowed)
		}
	}

	if req.CustomerName != nil {
		t.CustomerName = *req.CustomerName
		if t.Warranty != nil {
			t.Warranty.CustomerName = *req.CustomerName
		}
	}
	if req.CustomerPhone != nil {
		t.CustomerPhone = *req.CustomerPhone
	}
	if req.BrandPhone != nil || req.ModelPhone != nil {
		if req.BrandPhone != nil {
			t.BrandPhone = *req.BrandPhone
		}
		if req.ModelPhone != nil {
			t.ModelPhone = *req.ModelPhone
		}
		if t.Warranty != nil {
			t.Warranty.DeviceInfo = fmt.Sprintf("%s %s", t.BrandPhone, t.ModelPhone)
		}
	}
	if req.SerialNumber != nil {
		t.SerialNumber = *req.SerialNumber
	}

	if req.DamageDescription != nil {
		t.DamageDescription = *req.DamageDescription
	}
	if req.RepairAction != nil {
		t.RepairAction = *req.RepairAction
	}
	if req.Cost != nil {
		t.Cost = float64(*req.Cost)
	}
	if req.WarrantyDurationDays != nil {
		if err := t.UpdateWarrantyDuration(*req.WarrantyDurationDays); err != nil {
			return nil, fmt.Errorf("%w: %v", ErrInvalidInput, err)
		}
	}
	if req.PaymentStatus != nil {
		t.PaymentStatus = string(*req.PaymentStatus)
	}
	if req.PaymentMethod != nil {
		pm := string(*req.PaymentMethod)
		t.PaymentMethod = &pm
	}
	if req.Status != nil {
		t.Status = string(*req.Status)
	}

	if req.DevicePosition != nil {
		newPos := string(*req.DevicePosition)
		if newPos == PositionPickedUp && t.DevicePosition != PositionPickedUp {
			if err := t.ProcessPickup(time.Now()); err != nil {
				return nil, fmt.Errorf("%w: %v", ErrInvalidInput, err)
			}
		} else {
			t.DevicePosition = newPos
		}
	}

	if err := t.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidInput, err)
	}

	if err := s.repo.Update(ctx, tx, t); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return t, nil
}

func (s *ticketService) ListWarranties(ctx context.Context) ([]*Warranty, error) {
	return s.repo.ListWarranties(ctx, nil)
}

func (s *ticketService) EmergencyUpdateTicket(ctx context.Context, id string, req *api.TicketUpdate) (*Ticket, error) {
	tx, err := s.db.DB.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	t, err := s.repo.GetByIDWithLock(ctx, tx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrTicketNotFound
		}
		return nil, err
	}

	isReversal := t.DevicePosition == PositionPickedUp && req.DevicePosition != nil && *req.DevicePosition != PositionPickedUp

	if req.CustomerName != nil {
		t.CustomerName = *req.CustomerName
		if t.Warranty != nil {
			t.Warranty.CustomerName = *req.CustomerName
		}
	}
	if req.CustomerPhone != nil {
		t.CustomerPhone = *req.CustomerPhone
	}
	if req.BrandPhone != nil || req.ModelPhone != nil {
		if req.BrandPhone != nil {
			t.BrandPhone = *req.BrandPhone
		}
		if req.ModelPhone != nil {
			t.ModelPhone = *req.ModelPhone
		}
		if t.Warranty != nil {
			t.Warranty.DeviceInfo = fmt.Sprintf("%s %s", t.BrandPhone, t.ModelPhone)
		}
	}
	if req.SerialNumber != nil {
		t.SerialNumber = *req.SerialNumber
	}

	if req.DamageDescription != nil {
		t.DamageDescription = *req.DamageDescription
	}
	if req.RepairAction != nil {
		t.RepairAction = *req.RepairAction
	}
	if req.Cost != nil {
		t.Cost = float64(*req.Cost)
	}
	if req.WarrantyDurationDays != nil {
		if err := t.EmergencyUpdateWarrantyDuration(*req.WarrantyDurationDays); err != nil {
			return nil, fmt.Errorf("%w: %v", ErrInvalidInput, err)
		}
	}
	if req.PaymentStatus != nil {
		t.PaymentStatus = string(*req.PaymentStatus)
	}
	if req.PaymentMethod != nil {
		pm := string(*req.PaymentMethod)
		t.PaymentMethod = &pm
	}
	if req.Status != nil {
		t.Status = string(*req.Status)
	}

	if req.DevicePosition != nil {
		newPos := string(*req.DevicePosition)
		if newPos == PositionPickedUp && t.DevicePosition != PositionPickedUp {
			if err := t.ProcessPickup(time.Now()); err != nil {
				return nil, fmt.Errorf("%w: %v", ErrInvalidInput, err)
			}
		} else if isReversal {
			t.ReversePickupLocation()
		} else {
			t.DevicePosition = newPos
		}
	}

	if err := t.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidInput, err)
	}

	if err := s.repo.Update(ctx, tx, t); err != nil {
		return nil, err
	}

	if isReversal {
		if err := s.repo.DeleteWarrantyByTicketID(ctx, tx, t.ID); err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return t, nil
}

func (s *ticketService) GetTicketByNumber(ctx context.Context, ticketNumber string) (*Ticket, error) {
	t, err := s.repo.GetByTicketNumber(ctx, nil, ticketNumber)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrTicketNotFound
		}
		return nil, err
	}
	return t, nil
}

func (s *ticketService) generateTicketNumber(ctx context.Context, tx *sqlx.Tx) (string, error) {
	prefix := fmt.Sprintf("OB-%s-", time.Now().Format("200601"))
	maxNum, err := s.repo.GetMaxTicketNumberByPrefix(ctx, tx, prefix)
	if err != nil {
		return "", err
	}

	count := 1
	if maxNum != "" {
		parts := strings.Split(maxNum, "-")
		if len(parts) >= 3 {
			seqPart := parts[2]
			if val, err := strconv.Atoi(seqPart); err == nil {
				count = val + 1
			}
		}
	}

	// Generate an 8-character random alphanumeric suffix
	suffix, err := generateRandomSuffix(8)
	if err != nil {
		return "", fmt.Errorf("failed to generate random suffix: %w", err)
	}

	return fmt.Sprintf("%s%04d-%s", prefix, count, suffix), nil
}

func generateRandomSuffix(length int) (string, error) {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		b[i] = charset[n.Int64()]
	}
	return string(b), nil
}
