import { json } from '@sveltejs/kit';
import { mockTickets, setMockTickets, mockWarrantyClaims, setMockWarrantyClaims } from './mockData';

export async function handleMockRequest(event: any): Promise<Response | null> {
    const { request, url } = event;
    const path = url.pathname;
    const method = request.method;

    // Matches /api/v1/tickets
    if (path === '/api/v1/tickets') {
        if (method === 'GET') {
            const page = parseInt(url.searchParams.get('page') || '1', 10);
            const limit = parseInt(url.searchParams.get('limit') || '20', 10);
            const search = (url.searchParams.get('search') || '').trim().toLowerCase();
            const status = (url.searchParams.get('status') || 'all').trim().toLowerCase();

            // Hitung statistik status counts berdasar filter search (tapi TIDAK memfilter status)
            const countFiltered = mockTickets.filter(t => {
                if (!search) return true;
                return (
                    t.id.toLowerCase().includes(search) ||
                    t.customer_name.toLowerCase().includes(search) ||
                    t.brand.toLowerCase().includes(search) ||
                    t.model.toLowerCase().includes(search) ||
                    t.issue.toLowerCase().includes(search)
                );
            });

            const statusCounts = {
                all: countFiltered.filter(t => t.status !== 'picked_up').length,
                service_in: countFiltered.filter(t => t.status === 'service_in').length,
                on_process: countFiltered.filter(t => t.status === 'on_process').length,
                waiting_confirmation: countFiltered.filter(t => t.status === 'waiting_confirmation').length,
                fixed: countFiltered.filter(t => t.status === 'fixed').length,
                picked_up: countFiltered.filter(t => t.status === 'picked_up').length,
                cancelled: countFiltered.filter(t => t.status === 'cancelled').length,
            };

            // Filter data untuk list response (memfilter status + search)
            let filtered = countFiltered;
            if (status === 'all') {
                filtered = filtered.filter(t => t.status !== 'picked_up');
            } else {
                filtered = filtered.filter(t => t.status === status);
            }

            // Slicing offset-based pagination
            const total = filtered.length;
            const totalPages = Math.ceil(total / limit) || 1;
            const start = (page - 1) * limit;
            const paginatedData = filtered.slice(start, start + limit);

            return json({
                code: 200,
                message: 'Success',
                data: paginatedData,
                total,
                total_pages: totalPages,
                page,
                limit,
                status_counts: statusCounts
            });
        }
        if (method === 'POST') {
            const data = await request.json();
            const newTicket = {
                id: crypto.randomUUID(),
                ...data,
                status: 'service_in',
                payment_status: 'unpaid',
                entry_date: new Date().toISOString()
            };
            setMockTickets([newTicket, ...mockTickets]);
            return json({ code: 201, message: 'Success', data: newTicket }, { status: 201 });
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
                return json({
                    type: 'https://openbench.denden.com/errors/not-found',
                    title: 'Not Found',
                    status: 404,
                    detail: 'Ticket not found',
                    instance: path
                }, { status: 404 });
            }
            const updatedTicket = { ...mockTickets[index], ...data };
            if (data.status === 'picked_up') {
                updatedTicket.exit_date = new Date().toISOString();
                updatedTicket.payment_status = 'paid';
            }
            const newTickets = [...mockTickets];
            newTickets[index] = updatedTicket;
            setMockTickets(newTickets);
            return json({ code: 200, message: 'Success', data: updatedTicket });
        }
        if (method === 'DELETE') {
            const index = mockTickets.findIndex(t => t.id === id);
            if (index === -1) {
                return json({
                    type: 'https://openbench.denden.com/errors/not-found',
                    title: 'Not Found',
                    status: 404,
                    detail: 'Ticket not found',
                    instance: path
                }, { status: 404 });
            }
            const newTickets = mockTickets.filter(t => t.id !== id);
            setMockTickets(newTickets);
            return json({ code: 200, message: 'Ticket deleted successfully', data: { deleted: true } });
        }
    }

    // Matches /api/v1/warranty-claims
    if (path === '/api/v1/warranty-claims') {
        if (method === 'GET') {
            const page = parseInt(url.searchParams.get('page') || '1', 10);
            const limit = parseInt(url.searchParams.get('limit') || '10', 10);
            const status = url.searchParams.get('status') || 'all';

            let filtered = [...mockWarrantyClaims];
            if (status !== 'all') {
                filtered = filtered.filter(c => c.status === status);
            }

            // Inject ticket details for convenience mapping
            const claimsWithTicket = filtered.map(c => {
                const ticket = mockTickets.find(t => t.id === c.ticket_id);
                return {
                    ...c,
                    originalTicket: ticket || null
                };
            });

            const total = claimsWithTicket.length;
            const totalPages = Math.ceil(total / limit) || 1;
            const start = (page - 1) * limit;
            const paginatedData = claimsWithTicket.slice(start, start + limit);

            return json({
                code: 200,
                message: 'Success',
                data: paginatedData,
                total,
                total_pages: totalPages,
                page,
                limit
            });
        }
        if (method === 'POST') {
            const data = await request.json();
            const originalTicket = mockTickets.find(t => t.id === data.ticket_id) as any;
            if (!originalTicket) {
                return json({
                    type: 'https://openbench.denden.com/errors/not-found',
                    title: 'Not Found',
                    status: 404,
                    detail: 'Original ticket not found',
                    instance: path
                }, { status: 404 });
            }
            if (!data.issue?.trim()) {
                return json({
                    type: 'https://openbench.denden.com/errors/validation-failed',
                    title: 'Validation Failed',
                    status: 400,
                    detail: 'Issue is required',
                    instance: path
                }, { status: 400 });
            }
            if (originalTicket.is_warranty) {
                return json({
                    type: 'https://openbench.denden.com/errors/validation-failed',
                    title: 'Validation Failed',
                    status: 400,
                    detail: 'Warranty claim tickets cannot spawn another warranty claim',
                    instance: path
                }, { status: 400 });
            }
            if (originalTicket.status !== 'picked_up' || !originalTicket.exit_date) {
                return json({
                    type: 'https://openbench.denden.com/errors/validation-failed',
                    title: 'Validation Failed',
                    status: 400,
                    detail: 'Ticket is not eligible for warranty claim',
                    instance: path
                }, { status: 400 });
            }

            const expiryTime = new Date(originalTicket.exit_date).getTime() + originalTicket.warranty_days * 86400000;
            if (Date.now() > expiryTime) {
                return json({
                    type: 'https://openbench.denden.com/errors/validation-failed',
                    title: 'Validation Failed',
                    status: 400,
                    detail: 'Warranty period has expired',
                    instance: path
                }, { status: 400 });
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
            return json({ code: 201, message: 'Success', data: claim }, { status: 201 });
        }
    }

    // Matches /api/v1/warranty-claims/[id]/approve
    const claimApproveMatch = path.match(/^\/api\/v1\/warranty-claims\/([^/]+)\/approve$/);
    if (claimApproveMatch && method === 'POST') {
        const id = claimApproveMatch[1];
        const claimIndex = mockWarrantyClaims.findIndex(c => c.id === id);
        if (claimIndex === -1) {
            return json({
                type: 'https://openbench.denden.com/errors/not-found',
                title: 'Not Found',
                status: 404,
                detail: 'Claim not found',
                instance: path
            }, { status: 404 });
        }
        const claim = mockWarrantyClaims[claimIndex];
        const originalTicket = mockTickets.find(t => t.id === claim.ticket_id) as any;
        if (!originalTicket) {
            return json({
                type: 'https://openbench.denden.com/errors/not-found',
                title: 'Not Found',
                status: 404,
                detail: 'Original ticket not found',
                instance: path
            }, { status: 404 });
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
            code: 200,
            message: 'Success',
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
            return json({
                type: 'https://openbench.denden.com/errors/validation-failed',
                title: 'Validation Failed',
                status: 400,
                detail: 'Void reason is required',
                instance: path
            }, { status: 400 });
        }

        const claimIndex = mockWarrantyClaims.findIndex(c => c.id === id);
        if (claimIndex === -1) {
            return json({
                type: 'https://openbench.denden.com/errors/not-found',
                title: 'Not Found',
                status: 404,
                detail: 'Claim not found',
                instance: path
            }, { status: 404 });
        }
        const claim = mockWarrantyClaims[claimIndex];
        const originalTicket = mockTickets.find(t => t.id === claim.ticket_id) as any;
        if (!originalTicket) {
            return json({
                type: 'https://openbench.denden.com/errors/not-found',
                title: 'Not Found',
                status: 404,
                detail: 'Original ticket not found',
                instance: path
            }, { status: 404 });
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
            code: 200,
            message: 'Success',
            data: {
                claim: updatedClaim,
                ticket: spawnedTicket
            }
        });
    }

    // 1. Matches /api/v1/public/tickets/[id]
    const publicTicketMatch = path.match(/^\/api\/v1\/public\/tickets\/([^/]+)$/);
    if (publicTicketMatch && method === 'GET') {
        const id = publicTicketMatch[1];
        if (id.length !== 36) {
            return json({
                type: 'https://openbench.denden.com/errors/validation-failed',
                title: 'Validation Failed',
                status: 400,
                detail: 'Invalid ticket ID format. Only full UUID is supported.',
                instance: path
            }, { status: 400 });
        }
        const ticket = mockTickets.find(t => t.id.toLowerCase() === id.toLowerCase());
        if (!ticket) {
            return json({
                type: 'https://openbench.denden.com/errors/not-found',
                title: 'Not Found',
                status: 404,
                detail: 'Tiket tidak ditemukan',
                instance: path
            }, { status: 404 });
        }
        const maskName = (name: string) => name.split(' ').map(n => n[0] + '*'.repeat(Math.max(0, n.length - 1))).join(' ');
        const maskPhone = (phone: string) => phone ? phone.slice(0, 4) + '*'.repeat(Math.max(0, phone.length - 6)) + phone.slice(-2) : '';

        return json({
            code: 200,
            message: 'Success',
            data: {
                id: ticket.id,
                customer_name_masked: maskName(ticket.customer_name),
                customer_phone_masked: maskPhone((ticket as any).customer_phone || ''),
                brand: ticket.brand,
                model: ticket.model,
                issue: ticket.issue,
                additional_description: (ticket as any).additional_description || null,
                accessories: (ticket as any).accessories || null,
                price: ticket.price,
                payment_status: ticket.payment_status,
                status: ticket.status,
                entry_date: ticket.entry_date,
                exit_date: (ticket as any).exit_date || null,
                warranty_days: ticket.warranty_days || 0
            }
        });
    }

    // 2. Matches /api/v1/public/track
    if (path === '/api/v1/public/track' && method === 'POST') {
        const data = await request.json();
        const { short_id, phone } = data;
        if (!short_id || !phone) {
            return json({
                type: 'https://openbench.denden.com/errors/validation-failed',
                title: 'Validation Failed',
                status: 400,
                detail: 'ID Tiket dan Nomor Telepon wajib diisi',
                instance: path
            }, { status: 400 });
        }
        const ticket = mockTickets.find(t => {
            const query = short_id.toLowerCase();
            const target = t.id.toLowerCase();
            let idMatch = false;
            if (query.length === 8) {
                idMatch = target.startsWith(query);
            } else if (query.length === 36) {
                idMatch = target === query;
            }
            const cleanPhoneT = ((t as any).customer_phone || '').replace(/\D/g, '');
            const cleanPhoneInput = phone.replace(/\D/g, '');
            return idMatch && cleanPhoneT === cleanPhoneInput;
        });
        if (!ticket) {
            return json({
                type: 'https://openbench.denden.com/errors/not-found',
                title: 'Not Found',
                status: 404,
                detail: 'Tiket tidak ditemukan atau nomor HP salah',
                instance: path
            }, { status: 404 });
        }
        return json({
            code: 200,
            message: 'Success',
            data: {
                ticket_id: ticket.id
            }
        });
    }

    return null;
}
