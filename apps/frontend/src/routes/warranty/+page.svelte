<script lang="ts">
  import { onMount } from 'svelte';
  import { 
    Search, ShieldCheck, ShieldAlert, 
    ArrowLeft, Loader2, Info, AlertOctagon, CheckCircle2, X,
    FileText, User, Clock, AlertTriangle, Play, RefreshCw
  } from 'lucide-svelte';

  interface Ticket {
    id: string;
    customer_name: string;
    customer_gender: string;
    brand: string;
    model: string;
    issue: string;
    additional_description?: string;
    price: number;
    status: string;
    payment_status: string;
    warranty_days: number;
    entry_date: string;
    exit_date?: string;
    is_warranty?: boolean;
    parent_ticket_id?: string;
  }

  interface Claim {
    id: string;
    ticket_id: string;
    claim_ticket_id: string | null;
    issue: string;
    additional_description: string;
    status: 'waiting_inspection' | 'approved' | 'void';
    void_reason: string | null;
    inspected_at: string | null;
    created_at: string;
    // Client-side enriched fields
    originalTicket?: Ticket;
  }

  let searchQuery = $state('');
  let isLoading = $state(false);
  let searchResult = $state<Ticket | null>(null);
  let notFound = $state(false);
  let searchInput = $state<HTMLInputElement | null>(null);

  // Queue List State
  let claimsQueue = $state<Claim[]>([]);
  let isQueueLoading = $state(false);
  let pendingClaims = $derived(claimsQueue.filter(c => c.status === 'waiting_inspection'));

  // Intake Form State
  let intakeIssue = $state('');
  let intakeNote = $state('');
  let isSubmittingIntake = $state(false);

  // Void Modal State
  let showVoidModal = $state(false);
  let selectedClaimForVoid = $state<Claim | null>(null);
  let voidReason = $state('');
  let isSubmittingVoid = $state(false);

  // Notification Banners
  let successMessage = $state('');
  let errorMessage = $state('');

  onMount(() => {
    searchInput?.focus();
    fetchClaims();
  });

  async function fetchClaims() {
    isQueueLoading = true;
    try {
      // 1. Fetch claims
      const claimsRes = await fetch('/api/v1/warranty-claims');
      const claimsPayload = await claimsRes.json();
      
      // 2. Fetch all tickets to map details
      const ticketsRes = await fetch('/api/v1/tickets');
      const ticketsPayload = await ticketsRes.json();

      if (claimsPayload.success && ticketsPayload.success) {
        claimsQueue = claimsPayload.data.map((claim: Claim) => {
          const originalTicket = ticketsPayload.data.find((t: Ticket) => t.id === claim.ticket_id);
          return { ...claim, originalTicket };
        });
      }
    } catch (err) {
      console.error('Error fetching claims queue:', err);
    } finally {
      isQueueLoading = false;
    }
  }

  function checkWarrantyExpiry(exitDateStr: string, warrantyDays: number) {
    const exitDate = new Date(exitDateStr);
    const expiryDate = new Date(exitDate.getTime() + warrantyDays * 24 * 60 * 60 * 1000);
    const today = new Date();
    const remainingDays = Math.ceil((expiryDate.getTime() - today.getTime()) / (24 * 60 * 60 * 1000));
    
    return {
      expiryDate,
      isValid: remainingDays >= 0,
      remainingDays: remainingDays >= 0 ? remainingDays : 0,
      formattedExpiry: expiryDate.toLocaleDateString('id-ID', { day: 'numeric', month: 'long', year: 'numeric' })
    };
  }

  async function handleSearch(e: SubmitEvent) {
    e.preventDefault();
    if (!searchQuery.trim()) return;
    
    isLoading = true;
    notFound = false;
    searchResult = null;
    successMessage = '';
    errorMessage = '';

    try {
      const res = await fetch(`/api/v1/tickets`);
      const payload = await res.json();
      if (payload.success && payload.data) {
        const needle = searchQuery.trim().toLowerCase();
        const match = payload.data.find((t: Ticket) => 
          t.id.toLowerCase() === needle ||
          t.customer_name.toLowerCase().includes(needle) ||
          t.brand.toLowerCase().includes(needle) ||
          t.model.toLowerCase().includes(needle)
        );

        if (match) {
          searchResult = match;
          intakeIssue = '';
          intakeNote = '';
        } else {
          notFound = true;
        }
      }
    } catch (err) {
      console.error(err);
      notFound = true;
    } finally {
      isLoading = false;
    }
  }

  async function submitIntake() {
    if (!searchResult) return;
    if (!intakeIssue.trim()) return;

    isSubmittingIntake = true;
    errorMessage = '';
    successMessage = '';

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
      if (payload.success) {
        successMessage = `Berhasil mendaftarkan klaim garansi untuk nota ${searchResult.id}. Perangkat siap diletakkan di antrean inspeksi teknisi.`;
        searchResult = null;
        searchQuery = '';
        fetchClaims();
      } else {
        errorMessage = payload.error || 'Pendaftaran klaim gagal.';
      }
    } catch (err) {
      console.error(err);
      errorMessage = 'Koneksi ke mock API gagal.';
    } finally {
      isSubmittingIntake = false;
    }
  }

  async function approveClaim(claimId: string) {
    if (!confirm('Setujui klaim garansi ini dan buat tiket servis gratis (Rp 0)?')) return;
    
    successMessage = '';
    errorMessage = '';
    try {
      const res = await fetch(`/api/v1/warranty-claims/${claimId}/approve`, {
        method: 'POST'
      });
      const payload = await res.json();
      if (payload.success) {
        successMessage = `Klaim disetujui! Tiket baru ${payload.data.ticket.id} dengan harga Rp 0 telah dibuat.`;
        fetchClaims();
      } else {
        errorMessage = payload.error || 'Gagal menyetujui klaim.';
      }
    } catch (err) {
      console.error(err);
      errorMessage = 'Koneksi ke mock API gagal.';
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
    errorMessage = '';
    successMessage = '';
    try {
      const res = await fetch(`/api/v1/warranty-claims/${selectedClaimForVoid.id}/void`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          void_reason: voidReason
        })
      });

      const payload = await res.json();
      if (payload.success) {
        successMessage = `Klaim ditolak (Void). Tiket baru ${payload.data.ticket.id} dibuat dengan status dibatalkan beserta alasannya.`;
        showVoidModal = false;
        fetchClaims();
      } else {
        errorMessage = payload.error || 'Gagal membatalkan klaim.';
      }
    } catch (err) {
      console.error(err);
      errorMessage = 'Koneksi ke mock API gagal.';
    } finally {
      isSubmittingVoid = false;
    }
  }

  function formatCurrency(val: number) {
    return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(val);
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

    <!-- Alert Notifications -->
    {#if successMessage}
      <div class="p-4 bg-emerald-50 dark:bg-emerald-950/30 border border-emerald-200 dark:border-emerald-900/50 rounded-2xl flex justify-between items-center shadow-sm animate-fade-in">
        <div class="flex items-center gap-3 text-emerald-800 dark:text-emerald-300">
          <CheckCircle2 size={20} />
          <p class="text-sm font-semibold">{successMessage}</p>
        </div>
        <button onclick={() => successMessage = ''} class="text-emerald-500 hover:text-emerald-700">
          <X size={16} />
        </button>
      </div>
    {/if}

    {#if errorMessage}
      <div class="p-4 bg-rose-50 dark:bg-rose-950/30 border border-rose-200 dark:border-rose-900/50 rounded-2xl flex justify-between items-center shadow-sm animate-fade-in">
        <div class="flex items-center gap-3 text-rose-800 dark:text-rose-300">
          <AlertOctagon size={20} />
          <p class="text-sm font-semibold">{errorMessage}</p>
        </div>
        <button onclick={() => errorMessage = ''} class="text-rose-500 hover:text-rose-700">
          <X size={16} />
        </button>
      </div>
    {/if}

    <!-- GRID LAYOUT FOR INTAKE & DETAIL -->
    <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
      
      <!-- Left Column: Search & Verification Intake -->
      <div class="lg:col-span-2 space-y-6">
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
                <Loader2 size={14} class="animate-spin" />
                Mencari...
              {:else}
                Cari
              {/if}
            </button>
          </form>
        </div>

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
              <div>
                <span class="text-[10px] font-bold text-slate-400 uppercase block">Pelanggan</span>
                <span class="font-semibold text-slate-800 dark:text-slate-200">{searchResult.customer_name} ({searchResult.customer_gender})</span>
              </div>
              <div>
                <span class="text-[10px] font-bold text-slate-400 uppercase block">Kerusakan Awal</span>
                <span class="font-semibold text-slate-800 dark:text-slate-200">{searchResult.issue}</span>
              </div>
              <div>
                <span class="text-[10px] font-bold text-slate-400 uppercase block">Biaya Servis</span>
                <span class="font-semibold text-slate-800 dark:text-slate-200">{formatCurrency(searchResult.price)}</span>
              </div>
              <div>
                <span class="text-[10px] font-bold text-slate-400 uppercase block">Status Pembayaran</span>
                <span class="font-semibold text-slate-800 dark:text-slate-200 uppercase">{searchResult.payment_status}</span>
              </div>
            </div>

            <!-- Expiry Timeline info -->
            {#if searchResult.status === 'picked_up' && warranty}
              <div class="p-4 bg-slate-50 dark:bg-slate-900/50 rounded-xl space-y-2 border border-slate-100 dark:border-slate-850 text-xs">
                <div class="flex justify-between">
                  <span class="text-slate-400">Tanggal Ambil Perangkat:</span>
                  <span class="font-semibold">{searchResult.exit_date ? new Date(searchResult.exit_date).toLocaleDateString('id-ID', { day: 'numeric', month: 'long', year: 'numeric' }) : ''}</span>
                </div>
                <div class="flex justify-between">
                  <span class="text-slate-400">Durasi Garansi:</span>
                  <span class="font-semibold">{searchResult.warranty_days} Hari</span>
                </div>
                <div class="flex justify-between border-t border-slate-200/50 dark:border-slate-800/50 pt-2 font-bold">
                  <span class="text-slate-500">Masa Kedaluwarsa:</span>
                  <span class={warranty.isValid ? 'text-emerald-600 dark:text-emerald-400' : 'text-rose-600 dark:text-rose-400'}>
                    {warranty.formattedExpiry} ({warranty.isValid ? `Sisa ${warranty.remainingDays} Hari` : 'Habis'})
                  </span>
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
                        <Loader2 size={12} class="animate-spin" />
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

      <!-- Right Column: Operational Info -->
      <div class="bg-gradient-to-br from-indigo-500 to-purple-600 rounded-3xl p-6 text-white shadow-lg space-y-4 self-start">
        <h3 class="font-bold text-base flex items-center gap-2">
          <FileText size={18} />
          Informasi Kasir & Antrean
        </h3>
        <p class="text-xs leading-relaxed opacity-90">
          Untuk melayani pelanggan dengan cepat pada kondisi toko ramai:
        </p>
        <ul class="text-xs space-y-2 list-disc list-inside opacity-90">
          <li><strong>Langkah 1:</strong> Cari nota lama pelanggan & daftarkan kerusakan barunya.</li>
          <li><strong>Langkah 2:</strong> Berikan tanda terima intake klaim garansi ke pelanggan. Pelanggan dapat meninggalkan konter.</li>
          <li><strong>Langkah 3:</strong> Perangkat ditaruh di rak inspeksi teknisi.</li>
          <li><strong>Langkah 4:</strong> Teknisi/Admin memeriksa secara offline dan mengklik <strong>Setujui / Tolak</strong> di tabel antrean.</li>
        </ul>
      </div>

    </div>

    <!-- BOTTOM ROW: WORK LIST QUEUE (TABLE ANTREAN KLAIM) -->
    <div class="bg-white dark:bg-slate-950 border border-slate-200 dark:border-slate-800 rounded-2xl overflow-hidden shadow-sm">
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

      <div class="overflow-x-auto">
        <table class="w-full text-left border-collapse text-xs">
          <thead>
            <tr class="bg-slate-50 dark:bg-slate-900/50 border-b border-slate-200 dark:border-slate-800 text-[10px] font-bold text-slate-500 uppercase tracking-wider">
              <th class="py-3 px-5">ID Nota</th>
              <th class="py-3 px-5">Perangkat</th>
              <th class="py-3 px-5">Pelanggan</th>
              <th class="py-3 px-5">Keluhan Garansi</th>
              <th class="py-3 px-5">Tgl Daftar</th>
              <th class="py-3 px-5 text-right">Keputusan Inspeksi</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-100 dark:divide-slate-850">
            {#if pendingClaims.length === 0}
              <tr>
                <td colspan="6" class="py-12 text-center text-slate-400">
                  <CheckCircle2 size={24} class="mx-auto text-emerald-500 mb-2 opacity-80" />
                  <p class="font-semibold">Semua antrean inspeksi bersih</p>
                  <p class="text-[10px] mt-0.5">Tidak ada klaim garansi yang tertunda saat ini.</p>
                </td>
              </tr>
            {:else}
              {#each pendingClaims as claim (claim.id)}
                <tr class="hover:bg-slate-50/50 dark:hover:bg-slate-900/20 transition-colors">
                  <td class="py-3.5 px-5 font-mono font-semibold text-indigo-600 dark:text-indigo-400">
                    {claim.ticket_id}
                  </td>
                  <td class="py-3.5 px-5 font-semibold text-slate-900 dark:text-white">
                    {claim.originalTicket?.brand || ''} {claim.originalTicket?.model || ''}
                  </td>
                  <td class="py-3.5 px-5 text-slate-500">
                    {claim.originalTicket?.customer_name || ''}
                  </td>
                  <td class="py-3.5 px-5">
                    <span class="font-medium text-slate-800 dark:text-slate-200">{claim.issue}</span>
                    {#if claim.additional_description}
                      <span class="block text-[10px] text-slate-400 mt-0.5 truncate max-w-[200px]" title={claim.additional_description}>
                        {claim.additional_description}
                      </span>
                    {/if}
                  </td>
                  <td class="py-3.5 px-5 text-slate-400">
                    {new Date(claim.created_at).toLocaleString('id-ID', { day: 'numeric', month: 'short', hour: '2-digit', minute: '2-digit' })}
                  </td>
                  <td class="py-3.5 px-5 text-right space-x-1.5 whitespace-nowrap">
                    <button 
                      onclick={() => approveClaim(claim.id)} 
                      class="px-3 py-1.5 bg-emerald-600 hover:bg-emerald-700 text-white font-bold rounded-lg transition-colors inline-flex items-center gap-1 cursor-pointer active:scale-95"
                    >
                      <ShieldCheck size={12} /> Setujui
                    </button>
                    <button 
                      onclick={() => openVoidModal(claim)} 
                      class="px-3 py-1.5 bg-rose-600 hover:bg-rose-700 text-white font-bold rounded-lg transition-colors inline-flex items-center gap-1 cursor-pointer active:scale-95"
                    >
                      <X size={12} /> Tolak (Void)
                    </button>
                  </td>
                </tr>
              {/each}
            {/if}
          </tbody>
        </table>
      </div>
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
            <Loader2 size={12} class="animate-spin" />
            Menolak...
          {:else}
            Void Garansi
          {/if}
        </button>
      </div>
    </div>
  </div>
{/if}
