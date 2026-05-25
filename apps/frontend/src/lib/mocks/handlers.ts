import { json } from '@sveltejs/kit';
import { mockTickets, setMockTickets } from './mockData';

export async function handleMockRequest(event: any): Promise<Response | null> {
    const { request, url } = event;
    const path = url.pathname;
    const method = request.method;

    // Matches /api/v1/tickets
    if (path === '/api/v1/tickets') {
        if (method === 'GET') {
            return json({ success: true, data: mockTickets });
        }
        if (method === 'POST') {
            const data = await request.json();
            const newTicket = {
                id: 'TCK-' + Math.floor(Math.random() * 1000).toString().padStart(3, '0'),
                ...data,
                status: 'service_in',
                payment_status: 'unpaid',
                entry_date: new Date().toISOString()
            };
            setMockTickets([newTicket, ...mockTickets]);
            return json({ success: true, data: newTicket });
        }
    }

    // Matches /api/v1/tickets/[id]
    const ticketIdMatch = path.match(/^\/api\/v1\/tickets\/([^/]+)$/);
    if (ticketIdMatch) {
        const id = ticketIdMatch[1];
        if (method === 'PATCH') {
            const data = await request.json();
            const index = mockTickets.findIndex(t => t.id === id);
            if (index === -1) {
                return json({ success: false, error: 'Ticket not found' }, { status: 404 });
            }
            const updatedTicket = { ...mockTickets[index], ...data };
            if (data.status === 'picked_up') {
                updatedTicket.exit_date = new Date().toISOString();
                updatedTicket.payment_status = 'paid';
            }
            const newTickets = [...mockTickets];
            newTickets[index] = updatedTicket;
            setMockTickets(newTickets);
            return json({ success: true, data: updatedTicket });
        }
        if (method === 'DELETE') {
            const index = mockTickets.findIndex(t => t.id === id);
            if (index === -1) {
                return json({ success: false, error: 'Ticket not found' }, { status: 404 });
            }
            const newTickets = mockTickets.filter(t => t.id !== id);
            setMockTickets(newTickets);
            return json({ success: true, data: { deleted: true } });
        }
    }

    return null;
}
