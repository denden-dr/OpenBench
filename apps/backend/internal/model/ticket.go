package model

import (
    "time"
)

type Ticket struct {
    ID               string    `db:"id" json:"id"`
    DeviceType       string    `db:"device_type" json:"device_type"`
    Brand            string    `db:"brand" json:"brand"`
    Model            string    `db:"model" json:"model"`
    IssueDescription string    `db:"issue_description" json:"issue_description"`
    Status           string    `db:"status" json:"status"`
    DiagnosisFee     float64   `db:"diagnosis_fee" json:"diagnosis_fee"`
    CreatedAt        time.Time `db:"created_at" json:"created_at"`
    UpdatedAt        time.Time `db:"updated_at" json:"updated_at"`
}
