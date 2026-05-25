package dto

import "time"

type CreateWarrantyClaimRequest struct {
	TicketID              string `json:"ticket_id" validate:"required"`
	Issue                 string `json:"issue" validate:"required"`
	AdditionalDescription string `json:"additional_description"`
}

type VoidWarrantyClaimRequest struct {
	VoidReason string `json:"void_reason" validate:"required"`
}

type WarrantyClaimResponse struct {
	ID                    string          `json:"id"`
	TicketID              string          `json:"ticket_id"`
	ClaimTicketID         *string         `json:"claim_ticket_id"`
	Issue                 string          `json:"issue"`
	AdditionalDescription *string         `json:"additional_description"`
	Status                string          `json:"status"`
	VoidReason            *string         `json:"void_reason"`
	InspectedAt           *time.Time      `json:"inspected_at"`
	CreatedAt             time.Time       `json:"created_at"`
	UpdatedAt             time.Time       `json:"updated_at"`
	OriginalTicket        *TicketResponse `json:"originalTicket,omitempty"`
}

type ClaimCreationResult struct {
	Claim  WarrantyClaimResponse `json:"claim"`
	Ticket *TicketResponse       `json:"ticket,omitempty"`
}
