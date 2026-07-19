import { NavLink, Outlet } from 'react-router-dom'
import { 
  LayoutDashboard, 
  Ticket, 
  Package, 
  CircleDollarSign, 
  ShieldCheck, 
  Store 
} from 'lucide-react'

function AdminLayout() {
  const navItems = [
    { to: '/dashboard', label: 'Dashboard Overview', icon: LayoutDashboard },
    { to: '/tickets', label: 'Repair Tickets', icon: Ticket },
    { to: '/inventory', label: 'Product Inventory', icon: Package },
    { to: '/pos', label: 'Point of Sale (POS)', icon: CircleDollarSign },
    { to: '/warranties', label: 'Warranty Claims', icon: ShieldCheck },
  ]

  return (
    <div className="min-h-screen bg-slate-50 flex flex-col text-slate-900 selection:bg-primary selection:text-white">
      {/* Top Header */}
      <header className="sticky top-0 z-40 bg-white/80 backdrop-blur-md border-b border-slate-200/80">
        <div className="mx-auto px-8 h-16 flex items-center justify-between">
          <div className="flex items-center gap-2">
            <div className="w-8 h-8 rounded-lg bg-primary flex items-center justify-center text-white font-bold text-lg shadow-md animate-pulse">
              OB
            </div>
            <span className="font-extrabold text-xl tracking-tight text-primary">
              OpenBench <span className="text-xs font-semibold px-2 py-0.5 rounded-full bg-primary/10 text-primary uppercase tracking-wider ml-1">Admin</span>
            </span>
          </div>
          
          <div className="flex items-center gap-6">
            <div className="flex items-center gap-2 text-xs font-semibold text-slate-500">
              <Store className="w-4 h-4 text-slate-400" />
              <span>Store: Main Branch</span>
            </div>
            <div className="w-8 h-8 rounded-full bg-slate-100 flex items-center justify-center font-bold text-slate-700 text-sm border border-slate-200 hover:bg-slate-200 transition-colors cursor-pointer">
              AD
            </div>
          </div>
        </div>
      </header>

      {/* Main Container */}
      <div className="flex-grow flex flex-col md:flex-row">
        {/* Sidebar */}
        <aside className="w-full md:w-64 bg-white border-r border-slate-200 p-6 space-y-8 shrink-0">
          <div className="space-y-2">
            <p className="text-xxs font-bold text-slate-400 uppercase tracking-wider px-3">Main Navigation</p>
            <nav className="space-y-1">
              {navItems.map((item) => {
                const IconComponent = item.icon
                return (
                  <NavLink
                    key={item.to}
                    to={item.to}
                    className={({ isActive }) =>
                      `w-full flex items-center gap-2.5 px-3 py-2 rounded-lg text-sm font-medium transition-colors ${
                        isActive 
                          ? 'bg-primary/10 text-primary' 
                          : 'text-slate-600 hover:bg-slate-50 hover:text-slate-950'
                      }`
                    }
                  >
                    <IconComponent className="w-4 h-4 shrink-0" />
                    <span>{item.label}</span>
                  </NavLink>
                )
              })}
            </nav>
          </div>
        </aside>

        {/* Workspace */}
        <main className="flex-grow p-8 bg-slate-50/50 overflow-x-hidden">
          <Outlet />
        </main>
      </div>
    </div>
  )
}

export default AdminLayout
