---
name: midtrans-webhook-integration
description: Use when implementing Midtrans payment gateway integration, handling payment webhooks, verifying webhook signatures, or processing async payment status updates in the Go Fiber backend
---

# Midtrans Webhook Integration

## Overview

Midtrans sends payment status updates via **asynchronous webhooks** to `POST /api/v1/webhooks/payment`. This endpoint must verify the signature, process idempotently, and update both `payments` and `tickets` tables.

## Webhook Flow

```
Customer → Midtrans (pays via QRIS/Bank Transfer)
Midtrans → POST /api/v1/webhooks/payment
  1. Parse notification JSON
  2. Verify signature (SHA-512)
  3. Check idempotency (skip if already processed)
  4. Update payment status
  5. If payment completed → update ticket status
  6. Return 200 OK to Midtrans
```

## Signature Verification

Midtrans signature: `SHA512(order_id + status_code + gross_amount + server_key)`

```go
func verifySignature(orderID, statusCode, grossAmount, serverKey, signature string) bool {
    raw := orderID + statusCode + grossAmount + serverKey
    hash := sha512.Sum512([]byte(raw))
    expected := hex.EncodeToString(hash[:])
    return hmac.Equal([]byte(expected), []byte(signature))
}
```

## Webhook Handler

```go
func (h *PaymentHandler) Webhook(c *fiber.Ctx) error {
    var notif MidtransNotification
    if err := c.BodyParser(&notif); err != nil {
        return c.Status(400).JSON(ErrorResponse{Code: "VALIDATION_ERROR"})
    }

    // 1. Verify signature
    if !verifySignature(notif.OrderID, notif.StatusCode, notif.GrossAmount, h.serverKey, notif.SignatureKey) {
        return c.Status(422).JSON(ErrorResponse{Code: "PAYMENT_VERIFICATION_FAILED"})
    }

    // 2. Map status
    status := mapMidtransStatus(notif.TransactionStatus, notif.FraudStatus)

    // 3. Idempotent upsert (skip if already completed)
    err := h.service.ProcessPaymentNotification(c.Context(), notif.OrderID, status)
    if err != nil {
        // Log error but still return 200 to Midtrans to prevent retries
        log.Error("webhook processing failed", "order_id", notif.OrderID, "error", err)
    }

    // 4. Always return 200 to Midtrans
    return c.SendStatus(200)
}
```

## Status Mapping

```go
func mapMidtransStatus(txStatus, fraudStatus string) string {
    switch txStatus {
    case "capture":
        if fraudStatus == "accept" { return "completed" }
        return "pending"
    case "settlement":
        return "completed"
    case "pending":
        return "pending"
    case "deny", "cancel", "expire":
        return "failed"
    default:
        return "pending"
    }
}
```

## Notification Struct

```go
type MidtransNotification struct {
    TransactionStatus string `json:"transaction_status"`
    OrderID           string `json:"order_id"`       // maps to payments.external_id
    StatusCode        string `json:"status_code"`
    GrossAmount       string `json:"gross_amount"`
    SignatureKey       string `json:"signature_key"`
    FraudStatus       string `json:"fraud_status"`
    PaymentType       string `json:"payment_type"`
}
```

## Idempotency Rules

- Use `external_id` (Midtrans order_id) as the idempotency key
- If payment already `completed` → skip update, return 200
- If payment already `failed` → skip update, return 200
- Only update `pending` → `completed` or `pending` → `failed`

## Service Layer Logic

```go
func (s *paymentService) ProcessPaymentNotification(ctx context.Context, orderID, status string) error {
    return s.db.WithTransaction(ctx, func(tx *sql.Tx) error {
        payment, err := s.repo.FindByExternalID(ctx, tx, orderID)
        if err != nil { return err }
        if payment.Status != "pending" { return nil } // idempotent

        payment.Status = status
        if err := s.repo.UpdateStatus(ctx, tx, payment); err != nil { return err }

        if status == "completed" {
            if err := s.ticketRepo.UpdateStatus(ctx, tx, payment.TicketID, "completed"); err != nil {
                return err
            }
            s.audit.Log(ctx, tx, "payment_completed", "payment", payment.ID, nil)
        }
        return nil
    })
}
```

## Configuration

```env
MIDTRANS_SERVER_KEY=your-server-key    # From Midtrans dashboard
MIDTRANS_IS_PRODUCTION=false           # true for production
```

## Testing Locally

- Use Midtrans **Sandbox** environment
- Use `ngrok` or similar to expose local webhook endpoint
- Midtrans Sandbox dashboard has a "Resend Notification" button

## Common Mistakes

| Mistake | Fix |
|---------|-----|
| Returning non-200 on processing error | Always return 200, Midtrans retries on non-200 |
| Not verifying signature | Always verify — prevents spoofed webhooks |
| Processing same notification twice | Check payment status before updating (idempotent) |
| Missing transaction around payment + ticket update | Wrap in DB transaction |
| Using client key instead of server key for signature | Server key only, never expose to client |
