import { mockDbService } from './db';
import { mockAuthService } from './auth';

let isInitialized = false;
let originalFetch: typeof window.fetch | null = null;

// Helper to construct a Standard Response Envelope
function createResponse(code: number, message: string, data: any): Response {
  const body: any = {
    code,
    message
  };

  // Backend requires data field for responses, even if empty array
  if (data !== null && data !== undefined) {
    body.data = data;
  }

  return new Response(
    JSON.stringify(body),
    {
      status: code,
      headers: {
        'Content-Type': 'application/json'
      }
    }
  );
}

export function initMockNetwork() {
  if (isInitialized || typeof window === 'undefined') return;
  
  originalFetch = window.fetch;
  
  window.fetch = async (input: RequestInfo | URL, init?: RequestInit): Promise<Response> => {
    const urlStr = typeof input === 'string' ? input : (input instanceof URL ? input.toString() : input.url);
    const method = init?.method?.toUpperCase() || 'GET';
    
    // We only intercept calls targeting our API /api/v1
    const match = urlStr.match(/\/api\/v1(\/.*)$/);
    if (!match) {
      return originalFetch!(input, init);
    }
    
    const path = match[1]; // e.g. "/auth/signin" or "/admin/tickets"
    
    try {
      // 1. Auth Endpoints
      if (path === '/auth/signin' && method === 'POST') {
        const body = JSON.parse(init?.body as string || '{}');
        const session = await mockAuthService.signIn(body.email, body.password);
        return createResponse(200, 'Successfully signed in', session);
      }
      
      if (path === '/auth/signout' && method === 'POST') {
        await mockAuthService.signOut();
        return createResponse(200, 'Successfully signed out', null);
      }
      
      if (path === '/auth/refresh' && method === 'POST') {
        return createResponse(200, 'Tokens rotated successfully', null);
      }
      
      if (path === '/auth/me' && method === 'GET') {
        const session = mockAuthService.getSession();
        if (!session) {
          return createResponse(401, 'Authentication required', null);
        }
        return createResponse(200, 'Success', session);
      }
      
      // 2. Tickets Endpoints
      if (path === '/admin/tickets' && method === 'GET') {
        const tickets = await mockDbService.getTickets();
        return createResponse(200, 'Success', tickets);
      }
      
      if (path === '/admin/tickets' && method === 'POST') {
        const body = JSON.parse(init?.body as string || '{}');
        const newTicket = await mockDbService.createTicket(body);
        return createResponse(201, 'Ticket created successfully', newTicket);
      }
      
      // GET /admin/tickets/:id
      const adminTicketMatch = path.match(/^\/admin\/tickets\/([a-f0-9-]+)$/i);
      if (adminTicketMatch && method === 'GET') {
        const id = adminTicketMatch[1];
        const ticket = await mockDbService.getTicket(id);
        if (!ticket) {
          return createResponse(404, 'Ticket not found', null);
        }
        return createResponse(200, 'Success', ticket);
      }
      
      // PATCH /admin/tickets/:id
      if (adminTicketMatch && method === 'PATCH') {
        const id = adminTicketMatch[1];
        const body = JSON.parse(init?.body as string || '{}');
        const updated = await mockDbService.updateTicket(id, body);
        return createResponse(200, 'Ticket updated successfully', updated);
      }

      // POST /admin/tickets/:id/emergency
      const emergencyMatch = path.match(/^\/admin\/tickets\/([a-f0-9-]+)\/emergency$/i);
      if (emergencyMatch && method === 'POST') {
        const id = emergencyMatch[1];
        const body = JSON.parse(init?.body as string || '{}');
        const updated = await mockDbService.updateTicket(id, body);
        return createResponse(200, 'Ticket updated successfully (emergency)', updated);
      }
      
      // Public Tracker Lookup: GET /tracker/:id
      const trackerMatch = path.match(/^\/tracker\/([a-f0-9-]+)$/i);
      if (trackerMatch && method === 'GET') {
        const id = trackerMatch[1];
        const ticket = await mockDbService.getTicket(id);
        if (!ticket) {
          return createResponse(404, 'Ticket not found', null);
        }

        const maskName = (name: string): string => {
          return name.split(' ').map(word => {
            if (!word) return '';
            const len = word.length;
            if (len <= 4) {
              return word[0] + '*'.repeat(len - 1);
            } else {
              return word.slice(0, 3) + '*'.repeat(len - 3);
            }
          }).join(' ');
        };

        const maskPhone = (phone: string): string => {
          const len = phone.length;
          if (len <= 3) return phone;
          return '*'.repeat(len - 3) + phone.slice(-3);
        };

        // Only return non-sensitive fields for the public tracker
        const publicInfo = {
          id: ticket.id,
          customer_name_masked: maskName(ticket.customer_name),
          customer_phone_masked: maskPhone(ticket.customer_phone),
          brand_phone: ticket.brand_phone,
          model_phone: ticket.model_phone,
          damage_description: ticket.damage_description,
          repair_action: ticket.repair_action,
          status: ticket.status,
          warranty_duration_days: ticket.warranty_duration_days,
          picked_up_at: ticket.picked_up_at,
          warranty_expiry_date: ticket.warranty_expiry_date,
          created_at: ticket.created_at
        };
        return createResponse(200, 'Success', publicInfo);
      }
      
      // 3. Inventory Endpoints
      if (path === '/admin/inventory' && method === 'GET') {
        const products = await mockDbService.getInventory();
        return createResponse(200, 'Success', products);
      }
      
      if (path === '/admin/inventory' && method === 'POST') {
        const body = JSON.parse(init?.body as string || '{}');
        const newProduct = await mockDbService.createProduct(body);
        return createResponse(201, 'Product created successfully', newProduct);
      }
      
      const adminInventoryMatch = path.match(/^\/admin\/inventory\/([a-z0-9-]+)$/i);
      if (adminInventoryMatch && method === 'GET') {
        const id = adminInventoryMatch[1];
        const product = await mockDbService.getProduct(id);
        if (!product) {
          return createResponse(404, 'Product not found', null);
        }
        return createResponse(200, 'Success', product);
      }
      
      if (adminInventoryMatch && method === 'PATCH') {
        const id = adminInventoryMatch[1];
        const body = JSON.parse(init?.body as string || '{}');
        const updated = await mockDbService.updateProduct(id, body);
        return createResponse(200, 'Product updated successfully', updated);
      }
      
      if (adminInventoryMatch && method === 'DELETE') {
        const id = adminInventoryMatch[1];
        await mockDbService.deleteProduct(id);
        return createResponse(200, 'Product deleted successfully', null);
      }
      
      // 4. Sales Endpoints
      if (path === '/admin/sales' && method === 'GET') {
        const sales = await mockDbService.getSales();
        return createResponse(200, 'Success', sales);
      }
      
      if (path === '/admin/sales' && method === 'POST') {
        const body = JSON.parse(init?.body as string || '{}');
        const newSale = await mockDbService.createSale(body);
        return createResponse(201, 'Sale logged successfully', newSale);
      }
      
      // 5. Warranties Endpoints
      if (path === '/admin/warranties' && method === 'GET') {
        const warranties = await mockDbService.getWarranties();
        return createResponse(200, 'Success', warranties);
      }
      
      // Match not found
      return originalFetch!(input, init);
    } catch (err: any) {
      return createResponse(500, err.message || 'Internal Server Error', null);
    }
  };
  
  isInitialized = true;
}
