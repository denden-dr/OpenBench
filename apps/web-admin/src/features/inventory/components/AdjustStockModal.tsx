import { useState, useEffect } from 'react'
import type { Product, AdjustStockRequest } from '@/types/pos'
import { inventoryService } from '@/services/inventoryService'
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

interface AdjustStockModalProps {
  isOpen: boolean
  onOpenChange: (open: boolean) => void
  product: Product | null
  onSuccess: (product: Product) => void
}

export function AdjustStockModal({ isOpen, onOpenChange, product, onSuccess }: AdjustStockModalProps) {
  const [adjustQty, setAdjustQty] = useState(0)
  const [loading, setLoading] = useState(false)

  useEffect(() => {
    if (isOpen) {
      setAdjustQty(0)
    }
  }, [isOpen])

  const handleAdjustSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!product) return
    setLoading(true)
    try {
      const req: AdjustStockRequest = { quantity_change: adjustQty }
      const updated = await inventoryService.adjustStock(product.id, req)
      onSuccess(updated)
      onOpenChange(false)
    } catch (err: any) {
      alert(err?.response?.data?.detail || 'Failed to adjust stock')
    } finally {
      setLoading(false)
    }
  }

  return (
    <Dialog open={isOpen} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-sm">
        <form onSubmit={handleAdjustSubmit}>
          <DialogHeader>
            <DialogTitle className="text-xl font-extrabold text-slate-900 dark:text-slate-100">Adjust Stock Count</DialogTitle>
            <DialogDescription className="text-slate-500 dark:text-slate-400">
              Modify stock count for <span className="font-bold text-slate-900 dark:text-slate-100">{product?.name}</span>.
            </DialogDescription>
          </DialogHeader>

          <div className="space-y-4 py-4">
            <div className="flex items-center justify-between text-slate-600 dark:text-slate-400 text-sm font-semibold">
              <span>Current Stock:</span>
              <span className="font-bold text-slate-900 dark:text-slate-100">{product?.stock}</span>
            </div>

            <div className="space-y-1">
              <label className="text-xs font-bold text-slate-500 dark:text-slate-400 uppercase">Quantity Change</label>
              <Input 
                required 
                type="number"
                placeholder="e.g. 5 or -3"
                value={adjustQty}
                onChange={e => setAdjustQty(parseInt(e.target.value) || 0)}
              />
              <p className="text-xxs text-slate-400 dark:text-slate-500">
                Use positive numbers to add to stock, negative numbers to decrease stock.
              </p>
            </div>
          </div>

          <DialogFooter className="gap-2 sm:gap-0 pt-2">
            <Button type="button" variant="outline" className="cursor-pointer" onClick={() => onOpenChange(false)}>Cancel</Button>
            <Button type="submit" disabled={loading} className="bg-primary hover:bg-secondary cursor-pointer">
              {loading ? 'Applying...' : 'Apply Adjustment'}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  )
}
