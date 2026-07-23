import { useState, useEffect, useCallback } from 'react'
import type { Product } from '@/types/pos'
import { inventoryService } from '@/services/inventoryService'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Search, Plus } from 'lucide-react'
import { AddProductModal } from '@/features/inventory/components/AddProductModal'
import { EditProductModal } from '@/features/inventory/components/EditProductModal'
import { AdjustStockModal } from '@/features/inventory/components/AdjustStockModal'
import { InventoryTable } from '@/features/inventory/components/InventoryTable'

function InventoryPage() {
  const [products, setProducts] = useState<Product[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [searchQuery, setSearchQuery] = useState('')
  const [isAddOpen, setIsAddOpen] = useState(false)
  const [isAdjustOpen, setIsAdjustOpen] = useState(false)
  const [isEditOpen, setIsEditOpen] = useState(false)
  const [selectedProduct, setSelectedProduct] = useState<Product | null>(null)
  const [selectedEditProduct, setSelectedEditProduct] = useState<Product | null>(null)

  const fetchProducts = useCallback(async () => {
    try {
      setError(null)
      const result = await inventoryService.getProducts({ limit: 50 })
      setProducts(result.data)
    } catch (err: any) {
      setError(err?.response?.data?.detail || err?.message || 'Failed to load products')
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    fetchProducts()
  }, [fetchProducts])

  const handleDelete = async (id: string) => {
    if (!confirm('Are you sure you want to delete this product?')) return
    try {
      await inventoryService.deleteProduct(id)
      setProducts(products.filter(p => p.id !== id))
    } catch (err: any) {
      alert(err?.response?.data?.detail || 'Failed to delete product')
    }
  }

  const filteredProducts = products.filter(p => 
    p.name.toLowerCase().includes(searchQuery.toLowerCase())
  )

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <p className="text-slate-500">Loading products...</p>
      </div>
    )
  }

  if (error) {
    return (
      <div className="flex flex-col items-center justify-center h-64 gap-4">
        <p className="text-red-500 font-semibold">{error}</p>
        <Button onClick={fetchProducts} variant="outline">Retry</Button>
      </div>
    )
  }

  return (
    <div className="space-y-8">
      {/* Title */}
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <h1 className="text-3xl font-extrabold text-slate-900 dark:text-slate-100 tracking-tight">
            Product Inventory
          </h1>
          <p className="text-slate-500 dark:text-slate-400 text-sm">Manage products, tracking stock levels, and adjusting counts for shop accessories.</p>
        </div>

        <Button 
          className="font-semibold bg-primary hover:bg-secondary cursor-pointer"
          onClick={() => setIsAddOpen(true)}
        >
          <Plus className="w-4 h-4 mr-1" />
          Add Product
        </Button>
      </div>

      {/* Toolbar */}
      <div className="flex items-center justify-between">
        <div className="relative w-full md:w-80">
          <Search className="absolute left-3 top-2.5 h-4 w-4 text-slate-400 dark:text-slate-500" />
          <Input
            placeholder="Search products by name..."
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            className="pl-9 bg-white dark:bg-slate-900 border-slate-200 dark:border-slate-800 text-slate-900 dark:text-slate-100 focus-visible:ring-primary/20"
          />
        </div>
      </div>

      {/* Inventory Table */}
      <InventoryTable 
        products={filteredProducts}
        onEdit={(p) => {
          setSelectedEditProduct(p)
          setIsEditOpen(true)
        }}
        onAdjust={(p) => {
          setSelectedProduct(p)
          setIsAdjustOpen(true)
        }}
        onDelete={handleDelete}
      />

      {/* Modals */}
      <AddProductModal 
        isOpen={isAddOpen}
        onOpenChange={setIsAddOpen}
        onSuccess={(created) => setProducts([created, ...products])}
      />

      <EditProductModal 
        isOpen={isEditOpen}
        onOpenChange={setIsEditOpen}
        product={selectedEditProduct}
        onSuccess={(updated) => setProducts(products.map(p => p.id === updated.id ? updated : p))}
      />

      <AdjustStockModal 
        isOpen={isAdjustOpen}
        onOpenChange={setIsAdjustOpen}
        product={selectedProduct}
        onSuccess={(updated) => setProducts(products.map(p => p.id === updated.id ? updated : p))}
      />
    </div>
  )
}

export default InventoryPage
