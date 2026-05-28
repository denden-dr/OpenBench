<script lang="ts">
  import { onMount } from 'svelte';
  import { 
    Search, ShieldCheck, ShieldAlert, 
    ArrowLeft, LoaderCircle, Info, AlertOctagon, CheckCircle2, X,
    FileText, User, Clock, AlertTriangle, Play, RefreshCw,
    BookOpen, ChevronDown, ChevronUp
  } from 'lucide-svelte';

  import type { Ticket, Claim } from '$lib/types/ticket';
  import { formatCurrency } from '$lib/utils/format';
  import { checkWarrantyExpiry } from '$lib/utils/warranty';
  import Pagination from '../_components/Pagination.svelte';

  function checkSuccess(res: Response, payload: any): boolean {
    return payload && (payload.success === true || (res.ok && (payload.code === undefined || (payload.code >= 200 && payload.code < 300))));
  }

  function getErrorMessage(payload: any, fallback: string): string {
    return payload?.detail || payload?.title || payload?.message || payload?.error || fallback;
  }

  let searchQuery = $state('');
  let isLoading = $state(false);
  let searchResult = $state<Ticket | null>(null);
  let searchCandidates = $state<Ticket[]>([]);
  let notFound = $state(false);
  let searchInput = $state<HTMLInputElement | null>(null);

  // Queue List State
  let claimsQueue = $state<Claim[]>([]);
  let isQueueLoading = $state(false);
  let queuePage = $state(1);
  let queueTotalPages = $state(1);
  let queueTotalItems = $state(0);

  // Intake Form State
  let intakeIssue = $state('');
  let intakeNote = $state('');
  let isSubmittingIntake = $state(false);

  // Void Modal State
  let showVoidModal = $state(false);
  let selectedClaimForVoid = $state<Claim | null>(null);
  let voidReason = $state('');
  let isSubmittingVoid = $state(false);

  // Operational Guide State
  let isGuideExpanded = $state(true);

  // Toast state
  interface Toast {
    id: string;
    message: string;
    type: 'success' | 'error';
  }
  let toasts = $state<Toast[]>([]);

  function addToast(message: string, type: 'success' | 'error') {
    const id = crypto.randomUUID();
    toasts = [...toasts, { id, message, type }];
    setTimeout(() => {
      toasts = toasts.filter(t => t.id !== id);
    }, 5000);
  }

  onMount(() => {
    searchInput?.focus();
  });

  async function fetchClaims() {
    isQueueLoading = true;
    try {
      const params = new URLSearchParams();
      params.set('status', 'waiting_inspection');
      params.set('page', String(queuePage));
      params.set('limit', '10');

      const claimsRes = await fetch(`/api/v1/warranty-claims?${params}`);
      const claimsPayload = await claimsRes.json();

      if (checkSuccess(claimsRes, claimsPayload)) {
        claimsQueue = claimsPayload.data || [];
        queueTotalPages = claimsPayload.total_pages || 1;
        queueTotalItems = claimsPayload.total ?? claimsQueue.length;
      }
    } catch (err) {
      console.error('Error fetching claims queue:', err);
    } finally {
      isQueueLoading = false;
    }
  }

  // Reactive effect for fetching claims queue when page changes
  $effect(() => {
    fetchClaims();
  });

  async function handleSearch(e: SubmitEvent) {
    e.preventDefault();
    if (!searchQuery.trim()) return;
    
    isLoading = true;
    notFound = false;
    searchResult = null;
    searchCandidates = [];

    try {
      const params = new URLSearchParams();
      params.set('status', 'picked_up');
      params.set('search', searchQuery.trim());
      params.set('limit', '5');

      const res = await fetch(`/api/v1/tickets?${params}`);
      const payload = await res.json();
      
      if (checkSuccess(res, payload) && payload.data && payload.data.length === 1) {
        searchResult = payload.data[0];
        intakeIssue = '';
        intakeNote = '';
      } else if (checkSuccess(res, payload) && payload.data && payload.data.length > 1) {
        searchCandidates = payload.data;
      } else {
        notFound = true;
      }
    } catch (err) {
      console.error('Gagal memverifikasi nota:', err);
      notFound = true;
    } finally {
      isLoading = false;
    }
  }

  function selectSearchCandidate(ticket: Ticket) {
    searchResult = ticket;
    searchCandidates = [];
    intakeIssue = '';
    intakeNote = '';
  }

  // Custom Approve Modal State
  let showApproveModal = $state(false);
  let selectedClaimForApprove = $state<Claim | null>(null);
  let isSubmittingApprove = $state(false);

  function triggerApproveModal(claim: Claim) {
    selectedClaimForApprove = claim;
    showApproveModal = true;
  }

  async function submitIntake() {
    if (!searchResult) return;
    if (!intakeIssue.trim()) return;

    isSubmittingIntake = true;

    try {
      const res = await fetch('/api/v1/warranty-claims', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          ticket_id: searchResult.id,
          issue: intakeIssue,
          additional_description: intakeNote
        })
      });

      const payload = await res.json();
      if (checkSuccess(res, payload)) {
        addToast(`Berhasil mendaftarkan klaim garansi untuk nota ${searchResult.id}. Perangkat siap diletakkan di antrean inspeksi teknisi.`, 'success');
        searchResult = null;
        searchQuery = '';
        fetchClaims();
      } else {
        addToast(getErrorMessage(payload, 'Pendaftaran klaim gagal.'), 'error');
      }
    } catch (err) {
      console.error(err);
      addToast('Koneksi ke mock API gagal.', 'error');
    } finally {
      isSubmittingIntake = false;
    }
  }

  async function approveClaim() {
    if (!selectedClaimForApprove) return;
    
    isSubmittingApprove = true;
    try {
      const claimId = selectedClaimForApprove.id;
      const res = await fetch(`/api/v1/warranty-claims/${claimId}/approve`, {
        method: 'POST'
      });
      const payload = await res.json();
      if (checkSuccess(res, payload)) {
        addToast(`Klaim disetujui! Tiket baru ${payload.data.ticket.id} dengan harga Rp 0 telah dibuat.`, 'success');
        showApproveModal = false;
        selectedClaimForApprove = null;
        fetchClaims();
      } else {
        addToast(getErrorMessage(payload, 'Gagal menyetujui klaim.'), 'error');
      }
    } catch (err) {
      console.error(err);
      addToast('Koneksi ke mock API gagal.', 'error');
    } finally {
      isSubmittingApprove = false;
    }
  }

  function openVoidModal(claim: Claim) {
    selectedClaimForVoid = claim;
    voidReason = '';
    showVoidModal = true;
  }

  async function submitVoid() {
    if (!selectedClaimForVoid || !voidReason.trim()) return;

    isSubmittingVoid = true;
    try {
      const res = await fetch(`/api/v1/warranty-claims/${selectedClaimForVoid.id}/void`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          void_reason: voidReason
        })
      });

      const payload = await res.json();
      if (checkSuccess(res, payload)) {
        addToast(`Klaim ditolak (Void). Tiket baru ${payload.data.ticket.id} dibuat dengan status dibatalkan beserta alasannya.`, 'success');
        showVoidModal = false;
        fetchClaims();
      } else {
        addToast(getErrorMessage(payload, 'Gagal membatalkan klaim.'), 'error');
      }
    } catch (err) {
      console.error(err);
      addToast('Koneksi ke mock API gagal.', 'error');
    } finally {
      isSubmittingVoid = false;
    }
  }
