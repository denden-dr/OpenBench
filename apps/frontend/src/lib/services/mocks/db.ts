import type { MockTicket, MockProduct, MockSale, MockWarranty, MockUser } from './types';
import type { TicketCreate } from '../ticket';
import type { components } from '$lib/api/openapi.gen';
import {
  initialTickets,
  initialInventory,
  initialSales,
  initialWarranties,
  initialUsers
} from './seed';

const KEYS = {
  TICKETS: 'openbench_mock_tickets',
  INVENTORY: 'openbench_mock_inventory',
  SALES: 'openbench_mock_sales',
  WARRANTIES: 'openbench_mock_warranties',
  USERS: 'openbench_mock_users'
};

function getTickets(): MockTicket[] {
  if (typeof window === 'undefined') return initialTickets;
  const stored = localStorage.getItem(KEYS.TICKETS);
  if (!stored) {
    localStorage.setItem(KEYS.TICKETS, JSON.stringify(initialTickets));
    return initialTickets;
  }
  try {
    return JSON.parse(stored);
  } catch {
    localStorage.setItem(KEYS.TICKETS, JSON.stringify(initialTickets));
    return initialTickets;
  }
}

function saveTickets(tickets: MockTicket[]): void {
  if (typeof window !== 'undefined') {
    localStorage.setItem(KEYS.TICKETS, JSON.stringify(tickets));
  }
}

function getInventory(): MockProduct[] {
  if (typeof window === 'undefined') return initialInventory;
  const stored = localStorage.getItem(KEYS.INVENTORY);
  if (!stored) {
    localStorage.setItem(KEYS.INVENTORY, JSON.stringify(initialInventory));
    return initialInventory;
  }
  try {
    return JSON.parse(stored);
  } catch {
    localStorage.setItem(KEYS.INVENTORY, JSON.stringify(initialInventory));
    return initialInventory;
  }
}

function saveInventory(inventory: MockProduct[]): void {
  if (typeof window !== 'undefined') {
    localStorage.setItem(KEYS.INVENTORY, JSON.stringify(inventory));
  }
}

function getSales(): MockSale[] {
  if (typeof window === 'undefined') return initialSales;
  const stored = localStorage.getItem(KEYS.SALES);
  if (!stored) {
    localStorage.setItem(KEYS.SALES, JSON.stringify(initialSales));
    return initialSales;
  }
  try {
    return JSON.parse(stored);
  } catch {
    localStorage.setItem(KEYS.SALES, JSON.stringify(initialSales));
    return initialSales;
  }
}

function saveSales(sales: MockSale[]): void {
  if (typeof window !== 'undefined') {
    localStorage.setItem(KEYS.SALES, JSON.stringify(sales));
  }
}

function getWarranties(): MockWarranty[] {
  if (typeof window === 'undefined') return initialWarranties;
  const stored = localStorage.getItem(KEYS.WARRANTIES);
  if (!stored) {
    localStorage.setItem(KEYS.WARRANTIES, JSON.stringify(initialWarranties));
    return initialWarranties;
  }
  try {
    return JSON.parse(stored);
  } catch {
    localStorage.setItem(KEYS.WARRANTIES, JSON.stringify(initialWarranties));
    return initialWarranties;
  }
}

function saveWarranties(warranties: MockWarranty[]): void {
  if (typeof window !== 'undefined') {
    localStorage.setItem(KEYS.WARRANTIES, JSON.stringify(warranties));
  }
}

function getUsers(): MockUser[] {
  if (typeof window === 'undefined') return initialUsers;
  const stored = localStorage.getItem(KEYS.USERS);
  if (!stored) {
    localStorage.setItem(KEYS.USERS, JSON.stringify(initialUsers));
    return initialUsers;
  }
  try {
    return JSON.parse(stored);
  } catch {
    localStorage.setItem(KEYS.USERS, JSON.stringify(initialUsers));
    return initialUsers;
  }
}

function saveUsers(users: MockUser[]): void {
  if (typeof window !== 'undefined') {
    localStorage.setItem(KEYS.USERS, JSON.stringify(users));
  }
}

// SIMULATE network latency
const LATENCY = 400;
const delay = () => new Promise(resolve => setTimeout(resolve, LATENCY));

