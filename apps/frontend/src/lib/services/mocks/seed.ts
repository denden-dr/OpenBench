import type { MockTicket, MockProduct, MockSale, MockWarranty, MockUser } from './types';

export const initialTickets: MockTicket[] = [
  {
    id: '550e8400-e29b-41d4-a716-446655440001',
    ticket_number: 'OB-202606-0001',
    customer_name: 'Denden Hidayat',
    customer_phone: '081234567890',
    brand_phone: 'Samsung',
    model_phone: 'Galaxy S23 Ultra',
    serial_number: 'SN-S23U-992812',
    damage_description: 'Shattered screen after falling from a motorcycle.',
    repair_action: 'Replacement of S23 Ultra Original LCD Screen module.',
    cost: 2800000,
    status: 'completed',
    device_position: 'warehouse',
    payment_status: 'requesting',
    warranty_duration_days: 30,
    created_at: '2026-06-12T10:00:00.000Z'
  },
  {
    id: '550e8400-e29b-41d4-a716-446655440002',
    ticket_number: 'OB-202606-0002',
    customer_name: 'John Doe',
    customer_phone: '087799228833',
    brand_phone: 'Apple',
    model_phone: 'iPhone 14 Pro',
    serial_number: 'IMEI-358921029381023',
    damage_description: 'Battery drains very fast and device turns off unexpectedly (Battery health 68%).',
    repair_action: 'iPhone 14 Pro Original Apple Battery replacement + Reset Battery Cycle.',
    cost: 950000,
    status: 'in_repair',
    device_position: 'warehouse',
    payment_status: 'none',
    warranty_duration_days: 90,
    created_at: '2026-06-13T14:30:00.000Z'
  },
  {
    id: '550e8400-e29b-41d4-a716-446655440003',
    ticket_number: 'OB-202606-0003',
    customer_name: 'Alice Cooper',
    customer_phone: '082199223344',
    brand_phone: 'Xiaomi',
    model_phone: 'Redmi Note 12',
    serial_number: 'IMEI-869281029281928',
    damage_description: 'Broken USB port, unable to charge.',
    repair_action: 'Cleaning and replacing USB charger connector sub-board.',
    cost: 350000,
    status: 'completed',
    device_position: 'picked_up',
    payment_status: 'paid',
    payment_method: 'qris',
    warranty_duration_days: 14,
    picked_up_at: '2026-06-14T00:00:00.000Z',
    created_at: '2026-06-10T09:00:00.000Z'
  }
];

export const initialInventory: MockProduct[] = [
  {
    id: 'prod-001',
    name: 'Charger 25W Fast Charging Type-C',
    category: 'retail',
    stock: 18,
    price: 245000,
    cost_price: 150000,
    min_stock: 5
  },
  {
    id: 'prod-002',
    name: 'Tempered Glass Ultra Clear iPhone 14 Pro',
    category: 'retail',
    stock: 3,
    price: 95000,
    cost_price: 40000,
    min_stock: 5
  },
  {
    id: 'prod-003',
    name: 'LCD Screen Module Samsung S23 Ultra (Original)',
    category: 'spare_part',
    stock: 4,
    price: 2600000,
    cost_price: 1900000,
    min_stock: 2
  },
  {
    id: 'prod-004',
    name: 'Battery Replacement iPhone 14 Pro (Original)',
    category: 'spare_part',
    stock: 12,
    price: 850000,
    cost_price: 500000,
    min_stock: 3
  },
  {
    id: 'prod-005',
    name: 'Soft Case Transparent anti-crack (Universal)',
    category: 'retail',
    stock: 25,
    price: 49000,
    cost_price: 15000,
    min_stock: 5
  }
];

export const initialSales: MockSale[] = [
  {
    id: 'sale-001',
    invoice_number: 'INV-202606-0001',
    items: [
      { productId: 'prod-001', name: 'Charger 25W Fast Charging Type-C', price: 245000, qty: 1 }
    ],
    subtotal: 245000,
    discount: 0,
    total: 245000,
    payment_method: 'cash',
    created_at: '2026-06-12T11:15:00.000Z'
  },
  {
    id: 'sale-002',
    invoice_number: 'INV-202606-0002',
    items: [
      { productId: 'prod-002', name: 'Tempered Glass Ultra Clear iPhone 14 Pro', price: 95000, qty: 1 },
      { productId: 'prod-005', name: 'Soft Case Transparent anti-crack (Universal)', price: 49000, qty: 1 }
    ],
    subtotal: 144000,
    discount: 10000,
    total: 134000,
    payment_method: 'qris',
    created_at: '2026-06-13T16:20:00.000Z'
  }
];

export const initialWarranties: MockWarranty[] = [
  {
    id: 'war-001',
    ticket_id: '550e8400-e29b-41d4-a716-446655440003',
    ticket_number: 'OB-202606-0003',
    customer_name: 'Alice Cooper',
    device_info: 'Xiaomi Redmi Note 12',
    start_date: '2026-06-14T00:00:00.000Z',
    end_date: '2026-07-14T00:00:00.000Z',
    status: 'active'
  }
];

export const initialUsers: MockUser[] = [
  {
    id: 'mock-user-id-admin-12345',
    email: 'admin@openbench.dev',
    role: 'admin',
    passwordHash: 'SecureAdminPassword123!'
  }
];
