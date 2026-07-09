package models

import "time"

type WarrantyStatus string

const (
	WarrantyStatusActive  WarrantyStatus = "ACTIVE"
	WarrantyStatusExpired WarrantyStatus = "EXPIRED"
	WarrantyStatusVoid    WarrantyStatus = "VOID"
)

type ClaimEvaluationStatus string

const (
	ClaimEvaluationPending  ClaimEvaluationStatus = "PENDING"
	ClaimEvaluationAccepted ClaimEvaluationStatus = "ACCEPTED"
	ClaimEvaluationRejected ClaimEvaluationStatus = "REJECTED"
	ClaimEvaluationVoid     ClaimEvaluationStatus = "VOID"
)

type Warranty struct {
	ID        string         `json:"id"`
	TicketID  string         `json:"ticket_id"`
	StartDate time.Time      `json:"start_date"`
	EndDate   time.Time      `json:"end_date"`
	Status    WarrantyStatus `json:"status"`
	Notes     *string        `json:"notes,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

type Claim struct {
	ID               string                `json:"claim_id"`
	ClaimNumber      string                `json:"claim_number"`
	WarrantyID       string                `json:"warranty_id"`
	Status           ServiceTicketStatus   `json:"status"`
	EvaluationStatus ClaimEvaluationStatus `json:"evaluation_status"`
	IssueDescription string                `json:"issue_description"`
	RepairAction     *string               `json:"repair_action"`
	Notes            *string               `json:"notes,omitempty"`
	EvaluationNotes  *string               `json:"evaluation_notes,omitempty"`
	CreatedAt        time.Time             `json:"created_at"`
	UpdatedAt        time.Time             `json:"updated_at"`
}
