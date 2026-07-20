import { useState } from 'react'
import type { TicketStatus } from '@/types/ticket'
import type { Warranty, WarrantyClaim, WarrantyStatus, ClaimEvaluationStatus } from '@/types/warranty'
import { Card, CardHeader, CardTitle, CardContent, CardDescription } from '@/components/ui/card'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Tabs, TabsList, TabsTrigger, TabsContent } from '@/components/ui/tabs'
import { 
  Dialog, 
  DialogContent, 
  DialogDescription, 
  DialogHeader, 
  DialogTitle, 
  DialogFooter
} from '@/components/ui/dialog'
import { 
  Search, 
  ShieldCheck, 
  ShieldAlert, 
  ShieldX, 
  Ban, 
  SlidersHorizontal, 
  Settings, 
  ChevronLeft, 
  ChevronRight,
  ClipboardList,
  AlertTriangle
} from 'lucide-react'

// Initial warranties mock data
const initialWarranties: Record<string, Warranty> = {
  'TKT-20260707-1234': {
    id: 'w-1',
    ticket_id: 'd290f1ee-6c54-4b01-90e6-d701748f0851',
    start_date: '2026-07-07T12:30:00Z',
    end_date: '2026-08-06T12:30:00Z',
    status: 'ACTIVE',
    notes: null
  },
  'TKT-20260601-5678': {
    id: 'w-2',
    ticket_id: 'tkt-expired-id',
    start_date: '2026-06-01T10:00:00Z',
    end_date: '2026-07-01T10:00:00Z',
    status: 'EXPIRED',
    notes: 'Warranty duration has completed'
  },
  'TKT-20260710-9999': {
    id: 'w-3',
    ticket_id: 'tkt-void-id',
    start_date: '2026-07-10T15:00:00Z',
    end_date: '2026-08-09T15:00:00Z',
    status: 'VOID',
    notes: 'Seal broken on inspection'
  }
}

// Initial claims mock data
const initialClaims: WarrantyClaim[] = [
  {
    claim_id: 'c-1',
    claim_number: 'CLM-20260714-0001',
    warranty_id: 'w-1',
    status: 'RECEIVED',
    evaluation_status: 'PENDING',
    issue_description: 'Layar sentuh tidak responsif di bagian pojok kiri atas setelah diganti minggu lalu',
    repair_action: null,
    notes: null,
    evaluation_notes: null,
    created_at: '2026-07-14T09:00:00Z',
    updated_at: '2026-07-14T09:00:00Z'
  },
  {
    claim_id: 'c-2',
    claim_number: 'CLM-20260715-0002',
    warranty_id: 'w-1',
    status: 'REPAIRING',
    evaluation_status: 'ACCEPTED',
    issue_description: 'Speaker pecah suaranya setelah perbaikan modul audio',
    repair_action: 'Ganti modul speaker internal baru',
    notes: 'Penggantian modul audio ditanggung garansi pengerjaan modul audio sebelumnya',
    evaluation_notes: 'Diterima, part digaransi penuh',
    created_at: '2026-07-15T11:00:00Z',
    updated_at: '2026-07-16T10:00:00Z'
  }
]

