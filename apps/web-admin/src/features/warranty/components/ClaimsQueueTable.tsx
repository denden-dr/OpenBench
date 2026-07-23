import { Card, CardContent } from '@/components/ui/card'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { Button } from '@/components/ui/button'
import { ChevronLeft, ChevronRight, Settings } from 'lucide-react'
import type { WarrantyClaim } from '@/types/warranty'
import { ClaimEvalBadge } from './WarrantyStatusBadge'

interface ClaimsQueueTableProps {
  claims: WarrantyClaim[]
  onEvaluate: (claim: WarrantyClaim) => void
  onEdit: (claim: WarrantyClaim) => void
}

export function ClaimsQueueTable({ claims, onEvaluate, onEdit }: ClaimsQueueTableProps) {
  return (
    <Card className="border-slate-200/80 dark:border-slate-800 bg-white dark:bg-slate-900 shadow-sm overflow-hidden">
      <CardContent className="p-0">
        <Table>
          <TableHeader className="bg-slate-50 dark:bg-slate-800/50 border-b border-slate-100 dark:border-slate-800">
            <TableRow className="hover:bg-transparent">
              <TableHead className="w-32 pl-6 font-bold uppercase tracking-wider text-xxs dark:text-slate-400">Claim Num</TableHead>
              <TableHead className="w-24 font-bold uppercase tracking-wider text-xxs dark:text-slate-400">Warranty ID</TableHead>
              <TableHead className="font-bold uppercase tracking-wider text-xxs dark:text-slate-400">Issue Description</TableHead>
              <TableHead className="w-32 font-bold uppercase tracking-wider text-xxs dark:text-slate-400">Evaluation Status</TableHead>
              <TableHead className="w-40 font-bold uppercase tracking-wider text-xxs dark:text-slate-400">Warranty Ticket Ref</TableHead>
              <TableHead className="text-center pr-6 w-36 font-bold uppercase tracking-wider text-xxs dark:text-slate-400">Actions</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody className="text-sm font-medium text-slate-700 dark:text-slate-300 divide-y divide-slate-100/50 dark:divide-slate-800">
            {claims.length > 0 ? (
              claims.map((claim) => (
                <TableRow key={claim.claim_id} className="border-slate-100/50 dark:border-slate-800 hover:bg-slate-50/30 dark:hover:bg-slate-800/50 transition-colors">
                  <TableCell className="pl-6 font-mono text-xs font-bold text-slate-900 dark:text-slate-100">{claim.claim_number}</TableCell>
                  <TableCell className="font-mono text-xs text-slate-500 dark:text-slate-400 font-semibold">{claim.warranty_id}</TableCell>
                  <TableCell className="max-w-xs truncate text-slate-600 dark:text-slate-300 font-semibold" title={claim.issue_description}>
                    {claim.issue_description}
                  </TableCell>
                  <TableCell><ClaimEvalBadge status={claim.evaluation_status} /></TableCell>
                  <TableCell>
                    {claim.warranty_ticket_ref_id ? (
                      <span className="font-mono text-xs font-bold text-blue-600 dark:text-blue-400 flex items-center gap-1">
                        {claim.warranty_ticket_ref_id}
                      </span>
                    ) : (
                      <span className="text-xs text-slate-400 italic">-</span>
                    )}
                  </TableCell>
                  <TableCell className="text-center pr-6">
                    <div className="flex items-center justify-center gap-1">
                      {claim.evaluation_status === 'PENDING' ? (
                        <Button
                          size="xs"
                          className="bg-primary hover:bg-secondary text-xxs font-bold px-2 py-1 h-7 cursor-pointer"
                          onClick={() => onEvaluate(claim)}
                        >
                          Evaluate
                        </Button>
                      ) : (
                        <Button
                          variant="ghost"
                          size="icon-xs"
                          className="h-7 w-7 text-slate-500 dark:text-slate-400 hover:text-primary dark:hover:text-primary hover:bg-slate-100 dark:hover:bg-slate-800 rounded-md cursor-pointer"
                          title="Edit Claim Notes"
                          onClick={() => onEdit(claim)}
                        >
                          <Settings className="w-3.5 h-3.5" />
                        </Button>
                      )}
                    </div>
                  </TableCell>
                </TableRow>
              ))
            ) : (
              <TableRow>
                <TableCell colSpan={6} className="h-32 text-center text-slate-400 dark:text-slate-500">
                  No claims registered.
                </TableCell>
              </TableRow>
            )}
          </TableBody>
        </Table>

        {/* Pagination Footer */}
        <div className="border-t border-slate-100 dark:border-slate-800 px-6 py-4 flex items-center justify-between bg-slate-50/50 dark:bg-slate-800/30">
          <span className="text-xs font-semibold text-slate-500 dark:text-slate-400">
            Showing {claims.length} claim tickets
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
