import { describe, it, expect, vi, afterEach } from 'vitest';
import { ticketService, type Ticket, type PublicTrackerTicket } from './ticket';

describe('ticketService Unit Tests', () => {
  afterEach(() => {
    vi.restoreAllMocks();
  });

  const createMockResponse = (ok: boolean, data: any = null, message: string = '') => {
    return {
      ok,
      json: async () => ({ code: ok ? 200 : 400, message, data })
    };
  };

  it('should get all tickets', async () => {
    const mockTickets = [{ id: 'ticket-1', ticket_number: 'OB-001' }];
    const fetchSpy = vi.fn().mockResolvedValue(createMockResponse(true, mockTickets));
    vi.stubGlobal('fetch', fetchSpy);

    const result = await ticketService.getTickets();
    
    expect(fetchSpy).toHaveBeenCalledWith(expect.stringContaining('/api/v1/admin/tickets'), { credentials: 'include' });
    expect(result).toEqual([{ ...mockTickets[0], ui_status: undefined, warranty_expiry_date: undefined }]);
  });

  it('should handle getTickets error', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue(createMockResponse(false, null, 'Error fetching tickets')));
    await expect(ticketService.getTickets()).rejects.toThrow('Error fetching tickets');
  });

  it('should get a ticket by id', async () => {
    const mockTicket = { id: 'ticket-1', ticket_number: 'OB-001' };
    const fetchSpy = vi.fn().mockResolvedValue(createMockResponse(true, mockTicket));
    vi.stubGlobal('fetch', fetchSpy);

    const result = await ticketService.getTicket('ticket-1');
    
    expect(fetchSpy).toHaveBeenCalledWith(expect.stringContaining('/api/v1/admin/tickets/ticket-1'), { credentials: 'include' });
    expect(result).toEqual({ ...mockTicket, ui_status: undefined, warranty_expiry_date: undefined });
  });

  it('should get public tracker ticket', async () => {
    const mockTicket = { ticket_number: 'OB-202606-0001-A9X2', status: 'in_repair' };
    const fetchSpy = vi.fn().mockResolvedValue(createMockResponse(true, mockTicket));
    vi.stubGlobal('fetch', fetchSpy);

    const result = await ticketService.getPublicTrackerTicket('OB-202606-0001-A9X2');
    
    expect(fetchSpy).toHaveBeenCalledWith(expect.stringContaining('/api/v1/tracker/OB-202606-0001-A9X2'), { signal: undefined });
    expect(result).toEqual({ ...mockTicket, ui_status: 'in_repair', warranty_expiry_date: undefined });
  });

  it('should create a ticket', async () => {
    const newTicket = { customer_name: 'John Doe' } as any;
    const createdTicket = { ...newTicket, id: 'ticket-1', ticket_number: 'OB-001' };
    const fetchSpy = vi.fn().mockResolvedValue(createMockResponse(true, createdTicket));
    vi.stubGlobal('fetch', fetchSpy);

    const result = await ticketService.createTicket(newTicket);
    
    expect(fetchSpy).toHaveBeenCalledWith(expect.stringContaining('/api/v1/admin/tickets'), expect.objectContaining({
      method: 'POST',
      body: JSON.stringify(newTicket)
    }));
    expect(result).toEqual({ ...createdTicket, ui_status: undefined, warranty_expiry_date: undefined });
  });

  it('should update a ticket', async () => {
    const updates: Partial<Ticket> = { status: 'completed' };
    const updatedTicket = { id: 'ticket-1', status: 'completed' };
    const fetchSpy = vi.fn().mockResolvedValue(createMockResponse(true, updatedTicket));
    vi.stubGlobal('fetch', fetchSpy);

    const result = await ticketService.updateTicket('ticket-1', updates);
    
    expect(fetchSpy).toHaveBeenCalledWith(expect.stringContaining('/api/v1/admin/tickets/ticket-1'), expect.objectContaining({
      method: 'PATCH',
      body: JSON.stringify(updates)
    }));
    expect(result).toEqual({ ...updatedTicket, ui_status: 'completed', warranty_expiry_date: undefined });
  });
});
