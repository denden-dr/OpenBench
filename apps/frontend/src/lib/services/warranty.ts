import { env } from '$env/dynamic/public';
import { apiFetch as fetch } from './api';

const getApiUrl = () => {
  try {
    return env.PUBLIC_API_URL || '';
  } catch {
    return '';
  }
};

export interface Warranty {
  id: string;
  ticket_id: string;
  ticket_number: string;
  customer_name: string;
  device_info: string;
  start_date: string;
  end_date: string;
  status: 'active' | 'expired';
}

export const warrantyService = {
  async getWarranties(): Promise<Warranty[]> {
    const res = await fetch(`${getApiUrl()}/api/v1/admin/warranties`, { credentials: 'include' });
    const body = await res.json();
    if (!res.ok) throw new Error(body.message || 'Failed to fetch warranties');
    return body.data ?? [];
  }
};
