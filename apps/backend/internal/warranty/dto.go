package warranty

import (
	"time"

	"github.com/denden-dr/OpenBench/apps/backend/internal/models"
)

type WarrantyResponse struct {
	ID        string                `json:"id"`
	TicketID  string                `json:"ticket_id"`
	StartDate time.Time             `json:"start_date"`
	EndDate   time.Time             `json:"end_date"`
	Status    models.WarrantyStatus `json:"status"`
	Notes     *string               `json:"notes,omitempty"`
}

type CreateClaimRequest struct {
	WarrantyID       string `json:"warranty_id" validate:"required"`
	IssueDescription string `json:"issue_description" validate:"required"`
}

type ClaimResponse struct {
	ClaimID          string                       `json:"claim_id"`
	ClaimNumber      string                       `json:"claim_number"`
	WarrantyID       string                       `json:"warranty_id"`
	Status           models.ServiceTicketStatus   `json:"status"`
	EvaluationStatus models.ClaimEvaluationStatus `json:"evaluation_status"`
	IssueDescription string                       `json:"issue_description"`
	RepairAction     *string                      `json:"repair_action"`
	Notes            *string                      `json:"notes,omitempty"`
	EvaluationNotes  *string                      `json:"evaluation_notes,omitempty"`
	CreatedAt        time.Time                    `json:"created_at"`
	UpdatedAt        time.Time                    `json:"updated_at"`
}

type ClaimListResponse struct {
	ClaimID          string                       `json:"claim_id"`
	ClaimNumber      string                       `json:"claim_number"`
	WarrantyID       string                       `json:"warranty_id"`
	Status           models.ServiceTicketStatus   `json:"status"`
	EvaluationStatus models.ClaimEvaluationStatus `json:"evaluation_status"`
	IssueDescription string                       `json:"issue_description"`
	CreatedAt        time.Time                    `json:"created_at"`
}

type ChangeClaimStatusRequest struct {
	Status models.ServiceTicketStatus `json:"status" validate:"required,oneof=RECEIVED REPAIRING PENDING_CONFIRMATION FIXED COMPLETED CANCELLED RETURNED"`
}

type ClaimStatusResponse struct {
	ClaimID   string                     `json:"claim_id"`
	Status    models.ServiceTicketStatus `json:"status"`
	UpdatedAt time.Time                  `json:"updated_at"`
}

type UpdateClaimRequest struct {
	IssueDescription string `json:"issue_description" validate:"required"`
	RepairAction     string `json:"repair_action"`
	Notes            string `json:"notes"`
}

type EvaluateClaimRequest struct {
	Status models.ClaimEvaluationStatus `json:"status" validate:"required,oneof=PENDING ACCEPTED REJECTED VOID"`
	Notes  string                       `json:"notes"`
}

type EvaluateClaimResponse struct {
	ClaimID          string                       `json:"claim_id"`
	EvaluationStatus models.ClaimEvaluationStatus `json:"evaluation_status"`
	UpdatedAt        time.Time                    `json:"updated_at"`
}

type UpdateWarrantyStatusRequest struct {
	Status models.WarrantyStatus `json:"status" validate:"required,oneof=ACTIVE EXPIRED VOID"`
	Notes  string                `json:"notes"`
}

type UpdateWarrantyStatusResponse struct {
	WarrantyID string                `json:"warranty_id"`
	Status     models.WarrantyStatus `json:"status"`
	UpdatedAt  time.Time             `json:"updated_at"`
}

func MapToWarrantyResponse(w *models.Warranty) WarrantyResponse {
	return WarrantyResponse{
		ID:        w.ID,
		TicketID:  w.TicketID,
		StartDate: w.StartDate,
		EndDate:   w.EndDate,
		Status:    w.Status,
		Notes:     w.Notes,
	}
}

func MapToClaimResponse(c *models.Claim) ClaimResponse {
	return ClaimResponse{
		ClaimID:          c.ID,
		ClaimNumber:      c.ClaimNumber,
		WarrantyID:       c.WarrantyID,
		Status:           c.Status,
		EvaluationStatus: c.EvaluationStatus,
		IssueDescription: c.IssueDescription,
		RepairAction:     c.RepairAction,
		Notes:            c.Notes,
		EvaluationNotes:  c.EvaluationNotes,
		CreatedAt:        c.CreatedAt,
		UpdatedAt:        c.UpdatedAt,
	}
}

func MapToClaimListResponse(c models.Claim) ClaimListResponse {
	return ClaimListResponse{
		ClaimID:          c.ID,
		ClaimNumber:      c.ClaimNumber,
		WarrantyID:       c.WarrantyID,
		Status:           c.Status,
		EvaluationStatus: c.EvaluationStatus,
		IssueDescription: c.IssueDescription,
		CreatedAt:        c.CreatedAt,
	}
}
