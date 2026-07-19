import { useState } from 'react'
import type { TicketListItem, TicketStatus } from '@/types/ticket'
import { Card, CardContent } from '@/components/ui/card'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Tabs, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { 
  Dialog, 
  DialogContent, 
  DialogDescription, 
  DialogHeader, 
  DialogTitle, 
  DialogTrigger, 
  DialogFooter
} from '@/components/ui/dialog'
import { Search, Plus, Eye, Edit, ShieldAlert } from 'lucide-react'

const mockTickets: TicketListItem[] = [
  { ticket_id: '1', ticket_number: 'TKT-20260719-0001', status: 'RECEIVED', customer_name: 'Budi Santoso', device_brand: 'Samsung', device_model: 'Galaxy S23', created_at: '2026-07-19T10:00:00Z' },
  { ticket_id: '2', ticket_number: 'TKT-20260719-0002', status: 'REPAIRING', customer_name: 'Jane Smith', device_brand: 'Apple', device_model: 'MacBook Air M1', created_at: '2026-07-19T11:30:00Z' },
  { ticket_id: '3', ticket_number: 'TKT-20260718-0003', status: 'PENDING_CONFIRMATION', customer_name: 'Alex Johnson', device_brand: 'Apple', device_model: 'iPad Pro', created_at: '2026-07-18T14:15:00Z' },
  { ticket_id: '4', ticket_number: 'TKT-20260718-0004', status: 'FIXED', customer_name: 'Sarah Connor', device_brand: 'Sony', device_model: 'WH-1000XM4', created_at: '2026-07-18T16:00:00Z' },
  { ticket_id: '5', ticket_number: 'TKT-20260717-0005', status: 'COMPLETED', customer_name: 'David Miller', device_brand: 'Xiaomi', device_model: 'Redmi Note 12', created_at: '2026-07-17T09:00:00Z' },
  { ticket_id: '6', ticket_number: 'TKT-20260717-0006', status: 'CANCELLED', customer_name: 'Emily Davis', device_brand: 'Google', device_model: 'Pixel 7', created_at: '2026-07-17T13:45:00Z' },
  { ticket_id: '7', ticket_number: 'TKT-20260716-0007', status: 'RETURNED', customer_name: 'Michael Brown', device_brand: 'Asus', device_model: 'ROG Phone 6', created_at: '2026-07-16T15:20:00Z' }
]

