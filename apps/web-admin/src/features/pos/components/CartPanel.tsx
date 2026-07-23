import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { ShoppingCart, Trash2, Plus, Minus, CreditCard } from 'lucide-react'
import type { CartItem, PaymentMethod } from '@/types/pos'

interface CartPanelProps {
  cart: CartItem[]
  onUpdateQuantity: (productId: string, delta: number) => void
  onRemoveItem: (productId: string) => void
  paymentMethod: PaymentMethod
  setPaymentMethod: (method: PaymentMethod) => void
  onCheckout: (e: React.FormEvent) => void
}

export function CartPanel({
  cart,
  onUpdateQuantity,
  onRemoveItem,
  paymentMethod,
  setPaymentMethod,
  onCheckout,
}: CartPanelProps) {
  const getCartTotal = () => {
    return cart.reduce((sum, item) => sum + (item.product.price * item.quantity), 0)
  }

  const formatPrice = (price: number) => {
    return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(price)
  }

  return (
    <Card className="border-slate-200/80 dark:border-slate-800 bg-white dark:bg-slate-900 shadow-sm flex flex-col h-[600px]">
      <CardHeader className="border-b border-slate-100 dark:border-slate-800 flex flex-row items-center justify-between pb-3">
        <div>
          <CardTitle className="font-extrabold text-slate-800 dark:text-slate-100 text-base flex items-center gap-1.5">
            <ShoppingCart className="w-4 h-4 text-primary" />
            Shopping Cart
          </CardTitle>
        </div>
        <Badge variant="outline" className="text-slate-600 dark:text-slate-400 border-slate-200 dark:border-slate-800 font-semibold">{cart.length} items</Badge>
      </CardHeader>

      <CardContent className="flex-grow overflow-y-auto p-4 space-y-4">
        {cart.length > 0 ? (
          cart.map(item => (
            <div key={item.product.id} className="flex justify-between items-start gap-3 pb-3 border-b border-slate-100/60 dark:border-slate-800 last:border-0 last:pb-0">
              <div className="flex-grow space-y-1">
                <p className="text-xs font-bold text-slate-900 dark:text-slate-100 line-clamp-1">{item.product.name}</p>
                <p className="text-xxs font-mono text-slate-500 dark:text-slate-400 font-semibold">{formatPrice(item.product.price)} x {item.quantity}</p>
              </div>
              
              <div className="flex items-center gap-2">
                <div className="flex items-center border border-slate-200 dark:border-slate-800 rounded-md">
                  <Button 
                    variant="ghost" 
                    size="icon-xs" 
                    className="h-6 w-6 rounded-r-none cursor-pointer"
                    onClick={() => onUpdateQuantity(item.product.id, -1)}
                  >
                    <Minus className="w-3 h-3" />
                  </Button>
                  <span className="w-6 text-center font-mono text-xs font-bold text-slate-700 dark:text-slate-300">{item.quantity}</span>
                  <Button 
                    variant="ghost" 
                    size="icon-xs" 
                    className="h-6 w-6 rounded-l-none cursor-pointer"
                    onClick={() => onUpdateQuantity(item.product.id, 1)}
                  >
                    <Plus className="w-3 h-3" />
                  </Button>
                </div>

                <Button 
                  variant="ghost" 
                  size="icon-xs" 
                  className="h-6 w-6 text-slate-400 hover:text-red-500 cursor-pointer"
                  onClick={() => onRemoveItem(item.product.id)}
                >
                  <Trash2 className="w-3.5 h-3.5" />
                </Button>
              </div>
            </div>
          ))
        ) : (
          <div className="h-full flex flex-col items-center justify-center text-center py-20 text-slate-400 dark:text-slate-500 gap-1.5">
            <ShoppingCart className="w-8 h-8 text-slate-300 dark:text-slate-600" />
            <span className="text-xs">Your shopping cart is empty.<br />Click products in catalog to add.</span>
          </div>
        )}
      </CardContent>

      {/* Checkout Form Footer */}
      <div className="border-t border-slate-100 dark:border-slate-800 p-4 bg-slate-50/50 dark:bg-slate-800/30 space-y-4 rounded-b-xl">
        <div className="space-y-1.5">
          <span className="text-xxs font-bold text-slate-400 dark:text-slate-500 uppercase tracking-wider">Payment Method</span>
          <div className="grid grid-cols-2 gap-2">
            <Button 
              type="button" 
              variant={paymentMethod === 'CASH' ? 'default' : 'outline'}
              className={`font-bold text-xs h-9 cursor-pointer ${paymentMethod === 'CASH' ? 'bg-primary hover:bg-secondary' : 'border-slate-200 dark:border-slate-800'}`}
              onClick={() => setPaymentMethod('CASH')}
            >
              Cash
            </Button>
            <Button 
              type="button" 
              variant={paymentMethod === 'QRIS' ? 'default' : 'outline'}
              className={`font-bold text-xs h-9 cursor-pointer ${paymentMethod === 'QRIS' ? 'bg-primary hover:bg-secondary' : 'border-slate-200 dark:border-slate-800'}`}
              onClick={() => setPaymentMethod('QRIS')}
            >
              QRIS Code
            </Button>
          </div>
        </div>

        <div className="flex justify-between items-center pt-2">
          <span className="text-xs font-bold text-slate-500 dark:text-slate-400">Grand Total:</span>
          <span className="font-mono text-lg font-extrabold text-slate-900 dark:text-slate-100">{formatPrice(getCartTotal())}</span>
        </div>

        <Button 
          onClick={onCheckout} 
          disabled={cart.length === 0}
          className="w-full font-bold bg-green-600 hover:bg-green-700 cursor-pointer disabled:bg-slate-200 dark:disabled:bg-slate-800 disabled:text-slate-400 dark:disabled:text-slate-600 disabled:cursor-not-allowed text-xs h-10"
        >
          <CreditCard className="w-4 h-4 mr-1.5" />
          Process Checkout
        </Button>
      </div>
    </Card>
  )
}
