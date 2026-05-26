package model

import "time"

type WarrantyClaimStatus string

const (
	ClaimWaitingInspection WarrantyClaimStatus = "waiting_inspection"
	ClaimApproved          WarrantyClaimStatus = "approved"
	ClaimVoid              WarrantyClaimStatus = "void"
)

type WarrantyClaim struct {
	ID                    string              `db:"id" json:"id"`
	TicketID              string              `db:"ticket_id" json:"ticket_id"`
	ClaimTicketID         *string             `db:"claim_ticket_id" json:"claim_ticket_id"`
	Issue                 string              `db:"issue" json:"issue"`
	AdditionalDescription *string             `db:"additional_description" json:"additional_description"`
	Status                WarrantyClaimStatus `db:"status" json:"status"`
	VoidReason            *string             `db:"void_reason" json:"void_reason"`
	InspectedAt           *time.Time          `db:"inspected_at" json:"inspected_at"`
	CreatedAt             time.Time           `db:"created_at" json:"created_at"`
	UpdatedAt             time.Time           `db:"updated_at" json:"updated_at"`
}
