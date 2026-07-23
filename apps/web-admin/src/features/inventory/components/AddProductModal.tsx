import { useState } from 'react'
import type { CreateProductRequest } from '@/types/pos'
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

interface AddProductModalProps {
  isOpen: boolean
  onOpenChange: (open: boolean) => void
  onSuccess: (product: any) => void
}

export function AddProductModal({ isOpen, onOpenChange, onSuccess }: AddProductModalProps) {
  const [newProduct, setNewProduct] = useState({
    name: '',
    price: 0,
    stock: 0
  })
  const [loading, setLoading] = useState(false)

  const handleAddSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setLoading(true)
    try {
      const req: CreateProductRequest = {
        name: newProduct.name,
        price: newProduct.price,
        stock: newProduct.stock,
      }
      const created = await inventoryService.createProduct(req)
      onSuccess(created)
      onOpenChange(false)
      setNewProduct({ name: '', price: 0, stock: 0 })
    } catch (err: any) {
      alert(err?.response?.data?.detail || 'Failed to create product')
    } finally {
      setLoading(false)
    }
  }

  return (
    <Dialog open={isOpen} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-md">
        <form onSubmit={handleAddSubmit}>
          <DialogHeader>
            <DialogTitle className="text-xl font-extrabold text-slate-900 dark:text-slate-100">Add New Product</DialogTitle>
            <DialogDescription className="text-slate-500 dark:text-slate-400">
              Enter details to register a new accessory or product in the inventory.
            </DialogDescription>
          </DialogHeader>

          <div className="space-y-4 py-4">
            <div className="space-y-1">
              <label className="text-xs font-bold text-slate-500 dark:text-slate-400 uppercase">Product Name</label>
              <Input
                required
                placeholder="e.g. Tempered Glass iPhone 15"
                value={newProduct.name}
                onChange={e => setNewProduct({ ...newProduct, name: e.target.value })}
              />
            </div>

            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-1">
                <label className="text-xs font-bold text-slate-500 dark:text-slate-400 uppercase">Price (Rp)</label>
                <Input
                  required
                  type="number"
                  value={newProduct.price}
                  onChange={e => setNewProduct({ ...newProduct, price: parseInt(e.target.value) || 0 })}
                />
              </div>

              <div className="space-y-1">
                <label className="text-xs font-bold text-slate-500 dark:text-slate-400 uppercase">Initial Stock</label>
                <Input
                  required
                  type="number"
                  value={newProduct.stock}
                  onChange={e => setNewProduct({ ...newProduct, stock: parseInt(e.target.value) || 0 })}
                />
              </div>
            </div>
          </div>

          <DialogFooter className="gap-2 sm:gap-0 pt-2">
            <Button type="button" variant="outline" className="cursor-pointer" onClick={() => onOpenChange(false)}>Cancel</Button>
            <Button type="submit" disabled={loading} className="bg-primary hover:bg-secondary cursor-pointer">
              {loading ? 'Saving...' : 'Save Product'}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  )
}
