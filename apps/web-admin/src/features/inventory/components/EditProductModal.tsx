import { useState, useEffect } from 'react'
import type { Product, UpdateProductRequest } from '@/types/pos'
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

interface EditProductModalProps {
  isOpen: boolean
  onOpenChange: (open: boolean) => void
  product: Product | null
  onSuccess: (product: Product) => void
}

export function EditProductModal({ isOpen, onOpenChange, product, onSuccess }: EditProductModalProps) {
  const [editProduct, setEditProduct] = useState({
    name: '',
    price: 0
  })
  const [loading, setLoading] = useState(false)

  useEffect(() => {
    if (product) {
      setEditProduct({
        name: product.name,
        price: product.price
      })
    }
  }, [product])

  const handleEditSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!product) return
    setLoading(true)
    try {
      const req: UpdateProductRequest = {
        name: editProduct.name,
        price: editProduct.price,
        stock: product.stock,
      }
      const updated = await inventoryService.updateProduct(product.id, req)
      onSuccess(updated)
      onOpenChange(false)
    } catch (err: any) {
      alert(err?.response?.data?.detail || 'Failed to update product')
    } finally {
      setLoading(false)
    }
  }

  return (
    <Dialog open={isOpen} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-md">
        <form onSubmit={handleEditSubmit}>
          <DialogHeader>
            <DialogTitle className="text-xl font-extrabold text-slate-900 dark:text-slate-100">Edit Product Details</DialogTitle>
            <DialogDescription className="text-slate-500 dark:text-slate-400">
              Update basic details for product <span className="font-mono font-bold text-slate-900 dark:text-slate-100">{product?.id}</span>.
            </DialogDescription>
          </DialogHeader>

          <div className="space-y-4 py-4">
            <div className="space-y-1">
              <label className="text-xs font-bold text-slate-500 dark:text-slate-400 uppercase">Product Name</label>
              <Input 
                required 
                value={editProduct.name} 
                onChange={e => setEditProduct({...editProduct, name: e.target.value})} 
              />
            </div>

            <div className="space-y-1">
              <label className="text-xs font-bold text-slate-500 dark:text-slate-400 uppercase">Price (Rp)</label>
              <Input 
                required 
                type="number"
                value={editProduct.price} 
                onChange={e => setEditProduct({...editProduct, price: parseInt(e.target.value) || 0})} 
              />
            </div>
          </div>

          <DialogFooter className="gap-2 sm:gap-0 pt-2">
            <Button type="button" variant="outline" className="cursor-pointer" onClick={() => onOpenChange(false)}>Cancel</Button>
            <Button type="submit" disabled={loading} className="bg-primary hover:bg-secondary cursor-pointer">
              {loading ? 'Saving...' : 'Save Changes'}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  )
}
