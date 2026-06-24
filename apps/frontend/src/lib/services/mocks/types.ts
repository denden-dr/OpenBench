export interface MockTicket {
  id: string; // UUID
  ticket_number: string; // Structured: OB-YYYYMM-XXXX
  customer_name: string;
  customer_phone: string;
  brand_phone: string;
  model_phone: string;
  serial_number: string;
  damage_description: string;
  repair_action: string;
  cost: number;
  status: 'received' | 'in_repair' | 'completed' | 'cancelled';
  ui_status?: 'received' | 'in_repair' | 'ready_for_pickup' | 'completed' | 'cancelled';
  device_position: 'warehouse' | 'picked_up';
  payment_status: 'none' | 'requesting' | 'paid';
  payment_method?: 'cash' | 'qris';
  warranty_duration_days: number;
  picked_up_at?: string;
  warranty_expiry_date?: string;
  created_at: string;
}

export interface MockProduct {
  id: string;
  name: string;
  category: 'retail' | 'spare_part';
  stock: number;
  price: number;
  cost_price: number;
  min_stock: number;
}

export interface MockSaleItem {
  productId: string;
  name: string;
  price: number;
  qty: number;
}

export interface MockSale {
  id: string;
  invoice_number: string;
  items: MockSaleItem[];
  subtotal: number;
  discount: number;
  total: number;
  payment_method: 'cash' | 'qris';
  created_at: string;
}

export interface MockSaleCreateItem {
  productId: string;
  qty: number;
}

export interface MockSaleCreate {
  items: MockSaleCreateItem[];
  discount: number;
  payment_method: 'cash' | 'qris';
}

export interface MockWarranty {
  id: string;
  ticket_id: string;
  ticket_number: string;
  customer_name: string;
  device_info: string;
  start_date: string;
  end_date: string;
  status: 'active' | 'expired';
}

export interface MockUser {
  id: string;
  email: string;
  role: 'admin' | 'user';
  passwordHash: string;
}

export interface UserSession {
  email: string;
  role: 'admin' | 'user';
  user_id: string;
}
