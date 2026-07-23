import { Badge } from '@/components/ui/badge'
import type { TicketStatus } from '@/types/ticket'

export function TicketStatusBadge({ status }: { status: TicketStatus }) {
  switch (status) {
    case 'RECEIVED':
      return <Badge variant="outline" className="bg-slate-100/50 text-slate-600 border-slate-200 font-semibold">Received</Badge>
    case 'REPAIRING':
      return <Badge className="bg-blue-500/10 text-blue-600 border-none font-semibold hover:bg-blue-500/15">Repairing</Badge>
    case 'PENDING_CONFIRMATION':
      return <Badge className="bg-orange-500/10 text-orange-600 border-none font-semibold hover:bg-orange-500/15">Pending Confirm</Badge>
    case 'FIXED':
      return <Badge className="bg-purple-500/10 text-purple-600 border-none font-semibold hover:bg-purple-500/15">Fixed</Badge>
    case 'COMPLETED':
      return <Badge className="bg-green-500/10 text-green-600 border-none font-semibold hover:bg-green-500/15">Completed</Badge>
    case 'CANCELLED':
      return <Badge className="bg-red-500/10 text-red-600 border-none font-semibold hover:bg-red-500/15">Cancelled</Badge>
    case 'RETURNED':
      return <Badge className="bg-gray-500/10 text-gray-600 border-none font-semibold hover:bg-gray-500/15">Returned</Badge>
    default:
      return <Badge variant="outline">{status}</Badge>
  }
}
