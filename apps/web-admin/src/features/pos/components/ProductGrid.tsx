import { Search } from 'lucide-react'
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import type { Product } from '@/types/pos'

interface ProductGridProps {
  products: Product[]
  searchQuery: string
  setSearchQuery: (query: string) => void
  onAddToCart: (product: Product) => void
}

export function ProductGrid({
  products,
  searchQuery,
  setSearchQuery,
  onAddToCart,
}: ProductGridProps) {
  const filteredProducts = products.filter(p => 
    p.name.toLowerCase().includes(searchQuery.toLowerCase()) && p.stock > 0
  )

  const formatPrice = (price: number) => {
    return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(price)
  }

  return (
    <div className="space-y-6">
      <div className="relative">
        <Search className="absolute left-3 top-2.5 h-4 w-4 text-slate-400 dark:text-slate-500" />
        <Input
          placeholder="Search accessories in stock..."
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
          className="pl-9 bg-white dark:bg-slate-900 border-slate-200 dark:border-slate-800 text-slate-900 dark:text-slate-100 focus-visible:ring-primary/20"
        />
      </div>

      <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-4">
        {filteredProducts.length > 0 ? (
          filteredProducts.map(p => (
            <Card 
              key={p.id} 
              className="hover:shadow-md hover:border-primary/40 transition-all border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 cursor-pointer select-none"
              onClick={() => onAddToCart(p)}
            >
              <CardHeader className="p-4 pb-2">
                <div className="flex justify-between items-start gap-2">
                  <CardTitle className="text-sm font-bold text-slate-800 dark:text-slate-100 line-clamp-2">{p.name}</CardTitle>
                </div>
              </CardHeader>
              <CardContent className="p-4 pt-0 flex justify-between items-center mt-2">
                <span className="font-mono text-sm font-bold text-slate-900 dark:text-slate-100">{formatPrice(p.price)}</span>
                <span className="text-xxs font-bold px-2 py-0.5 rounded-full bg-slate-100 dark:bg-slate-800 text-slate-600 dark:text-slate-400">Stock: {p.stock}</span>
              </CardContent>
            </Card>
          ))
        ) : (
          <div className="col-span-full py-16 text-center text-slate-400 dark:text-slate-500">
            No active accessories found matching your query.
          </div>
        )}
      </div>
    </div>
  )
}
