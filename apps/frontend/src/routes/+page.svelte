<script lang="ts">
  import { onMount } from "svelte";
  import { LoaderCircle, Search, Archive } from "lucide-svelte";
  import type { Ticket } from "$lib/types/ticket";
  import { checkSuccess, getErrorMessage } from "$lib/utils/api";
  
  import ControlBar from "./_components/ControlBar.svelte";
  import StatusFilters from "./_components/StatusFilters.svelte";
  import TicketTable from "./_components/TicketTable.svelte";
  import TicketCard from "./_components/TicketCard.svelte";
  import CreateModal from "./_components/CreateModal.svelte";
  import EditDrawer from "./_components/EditDrawer.svelte";
  import QRModal from "./_components/QRModal.svelte";
  import DeleteConfirmModal from "./_components/DeleteConfirmModal.svelte";
  import SkeletonTable from "./_components/SkeletonTable.svelte";
  import Pagination from "./_components/Pagination.svelte";

  // Svelte 5 Runes state
  let tickets = $state<Ticket[]>([]);
  let searchQuery = $state("");
  let localSearchInput = $state("");
  let statusFilter = $state("all");
  let page = $state(1);
  let totalPages = $state(1);
  let totalItems = $state(0);
  let isLoading = $state(true);
  let toastMessage = $state("");
  let statusCounts = $state<Record<string, number>>({});
  let hasMounted = $state(false);

  // Modals & Action loadings
  let showCreateModal = $state(false);
  let isCreating = $state(false);
  let selectedTicket = $state<Ticket | null>(null);
  let isUpdating = $state(false);
  let isDeleting = $state(false);
  let showQRModal = $state(false);
  let qrUrl = $state("");
  let isActionLoading = $state<Record<string, boolean>>({});
  let showDeleteModal = $state(false);

  // Reset page saat status filter berganti
  $effect(() => {
    const _ = statusFilter;
    page = 1;
  });

  // Logika Debouncing untuk Search Input (300ms setelah tidak mengetik)
  let debounceTimeout: ReturnType<typeof setTimeout>;
  $effect(() => {
    const input = localSearchInput;
    clearTimeout(debounceTimeout);
    debounceTimeout = setTimeout(() => {
      if (searchQuery !== input) {
        searchQuery = input;
        page = 1;
      }
    }, 300);
    return () => clearTimeout(debounceTimeout);
  });

  // Memicu fetch ulang secara otomatis ketika parameter utama berubah
  $effect(() => {
    if (!hasMounted) return;
    fetchTickets();
  });

  onMount(() => {
    hasMounted = true;
  });

  let activeAbortController: AbortController | null = null;

  async function fetchTickets() {
    if (activeAbortController) {
      activeAbortController.abort();
    }
    
    activeAbortController = new AbortController();
    const { signal } = activeAbortController;

    isLoading = true;
    try {
      const params = new URLSearchParams();
      if (statusFilter !== "all") params.set("status", statusFilter);
      if (searchQuery.trim()) params.set("search", searchQuery.trim());
      params.set("page", String(page));
      params.set("limit", "20");

      const res = await fetch(`/api/v1/tickets?${params}`, { signal });
      const json = await res.json();
      if (checkSuccess(res, json)) {
        const rawTickets = json.data || [];
        
        // Robustness fallback: if backend is unpaginated (does not return total count)
        if (json.total === undefined) {
          let filtered = rawTickets;
          // 1. Filter status
          if (statusFilter === "all") {
            filtered = filtered.filter((t: Ticket) => t.status !== "picked_up");
          } else {
            filtered = filtered.filter((t: Ticket) => t.status === statusFilter);
          }
          // 2. Filter search
          if (searchQuery.trim()) {
            const query = searchQuery.trim().toLowerCase();
            filtered = filtered.filter((t: Ticket) =>
              t.customer_name.toLowerCase().includes(query) ||
              t.model.toLowerCase().includes(query) ||
              t.issue.toLowerCase().includes(query) ||
              t.brand.toLowerCase().includes(query)
            );
          }
          
          tickets = filtered;
          totalPages = 1;
          totalItems = filtered.length;

          // Local status counts calculation
          const counts: Record<string, number> = {
            all: rawTickets.filter((t: Ticket) => t.status !== "picked_up").length,
            service_in: 0,
            on_process: 0,
            waiting_confirmation: 0,
            fixed: 0,
            picked_up: 0,
            cancelled: 0,
          };
          for (const t of rawTickets) {
            if (counts[t.status] !== undefined) {
              counts[t.status]++;
            }
          }
          statusCounts = counts;
        } else {
          // If server supports pagination and status filtration
          tickets = rawTickets;
          totalPages = json.total_pages || 1;
          totalItems = json.total ?? tickets.length;
          statusCounts = json.status_counts ?? {
            all: tickets.length,
            service_in: 0,
            on_process: 0,
            waiting_confirmation: 0,
            fixed: 0,
            picked_up: 0,
            cancelled: 0,
          };
        }

        if (page > totalPages) {
          page = totalPages;
          return;
        }
      } else {
        toastMessage = json.error || "Gagal memuat tiket.";
        setTimeout(() => toastMessage = "", 3000);
      }
    } catch (e: any) {
      if (e.name === 'AbortError') {
        return;
      }
      console.error("Error fetching tickets:", e);
      toastMessage = "Koneksi gagal. Silakan coba lagi.";
      setTimeout(() => toastMessage = "", 3000);
    } finally {
      if (!signal.aborted) {
        isLoading = false;
      }
    }
  }

  // Idempotency Key logic
  let createIdempotencyKey = $state("");

  function newIdempotencyKey() {
    return crypto.randomUUID();
  }

  function regenerateCreateIdempotencyKey() {
    createIdempotencyKey = newIdempotencyKey();
  }

  async function handleCreateTicket(formData: any) {
    if (isCreating) return;
    if (!createIdempotencyKey) regenerateCreateIdempotencyKey();
    isCreating = true;
    try {
      const res = await fetch("/api/v1/tickets", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          "X-Idempotency-Key": createIdempotencyKey,
        },
        body: JSON.stringify(formData),
      });
      const data = await res.json();
      if (checkSuccess(res, data)) {
        toastMessage = "Servis baru berhasil didaftarkan!";
        showCreateModal = false;
        await fetchTickets();
        setTimeout(() => (toastMessage = ""), 3000);
      } else {
        toastMessage = getErrorMessage(data, "Gagal mendaftarkan servis.");
        regenerateCreateIdempotencyKey();
        setTimeout(() => (toastMessage = ""), 3000);
      }
    } catch (e) {
      console.error("Error creating ticket:", e);
      toastMessage = "Terjadi kesalahan koneksi.";
      setTimeout(() => (toastMessage = ""), 3000);
    } finally {
      isCreating = false;
    }
  }

  async function handleUpdateTicket(formData: any) {
    if (!selectedTicket || isUpdating) return;
    isUpdating = true;
    try {
      const res = await fetch(`/api/v1/tickets/${selectedTicket.id}`, {
        method: "PATCH",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(formData),
      });
      const data = await res.json();
      if (checkSuccess(res, data)) {
        toastMessage = "Detail tiket berhasil diperbarui!";
        selectedTicket = null;
        await fetchTickets();
        setTimeout(() => (toastMessage = ""), 3000);
      } else {
        toastMessage = getErrorMessage(data, "Gagal memperbarui tiket.");
        setTimeout(() => (toastMessage = ""), 3000);
      }
    } catch (e) {
      console.error("Error updating ticket:", e);
      toastMessage = "Terjadi kesalahan koneksi.";
      setTimeout(() => (toastMessage = ""), 3000);
    } finally {
      isUpdating = false;
    }
  }

  function triggerDeleteConfirmation() {
    showDeleteModal = true;
  }

  async function executeDeleteTicket() {
    if (!selectedTicket || isDeleting || isUpdating) return;
    isDeleting = true;
    try {
      const res = await fetch(`/api/v1/tickets/${selectedTicket.id}`, {
        method: "DELETE",
      });
      const data = await res.json();
      if (checkSuccess(res, data)) {
        toastMessage = "Tiket berhasil dihapus!";
        showDeleteModal = false;
        selectedTicket = null;
        await fetchTickets();
        setTimeout(() => (toastMessage = ""), 3000);
      } else {
        toastMessage = getErrorMessage(data, "Gagal menghapus tiket.");
        setTimeout(() => (toastMessage = ""), 3000);
      }
    } catch (e) {
      console.error("Error deleting ticket:", e);
      toastMessage = "Terjadi kesalahan koneksi.";
      setTimeout(() => (toastMessage = ""), 3000);
    } finally {
      isDeleting = false;
    }
  }

  async function handleQuickStatusUpdate(ticketId: string, currentStatus: string) {
    if (isActionLoading[ticketId]) return;

    let nextStatus = "";
    if (currentStatus === "service_in") nextStatus = "on_process";
    else if (currentStatus === "on_process") nextStatus = "fixed";
    else if (currentStatus === "fixed") nextStatus = "picked_up";

    if (!nextStatus) return;
    isActionLoading[ticketId] = true;
    try {
      const res = await fetch(`/api/v1/tickets/${ticketId}`, {
        method: "PATCH",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ status: nextStatus }),
      });
      const data = await res.json();
      if (checkSuccess(res, data)) {
        toastMessage = `Status berhasil diperbarui ke: ${getStatusLabel(nextStatus, false)}`;
        await fetchTickets();
        setTimeout(() => (toastMessage = ""), 3000);
      } else {
        toastMessage = getErrorMessage(data, "Gagal memperbarui status.");
        setTimeout(() => (toastMessage = ""), 3000);
      }
    } catch (e) {
      console.error("Error quick updating status:", e);
    } finally {
      isActionLoading[ticketId] = false;
    }
  }

  async function submitIssue(issueType: string, newPrice: number, newDiagnosis: string) {
    if (!selectedTicket || isUpdating) return;
    isUpdating = true;
    try {
      let payload: any = {};
      if (issueType === "unrepairable") {
        payload = {
          status: "cancelled",
        };
      } else {
        payload = {
          status: "waiting_confirmation",
          price: newPrice,
          additional_description: newDiagnosis,
        };
      }

      const res = await fetch(`/api/v1/tickets/${selectedTicket.id}`, {
        method: "PATCH",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(payload),
      });
      const data = await res.json();
      if (checkSuccess(res, data)) {
        toastMessage = "Kendala teknisi berhasil dilaporkan!";
        selectedTicket = null;
        await fetchTickets();
        setTimeout(() => (toastMessage = ""), 3000);
      } else {
        toastMessage = getErrorMessage(data, "Gagal mengirim laporan kendala.");
        setTimeout(() => (toastMessage = ""), 3000);
      }
    } catch (e) {
      console.error("Error reporting issue:", e);
      toastMessage = "Terjadi kesalahan koneksi.";
      setTimeout(() => (toastMessage = ""), 3000);
    } finally {
      isUpdating = false;
    }
  }

  async function approveByCustomer() {
    if (!selectedTicket || isUpdating) return;
    isUpdating = true;
    try {
      const res = await fetch(`/api/v1/tickets/${selectedTicket.id}`, {
        method: "PATCH",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ status: "on_process" }),
      });
      const data = await res.json();
      if (checkSuccess(res, data)) {
        toastMessage = "Persetujuan pelanggan berhasil disimpan!";
        selectedTicket = null;
        await fetchTickets();
        setTimeout(() => (toastMessage = ""), 3000);
      } else {
        toastMessage = getErrorMessage(data, "Gagal menyimpan persetujuan.");
        setTimeout(() => (toastMessage = ""), 3000);
      }
    } catch (e) {
      console.error("Error approving:", e);
      toastMessage = "Terjadi kesalahan koneksi.";
      setTimeout(() => (toastMessage = ""), 3000);
    } finally {
      isUpdating = false;
    }
  }

  async function rejectByCustomer() {
    if (!selectedTicket || isUpdating) return;
    isUpdating = true;
    try {
      const res = await fetch(`/api/v1/tickets/${selectedTicket.id}`, {
        method: "PATCH",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ status: "cancelled" }),
      });
      const data = await res.json();
      if (checkSuccess(res, data)) {
        toastMessage = "Penolakan pelanggan berhasil disimpan (servis batal)!";
        selectedTicket = null;
        await fetchTickets();
        setTimeout(() => (toastMessage = ""), 3000);
      } else {
        toastMessage = getErrorMessage(data, "Gagal menyimpan penolakan.");
        setTimeout(() => (toastMessage = ""), 3000);
      }
    } catch (e) {
      console.error("Error rejecting:", e);
      toastMessage = "Terjadi kesalahan koneksi.";
      setTimeout(() => (toastMessage = ""), 3000);
    } finally {
      isUpdating = false;
    }
  }

  function getStatusBadgeClass(status: string, isWarranty: boolean) {
    if (isWarranty) {
      return "bg-purple-50 text-purple-700 border-purple-200 dark:bg-purple-950/30 dark:text-purple-400 dark:border-purple-900";
    }
    switch (status) {
      case "service_in":
        return "bg-blue-50 text-blue-700 border-blue-200 dark:bg-blue-950/30 dark:text-blue-400 dark:border-blue-900";
      case "on_process":
        return "bg-amber-50 text-amber-700 border-amber-200 dark:bg-amber-950/30 dark:text-amber-400 dark:border-amber-900";
      case "waiting_confirmation":
        return "bg-rose-50 text-rose-700 border-rose-200 dark:bg-rose-950/30 dark:text-rose-400 dark:border-rose-900";
      case "cancelled":
        return "bg-slate-100 text-slate-600 border-slate-300 dark:bg-slate-900 dark:text-slate-400 dark:border-slate-800";
      case "fixed":
        return "bg-emerald-50 text-emerald-700 border-emerald-200 dark:bg-emerald-950/30 dark:text-emerald-400 dark:border-emerald-900";
      case "picked_up":
        return "bg-indigo-50 text-indigo-700 border-indigo-200 dark:bg-indigo-950/30 dark:text-indigo-400 dark:border-indigo-900";
      default:
        return "bg-slate-50 text-slate-700 border-slate-200";
    }
  }

  function getStatusLabel(status: string, isWarranty: boolean) {
    if (isWarranty) return "Klaim Garansi";
    switch (status) {
      case "service_in":
        return "Masuk";
      case "on_process":
        return "Diproses";
      case "waiting_confirmation":
        return "Kendala (Menunggu)";
      case "cancelled":
        return "Batal";
      case "fixed":
        return "Selesai";
      case "picked_up":
        return "Sudah Diambil";
      default:
        return status;
    }
  }

  function copyTrackingLink(id: string) {
    const link = `${window.location.origin}/track/${id}`;
    navigator.clipboard.writeText(link);
    toastMessage = "Link pelacakan berhasil disalin!";
    setTimeout(() => (toastMessage = ""), 3000);
  }

  function openQRModal(id: string) {
    qrUrl = `${window.location.origin}/track/${id}`;
    showQRModal = true;
  }

  function copyQRUrl() {
    navigator.clipboard.writeText(qrUrl);
    toastMessage = "URL pelacakan berhasil disalin!";
    setTimeout(() => (toastMessage = ""), 3000);
  }
