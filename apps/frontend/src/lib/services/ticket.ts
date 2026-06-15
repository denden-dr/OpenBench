import { env } from '$env/dynamic/public';

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
  status: 'received' | 'diagnosing' | 'in_repair' | 'ready_for_pickup' | 'picked_up' | 'cancelled';
  device_position: 'warehouse' | 'picked_up';
  payment_status: 'none' | 'requesting' | 'paid';
  payment_method?: 'cash' | 'qris';
  warranty_expiry_date?: string;
  created_at: string;
}

export const ticketService = {
  async getTickets(): Promise<Ticket[]> {
    const res = await fetch(`${getApiUrl()}/api/v1/admin/tickets`, { credentials: 'include' });
    const body = await res.json();
    if (!res.ok) throw new Error(body.message || 'Failed to fetch tickets');
    return body.data;
  },

  async getTicket(id: string): Promise<Ticket> {
    const res = await fetch(`${getApiUrl()}/api/v1/admin/tickets/${id}`, { credentials: 'include' });
    const body = await res.json();
    if (!res.ok) throw new Error(body.message || 'Failed to fetch ticket');
    return body.data;
  },

  async getPublicTrackerTicket(id: string): Promise<Ticket> {
    const res = await fetch(`${getApiUrl()}/api/v1/tracker/${id}`);
    const body = await res.json();
    if (!res.ok) throw new Error(body.message || 'Ticket not found');
    return body.data;
  },

  async createTicket(ticket: Omit<Ticket, 'id' | 'ticket_number' | 'created_at'>): Promise<Ticket> {
    const res = await fetch(`${getApiUrl()}/api/v1/admin/tickets`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify(ticket)
    });
    const body = await res.json();
    if (!res.ok) throw new Error(body.message || 'Failed to create ticket');
    return body.data;
  },

  async updateTicket(id: string, updates: Partial<Ticket>): Promise<Ticket> {
    const res = await fetch(`${getApiUrl()}/api/v1/admin/tickets/${id}`, {
      method: 'PATCH',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify(updates)
    });
    const body = await res.json();
    if (!res.ok) throw new Error(body.message || 'Failed to update ticket');
    return body.data;
  }
};
