import { useState } from 'react'
import type { Product, CartItem, POSTransaction, PaymentMethod, POSCheckoutItem, POSCheckoutRequest } from '@/types/pos'
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Tabs, TabsList, TabsTrigger, TabsContent } from '@/components/ui/tabs'
import { 
  Dialog, 
  DialogContent, 
  DialogDescription, 
  DialogHeader, 
  DialogTitle, 
  DialogFooter
} from '@/components/ui/dialog'
import { Search, ShoppingCart, Trash2, Plus, Minus, CreditCard, ShieldCheck, ChevronLeft, ChevronRight } from 'lucide-react'

// Initial products catalogue for cashier
const initialProducts: Product[] = [
  { id: 'p-1', name: 'Tempered Glass iPhone 15 Pro Max', price: 75000, stock: 18, created_at: '2026-07-15T09:00:00Z' },
  { id: 'p-2', name: 'Silicon Case iPhone 15', price: 120000, stock: 4, created_at: '2026-07-15T09:15:00Z' },
  { id: 'p-3', name: 'USB-C Charger Adapter 20W', price: 299000, stock: 2, created_at: '2026-07-16T11:00:00Z' },
  { id: 'p-4', name: 'MicroUSB Cable 1m', price: 35000, stock: 25, created_at: '2026-07-16T11:30:00Z' },
  { id: 'p-5', name: 'Lightning to USB-C Cable 2m', price: 150000, stock: 0, created_at: '2026-07-17T14:00:00Z' }
]

// Mock transaction history
const initialTransactions: POSTransaction[] = [
  {
    id: 'tx-1',
    payment_method: 'QRIS',
    total_amount: 150000,
    created_at: '2026-07-18T10:30:00Z',
    items: [
      { id: 'item-1', product_id: 'p-1', product_name: 'Tempered Glass iPhone 15 Pro Max', quantity: 2, price: 75000 }
    ]
  },
  {
    id: 'tx-2',
    payment_method: 'CASH',
    total_amount: 320000,
    created_at: '2026-07-18T14:20:00Z',
    items: [
      { id: 'item-2', product_id: 'p-2', product_name: 'Silicon Case iPhone 15', quantity: 1, price: 120000 },
      { id: 'item-3', product_id: 'p-3', product_name: 'USB-C Charger Adapter 20W', quantity: 1, price: 299000 }
    ]
  }
]

