import { Card, CardHeader, CardTitle, CardContent, CardDescription } from '@/components/ui/card'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { ArrowUpRight } from 'lucide-react'

function DashboardPage() {
  const metrics = [
    { title: 'Active Tickets', value: '14', change: '+2 new', color: 'text-primary' },
    { title: 'Pending Diagnoses', value: '5', change: 'needs review', color: 'text-tertiary' },
    { title: 'Sales Today', value: '$849.50', change: '+12% vs yesterday', color: 'text-secondary' },
    { title: 'Active Warranties', value: '142', change: '99% compliance', color: 'text-green-600' },
  ]

  const recentTickets = [
    { id: 'TKT-89304', customer: 'John Doe', device: 'iPhone 13 Pro', status: 'In Progress', cost: '$189.00' },
    { id: 'TKT-89303', customer: 'Jane Smith', device: 'MacBook Air M1', status: 'Ready for Pickup', cost: '$249.00' },
    { id: 'TKT-89302', customer: 'Alex Johnson', device: 'iPad Pro', status: 'Waiting for Parts', cost: '$95.00' },
    { id: 'TKT-89301', customer: 'Sarah Connor', device: 'Sony WH-1000XM4', status: 'Diagnosed', cost: '$75.00' },
  ]

  const getStatusBadge = (status: string) => {
    switch (status) {
      case 'Ready for Pickup':
        return <Badge className="bg-green-500/10 text-green-600 border-none font-semibold hover:bg-green-500/15">Ready for Pickup</Badge>
      case 'In Progress':
        return <Badge className="bg-blue-500/10 text-blue-600 border-none font-semibold hover:bg-blue-500/15">In Progress</Badge>
      case 'Waiting for Parts':
        return <Badge className="bg-orange-500/10 text-orange-600 border-none font-semibold hover:bg-orange-500/15">Waiting for Parts</Badge>
      default:
        return <Badge variant="outline" className="text-slate-600 font-semibold">{status}</Badge>
    }
  }

  return (
    <div className="space-y-8">
      {/* Title */}
      <div>
        <h1 className="text-3xl font-extrabold text-slate-900 tracking-tight">
          Dashboard Overview
        </h1>
        <p className="text-slate-500 text-sm">Welcome back. Here is what is happening at the workshop today.</p>
      </div>

      {/* Metrics grid */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
        {metrics.map((metric, idx) => (
          <Card key={idx} className="hover:shadow-md transition-shadow border-slate-200/80 bg-white">
            <CardHeader className="flex flex-row items-center justify-between pb-2">
              <CardTitle className="text-xs font-bold text-slate-500 uppercase tracking-wider">{metric.title}</CardTitle>
            </CardHeader>
            <CardContent className="flex items-baseline justify-between">
              <span className="text-3xl font-extrabold text-slate-900 tracking-tight">{metric.value}</span>
              <span className={`text-xxs font-semibold px-2 py-0.5 rounded-full bg-slate-100 ${metric.color}`}>{metric.change}</span>
            </CardContent>
          </Card>
        ))}
      </div>

      {/* Table container using shadcn Table */}
      <Card className="border-slate-200/80 bg-white/70 backdrop-blur-xl shadow-sm">
        <CardHeader className="flex flex-row items-center justify-between pb-4 border-b border-slate-100">
          <div>
            <CardTitle className="font-extrabold text-slate-800 text-lg">Recent Tickets</CardTitle>
            <CardDescription className="text-slate-500 text-xs">A list of the 4 most recently updated tickets.</CardDescription>
          </div>
          <Button variant="outline" size="sm" className="font-semibold text-xs gap-1 border-slate-200 hover:bg-slate-50">
            <span>View All Tickets</span>
            <ArrowUpRight className="w-3.5 h-3.5" />
          </Button>
        </CardHeader>
        <CardContent className="pt-4">
          <Table>
            <TableHeader>
              <TableRow className="border-slate-100 hover:bg-transparent">
                <TableHead className="w-24 pl-2 font-bold uppercase tracking-wider text-xxs">ID</TableHead>
                <TableHead className="font-bold uppercase tracking-wider text-xxs">Customer</TableHead>
                <TableHead className="font-bold uppercase tracking-wider text-xxs">Device</TableHead>
                <TableHead className="font-bold uppercase tracking-wider text-xxs">Status</TableHead>
                <TableHead className="text-right pr-2 font-bold uppercase tracking-wider text-xxs">Cost</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody className="text-sm font-medium text-slate-700 divide-y divide-slate-100/50">
              {recentTickets.map((t, idx) => (
                <TableRow key={idx} className="border-slate-100/50 hover:bg-slate-50/50 transition-colors">
                  <TableCell className="pl-2 font-mono text-xs font-bold text-slate-500">{t.id}</TableCell>
                  <TableCell>{t.customer}</TableCell>
                  <TableCell className="font-semibold text-slate-900">{t.device}</TableCell>
                  <TableCell>{getStatusBadge(t.status)}</TableCell>
                  <TableCell className="text-right pr-2 font-mono font-bold text-slate-900">{t.cost}</TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </CardContent>
      </Card>
    </div>
  )
}

export default DashboardPage
