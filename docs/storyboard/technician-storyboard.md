# Technician Storyboard

This document focuses on the workflows from the perspective of the Technician performing diagnoses and repairs.

## Staff Stories

### Story 9: Technician Daily Workflow

> A technician's daily routine in the system.

1. Technician logs in
2. Technician views **unassigned ticket queue** (sorted by created_at)
3. Technician claims a ticket → assigned to them `[diagnosing]`
4. Technician opens ticket, reviews customer info and issue description
5. Technician uploads "before" photos for documentation
6. Technician performs diagnosis, adds internal technical notes
7. Technician sets **estimated completion time**
8. Technician submits diagnosis with cost breakdown
9. _(Waits for customer confirmation)_
10. Technician logs parts used from inventory (POS-style form)
11. Technician performs repair, updating progress notes
12. Technician uploads "after" photos
13. Technician marks ticket as ready → gets **accessory return reminder**
14. Technician returns SIM, case, SD card etc. from checklist

---

### Technician Roles in Customer Flows

While Story 9 outlines the daily workflow, technicians are also directly involved in managing edge cases during repair:

#### Handling Unrepairable Devices (from Story 2)
1. Technician uploads "before" photos, performs diagnosis
2. Technician determines device is **unrepairable**
3. Technician submits diagnosis as unrepairable with notes

#### Discovering Critical Problems Mid-Repair (from Story 4 & 5)
1. Technician discovers a **critical/additional problem** during repair
2. Technician pauses repair, updates diagnosis with new findings `[diagnosing]` _(re-diagnosis)_
3. Technician submits updated cost estimate `[waiting_customer_confirm]`
4. If customer confirms, technician resumes and completes repair `[repairing]`

#### Handling Out-of-Stock Parts (from Story 6)
1. Technician identifies parts needed but **not in stock**
2. Technician updates status `[waiting_parts]`
3. Once Admin restocks parts, technician resumes, logs parts used `[repairing]`
