import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from '@/components/ui/dialog'
import type { WarrantyClaim, ClaimEvaluationStatus } from '@/types/warranty'

interface EvaluateClaimDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  claim: WarrantyClaim | null
  evalStatus: ClaimEvaluationStatus
  onEvalStatusChange: (status: ClaimEvaluationStatus) => void
  evalNotes: string
  onEvalNotesChange: (value: string) => void
  onSubmit: (e: React.FormEvent) => void
}

export function EvaluateClaimDialog({
  open,
  onOpenChange,
  claim,
  evalStatus,
  onEvalStatusChange,
  evalNotes,
  onEvalNotesChange,
  onSubmit,
}: EvaluateClaimDialogProps) {
  const requiresNotes = evalStatus === 'REJECTED' || evalStatus === 'VOID'

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-md">
        <form onSubmit={onSubmit}>
          <DialogHeader>
            <DialogTitle className="text-xl font-extrabold text-slate-900 dark:text-slate-100">Evaluate Warranty Claim</DialogTitle>
            <DialogDescription className="text-slate-500 dark:text-slate-400">
              Decide whether to accept, reject, or void claim <span className="font-mono font-bold text-slate-900 dark:text-slate-100">{claim?.claim_number}</span>.
            </DialogDescription>
          </DialogHeader>

          <div className="space-y-4 py-4">
            <div className="space-y-1">
              <label className="text-xs font-bold text-slate-500 dark:text-slate-400 uppercase">Decision</label>
              <select
                className="w-full bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 text-slate-900 dark:text-slate-100 rounded-md p-2 text-sm focus:outline-none focus:ring-2 focus:ring-primary/20"
                value={evalStatus}
                onChange={(e) => onEvalStatusChange(e.target.value as ClaimEvaluationStatus)}
              >
                <option value="ACCEPTED">ACCEPTED (Approve under warranty)</option>
                <option value="REJECTED">REJECTED (Decline claim)</option>
                <option value="VOID">VOID (Violates terms, invalidate contract)</option>
              </select>
            </div>

            <div className="space-y-1">
              <label className="text-xs font-bold text-slate-500 dark:text-slate-400 uppercase">
                Notes {requiresNotes && <span className="text-red-500">*</span>}
              </label>
              <Input
                required={requiresNotes}
                placeholder={evalStatus === 'ACCEPTED' ? 'e.g. Free parts replacement' : 'Reason is required...'}
                value={evalNotes}
                onChange={(e) => onEvalNotesChange(e.target.value)}
              />
              {evalStatus === 'VOID' && (
                <p className="text-xxs text-amber-600 dark:text-amber-400 font-semibold">
                  * WARNING: Setting status to VOID will automatically set the parent warranty contract status to VOID as well.
                </p>
              )}
            </div>
          </div>

          <DialogFooter className="gap-2 sm:gap-0 pt-2">
            <Button type="button" variant="outline" className="cursor-pointer" onClick={() => onOpenChange(false)}>Cancel</Button>
            <Button type="submit" className="bg-primary hover:bg-secondary cursor-pointer">Submit Decision</Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  )
}
