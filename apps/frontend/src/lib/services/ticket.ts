import { env } from '$env/dynamic/public';
import { apiFetch as fetch } from './api';

const getApiUrl = () => {
  try {
    return env.PUBLIC_API_URL || '';
  } catch {
    return '';
  }
};

export interface Ticket {
  id: string;
  ticket_number: string;
  customer_name: string;
  customer_phone: string;
  brand_phone: string;
  model_phone: string;
  serial_number: string;
  damage_description: string;
  repair_action: string;
  cost: number;
  status: 'received' | 'in_repair' | 'completed' | 'cancelled';
  ui_status?: 'received' | 'in_repair' | 'ready_for_pickup' | 'completed' | 'cancelled';
  device_position: 'warehouse' | 'picked_up';
  payment_status: 'none' | 'requesting' | 'paid';
  payment_method?: 'cash' | 'qris';
  warranty_duration_days: number;
  picked_up_at?: string;
  warranty_expiry_date?: string;
  created_at: string;
}

export interface PublicTrackerTicket {
  id: string;
  ticket_number: string;
  customer_name_masked: string;
  customer_phone_masked: string;
  brand_phone: string;
  model_phone: string;
  damage_description: string;
  repair_action?: string;
  status: 'received' | 'in_repair' | 'completed' | 'cancelled';
  ui_status?: 'received' | 'in_repair' | 'ready_for_pickup' | 'completed' | 'cancelled';
  warranty_duration_days: number;
  picked_up_at?: string;
  warranty_expiry_date?: string;
  created_at: string;
}

export interface TicketCreate {
  customer_name: string;
  customer_phone: string;
  brand_phone: string;
  model_phone: string;
  serial_number?: string;
  damage_description: string;
  repair_action?: string;
  cost?: number;
  warranty_duration_days: number;
}

function mapTicket(t: any): Ticket {
  if (!t) return t;
  let warranty_expiry_date = t.warranty_expiry_date;
  if (!warranty_expiry_date && t.picked_up_at && t.warranty_duration_days) {
    const d = new Date(t.picked_up_at);
    d.setDate(d.getDate() + t.warranty_duration_days);
    warranty_expiry_date = d.toISOString().split('T')[0];
  }
  return {
    ...t,
    ui_status: t.status,
    warranty_expiry_date
  };
}

function mapPublicTrackerTicket(t: any): PublicTrackerTicket {
  if (!t) return t;
  let warranty_expiry_date = t.warranty_expiry_date;
  if (!warranty_expiry_date && t.picked_up_at && t.warranty_duration_days) {
    const d = new Date(t.picked_up_at);
    d.setDate(d.getDate() + t.warranty_duration_days);
    warranty_expiry_date = d.toISOString().split('T')[0];
  }
  return {
    ...t,
    ui_status: t.status,
    warranty_expiry_date
  };
}

export const ticketService = {
  async getTickets(): Promise<Ticket[]> {
    const res = await fetch(`${getApiUrl()}/api/v1/admin/tickets`, { credentials: 'include' });
    const body = await res.json();
    if (!res.ok) throw new Error(body.message || 'Failed to fetch tickets');
    return (body.data ?? []).map(mapTicket);
  },

  async getTicket(id: string): Promise<Ticket> {
    const res = await fetch(`${getApiUrl()}/api/v1/admin/tickets/${id}`, { credentials: 'include' });
    const body = await res.json();
    if (!res.ok) throw new Error(body.message || 'Failed to fetch ticket');
    return mapTicket(body.data);
  },

  async getPublicTrackerTicket(ticketNumber: string, signal?: AbortSignal): Promise<PublicTrackerTicket> {
    const res = await fetch(`${getApiUrl()}/api/v1/tracker/${ticketNumber}`, { signal });
    const body = await res.json();
    if (!res.ok) throw new Error(body.message || 'Ticket not found');
    return mapPublicTrackerTicket(body.data);
  },

  async createTicket(ticket: TicketCreate): Promise<Ticket> {
    const res = await fetch(`${getApiUrl()}/api/v1/admin/tickets`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify(ticket)
    });
    const body = await res.json();
    if (!res.ok) throw new Error(body.message || 'Failed to create ticket');
    return mapTicket(body.data);
  },

  async updateTicket(id: string, updates: Partial<Ticket>): Promise<Ticket> {
    const payload = { ...updates };
    const res = await fetch(`${getApiUrl()}/api/v1/admin/tickets/${id}`, {
      method: 'PATCH',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify(payload)
    });
    const body = await res.json();
    if (!res.ok) throw new Error(body.message || 'Failed to update ticket');
    return mapTicket(body.data);
  },

  async emergencyUpdateTicket(id: string, updates: Partial<Ticket>): Promise<Ticket> {
    const payload = { ...updates };
    const res = await fetch(`${getApiUrl()}/api/v1/admin/tickets/${id}/emergency`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify(payload)
    });
    const body = await res.json();
    if (!res.ok) throw new Error(body.message || 'Failed to perform emergency update on ticket');
    return mapTicket(body.data);
  },

  async ListWarranties(): Promise<any[]> {
    const res = await fetch(`${getApiUrl()}/api/v1/admin/warranties`, { credentials: 'include' });
    const body = await res.json();
    if (!res.ok) throw new Error(body.message || 'Failed to fetch warranties');
    return body.data ?? [];
  },

  async getMyTickets(): Promise<Ticket[]> {
    const res = await fetch(`${getApiUrl()}/api/v1/tickets/my`, { credentials: 'include' });
    const body = await res.json();
    if (!res.ok) throw new Error(body.message || 'Failed to fetch my tickets');
    return (body.data ?? []).map(mapTicket);
  },

  async createRepairRequest(request: {
    customer_name: string;
    customer_phone: string;
    brand_phone: string;
    model_phone: string;
    serial_number?: string;
    damage_description: string;
  }): Promise<{ id: string; ticket_number: string }> {
    const res = await fetch(`${getApiUrl()}/api/v1/tickets/request`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(request)
    });
    const body = await res.json();
    if (!res.ok) throw new Error(body.message || 'Failed to submit repair request');
    return body.data;
  }
};
