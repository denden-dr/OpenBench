import { useState } from 'react'
import type { Product } from '@/types/pos'
import { Card, CardContent } from '@/components/ui/card'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { 
  Dialog, 
  DialogContent, 
  DialogDescription, 
  DialogHeader, 
  DialogTitle, 
  DialogTrigger, 
  DialogFooter
} from '@/components/ui/dialog'
import { Search, Plus, Trash2, SlidersHorizontal, Pencil, ChevronLeft, ChevronRight } from 'lucide-react'

const initialProducts: Product[] = [
  { id: 'p-1', name: 'Tempered Glass iPhone 15 Pro Max', price: 75000, stock: 18, created_at: '2026-07-15T09:00:00Z' },
  { id: 'p-2', name: 'Silicon Case iPhone 15', price: 120000, stock: 4, created_at: '2026-07-15T09:15:00Z' },
  { id: 'p-3', name: 'USB-C Charger Adapter 20W', price: 299000, stock: 2, created_at: '2026-07-16T11:00:00Z' },
  { id: 'p-4', name: 'MicroUSB Cable 1m', price: 35000, stock: 25, created_at: '2026-07-16T11:30:00Z' },
  { id: 'p-5', name: 'Lightning to USB-C Cable 2m', price: 150000, stock: 0, created_at: '2026-07-17T14:00:00Z' }
]

