export let mockTickets = [
    {
        id: 'f7e2f418-d713-49b6-b9fe-28e630b1282a',
        customer_name: 'Budi Santoso',
        customer_gender: 'Male',
        customer_phone: '081234567890',
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
        id: 'c3a8b9c1-4d5e-6f7a-8b9c-0d1e2f3a4b5c',
        customer_name: 'Siti Aminah',
        customer_gender: 'Female',
        customer_phone: '081234567891',
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
        id: '9d8e7f6a-5b4c-3d2e-1f0a-9b8c7d6e5f4a',
        customer_name: 'Andi Wijaya',
        customer_gender: 'Male',
        customer_phone: '081234567892',
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
        id: '1a2b3c4d-5e6f-7a8b-9c0d-1e2f3a4b5c6d',
        customer_name: 'Denden Rahman',
        customer_gender: 'Male',
        customer_phone: '081234567893',
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
        id: '2b3c4d5e-6f7a-8b9c-0d1e-2f3a4b5c6d7e',
        customer_name: 'Eka Wijaya',
        customer_gender: 'Female',
        customer_phone: '081234567894',
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
