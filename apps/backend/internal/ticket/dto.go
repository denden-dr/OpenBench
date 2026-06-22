package ticket

import (
	"strings"

	"github.com/denden-dr/openbench/apps/backend/internal/pkg/api"
	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// ToTicketAPI converts a domain Ticket to the API Ticket representation
func ToTicketAPI(t *Ticket) api.Ticket {
	costVal := float32(t.Cost)
	serialVal := t.SerialNumber
	repairVal := t.RepairAction

	var paymentMethod *api.TicketPaymentMethod
	if t.PaymentMethod != nil {
		pm := api.TicketPaymentMethod(*t.PaymentMethod)
		paymentMethod = &pm
	}

	idUUID, _ := uuid.Parse(t.ID)

	return api.Ticket{
		Id:                   idUUID,
		TicketNumber:         t.TicketNumber,
		CustomerName:         t.CustomerName,
		CustomerPhone:        t.CustomerPhone,
		BrandPhone:           t.BrandPhone,
		ModelPhone:           t.ModelPhone,
		SerialNumber:         &serialVal,
		DamageDescription:    t.DamageDescription,
		RepairAction:         &repairVal,
		Cost:                 &costVal,
		Status:               api.TicketStatus(t.Status),
		DevicePosition:       api.TicketDevicePosition(t.DevicePosition),
		PaymentStatus:        api.TicketPaymentStatus(t.PaymentStatus),
		PaymentMethod:        paymentMethod,
		WarrantyDurationDays: t.WarrantyDurationDays,
		PickedUpAt:           t.PickedUpAt,
		CreatedAt:            t.CreatedAt,
	}
}

// ToTicketListAPI converts a slice of domain Tickets to a slice of API Tickets
func ToTicketListAPI(tickets []*Ticket) []api.Ticket {
	res := make([]api.Ticket, len(tickets))
	for i, t := range tickets {
		res[i] = ToTicketAPI(t)
	}
	return res
}

// ToWarrantyAPI converts a domain Warranty to the API Warranty representation
func ToWarrantyAPI(w *Warranty) api.Warranty {
	idUUID, _ := uuid.Parse(w.ID)
	ticketUUID, _ := uuid.Parse(w.TicketID)

	return api.Warranty{
		Id:           idUUID,
		TicketId:     ticketUUID,
		TicketNumber: w.TicketNumber,
		CustomerName: w.CustomerName,
		DeviceInfo:   w.DeviceInfo,
		StartDate:    openapi_types.Date{Time: w.StartDate},
		EndDate:      openapi_types.Date{Time: w.EndDate},
		Status:       api.WarrantyStatus(w.Status),
	}
}

// ToWarrantyListAPI converts a slice of domain Warranties to a slice of API Warranties
func ToWarrantyListAPI(warranties []*Warranty) []api.Warranty {
	res := make([]api.Warranty, len(warranties))
	for i, w := range warranties {
		res[i] = ToWarrantyAPI(w)
	}
	return res
}

func maskName(name string) string {
	words := strings.Split(name, " ")
	for i, word := range words {
		if word == "" {
			continue
		}
		runes := []rune(word)
		length := len(runes)
		if length <= 4 {
			for j := 1; j < length; j++ {
				runes[j] = '*'
			}
		} else {
			for j := 3; j < length; j++ {
				runes[j] = '*'
			}
		}
		words[i] = string(runes)
	}
	return strings.Join(words, " ")
}

func maskPhone(phone string) string {
	runes := []rune(phone)
	length := len(runes)
	if length <= 3 {
		return phone
	}
	for i := 0; i < length-3; i++ {
		runes[i] = '*'
	}
	return string(runes)
}

// ToPublicTrackerTicketAPI converts a domain Ticket to the Public Tracker API representation with masked PII and omitted financial info
func ToPublicTrackerTicketAPI(t *Ticket) api.PublicTrackerTicket {
	repairVal := t.RepairAction

	idUUID, _ := uuid.Parse(t.ID)

	return api.PublicTrackerTicket{
		Id:                   idUUID,
		CustomerNameMasked:   maskName(t.CustomerName),
		CustomerPhoneMasked:  maskPhone(t.CustomerPhone),
		BrandPhone:           t.BrandPhone,
		ModelPhone:           t.ModelPhone,
		DamageDescription:    t.DamageDescription,
		RepairAction:         &repairVal,
		Status:               api.PublicTrackerTicketStatus(t.Status),
		WarrantyDurationDays: t.WarrantyDurationDays,
		PickedUpAt:           t.PickedUpAt,
		CreatedAt:            t.CreatedAt,
	}
}
