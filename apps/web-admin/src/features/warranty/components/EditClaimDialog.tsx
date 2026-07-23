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
import type { WarrantyClaim } from '@/types/warranty'

interface EditClaimDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  claim: WarrantyClaim | null
  formData: {
    issue_description: string
    notes: string
  }
  onFormDataChange: (data: { issue_description: string; notes: string }) => void
  onSubmit: (e: React.FormEvent) => void
}

export function EditClaimDialog({
  open,
  onOpenChange,
  claim,
  formData,
  onFormDataChange,
  onSubmit,
}: EditClaimDialogProps) {
  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-md">
        <form onSubmit={onSubmit}>
          <DialogHeader>
            <DialogTitle className="text-xl font-extrabold text-slate-900 dark:text-slate-100">Update Claim Info</DialogTitle>
            <DialogDescription className="text-slate-500 dark:text-slate-400">
              Edit claim issue and notes for <span className="font-mono font-bold text-slate-900 dark:text-slate-100">{claim?.claim_number}</span>.
            </DialogDescription>
          </DialogHeader>

          <div className="space-y-4 py-4">
            <div className="space-y-1">
              <label className="text-xs font-bold text-slate-500 dark:text-slate-400 uppercase">Issue Description</label>
              <Input
                required
                value={formData.issue_description}
                onChange={(e) => onFormDataChange({ ...formData, issue_description: e.target.value })}
              />
            </div>

            <div className="space-y-1">
              <label className="text-xs font-bold text-slate-500 dark:text-slate-400 uppercase">Notes</label>
              <Input
                value={formData.notes}
                onChange={(e) => onFormDataChange({ ...formData, notes: e.target.value })}
              />
            </div>
          </div>

          <DialogFooter className="gap-2 sm:gap-0 pt-2">
            <Button type="button" variant="outline" className="cursor-pointer" onClick={() => onOpenChange(false)}>Cancel</Button>
            <Button type="submit" className="bg-primary hover:bg-secondary cursor-pointer">Save Changes</Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  )
}
