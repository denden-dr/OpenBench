import { describe, it, expect, vi, beforeEach } from 'vitest'
import { ticketService } from '../ticketService'
import api from '@/lib/api'
import type { 
  TicketListItem, 
  TicketDetail, 
  CreateTicketRequest, 
  TicketStatus 
} from '@/types/ticket'

vi.mock('@/lib/api', () => ({
  default: {
    get: vi.fn(),
    post: vi.fn(),
    patch: vi.fn(),
    put: vi.fn(),
  },
}))

describe('ticketService', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  describe('getTickets', () => {
    it('fetches list of tickets with parameters', async () => {
      const mockTickets: TicketListItem[] = [
        {
          ticket_id: 't-1',
          ticket_number: 'TKT-001',
          status: 'RECEIVED',
          customer_name: 'Budi',
          device_brand: 'Samsung',
          device_model: 'S23',
          created_at: '2026-07-22T10:00:00Z',
        },
      ]
      const mockResponse = { data: { data: mockTickets, meta: { limit: 10 } } }
      vi.mocked(api.get).mockResolvedValueOnce(mockResponse)

      const result = await ticketService.getTickets({ search: 'Budi', status: 'RECEIVED' })

      expect(api.get).toHaveBeenCalledWith('/admin/services', {
        params: { search: 'Budi', status: 'RECEIVED' },
      })
      expect(result).toEqual(mockResponse.data)
    })
  })

  describe('getTicketByID', () => {
    it('fetches single ticket details by ID', async () => {
      const mockTicketDetail: TicketDetail = {
        ticket_id: 't-1',
        ticket_number: 'TKT-001',
        status: 'RECEIVED',
        customer_name: 'Budi',
        customer_phone: '0812345',
        device_brand: 'Samsung',
        device_model: 'S23',
        issue_description: 'Screen broken',
        cost: 1000000,
        warranty_days: 30,
        created_at: '2026-07-22T10:00:00Z',
      }
      vi.mocked(api.get).mockResolvedValueOnce({ data: { data: mockTicketDetail } })

      const result = await ticketService.getTicketByID('t-1')

      expect(api.get).toHaveBeenCalledWith('/admin/services/t-1')
      expect(result).toEqual(mockTicketDetail)
    })
  })

  describe('createTicket', () => {
    it('sends POST request to create ticket', async () => {
      const createReq: CreateTicketRequest = {
        customer_name: 'Jane',
        customer_phone: '08999',
        device_brand: 'Apple',
        device_model: 'iPhone 13',
        issue_description: 'Battery drop',
        repair_action: 'Replace battery',
        cost: 500000,
        warranty_days: 14,
      }
      const createdDetail: TicketDetail = {
        ticket_id: 't-2',
        ticket_number: 'TKT-002',
        status: 'RECEIVED',
        ...createReq,
        created_at: '2026-07-22T10:00:00Z',
      }
      vi.mocked(api.post).mockResolvedValueOnce({ data: { data: createdDetail } })

      const result = await ticketService.createTicket(createReq)

      expect(api.post).toHaveBeenCalledWith('/admin/services', createReq)
      expect(result).toEqual(createdDetail)
    })
  })

  describe('updateTicketStatus', () => {
    it('sends PATCH request to update ticket status', async () => {
      const updatedStatus: TicketStatus = 'REPAIRING'
      const mockRes = { ticket_id: 't-1', status: updatedStatus, updated_at: '2026-07-22T10:30:00Z' }
      vi.mocked(api.patch).mockResolvedValueOnce({ data: { data: mockRes } })

      const result = await ticketService.updateTicketStatus('t-1', updatedStatus)

      expect(api.patch).toHaveBeenCalledWith('/admin/services/t-1/status', { status: updatedStatus })
      expect(result).toEqual(mockRes)
    })
  })

  describe('updateTicketDetails', () => {
    it('sends PUT request to update ticket details', async () => {
      const updateData = { cost: 600000, notes: 'Discount applied' }
      const updatedDetail: TicketDetail = {
        ticket_id: 't-1',
        status: 'RECEIVED',
        customer_name: 'Budi',
        customer_phone: '0812345',
        device_brand: 'Samsung',
        device_model: 'S23',
        issue_description: 'Screen broken',
        cost: 600000,
        warranty_days: 30,
        notes: 'Discount applied',
        created_at: '2026-07-22T10:00:00Z',
      }
      vi.mocked(api.put).mockResolvedValueOnce({ data: { data: updatedDetail } })

      const result = await ticketService.updateTicketDetails('t-1', updateData)

      expect(api.put).toHaveBeenCalledWith('/admin/services/t-1', updateData)
      expect(result).toEqual(updatedDetail)
    })
  })
})
