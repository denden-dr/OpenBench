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
    }
];

export function setMockTickets(newTickets: any[]) {
    mockTickets = newTickets;
}
