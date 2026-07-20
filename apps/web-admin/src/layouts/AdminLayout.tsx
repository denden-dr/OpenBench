import { NavLink, Outlet, useNavigate } from 'react-router-dom'
import { 
  LayoutDashboard, 
  Ticket, 
  Package, 
  CircleDollarSign, 
  ShieldCheck, 
  Store,
  LogOut
} from 'lucide-react'
import { useAuthStore } from '@/stores/authStore'
import { ThemeToggle } from '@/components/ThemeToggle'

function AdminLayout() {
  const { user, logout } = useAuthStore()
  const navigate = useNavigate()

  const handleLogout = async () => {
    await logout()
    navigate('/login')
  }

  const userInitial = user?.email ? user.email[0].toUpperCase() : 'A'

  const navItems = [
    { to: '/dashboard', label: 'Dashboard Overview', icon: LayoutDashboard },
    { to: '/tickets', label: 'Repair Tickets', icon: Ticket },
    { to: '/inventory', label: 'Product Inventory', icon: Package },
    { to: '/pos', label: 'Point of Sale (POS)', icon: CircleDollarSign },
    { to: '/warranties', label: 'Warranty Claims', icon: ShieldCheck },
  ]

  return (
    <div className="min-h-screen bg-slate-50 dark:bg-slate-950 flex flex-col text-slate-900 dark:text-slate-100 selection:bg-primary selection:text-white transition-colors duration-200">
      {/* Top Header */}
      <header className="sticky top-0 z-40 bg-white/80 dark:bg-slate-900/80 backdrop-blur-md border-b border-slate-200/80 dark:border-slate-800">
        <div className="mx-auto px-8 h-16 flex items-center justify-between">
          <div className="flex items-center gap-2">
            <div className="w-8 h-8 rounded-lg bg-primary flex items-center justify-center text-white font-bold text-lg shadow-md">
              OB
            </div>
            <span className="font-extrabold text-xl tracking-tight text-primary dark:text-white">
              OpenBench <span className="text-xs font-semibold px-2 py-0.5 rounded-full bg-primary/10 dark:bg-primary/20 text-primary dark:text-accent uppercase tracking-wider ml-1">Admin</span>
            </span>
          </div>
          
          <div className="flex items-center gap-4 sm:gap-6">
            <div className="hidden md:flex items-center gap-2 text-xs font-semibold text-slate-500 dark:text-slate-400">
              <Store className="w-4 h-4 text-slate-400 dark:text-slate-500" />
              <span>Store: Main Branch</span>
            </div>

            <div className="flex items-center gap-3 pl-4 border-l border-slate-200 dark:border-slate-800">
              <ThemeToggle />

              <div className="flex items-center gap-2.5">
                <div className="w-8 h-8 rounded-full bg-primary/10 text-primary dark:bg-primary/20 dark:text-accent flex items-center justify-center font-bold text-sm border border-primary/20 dark:border-primary/30">
                  {userInitial}
                </div>
                <div className="hidden sm:block text-left">
                  <p className="text-xs font-semibold text-slate-800 dark:text-slate-200 truncate max-w-[150px]">{user?.email || 'Admin User'}</p>
                  <p className="text-[10px] font-medium text-slate-500 dark:text-slate-400 uppercase tracking-wider">{user?.role || 'ADMIN'}</p>
                </div>
              </div>

              <button
                onClick={handleLogout}
                title="Sign Out"
                className="p-1.5 rounded-lg text-slate-400 hover:text-destructive hover:bg-destructive/10 transition-colors cursor-pointer"
              >
                <LogOut className="w-4 h-4" />
              </button>
            </div>
          </div>
        </div>
      </header>

      {/* Main Container */}
      <div className="flex-grow flex flex-col md:flex-row">
        {/* Sidebar */}
        <aside className="w-full md:w-64 bg-white dark:bg-slate-900 border-r border-slate-200 dark:border-slate-800 p-6 space-y-8 shrink-0">
          <div className="space-y-2">
            <p className="text-xxs font-bold text-slate-400 dark:text-slate-500 uppercase tracking-wider px-3">Main Navigation</p>
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
                          ? 'bg-primary/10 text-primary dark:bg-primary/20 dark:text-accent font-semibold' 
                          : 'text-slate-600 dark:text-slate-400 hover:bg-slate-100 dark:hover:bg-slate-800 hover:text-slate-950 dark:hover:text-white'
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
        <main className="flex-grow p-8 bg-slate-50/50 dark:bg-slate-950/50 overflow-x-hidden">
          <Outlet />
        </main>
      </div>
    </div>
  )
}

export default AdminLayout
