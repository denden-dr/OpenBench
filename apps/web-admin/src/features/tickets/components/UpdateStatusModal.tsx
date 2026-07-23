import { useState, useEffect } from 'react'
import type { TicketListItem, TicketStatus } from '@/types/ticket'
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

interface UpdateStatusModalProps {
  isOpen: boolean
  onOpenChange: (open: boolean) => void
  ticket: TicketListItem | null
  onSuccess: () => void
}

export function UpdateStatusModal({ isOpen, onOpenChange, ticket, onSuccess }: UpdateStatusModalProps) {
  const [statusSubmitting, setStatusSubmitting] = useState(false)
  const [newStatus, setNewStatus] = useState<TicketStatus>('RECEIVED')

  useEffect(() => {
    if (ticket) {
      setNewStatus(ticket.status)
    }
  }, [ticket])

  const handleStatusSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!ticket) return
    setStatusSubmitting(true)
    try {
      await ticketService.updateTicketStatus(ticket.ticket_id, newStatus)
      onOpenChange(false)
      onSuccess()
    } catch (err) {
      console.error('Failed to update status:', err)
      alert('Gagal memperbarui status tiket.')
    } finally {
      setStatusSubmitting(false)
    }
  }

  return (
    <Dialog open={isOpen} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-md bg-white dark:bg-slate-900 border-slate-200 dark:border-slate-800 text-slate-900 dark:text-slate-100">
        <form onSubmit={handleStatusSubmit}>
          <DialogHeader>
            <DialogTitle className="text-xl font-extrabold">Update Ticket Status</DialogTitle>
            <DialogDescription className="text-slate-500 dark:text-slate-400">
              Change progress status for {ticket?.ticket_number}
            </DialogDescription>
          </DialogHeader>

          <div className="py-4 space-y-3">
            <label className="text-xs font-bold text-slate-500 dark:text-slate-400 uppercase">Select Status</label>
            <select
              value={newStatus}
              onChange={(e) => setNewStatus(e.target.value as TicketStatus)}
              className="w-full p-2.5 rounded-md bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 text-slate-900 dark:text-slate-100 font-semibold focus:outline-none focus:ring-2 focus:ring-primary/20"
            >
              <option value="RECEIVED">RECEIVED (Diterima)</option>
              <option value="REPAIRING">REPAIRING (Sedang Dikerjakan)</option>
              <option value="PENDING_CONFIRMATION">PENDING_CONFIRMATION (Butuh Konfirmasi)</option>
              <option value="FIXED">FIXED (Selesai Dikerjakan)</option>
              <option value="COMPLETED">COMPLETED (Diambil/Lunas)</option>
              <option value="CANCELLED">CANCELLED (Dibatalkan)</option>
              <option value="RETURNED">RETURNED (Dikembalikan)</option>
            </select>
          </div>

          <DialogFooter className="gap-2 sm:gap-0 pt-2">
            <Button type="button" variant="outline" disabled={statusSubmitting} onClick={() => onOpenChange(false)}>Cancel</Button>
            <Button type="submit" disabled={statusSubmitting} className="bg-primary hover:bg-secondary">
              {statusSubmitting ? <Loader2 className="w-4 h-4 animate-spin mr-1" /> : null}
              Update Status
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  )
}