export const mockDbService = {
  // TICKETS API
  async getTickets(): Promise<MockTicket[]> {
    await delay();
    return getTickets();
  },

  async getTicket(id: string): Promise<MockTicket | null> {
    await delay();
    const tickets = getTickets();
    return tickets.find(t => t.id === id) || null;
  },

  async getTicketByNumber(ticketNumber: string): Promise<MockTicket | null> {
    await delay();
    const tickets = getTickets();
    return tickets.find(t => t.ticket_number.toLowerCase() === ticketNumber.toLowerCase()) || null;
  },

  async createTicket(ticket: TicketCreate): Promise<MockTicket> {
    await delay();
    const tickets = getTickets();

    // Generate simple sequential ticket number with random 4-char alphanumeric suffix
    const count = tickets.length + 1;
    const chars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789';
    let randomSuffix = '';
    for (let i = 0; i < 4; i++) {
      randomSuffix += chars.charAt(Math.floor(Math.random() * chars.length));
    }
    const ticket_number = `OB-202606-${count.toString().padStart(4, '0')}-${randomSuffix}`;

    const newTicket: MockTicket = {
      status: 'received',
      device_position: 'warehouse',
      payment_status: 'none',
      ...ticket,
      cost: ticket.cost !== undefined ? ticket.cost : 0,
      repair_action: ticket.repair_action || '',
      serial_number: ticket.serial_number || 'N/A',
      id: crypto.randomUUID(),
      ticket_number,
      created_at: new Date().toISOString()
    };

    tickets.unshift(newTicket);
    saveTickets(tickets);
    return newTicket;
  },

  async updateTicket(id: string, updates: Partial<MockTicket>): Promise<MockTicket> {
    await delay();
    const tickets = getTickets();
    const idx = tickets.findIndex(t => t.id === id);
    if (idx === -1) throw new Error('Ticket not found.');

    const oldTicket = tickets[idx];
    const updatedTicket = { ...oldTicket, ...updates };

    const isReversal = oldTicket.device_position === 'picked_up' && updates.device_position !== undefined && updates.device_position !== 'picked_up';
    if (isReversal) {
      delete updatedTicket.picked_up_at;
      delete updatedTicket.warranty_expiry_date;
      const warranties = getWarranties().filter(w => w.ticket_id !== id);
      saveWarranties(warranties);
    }

    if (updates.warranty_duration_days !== undefined && updatedTicket.device_position === 'picked_up' && updatedTicket.picked_up_at) {
      const start = new Date(updatedTicket.picked_up_at);
      const end = new Date(start);
      end.setDate(start.getDate() + updates.warranty_duration_days);
      updatedTicket.warranty_expiry_date = end.toISOString();

      const warranties = getWarranties();
      const warIdx = warranties.findIndex(w => w.ticket_id === id);
      if (warIdx !== -1) {
        warranties[warIdx].end_date = end.toISOString();
        if (updates.customer_name) {
          warranties[warIdx].customer_name = updates.customer_name;
        }
        warranties[warIdx].device_info = `${updatedTicket.brand_phone} ${updatedTicket.model_phone}`;
        saveWarranties(warranties);
      }
    }

    // If device_position changes to 'picked_up', automatically trigger warranty and update payment/device details
    if (updates.device_position === 'picked_up' && oldTicket.device_position !== 'picked_up') {
      updatedTicket.device_position = 'picked_up';
      if (updatedTicket.status !== 'completed' && updatedTicket.status !== 'cancelled') {
        updatedTicket.status = 'completed';
      }
      if (updatedTicket.status === 'completed') {
        updatedTicket.payment_status = 'paid';
      }

      const durationDays = updates.warranty_duration_days !== undefined
        ? updates.warranty_duration_days
        : (oldTicket.warranty_duration_days !== undefined ? oldTicket.warranty_duration_days : 30);
      const startDate = new Date();
      const endDate = new Date();
      endDate.setDate(startDate.getDate() + durationDays);

      updatedTicket.picked_up_at = startDate.toISOString();
      updatedTicket.warranty_expiry_date = endDate.toISOString();

      if (updatedTicket.status === 'completed') {
        // Push warranty record
        const warranties = getWarranties();
        const warranty: MockWarranty = {
          id: `war-${warranties.length + 1}`,
          ticket_id: id,
          ticket_number: oldTicket.ticket_number,
          customer_name: oldTicket.customer_name,
          device_info: `${oldTicket.brand_phone} ${oldTicket.model_phone}`,
          start_date: startDate.toISOString(),
          end_date: endDate.toISOString(),
          status: 'active'
        };

        warranties.unshift(warranty);
        saveWarranties(warranties);
      }
    }

    tickets[idx] = updatedTicket;
    saveTickets(tickets);
    return updatedTicket;
  },

  // INVENTORY API
  async getInventory(): Promise<MockProduct[]> {
    await delay();
    return getInventory();
  },

  async getProduct(id: string): Promise<MockProduct | null> {
    await delay();
    return getInventory().find(p => p.id === id) || null;
  },

  async createProduct(product: Omit<MockProduct, 'id'>): Promise<MockProduct> {
    await delay();
    const inventory = getInventory();
    const newProduct: MockProduct = {
      ...product,
      id: `prod-${inventory.length + 1}`
    };
    inventory.push(newProduct);
    saveInventory(inventory);
    return newProduct;
  },

  async updateProduct(id: string, updates: Partial<MockProduct>): Promise<MockProduct> {
    await delay();
    const inventory = getInventory();
    const idx = inventory.findIndex(p => p.id === id);
    if (idx === -1) throw new Error('Product not found.');

    const updated = { ...inventory[idx], ...updates };
    inventory[idx] = updated;
    saveInventory(inventory);
    return updated;
  },

  async deleteProduct(id: string): Promise<void> {
    await delay();
    let inventory = getInventory();
    inventory = inventory.filter(p => p.id !== id);
    saveInventory(inventory);
  },

  // SALES (POS) API
  async getSales(): Promise<MockSale[]> {
    await delay();
    return getSales();
  },

  async createSale(sale: components['schemas']['SaleCreate']): Promise<MockSale> {
    await delay();
    const sales = getSales();
    const inventory = getInventory();

    let subtotal = 0;
    const finalItems = [];

    // 0. Accumulate quantities
    const itemQtyMap: Record<string, number> = {};
    for (const item of sale.items) {
      itemQtyMap[item.product_id] = (itemQtyMap[item.product_id] || 0) + item.qty;
    }

    // 1. Check stock and compute subtotal
    for (const [prodId, qty] of Object.entries(itemQtyMap)) {
      const pIdx = inventory.findIndex(p => p.id === prodId);
      if (pIdx === -1) {
        throw new Error(`Product not found`);
      }
      const product = inventory[pIdx];
      if (product.stock < qty) {
        throw new Error(`Insufficient stock: ${product.name} (available: ${product.stock}, requested: ${qty})`);
      }
      subtotal += product.price * qty;
      finalItems.push({
        product_id: product.id,
        name: product.name,
        price: product.price,
        qty: qty
      });
    }

    if (sale.discount > subtotal) {
      throw new Error('Discount cannot exceed subtotal');
    }

    // 2. Deduct inventory stocks
    for (const [prodId, qty] of Object.entries(itemQtyMap)) {
      const pIdx = inventory.findIndex(p => p.id === prodId);
      inventory[pIdx].stock -= qty;
    }
    saveInventory(inventory);

    // 3. Generate Invoice Number
    const count = sales.length + 1;
    const invoice_number = `INV-202606-${count.toString().padStart(4, '0')}`;

    const newSale: MockSale = {
      ...sale,
      id: `sale-${count}`,
      invoice_number,
      subtotal,
      total: subtotal - sale.discount,
      items: finalItems,
      created_at: new Date().toISOString()
    };

    sales.unshift(newSale);
    saveSales(sales);
    return newSale;
  },

  // WARRANTIES API
  async getWarranties(): Promise<MockWarranty[]> {
    await delay();
    return getWarranties();
  },

  // USERS / AUTH API
  async getUserByEmail(email: string): Promise<MockUser | null> {
    await delay();
    const users = getUsers();
    return users.find(u => u.email.toLowerCase() === email.trim().toLowerCase()) || null;
  },

  async createUser(user: Omit<MockUser, 'id' | 'role'> & { role?: 'admin' | 'user' }): Promise<MockUser> {
    await delay();
    const users = getUsers();
    const existing = users.find(u => u.email.toLowerCase() === user.email.trim().toLowerCase());
    if (existing) {
      throw new Error('User with this email already exists.');
    }
    const newUser: MockUser = {
      ...user,
      id: crypto.randomUUID(),
      role: user.role || 'user'
    };
    users.push(newUser);
    saveUsers(users);
    return newUser;
  },

  async updateUser(id: string, updates: Partial<MockUser>): Promise<MockUser> {
    await delay();
    const users = getUsers();
    const idx = users.findIndex(u => u.id === id);
    if (idx === -1) throw new Error('User not found.');

    const updatedUser = { ...users[idx], ...updates };
    users[idx] = updatedUser;
    saveUsers(users);

    // If this is the logged-in user, update the active session cache
    const session = this.getActiveSession();
    if (session && session.userId === id) {
      const updatedSession = {
        ...session,
        username: updatedUser.username,
        full_name: updatedUser.full_name,
        phone_number: updatedUser.phone_number
      };
      this.saveActiveSession(updatedSession);
    }

    return updatedUser;
  },

  async getUserTickets(userId: string): Promise<MockTicket[]> {
    await delay();
    const tickets = getTickets();
    return tickets.filter(t => t.user_id === userId);
  },

  saveActiveSession(session: any): void {
    if (typeof window !== 'undefined') {
      sessionStorage.setItem('openbench_session', JSON.stringify(session));
    }
  },

  getActiveSession(): any | null {
    if (typeof window === 'undefined') return null;
    const data = sessionStorage.getItem('openbench_session');
    if (!data) return null;
    try {
      return JSON.parse(data);
    } catch {
      return null;
    }
  },

  clearActiveSession(): void {
    if (typeof window !== 'undefined') {
      sessionStorage.removeItem('openbench_session');
    }
  }
};
