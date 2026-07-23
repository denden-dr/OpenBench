import { useState, useEffect, useCallback } from 'react'
import type { TicketListItem } from '@/types/ticket'
import { ticketService } from '@/services/ticketService'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Tabs, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Search, Plus, ShieldAlert, RefreshCw } from 'lucide-react'
import { CreateTicketModal } from '@/features/tickets/components/CreateTicketModal'
import { TicketDetailsModal } from '@/features/tickets/components/TicketDetailsModal'
import { UpdateStatusModal } from '@/features/tickets/components/UpdateStatusModal'
import { TicketsTable } from '@/features/tickets/components/TicketsTable'

function TicketsPage() {
  const [searchQuery, setSearchQuery] = useState('')
  const [activeFilterTab, setActiveFilterTab] = useState('all')
  const [tickets, setTickets] = useState<TicketListItem[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  // Dialog / Modal Visibility States
  const [isCreateOpen, setIsCreateOpen] = useState(false)
  const [isViewOpen, setIsViewOpen] = useState(false)
  const [isStatusOpen, setIsStatusOpen] = useState(false)

  // Selected ticket tracking
  const [selectedTicketId, setSelectedTicketId] = useState<string | null>(null)
  const [selectedTicketForStatus, setSelectedTicketForStatus] = useState<TicketListItem | null>(null)

  const fetchTickets = useCallback(async () => {
    setLoading(true)
    setError(null)
    try {
      const response = await ticketService.getTickets({
        search: searchQuery.trim() || undefined
      })
      setTickets(response.data || [])
    } catch (err: unknown) {
      console.error('Failed to fetch tickets:', err)
      setError('Gagal memuat data tiket servis. Silakan coba lagi.')
    } finally {
      setLoading(false)
    }
  }, [searchQuery])

  useEffect(() => {
    const timer = setTimeout(() => {
      fetchTickets()
    }, 300)
    return () => clearTimeout(timer)
  }, [fetchTickets])

  const filteredTickets = tickets.filter(t => {
    if (activeFilterTab === 'active') {
      return ['RECEIVED', 'REPAIRING', 'PENDING_CONFIRMATION', 'FIXED', 'CANCELLED'].includes(t.status)
    } else if (activeFilterTab === 'closed') {
      return ['COMPLETED', 'RETURNED'].includes(t.status)
    }
    return true
  })

  const handleOpenView = (ticketId: string) => {
    setSelectedTicketId(ticketId)
    setIsViewOpen(true)
  }

  const handleOpenStatus = (ticket: TicketListItem) => {
    setSelectedTicketForStatus(ticket)
    setIsStatusOpen(true)
  }

  return (
    <div className="space-y-8">
      {/* Title */}
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <h1 className="text-3xl font-extrabold text-slate-900 dark:text-slate-100 tracking-tight">
            Repair Tickets
          </h1>
          <p className="text-slate-500 dark:text-slate-400 text-sm">Manage, track, and update service tickets for incoming repair devices.</p>
        </div>

        <Button 
          className="font-semibold bg-primary hover:bg-secondary cursor-pointer"
          onClick={() => setIsCreateOpen(true)}
        >
          <Plus className="w-4 h-4 mr-1" />
          New Ticket
        </Button>
      </div>

      {/* Filters and Search toolbar */}
      <div className="flex flex-col md:flex-row gap-4 items-center justify-between">
        {/* Search */}
        <div className="relative w-full md:w-80 flex items-center gap-2">
          <div className="relative w-full">
            <Search className="absolute left-3 top-2.5 h-4 w-4 text-slate-400" />
            <Input
              placeholder="Search number, name, device..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className="pl-9 bg-white dark:bg-slate-900 border-slate-200 dark:border-slate-800 text-slate-900 dark:text-slate-100 placeholder:text-slate-400 dark:placeholder:text-slate-500 focus-visible:ring-primary/20"
            />
          </div>
          <Button variant="outline" size="icon-sm" onClick={fetchTickets} className="h-9 w-9 border-slate-200 dark:border-slate-800 shrink-0" title="Refresh">
            <RefreshCw className={`w-4 h-4 ${loading ? 'animate-spin' : ''}`} />
          </Button>
        </div>

        {/* Quick Filters (Tabs) */}
        <Tabs value={activeFilterTab} onValueChange={setActiveFilterTab} className="w-full md:w-auto">
          <TabsList className="bg-slate-100 dark:bg-slate-900 border border-slate-200/50 dark:border-slate-800 p-1">
            <TabsTrigger value="all" className="data-[state=active]:bg-white dark:data-[state=active]:bg-slate-800 data-[state=active]:text-slate-900 dark:data-[state=active]:text-slate-100 text-slate-600 dark:text-slate-400 font-semibold text-xs px-4">All Tickets</TabsTrigger>
            <TabsTrigger value="active" className="data-[state=active]:bg-white dark:data-[state=active]:bg-slate-800 data-[state=active]:text-slate-900 dark:data-[state=active]:text-slate-100 text-slate-600 dark:text-slate-400 font-semibold text-xs px-4">Active Repairs</TabsTrigger>
            <TabsTrigger value="closed" className="data-[state=active]:bg-white dark:data-[state=active]:bg-slate-800 data-[state=active]:text-slate-900 dark:data-[state=active]:text-slate-100 text-slate-600 dark:text-slate-400 font-semibold text-xs px-4">Closed</TabsTrigger>
          </TabsList>
        </Tabs>
      </div>

      {/* Error alert */}
      {error && (
        <div className="bg-red-50 dark:bg-red-950/40 border border-red-200 dark:border-red-800 text-red-700 dark:text-red-300 p-4 rounded-lg flex items-center justify-between">
          <div className="flex items-center gap-2 text-sm">
            <ShieldAlert className="w-5 h-5 shrink-0 text-red-500" />
            <span>{error}</span>
          </div>
          <Button variant="outline" size="sm" onClick={fetchTickets} className="border-red-200 text-red-700 hover:bg-red-100 dark:border-red-800 dark:text-red-300 dark:hover:bg-red-900/50">
            Coba Lagi
          </Button>
        </div>
      )}

      {/* Main Table */}
      <TicketsTable 
        tickets={filteredTickets} 
        loading={loading} 
        onViewDetails={handleOpenView} 
        onUpdateStatus={handleOpenStatus} 
      />

      {/* Modals */}
      <CreateTicketModal 
        isOpen={isCreateOpen} 
        onOpenChange={setIsCreateOpen} 
        onSuccess={fetchTickets} 
      />

      <TicketDetailsModal 
        isOpen={isViewOpen} 
        onOpenChange={setIsViewOpen} 
        ticketId={selectedTicketId} 
      />

      <UpdateStatusModal 
        isOpen={isStatusOpen} 
        onOpenChange={setIsStatusOpen} 
        ticket={selectedTicketForStatus} 
        onSuccess={fetchTickets} 
      />
    </div>
  )
}
export default TicketsPage
