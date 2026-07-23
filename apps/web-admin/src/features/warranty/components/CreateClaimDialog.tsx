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
import { ClipboardList } from 'lucide-react'
import type { Warranty } from '@/types/warranty'

interface CreateClaimDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  warranty: Warranty | null
  issueDescription: string
  onIssueChange: (value: string) => void
  onSubmit: (e: React.FormEvent) => void
}

export function CreateClaimDialog({
  open,
  onOpenChange,
  warranty,
  issueDescription,
  onIssueChange,
  onSubmit,
}: CreateClaimDialogProps) {
  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-md">
        <form onSubmit={onSubmit}>
          <DialogHeader>
            <DialogTitle className="text-xl font-extrabold text-slate-900 dark:text-slate-100 flex items-center gap-2">
              <ClipboardList className="w-5 h-5 text-primary" />
              Create Warranty Claim
            </DialogTitle>
            <DialogDescription className="text-slate-500 dark:text-slate-400">
              Register a new claim under warranty <span className="font-mono font-bold text-slate-900 dark:text-slate-100">{warranty?.id}</span>.
            </DialogDescription>
          </DialogHeader>

          <div className="space-y-4 py-4">
            <div className="space-y-1">
              <label className="text-xs font-bold text-slate-500 dark:text-slate-400 uppercase">Issue Description</label>
              <Input
                required
                placeholder="Describe the repeating or new issue..."
                value={issueDescription}
                onChange={(e) => onIssueChange(e.target.value)}
              />
            </div>
          </div>

          <DialogFooter className="gap-2 sm:gap-0 pt-2">
            <Button type="button" variant="outline" className="cursor-pointer" onClick={() => onOpenChange(false)}>Cancel</Button>
            <Button type="submit" className="bg-primary hover:bg-secondary cursor-pointer">Register Claim</Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  )
}