</script>

<div class="min-h-screen bg-slate-50 dark:bg-slate-900 text-slate-900 dark:text-slate-100 p-6">
  <div class="max-w-4xl mx-auto space-y-6">
    
    <!-- Top Header -->
    <div class="flex flex-col sm:flex-row sm:justify-between sm:items-center gap-4">
      <a href="/" class="flex items-center gap-1.5 text-sm font-bold text-slate-600 dark:text-slate-400 hover:text-blue-600 dark:hover:text-blue-400 transition-colors">
        <ArrowLeft size={16} />
        Kembali ke Dashboard
      </a>
      <h1 class="font-extrabold text-xl tracking-tight text-slate-900 dark:text-white flex items-center gap-2">
        <ShieldCheck class="text-indigo-600 animate-pulse" size={24} />
        Verifikasi Nota & Klaim Garansi
      </h1>
    </div>



    <!-- Collapsible Operational Guide -->
    <div class="bg-white dark:bg-slate-950 border border-slate-200 dark:border-slate-800 rounded-2xl overflow-hidden shadow-sm transition-all duration-300">
      <button 
        onclick={() => isGuideExpanded = !isGuideExpanded}
        class="w-full px-6 py-4 flex justify-between items-center bg-slate-50/50 dark:bg-slate-900/10 hover:bg-slate-50 dark:hover:bg-slate-900 transition-colors text-left cursor-pointer"
      >
        <div class="flex items-center gap-2 text-slate-800 dark:text-slate-200">
          <BookOpen size={18} class="text-indigo-600" />
          <h3 class="font-bold text-sm">Panduan Alur Kerja Klaim Garansi (Kasir & Antrean)</h3>
        </div>
        <div class="text-slate-400">
          {#if isGuideExpanded}
            <ChevronUp size={16} />
          {:else}
            <ChevronDown size={16} />
          {/if}
        </div>
      </button>

      {#if isGuideExpanded}
        <div class="p-6 border-t border-slate-100 dark:border-slate-900 grid grid-cols-1 sm:grid-cols-2 md:grid-cols-4 gap-4 animate-fade-in">
          <!-- Step 1 -->
          <div class="p-4 bg-slate-50 dark:bg-slate-900/40 rounded-xl border border-slate-100 dark:border-slate-800/60 relative">
            <span class="absolute top-3 right-3 text-xl font-black text-indigo-500/10 font-mono tracking-tighter">01</span>
            <h4 class="font-bold text-xs text-indigo-650 dark:text-indigo-400 mb-1.5 uppercase tracking-wider">Cari Nota</h4>
            <p class="text-[11px] text-slate-500 leading-relaxed">Cari nota lama pelanggan di form verifikasi di bawah.</p>
          </div>

          <!-- Step 2 -->
          <div class="p-4 bg-slate-50 dark:bg-slate-900/40 rounded-xl border border-slate-100 dark:border-slate-800/60 relative">
            <span class="absolute top-3 right-3 text-xl font-black text-indigo-500/10 font-mono tracking-tighter">02</span>
            <h4 class="font-bold text-xs text-indigo-650 dark:text-indigo-400 mb-1.5 uppercase tracking-wider">Daftar Klaim</h4>
            <p class="text-[11px] text-slate-500 leading-relaxed">Isi kerusakan baru dan daftarkan klaim ke antrean.</p>
          </div>

          <!-- Step 3 -->
          <div class="p-4 bg-slate-50 dark:bg-slate-900/40 rounded-xl border border-slate-100 dark:border-slate-800/60 relative">
            <span class="absolute top-3 right-3 text-xl font-black text-indigo-500/10 font-mono tracking-tighter">03</span>
            <h4 class="font-bold text-xs text-indigo-650 dark:text-indigo-400 mb-1.5 uppercase tracking-wider">Rak Inspeksi</h4>
            <p class="text-[11px] text-slate-500 leading-relaxed">Letakkan perangkat fisik di rak antrean inspeksi teknisi.</p>
          </div>

          <!-- Step 4 -->
          <div class="p-4 bg-slate-50 dark:bg-slate-900/40 rounded-xl border border-slate-100 dark:border-slate-800/60 relative">
            <span class="absolute top-3 right-3 text-xl font-black text-indigo-500/10 font-mono tracking-tighter">04</span>
            <h4 class="font-bold text-xs text-indigo-650 dark:text-indigo-400 mb-1.5 uppercase tracking-wider">Setujui/Tolak</h4>
            <p class="text-[11px] text-slate-500 leading-relaxed">Teknisi memeriksa secara offline lalu menyetujui/menolak.</p>
          </div>
        </div>
      {/if}
    </div>

    <!-- Search & Verification Intake -->
    <div class="space-y-6">
        <!-- Search Form -->
        <div class="bg-white dark:bg-slate-950 border border-slate-200 dark:border-slate-800 rounded-2xl p-5 shadow-sm space-y-4">
          <h3 class="font-bold text-xs text-slate-400 uppercase tracking-wider block">Verifikasi Nota / Cari Nota Servis</h3>
          <form onsubmit={handleSearch} class="flex gap-2">
            <div class="relative flex-1">
              <span class="absolute inset-y-0 left-0 pl-3.5 flex items-center text-slate-400">
                <Search size={18} />
              </span>
              <input
                type="text"
                bind:this={searchInput}
                bind:value={searchQuery}
                required
                placeholder="ID Tiket, nama pelanggan, model..."
                class="w-full bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-xl py-2.5 pl-10 pr-4 text-sm focus:outline-none focus:border-indigo-600 dark:focus:border-indigo-500 transition-colors"
              />
            </div>
            <button type="submit" disabled={isLoading} class="px-5 py-2.5 bg-indigo-600 hover:bg-indigo-700 disabled:bg-slate-100 text-white font-bold text-xs uppercase tracking-wider rounded-xl transition-all shadow-sm active:scale-95 inline-flex items-center gap-1.5">
              {#if isLoading}
                <LoaderCircle size={14} class="animate-spin" />
                Mencari...
              {:else}
                Cari
              {/if}
            </button>
          </form>
        </div>

        <!-- Search Candidates Disambiguation -->
        {#if searchCandidates.length > 0}
          <div class="bg-white dark:bg-slate-950 border border-slate-200 dark:border-slate-800 rounded-2xl p-6 shadow-sm space-y-4 animate-fade-in">
            <h4 class="font-bold text-xs text-slate-400 uppercase tracking-wider block">Beberapa Nota Ditemukan</h4>
            <p class="text-xs text-slate-500">Pilih nota yang sesuai dengan perangkat pelanggan:</p>
            <div class="space-y-2.5">
              {#each searchCandidates as candidate}
                <button
                  onclick={() => selectSearchCandidate(candidate)}
                  class="w-full text-left p-3.5 bg-slate-50 dark:bg-slate-900 hover:bg-indigo-50/50 dark:hover:bg-indigo-950/20 border border-slate-200 dark:border-slate-800 hover:border-indigo-200 dark:hover:border-indigo-900 rounded-xl transition-all duration-200 cursor-pointer active:scale-[0.99] flex flex-col sm:flex-row justify-between items-start sm:items-center gap-2"
                >
                  <div>
                    <span class="font-bold text-sm text-slate-900 dark:text-white">{candidate.brand} {candidate.model}</span>
                    <span class="block text-[10px] font-mono text-slate-400 mt-0.5">ID: {candidate.id}</span>
                  </div>
                  <div class="text-right">
                    <span class="text-xs font-semibold text-slate-700 dark:text-slate-300 block">{candidate.customer_name}</span>
                    <span class="text-[10px] text-slate-400">{candidate.entry_date ? new Date(candidate.entry_date).toLocaleDateString('id-ID') : ''}</span>
                  </div>
                </button>
              {/each}
            </div>
          </div>
        {/if}

        <!-- Search Result & Intake Form -->
        {#if searchResult}
          {@const warranty = searchResult.exit_date ? checkWarrantyExpiry(searchResult.exit_date, searchResult.warranty_days) : null}
          <div class="bg-white dark:bg-slate-950 border border-slate-200 dark:border-slate-800 rounded-2xl overflow-hidden shadow-sm space-y-5 p-6 animate-fade-in">
            <div class="flex justify-between items-start">
              <div>
                <h2 class="font-bold text-lg text-slate-900 dark:text-white">{searchResult.brand} {searchResult.model}</h2>
                <p class="text-xs font-mono text-slate-400 mt-0.5">ID Nota: {searchResult.id}</p>
              </div>
              
              {#if searchResult.status !== 'picked_up'}
                <span class="inline-flex items-center gap-1 px-3 py-1 rounded-full text-xs font-bold bg-slate-100 text-slate-500 dark:bg-slate-900 dark:text-slate-400 border border-slate-200 dark:border-slate-800">
                  <Info size={14} /> Belum Diambil
                </span>
              {:else if warranty && warranty.isValid}
                <span class="inline-flex items-center gap-1 px-3 py-1 rounded-full text-xs font-bold bg-emerald-50 text-emerald-700 dark:bg-emerald-950/30 dark:text-emerald-400 border border-emerald-200 dark:border-emerald-900/50 animate-bounce">
                  <ShieldCheck size={14} /> Garansi Valid
                </span>
              {:else}
                <span class="inline-flex items-center gap-1 px-3 py-1 rounded-full text-xs font-bold bg-rose-50 text-rose-700 dark:bg-rose-950/30 dark:text-rose-400 border border-rose-200 dark:border-rose-900/50">
                  <ShieldAlert size={14} /> Garansi Expired
                </span>
              {/if}
            </div>

            <!-- Specifications Grid -->
            <div class="grid grid-cols-2 gap-4 text-xs border-t border-slate-100 dark:border-slate-850 pt-4">
              <div class="p-3 bg-slate-50 dark:bg-slate-900/40 rounded-xl border border-slate-100 dark:border-slate-800/60">
                <span class="text-[9px] font-bold text-slate-400 uppercase tracking-wider block mb-1">Pelanggan</span>
                <span class="font-bold text-slate-800 dark:text-slate-200 flex items-center gap-1.5">
                  <User size={14} class="text-slate-400" />
                  {searchResult.customer_name} ({searchResult.customer_gender})
                </span>
              </div>
              <div class="p-3 bg-slate-50 dark:bg-slate-900/40 rounded-xl border border-slate-100 dark:border-slate-800/60">
                <span class="text-[9px] font-bold text-slate-400 uppercase tracking-wider block mb-1">Status Pembayaran</span>
                <span class="font-bold text-slate-800 dark:text-slate-200 flex items-center gap-1.5 uppercase">
                  <span class="w-1.5 h-1.5 rounded-full {searchResult.payment_status === 'paid' ? 'bg-emerald-500' : 'bg-amber-500'}"></span>
                  {searchResult.payment_status}
                </span>
              </div>
            </div>

            <!-- Timeline visual -->
            {#if searchResult.status === 'picked_up' && warranty}
              <div class="border-t border-slate-100 dark:border-slate-850 pt-4 space-y-4">
                <span class="text-[10px] font-bold text-slate-400 uppercase tracking-wider block">Garansi Timeline</span>
                
                <div class="relative py-2">
                  <!-- Connecting Line -->
                  <div class="absolute left-4 right-4 top-1/2 -translate-y-1/2 h-0.5 bg-slate-200 dark:bg-slate-800"></div>
                  <div 
                    class="absolute left-4 top-1/2 -translate-y-1/2 h-0.5 bg-indigo-600 transition-all duration-500"
                    style="width: {warranty.isValid ? '75%' : '100%'}"
                  ></div>

                  <!-- Timeline Nodes -->
                  <div class="relative flex justify-between text-[10px]">
                    <!-- Node 1: Entry -->
                    <div class="flex flex-col items-center">
                      <div class="w-6 h-6 rounded-full bg-indigo-600 text-white flex items-center justify-center font-bold relative z-10 shadow border-4 border-white dark:border-slate-950">
                        ✓
                      </div>
                      <span class="font-semibold text-slate-800 dark:text-slate-200 mt-1">Masuk</span>
                      <span class="text-[9px] text-slate-400 mt-0.5">{new Date(searchResult.entry_date).toLocaleDateString('id-ID', {day: 'numeric', month: 'short'})}</span>
                    </div>

                    <!-- Node 2: Picked Up -->
                    <div class="flex flex-col items-center">
                      <div class="w-6 h-6 rounded-full bg-indigo-600 text-white flex items-center justify-center font-bold relative z-10 shadow border-4 border-white dark:border-slate-950">
                        ✓
                      </div>
                      <span class="font-semibold text-slate-800 dark:text-slate-200 mt-1">Diambil</span>
                      <span class="text-[9px] text-slate-400 mt-0.5">{searchResult.exit_date ? new Date(searchResult.exit_date).toLocaleDateString('id-ID', {day: 'numeric', month: 'short'}) : ''}</span>
                    </div>

                    <!-- Node 3: Warranty Expiry -->
                    <div class="flex flex-col items-center">
                      <div 
                        class="w-6 h-6 rounded-full flex items-center justify-center font-bold relative z-10 shadow border-4 border-white dark:border-slate-950 text-white
                          {warranty.isValid ? 'bg-emerald-500' : 'bg-rose-500'}"
                      >
                        🛡️
                      </div>
                      <span class="font-semibold text-slate-800 dark:text-slate-200 mt-1">Masa Garansi</span>
                      <span class="text-[9px] text-slate-400 mt-0.5">{new Date(warranty.expiryDate).toLocaleDateString('id-ID', {day: 'numeric', month: 'short'})}</span>
                    </div>
                  </div>
                </div>

                <!-- Warranty Status Card -->
                <div 
                  class="p-4 rounded-2xl flex items-center justify-between border
                    {warranty.isValid 
                      ? 'bg-emerald-50/30 dark:bg-emerald-950/10 border-emerald-200/60 dark:border-emerald-900/30 text-emerald-900 dark:text-emerald-300' 
                      : 'bg-rose-50/30 dark:bg-rose-950/10 border-rose-200/60 dark:border-rose-900/30 text-rose-900 dark:text-rose-300'}"
                >
                  <div class="flex items-center gap-3">
                    {#if warranty.isValid}
                      <ShieldCheck size={24} class="text-emerald-500" />
                    {:else}
                      <ShieldAlert size={24} class="text-rose-500" />
                    {/if}
                    <div>
                      <span class="text-[9px] font-bold uppercase tracking-wider block opacity-70">Status Masa Garansi</span>
                      <span class="text-xs font-bold">{warranty.isValid ? `Valid - Sisa ${warranty.remainingDays} Hari` : 'Sudah Berakhir (Expired)'}</span>
                    </div>
                  </div>
                  <span class="text-[10px] font-bold opacity-80">{warranty.formattedExpiry}</span>
                </div>
              </div>
            {/if}

            <!-- Intake Registration Section -->
            {#if searchResult.status === 'picked_up' && warranty && warranty.isValid}
              <div class="border-t border-slate-100 dark:border-slate-850 pt-5 space-y-4">
                <h4 class="font-bold text-xs uppercase tracking-wider text-indigo-600">Pendaftaran Klaim Garansi Baru</h4>
                <div class="space-y-3">
                  <div class="space-y-1">
                    <label for="intake-issue" class="text-xs font-bold text-slate-500">Kerusakan / Masalah Baru *</label>
                    <input
                      id="intake-issue"
                      type="text"
                      bind:value={intakeIssue}
                      required
                      placeholder="Contoh: Layar sentuh macet, Kamera buram"
                      class="w-full bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-xl py-2 px-3 text-sm focus:outline-none focus:border-indigo-600 transition-colors"
                    />
                  </div>
                  <div class="space-y-1">
                    <label for="intake-note" class="text-xs font-bold text-slate-500">Catatan Tambahan (Kondisi Fisik / Kelengkapan)</label>
                    <textarea
                      id="intake-note"
                      bind:value={intakeNote}
                      rows="2"
                      placeholder="Contoh: Casing lecet pemakaian, segel masih utuh tanpa modifikasi"
                      class="w-full bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-xl py-2 px-3 text-sm focus:outline-none focus:border-indigo-600 transition-colors resize-none"
                    ></textarea>
                  </div>
                  <div class="flex justify-end pt-2">
                    <button 
                      onclick={submitIntake} 
                      disabled={isSubmittingIntake || !intakeIssue.trim()} 
                      class="px-6 py-2.5 bg-indigo-600 hover:bg-indigo-700 disabled:bg-slate-100 disabled:text-slate-400 text-white font-bold text-xs uppercase tracking-wider rounded-xl transition-colors shadow-sm inline-flex items-center gap-1.5"
                    >
                      {#if isSubmittingIntake}
                        <LoaderCircle size={12} class="animate-spin" />
                        Mendaftarkan...
                      {:else}
                        Daftarkan Antrean Inspeksi
                      {/if}
                    </button>
                  </div>
                </div>
              </div>
            {/if}
          </div>
        {/if}

        <!-- Not Found -->
        {#if notFound}
          <div class="p-8 text-center bg-white dark:bg-slate-950 border border-slate-200 dark:border-slate-800 rounded-2xl shadow-sm space-y-2 animate-fade-in">
            <AlertOctagon class="mx-auto text-slate-300 dark:text-slate-700" size={40} />
            <h3 class="font-bold text-slate-800 dark:text-slate-200">Nota Tidak Ditemukan</h3>
            <p class="text-xs text-slate-500 dark:text-slate-400">Silakan periksa kembali ID nota atau nama pelanggan yang dimasukkan.</p>
          </div>
        {/if}
      </div>



    <!-- BOTTOM ROW: WORK LIST QUEUE (TABLE ANTREAN KLAIM) -->
    <div class="bg-white dark:bg-slate-950 border border-slate-200 dark:border-slate-800 rounded-2xl overflow-hidden shadow-sm" data-claim-queue>
      <div class="px-6 py-4 border-b border-slate-200 dark:border-slate-800 flex justify-between items-center bg-slate-50/50 dark:bg-slate-900/10">
        <div>
          <h3 class="font-bold text-sm text-slate-900 dark:text-white flex items-center gap-2">
            <Clock size={16} class="text-indigo-600" />
            Antrean Klaim Menunggu Inspeksi
          </h3>
          <p class="text-xs text-slate-500 mt-0.5">Daftar perangkat garansi yang menunggu keputusan dari teknisi</p>
        </div>
        <button onclick={fetchClaims} disabled={isQueueLoading} class="p-2 border border-slate-200 dark:border-slate-800 rounded-xl hover:bg-slate-50 dark:hover:bg-slate-900 text-slate-600 dark:text-slate-400 transition-colors">
          <RefreshCw size={14} class={isQueueLoading ? 'animate-spin' : ''} />
        </button>
      </div>

      {#if claimsQueue.length === 0}
        <div class="py-16 text-center text-slate-400 dark:text-slate-600 bg-slate-50/50 dark:bg-slate-900/10 rounded-2xl border border-dashed border-slate-200 dark:border-slate-800 m-6 animate-fade-in">
          <CheckCircle2 size={36} class="mx-auto text-emerald-500 mb-3 opacity-80" />
          <p class="font-bold text-sm text-slate-700 dark:text-slate-300">Semua antrean inspeksi bersih</p>
          <p class="text-xs mt-0.5 opacity-80">Tidak ada klaim garansi yang tertunda saat ini.</p>
        </div>
      {:else}
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4 p-6">
          {#each claimsQueue as claim (claim.id)}
            <div class="bg-white dark:bg-slate-950 border border-slate-200 dark:border-slate-800 hover:border-indigo-300 dark:hover:border-indigo-900 rounded-2xl p-5 shadow-sm hover:shadow transition-all duration-300 flex flex-col justify-between gap-4 animate-fade-in">
              <!-- Card Header -->
              <div class="flex justify-between items-start gap-2">
                <div>
                  <div class="flex items-center gap-1.5">
                    <span class="text-[10px] font-bold text-slate-400 dark:text-slate-500 uppercase tracking-wider font-mono">Nota ID</span>
                    <span class="px-2 py-0.5 bg-slate-100 dark:bg-slate-900 text-slate-700 dark:text-slate-300 rounded font-mono font-bold text-[10px]">{claim.ticket_id}</span>
                  </div>
                  <h4 class="font-extrabold text-sm text-slate-900 dark:text-white mt-1">
                    {claim.originalTicket?.brand || ''} {claim.originalTicket?.model || ''}
                  </h4>
                </div>
                
                <span class="text-[10px] text-slate-400 flex items-center gap-1 shrink-0">
                  <Clock size={12} />
                  {new Date(claim.created_at).toLocaleDateString('id-ID', {day: 'numeric', month: 'short', hour: '2-digit', minute: '2-digit'})}
                </span>
              </div>

              <!-- Card Body -->
              <div class="space-y-2 border-t border-slate-100 dark:border-slate-900 pt-3 text-xs">
                <div class="flex items-center gap-2">
                  <User size={14} class="text-slate-400 shrink-0" />
                  <span class="font-semibold text-slate-700 dark:text-slate-350">{claim.originalTicket?.customer_name || ''}</span>
                </div>
                <div class="p-3 bg-indigo-50/30 dark:bg-indigo-950/10 border border-indigo-100/30 dark:border-indigo-900/20 rounded-xl space-y-1">
                  <span class="text-[9px] font-bold text-indigo-600 dark:text-indigo-400 uppercase tracking-wider block">Kerusakan Klaim</span>
                  <p class="font-bold text-slate-800 dark:text-slate-200">{claim.issue}</p>
                  {#if claim.additional_description}
                    <p class="text-[10px] text-slate-400 mt-1 leading-relaxed border-t border-indigo-100/20 pt-1">
                      <span class="font-semibold">Catatan:</span> {claim.additional_description}
                    </p>
                  {/if}
                </div>
              </div>

              <!-- Card Actions Footer -->
              <div class="flex items-center justify-end gap-2 border-t border-slate-100 dark:border-slate-900 pt-3">
                <button 
                  onclick={() => openVoidModal(claim)} 
                  class="px-4 py-2 border border-rose-200 hover:bg-rose-50 dark:border-rose-900/30 dark:hover:bg-rose-950/20 text-rose-600 dark:text-rose-400 font-bold text-xs rounded-xl transition-all inline-flex items-center gap-1.5 active:scale-95 cursor-pointer"
                >
                  <X size={14} /> Tolak (Void)
                </button>
                <button 
                  onclick={() => triggerApproveModal(claim)} 
                  class="px-4 py-2 bg-emerald-600 hover:bg-emerald-700 text-white font-bold text-xs rounded-xl transition-all shadow-sm shadow-emerald-500/10 inline-flex items-center gap-1.5 active:scale-95 cursor-pointer"
                >
                  <ShieldCheck size={14} /> Setujui Klaim
                </button>
              </div>
            </div>
          {/each}
        </div>
      {/if}

      {#if !isQueueLoading && claimsQueue.length > 0}
        <div class="px-4">
          <Pagination
            bind:currentPage={queuePage}
            totalPages={queueTotalPages}
            totalItems={queueTotalItems}
            limit={10}
            onPageChange={() => {
              document.querySelector('[data-claim-queue]')?.scrollIntoView({ behavior: 'smooth', block: 'start' });
            }}
          />
        </div>
      {/if}
    </div>

  </div>
</div>

<!-- Void Reason Modal Prompt -->
{#if showVoidModal && selectedClaimForVoid}
  <div class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-900/60 backdrop-blur-sm animate-fade-in">
    <div class="bg-white dark:bg-slate-950 w-full max-w-md rounded-2xl border border-slate-200 dark:border-slate-800 shadow-2xl overflow-hidden flex flex-col">
      <div class="flex justify-between items-center px-6 py-4 border-b border-slate-200 dark:border-slate-800 bg-slate-50 dark:bg-slate-900/30">
        <h3 class="font-bold text-sm text-slate-900 dark:text-white flex items-center gap-2">
          <AlertTriangle size={16} class="text-rose-600" />
          Tolak Klaim Garansi ({selectedClaimForVoid.ticket_id})
        </h3>
        <button onclick={() => showVoidModal = false} disabled={isSubmittingVoid} class="p-1.5 rounded-lg text-slate-400 hover:bg-slate-100 dark:hover:bg-slate-900 transition-colors">
          <X size={16} />
        </button>
      </div>

      <div class="p-6 space-y-4">
        <p class="text-xs text-slate-500 leading-relaxed">
          Klaim garansi ini akan ditolak (Void). Harap masukkan alasan keputusan penolakan ini secara rinci agar terdokumentasi dengan baik di sistem.
        </p>

        <div class="space-y-1.5">
          <label for="void-input-reason" class="text-xs font-bold text-slate-500">Alasan Void *</label>
          <input
            id="void-input-reason"
            type="text"
            bind:value={voidReason}
            required
            placeholder="Contoh: Segel rusak akibat dibongkar luar, bekas korosi air"
            class="w-full bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-xl py-2 px-3 text-xs focus:outline-none focus:border-rose-600 transition-colors"
          />
        </div>
      </div>

      <div class="px-6 py-4 border-t border-slate-200 dark:border-slate-800 bg-slate-50 dark:bg-slate-900/30 flex justify-end gap-2">
        <button onclick={() => showVoidModal = false} disabled={isSubmittingVoid} class="px-4 py-2 border border-slate-200 dark:border-slate-800 text-slate-700 dark:text-slate-300 font-bold text-xs uppercase tracking-wider rounded-lg hover:bg-slate-100 dark:hover:bg-slate-900 transition-colors">
          Batal
        </button>
        <button onclick={submitVoid} disabled={isSubmittingVoid || !voidReason.trim()} class="px-5 py-2 bg-rose-600 hover:bg-rose-700 disabled:bg-slate-100 disabled:text-slate-400 text-white font-bold text-xs uppercase tracking-wider rounded-lg transition-colors shadow-sm inline-flex items-center gap-1">
          {#if isSubmittingVoid}
            <LoaderCircle size={12} class="animate-spin" />
            Menolak...
          {:else}
            Void Garansi
          {/if}
        </button>
      </div>
    </div>
  </div>
{/if}

<!-- Custom Approve Confirmation Modal -->
{#if showApproveModal && selectedClaimForApprove}
  <div class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-900/60 backdrop-blur-sm animate-fade-in">
    <div class="bg-white dark:bg-slate-950 w-full max-w-lg rounded-3xl border border-slate-200 dark:border-slate-800 shadow-2xl overflow-hidden flex flex-col">
      <!-- Header -->
      <div class="flex justify-between items-center px-6 py-4 border-b border-slate-200 dark:border-slate-800 bg-slate-50 dark:bg-slate-900/30">
        <h3 class="font-extrabold text-sm text-slate-900 dark:text-white flex items-center gap-2">
          <ShieldCheck size={18} class="text-emerald-500 animate-pulse" />
          Setujui Klaim Garansi ({selectedClaimForApprove.ticket_id})
        </h3>
        <button onclick={() => showApproveModal = false} disabled={isSubmittingApprove} class="p-1.5 rounded-lg text-slate-400 hover:bg-slate-100 dark:hover:bg-slate-900 transition-colors">
          <X size={16} />
        </button>
      </div>

      <!-- Body -->
      <div class="p-6 space-y-4">
        <p class="text-xs text-slate-500 leading-relaxed">
          Periksa dan bandingkan data kerusakan sebelum menyetujui klaim garansi ini:
        </p>

        <!-- Side-by-Side Comparison -->
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <!-- Left Side: Original Ticket Details -->
          <div class="p-4 bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-2xl space-y-2">
            <span class="text-[9px] font-extrabold text-slate-400 uppercase tracking-wider block">Servis Awal</span>
            <div class="space-y-1 text-xs">
              <p class="font-bold text-slate-800 dark:text-slate-200">
                {selectedClaimForApprove.originalTicket?.brand || ''} {selectedClaimForApprove.originalTicket?.model || ''}
              </p>
              <p class="text-[11px] text-slate-500">
                <span class="font-medium text-slate-400">Kerusakan:</span> {selectedClaimForApprove.originalTicket?.issue || ''}
              </p>
              <p class="text-[11px] text-slate-500">
                <span class="font-medium text-slate-400">Biaya Awal:</span> {formatCurrency(selectedClaimForApprove.originalTicket?.price || 0)}
              </p>
            </div>
          </div>

          <!-- Right Side: New Claim Details -->
          <div class="p-4 bg-indigo-50/20 dark:bg-indigo-950/10 border border-indigo-100/30 dark:border-indigo-900/20 rounded-2xl space-y-2">
            <span class="text-[9px] font-extrabold text-indigo-600 dark:text-indigo-400 uppercase tracking-wider block">Klaim Garansi Baru</span>
            <div class="space-y-1 text-xs">
              <p class="font-bold text-indigo-700 dark:text-indigo-350">
                {selectedClaimForApprove.issue}
              </p>
              {#if selectedClaimForApprove.additional_description}
                <p class="text-[11px] text-slate-500">
                  <span class="font-medium text-slate-400">Catatan:</span> {selectedClaimForApprove.additional_description}
                </p>
              {/if}
            </div>
          </div>
        </div>

        <!-- Price & Free Ticket Info -->
        <div class="p-3 bg-emerald-50/50 dark:bg-emerald-950/10 border border-emerald-200/50 dark:border-emerald-900/20 rounded-2xl flex items-start gap-2.5 text-xs text-emerald-800 dark:text-emerald-350">
          <Info size={16} class="text-emerald-500 shrink-0 mt-0.5" />
          <p class="leading-relaxed">
            Dengan menyetujui klaim, sistem akan **menerbitkan tiket perbaikan baru senilai Rp 0** (Servis Garansi Gratis) dan statusnya otomatis diset ke antrean teknisi.
          </p>
        </div>
      </div>

      <!-- Footer -->
      <div class="px-6 py-4 border-t border-slate-200 dark:border-slate-800 bg-slate-50 dark:bg-slate-900/30 flex justify-end gap-2">
        <button onclick={() => showApproveModal = false} disabled={isSubmittingApprove} class="px-4 py-2 border border-slate-200 dark:border-slate-800 text-slate-700 dark:text-slate-300 font-bold text-xs uppercase tracking-wider rounded-xl hover:bg-slate-100 dark:hover:bg-slate-900 transition-colors cursor-pointer">
          Batal
        </button>
        <button onclick={approveClaim} disabled={isSubmittingApprove} class="px-5 py-2 bg-emerald-600 hover:bg-emerald-700 disabled:bg-slate-100 disabled:text-slate-400 text-white font-bold text-xs uppercase tracking-wider rounded-xl transition-colors shadow-sm inline-flex items-center gap-1.5 cursor-pointer">
          {#if isSubmittingApprove}
            <LoaderCircle size={14} class="animate-spin" />
            Menyetujui...
          {:else}
            Setujui Klaim
          {/if}
        </button>
      </div>
    </div>
  </div>
{/if}

<!-- Toast Stack -->
<div class="fixed top-6 right-6 z-50 flex flex-col gap-3 max-w-sm w-full pointer-events-none">
  {#each toasts as toast (toast.id)}
    <div 
      class="p-4 rounded-2xl shadow-xl border backdrop-blur-md flex justify-between items-center transition-all duration-300 pointer-events-auto animate-fade-in
        {toast.type === 'success' 
          ? 'bg-emerald-50/95 dark:bg-emerald-950/95 border-emerald-200 dark:border-emerald-900/50 text-emerald-800 dark:text-emerald-300' 
          : 'bg-rose-50/95 dark:bg-rose-950/95 border-rose-200 dark:border-rose-900/50 text-rose-800 dark:text-rose-300'}"
    >
      <div class="flex items-center gap-3">
        {#if toast.type === 'success'}
          <CheckCircle2 size={18} class="shrink-0" />
        {:else}
          <AlertOctagon size={18} class="shrink-0" />
        {/if}
        <p class="text-xs font-semibold leading-relaxed">{toast.message}</p>
      </div>
      <button 
        onclick={() => toasts = toasts.filter(t => t.id !== toast.id)} 
        class="p-1 rounded-lg hover:bg-black/5 dark:hover:bg-white/5 text-current transition-colors ml-3 shrink-0"
      >
        <X size={14} />
      </button>
    </div>
  {/each}
</div>
