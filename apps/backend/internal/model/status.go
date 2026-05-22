// apps/backend/internal/model/status.go
package model

const (
	StatusReceived        = "received"
	StatusDiagnosing      = "diagnosing"
	StatusPendingApproval = "pending_approval"
	StatusRepairing       = "repairing"
	StatusReady           = "ready"
	StatusCancelled       = "cancelled"
	StatusPickedUp        = "picked_up"
)

// ValidTransitions maps each status to its allowed next statuses.
var ValidTransitions = map[string][]string{
	StatusReceived:        {StatusDiagnosing},
	StatusDiagnosing:      {StatusPendingApproval},
	StatusPendingApproval: {StatusRepairing, StatusCancelled},
	StatusRepairing:       {StatusReady},
	StatusReady:           {StatusPickedUp},
	StatusCancelled:       {StatusPickedUp},
}
