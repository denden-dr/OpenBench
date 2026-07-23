import { useState, useEffect, useCallback } from 'react'
import type { Product, CartItem, POSTransaction, PaymentMethod, POSCheckoutRequest } from '@/types/pos'
import { inventoryService } from '@/services/inventoryService'
import { posService } from '@/services/posService'
import { Button } from '@/components/ui/button'
import { Tabs, TabsList, TabsTrigger, TabsContent } from '@/components/ui/tabs'
import { ProductGrid } from '@/features/pos/components/ProductGrid'
import { CartPanel } from '@/features/pos/components/CartPanel'
import { TransactionsTable } from '@/features/pos/components/TransactionsTable'
import { TransactionDetailsModal } from '@/features/pos/components/TransactionDetailsModal'

function POSPage() {
  const [products, setProducts] = useState<Product[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [cart, setCart] = useState<CartItem[]>([])
  const [searchQuery, setSearchQuery] = useState('')
  const [paymentMethod, setPaymentMethod] = useState<PaymentMethod>('CASH')
  const [transactions, setTransactions] = useState<POSTransaction[]>([])
  const [loadingTx, setLoadingTx] = useState(false)
  
  // Dialog Details
  const [selectedTx, setSelectedTx] = useState<POSTransaction | null>(null)
  const [isDetailOpen, setIsDetailOpen] = useState(false)

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

  const fetchTransactions = useCallback(async () => {
    try {
      setLoadingTx(true)
      const result = await posService.getTransactions({ limit: 50 })
      setTransactions(result.data)
    } catch (err: any) {
      console.error('Failed to load transactions', err)
    } finally {
      setLoadingTx(false)
    }
  }, [])

  useEffect(() => {
    fetchProducts()
    fetchTransactions()
  }, [fetchProducts, fetchTransactions])

  const handleAddToCart = (product: Product) => {
    const existing = cart.find(item => item.product.id === product.id)
    if (existing) {
      if (existing.quantity >= product.stock) {
        alert('Cannot add more items than available in stock.')
        return
      }
      setCart(cart.map(item => 
        item.product.id === product.id 
          ? { ...item, quantity: item.quantity + 1 }
          : item
      ))
    } else {
      setCart([...cart, { product, quantity: 1 }])
    }
  }

  const handleUpdateQuantity = (productId: string, delta: number) => {
    const existing = cart.find(item => item.product.id === productId)
    if (!existing) return

    const newQty = existing.quantity + delta
    if (newQty <= 0) {
      handleRemoveItem(productId)
    } else {
      if (delta > 0 && newQty > existing.product.stock) {
        alert('Cannot add more items than available in stock.')
        return
      }
      setCart(cart.map(item => 
        item.product.id === productId 
          ? { ...item, quantity: newQty }
          : item
      ))
    }
  }

  const handleRemoveItem = (productId: string) => {
    setCart(cart.filter(item => item.product.id !== productId))
  }

  const handleCheckout = async (e: React.FormEvent) => {
    e.preventDefault()
    if (cart.length === 0) return

    const checkoutRequest: POSCheckoutRequest = {
      payment_method: paymentMethod,
      items: cart.map(item => ({
        product_id: item.product.id,
        quantity: item.quantity,
      })),
    }

    try {
      const tx = await posService.checkout(checkoutRequest)
      setTransactions([tx, ...transactions])
      await fetchProducts()
      setCart([])
      alert('Checkout completed successfully!')
    } catch (err: any) {
      alert(err?.response?.data?.detail || 'Checkout failed')
    }
  }

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
      <div>
        <h1 className="text-3xl font-extrabold text-slate-900 dark:text-slate-100 tracking-tight">
          Point of Sale (POS)
        </h1>
        <p className="text-slate-500 dark:text-slate-400 text-sm">Sell accessories, process checkout transactions, and review cashier history.</p>
      </div>

      <Tabs defaultValue="cashier" className="w-full">
        <TabsList className="bg-slate-100 dark:bg-slate-800/60 border border-slate-200/50 dark:border-slate-800 p-1 mb-6">
          <TabsTrigger value="cashier" className="data-[state=active]:bg-white dark:data-[state=active]:bg-slate-900 dark:data-[state=active]:text-slate-100 data-[state=active]:shadow-sm font-semibold text-xs px-6">Cashier Screen</TabsTrigger>
          <TabsTrigger value="history" className="data-[state=active]:bg-white dark:data-[state=active]:bg-slate-900 dark:data-[state=active]:text-slate-100 data-[state=active]:shadow-sm font-semibold text-xs px-6">Transaction History</TabsTrigger>
        </TabsList>

        {/* CASHIER CONTENT */}
        <TabsContent value="cashier" className="grid grid-cols-1 lg:grid-cols-12 gap-8 outline-none">
          {/* LEFT: Products Grid (8 columns) */}
          <div className="lg:col-span-8">
            <ProductGrid 
              products={products}
              searchQuery={searchQuery}
              setSearchQuery={setSearchQuery}
              onAddToCart={handleAddToCart}
            />
          </div>

          {/* RIGHT: Cart Drawer (4 columns) */}
          <div className="lg:col-span-4">
            <CartPanel 
              cart={cart}
              onUpdateQuantity={handleUpdateQuantity}
              onRemoveItem={handleRemoveItem}
              paymentMethod={paymentMethod}
              setPaymentMethod={setPaymentMethod}
              onCheckout={handleCheckout}
            />
          </div>
        </TabsContent>

        {/* TRANSACTIONS HISTORY CONTENT */}
        <TabsContent value="history" className="outline-none">
          <TransactionsTable 
            transactions={transactions}
            loading={loadingTx}
            onViewDetails={(tx) => {
              setSelectedTx(tx)
              setIsDetailOpen(true)
            }}
          />
        </TabsContent>
      </Tabs>

      <TransactionDetailsModal 
        transaction={selectedTx}
        isOpen={isDetailOpen}
        onClose={() => {
          setIsDetailOpen(false)
          setSelectedTx(null)
        }}
      />
    </div>
  )
}

export default POSPage
