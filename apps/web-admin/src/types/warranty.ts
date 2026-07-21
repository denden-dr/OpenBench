export type WarrantyStatus = 'ACTIVE' | 'VOID' | 'EXPIRED'

export type ClaimEvaluationStatus = 'PENDING' | 'ACCEPTED' | 'REJECTED' | 'VOID'

export interface Warranty {
  id: string
  ticket_id: string
  start_date: string
  end_date: string
  status: WarrantyStatus
  notes: string | null
}

export interface WarrantyClaim {
  claim_id: string
  claim_number: string
  warranty_id: string
  warranty_ticket_ref_id?: string | null
  evaluation_status: ClaimEvaluationStatus
  issue_description: string
  notes: string | null
  evaluation_notes: string | null
  created_at: string
  updated_at: string
}
