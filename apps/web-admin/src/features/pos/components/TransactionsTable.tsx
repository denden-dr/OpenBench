import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Card, CardContent } from '@/components/ui/card'
import { ChevronLeft, ChevronRight } from 'lucide-react'
import type { POSTransaction } from '@/types/pos'

interface TransactionsTableProps {
  transactions: POSTransaction[]
  loading: boolean
  onViewDetails: (tx: POSTransaction) => void
}

export function TransactionsTable({
  transactions,
  loading,
  onViewDetails,
}: TransactionsTableProps) {
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
    <div>
      {loading && (
        <div className="flex justify-center py-8">
          <p className="text-slate-500 text-sm">Loading transactions...</p>
        </div>
      )}
      <Card className="border-slate-200/80 dark:border-slate-800 bg-white dark:bg-slate-900 shadow-sm overflow-hidden">
        <CardContent className="p-0">
          <Table>
            <TableHeader className="bg-slate-50 dark:bg-slate-800/50 border-b border-slate-100 dark:border-slate-800">
              <TableRow className="hover:bg-transparent">
                <TableHead className="w-48 pl-6 font-bold uppercase tracking-wider text-xxs dark:text-slate-400">Transaction ID</TableHead>
                <TableHead className="font-bold uppercase tracking-wider text-xxs dark:text-slate-400">Date & Time</TableHead>
                <TableHead className="font-bold uppercase tracking-wider text-xxs dark:text-slate-400">Payment Method</TableHead>
                <TableHead className="font-bold uppercase tracking-wider text-xxs dark:text-slate-400">Total Amount</TableHead>
                <TableHead className="text-center pr-6 font-bold uppercase tracking-wider text-xxs dark:text-slate-400">Actions</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody className="text-sm font-medium text-slate-700 dark:text-slate-300 divide-y divide-slate-100/50 dark:divide-slate-800">
              {transactions.length > 0 ? (
                transactions.map((tx) => (
                  <TableRow key={tx.id} className="border-slate-100/50 dark:border-slate-800 hover:bg-slate-50/30 dark:hover:bg-slate-800/50 transition-colors">
                    <TableCell className="pl-6 font-mono text-xs font-bold text-slate-600 dark:text-slate-400">{tx.id}</TableCell>
                    <TableCell className="text-slate-500 dark:text-slate-400 font-semibold">{formatDate(tx.created_at)}</TableCell>
                    <TableCell>
                      <Badge variant="outline" className={`font-semibold ${tx.payment_method === 'QRIS' ? 'bg-purple-50 dark:bg-purple-950/40 text-purple-600 dark:text-purple-400 border-purple-200 dark:border-purple-800' : 'bg-green-50 dark:bg-green-950/40 text-green-600 dark:text-green-400 border-green-200 dark:border-green-800'}`}>
                        {tx.payment_method}
                      </Badge>
                    </TableCell>
                    <TableCell className="font-mono font-bold text-slate-900 dark:text-slate-100">{formatPrice(tx.total_amount)}</TableCell>
                    <TableCell className="text-center pr-6">
                      <Button 
                        variant="ghost" 
                        size="sm" 
                        className="text-xs font-semibold text-slate-500 dark:text-slate-400 hover:text-primary dark:hover:text-primary hover:bg-slate-100 dark:hover:bg-slate-800 rounded-md cursor-pointer px-3 h-7"
                        onClick={() => onViewDetails(tx)}
                      >
                        Details
                      </Button>
                    </TableCell>
                  </TableRow>
                ))
              ) : (
                <TableRow>
                  <TableCell colSpan={5} className="h-32 text-center text-slate-400 dark:text-slate-500">
                    No transactions registered yet.
                  </TableCell>
                </TableRow>
              )}
            </TableBody>
          </Table>

          {/* Pagination Footer */}
          <div className="border-t border-slate-100 dark:border-slate-800 px-6 py-4 flex items-center justify-between bg-slate-50/50 dark:bg-slate-800/30">
            <span className="text-xs font-semibold text-slate-500 dark:text-slate-400">
              Showing {transactions.length} transactions
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
    </div>
  )
}