function InventoryPage() {
  const [products, setProducts] = useState<Product[]>(initialProducts)
  const [searchQuery, setSearchQuery] = useState('')
  const [isAddOpen, setIsAddOpen] = useState(false)
  const [isAdjustOpen, setIsAdjustOpen] = useState(false)
  const [isEditOpen, setIsEditOpen] = useState(false)
  const [selectedProduct, setSelectedProduct] = useState<Product | null>(null)
  const [selectedEditProduct, setSelectedEditProduct] = useState<Product | null>(null)
  
  // Add Form State
  const [newProduct, setNewProduct] = useState({
    name: '',
    price: 0,
    stock: 0
  })

  // Edit Form State
  const [editProduct, setEditProduct] = useState({
    name: '',
    price: 0
  })

  // Adjust Form State
  const [adjustQty, setAdjustQty] = useState(0)

  const handleAddSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    const productItem: Product = {
      id: `p-${Math.random().toString(36).substr(2, 9)}`,
      name: newProduct.name,
      price: newProduct.price,
      stock: newProduct.stock,
      created_at: new Date().toISOString()
    }
    setProducts([productItem, ...products])
    setIsAddOpen(false)
    setNewProduct({ name: '', price: 0, stock: 0 })
  }

  const handleEditSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    if (!selectedEditProduct) return
    setProducts(products.map(p => {
      if (p.id === selectedEditProduct.id) {
        return {
          ...p,
          name: editProduct.name,
          price: editProduct.price
        }
      }
      return p
    }))
    setIsEditOpen(false)
    setSelectedEditProduct(null)
  }

  const handleAdjustSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    if (!selectedProduct) return
    setProducts(products.map(p => {
      if (p.id === selectedProduct.id) {
        return {
          ...p,
          stock: Math.max(0, p.stock + adjustQty)
        }
      }
      return p
    }))
    setIsAdjustOpen(false)
    setSelectedProduct(null)
    setAdjustQty(0)
  }

  const handleDelete = (id: string) => {
    if (confirm('Are you sure you want to delete this product?')) {
      setProducts(products.filter(p => p.id !== id))
    }
  }

  const getStockBadge = (stock: number) => {
    if (stock === 0) {
      return <Badge className="bg-red-500/10 text-red-600 border-none font-semibold">Out of Stock</Badge>
    }
    if (stock < 5) {
      return <Badge className="bg-orange-500/10 text-orange-600 border-none font-semibold">Low Stock ({stock})</Badge>
    }
    return <Badge className="bg-green-500/10 text-green-600 border-none font-semibold">In Stock ({stock})</Badge>
  }

  const filteredProducts = products.filter(p => 
    p.name.toLowerCase().includes(searchQuery.toLowerCase())
  )

  const formatPrice = (price: number) => {
    return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(price)
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

        {/* Add Product Dialog */}
        <Dialog open={isAddOpen} onOpenChange={setIsAddOpen}>
          <DialogTrigger render={
            <Button className="font-semibold bg-primary hover:bg-secondary cursor-pointer">
              <Plus className="w-4 h-4 mr-1" />
              Add Product
            </Button>
          } />
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
                    onChange={e => setNewProduct({...newProduct, name: e.target.value})} 
                  />
                </div>
                
                <div className="grid grid-cols-2 gap-4">
                  <div className="space-y-1">
                    <label className="text-xs font-bold text-slate-500 dark:text-slate-400 uppercase">Price (Rp)</label>
                    <Input 
                      required 
                      type="number"
                      value={newProduct.price} 
                      onChange={e => setNewProduct({...newProduct, price: parseInt(e.target.value) || 0})} 
                    />
                  </div>
                  
                  <div className="space-y-1">
                    <label className="text-xs font-bold text-slate-500 dark:text-slate-400 uppercase">Initial Stock</label>
                    <Input 
                      required 
                      type="number"
                      value={newProduct.stock} 
                      onChange={e => setNewProduct({...newProduct, stock: parseInt(e.target.value) || 0})} 
                    />
                  </div>
                </div>
              </div>

              <DialogFooter className="gap-2 sm:gap-0 pt-2">
                <Button type="button" variant="outline" className="cursor-pointer" onClick={() => setIsAddOpen(false)}>Cancel</Button>
                <Button type="submit" className="bg-primary hover:bg-secondary cursor-pointer">Save Product</Button>
              </DialogFooter>
            </form>
          </DialogContent>
        </Dialog>
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
      <Card className="border-slate-200/80 dark:border-slate-800 bg-white dark:bg-slate-900 shadow-sm overflow-hidden">
        <CardContent className="p-0">
          <Table>
            <TableHeader className="bg-slate-50 dark:bg-slate-800/50 border-b border-slate-100 dark:border-slate-800">
              <TableRow className="hover:bg-transparent">
                <TableHead className="w-24 pl-6 font-bold uppercase tracking-wider text-xxs dark:text-slate-400">ID</TableHead>
                <TableHead className="font-bold uppercase tracking-wider text-xxs dark:text-slate-400">Product Name</TableHead>
                <TableHead className="font-bold uppercase tracking-wider text-xxs dark:text-slate-400">Price</TableHead>
                <TableHead className="font-bold uppercase tracking-wider text-xxs dark:text-slate-400">Stock Status</TableHead>
                <TableHead className="text-center pr-6 font-bold uppercase tracking-wider text-xxs dark:text-slate-400">Actions</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody className="text-sm font-medium text-slate-700 dark:text-slate-300 divide-y divide-slate-100/50 dark:divide-slate-800">
              {filteredProducts.length > 0 ? (
                filteredProducts.map((p) => (
                  <TableRow key={p.id} className="border-slate-100/50 dark:border-slate-800 hover:bg-slate-50/30 dark:hover:bg-slate-800/50 transition-colors">
                    <TableCell className="pl-6 font-mono text-xs font-bold text-slate-500 dark:text-slate-400">{p.id}</TableCell>
                    <TableCell className="font-semibold text-slate-900 dark:text-slate-100">{p.name}</TableCell>
                    <TableCell className="font-mono font-bold text-slate-900 dark:text-slate-100">{formatPrice(p.price)}</TableCell>
                    <TableCell>{getStockBadge(p.stock)}</TableCell>
                    <TableCell className="text-center pr-6">
                      <div className="flex items-center justify-center gap-1.5">
                        <Button 
                          variant="ghost" 
                          size="icon-sm" 
                          className="h-7 w-7 text-slate-500 dark:text-slate-400 hover:text-primary dark:hover:text-primary hover:bg-slate-100 dark:hover:bg-slate-800 rounded-md cursor-pointer"
                          title="Edit Product"
                          onClick={() => {
                            setSelectedEditProduct(p)
                            setEditProduct({ name: p.name, price: p.price })
                            setIsEditOpen(true)
                          }}
                        >
                          <Pencil className="w-4 h-4" />
                        </Button>
                        <Button 
                          variant="ghost" 
                          size="icon-sm" 
                          className="h-7 w-7 text-slate-500 dark:text-slate-400 hover:text-primary dark:hover:text-primary hover:bg-slate-100 dark:hover:bg-slate-800 rounded-md cursor-pointer"
                          title="Adjust Stock"
                          onClick={() => {
                            setSelectedProduct(p)
                            setIsAdjustOpen(true)
                          }}
                        >
                          <SlidersHorizontal className="w-4 h-4" />
                        </Button>
                        <Button 
                          variant="ghost" 
                          size="icon-sm" 
                          className="h-7 w-7 text-slate-500 dark:text-slate-400 hover:text-red-600 hover:bg-red-50 dark:hover:bg-red-950/40 rounded-md cursor-pointer"
                          title="Delete Product"
                          onClick={() => handleDelete(p.id)}
                        >
                          <Trash2 className="w-4 h-4" />
                        </Button>
                      </div>
                    </TableCell>
                  </TableRow>
                ))
              ) : (
                <TableRow>
                  <TableCell colSpan={5} className="h-32 text-center text-slate-400 dark:text-slate-500">
                    No products found matching your search.
                  </TableCell>
                </TableRow>
              )}
            </TableBody>
          </Table>

          {/* Pagination Footer */}
          <div className="border-t border-slate-100 dark:border-slate-800 px-6 py-4 flex items-center justify-between bg-slate-50/50 dark:bg-slate-800/30">
            <span className="text-xs font-semibold text-slate-500 dark:text-slate-400">
              Showing {filteredProducts.length} products
            </span>
            <div className="flex items-center gap-2">
              <Button variant="outline" size="sm" className="h-8 text-xs font-semibold border-slate-200 dark:border-slate-800 text-slate-500 dark:text-slate-400 hover:text-slate-900 dark:hover:text-slate-100 cursor-not-allowed" disabled>
                <ChevronLeft className="w-3.5 h-3.5 mr-1" />
                Previous
              </Button>
              <Button variant="outline" size="sm" className="h-8 text-xs font-semibold border-slate-200 dark:border-slate-800 text-slate-500 dark:text-slate-400 hover:text-slate-900 dark:hover:text-slate-100 cursor-not-allowed" disabled>
                Next
                <ChevronRight className="w-3.5 h-3.5 ml-1" />
              </Button>
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Adjust Stock Dialog */}
      <Dialog open={isAdjustOpen} onOpenChange={(open) => {
        setIsAdjustOpen(open)
        if(!open) setSelectedProduct(null)
      }}>
        <DialogContent className="max-w-sm">
          <form onSubmit={handleAdjustSubmit}>
            <DialogHeader>
              <DialogTitle className="text-xl font-extrabold text-slate-900 dark:text-slate-100">Adjust Stock Count</DialogTitle>
              <DialogDescription className="text-slate-500 dark:text-slate-400">
                Modify stock count for <span className="font-bold text-slate-900 dark:text-slate-100">{selectedProduct?.name}</span>.
              </DialogDescription>
            </DialogHeader>

            <div className="space-y-4 py-4">
              <div className="flex items-center justify-between text-slate-600 dark:text-slate-400 text-sm font-semibold">
                <span>Current Stock:</span>
                <span className="font-bold text-slate-900 dark:text-slate-100">{selectedProduct?.stock}</span>
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
              <Button type="button" variant="outline" className="cursor-pointer" onClick={() => setIsAdjustOpen(false)}>Cancel</Button>
              <Button type="submit" className="bg-primary hover:bg-secondary cursor-pointer">Apply Adjustment</Button>
            </DialogFooter>
          </form>
        </DialogContent>
      </Dialog>

      {/* Edit Product Dialog */}
      <Dialog open={isEditOpen} onOpenChange={(open) => {
        setIsEditOpen(open)
        if(!open) setSelectedEditProduct(null)
      }}>
        <DialogContent className="max-w-md">
          <form onSubmit={handleEditSubmit}>
            <DialogHeader>
              <DialogTitle className="text-xl font-extrabold text-slate-900 dark:text-slate-100">Edit Product Details</DialogTitle>
              <DialogDescription className="text-slate-500 dark:text-slate-400">
                Update basic details for product <span className="font-mono font-bold text-slate-900 dark:text-slate-100">{selectedEditProduct?.id}</span>.
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
              <Button type="button" variant="outline" className="cursor-pointer" onClick={() => setIsEditOpen(false)}>Cancel</Button>
              <Button type="submit" className="bg-primary hover:bg-secondary cursor-pointer">Save Changes</Button>
            </DialogFooter>
          </form>
        </DialogContent>
      </Dialog>
    </div>
  )
}

export default InventoryPage
