import { useState } from 'react'

function App() {
  const [activeTab, setActiveTab] = useState('overview')

  const metrics = [
    { title: 'Active Tickets', value: '14', change: '+2 new', color: 'bg-primary' },
    { title: 'Pending Diagnoses', value: '5', change: 'needs review', color: 'bg-tertiary' },
    { title: 'Sales Today', value: '$849.50', change: '+12% from yesterday', color: 'bg-secondary' },
    { title: 'Active Warranties', value: '142', change: '99% compliance', color: 'bg-accent text-slate-800' },
  ]

  const recentTickets = [
    { id: 'TKT-89304', customer: 'John Doe', device: 'iPhone 13 Pro', status: 'In Progress', cost: '$189.00' },
    { id: 'TKT-89303', customer: 'Jane Smith', device: 'MacBook Air M1', status: 'Ready for Pickup', cost: '$249.00' },
    { id: 'TKT-89302', customer: 'Alex Johnson', device: 'iPad Pro', status: 'Waiting for Parts', cost: '$95.00' },
    { id: 'TKT-89301', customer: 'Sarah Connor', device: 'Sony WH-1000XM4', status: 'Diagnosed', cost: '$75.00' },
  ]

  return (
    <div className="min-h-screen bg-slate-50 flex flex-col text-slate-900 selection:bg-primary selection:text-white">
      {/* Top Header */}
      <header className="sticky top-0 z-50 bg-white/80 backdrop-blur-md border-b border-slate-200/80">
        <div className="mx-auto px-8 h-16 flex items-center justify-between">
          <div className="flex items-center gap-2">
            <div className="w-8 h-8 rounded-lg bg-primary flex items-center justify-center text-white font-bold text-lg shadow-md">
              OB
            </div>
            <span className="font-extrabold text-xl tracking-tight text-primary">OpenBench <span className="text-xs font-semibold px-2 py-0.5 rounded-full bg-primary/10 text-primary uppercase tracking-wider ml-1">Admin</span></span>
          </div>
          
          <div className="flex items-center gap-4">
            <span className="text-xs font-semibold text-slate-500">Store: Main Branch</span>
            <div className="w-8 h-8 rounded-full bg-slate-200 flex items-center justify-center font-bold text-slate-700 text-sm border border-slate-300">
              AD
            </div>
          </div>
        </div>
      </header>

      {/* Main Container */}
      <div className="flex-grow flex flex-col md:flex-row">
        {/* Sidebar */}
        <aside className="w-full md:w-64 bg-white border-r border-slate-200 p-6 space-y-8">
          <div className="space-y-1">
            <p className="text-xxs font-bold text-slate-400 uppercase tracking-wider px-3">Main Navigation</p>
            <nav className="space-y-1">
              <button
                onClick={() => setActiveTab('overview')}
                className={`w-full text-left px-3 py-2 rounded-lg text-sm font-medium transition-colors ${
                  activeTab === 'overview' ? 'bg-primary/10 text-primary' : 'text-slate-600 hover:bg-slate-100 hover:text-slate-950'
                }`}
              >
                📊 Dashboard Overview
              </button>
              <button
                onClick={() => setActiveTab('tickets')}
                className={`w-full text-left px-3 py-2 rounded-lg text-sm font-medium transition-colors ${
                  activeTab === 'tickets' ? 'bg-primary/10 text-primary' : 'text-slate-600 hover:bg-slate-100 hover:text-slate-950'
                }`}
              >
                🎫 Repair Tickets
              </button>
              <button
                onClick={() => setActiveTab('inventory')}
                className={`w-full text-left px-3 py-2 rounded-lg text-sm font-medium transition-colors ${
                  activeTab === 'inventory' ? 'bg-primary/10 text-primary' : 'text-slate-600 hover:bg-slate-100 hover:text-slate-950'
                }`}
              >
                📦 Product Inventory
              </button>
              <button
                onClick={() => setActiveTab('pos')}
                className={`w-full text-left px-3 py-2 rounded-lg text-sm font-medium transition-colors ${
                  activeTab === 'pos' ? 'bg-primary/10 text-primary' : 'text-slate-600 hover:bg-slate-100 hover:text-slate-950'
                }`}
              >
                💰 Point of Sale (POS)
              </button>
              <button
                onClick={() => setActiveTab('warranties')}
                className={`w-full text-left px-3 py-2 rounded-lg text-sm font-medium transition-colors ${
                  activeTab === 'warranties' ? 'bg-primary/10 text-primary' : 'text-slate-600 hover:bg-slate-100 hover:text-slate-950'
                }`}
              >
                🛡️ Warranty Claims
              </button>
            </nav>
          </div>
        </aside>

        {/* Workspace */}
        <main className="flex-grow p-8 space-y-8 bg-slate-50/50">
          <div>
            <h1 className="text-3xl font-extrabold text-slate-900 tracking-tight">
              Dashboard Overview
            </h1>
            <p className="text-slate-500 text-sm">Welcome back. Here is what is happening at the workshop today.</p>
          </div>

          {/* Metrics grid */}
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
            {metrics.map((metric, idx) => (
              <div key={idx} className="p-6 rounded-2xl bg-white shadow-sm border border-slate-200/80 flex flex-col justify-between h-32 hover:shadow-md transition-shadow">
                <span className="text-xs font-bold text-slate-500 uppercase tracking-wider">{metric.title}</span>
                <div className="flex items-baseline justify-between mt-2">
                  <span className="text-3xl font-extrabold text-slate-900 tracking-tight">{metric.value}</span>
                  <span className="text-xxs font-semibold text-slate-500 bg-slate-100 px-2 py-0.5 rounded-full">{metric.change}</span>
                </div>
              </div>
            ))}
          </div>

          {/* Table container */}
          <div className="p-6 rounded-2xl bg-white/70 backdrop-blur-xl border border-slate-200/80 shadow-sm space-y-4">
            <div className="flex justify-between items-center pb-2">
              <div>
                <h3 className="font-extrabold text-slate-800 text-lg">Recent Tickets</h3>
                <p className="text-slate-500 text-xs">A list of the 4 most recently updated tickets.</p>
              </div>
              <button className="px-3 py-1.5 text-xs font-semibold text-primary hover:bg-primary/5 border border-primary/20 rounded-lg transition-colors">
                View All Tickets
              </button>
            </div>

            <div className="overflow-x-auto">
              <table className="w-full text-left border-collapse">
                <thead>
                  <tr className="border-b border-slate-100 text-slate-400 text-xxs font-bold uppercase tracking-wider">
                    <th className="pb-3 pl-2">ID</th>
                    <th className="pb-3">Customer</th>
                    <th className="pb-3">Device</th>
                    <th className="pb-3">Status</th>
                    <th className="pb-3 text-right pr-2">Cost</th>
                  </tr>
                </thead>
                <tbody className="text-sm font-medium text-slate-700 divide-y divide-slate-100/50">
                  {recentTickets.map((t, idx) => (
                    <tr key={idx} className="hover:bg-slate-50/50 transition-colors">
                      <td className="py-3 pl-2 font-mono text-xs font-bold text-slate-500">{t.id}</td>
                      <td className="py-3">{t.customer}</td>
                      <td className="py-3 font-semibold text-slate-900">{t.device}</td>
                      <td className="py-3">
                        <span className={`px-2 py-0.5 rounded-full text-xxs font-bold ${
                          t.status === 'Ready for Pickup' ? 'bg-green-100 text-green-700' :
                          t.status === 'In Progress' ? 'bg-blue-100 text-blue-700' :
                          t.status === 'Waiting for Parts' ? 'bg-orange-100 text-orange-700' : 'bg-slate-100 text-slate-700'
                        }`}>
                          {t.status}
                        </span>
                      </td>
                      <td className="py-3 text-right pr-2 font-mono font-bold text-slate-900">{t.cost}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </div>
        </main>
      </div>
    </div>
  )
}

export default App
