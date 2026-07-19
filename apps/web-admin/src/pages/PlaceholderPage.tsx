import { ShieldAlert } from 'lucide-react'

interface PlaceholderPageProps {
  title: string
  description?: string
}

function PlaceholderPage({ title, description }: PlaceholderPageProps) {
  return (
    <div className="flex flex-col items-center justify-center min-h-[50vh] text-center p-8 space-y-4 bg-white/50 border border-dashed border-slate-200 rounded-2xl">
      <div className="w-12 h-12 rounded-full bg-primary/10 flex items-center justify-center text-primary">
        <ShieldAlert className="w-6 h-6" />
      </div>
      <div className="space-y-2">
        <h1 className="text-2xl font-extrabold text-slate-900 tracking-tight">{title} Shell</h1>
        <p className="text-slate-500 text-sm max-w-sm">
          {description || `The ${title} management interface shell is active. Backend API integration will be implemented next.`}
        </p>
      </div>
    </div>
  )
}

export default PlaceholderPage
