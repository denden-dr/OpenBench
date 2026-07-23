import { useState, useEffect } from 'react'
import type { TicketDetail } from '@/types/ticket'
import { ticketService } from '@/services/ticketService'
import { Button } from '@/components/ui/button'
import { 
  Dialog, 
  DialogContent, 
  DialogDescription, 
  DialogHeader, 
  DialogTitle, 
  DialogFooter
} from '@/components/ui/dialog'
import { Loader2 } from 'lucide-react'
import { TicketStatusBadge } from './TicketStatusBadge'

interface TicketDetailsModalProps {
  isOpen: boolean
  onOpenChange: (open: boolean) => void
  ticketId: string | null
}

export function TicketDetailsModal({ isOpen, onOpenChange, ticketId }: TicketDetailsModalProps) {
  const [viewLoading, setViewLoading] = useState(false)
  const [selectedTicketDetail, setSelectedTicketDetail] = useState<TicketDetail | null>(null)

  useEffect(() => {
    if (!isOpen || !ticketId) return

    const fetchDetail = async () => {
      setViewLoading(true)
      setSelectedTicketDetail(null)
      try {
        const detail = await ticketService.getTicketByID(ticketId)
        setSelectedTicketDetail(detail)
      } catch (err) {
        console.error('Failed to fetch ticket detail:', err)
        alert('Gagal memuat detail tiket.')
        onOpenChange(false)
      } finally {
        setViewLoading(false)
      }
    }
    fetchDetail()
  }, [isOpen, ticketId, onOpenChange])

  const formatDate = (isoString: string) => {
    if (!isoString) return '-'
    return new Date(isoString).toLocaleDateString('id-ID', {
      day: 'numeric',
      month: 'short',
      year: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    })
  }

  return (
    <Dialog open={isOpen} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-lg max-h-[90vh] overflow-y-auto bg-white dark:bg-slate-900 border-slate-200 dark:border-slate-800 text-slate-900 dark:text-slate-100">
        <DialogHeader>
          <DialogTitle className="text-xl font-extrabold flex items-center justify-between pr-4">
            <span>Ticket Details</span>
            {selectedTicketDetail && <TicketStatusBadge status={selectedTicketDetail.status} />}
          </DialogTitle>
          <DialogDescription className="text-slate-500 dark:text-slate-400">
            {selectedTicketDetail?.ticket_number || 'Loading details...'}
          </DialogDescription>
        </DialogHeader>

        {viewLoading ? (
          <div className="flex flex-col items-center justify-center py-12 gap-2 text-slate-400">
            <Loader2 className="w-6 h-6 animate-spin text-primary" />
            <span>Memuat rincian tiket...</span>
          </div>
        ) : selectedTicketDetail ? (
          <div className="space-y-4 py-2 text-sm">
            {/* Customer */}
            <div className="bg-slate-50 dark:bg-slate-800/50 p-3 rounded-lg space-y-1">
              <div className="text-xs font-bold uppercase text-slate-400">Customer</div>
              <div className="font-bold text-slate-900 dark:text-slate-100">{selectedTicketDetail.customer_name}</div>
              <div className="text-xs text-slate-500">{selectedTicketDetail.customer_phone}</div>
            </div>

            {/* Device */}
            <div className="bg-slate-50 dark:bg-slate-800/50 p-3 rounded-lg space-y-1">
              <div className="text-xs font-bold uppercase text-slate-400">Device</div>
              <div className="font-bold text-slate-900 dark:text-slate-100">
                {selectedTicketDetail.device_brand} {selectedTicketDetail.device_model}
              </div>
              {selectedTicketDetail.device_passcode && (
                <div className="text-xs text-slate-500">Passcode: <span className="font-mono">{selectedTicketDetail.device_passcode}</span></div>
              )}
            </div>

            {/* Diagnosis */}
            <div className="space-y-2 border-t border-slate-100 dark:border-slate-800 pt-2">
              <div>
                <div className="text-xs font-bold uppercase text-slate-400">Issue Description</div>
                <div className="text-slate-800 dark:text-slate-200 mt-0.5">{selectedTicketDetail.issue_description}</div>
              </div>
              {selectedTicketDetail.repair_action && (
                <div>
                  <div className="text-xs font-bold uppercase text-slate-400">Repair Action</div>
                  <div className="text-slate-800 dark:text-slate-200 mt-0.5">{selectedTicketDetail.repair_action}</div>
                </div>
              )}
            </div>

            {/* Pricing & Warranty */}
            <div className="grid grid-cols-2 gap-3 border-t border-slate-100 dark:border-slate-800 pt-2">
              <div>
                <div className="text-xs font-bold uppercase text-slate-400">Cost</div>
                <div className="font-bold text-slate-900 dark:text-slate-100 text-base">
                  Rp {selectedTicketDetail.cost.toLocaleString('id-ID')}
                </div>
              </div>
              <div>
                <div className="text-xs font-bold uppercase text-slate-400">Warranty</div>
                <div className="font-bold text-slate-900 dark:text-slate-100 text-base">
                  {selectedTicketDetail.warranty_days} Days
                </div>
              </div>
            </div>

            {/* Notes */}
            {selectedTicketDetail.notes && (
              <div className="border-t border-slate-100 dark:border-slate-800 pt-2">
                <div className="text-xs font-bold uppercase text-slate-400">Notes</div>
                <div className="text-xs text-slate-600 dark:text-slate-300 mt-1 italic">{selectedTicketDetail.notes}</div>
              </div>
            )}

            {/* Timestamps */}
            <div className="text-xxs text-slate-400 border-t border-slate-100 dark:border-slate-800 pt-2 flex justify-between">
              <span>Created: {formatDate(selectedTicketDetail.created_at)}</span>
              {selectedTicketDetail.updated_at && <span>Updated: {formatDate(selectedTicketDetail.updated_at)}</span>}
            </div>
          </div>
        ) : null}

        <DialogFooter>
          <Button variant="outline" onClick={() => onOpenChange(false)}>Close</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
