import { ShieldAlert } from 'lucide-react'

interface PlaceholderPageProps {
  title: string
  description?: string
}

function PlaceholderPage({ title, description }: PlaceholderPageProps) {
  return (
    <div className="flex flex-col items-center justify-center min-h-[50vh] text-center p-8 space-y-4 bg-white/50 dark:bg-slate-900/50 border border-dashed border-slate-200 dark:border-slate-800 rounded-2xl">
      <div className="w-12 h-12 rounded-full bg-primary/10 dark:bg-primary/20 flex items-center justify-center text-primary dark:text-accent">
        <ShieldAlert className="w-6 h-6" />
      </div>
      <div className="space-y-2">
        <h1 className="text-2xl font-extrabold text-slate-900 dark:text-slate-100 tracking-tight">{title} Shell</h1>
        <p className="text-slate-500 dark:text-slate-400 text-sm max-w-sm">
          {description || `The ${title} management interface shell is active. Backend API integration will be implemented next.`}
        </p>
      </div>
    </div>
  )
}

export default PlaceholderPage