function TicketsPage() {
  const [searchQuery, setSearchQuery] = useState('')
  const [activeFilterTab, setActiveFilterTab] = useState('all')
  const [tickets, setTickets] = useState<TicketListItem[]>(mockTickets)
  const [isCreateOpen, setIsCreateOpen] = useState(false)

  // Form State matching POST /api/v1/admin/services
  const [newTicket, setNewTicket] = useState({
    customer_name: '',
    customer_phone: '',
    device_brand: '',
    device_model: '',
    device_passcode: '',
    issue_description: '',
    repair_action: '',
    cost: 0,
    warranty_days: 30
  })

  const getStatusBadge = (status: TicketStatus) => {
    switch (status) {
      case 'RECEIVED':
        return <Badge variant="outline" className="bg-slate-100/50 text-slate-600 border-slate-200 font-semibold">Received</Badge>
      case 'REPAIRING':
        return <Badge className="bg-blue-500/10 text-blue-600 border-none font-semibold hover:bg-blue-500/15">Repairing</Badge>
      case 'PENDING_CONFIRMATION':
        return <Badge className="bg-orange-500/10 text-orange-600 border-none font-semibold hover:bg-orange-500/15">Pending Confirm</Badge>
      case 'FIXED':
        return <Badge className="bg-purple-500/10 text-purple-600 border-none font-semibold hover:bg-purple-500/15">Fixed</Badge>
      case 'COMPLETED':
        return <Badge className="bg-green-500/10 text-green-600 border-none font-semibold hover:bg-green-500/15">Completed</Badge>
      case 'CANCELLED':
        return <Badge className="bg-red-500/10 text-red-600 border-none font-semibold hover:bg-red-500/15">Cancelled</Badge>
      case 'RETURNED':
        return <Badge className="bg-gray-500/10 text-gray-600 border-none font-semibold hover:bg-gray-500/15">Returned</Badge>
      default:
        return <Badge variant="outline">{status}</Badge>
    }
  }

  // Filter logic
  const filteredTickets = tickets.filter(t => {
    // 1. Search Query filter (matches Ticket Number, Customer Name, Brand, Model)
    const matchesSearch = 
      t.ticket_number.toLowerCase().includes(searchQuery.toLowerCase()) ||
      t.customer_name.toLowerCase().includes(searchQuery.toLowerCase()) ||
      t.device_brand.toLowerCase().includes(searchQuery.toLowerCase()) ||
      t.device_model.toLowerCase().includes(searchQuery.toLowerCase())

    // 2. Status Tab filter
    let matchesStatus = true
    if (activeFilterTab === 'active') {
      matchesStatus = ['RECEIVED', 'REPAIRING', 'PENDING_CONFIRMATION', 'FIXED'].includes(t.status)
    } else if (activeFilterTab === 'closed') {
      matchesStatus = ['COMPLETED', 'CANCELLED', 'RETURNED'].includes(t.status)
    }

    return matchesSearch && matchesStatus
  })

  const formatDate = (isoString: string) => {
    return new Date(isoString).toLocaleDateString('id-ID', {
      day: 'numeric',
      month: 'short',
      year: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    })
  }

  const handleCreateSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    // Generate a random ticket ID and number
    const randomId = Math.random().toString(36).substr(2, 9)
    const todayStr = new Date().toISOString().slice(0,10).replace(/-/g,"")
    const newTicketItem: TicketListItem = {
      ticket_id: randomId,
      ticket_number: `TKT-${todayStr}-${Math.floor(1000 + Math.random() * 9000)}`,
      status: 'RECEIVED',
      customer_name: newTicket.customer_name,
      device_brand: newTicket.device_brand,
      device_model: newTicket.device_model,
      created_at: new Date().toISOString()
    }

    setTickets([newTicketItem, ...tickets])
    setIsCreateOpen(false)
    // Reset Form
    setNewTicket({
      customer_name: '',
      customer_phone: '',
      device_brand: '',
      device_model: '',
      device_passcode: '',
      issue_description: '',
      repair_action: '',
      cost: 0,
      warranty_days: 30
    })
  }

  return (
    <div className="space-y-8">
      {/* Title */}
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <h1 className="text-3xl font-extrabold text-slate-900 tracking-tight">
            Repair Tickets
          </h1>
          <p className="text-slate-500 text-sm">Manage, track, and update service tickets for incoming repair devices.</p>
        </div>

        {/* Dialog Form for New Ticket */}
        <Dialog open={isCreateOpen} onOpenChange={setIsCreateOpen}>
          <DialogTrigger render={
            <Button className="font-semibold bg-primary hover:bg-secondary cursor-pointer">
              <Plus className="w-4 h-4 mr-1" />
              New Ticket
            </Button>
          } />
          <DialogContent className="max-w-xl max-h-[90vh] overflow-y-auto">
            <form onSubmit={handleCreateSubmit}>
              <DialogHeader>
                <DialogTitle className="text-xl font-extrabold text-slate-900">Create Service Ticket</DialogTitle>
                <DialogDescription>
                  Enter customer information, device details, and diagnostic notes to create a new ticket.
                </DialogDescription>
              </DialogHeader>

              <div className="grid grid-cols-1 sm:grid-cols-2 gap-4 py-4">
                {/* Customer Details */}
                <div className="sm:col-span-2 border-b border-slate-100 pb-2">
                  <h4 className="text-sm font-bold text-slate-800">Customer Information</h4>
                </div>
                
                <div className="space-y-1">
                  <label className="text-xs font-bold text-slate-500 uppercase">Customer Name</label>
                  <Input 
                    required 
                    placeholder="e.g. John Doe"
                    value={newTicket.customer_name} 
                    onChange={e => setNewTicket({...newTicket, customer_name: e.target.value})} 
                  />
                </div>
                
                <div className="space-y-1">
                  <label className="text-xs font-bold text-slate-500 uppercase">Customer Phone</label>
                  <Input 
                    required 
                    placeholder="e.g. 08123456789"
                    value={newTicket.customer_phone} 
                    onChange={e => setNewTicket({...newTicket, customer_phone: e.target.value})} 
                  />
                </div>

                {/* Device Details */}
                <div className="sm:col-span-2 border-b border-slate-100 pt-2 pb-2">
                  <h4 className="text-sm font-bold text-slate-800">Device details</h4>
                </div>

                <div className="space-y-1">
                  <label className="text-xs font-bold text-slate-500 uppercase">Brand</label>
                  <Input 
                    required 
                    placeholder="e.g. Apple"
                    value={newTicket.device_brand} 
                    onChange={e => setNewTicket({...newTicket, device_brand: e.target.value})} 
                  />
                </div>

                <div className="space-y-1">
                  <label className="text-xs font-bold text-slate-500 uppercase">Model</label>
                  <Input 
                    required 
                    placeholder="e.g. iPhone 13 Pro"
                    value={newTicket.device_model} 
                    onChange={e => setNewTicket({...newTicket, device_model: e.target.value})} 
                  />
                </div>

                <div className="sm:col-span-2 space-y-1">
                  <label className="text-xs font-bold text-slate-500 uppercase">Screen Passcode (Optional)</label>
                  <Input 
                    placeholder="e.g. pattern Letter-L, PIN 1234"
                    value={newTicket.device_passcode} 
                    onChange={e => setNewTicket({...newTicket, device_passcode: e.target.value})} 
                  />
                </div>

                {/* Diagnosis Details */}
                <div className="sm:col-span-2 border-b border-slate-100 pt-2 pb-2">
                  <h4 className="text-sm font-bold text-slate-800">Initial Diagnostic & Pricing</h4>
                </div>

                <div className="sm:col-span-2 space-y-1">
                  <label className="text-xs font-bold text-slate-500 uppercase">Issue Description</label>
                  <Input 
                    required 
                    placeholder="e.g. Cracked LCD, Touch unresponsive"
                    value={newTicket.issue_description} 
                    onChange={e => setNewTicket({...newTicket, issue_description: e.target.value})} 
                  />
                </div>

                <div className="sm:col-span-2 space-y-1">
                  <label className="text-xs font-bold text-slate-500 uppercase">Proposed Repair Action</label>
                  <Input 
                    required 
                    placeholder="e.g. Replacement LCD OLED Screen"
                    value={newTicket.repair_action} 
                    onChange={e => setNewTicket({...newTicket, repair_action: e.target.value})} 
                  />
                </div>

                <div className="space-y-1">
                  <label className="text-xs font-bold text-slate-500 uppercase">Estimated Cost (Rp)</label>
                  <Input 
                    required 
                    type="number"
                    value={newTicket.cost} 
                    onChange={e => setNewTicket({...newTicket, cost: parseInt(e.target.value) || 0})} 
                  />
                </div>

                <div className="space-y-1">
                  <label className="text-xs font-bold text-slate-500 uppercase">Warranty Period (Days)</label>
                  <Input 
                    required 
                    type="number"
                    value={newTicket.warranty_days} 
                    onChange={e => setNewTicket({...newTicket, warranty_days: parseInt(e.target.value) || 0})} 
                  />
                </div>
              </div>

              <DialogFooter className="gap-2 sm:gap-0 pt-2">
                <Button type="button" variant="outline" className="cursor-pointer" onClick={() => setIsCreateOpen(false)}>Cancel</Button>
                <Button type="submit" className="bg-primary hover:bg-secondary cursor-pointer">Create Ticket</Button>
              </DialogFooter>
            </form>
          </DialogContent>
        </Dialog>
      </div>

      {/* Filters and Search toolbar */}
      <div className="flex flex-col md:flex-row gap-4 items-center justify-between">
        {/* Search */}
        <div className="relative w-full md:w-80">
          <Search className="absolute left-3 top-2.5 h-4 w-4 text-slate-400" />
          <Input
            placeholder="Search number, name, device..."
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            className="pl-9 bg-white border-slate-200 focus-visible:ring-primary/20"
          />
        </div>

        {/* Quick Filters (Tabs) */}
        <Tabs value={activeFilterTab} onValueChange={setActiveFilterTab} className="w-full md:w-auto">
          <TabsList className="bg-slate-100 border border-slate-200/50 p-1">
            <TabsTrigger value="all" className="data-[state=active]:bg-white data-[state=active]:shadow-sm font-semibold text-xs px-4">All Tickets</TabsTrigger>
            <TabsTrigger value="active" className="data-[state=active]:bg-white data-[state=active]:shadow-sm font-semibold text-xs px-4">Active Repairs</TabsTrigger>
            <TabsTrigger value="closed" className="data-[state=active]:bg-white data-[state=active]:shadow-sm font-semibold text-xs px-4">Closed</TabsTrigger>
          </TabsList>
        </Tabs>
      </div>

      {/* Main Table */}
      <Card className="border-slate-200/80 bg-white shadow-sm overflow-hidden">
        <CardContent className="p-0">
          <Table>
            <TableHeader className="bg-slate-50 border-b border-slate-100">
              <TableRow className="hover:bg-transparent">
                <TableHead className="w-32 pl-6 font-bold uppercase tracking-wider text-xxs">Ticket Number</TableHead>
                <TableHead className="font-bold uppercase tracking-wider text-xxs">Date Registered</TableHead>
                <TableHead className="font-bold uppercase tracking-wider text-xxs">Customer Name</TableHead>
                <TableHead className="font-bold uppercase tracking-wider text-xxs">Device details</TableHead>
                <TableHead className="font-bold uppercase tracking-wider text-xxs">Status</TableHead>
                <TableHead className="text-center pr-6 font-bold uppercase tracking-wider text-xxs">Actions</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody className="text-sm font-medium text-slate-700 divide-y divide-slate-100/50">
              {filteredTickets.length > 0 ? (
                filteredTickets.map((t) => (
                  <TableRow key={t.ticket_id} className="border-slate-100/50 hover:bg-slate-50/30 transition-colors">
                    <TableCell className="pl-6 font-mono text-xs font-bold text-slate-600">{t.ticket_number}</TableCell>
                    <TableCell className="text-slate-500 font-semibold">{formatDate(t.created_at)}</TableCell>
                    <TableCell className="font-semibold text-slate-900">{t.customer_name}</TableCell>
                    <TableCell>
                      <span className="font-bold text-slate-900">{t.device_brand}</span>{' '}
                      <span className="text-slate-500">{t.device_model}</span>
                    </TableCell>
                    <TableCell>{getStatusBadge(t.status)}</TableCell>
                    <TableCell className="text-center pr-6">
                      <div className="flex items-center justify-center gap-1.5">
                        <Button variant="ghost" size="icon-sm" className="h-7 w-7 text-slate-500 hover:text-primary hover:bg-slate-100 rounded-md cursor-pointer" title="View details">
                          <Eye className="w-4 h-4" />
                        </Button>
                        <Button variant="ghost" size="icon-sm" className="h-7 w-7 text-slate-500 hover:text-tertiary hover:bg-slate-100 rounded-md cursor-pointer" title="Update status">
                          <Edit className="w-4 h-4" />
                        </Button>
                      </div>
                    </TableCell>
                  </TableRow>
                ))
              ) : (
                <TableRow>
                  <TableCell colSpan={6} className="h-32 text-center text-slate-400">
                    <div className="flex flex-col items-center justify-center gap-1">
                      <ShieldAlert className="w-6 h-6 text-slate-300" />
                      <span>No tickets found matching your query</span>
                    </div>
                  </TableCell>
                </TableRow>
              )}
            </TableBody>
          </Table>
        </CardContent>
      </Card>
    </div>
  )
}

export default TicketsPage
