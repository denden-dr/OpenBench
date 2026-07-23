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
import { ShieldX } from 'lucide-react'
import type { Warranty } from '@/types/warranty'

interface VoidWarrantyDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  warranty: Warranty | null
  notes: string
  onNotesChange: (value: string) => void
  onSubmit: (e: React.FormEvent) => void
}

export function VoidWarrantyDialog({
  open,
  onOpenChange,
  warranty,
  notes,
  onNotesChange,
  onSubmit,
}: VoidWarrantyDialogProps) {
  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-sm">
        <form onSubmit={onSubmit}>
          <DialogHeader>
            <DialogTitle className="text-xl font-extrabold text-slate-900 dark:text-slate-100 flex items-center gap-2">
              <ShieldX className="w-5 h-5 text-red-500" />
              Void Warranty Contract
            </DialogTitle>
            <DialogDescription className="text-slate-500 dark:text-slate-400">
              Are you sure you want to invalidate warranty <span className="font-mono font-bold text-slate-900 dark:text-slate-100">{warranty?.id}</span>? This action is irreversible.
            </DialogDescription>
          </DialogHeader>

          <div className="space-y-4 py-4">
            <div className="space-y-1">
              <label className="text-xs font-bold text-slate-500 dark:text-slate-400 uppercase">Reason / Void Notes</label>
              <Input
                required
                placeholder="e.g. Broken warranty seal, water damage detected."
                value={notes}
                onChange={(e) => onNotesChange(e.target.value)}
              />
            </div>
          </div>

          <DialogFooter className="gap-2 sm:gap-0 pt-2">
            <Button type="button" variant="outline" className="cursor-pointer" onClick={() => onOpenChange(false)}>Cancel</Button>
            <Button type="submit" className="bg-red-600 hover:bg-red-700 text-white cursor-pointer">Invalidate (Void)</Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  )
}