function WarrantyPage() {
  const [warranties, setWarranties] = useState<Record<string, Warranty>>(initialWarranties)
  const [claims, setClaims] = useState<WarrantyClaim[]>(initialClaims)
  
  // Search & Status state
  const [searchTicket, setSearchTicket] = useState('')
  const [foundWarranty, setFoundWarranty] = useState<Warranty | null>(null)
  const [searchedTicketNum, setSearchedTicketNum] = useState('')
  const [hasSearched, setHasSearched] = useState(false)

  // Dialog open states
  const [isVoidOpen, setIsVoidOpen] = useState(false)
  const [isClaimOpen, setIsClaimOpen] = useState(false)
  const [isEvaluateOpen, setIsEvaluateOpen] = useState(false)
  const [isStatusOpen, setIsStatusOpen] = useState(false)
  const [isEditClaimOpen, setIsEditClaimOpen] = useState(false)
  
  // Selected items
  const [selectedWarranty, setSelectedWarranty] = useState<Warranty | null>(null)
  const [selectedClaim, setSelectedClaim] = useState<WarrantyClaim | null>(null)

  // Form states
  const [voidNotes, setVoidNotes] = useState('')
  const [claimIssue, setClaimIssue] = useState('')
  const [evalStatus, setEvalStatus] = useState<ClaimEvaluationStatus>('ACCEPTED')
  const [evalNotes, setEvalNotes] = useState('')
  
  const [claimStatus, setClaimStatus] = useState<TicketStatus>('RECEIVED')
  
  const [editClaimData, setEditClaimData] = useState({
    issue_description: '',
    repair_action: '',
    notes: ''
  })

  // Tab State
  const [activeTab, setActiveTab] = useState<'check' | 'queue'>('check')

  // Search Logic
  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault()
    if (!searchTicket.trim()) return
    
    const warranty = warranties[searchTicket.trim()]
    setFoundWarranty(warranty || null)
    setSearchedTicketNum(searchTicket.trim())
    setHasSearched(true)
  }

  // Void Warranty Logic
  const handleVoidWarrantySubmit = (e: React.FormEvent) => {
    e.preventDefault()
    if (!selectedWarranty || !searchedTicketNum) return

    setWarranties({
      ...warranties,
      [searchedTicketNum]: {
        ...selectedWarranty,
        status: 'VOID',
        notes: voidNotes
      }
    })
    
    // Update local state details if they match
    if (foundWarranty && foundWarranty.id === selectedWarranty.id) {
      setFoundWarranty({
        ...foundWarranty,
        status: 'VOID',
        notes: voidNotes
      })
    }

    setIsVoidOpen(false)
    setSelectedWarranty(null)
    setVoidNotes('')
  }

  // Create Claim Logic
  const handleCreateClaimSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    if (!foundWarranty || !searchedTicketNum) return

    const newClaim: WarrantyClaim = {
      claim_id: `c-${Math.random().toString(36).substr(2, 9)}`,
      claim_number: `CLM-${new Date().toISOString().slice(0, 10).replace(/-/g, '')}-${Math.floor(1000 + Math.random() * 9000)}`,
      warranty_id: foundWarranty.id,
      status: 'RECEIVED',
      evaluation_status: 'PENDING',
      issue_description: claimIssue,
      repair_action: null,
      notes: null,
      evaluation_notes: null,
      created_at: new Date().toISOString(),
      updated_at: new Date().toISOString()
    }

    // Log the payload as per API contract
    console.log('Sending Create Claim Payload:', {
      ticket_number: searchedTicketNum,
      issue_description: claimIssue
    })

    setClaims([newClaim, ...claims])
    setIsClaimOpen(false)
    setClaimIssue('')
    // Switch to queue to show the added claim
    setActiveTab('queue')
  }

  // Evaluate Claim Logic
  const handleEvaluateSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    if (!selectedClaim) return

    // Validate that notes are filled if status is REJECTED or VOID
    if ((evalStatus === 'REJECTED' || evalStatus === 'VOID') && !evalNotes.trim()) {
      alert('Notes are required when rejecting or voiding a claim')
      return
    }

    // Update claim
    const updatedClaims = claims.map(c => {
      if (c.claim_id === selectedClaim.claim_id) {
        return {
          ...c,
          evaluation_status: evalStatus,
          evaluation_notes: evalNotes,
          status: evalStatus === 'ACCEPTED' ? ('REPAIRING' as TicketStatus) : c.status,
          updated_at: new Date().toISOString()
        }
      }
      return c
    })
    setClaims(updatedClaims)

    // If evaluated to VOID, void the parent warranty as well
    if (evalStatus === 'VOID') {
      const updatedWarranties = { ...warranties }
      Object.keys(updatedWarranties).forEach(key => {
        if (updatedWarranties[key].id === selectedClaim.warranty_id) {
          updatedWarranties[key] = {
            ...updatedWarranties[key],
            status: 'VOID',
            notes: `Claim ${selectedClaim.claim_number} evaluated as VOID: ${evalNotes}`
          }
        }
      })
      setWarranties(updatedWarranties)

      // Also sync if it is currently viewed in found warranty
      if (foundWarranty && foundWarranty.id === selectedClaim.warranty_id) {
        setFoundWarranty({
          ...foundWarranty,
          status: 'VOID',
          notes: `Claim ${selectedClaim.claim_number} evaluated as VOID: ${evalNotes}`
        })
      }
    }

    setIsEvaluateOpen(false)
    setSelectedClaim(null)
    setEvalNotes('')
  }

  // Update Progres Perbaikan Status Logic
  const handleStatusSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    if (!selectedClaim) return

    setClaims(claims.map(c => {
      if (c.claim_id === selectedClaim.claim_id) {
        return {
          ...c,
          status: claimStatus,
          updated_at: new Date().toISOString()
        }
      }
      return c
    }))

    setIsStatusOpen(false)
    setSelectedClaim(null)
  }

  // Update claim info logic
  const handleEditClaimSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    if (!selectedClaim) return

    setClaims(claims.map(c => {
      if (c.claim_id === selectedClaim.claim_id) {
        return {
          ...c,
          issue_description: editClaimData.issue_description,
          repair_action: editClaimData.repair_action,
          notes: editClaimData.notes,
          updated_at: new Date().toISOString()
        }
      }
      return c
    }))

    setIsEditClaimOpen(false)
    setSelectedClaim(null)
  }

  const formatDate = (dateStr: string) => {
    return new Date(dateStr).toLocaleString('id-ID', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    })
  }

  const getWarrantyBadge = (status: WarrantyStatus) => {
    switch (status) {
      case 'ACTIVE':
        return <Badge className="bg-green-500/10 text-green-600 border-none font-semibold px-2.5 py-1">ACTIVE</Badge>
      case 'VOID':
        return <Badge className="bg-red-500/10 text-red-600 border-none font-semibold px-2.5 py-1">VOIDED</Badge>
      case 'EXPIRED':
        return <Badge className="bg-slate-500/10 text-slate-600 border-none font-semibold px-2.5 py-1">EXPIRED</Badge>
    }
  }

  const getEvalBadge = (status: ClaimEvaluationStatus) => {
    switch (status) {
      case 'PENDING':
        return <Badge className="bg-amber-500/10 text-amber-600 border-none font-semibold">PENDING</Badge>
      case 'ACCEPTED':
        return <Badge className="bg-green-500/10 text-green-600 border-none font-semibold">ACCEPTED</Badge>
      case 'REJECTED':
        return <Badge className="bg-red-500/10 text-red-600 border-none font-semibold">REJECTED</Badge>
      case 'VOID':
        return <Badge className="bg-rose-950/20 text-rose-800 border-none font-semibold">VOID</Badge>
    }
  }

  const getStatusBadge = (status: TicketStatus) => {
    switch (status) {
      case 'RECEIVED':
        return <Badge className="bg-blue-500/10 text-blue-600 border-none font-semibold">DITERIMA</Badge>
      case 'REPAIRING':
        return <Badge className="bg-amber-500/10 text-amber-600 border-none font-semibold">DIPERBAIKI</Badge>
      case 'PENDING_CONFIRMATION':
        return <Badge className="bg-purple-500/10 text-purple-600 border-none font-semibold">BUTUH KONFIRMASI</Badge>
      case 'FIXED':
        return <Badge className="bg-indigo-500/10 text-indigo-600 border-none font-semibold">SELESAI DIKERJAKAN</Badge>
      case 'COMPLETED':
        return <Badge className="bg-green-500/10 text-green-600 border-none font-semibold">DIAMBIL (TUNTAS)</Badge>
      case 'CANCELLED':
        return <Badge className="bg-red-500/10 text-red-600 border-none font-semibold">DIBATALKAN</Badge>
      case 'RETURNED':
        return <Badge className="bg-slate-500/10 text-slate-600 border-none font-semibold">DIKEMBALIKAN (BATAL)</Badge>
    }
  }

  return (
    <div className="space-y-8">
      {/* Title */}
      <div>
        <h1 className="text-3xl font-extrabold text-slate-900 tracking-tight">
          Warranty & Claim Ticketing
        </h1>
        <p className="text-slate-500 text-sm">Validate customer warranty cards, void contracts, and process claims separately from regular ticketing.</p>
      </div>

      <Tabs value={activeTab} onValueChange={(v) => setActiveTab(v as 'check' | 'queue')} className="space-y-6">
        <TabsList className="bg-slate-100 p-1 rounded-xl">
          <TabsTrigger value="check" className="rounded-lg px-4 py-2 text-xs font-semibold cursor-pointer">
            Check & Claim Warranty
          </TabsTrigger>
          <TabsTrigger value="queue" className="rounded-lg px-4 py-2 text-xs font-semibold cursor-pointer">
            Claims Queue
          </TabsTrigger>
        </TabsList>

        {/* Tab 1: Check Warranty */}
        <TabsContent value="check" className="space-y-6">
          <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
            {/* Checker Form */}
            <Card className="lg:col-span-1 border-slate-200/80 bg-white shadow-sm h-fit">
              <CardHeader>
                <CardTitle className="text-lg font-bold text-slate-900">Warranty Checker</CardTitle>
                <CardDescription>Enter the ticket number of the previous repair to verify warranty status.</CardDescription>
              </CardHeader>
              <CardContent>
                <form onSubmit={handleSearch} className="space-y-4">
                  <div className="space-y-1">
                    <label className="text-xs font-bold text-slate-500 uppercase">Original Ticket ID</label>
                    <div className="relative">
                      <Search className="absolute left-3 top-2.5 h-4 w-4 text-slate-400" />
                      <Input
                        required
                        placeholder="e.g. TKT-20260707-1234"
                        value={searchTicket}
                        onChange={(e) => setSearchTicket(e.target.value)}
                        className="pl-9 bg-white border-slate-200 focus-visible:ring-primary/20"
                      />
                    </div>
                    <span className="text-[10px] text-slate-400">
                      Try: <span className="font-mono text-slate-600 font-semibold cursor-pointer underline" onClick={() => setSearchTicket('TKT-20260707-1234')}>TKT-20260707-1234</span> (Active), or <span className="font-mono text-slate-600 font-semibold cursor-pointer underline" onClick={() => setSearchTicket('TKT-20260601-5678')}>TKT-20260601-5678</span> (Expired)
                    </span>
                  </div>

                  <Button type="submit" className="w-full font-semibold bg-primary hover:bg-secondary cursor-pointer">
                    Verify Warranty
                  </Button>
                </form>
              </CardContent>
            </Card>

            {/* Checker Results */}
            <div className="lg:col-span-2 space-y-6">
              {hasSearched ? (
                foundWarranty ? (
                  <Card className="border-slate-200/80 bg-white/70 backdrop-blur-xl shadow-lg relative overflow-hidden">
                    <div className="absolute right-0 top-0 w-32 h-32 bg-primary/5 rounded-full blur-3xl -z-10" />
                    
                    <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-4 border-b border-slate-100">
                      <div>
                        <CardTitle className="text-xl font-extrabold text-slate-900">Warranty Contract Details</CardTitle>
                        <CardDescription className="font-mono text-xs">{foundWarranty.id}</CardDescription>
                      </div>
                      {getWarrantyBadge(foundWarranty.status)}
                    </CardHeader>

                    <CardContent className="pt-6 space-y-6 text-sm">
                      <div className="grid grid-cols-2 gap-6">
                        <div className="space-y-1">
                          <span className="text-xxs font-bold text-slate-400 uppercase tracking-wider block">Start Date</span>
                          <span className="font-semibold text-slate-700">{formatDate(foundWarranty.start_date)}</span>
                        </div>
                        <div className="space-y-1">
                          <span className="text-xxs font-bold text-slate-400 uppercase tracking-wider block">End Date (Expiration)</span>
                          <span className="font-semibold text-slate-700">{formatDate(foundWarranty.end_date)}</span>
                        </div>
                        <div className="space-y-1 col-span-2">
                          <span className="text-xxs font-bold text-slate-400 uppercase tracking-wider block">Original Ticket Reference</span>
                          <span className="font-mono font-bold text-slate-900">{searchedTicketNum}</span>
                        </div>
                        {foundWarranty.notes && (
                          <div className="space-y-1 col-span-2 p-3 bg-slate-50 border border-slate-150 rounded-lg">
                            <span className="text-xxs font-bold text-slate-400 uppercase tracking-wider block">Status Notes</span>
                            <span className="text-slate-600 italic">"{foundWarranty.notes}"</span>
                          </div>
                        )}
                      </div>

                      {/* Action buttons based on status */}
                      <div className="flex items-center gap-3 pt-4 border-t border-slate-100">
                        {foundWarranty.status === 'ACTIVE' ? (
                          <>
                            <Button 
                              onClick={() => setIsClaimOpen(true)}
                              className="font-semibold bg-primary hover:bg-secondary cursor-pointer"
                            >
                              <ClipboardList className="w-4 h-4 mr-1.5" />
                              Submit Warranty Claim
                            </Button>
                            <Button 
                              variant="outline" 
                              onClick={() => {
                                setSelectedWarranty(foundWarranty)
                                setIsVoidOpen(true)
                              }}
                              className="font-semibold text-red-600 hover:text-red-700 border-red-200 hover:bg-red-50 hover:border-red-300 cursor-pointer"
                            >
                              <Ban className="w-4 h-4 mr-1.5" />
                              Void Warranty
                            </Button>
                          </>
                        ) : (
                          <div className="flex items-center gap-2 text-slate-400 bg-slate-50 border border-slate-100 rounded-lg p-3 w-full">
                            <AlertTriangle className="w-5 h-5 text-amber-500 flex-shrink-0" />
                            <span className="text-xs font-semibold">
                              This warranty is {foundWarranty.status.toLowerCase()}. You cannot submit new claims or modify its status further.
                            </span>
                          </div>
                        )}
                      </div>
                    </CardContent>
                  </Card>
                ) : (
                  <Card className="border-slate-200/80 bg-white shadow-sm flex flex-col items-center justify-center p-12 text-center">
                    <ShieldAlert className="w-12 h-12 text-red-400 mb-4" />
                    <h3 className="text-lg font-bold text-slate-900 mb-1">No Warranty Found</h3>
                    <p className="text-slate-400 text-sm max-w-sm mb-4">
                      We couldn't find a warranty contract associated with the ticket number <span className="font-mono font-bold text-slate-900">"{searchedTicketNum}"</span>.
                    </p>
                  </Card>
                )
              ) : (
                <Card className="border-slate-200/80 bg-white/70 backdrop-blur-xl shadow-sm flex flex-col items-center justify-center p-16 text-center border-dashed">
                  <ShieldCheck className="w-14 h-14 text-slate-300 mb-4" />
                  <h3 className="text-lg font-bold text-slate-900 mb-1">Awaiting Warranty Check</h3>
                  <p className="text-slate-400 text-sm max-w-sm">
                    Enter a ticket number in the panel to check validity, void contracts, or launch a claim workflow.
                  </p>
                </Card>
              )}
            </div>
          </div>
        </TabsContent>

        {/* Tab 2: Claims Queue */}
        <TabsContent value="queue" className="space-y-6">
          <Card className="border-slate-200/80 bg-white shadow-sm overflow-hidden">
            <CardContent className="p-0">
              <Table>
                <TableHeader className="bg-slate-50 border-b border-slate-100">
                  <TableRow className="hover:bg-transparent">
                    <TableHead className="w-32 pl-6 font-bold uppercase tracking-wider text-xxs">Claim Num</TableHead>
                    <TableHead className="w-24 font-bold uppercase tracking-wider text-xxs">Warranty ID</TableHead>
                    <TableHead className="font-bold uppercase tracking-wider text-xxs">Issue Description</TableHead>
                    <TableHead className="w-32 font-bold uppercase tracking-wider text-xxs">Evaluation Status</TableHead>
                    <TableHead className="w-32 font-bold uppercase tracking-wider text-xxs">Repair Status</TableHead>
                    <TableHead className="text-center pr-6 w-36 font-bold uppercase tracking-wider text-xxs">Actions</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody className="text-sm font-medium text-slate-700 divide-y divide-slate-100/50">
                  {claims.length > 0 ? (
                    claims.map((claim) => (
                      <TableRow key={claim.claim_id} className="border-slate-100/50 hover:bg-slate-50/30 transition-colors">
                        <TableCell className="pl-6 font-mono text-xs font-bold text-slate-900">{claim.claim_number}</TableCell>
                        <TableCell className="font-mono text-xs text-slate-500 font-semibold">{claim.warranty_id}</TableCell>
                        <TableCell className="max-w-xs truncate text-slate-600 font-semibold" title={claim.issue_description}>
                          {claim.issue_description}
                        </TableCell>
                        <TableCell>{getEvalBadge(claim.evaluation_status)}</TableCell>
                        <TableCell>{getStatusBadge(claim.status)}</TableCell>
                        <TableCell className="text-center pr-6">
                          <div className="flex items-center justify-center gap-1">
                            {/* Evaluation Action */}
                            {claim.evaluation_status === 'PENDING' ? (
                              <Button
                                size="xs"
                                className="bg-primary hover:bg-secondary text-xxs font-bold px-2 py-1 h-7 cursor-pointer"
                                onClick={() => {
                                  setSelectedClaim(claim)
                                  setEvalStatus('ACCEPTED')
                                  setIsEvaluateOpen(true)
                                }}
                              >
                                Evaluate
                              </Button>
                            ) : (
                              <>
                                <Button
                                  variant="ghost"
                                  size="icon-xs"
                                  className="h-7 w-7 text-slate-500 hover:text-primary hover:bg-slate-100 rounded-md cursor-pointer"
                                  title="Update Progress Status"
                                  onClick={() => {
                                    setSelectedClaim(claim)
                                    setClaimStatus(claim.status)
                                    setIsStatusOpen(true)
                                  }}
                                >
                                  <SlidersHorizontal className="w-3.5 h-3.5" />
                                </Button>
                                <Button
                                  variant="ghost"
                                  size="icon-xs"
                                  className="h-7 w-7 text-slate-500 hover:text-primary hover:bg-slate-100 rounded-md cursor-pointer"
                                  title="Update Technician Info"
                                  onClick={() => {
                                    setSelectedClaim(claim)
                                    setEditClaimData({
                                      issue_description: claim.issue_description,
                                      repair_action: claim.repair_action || '',
                                      notes: claim.notes || ''
                                    })
                                    setIsEditClaimOpen(true)
                                  }}
                                >
                                  <Settings className="w-3.5 h-3.5" />
                                </Button>
                              </>
                            )}
                          </div>
                        </TableCell>
                      </TableRow>
                    ))
                  ) : (
                    <TableRow>
                      <TableCell colSpan={6} className="h-32 text-center text-slate-400">
                        No claims registered.
                      </TableCell>
                    </TableRow>
                  )}
                </TableBody>
              </Table>

              {/* Pagination Footer */}
              <div className="border-t border-slate-100 px-6 py-4 flex items-center justify-between bg-slate-50/50">
                <span className="text-xs font-semibold text-slate-500">
                  Showing {claims.length} claim tickets
                </span>
                <div className="flex items-center gap-2">
                  <Button variant="outline" size="sm" className="h-8 text-xs font-semibold border-slate-200 text-slate-500 hover:text-slate-900 cursor-not-allowed" disabled>
                    <ChevronLeft className="w-3.5 h-3.5 mr-1" />
                    Previous
                  </Button>
                  <Button variant="outline" size="sm" className="h-8 text-xs font-semibold border-slate-200 text-slate-500 hover:text-slate-900 cursor-not-allowed" disabled>
                    Next
                    <ChevronRight className="w-3.5 h-3.5 ml-1" />
                  </Button>
                </div>
              </div>
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>

      {/* Dialog: Void Warranty */}
      <Dialog open={isVoidOpen} onOpenChange={(open) => {
        setIsVoidOpen(open)
        if(!open) {
          setSelectedWarranty(null)
          setVoidNotes('')
        }
      }}>
        <DialogContent className="max-w-sm">
          <form onSubmit={handleVoidWarrantySubmit}>
            <DialogHeader>
              <DialogTitle className="text-xl font-extrabold text-slate-900 flex items-center gap-2">
                <ShieldX className="w-5 h-5 text-red-500" />
                Void Warranty Contract
              </DialogTitle>
              <DialogDescription>
                Are you sure you want to invalidate warranty <span className="font-mono font-bold text-slate-900">{selectedWarranty?.id}</span>? This action is irreversible.
              </DialogDescription>
            </DialogHeader>

            <div className="space-y-4 py-4">
              <div className="space-y-1">
                <label className="text-xs font-bold text-slate-500 uppercase">Reason / Void Notes</label>
                <Input 
                  required
                  placeholder="e.g. Broken warranty seal, water damage detected."
                  value={voidNotes}
                  onChange={(e) => setVoidNotes(e.target.value)}
                />
              </div>
            </div>

            <DialogFooter className="gap-2 sm:gap-0 pt-2">
              <Button type="button" variant="outline" className="cursor-pointer" onClick={() => setIsVoidOpen(false)}>Cancel</Button>
              <Button type="submit" className="bg-red-600 hover:bg-red-700 text-white cursor-pointer">Invalidate (Void)</Button>
            </DialogFooter>
          </form>
        </DialogContent>
      </Dialog>

      {/* Dialog: Submit Claim */}
      <Dialog open={isClaimOpen} onOpenChange={(open) => {
        setIsClaimOpen(open)
        if(!open) {
          setClaimIssue('')
        }
      }}>
        <DialogContent className="max-w-md">
          <form onSubmit={handleCreateClaimSubmit}>
            <DialogHeader>
              <DialogTitle className="text-xl font-extrabold text-slate-900 flex items-center gap-2">
                <ClipboardList className="w-5 h-5 text-primary" />
                Create Warranty Claim
              </DialogTitle>
              <DialogDescription>
                Register a new claim under warranty <span className="font-mono font-bold text-slate-900">{foundWarranty?.id}</span>.
              </DialogDescription>
            </DialogHeader>

            <div className="space-y-4 py-4">
              <div className="space-y-1">
                <label className="text-xs font-bold text-slate-500 uppercase">Issue Description</label>
                <Input 
                  required
                  placeholder="Describe the repeating or new issue..."
                  value={claimIssue}
                  onChange={(e) => setClaimIssue(e.target.value)}
                />
              </div>
            </div>

            <DialogFooter className="gap-2 sm:gap-0 pt-2">
              <Button type="button" variant="outline" className="cursor-pointer" onClick={() => setIsClaimOpen(false)}>Cancel</Button>
              <Button type="submit" className="bg-primary hover:bg-secondary cursor-pointer">Register Claim</Button>
            </DialogFooter>
          </form>
        </DialogContent>
      </Dialog>

      {/* Dialog: Evaluate Claim */}
      <Dialog open={isEvaluateOpen} onOpenChange={(open) => {
        setIsEvaluateOpen(open)
        if(!open) {
          setSelectedClaim(null)
          setEvalNotes('')
        }
      }}>
        <DialogContent className="max-w-md">
          <form onSubmit={handleEvaluateSubmit}>
            <DialogHeader>
              <DialogTitle className="text-xl font-extrabold text-slate-900">Evaluate Warranty Claim</DialogTitle>
              <DialogDescription>
                Decide whether to accept, reject, or void claim <span className="font-mono font-bold text-slate-900">{selectedClaim?.claim_number}</span>.
              </DialogDescription>
            </DialogHeader>

            <div className="space-y-4 py-4">
              <div className="space-y-1">
                <label className="text-xs font-bold text-slate-500 uppercase">Decision</label>
                <select 
                  className="w-full bg-white border border-slate-200 rounded-md p-2 text-sm focus:outline-none focus:ring-2 focus:ring-primary/20"
                  value={evalStatus}
                  onChange={(e) => setEvalStatus(e.target.value as ClaimEvaluationStatus)}
                >
                  <option value="ACCEPTED">ACCEPTED (Approve under warranty)</option>
                  <option value="REJECTED">REJECTED (Decline claim)</option>
                  <option value="VOID">VOID (Violates terms, invalidate contract)</option>
                </select>
              </div>

              <div className="space-y-1">
                <label className="text-xs font-bold text-slate-500 uppercase">
                  Notes { (evalStatus === 'REJECTED' || evalStatus === 'VOID') && <span className="text-red-500">*</span> }
                </label>
                <Input 
                  required={evalStatus === 'REJECTED' || evalStatus === 'VOID'}
                  placeholder={evalStatus === 'ACCEPTED' ? "e.g. Free parts replacement" : "Reason is required..."}
                  value={evalNotes}
                  onChange={(e) => setEvalNotes(e.target.value)}
                />
                {evalStatus === 'VOID' && (
                  <p className="text-xxs text-amber-600 font-semibold">
                    * WARNING: Setting status to VOID will automatically set the parent warranty contract status to VOID as well.
                  </p>
                )}
              </div>
            </div>

            <DialogFooter className="gap-2 sm:gap-0 pt-2">
              <Button type="button" variant="outline" className="cursor-pointer" onClick={() => setIsEvaluateOpen(false)}>Cancel</Button>
              <Button type="submit" className="bg-primary hover:bg-secondary cursor-pointer">Submit Decision</Button>
            </DialogFooter>
          </form>
        </DialogContent>
      </Dialog>

      {/* Dialog: Update Status */}
      <Dialog open={isStatusOpen} onOpenChange={(open) => {
        setIsStatusOpen(open)
        if(!open) setSelectedClaim(null)
      }}>
        <DialogContent className="max-w-sm">
          <form onSubmit={handleStatusSubmit}>
            <DialogHeader>
              <DialogTitle className="text-xl font-extrabold text-slate-900">Update Repair Status</DialogTitle>
              <DialogDescription>
                Modify lifecycle status for claim <span className="font-mono font-bold text-slate-900">{selectedClaim?.claim_number}</span>.
              </DialogDescription>
            </DialogHeader>

            <div className="space-y-4 py-4">
              <div className="space-y-1">
                <label className="text-xs font-bold text-slate-500 uppercase">Status</label>
                <select 
                  className="w-full bg-white border border-slate-200 rounded-md p-2 text-sm focus:outline-none focus:ring-2 focus:ring-primary/20"
                  value={claimStatus}
                  onChange={(e) => setClaimStatus(e.target.value as TicketStatus)}
                >
                  <option value="RECEIVED">DITERIMA (RECEIVED)</option>
                  <option value="REPAIRING">DIPERBAIKI (REPAIRING)</option>
                  <option value="PENDING_CONFIRMATION">BUTUH KONFIRMASI (PENDING CONFIRMATION)</option>
                  <option value="FIXED">SELESAI DIKERJAKAN (FIXED)</option>
                  <option value="COMPLETED">DIAMBIL (COMPLETED)</option>
                  <option value="CANCELLED">DIBATALKAN (CANCELLED)</option>
                  <option value="RETURNED">DIKEMBALIKAN (RETURNED)</option>
                </select>
              </div>
            </div>

            <DialogFooter className="gap-2 sm:gap-0 pt-2">
              <Button type="button" variant="outline" className="cursor-pointer" onClick={() => setIsStatusOpen(false)}>Cancel</Button>
              <Button type="submit" className="bg-primary hover:bg-secondary cursor-pointer">Update Status</Button>
            </DialogFooter>
          </form>
        </DialogContent>
      </Dialog>

      {/* Dialog: Update Claim Details */}
      <Dialog open={isEditClaimOpen} onOpenChange={(open) => {
        setIsEditClaimOpen(open)
        if(!open) setSelectedClaim(null)
      }}>
        <DialogContent className="max-w-md">
          <form onSubmit={handleEditClaimSubmit}>
            <DialogHeader>
              <DialogTitle className="text-xl font-extrabold text-slate-900">Update Repair Info</DialogTitle>
              <DialogDescription>
                Edit repair data for claim <span className="font-mono font-bold text-slate-900">{selectedClaim?.claim_number}</span>.
              </DialogDescription>
            </DialogHeader>

            <div className="space-y-4 py-4">
              <div className="space-y-1">
                <label className="text-xs font-bold text-slate-500 uppercase">Issue Description</label>
                <Input 
                  required
                  value={editClaimData.issue_description}
                  onChange={(e) => setEditClaimData({...editClaimData, issue_description: e.target.value})}
                />
              </div>

              <div className="space-y-1">
                <label className="text-xs font-bold text-slate-500 uppercase">Repair Action</label>
                <Input 
                  placeholder="e.g. Reconnected LCD flexible connector"
                  value={editClaimData.repair_action}
                  onChange={(e) => setEditClaimData({...editClaimData, repair_action: e.target.value})}
                />
              </div>

              <div className="space-y-1">
                <label className="text-xs font-bold text-slate-500 uppercase">Notes (Technician)</label>
                <Input 
                  placeholder="e.g. Audio checked, no parts replaced"
                  value={editClaimData.notes}
                  onChange={(e) => setEditClaimData({...editClaimData, notes: e.target.value})}
                />
              </div>
            </div>

            <DialogFooter className="gap-2 sm:gap-0 pt-2">
              <Button type="button" variant="outline" className="cursor-pointer" onClick={() => setIsEditClaimOpen(false)}>Cancel</Button>
              <Button type="submit" className="bg-primary hover:bg-secondary cursor-pointer">Save Changes</Button>
            </DialogFooter>
          </form>
        </DialogContent>
      </Dialog>
    </div>
  )
}

export default WarrantyPage