</script>

<div class="container mx-auto px-4 py-8 max-w-7xl animate-fade-in font-sans">
  <!-- Control Bar -->
  <ControlBar
    bind:searchQuery={localSearchInput}
    onAddTicket={() => {
      regenerateCreateIdempotencyKey();
      showCreateModal = true;
    }}
  />

  <!-- Table Header Controls & Special Archive Button -->
  <div class="flex flex-col lg:flex-row lg:items-center justify-between gap-4 mb-4 mt-2 px-1">
    <!-- Status Filters (menempel pada tabel) -->
    <div class="flex-grow min-w-0">
      <StatusFilters
        bind:selectedStatus={statusFilter}
        {statusCounts}
      />
    </div>

    <!-- Special Archive Button & Status Label -->
    <div class="flex items-center gap-3 shrink-0 self-end lg:self-auto">
      <div class="text-xs font-semibold uppercase tracking-wider text-slate-400 dark:text-slate-500">
        {#if statusFilter === "picked_up"}
          <span class="text-indigo-600 dark:text-indigo-400 font-bold">Arsip Tiket (Sudah Diambil)</span>
        {:else}
          <span>Antrean Perbaikan Aktif</span>
        {/if}
      </div>

      <button
        onclick={() => {
          if (statusFilter === "picked_up") {
            statusFilter = "all";
          } else {
            statusFilter = "picked_up";
          }
        }}
        class="inline-flex items-center gap-2 px-4 py-2 text-xs font-bold border rounded-xl transition-all duration-200 cursor-pointer select-none active:scale-95
          {statusFilter === 'picked_up'
            ? 'bg-indigo-600 border-indigo-600 text-white shadow-sm ring-2 ring-offset-2 ring-indigo-500/20 dark:ring-offset-slate-950'
            : 'bg-white dark:bg-slate-950 border-slate-200 dark:border-slate-800 text-slate-700 dark:text-slate-300 hover:bg-slate-50 dark:hover:bg-slate-900 shadow-sm'}"
      >
        <Archive size={14} class={statusFilter === 'picked_up' ? 'text-white' : 'text-indigo-600'} />
        <span>Lihat Sudah Diambil</span>
      </button>
    </div>
  </div>

  <!-- Tickets Content -->
  <div class="relative min-h-[300px]">
    {#if isLoading}
      <!-- Viewport-Fixed Premium Circular Spinner Overlay -->
      <div class="fixed inset-0 bg-slate-900/10 dark:bg-slate-950/20 backdrop-blur-[0.5px] flex items-center justify-center z-50 transition-all duration-300">
        <div class="flex flex-col items-center gap-3 bg-white dark:bg-slate-900 px-6 py-4 rounded-2xl shadow-xl border border-slate-200/80 dark:border-slate-800/80 backdrop-blur-md">
          <LoaderCircle class="animate-spin text-blue-600 dark:text-blue-400" size={32} />
          <span class="text-xs font-semibold text-slate-600 dark:text-slate-300 uppercase tracking-wider animate-pulse">Memuat Data...</span>
        </div>
      </div>
    {/if}

    {#if tickets.length === 0 && !isLoading}
      <div class="bg-white dark:bg-slate-950 border border-slate-200 dark:border-slate-800 rounded-2xl p-16 text-center shadow-sm">
        <div class="w-12 h-12 bg-slate-50 dark:bg-slate-900 rounded-full flex items-center justify-center mx-auto mb-4 border border-slate-100 dark:border-slate-800 text-slate-400">
          <Search size={22} />
        </div>
        <h3 class="text-lg font-bold text-slate-900 dark:text-white mb-1">Tidak ada tiket ditemukan</h3>
        <p class="text-sm text-slate-500 dark:text-slate-400 max-w-sm mx-auto">
          Try modifying your search or filters, or add a new repair ticket to get started.
        </p>
      </div>
    {:else if tickets.length > 0}
      <!-- Table Layout (Desktop) -->
      <TicketTable
        tickets={tickets}
        {isActionLoading}
        onSelectTicket={(ticket) => (selectedTicket = ticket)}
        onQuickAction={handleQuickStatusUpdate}
        {getStatusBadgeClass}
        {getStatusLabel}
      />

      <!-- Mobile Cards List (Mobile) -->
      <div class="grid grid-cols-1 gap-4 md:hidden">
        {#each tickets as ticket (ticket.id)}
          <TicketCard
            {ticket}
            isActionLoading={Boolean(isActionLoading[ticket.id])}
            onSelectTicket={() => (selectedTicket = ticket)}
            onQuickAction={() => handleQuickStatusUpdate(ticket.id, ticket.status)}
            {getStatusBadgeClass}
            {getStatusLabel}
          />
        {/each}
      </div>

      {#if tickets.length > 0}
        <Pagination
          bind:currentPage={page}
          {totalPages}
          {totalItems}
          limit={20}
          onPageChange={() => {
            window.scrollTo({ top: 0, behavior: 'smooth' });
          }}
        />
      {/if}
    {/if}
  </div>
</div>

<!-- Create Modal -->
<CreateModal
  isOpen={showCreateModal}
  {isCreating}
  onSubmit={handleCreateTicket}
  onClose={() => (showCreateModal = false)}
/>

<!-- Edit Drawer -->
<EditDrawer
  ticket={selectedTicket}
  {isUpdating}
  {isDeleting}
  onClose={() => (selectedTicket = null)}
  onSubmitUpdate={handleUpdateTicket}
  onDelete={triggerDeleteConfirmation}
  {copyTrackingLink}
  {openQRModal}
  {submitIssue}
  {approveByCustomer}
  {rejectByCustomer}
  {getStatusBadgeClass}
  {getStatusLabel}
/>

<!-- QR Modal -->
<QRModal
  isOpen={showQRModal}
  {qrUrl}
  onClose={() => (showQRModal = false)}
  onCopyQRUrl={copyQRUrl}
/>

<!-- Delete Confirm Modal -->
<DeleteConfirmModal
  isOpen={showDeleteModal}
  ticket={selectedTicket}
  {isDeleting}
  onConfirm={executeDeleteTicket}
  onClose={() => (showDeleteModal = false)}
/>

{#if toastMessage}
  <div class="fixed bottom-6 right-6 z-50 animate-slide-in bg-slate-900 text-white px-5 py-4 rounded-xl shadow-2xl flex items-center gap-3 border border-slate-700">
    <span class="text-sm font-semibold">{toastMessage}</span>
  </div>
{/if}

<style>
  @keyframes fade-in {
    from { opacity: 0; }
    to { opacity: 1; }
  }
  @keyframes slide-in {
    from { transform: translateX(100%); }
    to { transform: translateX(0); }
  }
  .animate-fade-in {
    animation: fade-in 0.25s ease-out forwards;
  }
  .animate-slide-in {
    animation: slide-in 0.3s cubic-bezier(0.16, 1, 0.3, 1) forwards;
  }
</style>
