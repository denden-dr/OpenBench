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
	ID        string         `json:"id" db:"id"`
	TicketID  string         `json:"ticket_id" db:"ticket_id"`
	StartDate time.Time      `json:"start_date" db:"start_date"`
	EndDate   time.Time      `json:"end_date" db:"end_date"`
	Status    WarrantyStatus `json:"status" db:"status"`
	Notes     *string        `json:"notes,omitempty" db:"notes"`
	CreatedAt time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt time.Time      `json:"updated_at" db:"updated_at"`
}

type Claim struct {
	ID               string                `json:"claim_id" db:"id"`
	ClaimNumber      string                `json:"claim_number" db:"claim_number"`
	WarrantyID       string                `json:"warranty_id" db:"warranty_id"`
	Status           ServiceTicketStatus   `json:"status" db:"status"`
	EvaluationStatus ClaimEvaluationStatus `json:"evaluation_status" db:"evaluation_status"`
	IssueDescription string                `json:"issue_description" db:"issue_description"`
	RepairAction     *string               `json:"repair_action" db:"repair_action"`
	Notes            *string               `json:"notes,omitempty" db:"notes"`
	EvaluationNotes  *string               `json:"evaluation_notes,omitempty" db:"evaluation_notes"`
	CreatedAt        time.Time             `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time             `json:"updated_at" db:"updated_at"`
}
