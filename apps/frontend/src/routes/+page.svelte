<script lang="ts">
  import { onMount } from "svelte";
  import { LoaderCircle, Search } from "lucide-svelte";
  import type { Ticket } from "$lib/types/ticket";
  
  import ControlBar from "./_components/ControlBar.svelte";
  import TicketTable from "./_components/TicketTable.svelte";
  import TicketCard from "./_components/TicketCard.svelte";
  import CreateModal from "./_components/CreateModal.svelte";
  import EditDrawer from "./_components/EditDrawer.svelte";
  import QRModal from "./_components/QRModal.svelte";
  import DeleteConfirmModal from "./_components/DeleteConfirmModal.svelte";
  import SkeletonTable from "./_components/SkeletonTable.svelte";

  // Svelte 5 Runes state
  let tickets = $state<Ticket[]>([]);
  let searchQuery = $state("");
  let statusFilter = $state("all");
  let isLoading = $state(true);
  let toastMessage = $state("");

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

  // Status counts computed properties
  let statusCounts = $derived.by(() => {
    const counts: Record<string, number> = {
      all: tickets.length,
      service_in: 0,
      on_process: 0,
      waiting_confirmation: 0,
      fixed: 0,
      picked_up: 0,
      cancelled: 0,
    };
    for (const t of tickets) {
      if (counts[t.status] !== undefined) {
        counts[t.status]++;
      }
    }
    return counts;
  });

  let filteredTickets = $derived(
    tickets.filter((t) => {
      const matchStatus = statusFilter === "all" || t.status === statusFilter;
      const matchSearch =
        searchQuery === "" ||
        t.customer_name.toLowerCase().includes(searchQuery.toLowerCase()) ||
        t.model.toLowerCase().includes(searchQuery.toLowerCase()) ||
        t.issue.toLowerCase().includes(searchQuery.toLowerCase()) ||
        t.brand.toLowerCase().includes(searchQuery.toLowerCase());
      return matchStatus && matchSearch;
    })
  );

  // Load tickets on mount
  onMount(async () => {
    await fetchTickets();
  });

  async function fetchTickets() {
    isLoading = true;
    try {
      const res = await fetch("/api/v1/tickets");
      const data = await res.json();
      if (data.success) {
        tickets = data.data || [];
      } else {
        toastMessage = data.error || "Gagal memuat tiket.";
        setTimeout(() => toastMessage = "", 3000);
      }
    } catch (e) {
      console.error("Error fetching tickets:", e);
      toastMessage = "Koneksi gagal. Silakan coba lagi.";
      setTimeout(() => toastMessage = "", 3000);
    } finally {
      isLoading = false;
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
      if (data.success) {
        toastMessage = "Servis baru berhasil didaftarkan!";
        showCreateModal = false;
        await fetchTickets();
        setTimeout(() => (toastMessage = ""), 3000);
      } else {
        toastMessage = data.error || "Gagal mendaftarkan servis.";
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
      if (data.success) {
        toastMessage = "Detail tiket berhasil diperbarui!";
        selectedTicket = null;
        await fetchTickets();
        setTimeout(() => (toastMessage = ""), 3000);
      } else {
        toastMessage = data.error || "Gagal memperbarui tiket.";
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
      if (data.success) {
        toastMessage = "Tiket berhasil dihapus!";
        showDeleteModal = false;
        selectedTicket = null;
        await fetchTickets();
        setTimeout(() => (toastMessage = ""), 3000);
      } else {
        toastMessage = data.error || "Gagal menghapus tiket.";
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
      if (data.success) {
        toastMessage = `Status berhasil diperbarui ke: ${getStatusLabel(nextStatus, false)}`;
        await fetchTickets();
        setTimeout(() => (toastMessage = ""), 3000);
      } else {
        toastMessage = data.error || "Gagal memperbarui status.";
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
      if (data.success) {
        toastMessage = "Kendala teknisi berhasil dilaporkan!";
        selectedTicket = null;
        await fetchTickets();
        setTimeout(() => (toastMessage = ""), 3000);
      } else {
        toastMessage = data.error || "Gagal mengirim laporan kendala.";
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
      if (data.success) {
        toastMessage = "Persetujuan pelanggan berhasil disimpan!";
        selectedTicket = null;
        await fetchTickets();
        setTimeout(() => (toastMessage = ""), 3000);
      } else {
        toastMessage = data.error || "Gagal menyimpan persetujuan.";
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
      if (data.success) {
        toastMessage = "Penolakan pelanggan berhasil disimpan (servis batal)!";
        selectedTicket = null;
        await fetchTickets();
        setTimeout(() => (toastMessage = ""), 3000);
      } else {
        toastMessage = data.error || "Gagal menyimpan penolakan.";
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
    bind:searchQuery
    bind:statusFilter
    {statusCounts}
    onAddTicket={() => {
      regenerateCreateIdempotencyKey();
      showCreateModal = true;
    }}
  />

  <!-- Tickets Content -->
  {#if isLoading}
    <SkeletonTable />
  {:else if filteredTickets.length === 0}
    <div class="bg-white dark:bg-slate-950 border border-slate-200 dark:border-slate-800 rounded-2xl p-16 text-center shadow-sm">
      <div class="w-12 h-12 bg-slate-50 dark:bg-slate-900 rounded-full flex items-center justify-center mx-auto mb-4 border border-slate-100 dark:border-slate-800 text-slate-400">
        <Search size={22} />
      </div>
      <h3 class="text-lg font-bold text-slate-900 dark:text-white mb-1">Tidak ada tiket ditemukan</h3>
      <p class="text-sm text-slate-500 dark:text-slate-400 max-w-sm mx-auto">
        Try modifying your search or filters, or add a new repair ticket to get started.
      </p>
    </div>
  {:else}
    <!-- Table Layout (Desktop) -->
    <TicketTable
      tickets={filteredTickets}
      {isActionLoading}
      onSelectTicket={(ticket) => (selectedTicket = ticket)}
      onQuickAction={handleQuickStatusUpdate}
      {getStatusBadgeClass}
      {getStatusLabel}
    />

    <!-- Mobile Cards List (Mobile) -->
    <div class="grid grid-cols-1 gap-4 md:hidden">
      {#each filteredTickets as ticket (ticket.id)}
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
  {/if}
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