function POSPage() {
  const [products, setProducts] = useState<Product[]>(initialProducts)
  const [cart, setCart] = useState<CartItem[]>([])
  const [searchQuery, setSearchQuery] = useState('')
  const [paymentMethod, setPaymentMethod] = useState<PaymentMethod>('CASH')
  const [transactions, setTransactions] = useState<POSTransaction[]>(initialTransactions)
  
  // Dialog Details
  const [selectedTx, setSelectedTx] = useState<POSTransaction | null>(null)
  const [isDetailOpen, setIsDetailOpen] = useState(false)

  // Filter Catalog
  const filteredProducts = products.filter(p => 
    p.name.toLowerCase().includes(searchQuery.toLowerCase()) && p.stock > 0
  )

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

  const getCartTotal = () => {
    return cart.reduce((sum, item) => sum + (item.product.price * item.quantity), 0)
  }

  const handleCheckout = (e: React.FormEvent) => {
    e.preventDefault()
    if (cart.length === 0) return

    // Construct the payload matching POSCheckoutRequest for API contract
    const checkoutRequest: POSCheckoutRequest = {
      payment_method: paymentMethod,
      items: cart.map(item => ({
        product_id: item.product.id,
        quantity: item.quantity
      }))
    }
    
    // Log request payload for verification (simulating API POST payload)
    console.log('Sending POS Checkout Request Payload:', checkoutRequest)

    const cartTotal = getCartTotal()
    
    // Create new transaction matching v1 specs for local state display
    const newTxItems: POSCheckoutItem[] = cart.map((item) => ({
      id: `item-${Math.random().toString(36).substr(2, 9)}`,
      product_id: item.product.id,
      product_name: item.product.name,
      quantity: item.quantity,
      price: item.product.price
    }))

    const newTx: POSTransaction = {
      id: `tx-${Math.random().toString(36).substr(2, 9)}`,
      payment_method: paymentMethod,
      total_amount: cartTotal,
      created_at: new Date().toISOString(),
      items: newTxItems
    }

    // Deduct stock
    setProducts(products.map(p => {
      const cartItem = cart.find(item => item.product.id === p.id)
      if (cartItem) {
        return {
          ...p,
          stock: Math.max(0, p.stock - cartItem.quantity)
        }
      }
      return p
    }))

    // Save transaction
    setTransactions([newTx, ...transactions])
    setCart([])
    alert('Checkout completed successfully!')
  }

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
    <div className="space-y-8">
      {/* Title */}
      <div>
        <h1 className="text-3xl font-extrabold text-slate-900 tracking-tight">
          Point of Sale (POS)
        </h1>
        <p className="text-slate-500 text-sm">Sell accessories, process checkout transactions, and review cashier history.</p>
      </div>

      <Tabs defaultValue="cashier" className="w-full">
        <TabsList className="bg-slate-100 border border-slate-200/50 p-1 mb-6">
          <TabsTrigger value="cashier" className="data-[state=active]:bg-white data-[state=active]:shadow-sm font-semibold text-xs px-6">Cashier Screen</TabsTrigger>
          <TabsTrigger value="history" className="data-[state=active]:bg-white data-[state=active]:shadow-sm font-semibold text-xs px-6">Transaction History</TabsTrigger>
        </TabsList>

        {/* CASHIER CONTENT */}
        <TabsContent value="cashier" className="grid grid-cols-1 lg:grid-cols-12 gap-8 outline-none">
          {/* LEFT: Products Grid (8 columns) */}
          <div className="lg:col-span-8 space-y-6">
            <div className="relative">
              <Search className="absolute left-3 top-2.5 h-4 w-4 text-slate-400" />
              <Input
                placeholder="Search accessories in stock..."
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                className="pl-9 bg-white border-slate-200 focus-visible:ring-primary/20"
              />
            </div>

            <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-4">
              {filteredProducts.length > 0 ? (
                filteredProducts.map(p => (
                  <Card 
                    key={p.id} 
                    className="hover:shadow-md hover:border-primary/40 transition-all border-slate-200 bg-white cursor-pointer select-none"
                    onClick={() => handleAddToCart(p)}
                  >
                    <CardHeader className="p-4 pb-2">
                      <div className="flex justify-between items-start gap-2">
                        <CardTitle className="text-sm font-bold text-slate-800 line-clamp-2">{p.name}</CardTitle>
                      </div>
                    </CardHeader>
                    <CardContent className="p-4 pt-0 flex justify-between items-center mt-2">
                      <span className="font-mono text-sm font-bold text-slate-900">{formatPrice(p.price)}</span>
                      <span className="text-xxs font-bold px-2 py-0.5 rounded-full bg-slate-100 text-slate-600">Stock: {p.stock}</span>
                    </CardContent>
                  </Card>
                ))
              ) : (
                <div className="col-span-full py-16 text-center text-slate-400">
                  No active accessories found matching your query.
                </div>
              )}
            </div>
          </div>

          {/* RIGHT: Cart Drawer (4 columns) */}
          <div className="lg:col-span-4">
            <Card className="border-slate-200/80 bg-white shadow-sm flex flex-col h-[600px]">
              <CardHeader className="border-b border-slate-100 flex flex-row items-center justify-between pb-3">
                <div>
                  <CardTitle className="font-extrabold text-slate-800 text-base flex items-center gap-1.5">
                    <ShoppingCart className="w-4 h-4 text-primary" />
                    Shopping Cart
                  </CardTitle>
                </div>
                <Badge variant="outline" className="text-slate-600 font-semibold">{cart.length} items</Badge>
              </CardHeader>

              <CardContent className="flex-grow overflow-y-auto p-4 space-y-4">
                {cart.length > 0 ? (
                  cart.map(item => (
                    <div key={item.product.id} className="flex justify-between items-start gap-3 pb-3 border-b border-slate-100/60 last:border-0 last:pb-0">
                      <div className="flex-grow space-y-1">
                        <p className="text-xs font-bold text-slate-900 line-clamp-1">{item.product.name}</p>
                        <p className="text-xxs font-mono text-slate-500 font-semibold">{formatPrice(item.product.price)} x {item.quantity}</p>
                      </div>
                      
                      <div className="flex items-center gap-2">
                        <div className="flex items-center border border-slate-200 rounded-md">
                          <Button 
                            variant="ghost" 
                            size="icon-xs" 
                            className="h-6 w-6 rounded-r-none cursor-pointer"
                            onClick={() => handleUpdateQuantity(item.product.id, -1)}
                          >
                            <Minus className="w-3 h-3" />
                          </Button>
                          <span className="w-6 text-center font-mono text-xs font-bold text-slate-700">{item.quantity}</span>
                          <Button 
                            variant="ghost" 
                            size="icon-xs" 
                            className="h-6 w-6 rounded-l-none cursor-pointer"
                            onClick={() => handleUpdateQuantity(item.product.id, 1)}
                          >
                            <Plus className="w-3 h-3" />
                          </Button>
                        </div>

                        <Button 
                          variant="ghost" 
                          size="icon-xs" 
                          className="h-6 w-6 text-slate-400 hover:text-red-500 cursor-pointer"
                          onClick={() => handleRemoveItem(item.product.id)}
                        >
                          <Trash2 className="w-3.5 h-3.5" />
                        </Button>
                      </div>
                    </div>
                  ))
                ) : (
                  <div className="h-full flex flex-col items-center justify-center text-center py-20 text-slate-400 gap-1.5">
                    <ShoppingCart className="w-8 h-8 text-slate-300" />
                    <span className="text-xs">Your shopping cart is empty.<br />Click products in catalog to add.</span>
                  </div>
                )}
              </CardContent>

              {/* Checkout Form Footer */}
              <div className="border-t border-slate-100 p-4 bg-slate-50/50 space-y-4 rounded-b-xl">
                <div className="space-y-1.5">
                  <span className="text-xxs font-bold text-slate-400 uppercase tracking-wider">Payment Method</span>
                  <div className="grid grid-cols-2 gap-2">
                    <Button 
                      type="button" 
                      variant={paymentMethod === 'CASH' ? 'default' : 'outline'}
                      className={`font-bold text-xs h-9 cursor-pointer ${paymentMethod === 'CASH' ? 'bg-primary hover:bg-secondary' : 'border-slate-200'}`}
                      onClick={() => setPaymentMethod('CASH')}
                    >
                      Cash
                    </Button>
                    <Button 
                      type="button" 
                      variant={paymentMethod === 'QRIS' ? 'default' : 'outline'}
                      className={`font-bold text-xs h-9 cursor-pointer ${paymentMethod === 'QRIS' ? 'bg-primary hover:bg-secondary' : 'border-slate-200'}`}
                      onClick={() => setPaymentMethod('QRIS')}
                    >
                      QRIS Code
                    </Button>
                  </div>
                </div>

                <div className="flex justify-between items-center pt-2">
                  <span className="text-xs font-bold text-slate-500">Grand Total:</span>
                  <span className="font-mono text-lg font-extrabold text-slate-900">{formatPrice(getCartTotal())}</span>
                </div>

                <Button 
                  onClick={handleCheckout} 
                  disabled={cart.length === 0}
                  className="w-full font-bold bg-green-600 hover:bg-green-700 cursor-pointer disabled:bg-slate-200 disabled:text-slate-400 disabled:cursor-not-allowed text-xs h-10"
                >
                  <CreditCard className="w-4 h-4 mr-1.5" />
                  Process Checkout
                </Button>
              </div>
            </Card>
          </div>
        </TabsContent>

        {/* TRANSACTIONS HISTORY CONTENT */}
        <TabsContent value="history" className="outline-none">
          <Card className="border-slate-200/80 bg-white shadow-sm overflow-hidden">
            <CardContent className="p-0">
              <Table>
                <TableHeader className="bg-slate-50 border-b border-slate-100">
                  <TableRow className="hover:bg-transparent">
                    <TableHead className="w-48 pl-6 font-bold uppercase tracking-wider text-xxs">Transaction ID</TableHead>
                    <TableHead className="font-bold uppercase tracking-wider text-xxs">Date & Time</TableHead>
                    <TableHead className="font-bold uppercase tracking-wider text-xxs">Payment Method</TableHead>
                    <TableHead className="font-bold uppercase tracking-wider text-xxs">Total Amount</TableHead>
                    <TableHead className="text-center pr-6 font-bold uppercase tracking-wider text-xxs">Actions</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody className="text-sm font-medium text-slate-700 divide-y divide-slate-100/50">
                  {transactions.length > 0 ? (
                    transactions.map((tx) => (
                      <TableRow key={tx.id} className="border-slate-100/50 hover:bg-slate-50/30 transition-colors">
                        <TableCell className="pl-6 font-mono text-xs font-bold text-slate-600">{tx.id}</TableCell>
                        <TableCell className="text-slate-500 font-semibold">{formatDate(tx.created_at)}</TableCell>
                        <TableCell>
                          <Badge variant="outline" className={`font-semibold ${tx.payment_method === 'QRIS' ? 'bg-purple-50 text-purple-600 border-purple-200' : 'bg-green-50 text-green-600 border-green-200'}`}>
                            {tx.payment_method}
                          </Badge>
                        </TableCell>
                        <TableCell className="font-mono font-bold text-slate-900">{formatPrice(tx.total_amount)}</TableCell>
                        <TableCell className="text-center pr-6">
                          <Button 
                            variant="ghost" 
                            size="sm" 
                            className="text-xs font-semibold text-slate-500 hover:text-primary hover:bg-slate-100 rounded-md cursor-pointer px-3 h-7"
                            onClick={() => {
                              setSelectedTx(tx)
                              setIsDetailOpen(true)
                            }}
                          >
                            Details
                          </Button>
                        </TableCell>
                      </TableRow>
                    ))
                  ) : (
                    <TableRow>
                      <TableCell colSpan={5} className="h-32 text-center text-slate-400">
                        No transactions registered yet.
                      </TableCell>
                    </TableRow>
                  )}
                </TableBody>
              </Table>

              {/* Pagination Footer */}
              <div className="border-t border-slate-100 px-6 py-4 flex items-center justify-between bg-slate-50/50">
                <span className="text-xs font-semibold text-slate-500">
                  Showing {transactions.length} transactions
                </span>
                <div className="flex items-center gap-2">
                  <Button variant="outline" size="sm" className="h-8 text-xs font-semibold border-slate-200 text-slate-500 hover:text-slate-900 cursor-not-allowed" disabled>
                    <ChevronLeft className="w-3.5 h-3.5 mr-1" />
                    Previous
                  </Button>
                  <Button variant="outline" size="sm" className="h-8 text-xs font-semibold border-slate-200 text-slate-500 hover:text-slate-900 cursor-not-allowed" disabled>
                    Next
                    <ChevronRight className="w-3.5 h-3.5 ml-1" />
                  </Button>
                </div>
              </div>
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>

      {/* Transaction Details Dialog */}
      <Dialog open={isDetailOpen} onOpenChange={(open) => {
        setIsDetailOpen(open)
        if(!open) setSelectedTx(null)
      }}>
        <DialogContent className="max-w-md">
          <DialogHeader>
            <DialogTitle className="text-xl font-extrabold text-slate-900">Transaction Details</DialogTitle>
            <DialogDescription>
              Purchased items and billing summary for transaction <span className="font-mono font-bold text-slate-900">{selectedTx?.id}</span>.
            </DialogDescription>
          </DialogHeader>

          <div className="space-y-4 py-4">
            <div className="flex justify-between items-center text-xs font-semibold text-slate-500 border-b border-slate-100 pb-2">
              <span>Date: {selectedTx && formatDate(selectedTx.created_at)}</span>
              <span>Method: <span className="font-bold text-slate-950">{selectedTx?.payment_method}</span></span>
            </div>

            <div className="space-y-3">
              <span className="text-xxs font-bold text-slate-400 uppercase tracking-wider">Items Purchased</span>
              <div className="space-y-2.5 max-h-48 overflow-y-auto">
                {selectedTx?.items.map((item) => (
                  <div key={item.id} className="flex justify-between text-slate-700 text-xs">
                    <span>{item.product_name} <span className="text-slate-400 font-medium">x {item.quantity}</span></span>
                    <span className="font-mono font-bold">{formatPrice(item.price * item.quantity)}</span>
                  </div>
                ))}
              </div>
            </div>

            <div className="border-t border-slate-100 pt-3 flex justify-between items-center">
              <span className="text-xs font-bold text-slate-500">Total Paid:</span>
              <span className="font-mono text-base font-extrabold text-slate-900">{selectedTx && formatPrice(selectedTx.total_amount)}</span>
            </div>
          </div>

          <DialogFooter className="pt-2">
            <Button type="button" className="w-full bg-primary hover:bg-secondary cursor-pointer text-xs" onClick={() => setIsDetailOpen(false)}>
              <ShieldCheck className="w-4 h-4 mr-1.5" />
              Close Invoice
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  )
}

export default POSPage
