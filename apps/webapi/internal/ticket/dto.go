package ticket

import (
	"time"

	"github.com/denden-dr/OpenBench/internal/models"
)

type CreateTicketRequest struct {
	CustomerName     string `json:"customer_name" validate:"required"`
	CustomerPhone    string `json:"customer_phone" validate:"required"`
	DeviceBrand      string `json:"device_brand" validate:"required"`
	DeviceModel      string `json:"device_model" validate:"required"`
	DevicePasscode   string `json:"device_passcode"`
	IssueDescription string `json:"issue_description" validate:"required"`
	RepairAction     string `json:"repair_action"`
	Cost             int64  `json:"cost" validate:"min=0"`
	WarrantyDays     int    `json:"warranty_days" validate:"min=0"`
}

type UpdateTicketRequest struct {
	CustomerName     string `json:"customer_name" validate:"required"`
	CustomerPhone    string `json:"customer_phone" validate:"required"`
	IssueDescription string `json:"issue_description" validate:"required"`
	RepairAction     string `json:"repair_action"`
	Cost             int64  `json:"cost" validate:"min=0"`
	WarrantyDays     int    `json:"warranty_days" validate:"min=0"`
	Notes            string `json:"notes"`
}

type ChangeStatusRequest struct {
	Status models.ServiceTicketStatus `json:"status" validate:"required,oneof=RECEIVED REPAIRING PENDING_CONFIRMATION FIXED COMPLETED CANCELLED RETURNED"`
}

type EmergencyUpdateTicketRequest struct {
	CustomerName     string                     `json:"customer_name" validate:"required"`
	CustomerPhone    string                     `json:"customer_phone" validate:"required"`
	DeviceBrand      string                     `json:"device_brand" validate:"required"`
	DeviceModel      string                     `json:"device_model" validate:"required"`
	DevicePasscode   string                     `json:"device_passcode"`
	Status           models.ServiceTicketStatus `json:"status" validate:"required,oneof=RECEIVED REPAIRING PENDING_CONFIRMATION FIXED COMPLETED CANCELLED RETURNED"`
	IssueDescription string                     `json:"issue_description" validate:"required"`
	RepairAction     string                     `json:"repair_action"`
	Cost             int64                      `json:"cost" validate:"min=0"`
	WarrantyDays     int                        `json:"warranty_days" validate:"min=0"`
	Notes            string                     `json:"notes"`
}

type TicketSearchRequest struct {
	Search    string `json:"search"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	ExactDate string `json:"exact_date"`
	IsActive  *bool  `json:"is_active"`
	Limit     int    `json:"limit"`
	Cursor    string `json:"cursor"`
}

type TicketResponse struct {
	TicketID         string                     `json:"ticket_id"`
	TicketNumber     string                     `json:"ticket_number,omitempty"`
	Status           models.ServiceTicketStatus `json:"status"`
	CustomerName     string                     `json:"customer_name"`
	CustomerPhone    string                     `json:"customer_phone,omitempty"`
	DeviceBrand      string                     `json:"device_brand"`
	DeviceModel      string                     `json:"device_model"`
	DevicePasscode   string                     `json:"device_passcode,omitempty"`
	IssueDescription string                     `json:"issue_description,omitempty"`
	RepairAction     *string                    `json:"repair_action,omitempty"`
	Cost             int64                      `json:"cost,omitempty"`
	WarrantyDays     int                        `json:"warranty_days,omitempty"`
	Notes            *string                    `json:"notes,omitempty"`
	CreatedAt        time.Time                  `json:"created_at"`
	UpdatedAt        time.Time                  `json:"updated_at,omitempty"`
}

type TicketSummaryResponse struct {
	TicketID     string                     `json:"ticket_id"`
	TicketNumber string                     `json:"ticket_number"`
	Status       models.ServiceTicketStatus `json:"status"`
	CustomerName string                     `json:"customer_name"`
	DeviceBrand  string                     `json:"device_brand"`
	DeviceModel  string                     `json:"device_model"`
	CreatedAt    time.Time                  `json:"created_at"`
}

type TicketStatusResponse struct {
	TicketID  string                     `json:"ticket_id"`
	Status    models.ServiceTicketStatus `json:"status"`
	UpdatedAt time.Time                  `json:"updated_at"`
}

func MapToTicketResponse(t *models.ServiceTicket) TicketResponse {
	return TicketResponse{
		TicketID:         t.ID,
		TicketNumber:     t.TicketNumber,
		Status:           t.Status,
		CustomerName:     t.CustomerName,
		CustomerPhone:    t.CustomerPhone,
		DeviceBrand:      t.DeviceBrand,
		DeviceModel:      t.DeviceModel,
		DevicePasscode:   t.DevicePasscode,
		IssueDescription: t.IssueDescription,
		RepairAction:     t.RepairAction,
		Cost:             t.Cost,
		WarrantyDays:     t.WarrantyDays,
		Notes:            t.Notes,
		CreatedAt:        t.CreatedAt,
		UpdatedAt:        t.UpdatedAt,
	}
}

func MapToTicketSummaryResponse(t models.ServiceTicket) TicketSummaryResponse {
	return TicketSummaryResponse{
		TicketID:     t.ID,
		TicketNumber: t.TicketNumber,
		Status:       t.Status,
		CustomerName: t.CustomerName,
		DeviceBrand:  t.DeviceBrand,
		DeviceModel:  t.DeviceModel,
		CreatedAt:    t.CreatedAt,
	}
}
