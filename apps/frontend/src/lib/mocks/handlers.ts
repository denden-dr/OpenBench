import { json } from '@sveltejs/kit';
import { mockTickets, setMockTickets, mockWarrantyClaims, setMockWarrantyClaims } from './mockData';

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

    // Matches /api/v1/warranty-claims
    if (path === '/api/v1/warranty-claims') {
        if (method === 'GET') {
            return json({ success: true, data: mockWarrantyClaims });
        }
        if (method === 'POST') {
            const data = await request.json();
            const originalTicket = mockTickets.find(t => t.id === data.ticket_id) as any;
            if (!originalTicket) {
                return json({ success: false, error: 'Original ticket not found' }, { status: 404 });
            }
            if (!data.issue?.trim()) {
                return json({ success: false, error: 'Issue is required' }, { status: 400 });
            }
            if (originalTicket.is_warranty) {
                return json({ success: false, error: 'Warranty claim tickets cannot spawn another warranty claim' }, { status: 400 });
            }
            if (originalTicket.status !== 'picked_up' || !originalTicket.exit_date) {
                return json({ success: false, error: 'Ticket is not eligible for warranty claim' }, { status: 400 });
            }

            const expiryTime = new Date(originalTicket.exit_date).getTime() + originalTicket.warranty_days * 86400000;
            if (Date.now() > expiryTime) {
                return json({ success: false, error: 'Warranty period has expired' }, { status: 400 });
            }

            const now = new Date().toISOString();
            const claim = {
                id: 'CLM-' + Date.now(),
                ticket_id: originalTicket.id,
                claim_ticket_id: null,
                issue: data.issue,
                additional_description: data.additional_description || '',
                status: 'waiting_inspection' as const,
                void_reason: null,
                inspected_at: null,
                created_at: now
            };

            setMockWarrantyClaims([claim, ...mockWarrantyClaims]);
            return json({ success: true, data: claim }, { status: 201 });
        }
    }

    // Matches /api/v1/warranty-claims/[id]/approve
    const claimApproveMatch = path.match(/^\/api\/v1\/warranty-claims\/([^/]+)\/approve$/);
    if (claimApproveMatch && method === 'POST') {
        const id = claimApproveMatch[1];
        const claimIndex = mockWarrantyClaims.findIndex(c => c.id === id);
        if (claimIndex === -1) {
            return json({ success: false, error: 'Claim not found' }, { status: 404 });
        }
        const claim = mockWarrantyClaims[claimIndex];
        const originalTicket = mockTickets.find(t => t.id === claim.ticket_id) as any;
        if (!originalTicket) {
            return json({ success: false, error: 'Original ticket not found' }, { status: 404 });
        }

        const now = new Date().toISOString();
        const spawnedTicketId = 'TCK-W' + Date.now().toString().slice(-6);

        const spawnedTicket = {
            id: spawnedTicketId,
            customer_name: originalTicket.customer_name,
            customer_gender: originalTicket.customer_gender,
            brand: originalTicket.brand,
            model: originalTicket.model,
            issue: '[Klaim Garansi] ' + claim.issue,
            additional_description: claim.additional_description || '',
            price: 0,
            status: 'on_process',
            payment_status: 'paid',
            warranty_days: 0,
            is_warranty: true,
            parent_ticket_id: originalTicket.id,
            entry_date: now
        };

        const updatedClaim = {
            ...claim,
            status: 'approved' as const,
            claim_ticket_id: spawnedTicketId,
            inspected_at: now
        };

        const newClaims = [...mockWarrantyClaims];
        newClaims[claimIndex] = updatedClaim;
        setMockWarrantyClaims(newClaims);

        setMockTickets([spawnedTicket, ...mockTickets]);

        return json({
            success: true,
            data: {
                claim: updatedClaim,
                ticket: spawnedTicket
            }
        });
    }

    // Matches /api/v1/warranty-claims/[id]/void
    const claimVoidMatch = path.match(/^\/api\/v1\/warranty-claims\/([^/]+)\/void$/);
    if (claimVoidMatch && method === 'POST') {
        const id = claimVoidMatch[1];
        const data = await request.json();
        if (!data.void_reason?.trim()) {
            return json({ success: false, error: 'Void reason is required' }, { status: 400 });
        }

        const claimIndex = mockWarrantyClaims.findIndex(c => c.id === id);
        if (claimIndex === -1) {
            return json({ success: false, error: 'Claim not found' }, { status: 404 });
        }
        const claim = mockWarrantyClaims[claimIndex];
        const originalTicket = mockTickets.find(t => t.id === claim.ticket_id) as any;
        if (!originalTicket) {
            return json({ success: false, error: 'Original ticket not found' }, { status: 404 });
        }

        const now = new Date().toISOString();
        const spawnedTicketId = 'TCK-W' + Date.now().toString().slice(-6);

		const spawnedTicket = {
            id: spawnedTicketId,
            customer_name: originalTicket.customer_name,
            customer_gender: originalTicket.customer_gender,
            brand: originalTicket.brand,
            model: originalTicket.model,
            issue: '[Klaim Ditolak] ' + claim.issue,
            additional_description: 'Klaim Garansi Ditolak. Alasan: ' + data.void_reason,
            price: 0,
            status: 'cancelled',
            payment_status: 'paid',
            warranty_days: 0,
            is_warranty: true,
            parent_ticket_id: originalTicket.id,
            entry_date: now
        };

        const updatedClaim = {
            ...claim,
            status: 'void' as const,
            claim_ticket_id: spawnedTicketId,
            void_reason: data.void_reason,
            inspected_at: now
        };

        const newClaims = [...mockWarrantyClaims];
        newClaims[claimIndex] = updatedClaim;
        setMockWarrantyClaims(newClaims);

        setMockTickets([spawnedTicket, ...mockTickets]);

        return json({
            success: true,
            data: {
                claim: updatedClaim,
                ticket: spawnedTicket
            }
        });
    }

    return null;
}
