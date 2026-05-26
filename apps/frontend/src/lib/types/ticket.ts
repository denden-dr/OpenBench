export interface Ticket {
  id: string;
  customer_name: string;
  customer_phone?: string;
  customer_gender: string;
  brand: string;
  model: string;
  issue: string;
  additional_description?: string;
  accessories?: string;
  price: number;
  status: string;
  payment_status: string;
  warranty_days: number;
  entry_date: string;
  exit_date?: string;
  warranty_expiry_date?: string;
  is_warranty?: boolean;
  parent_ticket_id?: string;
}

export interface Claim {
  id: string;
  ticket_id: string;
  claim_ticket_id: string | null;
  issue: string;
  additional_description: string;
  status: 'waiting_inspection' | 'approved' | 'void';
  void_reason: string | null;
  inspected_at: string | null;
  created_at: string;
  originalTicket?: Ticket;
}
