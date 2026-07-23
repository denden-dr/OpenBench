import { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { Card, CardHeader, CardTitle, CardContent, CardDescription } from '@/components/ui/card'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { Button } from '@/components/ui/button'
import { ArrowUpRight } from 'lucide-react'
import { dashboardService, type DashboardData } from '@/services/dashboardService'
import { TicketStatusBadge } from '@/features/tickets/components/TicketStatusBadge'

function DashboardPage() {
  const navigate = useNavigate()
  const [data, setData] = useState<DashboardData | null>(null)
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    async function loadDashboard() {
      try {
        setIsLoading(true)
        const dbData = await dashboardService.getDashboard()
        setData(dbData)
      } catch (err: any) {
        console.error('Failed to load dashboard data', err)
        setError('Failed to fetch dashboard data. Please try again later.')
      } finally {
        setIsLoading(false)
      }
    }
    loadDashboard()
  }, [])

  const formatPrice = (price: number) => {
    return new Intl.NumberFormat('id-ID', {
      style: 'currency',
      currency: 'IDR',
      minimumFractionDigits: 0,
    }).format(price)
  }

  const renderMetrics = () => {
    if (isLoading || !data) {
      return Array.from({ length: 4 }).map((_, idx) => (
        <Card key={idx} className="border-slate-200/80 dark:border-slate-800 bg-white dark:bg-slate-900 animate-pulse">
          <CardHeader className="pb-2">
            <div className="h-3 w-24 bg-slate-200 dark:bg-slate-850 rounded"></div>
          </CardHeader>
          <CardContent className="flex items-baseline justify-between pt-2">
            <div className="h-8 w-16 bg-slate-200 dark:bg-slate-850 rounded"></div>
            <div className="h-4 w-12 bg-slate-200 dark:bg-slate-850 rounded-full"></div>
          </CardContent>
        </Card>
      ))
    }

    const { metrics } = data
    const metricItems = [
      { title: 'Active Tickets', value: String(metrics.active_tickets), change: 'in shop', color: 'text-primary dark:text-accent' },
      { title: 'Pending Diagnoses', value: String(metrics.pending_diagnoses), change: 'needs review', color: 'text-tertiary' },
      { title: 'Sales Today', value: formatPrice(metrics.sales_today), change: 'today', color: 'text-secondary dark:text-slate-300' },
      { title: 'Active Warranties', value: String(metrics.active_warranties), change: 'covered', color: 'text-green-600 dark:text-green-400' },
    ]

    return metricItems.map((metric, idx) => (
      <Card key={idx} className="hover:shadow-md transition-shadow border-slate-200/80 dark:border-slate-800 bg-white dark:bg-slate-900">
        <CardHeader className="flex flex-row items-center justify-between pb-2">
          <CardTitle className="text-xs font-bold text-slate-500 dark:text-slate-400 uppercase tracking-wider">{metric.title}</CardTitle>
        </CardHeader>
        <CardContent className="flex items-baseline justify-between">
          <span className="text-2xl font-extrabold text-slate-900 dark:text-slate-100 tracking-tight">{metric.value}</span>
          <span className={`text-xxs font-semibold px-2 py-0.5 rounded-full bg-slate-100 dark:bg-slate-800 ${metric.color}`}>{metric.change}</span>
        </CardContent>
      </Card>
    ))
  }

  const renderRecentTickets = () => {
    if (isLoading || !data) {
      return Array.from({ length: 4 }).map((_, idx) => (
        <TableRow key={idx} className="animate-pulse">
          <TableCell className="pl-2"><div className="h-4 w-20 bg-slate-200 dark:bg-slate-850 rounded"></div></TableCell>
          <TableCell><div className="h-4 w-24 bg-slate-200 dark:bg-slate-850 rounded"></div></TableCell>
          <TableCell><div className="h-4 w-32 bg-slate-200 dark:bg-slate-850 rounded"></div></TableCell>
          <TableCell><div className="h-6 w-16 bg-slate-200 dark:bg-slate-850 rounded-full"></div></TableCell>
          <TableCell className="text-right pr-2"><div className="h-4 w-16 bg-slate-200 dark:bg-slate-850 rounded ml-auto"></div></TableCell>
        </TableRow>
      ))
    }

    if (!data.recent_tickets || data.recent_tickets.length === 0) {
      return (
        <TableRow>
          <TableCell colSpan={5} className="text-center py-6 text-slate-500 dark:text-slate-400">
            No recent tickets found.
          </TableCell>
        </TableRow>
      )
    }

    return data.recent_tickets.map((t) => (
      <TableRow key={t.ticket_id} className="border-slate-100/50 dark:border-slate-800/50 hover:bg-slate-50/50 dark:hover:bg-slate-800/50 transition-colors">
        <TableCell className="pl-2 font-mono text-xs font-bold text-slate-500 dark:text-slate-400">{t.ticket_number}</TableCell>
        <TableCell>{t.customer_name}</TableCell>
        <TableCell className="font-semibold text-slate-900 dark:text-slate-100">{t.device_brand} {t.device_model}</TableCell>
        <TableCell><TicketStatusBadge status={t.status} /></TableCell>
        <TableCell className="text-right pr-2 font-mono font-bold text-slate-900 dark:text-slate-100">{formatPrice(t.cost)}</TableCell>
      </TableRow>
    ))
  }

  return (
    <div className="space-y-8">
      {/* Title */}
      <div>
        <h1 className="text-3xl font-extrabold text-slate-900 dark:text-slate-100 tracking-tight">
          Dashboard Overview
        </h1>
        <p className="text-slate-500 dark:text-slate-400 text-sm">Welcome back. Here is what is happening at the workshop today.</p>
      </div>

      {error && (
        <div className="bg-red-50 dark:bg-red-950/30 border border-red-200 dark:border-red-900/50 rounded-lg p-4 text-sm text-red-600 dark:text-red-400">
          {error}
        </div>
      )}

      {/* Metrics grid */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
        {renderMetrics()}
      </div>

      {/* Table container using shadcn Table */}
      <Card className="border-slate-200/80 dark:border-slate-800 bg-white/70 dark:bg-slate-900/70 backdrop-blur-xl shadow-sm">
        <CardHeader className="flex flex-row items-center justify-between pb-4 border-b border-slate-100 dark:border-slate-800">
          <div>
            <CardTitle className="font-extrabold text-slate-800 dark:text-slate-200 text-lg">Recent Tickets</CardTitle>
            <CardDescription className="text-slate-500 dark:text-slate-400 text-xs">A list of the 4 most recently updated tickets.</CardDescription>
          </div>
          <Button
            onClick={() => navigate('/tickets')}
            variant="outline"
            size="sm"
            className="font-semibold text-xs gap-1 border-slate-200 dark:border-slate-700 hover:bg-slate-50 dark:hover:bg-slate-800 text-slate-700 dark:text-slate-200 cursor-pointer"
          >
            <span>View All Tickets</span>
            <ArrowUpRight className="w-3.5 h-3.5" />
          </Button>
        </CardHeader>
        <CardContent className="pt-4">
          <Table>
            <TableHeader>
              <TableRow className="border-slate-100 dark:border-slate-800 hover:bg-transparent">
                <TableHead className="w-24 pl-2 font-bold uppercase tracking-wider text-xxs text-slate-500 dark:text-slate-400">ID</TableHead>
                <TableHead className="font-bold uppercase tracking-wider text-xxs text-slate-500 dark:text-slate-400">Customer</TableHead>
                <TableHead className="font-bold uppercase tracking-wider text-xxs text-slate-500 dark:text-slate-400">Device</TableHead>
                <TableHead className="font-bold uppercase tracking-wider text-xxs text-slate-500 dark:text-slate-400">Status</TableHead>
                <TableHead className="text-right pr-2 font-bold uppercase tracking-wider text-xxs text-slate-500 dark:text-slate-400">Cost</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody className="text-sm font-medium text-slate-700 dark:text-slate-300 divide-y divide-slate-100/50 dark:divide-slate-800">
              {renderRecentTickets()}
            </TableBody>
          </Table>
        </CardContent>
      </Card>
    </div>
  )
}

export default DashboardPage
