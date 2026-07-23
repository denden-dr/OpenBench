import { describe, it, expect, vi, beforeEach } from 'vitest'
import { warrantyService } from '../warrantyService'
import api from '@/lib/api'
import type { Warranty, WarrantyClaim } from '@/types/warranty'

vi.mock('@/lib/api', () => ({
  default: {
    get: vi.fn(),
    post: vi.fn(),
    patch: vi.fn(),
    put: vi.fn(),
  },
}))

describe('warrantyService', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  describe('getWarrantyByTicketNumber', () => {
    it('fetches warranty by ticket number', async () => {
      const mockWarranty: Warranty = { id: 'w-1', ticket_id: 't-1', start_date: '2026-07-01T00:00:00Z', end_date: '2026-08-01T00:00:00Z', status: 'ACTIVE', notes: null }
      vi.mocked(api.get).mockResolvedValueOnce({ data: { data: mockWarranty } })

      const result = await warrantyService.getWarrantyByTicketNumber('TKT-001')

      expect(api.get).toHaveBeenCalledWith('/admin/warranties/by-ticket-number/TKT-001')
      expect(result).toEqual(mockWarranty)
    })
  })

  describe('updateWarrantyStatus', () => {
    it('sends PATCH to update warranty status', async () => {
      const updated: Warranty = { id: 'w-1', ticket_id: 't-1', start_date: '2026-07-01T00:00:00Z', end_date: '2026-08-01T00:00:00Z', status: 'VOID', notes: 'voided' }
      vi.mocked(api.patch).mockResolvedValueOnce({ data: { data: updated } })

      const result = await warrantyService.updateWarrantyStatus('w-1', { status: 'VOID', notes: 'voided' })

      expect(api.patch).toHaveBeenCalledWith('/admin/warranties/w-1/status', { status: 'VOID', notes: 'voided' })
      expect(result).toEqual(updated)
    })
  })

  describe('createClaim', () => {
    it('sends POST to create a claim', async () => {
      const claim: WarrantyClaim = { claim_id: 'c-1', claim_number: 'CLM-001', warranty_id: 'w-1', warranty_ticket_ref_id: null, evaluation_status: 'PENDING', issue_description: 'Screen issue', notes: null, evaluation_notes: null, created_at: '2026-07-15T00:00:00Z', updated_at: '2026-07-15T00:00:00Z' }
      vi.mocked(api.post).mockResolvedValueOnce({ data: { data: claim } })

      const result = await warrantyService.createClaim({ ticket_number: 'TKT-001', issue_description: 'Screen issue' })

      expect(api.post).toHaveBeenCalledWith('/admin/claims', { ticket_number: 'TKT-001', issue_description: 'Screen issue' })
      expect(result).toEqual(claim)
    })
  })

  describe('getClaims', () => {
    it('fetches claims list with params', async () => {
      const claims: WarrantyClaim[] = [{ claim_id: 'c-1', claim_number: 'CLM-001', warranty_id: 'w-1', warranty_ticket_ref_id: null, evaluation_status: 'PENDING', issue_description: 'Test', notes: null, evaluation_notes: null, created_at: '2026-07-15T00:00:00Z', updated_at: '2026-07-15T00:00:00Z' }]
      const mockResponse = { data: { data: claims, meta: { limit: 10 } } }
      vi.mocked(api.get).mockResolvedValueOnce(mockResponse)

      const result = await warrantyService.getClaims({ status: 'PENDING', limit: 10 })

      expect(api.get).toHaveBeenCalledWith('/admin/claims', { params: { status: 'PENDING', limit: 10 } })
      expect(result).toEqual(mockResponse.data)
    })
  })

  describe('getClaimByID', () => {
    it('fetches single claim by ID', async () => {
      const claim: WarrantyClaim = { claim_id: 'c-1', claim_number: 'CLM-001', warranty_id: 'w-1', warranty_ticket_ref_id: null, evaluation_status: 'PENDING', issue_description: 'Test', notes: null, evaluation_notes: null, created_at: '2026-07-15T00:00:00Z', updated_at: '2026-07-15T00:00:00Z' }
      vi.mocked(api.get).mockResolvedValueOnce({ data: { data: claim } })

      const result = await warrantyService.getClaimByID('c-1')

      expect(api.get).toHaveBeenCalledWith('/admin/claims/c-1')
      expect(result).toEqual(claim)
    })
  })

  describe('updateClaim', () => {
    it('sends PUT to update a claim', async () => {
      const updated: WarrantyClaim = { claim_id: 'c-1', claim_number: 'CLM-001', warranty_id: 'w-1', warranty_ticket_ref_id: null, evaluation_status: 'PENDING', issue_description: 'Updated description', notes: 'updated note', evaluation_notes: null, created_at: '2026-07-15T00:00:00Z', updated_at: '2026-07-16T00:00:00Z' }
      vi.mocked(api.put).mockResolvedValueOnce({ data: { data: updated } })

      const result = await warrantyService.updateClaim('c-1', { issue_description: 'Updated description', notes: 'updated note' })

      expect(api.put).toHaveBeenCalledWith('/admin/claims/c-1', { issue_description: 'Updated description', notes: 'updated note' })
      expect(result).toEqual(updated)
    })
  })

  describe('evaluateClaim', () => {
    it('sends POST to evaluate a claim', async () => {
      const evaluated: WarrantyClaim = { claim_id: 'c-1', claim_number: 'CLM-001', warranty_id: 'w-1', warranty_ticket_ref_id: null, evaluation_status: 'ACCEPTED', issue_description: 'Test', notes: null, evaluation_notes: 'Approved', created_at: '2026-07-15T00:00:00Z', updated_at: '2026-07-16T00:00:00Z' }
      vi.mocked(api.post).mockResolvedValueOnce({ data: { data: evaluated } })

      const result = await warrantyService.evaluateClaim('c-1', { status: 'ACCEPTED', notes: 'Approved' })

      expect(api.post).toHaveBeenCalledWith('/admin/claims/c-1/evaluate', { status: 'ACCEPTED', notes: 'Approved' })
      expect(result).toEqual(evaluated)
    })
  })
})
