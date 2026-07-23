import { Card, CardHeader, CardTitle, CardContent, CardDescription } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Search, ShieldCheck, ShieldAlert, ClipboardList, Ban, AlertTriangle } from 'lucide-react'
import type { Warranty } from '@/types/warranty'
import { WarrantyStatusBadge } from './WarrantyStatusBadge'

interface WarrantyCheckerPanelProps {
  searchTicket: string
  onSearchChange: (value: string) => void
  onSearch: () => void
  searching: boolean
  searchError: string | null
  hasSearched: boolean
  foundWarranty: Warranty | null
  searchedTicketNum: string
  onOpenClaimDialog: () => void
  onOpenVoidDialog: (warranty: Warranty) => void
}

function formatDate(dateStr: string) {
  return new Date(dateStr).toLocaleString('id-ID', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  })
}

export function WarrantyCheckerPanel({
  searchTicket,
  onSearchChange,
  onSearch,
  searching,
  searchError,
  hasSearched,
  foundWarranty,
  searchedTicketNum,
  onOpenClaimDialog,
  onOpenVoidDialog,
}: WarrantyCheckerPanelProps) {
  return (
    <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
      {/* Checker Form */}
      <Card className="lg:col-span-1 border-slate-200/80 dark:border-slate-800 bg-white dark:bg-slate-900 shadow-sm h-fit">
        <CardHeader>
          <CardTitle className="text-lg font-bold text-slate-900 dark:text-slate-100">Warranty Checker</CardTitle>
          <CardDescription className="text-slate-500 dark:text-slate-400">Enter the ticket number of the previous repair to verify warranty status.</CardDescription>
        </CardHeader>
        <CardContent>
          <form onSubmit={(e) => { e.preventDefault(); onSearch() }} className="space-y-4">
            <div className="space-y-1">
              <label className="text-xs font-bold text-slate-500 dark:text-slate-400 uppercase">Original Ticket ID</label>
              <div className="relative">
                <Search className="absolute left-3 top-2.5 h-4 w-4 text-slate-400 dark:text-slate-500" />
                <Input
                  required
                  placeholder="e.g. TKT-20260707-1234"
                  value={searchTicket}
                  onChange={(e) => onSearchChange(e.target.value)}
                  className="pl-9 bg-white dark:bg-slate-900 border-slate-200 dark:border-slate-800 text-slate-900 dark:text-slate-100 focus-visible:ring-primary/20"
                />
              </div>
              <span className="text-[10px] text-slate-400 dark:text-slate-500">
                Try:{' '}
                <span className="font-mono text-slate-600 dark:text-slate-300 font-semibold cursor-pointer underline" onClick={() => onSearchChange('TKT-20260707-1234')}>TKT-20260707-1234</span>{' '}
                (Active), or{' '}
                <span className="font-mono text-slate-600 dark:text-slate-300 font-semibold cursor-pointer underline" onClick={() => onSearchChange('TKT-20260601-5678')}>TKT-20260601-5678</span>{' '}
                (Expired)
              </span>
            </div>
            <button type="submit" className="w-full py-2 px-4 font-semibold bg-primary text-white hover:bg-secondary rounded-lg cursor-pointer">
              Verify Warranty
            </button>
          </form>
        </CardContent>
      </Card>

      {/* Checker Results */}
      <div className="lg:col-span-2 space-y-6">
        {searching ? (
          <Card className="border-slate-200/80 dark:border-slate-800 bg-white/70 dark:bg-slate-900/70 shadow-sm flex flex-col items-center justify-center p-16 text-center border-dashed">
            <p className="text-slate-500 dark:text-slate-400">Searching warranty...</p>
          </Card>
        ) : searchError ? (
          <Card className="border-slate-200/80 dark:border-slate-800 bg-white dark:bg-slate-900 shadow-sm flex flex-col items-center justify-center p-12 text-center">
            <ShieldAlert className="w-12 h-12 text-red-400 mb-4" />
            <h3 className="text-lg font-bold text-slate-900 dark:text-slate-100 mb-1">No Warranty Found</h3>
            <p className="text-slate-400 dark:text-slate-500 text-sm max-w-sm mb-4">
              We couldn't find a warranty contract associated with the ticket number <span className="font-mono font-bold text-slate-900 dark:text-slate-100">"{searchedTicketNum}"</span>.
            </p>
            <p className="text-xs text-slate-400 dark:text-slate-500">{searchError}</p>
          </Card>
        ) : hasSearched ? (
          foundWarranty ? (
            <Card className="border-slate-200/80 dark:border-slate-800 bg-white/70 dark:bg-slate-900/70 backdrop-blur-xl shadow-lg relative overflow-hidden">
              <div className="absolute right-0 top-0 w-32 h-32 bg-primary/5 rounded-full blur-3xl -z-10" />

              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-4 border-b border-slate-100 dark:border-slate-800">
                <div>
                  <CardTitle className="text-xl font-extrabold text-slate-900 dark:text-slate-100">Warranty Contract Details</CardTitle>
                  <CardDescription className="font-mono text-xs text-slate-500 dark:text-slate-400">{foundWarranty.id}</CardDescription>
                </div>
                <WarrantyStatusBadge status={foundWarranty.status} />
              </CardHeader>

              <CardContent className="pt-6 space-y-6 text-sm">
                <div className="grid grid-cols-2 gap-6">
                  <div className="space-y-1">
                    <span className="text-xxs font-bold text-slate-400 dark:text-slate-500 uppercase tracking-wider block">Start Date</span>
                    <span className="font-semibold text-slate-700 dark:text-slate-300">{formatDate(foundWarranty.start_date)}</span>
                  </div>
                  <div className="space-y-1">
                    <span className="text-xxs font-bold text-slate-400 dark:text-slate-500 uppercase tracking-wider block">End Date (Expiration)</span>
                    <span className="font-semibold text-slate-700 dark:text-slate-300">{formatDate(foundWarranty.end_date)}</span>
                  </div>
                  <div className="space-y-1 col-span-2">
                    <span className="text-xxs font-bold text-slate-400 dark:text-slate-500 uppercase tracking-wider block">Original Ticket Reference</span>
                    <span className="font-mono font-bold text-slate-900 dark:text-slate-100">{searchedTicketNum}</span>
                  </div>
                  {foundWarranty.notes && (
                    <div className="space-y-1 col-span-2 p-3 bg-slate-50 dark:bg-slate-800/50 border border-slate-150 dark:border-slate-800 rounded-lg">
                      <span className="text-xxs font-bold text-slate-400 dark:text-slate-500 uppercase tracking-wider block">Status Notes</span>
                      <span className="text-slate-600 dark:text-slate-400 italic">"{foundWarranty.notes}"</span>
                    </div>
                  )}
                </div>

                {/* Action buttons based on status */}
                <div className="flex items-center gap-3 pt-4 border-t border-slate-100 dark:border-slate-800">
                  {foundWarranty.status === 'ACTIVE' ? (
                    <>
                      <Button
                        onClick={onOpenClaimDialog}
                        className="font-semibold bg-primary hover:bg-secondary cursor-pointer"
                      >
                        <ClipboardList className="w-4 h-4 mr-1.5" />
                        Submit Warranty Claim
                      </Button>
                      <Button
                        variant="outline"
                        onClick={() => onOpenVoidDialog(foundWarranty)}
                        className="font-semibold text-red-600 dark:text-red-400 hover:text-red-700 dark:hover:text-red-300 border-red-200 dark:border-red-900/50 hover:bg-red-50 dark:hover:bg-red-950/40 cursor-pointer"
                      >
                        <Ban className="w-4 h-4 mr-1.5" />
                        Void Warranty
                      </Button>
                    </>
                  ) : (
                    <div className="flex items-center gap-2 text-slate-400 dark:text-slate-500 bg-slate-50 dark:bg-slate-800/50 border border-slate-100 dark:border-slate-800 rounded-lg p-3 w-full">
                      <AlertTriangle className="w-5 h-5 text-amber-500 flex-shrink-0" />
                      <span className="text-xs font-semibold text-slate-600 dark:text-slate-300">
                        This warranty is {foundWarranty.status.toLowerCase()}. You cannot submit new claims or modify its status further.
                      </span>
                    </div>
                  )}
                </div>
              </CardContent>
            </Card>
          ) : (
            <Card className="border-slate-200/80 dark:border-slate-800 bg-white dark:bg-slate-900 shadow-sm flex flex-col items-center justify-center p-12 text-center">
              <ShieldAlert className="w-12 h-12 text-red-400 mb-4" />
              <h3 className="text-lg font-bold text-slate-900 dark:text-slate-100 mb-1">No Warranty Found</h3>
              <p className="text-slate-400 dark:text-slate-500 text-sm max-w-sm mb-4">
                We couldn't find a warranty contract associated with the ticket number <span className="font-mono font-bold text-slate-900 dark:text-slate-100">"{searchedTicketNum}"</span>.
              </p>
            </Card>
          )
        ) : (
          <Card className="border-slate-200/80 dark:border-slate-800 bg-white/70 dark:bg-slate-900/70 backdrop-blur-xl shadow-sm flex flex-col items-center justify-center p-16 text-center border-dashed">
            <ShieldCheck className="w-14 h-14 text-slate-300 dark:text-slate-600 mb-4" />
            <h3 className="text-lg font-bold text-slate-900 dark:text-slate-100 mb-1">Awaiting Warranty Check</h3>
            <p className="text-slate-400 dark:text-slate-500 text-sm max-w-sm">
              Enter a ticket number in the panel to check validity, void contracts, or launch a claim workflow.
            </p>
          </Card>
        )}
      </div>
    </div>
  )
}
