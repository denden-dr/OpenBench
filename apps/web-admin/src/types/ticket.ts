export type TicketStatus = 
  | 'RECEIVED' 
  | 'REPAIRING' 
  | 'PENDING_CONFIRMATION' 
  | 'FIXED' 
  | 'COMPLETED' 
  | 'CANCELLED' 
  | 'RETURNED';

export interface TicketListItem {
  ticket_id: string;
  ticket_number: string;
  status: TicketStatus;
  customer_name: string;
  device_brand: string;
  device_model: string;
  created_at: string;
}

export interface TicketDetail {
  ticket_id: string;
  ticket_number?: string;
  status: TicketStatus;
  customer_name: string;
  customer_phone: string;
  device_brand: string;
  device_model: string;
  device_passcode?: string;
  issue_description: string;
  repair_action?: string;
  cost: number;
  warranty_days: number;
  notes?: string;
  created_at: string;
  updated_at?: string;
}

export interface CreateTicketRequest {
  customer_name: string;
  customer_phone: string;
  device_brand: string;
  device_model: string;
  device_passcode?: string;
  issue_description: string;
  repair_action: string;
  cost: number;
  warranty_days: number;
}

export interface UpdateTicketStatusRequest {
  status: TicketStatus;
}

export interface UpdateTicketDetailsRequest {
  customer_name?: string;
  customer_phone?: string;
  issue_description?: string;
  repair_action?: string;
  cost?: number;
  warranty_days?: number;
  notes?: string;
}

export interface PaginationMeta {
  limit: number;
  next_cursor?: string;
}

export interface TicketListResponse {
  data: TicketListItem[];
  meta?: PaginationMeta;
}

export interface TicketDetailResponse {
  data: TicketDetail;
}
