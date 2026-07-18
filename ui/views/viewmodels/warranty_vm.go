package viewmodels

import (
	"fmt"

	"github.com/denden-dr/OpenBench/internal/models"
)

type ClaimVM struct {
	ClaimID            string
	ClaimNumber        string
	WarrantyID         string
	WarrantyStatus     models.WarrantyStatus
	TicketNumber       string
	CustomerName       string
	DeviceSummary      string
	EvaluationStatus   models.ClaimEvaluationStatus
	StatusBadgeClass   string
	CreatedAtFormatted string
	IssueDescription   string
}

func NewClaimVM(c models.ClaimSummary) ClaimVM {
	return ClaimVM{
		ClaimID:            c.ClaimID,
		ClaimNumber:        c.ClaimNumber,
		WarrantyID:         c.WarrantyID,
		WarrantyStatus:     c.WarrantyStatus,
		TicketNumber:       c.TicketNumber,
		CustomerName:       c.CustomerName,
		DeviceSummary:      fmt.Sprintf("%s %s", c.DeviceBrand, c.DeviceModel),
		EvaluationStatus:   c.EvaluationStatus,
		StatusBadgeClass:   getEvaluationBadgeClass(c.EvaluationStatus),
		CreatedAtFormatted: c.CreatedAt.Format("02 Jan 2006"),
		IssueDescription:   c.IssueDescription,
	}
}

func NewClaimVMs(claims []models.ClaimSummary) []ClaimVM {
	vms := make([]ClaimVM, len(claims))
	for i, c := range claims {
		vms[i] = NewClaimVM(c)
	}
	return vms
}

func getEvaluationBadgeClass(status models.ClaimEvaluationStatus) string {
	baseClass := "inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium border"
	switch status {
	case models.ClaimEvaluationPending:
		return baseClass + " bg-orange-100 text-orange-800 border-orange-200"
	case models.ClaimEvaluationAccepted:
		return baseClass + " bg-emerald-100 text-emerald-800 border-emerald-200"
	case models.ClaimEvaluationRejected, models.ClaimEvaluationVoid:
		return baseClass + " bg-red-100 text-red-800 border-red-200"
	default:
		return baseClass + " bg-slate-100 text-slate-800 border-slate-200"
	}
}
