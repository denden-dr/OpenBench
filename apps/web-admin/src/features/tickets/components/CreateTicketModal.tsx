import React, { useState } from 'react'
import { ticketService } from '@/services/ticketService'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { 
  Dialog, 
  DialogContent, 
  DialogDescription, 
  DialogHeader, 
  DialogTitle, 
  DialogFooter
} from '@/components/ui/dialog'
import { Loader2 } from 'lucide-react'

interface CreateTicketModalProps {
  isOpen: boolean
  onOpenChange: (open: boolean) => void
  onSuccess: () => void
}

export function CreateTicketModal({ isOpen, onOpenChange, onSuccess }: CreateTicketModalProps) {
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [newTicket, setNewTicket] = useState({
    customer_name: '',
    customer_phone: '',
    device_brand: '',
    device_model: '',
    device_passcode: '',
    issue_description: '',
    repair_action: '',
    cost: 0,
    warranty_days: 30
  })

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setIsSubmitting(true)
    try {
      await ticketService.createTicket({
        customer_name: newTicket.customer_name,
        customer_phone: newTicket.customer_phone,
        device_brand: newTicket.device_brand,
        device_model: newTicket.device_model,
        device_passcode: newTicket.device_passcode || undefined,
        issue_description: newTicket.issue_description,
        repair_action: newTicket.repair_action,
        cost: Number(newTicket.cost),
        warranty_days: Number(newTicket.warranty_days)
      })

      onOpenChange(false)
      setNewTicket({
        customer_name: '',
        customer_phone: '',
        device_brand: '',
        device_model: '',
        device_passcode: '',
        issue_description: '',
        repair_action: '',
        cost: 0,
        warranty_days: 30
      })
      onSuccess()
    } catch (err: unknown) {
      console.error('Failed to create ticket:', err)
      alert('Gagal membuat tiket servis. Periksa kembali input Anda.')
    } finally {
      setIsSubmitting(false)
    }
  }

  return (
    <Dialog open={isOpen} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-xl max-h-[90vh] overflow-y-auto bg-white dark:bg-slate-900 border-slate-200 dark:border-slate-800 text-slate-900 dark:text-slate-100">
        <form onSubmit={handleSubmit}>
          <DialogHeader>
            <DialogTitle className="text-xl font-extrabold text-slate-900 dark:text-slate-100">Create Service Ticket</DialogTitle>
            <DialogDescription className="text-slate-500 dark:text-slate-400">
              Enter customer information, device details, and diagnostic notes to create a new ticket.
            </DialogDescription>
          </DialogHeader>

          <div className="grid grid-cols-1 sm:grid-cols-2 gap-4 py-4">
            {/* Customer Details */}
            <div className="sm:col-span-2 border-b border-slate-100 dark:border-slate-800 pb-2">
              <h4 className="text-sm font-bold text-slate-800 dark:text-slate-200">Customer Information</h4>
            </div>
            
            <div className="space-y-1">
              <label className="text-xs font-bold text-slate-500 dark:text-slate-400 uppercase">Customer Name</label>
              <Input 
                required 
                placeholder="e.g. John Doe"
                value={newTicket.customer_name} 
                onChange={e => setNewTicket({...newTicket, customer_name: e.target.value})} 
              />
            </div>
            
            <div className="space-y-1">
              <label className="text-xs font-bold text-slate-500 dark:text-slate-400 uppercase">Customer Phone</label>
              <Input 
                required 
                placeholder="e.g. 08123456789"
                value={newTicket.customer_phone} 
                onChange={e => setNewTicket({...newTicket, customer_phone: e.target.value})} 
              />
            </div>

            {/* Device Details */}
            <div className="sm:col-span-2 border-b border-slate-100 dark:border-slate-800 pt-2 pb-2">
              <h4 className="text-sm font-bold text-slate-800 dark:text-slate-200">Device Details</h4>
            </div>

            <div className="space-y-1">
              <label className="text-xs font-bold text-slate-500 dark:text-slate-400 uppercase">Brand</label>
              <Input 
                required 
                placeholder="e.g. Apple"
                value={newTicket.device_brand} 
                onChange={e => setNewTicket({...newTicket, device_brand: e.target.value})} 
              />
            </div>

            <div className="space-y-1">
              <label className="text-xs font-bold text-slate-500 dark:text-slate-400 uppercase">Model</label>
              <Input 
                required 
                placeholder="e.g. iPhone 13 Pro"
                value={newTicket.device_model} 
                onChange={e => setNewTicket({...newTicket, device_model: e.target.value})} 
              />
            </div>

            <div className="sm:col-span-2 space-y-1">
              <label className="text-xs font-bold text-slate-500 dark:text-slate-400 uppercase">Screen Passcode (Optional)</label>
              <Input 
                placeholder="e.g. pattern Letter-L, PIN 1234"
                value={newTicket.device_passcode} 
                onChange={e => setNewTicket({...newTicket, device_passcode: e.target.value})} 
              />
            </div>

            {/* Diagnosis Details */}
            <div className="sm:col-span-2 border-b border-slate-100 dark:border-slate-800 pt-2 pb-2">
              <h4 className="text-sm font-bold text-slate-800 dark:text-slate-200">Initial Diagnostic & Pricing</h4>
            </div>

            <div className="sm:col-span-2 space-y-1">
              <label className="text-xs font-bold text-slate-500 dark:text-slate-400 uppercase">Issue Description</label>
              <Input 
                required 
                placeholder="e.g. Cracked LCD, Touch unresponsive"
                value={newTicket.issue_description} 
                onChange={e => setNewTicket({...newTicket, issue_description: e.target.value})} 
              />
            </div>

            <div className="sm:col-span-2 space-y-1">
              <label className="text-xs font-bold text-slate-500 dark:text-slate-400 uppercase">Proposed Repair Action</label>
              <Input 
                required 
                placeholder="e.g. Replacement LCD OLED Screen"
                value={newTicket.repair_action} 
                onChange={e => setNewTicket({...newTicket, repair_action: e.target.value})} 
              />
            </div>

            <div className="space-y-1">
              <label className="text-xs font-bold text-slate-500 dark:text-slate-400 uppercase">Estimated Cost (Rp)</label>
              <Input 
                required 
                type="number"
                value={newTicket.cost} 
                onChange={e => setNewTicket({...newTicket, cost: parseInt(e.target.value) || 0})} 
              />
            </div>

            <div className="space-y-1">
              <label className="text-xs font-bold text-slate-500 dark:text-slate-400 uppercase">Warranty Period (Days)</label>
              <Input 
                required 
                type="number"
                value={newTicket.warranty_days} 
                onChange={e => setNewTicket({...newTicket, warranty_days: parseInt(e.target.value) || 0})} 
              />
            </div>
          </div>

          <DialogFooter className="gap-2 sm:gap-0 pt-2">
            <Button type="button" variant="outline" disabled={isSubmitting} className="cursor-pointer border-slate-200 dark:border-slate-700" onClick={() => onOpenChange(false)}>Cancel</Button>
            <Button type="submit" disabled={isSubmitting} className="bg-primary hover:bg-secondary cursor-pointer">
              {isSubmitting ? <Loader2 className="w-4 h-4 animate-spin mr-1" /> : null}
              Create Ticket
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  )
}
