import { useState, useEffect, useCallback } from 'react'
import type { Warranty, WarrantyClaim, ClaimEvaluationStatus } from '@/types/warranty'
import { warrantyService } from '@/services/warrantyService'
import { Tabs, TabsList, TabsTrigger, TabsContent } from '@/components/ui/tabs'

import { WarrantyCheckerPanel } from '@/features/warranty/components/WarrantyCheckerPanel'
import { ClaimsQueueTable } from '@/features/warranty/components/ClaimsQueueTable'
import { VoidWarrantyDialog } from '@/features/warranty/components/VoidWarrantyDialog'
import { CreateClaimDialog } from '@/features/warranty/components/CreateClaimDialog'
import { EvaluateClaimDialog } from '@/features/warranty/components/EvaluateClaimDialog'
import { EditClaimDialog } from '@/features/warranty/components/EditClaimDialog'

function WarrantyPage() {
  const [claims, setClaims] = useState<WarrantyClaim[]>([])
  const [searching, setSearching] = useState(false)
  const [searchError, setSearchError] = useState<string | null>(null)

  // Search & Status state
  const [searchTicket, setSearchTicket] = useState('')
  const [foundWarranty, setFoundWarranty] = useState<Warranty | null>(null)
  const [searchedTicketNum, setSearchedTicketNum] = useState('')
  const [hasSearched, setHasSearched] = useState(false)

  // Dialog open states
  const [isVoidOpen, setIsVoidOpen] = useState(false)
  const [isClaimOpen, setIsClaimOpen] = useState(false)
  const [isEvaluateOpen, setIsEvaluateOpen] = useState(false)
  const [isEditClaimOpen, setIsEditClaimOpen] = useState(false)

  // Selected items
  const [selectedWarranty, setSelectedWarranty] = useState<Warranty | null>(null)
  const [selectedClaim, setSelectedClaim] = useState<WarrantyClaim | null>(null)

  // Form states
  const [voidNotes, setVoidNotes] = useState('')
  const [claimIssue, setClaimIssue] = useState('')
  const [evalStatus, setEvalStatus] = useState<ClaimEvaluationStatus>('ACCEPTED')
  const [evalNotes, setEvalNotes] = useState('')

  const [editClaimData, setEditClaimData] = useState({
    issue_description: '',
    notes: '',
  })

  // Tab State
  const [activeTab, setActiveTab] = useState<'check' | 'queue'>('check')

  const fetchClaims = useCallback(async () => {
    try {
      const result = await warrantyService.getClaims({ limit: 50 })
      setClaims(result.data)
    } catch (err: any) {
      console.error('Failed to load claims', err)
    }
  }, [])

  useEffect(() => {
    fetchClaims()
  }, [fetchClaims])

  // Search Logic
  const handleSearch = async (overrideTicket?: string) => {
    const ticketToSearch = overrideTicket !== undefined ? overrideTicket : searchTicket
    if (!ticketToSearch.trim()) return

    setSearching(true)
    setSearchError(null)
    try {
      const warranty = await warrantyService.getWarrantyByTicketNumber(ticketToSearch.trim())
      setFoundWarranty(warranty)
      setSearchedTicketNum(ticketToSearch.trim())
      setHasSearched(true)
    } catch (err: any) {
      setFoundWarranty(null)
      setSearchedTicketNum(ticketToSearch.trim())
      setHasSearched(true)
      setSearchError(err?.response?.data?.detail || err?.message || 'Warranty not found')
    } finally {
      setSearching(false)
    }
  }

  // Void Warranty Logic
  const handleVoidWarrantySubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!selectedWarranty) return

    try {
      const updated = await warrantyService.updateWarrantyStatus(selectedWarranty.id, {
        status: 'VOID',
        notes: voidNotes,
      })
      if (foundWarranty && foundWarranty.id === updated.id) {
        setFoundWarranty(updated)
      }
      setIsVoidOpen(false)
      setSelectedWarranty(null)
      setVoidNotes('')
    } catch (err: any) {
      alert(err?.response?.data?.detail || 'Failed to void warranty')
    }
  }

  // Create Claim Logic
  const handleCreateClaimSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!searchedTicketNum) return

    try {
      const newClaim = await warrantyService.createClaim({
        ticket_number: searchedTicketNum,
        issue_description: claimIssue,
      })
      setClaims([newClaim, ...claims])
      setIsClaimOpen(false)
      setClaimIssue('')
      setActiveTab('queue')
    } catch (err: any) {
      alert(err?.response?.data?.detail || 'Failed to create claim')
    }
  }

  // Evaluate Claim Logic
  const handleEvaluateSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!selectedClaim) return

    if ((evalStatus === 'REJECTED' || evalStatus === 'VOID') && !evalNotes.trim()) {
      alert('Notes are required when rejecting or voiding a claim')
      return
    }

    try {
      const updated = await warrantyService.evaluateClaim(selectedClaim.claim_id, {
        status: evalStatus,
        notes: evalStatus === 'ACCEPTED' ? evalNotes || undefined : evalNotes,
      })
      setClaims(claims.map(c => (c.claim_id === updated.claim_id ? updated : c)))

      if (evalStatus === 'VOID' && foundWarranty && foundWarranty.id === selectedClaim.warranty_id) {
        const refetched = await warrantyService.getWarrantyByTicketNumber(searchedTicketNum)
        setFoundWarranty(refetched)
      }

      setIsEvaluateOpen(false)
      setSelectedClaim(null)
      setEvalNotes('')
    } catch (err: any) {
      alert(err?.response?.data?.detail || 'Failed to evaluate claim')
    }
  }

  // Update claim info logic
  const handleEditClaimSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!selectedClaim) return

    try {
      const updated = await warrantyService.updateClaim(selectedClaim.claim_id, {
        issue_description: editClaimData.issue_description,
        notes: editClaimData.notes || undefined,
      })
      setClaims(claims.map(c => (c.claim_id === updated.claim_id ? updated : c)))
      setIsEditClaimOpen(false)
      setSelectedClaim(null)
    } catch (err: any) {
      alert(err?.response?.data?.detail || 'Failed to update claim')
    }
  }

  return (
    <div className="space-y-8">
      {/* Title */}
      <div>
        <h1 className="text-3xl font-extrabold text-slate-900 dark:text-slate-100 tracking-tight">
          Warranty & Claim Ticketing
        </h1>
        <p className="text-slate-500 dark:text-slate-400 text-sm">
          Validate customer warranty cards, void contracts, and process claims separately from regular ticketing.
        </p>
      </div>

      <Tabs value={activeTab} onValueChange={(v) => setActiveTab(v as 'check' | 'queue')} className="space-y-6">
        <TabsList className="bg-slate-100 dark:bg-slate-800/60 p-1 rounded-xl">
          <TabsTrigger
            value="check"
            className="rounded-lg px-4 py-2 text-xs font-semibold cursor-pointer dark:data-[state=active]:bg-slate-900 dark:data-[state=active]:text-slate-100"
          >
            Check & Claim Warranty
          </TabsTrigger>
          <TabsTrigger
            value="queue"
            className="rounded-lg px-4 py-2 text-xs font-semibold cursor-pointer dark:data-[state=active]:bg-slate-900 dark:data-[state=active]:text-slate-100"
          >
            Claims Queue
          </TabsTrigger>
        </TabsList>

        {/* Tab 1: Check Warranty */}
        <TabsContent value="check" className="space-y-6">
          <WarrantyCheckerPanel
            searchTicket={searchTicket}
            onSearchChange={(val) => {
              setSearchTicket(val)
            }}
            onSearch={handleSearch}
            searching={searching}
            searchError={searchError}
            hasSearched={hasSearched}
            foundWarranty={foundWarranty}
            searchedTicketNum={searchedTicketNum}
            onOpenClaimDialog={() => setIsClaimOpen(true)}
            onOpenVoidDialog={(w) => {
              setSelectedWarranty(w)
              setIsVoidOpen(true)
            }}
          />
        </TabsContent>

        {/* Tab 2: Claims Queue */}
        <TabsContent value="queue" className="space-y-6">
          <ClaimsQueueTable
            claims={claims}
            onEvaluate={(claim) => {
              setSelectedClaim(claim)
              setEvalStatus('ACCEPTED')
              setIsEvaluateOpen(true)}
            }
            onEdit={(claim) => {
              setSelectedClaim(claim)
              setEditClaimData({
                issue_description: claim.issue_description,
                notes: claim.notes || '',
              })
              setIsEditClaimOpen(true)
            }}
          />
        </TabsContent>
      </Tabs>

      {/* Dialogs */}
      <VoidWarrantyDialog
        open={isVoidOpen}
        onOpenChange={(open) => {
          setIsVoidOpen(open)
          if (!open) {
            setSelectedWarranty(null)
            setVoidNotes('')
          }
        }}
        warranty={selectedWarranty}
        notes={voidNotes}
        onNotesChange={setVoidNotes}
        onSubmit={handleVoidWarrantySubmit}
      />

      <CreateClaimDialog
        open={isClaimOpen}
        onOpenChange={(open) => {
          setIsClaimOpen(open)
          if (!open) {
            setClaimIssue('')
          }
        }}
        warranty={foundWarranty}
        issueDescription={claimIssue}
        onIssueChange={setClaimIssue}
        onSubmit={handleCreateClaimSubmit}
      />

      <EvaluateClaimDialog
        open={isEvaluateOpen}
        onOpenChange={(open) => {
          setIsEvaluateOpen(open)
          if (!open) {
            setSelectedClaim(null)
            setEvalNotes('')
          }
        }}
        claim={selectedClaim}
        evalStatus={evalStatus}
        onEvalStatusChange={setEvalStatus}
        evalNotes={evalNotes}
        onEvalNotesChange={setEvalNotes}
        onSubmit={handleEvaluateSubmit}
      />

      <EditClaimDialog
        open={isEditClaimOpen}
        onOpenChange={(open) => {
          setIsEditClaimOpen(open)
          if (!open) setSelectedClaim(null)
        }}
        claim={selectedClaim}
        formData={editClaimData}
        onFormDataChange={setEditClaimData}
        onSubmit={handleEditClaimSubmit}
      />
    </div>
  )
}

export default WarrantyPage
