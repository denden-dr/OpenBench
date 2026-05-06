# User Storyboard

This document focuses on the workflows from the perspective of the Customer and Public Users.

## Customer Stories

### Story 1: Happy Path — Book, Diagnose, Repair, Pay, Done

> Customer books a device, technician diagnoses and repairs it, customer pays and picks up.

1. Customer opens landing page
2. Customer navigates to booking form
3. Customer fills booking: device type, brand, model, issue, access info, accessories
4. Customer acknowledges diagnosis fee (mandatory checkbox)
5. Customer submits booking → receives **Ticket ID + QR Code (PDF receipt)** `[received]`
6. Customer opens ticket details via Ticket ID or QR Code
7. Customer sees ticket status: **received**, waiting for technician
8. Technician claims ticket from queue `[diagnosing]`
9. Technician uploads "before" photos
10. Technician performs diagnosis, adds notes, sets estimated completion time
11. Technician submits diagnosis result & fee breakdown `[waiting_customer_confirm]`
12. Customer receives notification: diagnosis complete
13. Customer views diagnosis details: findings, fee breakdown (diagnosis + labor + parts)
14. Customer confirms to proceed with repair
15. Technician logs parts used (inventory deducted) `[repairing]`
16. Technician completes repair, uploads "after" photos
17. Technician marks device ready, gets accessory return reminder `[ready]`
18. Customer receives notification: device ready for pickup
19. Customer views before & after photos, warranty info
20. Customer pays (cash or online via Midtrans) `[completed]`
21. System calculates warranty expiry based on part grades
22. Customer downloads invoice

---

### Story 2: Unrepairable — Technician Cancels After Diagnosis

> Technician determines device cannot be repaired. **No charge to customer.**

1. Customer submits booking `[received]`
2. Technician claims ticket `[diagnosing]`
3. Technician uploads "before" photos, performs diagnosis
4. Technician determines device is **unrepairable**
5. Technician submits diagnosis as unrepairable with notes
6. Admin cancels ticket with reason logged `[cancelled]`
7. Customer receives notification: repair not possible
8. Customer views diagnosis details and "before" photos
9. Customer picks up device — **no payment required** (diagnosis fee waived for unrepairable devices)
10. Accessories returned per checklist
11. Customer downloads invoice (zero-amount / record only)

---

### Story 3: Customer Cancels After Diagnosis

> Customer decides not to proceed after seeing the diagnosis cost.

1. Customer submits booking `[received]`
2. Technician claims ticket `[diagnosing]`
3. Technician performs diagnosis, submits results `[waiting_customer_confirm]`
4. Customer views diagnosis details and cost breakdown
5. Customer **declines** to proceed
6. Admin cancels ticket `[cancelled]`
7. Customer pays **diagnosis fee only**
8. Customer picks up device
9. Customer downloads invoice

---

### Story 4: Critical Problem Found During Repair → Customer Cancels

> Technician discovers additional damage mid-repair, customer decides to stop.

1. Customer submits booking `[received]`
2. Technician claims and diagnoses `[diagnosing]`
3. Customer confirms to proceed `[repairing]`
4. Technician discovers a **critical/additional problem** during repair
5. Technician pauses repair, updates diagnosis with new findings `[diagnosing]` _(re-diagnosis)_
6. Technician submits updated cost estimate `[waiting_customer_confirm]`
7. Customer views updated diagnosis and new cost
8. Customer **declines** to continue
9. Admin cancels ticket `[cancelled]`
10. Customer pays diagnosis fee + partial labor (if applicable)
11. Customer picks up device
12. Customer downloads invoice

---

### Story 5: Critical Problem Found → Customer Confirms (Repeatable)

> Same as Story 4, but the customer agrees to continue. This cycle can repeat.

1. Customer submits booking `[received]`
2. Technician claims and diagnoses `[diagnosing]`
3. Customer confirms to proceed `[repairing]`
4. Technician discovers a **critical/additional problem**
5. Technician pauses, re-diagnoses `[diagnosing]`
6. Technician submits updated estimate `[waiting_customer_confirm]`
7. Customer confirms to continue _(this cycle can repeat from step 3)_
8. Technician resumes and completes repair `[repairing]`
9. Technician uploads "after" photos, marks ready `[ready]`
10. Customer pays full amount `[completed]`
11. Customer downloads invoice with warranty info

---

### Story 6: Waiting for Parts

> Required parts are out of stock; repair is delayed.

1. Customer submits booking `[received]`
2. Technician claims and diagnoses `[diagnosing]`
3. Technician identifies parts needed but **not in stock**
4. Technician updates status `[waiting_parts]`
5. Customer receives notification: waiting for parts, updated ETA
6. Admin restocks parts (inventory updated)
7. Technician resumes, logs parts used `[repairing]`
8. Continues to standard completion flow (Story 1, steps 16–22)

---

## Public Stories

### Story 7: Guest Tracking (No Account)

> A guest tracks their ticket without logging in.

1. Guest opens landing page
2. Guest navigates to "Track Repair" section
3. Guest enters **Ticket ID + Phone Number**
4. System returns sanitized ticket info: status, device brand/model, estimated ready date
5. Guest cannot see internal notes, photos, or payment details

---

### Story 7a: Guest Fills Booking → Deferred Sign-Up → Auto-Submit

> A guest starts the booking form without an account. On confirm, they're prompted to sign up. After signing up, their form data is preserved and the booking is submitted automatically.

1. Guest opens landing page
2. Guest navigates to booking form _(accessible without login)_
3. Guest fills booking form: device type, brand, model, issue, access info, accessories
4. Guest acknowledges diagnosis fee (mandatory checkbox)
5. Guest clicks **"Confirm Booking"**
6. System detects guest is **not authenticated**
7. System saves form data to **localStorage** before redirect
8. System redirects guest to sign-up page with a banner: _"Create an account to submit your booking"_
9. Sign-up page shows **social login (Google) as recommended** option, with email+password as alternative
10. Guest signs up via social login or email+password
11. Supabase Auth creates account, returns session
12. System detects **pending booking data** in localStorage
13. System redirects to booking form — **pre-filled with saved data**
14. System shows confirmation summary: _"Review your booking details"_
15. Customer _(now authenticated)_ confirms → booking submitted automatically `[received]`
16. Customer receives **Ticket ID + QR Code (PDF receipt)**
17. System clears pending booking data from localStorage

#### Edge Cases
- **Guest already has an account** → step 9 shows login instead, same flow from step 11
- **localStorage cleared / different device** → form is empty after auth, customer re-fills (graceful degradation)
- **Auth fails / user cancels** → guest returns to booking form with data still intact (localStorage persists)
- **Session expires during form fill** → same deferred auth flow triggers again

---

### Story 8: Public Status Board

> Anyone views the live anonymized repair queue.

1. Visitor opens Public Status Board page
2. System displays Kanban-style board with active tickets
3. Each card shows: **masked Ticket ID**, brand, model, current status
4. No customer names, phone numbers, or identifying info displayed
5. Board updates in real-time

---

## Edge Case Stories

### Story 11: Online Payment Failure & Retry

> Customer's online payment fails.

1. Customer selects online payment (Midtrans)
2. System creates payment intent `[payment: pending]`
3. Payment fails (insufficient funds, timeout, etc.) `[payment: failed]`
4. Customer receives notification: payment failed
5. Customer retries payment
6. Payment succeeds via webhook `[payment: completed]`
7. Ticket moves to `[completed]`
