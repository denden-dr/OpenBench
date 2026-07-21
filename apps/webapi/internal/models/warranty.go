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
	ID                  string                `json:"claim_id" db:"id"`
	ClaimNumber         string                `json:"claim_number" db:"claim_number"`
	WarrantyID          string                `json:"warranty_id" db:"warranty_id"`
	WarrantyTicketRefID *string               `json:"warranty_ticket_ref_id,omitempty" db:"warranty_ticket_ref_id"`
	EvaluationStatus    ClaimEvaluationStatus `json:"evaluation_status" db:"evaluation_status"`
	IssueDescription    string                `json:"issue_description" db:"issue_description"`
	Notes               *string               `json:"notes,omitempty" db:"notes"`
	EvaluationNotes     *string               `json:"evaluation_notes,omitempty" db:"evaluation_notes"`
	CreatedAt           time.Time             `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time             `json:"updated_at" db:"updated_at"`
}

type ClaimSummary struct {
	ClaimID             string                `json:"claim_id" db:"claim_id"`
	ClaimNumber         string                `json:"claim_number" db:"claim_number"`
	WarrantyID          string                `json:"warranty_id" db:"warranty_id"`
	WarrantyTicketRefID *string               `json:"warranty_ticket_ref_id,omitempty" db:"warranty_ticket_ref_id"`
	WarrantyStatus      WarrantyStatus        `json:"warranty_status" db:"warranty_status"`
	TicketID            string                `json:"ticket_id" db:"ticket_id"`
	TicketNumber        string                `json:"ticket_number" db:"ticket_number"`
	CustomerName        string                `json:"customer_name" db:"customer_name"`
	DeviceBrand         string                `json:"device_brand" db:"device_brand"`
	DeviceModel         string                `json:"device_model" db:"device_model"`
	EvaluationStatus    ClaimEvaluationStatus `json:"evaluation_status" db:"evaluation_status"`
	IssueDescription    string                `json:"issue_description" db:"issue_description"`
	CreatedAt           time.Time             `json:"created_at" db:"created_at"`
}
