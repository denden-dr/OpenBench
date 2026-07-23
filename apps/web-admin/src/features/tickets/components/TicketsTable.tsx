import type { TicketListItem } from '@/types/ticket'
import { Card, CardContent } from '@/components/ui/card'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { Button } from '@/components/ui/button'
import { Eye, Edit, ShieldAlert, Loader2 } from 'lucide-react'
import { TicketStatusBadge } from './TicketStatusBadge'

interface TicketsTableProps {
  tickets: TicketListItem[]
  loading: boolean
  onViewDetails: (ticketId: string) => void
  onUpdateStatus: (ticket: TicketListItem) => void
}

export function TicketsTable({ tickets, loading, onViewDetails, onUpdateStatus }: TicketsTableProps) {
  const formatDate = (isoString: string) => {
    if (!isoString) return '-'
    return new Date(isoString).toLocaleDateString('id-ID', {
      day: 'numeric',
      month: 'short',
      year: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    })
  }

  return (
    <Card className="border-slate-200/80 dark:border-slate-800 bg-white dark:bg-slate-900 shadow-sm overflow-hidden">
      <CardContent className="p-0">
        <Table>
          <TableHeader className="bg-slate-50 dark:bg-slate-950/50 border-b border-slate-100 dark:border-slate-800">
            <TableRow className="hover:bg-transparent">
              <TableHead className="w-36 pl-6 font-bold uppercase tracking-wider text-xxs text-slate-500 dark:text-slate-400">Ticket Number</TableHead>
              <TableHead className="font-bold uppercase tracking-wider text-xxs text-slate-500 dark:text-slate-400">Date Registered</TableHead>
              <TableHead className="font-bold uppercase tracking-wider text-xxs text-slate-500 dark:text-slate-400">Customer Name</TableHead>
              <TableHead className="font-bold uppercase tracking-wider text-xxs text-slate-500 dark:text-slate-400">Device Details</TableHead>
              <TableHead className="font-bold uppercase tracking-wider text-xxs text-slate-500 dark:text-slate-400">Status</TableHead>
              <TableHead className="text-center pr-6 font-bold uppercase tracking-wider text-xxs text-slate-500 dark:text-slate-400">Actions</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody className="text-sm font-medium text-slate-700 dark:text-slate-300 divide-y divide-slate-100/50 dark:divide-slate-800">
            {loading ? (
              <TableRow>
                <TableCell colSpan={6} className="h-32 text-center text-slate-400">
                  <div className="flex flex-col items-center justify-center gap-2">
                    <Loader2 className="w-6 h-6 animate-spin text-primary" />
                    <span>Memuat daftar tiket servis...</span>
                  </div>
                </TableCell>
              </TableRow>
            ) : tickets.length > 0 ? (
              tickets.map((t) => (
                <TableRow key={t.ticket_id} className="border-slate-100/50 dark:border-slate-800/50 hover:bg-slate-50/30 dark:hover:bg-slate-800/30 transition-colors">
                  <TableCell className="pl-6 font-mono text-xs font-bold text-slate-600 dark:text-slate-400">{t.ticket_number}</TableCell>
                  <TableCell className="text-slate-500 dark:text-slate-400 font-semibold">{formatDate(t.created_at)}</TableCell>
                  <TableCell className="font-semibold text-slate-900 dark:text-slate-100">{t.customer_name}</TableCell>
                  <TableCell>
                    <span className="font-bold text-slate-900 dark:text-slate-100">{t.device_brand}</span>{' '}
                    <span className="text-slate-500 dark:text-slate-400">{t.device_model}</span>
                  </TableCell>
                  <TableCell>
                    <TicketStatusBadge status={t.status} />
                  </TableCell>
                  <TableCell className="text-center pr-6">
                    <div className="flex items-center justify-center gap-1.5">
                      <Button 
                        variant="ghost" 
                        size="icon-sm" 
                        className="h-7 w-7 text-slate-500 dark:text-slate-400 hover:text-primary hover:bg-slate-100 dark:hover:bg-slate-800 rounded-md cursor-pointer" 
                        title="View details"
                        onClick={() => onViewDetails(t.ticket_id)}
                      >
                        <Eye className="w-4 h-4" />
                      </Button>
                      <Button 
                        variant="ghost" 
                        size="icon-sm" 
                        className="h-7 w-7 text-slate-500 dark:text-slate-400 hover:text-tertiary hover:bg-slate-100 dark:hover:bg-slate-800 rounded-md cursor-pointer" 
                        title="Update status"
                        onClick={() => onUpdateStatus(t)}
                      >
                        <Edit className="w-4 h-4" />
                      </Button>
                    </div>
                  </TableCell>
                </TableRow>
              ))
            ) : (
              <TableRow>
                <TableCell colSpan={6} className="h-32 text-center text-slate-400">
                  <div className="flex flex-col items-center justify-center gap-1">
                    <ShieldAlert className="w-6 h-6 text-slate-300 dark:text-slate-600" />
                    <span>No tickets found matching your query</span>
                  </div>
                </TableCell>
              </TableRow>
            )}
          </TableBody>
        </Table>
      </CardContent>
    </Card>
  )
}
