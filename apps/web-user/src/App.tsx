import { useState } from 'react'

function App() {
  const [ticketId, setTicketId] = useState('')

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault()
    alert(`Searching for ticket: ${ticketId}`)
  }

  return (
    <div className="min-h-screen bg-slate-50 flex flex-col justify-between text-slate-900 selection:bg-primary selection:text-white">
      {/* Header */}
      <header className="sticky top-0 z-50 bg-white/70 backdrop-blur-md border-b border-slate-200/80">
        <div className="max-w-7xl mx-auto px-6 h-16 flex items-center justify-between">
          <div className="flex items-center gap-2">
            <div className="w-8 h-8 rounded-lg bg-primary flex items-center justify-center text-white font-bold text-lg shadow-md">
              OB
            </div>
            <span className="font-extrabold text-xl tracking-tight text-primary">OpenBench</span>
          </div>
          <nav className="flex items-center gap-6">
            <a href="#features" className="text-sm font-medium text-slate-600 hover:text-primary transition-colors">Features</a>
            <a href="#track" className="text-sm font-medium text-slate-600 hover:text-primary transition-colors">Track Repair</a>
            <button className="px-4 py-2 text-sm font-semibold text-white bg-primary hover:bg-secondary rounded-lg shadow-sm transition-all duration-200">
              New Ticket
            </button>
          </nav>
        </div>
      </header>

      {/* Hero Section */}
      <main className="flex-grow max-w-7xl mx-auto px-6 py-16 flex flex-col lg:flex-row items-center gap-12">
        <div className="flex-1 space-y-6 text-center lg:text-left">
          <div className="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-primary/10 text-primary text-xs font-semibold">
            <span>✨</span> Modern Repair Management System
          </div>
          <h1 className="text-4xl sm:text-5xl lg:text-6xl font-extrabold text-slate-900 tracking-tight leading-none">
            Reliable electronics repair, <span className="text-primary">tracked in real-time</span>.
          </h1>
          <p className="text-lg text-slate-600 max-w-xl mx-auto lg:mx-0">
            OpenBench bridges the gap between hardware repair shops and customers. Check your repair status, access service records, and verify warranties instantly.
          </p>

          <form onSubmit={handleSearch} className="max-w-md mx-auto lg:mx-0 flex gap-2">
            <input
              type="text"
              placeholder="Enter Ticket ID (e.g. TKT-12345)"
              value={ticketId}
              onChange={(e) => setTicketId(e.target.value)}
              className="flex-grow px-4 py-3 rounded-lg border border-slate-200 bg-white shadow-sm focus:outline-none focus:border-primary text-sm"
            />
            <button
              type="submit"
              className="px-5 py-3 bg-primary hover:bg-secondary text-white font-semibold rounded-lg text-sm shadow-md transition-all duration-200"
            >
              Track Status
            </button>
          </form>
        </div>

        {/* Visual Element */}
        <div className="flex-1 w-full max-w-md lg:max-w-none">
          <div className="relative p-8 rounded-2xl bg-white/70 backdrop-blur-xl border border-white/40 shadow-2xl space-y-6 overflow-hidden">
            <div className="absolute top-0 right-0 w-32 h-32 bg-accent/20 rounded-full blur-3xl -z-10"></div>
            <div className="absolute bottom-0 left-0 w-32 h-32 bg-primary/10 rounded-full blur-3xl -z-10"></div>

            <div className="flex justify-between items-center pb-4 border-b border-slate-100">
              <div>
                <p className="text-xs text-slate-500 font-medium">TICKET DETAILS</p>
                <p className="text-sm font-bold text-slate-800">TKT-89304</p>
              </div>
              <span className="px-2.5 py-1 rounded-full bg-tertiary/10 text-tertiary text-xs font-semibold">
                In Progress
              </span>
            </div>

            <div className="space-y-4">
              <div>
                <p className="text-xs text-slate-500">Device</p>
                <p className="text-sm font-semibold text-slate-800">iPhone 13 Pro Max</p>
              </div>
              <div>
                <p className="text-xs text-slate-500">Issue Reported</p>
                <p className="text-sm font-semibold text-slate-800">Cracked screen, battery replacement</p>
              </div>
              <div>
                <p className="text-xs text-slate-500">Estimated Cost</p>
                <p className="text-sm font-semibold text-slate-800">$189.00</p>
              </div>
            </div>

            <div className="pt-2">
              <div className="w-full bg-slate-100 h-2 rounded-full overflow-hidden">
                <div className="bg-primary h-full rounded-full" style={{ width: '60%' }}></div>
              </div>
              <div className="flex justify-between text-xs text-slate-500 mt-2 font-medium">
                <span>Received</span>
                <span className="text-primary font-semibold">Diagnosis</span>
                <span>Ready for Pickup</span>
              </div>
            </div>
          </div>
        </div>
      </main>

      {/* Footer */}
      <footer className="bg-slate-900 text-slate-400 py-8 border-t border-slate-800">
        <div className="max-w-7xl mx-auto px-6 flex flex-col md:flex-row items-center justify-between gap-4 text-sm">
          <p>&copy; {new Date().getFullYear()} OpenBench. All rights reserved.</p>
          <div className="flex gap-6">
            <a href="#" className="hover:text-white transition-colors">Privacy Policy</a>
            <a href="#" className="hover:text-white transition-colors">Terms of Service</a>
            <a href="#" className="hover:text-white transition-colors">Contact</a>
          </div>
        </div>
      </footer>
    </div>
  )
}

export default App
