import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import TicketsPage from '../TicketsPage'
import { ticketService } from '@/services/ticketService'
import type { TicketListItem, TicketDetail } from '@/types/ticket'

vi.mock('@/services/ticketService', () => ({
  ticketService: {
    getTickets: vi.fn(),
    getTicketByID: vi.fn(),
    createTicket: vi.fn(),
    updateTicketStatus: vi.fn(),
    updateTicketDetails: vi.fn(),
  },
}))

describe('TicketsPage', () => {
  const mockTickets: TicketListItem[] = [
    {
      ticket_id: 't-1',
      ticket_number: 'TKT-20260722-0001',
      status: 'RECEIVED',
      customer_name: 'Budi Santoso',
      device_brand: 'Samsung',
      device_model: 'Galaxy S23',
      created_at: '2026-07-22T10:00:00Z',
    },
    {
      ticket_id: 't-2',
      ticket_number: 'TKT-20260722-0002',
      status: 'COMPLETED',
      customer_name: 'Siti Rahma',
      device_brand: 'Apple',
      device_model: 'iPhone 13',
      created_at: '2026-07-22T11:00:00Z',
    },
  ]

  const mockDetail: TicketDetail = {
    ticket_id: 't-1',
    ticket_number: 'TKT-20260722-0001',
    status: 'RECEIVED',
    customer_name: 'Budi Santoso',
    customer_phone: '08123456789',
    device_brand: 'Samsung',
    device_model: 'Galaxy S23',
    device_passcode: '1234',
    issue_description: 'Layar Pecah',
    repair_action: 'Ganti LCD',
    cost: 1500000,
    warranty_days: 30,
    created_at: '2026-07-22T10:00:00Z',
  }

  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('renders loading state initially and fetches tickets', async () => {
    vi.mocked(ticketService.getTickets).mockResolvedValueOnce({ data: mockTickets })

    render(<TicketsPage />)

    expect(screen.getByText('Memuat daftar tiket servis...')).toBeInTheDocument()

    await waitFor(() => {
      expect(screen.getByText('TKT-20260722-0001')).toBeInTheDocument()
      expect(screen.getByText('Budi Santoso')).toBeInTheDocument()
      expect(screen.getByText('Siti Rahma')).toBeInTheDocument()
    })

    expect(ticketService.getTickets).toHaveBeenCalled()
  })

  it('displays error message when fetching tickets fails', async () => {
    vi.mocked(ticketService.getTickets).mockRejectedValueOnce(new Error('Network error'))

    render(<TicketsPage />)

    await waitFor(() => {
      expect(screen.getByText('Gagal memuat data tiket servis. Silakan coba lagi.')).toBeInTheDocument()
    })
  })

  it('opens view details dialog when Eye button is clicked', async () => {
    vi.mocked(ticketService.getTickets).mockResolvedValueOnce({ data: mockTickets })
    vi.mocked(ticketService.getTicketByID).mockResolvedValueOnce(mockDetail)

    render(<TicketsPage />)

    await waitFor(() => {
      expect(screen.getByText('TKT-20260722-0001')).toBeInTheDocument()
    })

    const viewButtons = screen.getAllByTitle('View details')
    fireEvent.click(viewButtons[0])

    await waitFor(() => {
      expect(screen.getByText('Ticket Details')).toBeInTheDocument()
      expect(screen.getByText('Layar Pecah')).toBeInTheDocument()
      expect(screen.getByText('Ganti LCD')).toBeInTheDocument()
      expect(screen.getByText('1234')).toBeInTheDocument()
    })

    expect(ticketService.getTicketByID).toHaveBeenCalledWith('t-1')
  })

  it('opens update status dialog and submits status change', async () => {
    vi.mocked(ticketService.getTickets).mockResolvedValue({ data: mockTickets })
    vi.mocked(ticketService.updateTicketStatus).mockResolvedValueOnce({
      ticket_id: 't-1',
      status: 'REPAIRING',
      updated_at: '2026-07-22T12:00:00Z',
    })

    render(<TicketsPage />)

    await waitFor(() => {
      expect(screen.getByText('TKT-20260722-0001')).toBeInTheDocument()
    })

    const editButtons = screen.getAllByTitle('Update status')
    fireEvent.click(editButtons[0])

    await waitFor(() => {
      expect(screen.getByText('Update Ticket Status')).toBeInTheDocument()
    })

    const select = screen.getByRole('combobox')
    fireEvent.change(select, { target: { value: 'REPAIRING' } })

    const submitBtn = screen.getByRole('button', { name: /Update Status/i })
    fireEvent.click(submitBtn)

    await waitFor(() => {
      expect(ticketService.updateTicketStatus).toHaveBeenCalledWith('t-1', 'REPAIRING')
    })
  })
})
