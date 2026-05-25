export let mockTickets = [
    {
        id: 'TCK-001',
        customer_name: 'Budi Santoso',
        customer_gender: 'Male',
        brand: 'Apple',
        model: 'iPhone 13',
        issue: 'Layar Retak',
        price: 1500000,
        status: 'service_in',
        payment_status: 'unpaid',
        warranty_days: 30,
        entry_date: new Date(Date.now() - 86400000).toISOString()
    },
    {
        id: 'TCK-002',
        customer_name: 'Siti Aminah',
        customer_gender: 'Female',
        brand: 'Samsung',
        model: 'Galaxy S22',
        issue: 'Baterai Drop',
        price: 800000,
        status: 'on_process',
        payment_status: 'unpaid',
        warranty_days: 30,
        entry_date: new Date().toISOString()
    },
    {
        id: 'TCK-003',
        customer_name: 'Andi Wijaya',
        customer_gender: 'Male',
        brand: 'Xiaomi',
        model: 'Redmi Note 10',
        issue: 'Mati Total',
        price: 450000,
        status: 'waiting_confirmation',
        additional_description: '[Kendala Teknisi]: Ternyata IC Power jebol.',
        payment_status: 'unpaid',
        warranty_days: 14,
        entry_date: new Date().toISOString()
    },
    {
        id: 'TCK-004',
        customer_name: 'Denden Rahman',
        customer_gender: 'Male',
        brand: 'Apple',
        model: 'iPhone 15 Pro',
        issue: 'Kerusakan LCD',
        price: 2500000,
        status: 'picked_up',
        payment_status: 'paid',
        warranty_days: 30,
        entry_date: new Date(Date.now() - 10 * 86400000).toISOString(),
        exit_date: new Date(Date.now() - 9 * 86400000).toISOString()
    },
    {
        id: 'TCK-005',
        customer_name: 'Eka Wijaya',
        customer_gender: 'Female',
        brand: 'Xiaomi',
        model: 'Mi 11 Ultra',
        issue: 'Kamera Error',
        price: 1200000,
        status: 'picked_up',
        payment_status: 'paid',
        warranty_days: 15,
        entry_date: new Date(Date.now() - 40 * 86400000).toISOString(),
        exit_date: new Date(Date.now() - 39 * 86400000).toISOString()
    }
];

export function setMockTickets(newTickets: any[]) {
    mockTickets = newTickets;
}

export interface MockClaim {
  id: string;
  ticket_id: string;
  claim_ticket_id: string | null;
  issue: string;
  additional_description: string;
  status: 'waiting_inspection' | 'approved' | 'void';
  void_reason: string | null;
  inspected_at: string | null;
  created_at: string;
}

export let mockWarrantyClaims: MockClaim[] = [];

export function setMockWarrantyClaims(newClaims: MockClaim[]) {
  mockWarrantyClaims = newClaims;
}
