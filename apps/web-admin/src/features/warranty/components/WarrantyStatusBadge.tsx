import { Badge } from '@/components/ui/badge'
import type { WarrantyStatus, ClaimEvaluationStatus } from '@/types/warranty'

export function WarrantyStatusBadge({ status }: { status: WarrantyStatus }) {
  switch (status) {
    case 'ACTIVE':
      return <Badge className="bg-green-500/10 text-green-600 border-none font-semibold px-2.5 py-1">ACTIVE</Badge>
    case 'VOID':
      return <Badge className="bg-red-500/10 text-red-600 border-none font-semibold px-2.5 py-1">VOIDED</Badge>
    case 'EXPIRED':
      return <Badge className="bg-slate-500/10 text-slate-600 border-none font-semibold px-2.5 py-1">EXPIRED</Badge>
  }
}

export function ClaimEvalBadge({ status }: { status: ClaimEvaluationStatus }) {
  switch (status) {
    case 'PENDING':
      return <Badge className="bg-amber-500/10 text-amber-600 border-none font-semibold">PENDING</Badge>
    case 'ACCEPTED':
      return <Badge className="bg-green-500/10 text-green-600 border-none font-semibold">ACCEPTED</Badge>
    case 'REJECTED':
      return <Badge className="bg-red-500/10 text-red-600 border-none font-semibold">REJECTED</Badge>
    case 'VOID':
      return <Badge className="bg-rose-950/20 text-rose-800 border-none font-semibold">VOID</Badge>
  }
}
