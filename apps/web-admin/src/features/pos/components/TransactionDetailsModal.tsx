import { 
  Dialog, 
  DialogContent, 
  DialogDescription, 
  DialogHeader, 
  DialogTitle, 
  DialogFooter
} from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import { ShieldCheck } from 'lucide-react'
import type { POSTransaction } from '@/types/pos'

interface TransactionDetailsModalProps {
  transaction: POSTransaction | null
  isOpen: boolean
  onClose: () => void
}

export function TransactionDetailsModal({
  transaction,
  isOpen,
  onClose,
}: TransactionDetailsModalProps) {
  const formatPrice = (price: number) => {
    return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(price)
  }

  const formatDate = (isoString: string) => {
    return new Date(isoString).toLocaleDateString('id-ID', {
      day: 'numeric',
      month: 'short',
      year: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    })
  }

  return (
    <Dialog open={isOpen} onOpenChange={(open) => {
      if(!open) onClose()
    }}>
      <DialogContent className="max-w-md">
        <DialogHeader>
          <DialogTitle className="text-xl font-extrabold text-slate-900 dark:text-slate-100">Transaction Details</DialogTitle>
          <DialogDescription className="text-slate-500 dark:text-slate-400">
            Purchased items and billing summary for transaction <span className="font-mono font-bold text-slate-900 dark:text-slate-100">{transaction?.id}</span>.
          </DialogDescription>
        </DialogHeader>

        <div className="space-y-4 py-4">
          <div className="flex justify-between items-center text-xs font-semibold text-slate-500 dark:text-slate-400 border-b border-slate-100 dark:border-slate-800 pb-2">
            <span>Date: {transaction && formatDate(transaction.created_at)}</span>
            <span>Method: <span className="font-bold text-slate-950 dark:text-slate-100">{transaction?.payment_method}</span></span>
          </div>

          <div className="space-y-3">
            <span className="text-xxs font-bold text-slate-400 dark:text-slate-500 uppercase tracking-wider">Items Purchased</span>
            <div className="space-y-2.5 max-h-48 overflow-y-auto">
              {transaction?.items.map((item) => (
                <div key={item.id} className="flex justify-between text-slate-700 dark:text-slate-300 text-xs">
                  <span>{item.product_name} <span className="text-slate-400 dark:text-slate-500 font-medium">x {item.quantity}</span></span>
                  <span className="font-mono font-bold">{formatPrice(item.price * item.quantity)}</span>
                </div>
              ))}
            </div>
          </div>

          <div className="border-t border-slate-100 dark:border-slate-800 pt-3 flex justify-between items-center">
            <span className="text-xs font-bold text-slate-500 dark:text-slate-400">Total Paid:</span>
            <span className="font-mono text-base font-extrabold text-slate-900 dark:text-slate-100">{transaction && formatPrice(transaction.total_amount)}</span>
          </div>
        </div>

        <DialogFooter className="pt-2">
          <Button type="button" className="w-full bg-primary hover:bg-secondary cursor-pointer text-xs" onClick={onClose}>
            <ShieldCheck className="w-4 h-4 mr-1.5" />
            Close Invoice
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
