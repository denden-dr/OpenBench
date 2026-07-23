import type { Product } from '@/types/pos'
import { Card, CardContent } from '@/components/ui/card'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Trash2, SlidersHorizontal, Pencil, ChevronLeft, ChevronRight } from 'lucide-react'

interface InventoryTableProps {
  products: Product[]
  onEdit: (product: Product) => void
  onAdjust: (product: Product) => void
  onDelete: (id: string) => void
}

export function InventoryTable({ products, onEdit, onAdjust, onDelete }: InventoryTableProps) {
  const getStockBadge = (stock: number) => {
    if (stock === 0) {
      return <Badge className="bg-red-500/10 text-red-600 border-none font-semibold">Out of Stock</Badge>
    }
    if (stock < 5) {
      return <Badge className="bg-orange-500/10 text-orange-600 border-none font-semibold">Low Stock ({stock})</Badge>
    }
    return <Badge className="bg-green-500/10 text-green-600 border-none font-semibold">In Stock ({stock})</Badge>
  }

  const formatPrice = (price: number) => {
    return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(price)
  }

  return (
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
            {products.length > 0 ? (
              products.map((p) => (
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
                        onClick={() => onEdit(p)}
                      >
                        <Pencil className="w-4 h-4" />
                      </Button>
                      <Button 
                        variant="ghost" 
                        size="icon-sm" 
                        className="h-7 w-7 text-slate-500 dark:text-slate-400 hover:text-primary dark:hover:text-primary hover:bg-slate-100 dark:hover:bg-slate-800 rounded-md cursor-pointer"
                        title="Adjust Stock"
                        onClick={() => onAdjust(p)}
                      >
                        <SlidersHorizontal className="w-4 h-4" />
                      </Button>
                      <Button 
                        variant="ghost" 
                        size="icon-sm" 
                        className="h-7 w-7 text-slate-500 dark:text-slate-400 hover:text-red-600 hover:bg-red-50 dark:hover:bg-red-950/40 rounded-md cursor-pointer"
                        title="Delete Product"
                        onClick={() => onDelete(p.id)}
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
            Showing {products.length} products
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
  )
}
