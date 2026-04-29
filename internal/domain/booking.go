package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// BookingStatus represents the lifecycle state of a repair booking.
type BookingStatus string

const (
	StatusPendingDiagnosis  BookingStatus = "PENDING_DIAGNOSIS"
	StatusDiagnosisComplete BookingStatus = "DIAGNOSIS_COMPLETED"
	StatusApproved          BookingStatus = "APPROVED"
	StatusCanceled          BookingStatus = "CANCELED"
	StatusInProgress        BookingStatus = "IN_PROGRESS"
	StatusCompleted         BookingStatus = "COMPLETED"
)

// Booking represents a repair request from a user.
type Booking struct {
	ID                  uuid.UUID     `db:"id" json:"id"`
	UserID              uuid.UUID     `db:"user_id" json:"user_id"`
	DeviceName          string        `db:"device_name" json:"device_name"`
	IssueDescription    string        `db:"issue_description" json:"issue_description"`
	Status              BookingStatus `db:"status" json:"status"`
	DiagnosisResult     *string       `db:"diagnosis_result" json:"diagnosis_result,omitempty"`
	EstimatedCost       *float64      `db:"estimated_cost" json:"estimated_cost,omitempty"`
	TechnicianID        *uuid.UUID    `db:"technician_id" json:"technician_id,omitempty"`
	EstimatedRepairTime *string       `db:"estimated_repair_time" json:"estimated_repair_time,omitempty"`
	CreatedAt           time.Time     `db:"created_at" json:"created_at"`
	UpdatedAt           time.Time     `db:"updated_at" json:"updated_at"`
}

// Custom domain errors for the Booking and Technician context
var (
	ErrBookingNotFound        = errors.New("booking not found")
	ErrInvalidStateTransition = errors.New("invalid state transition")
	ErrNotFound               = errors.New("not found")
	ErrConflict               = errors.New("conflict")
	ErrInvalidInput           = errors.New("invalid input")
)
