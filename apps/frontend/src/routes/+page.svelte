<script lang="ts">
  import { onMount } from 'svelte';
  import {
    Search,
    Plus,
    Wrench,
    TrendingUp,
    DollarSign,
    Clock,
    CheckCircle2,
    Trash2,
    X,
    Smartphone,
    User,
    Sparkles,
    Calendar,
    FileText,
    ArrowRight,
    Loader2
  } from 'lucide-svelte';

  interface Ticket {
    id: string;
    customer_name: string;
    customer_gender: string;
    brand: string;
    model: string;
    issue: string;
    additional_description?: string;
    accessories?: string;
    price: number;
    status: string;
    payment_status: string;
    warranty_days: number;
    entry_date: string;
    exit_date?: string;
    warranty_expiry_date?: string;
  }

  // Svelte 5 Runes state
  let tickets = $state<Ticket[]>([]);
  let searchQuery = $state('');
  let statusFilter = $state('all');
  let isLoading = $state(true);
  let isActionLoading = $state<Record<string, boolean>>({});

  // Modal & Drawer State
  let showCreateModal = $state(false);
  let selectedTicket = $state<Ticket | null>(null);
  let editFormElement = $state<HTMLFormElement | null>(null);

  // Form States
  let createForm = $state({
    customer_name: '',
    customer_gender: 'Male',
    brand: '',
    model: '',
    issue: '',
    additional_description: '',
    accessories: '',
    price: 0,
    warranty_days: 30
  });

  let editForm = $state({
    customer_name: '',
    customer_gender: 'Male',
    brand: '',
    model: '',
    issue: '',
    additional_description: '',
    accessories: '',
    price: 0,
    status: '',
    payment_status: '',
    warranty_days: 30
  });

  // Load tickets on mount
  onMount(async () => {
    await fetchTickets();
  });

  async function fetchTickets() {
    isLoading = true;
    try {
      const res = await fetch('/api/v1/tickets');
      const data = await res.json();
      if (data.success) {
        tickets = data.data || [];
      }
    } catch (e) {
      console.error('Error fetching tickets:', e);
    } finally {
      isLoading = false;
    }
  }

  async function handleCreateTicket(e: SubmitEvent) {
    e.preventDefault();
    try {
      const res = await fetch('/api/v1/tickets', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(createForm)
      });
      const data = await res.json();
      if (data.success) {
        showCreateModal = false;
        // Reset form
        createForm = {
          customer_name: '',
          customer_gender: 'Male',
          brand: '',
          model: '',
          issue: '',
          additional_description: '',
          accessories: '',
          price: 0,
          warranty_days: 30
        };
        await fetchTickets();
      } else {
        alert('Failed to create ticket: ' + (data.error || 'Unknown error'));
      }
    } catch (e) {
      console.error(e);
    }
  }

  async function handleQuickStatusUpdate(ticketId: string, currentStatus: string) {
    let nextStatus = '';
    if (currentStatus === 'service_in') nextStatus = 'on_process';
    else if (currentStatus === 'on_process') nextStatus = 'fixed';
    else if (currentStatus === 'fixed') nextStatus = 'picked_up';
    else return;

    isActionLoading[ticketId] = true;
    try {
      const res = await fetch(`/api/v1/tickets/${ticketId}`, {
        method: 'PATCH',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ status: nextStatus })
      });
      const data = await res.json();
      if (data.success) {
        await fetchTickets();
      }
    } catch (e) {
      console.error(e);
    } finally {
      isActionLoading[ticketId] = false;
    }
  }

  function openEditDrawer(ticket: Ticket) {
    selectedTicket = ticket;
    editForm = {
      customer_name: ticket.customer_name,
      customer_gender: ticket.customer_gender,
      brand: ticket.brand,
      model: ticket.model,
      issue: ticket.issue,
      additional_description: ticket.additional_description || '',
      accessories: ticket.accessories || '',
      price: Number(ticket.price),
      status: ticket.status,
      payment_status: ticket.payment_status,
      warranty_days: ticket.warranty_days
    };
  }

  function closeEditDrawer() {
    selectedTicket = null;
  }

  function submitEditForm() {
    editFormElement?.requestSubmit();
  }

  async function deleteSelectedTicket() {
    if (!selectedTicket) return;
    await handleDeleteTicket(selectedTicket.id);
  }

  async function handleUpdateTicket(e: SubmitEvent) {
    e.preventDefault();
    if (!selectedTicket) return;

    try {
      const res = await fetch(`/api/v1/tickets/${selectedTicket.id}`, {
        method: 'PATCH',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(editForm)
      });
      const data = await res.json();
      if (data.success) {
        closeEditDrawer();
        await fetchTickets();
      } else {
        alert('Failed to update ticket: ' + (data.error || 'Unknown error'));
      }
    } catch (e) {
      console.error(e);
    }
  }

  async function handleDeleteTicket(ticketId: string) {
    if (!confirm('Are you sure you want to delete this ticket?')) return;

    try {
      const res = await fetch(`/api/v1/tickets/${ticketId}`, {
        method: 'DELETE'
      });
      const data = await res.json();
      if (data.success) {
        closeEditDrawer();
        await fetchTickets();
      } else {
        alert('Failed to delete ticket: ' + (data.error || 'Unknown error'));
      }
    } catch (e) {
      console.error(e);
    }
  }

  // Derived statistics (Runes)
  let serviceInCount = $derived(tickets.filter(t => t.status === 'service_in').length);
  let onProcessCount = $derived(tickets.filter(t => t.status === 'on_process').length);
  let fixedCount = $derived(tickets.filter(t => t.status === 'fixed').length);
  let todaysRevenue = $derived(
    tickets
      .filter(t => {
        if (t.status !== 'picked_up' || !t.exit_date) return false;
        const exitDate = new Date(t.exit_date);
        const today = new Date();
        return (
          exitDate.getDate() === today.getDate() &&
          exitDate.getMonth() === today.getMonth() &&
          exitDate.getFullYear() === today.getFullYear()
        );
      })
      .reduce((acc, t) => acc + Number(t.price), 0)
  );

  // Filtered Tickets
  let filteredTickets = $derived(
    tickets.filter(t => {
      const matchesSearch =
        t.customer_name.toLowerCase().includes(searchQuery.toLowerCase()) ||
        t.brand.toLowerCase().includes(searchQuery.toLowerCase()) ||
        t.model.toLowerCase().includes(searchQuery.toLowerCase()) ||
        t.issue.toLowerCase().includes(searchQuery.toLowerCase());
      const matchesStatus = statusFilter === 'all' || t.status === statusFilter;
      return matchesSearch && matchesStatus;
    })
  );

  function getStatusLabel(status: string) {
    switch (status) {
      case 'service_in': return 'Service In';
      case 'on_process': return 'On Process';
      case 'fixed': return 'Fixed';
      case 'picked_up': return 'Picked Up';
      default: return status;
    }
  }

  function getStatusBadgeClass(status: string) {
    switch (status) {
      case 'service_in': return 'bg-blue-50 text-blue-700 border-blue-200 dark:bg-blue-900/30 dark:text-blue-300 dark:border-blue-800';
      case 'on_process': return 'bg-amber-50 text-amber-700 border-amber-200 dark:bg-amber-900/30 dark:text-amber-300 dark:border-amber-800';
      case 'fixed': return 'bg-emerald-50 text-emerald-700 border-emerald-200 dark:bg-emerald-900/30 dark:text-emerald-300 dark:border-emerald-800';
      case 'picked_up': return 'bg-slate-100 text-slate-700 border-slate-200 dark:bg-slate-800/80 dark:text-slate-300 dark:border-slate-700';
      default: return 'bg-slate-50 text-slate-600 border-slate-200';
    }
  }

  function formatCurrency(amount: number) {
    return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(amount);
  }

  function formatDate(dateStr: string) {
    if (!dateStr) return '-';
    return new Date(dateStr).toLocaleDateString('id-ID', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    });
  }
</script>

<div class="container mx-auto px-4 py-8 max-w-7xl animate-fade-in">
  <!-- Stats Section -->
  <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
    <div class="bg-white dark:bg-slate-950 p-6 rounded-2xl border border-slate-200 dark:border-slate-800 shadow-sm relative overflow-hidden transition-all hover:shadow-md">
      <div class="flex justify-between items-start mb-4">
        <div>
          <p class="text-xs font-semibold text-slate-500 uppercase tracking-wider">In for Service</p>
          <h3 class="text-3xl font-bold mt-1 text-slate-900 dark:text-white">{serviceInCount}</h3>
        </div>
        <div class="p-3 bg-blue-50 dark:bg-blue-900/20 text-blue-600 dark:text-blue-400 rounded-xl">
          <Clock size={20} />
        </div>
      </div>
      <div class="h-1 bg-blue-600 w-full absolute bottom-0 left-0 rounded-b-2xl"></div>
    </div>

    <div class="bg-white dark:bg-slate-950 p-6 rounded-2xl border border-slate-200 dark:border-slate-800 shadow-sm relative overflow-hidden transition-all hover:shadow-md">
      <div class="flex justify-between items-start mb-4">
        <div>
          <p class="text-xs font-semibold text-slate-500 uppercase tracking-wider">On Process</p>
          <h3 class="text-3xl font-bold mt-1 text-slate-900 dark:text-white">{onProcessCount}</h3>
        </div>
        <div class="p-3 bg-amber-50 dark:bg-amber-900/20 text-amber-600 dark:text-amber-400 rounded-xl">
          <Wrench size={20} />
        </div>
      </div>
      <div class="h-1 bg-amber-500 w-full absolute bottom-0 left-0 rounded-b-2xl"></div>
    </div>

    <div class="bg-white dark:bg-slate-950 p-6 rounded-2xl border border-slate-200 dark:border-slate-800 shadow-sm relative overflow-hidden transition-all hover:shadow-md">
      <div class="flex justify-between items-start mb-4">
        <div>
          <p class="text-xs font-semibold text-slate-500 uppercase tracking-wider">Fixed & Ready</p>
          <h3 class="text-3xl font-bold mt-1 text-slate-900 dark:text-white">{fixedCount}</h3>
        </div>
        <div class="p-3 bg-emerald-50 dark:bg-emerald-900/20 text-emerald-600 dark:text-emerald-400 rounded-xl">
          <CheckCircle2 size={20} />
        </div>
      </div>
      <div class="h-1 bg-emerald-500 w-full absolute bottom-0 left-0 rounded-b-2xl"></div>
    </div>

    <div class="bg-white dark:bg-slate-950 p-6 rounded-2xl border border-slate-200 dark:border-slate-800 shadow-sm relative overflow-hidden transition-all hover:shadow-md">
      <div class="flex justify-between items-start mb-4">
        <div>
          <p class="text-xs font-semibold text-slate-500 uppercase tracking-wider">Today's Revenue</p>
          <h3 class="text-2xl font-bold mt-1.5 text-slate-900 dark:text-white">{formatCurrency(todaysRevenue)}</h3>
        </div>
        <div class="p-3 bg-indigo-50 dark:bg-indigo-900/20 text-indigo-600 dark:text-indigo-400 rounded-xl">
          <DollarSign size={20} />
        </div>
      </div>
      <div class="h-1 bg-indigo-500 w-full absolute bottom-0 left-0 rounded-b-2xl"></div>
    </div>
  </div>

  <!-- Control Bar -->
  <div class="flex flex-col md:flex-row justify-between items-center gap-4 mb-6">
    <div class="flex flex-col sm:flex-row items-center w-full md:w-auto gap-4">
      <!-- Search Input -->
      <div class="relative w-full sm:w-80">
        <span class="absolute inset-y-0 left-0 flex items-center pl-3.5 pointer-events-none text-slate-400">
          <Search size={18} />
        </span>
        <input
          type="text"
          bind:value={searchQuery}
          placeholder="Search by customer, model, issue..."
          class="w-full bg-white dark:bg-slate-950 border border-slate-200 dark:border-slate-800 rounded-xl py-2.5 pl-10 pr-4 text-sm text-slate-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-blue-600/20 focus:border-blue-600 transition-all placeholder-slate-400"
        />
      </div>

      <!-- Status Filter -->
      <select
        bind:value={statusFilter}
        class="w-full sm:w-44 bg-white dark:bg-slate-950 border border-slate-200 dark:border-slate-800 rounded-xl py-2.5 px-3.5 text-sm text-slate-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-blue-600/20 focus:border-blue-600 transition-all cursor-pointer"
      >
        <option value="all">All Statuses</option>
        <option value="service_in">Service In</option>
        <option value="on_process">On Process</option>
        <option value="fixed">Fixed</option>
        <option value="picked_up">Picked Up</option>
      </select>
    </div>

    <!-- Create Button -->
    <button
      onclick={() => showCreateModal = true}
      class="w-full md:w-auto px-6 py-2.5 bg-blue-600 hover:bg-blue-700 text-white font-bold text-sm rounded-xl transition-all shadow-sm hover:shadow active:scale-95 inline-flex items-center justify-center gap-2"
    >
      <Plus size={16} />
      New Repair
    </button>
  </div>

  <!-- Tickets Content -->
  {#if isLoading}
    <div class="bg-white dark:bg-slate-950 border border-slate-200 dark:border-slate-800 rounded-2xl p-16 text-center shadow-sm">
      <Loader2 class="animate-spin text-blue-600 mx-auto mb-4" size={32} />
      <p class="text-sm font-semibold text-slate-500">Loading repair database...</p>
    </div>
  {:else if filteredTickets.length === 0}
    <div class="bg-white dark:bg-slate-950 border border-slate-200 dark:border-slate-800 rounded-2xl p-16 text-center shadow-sm">
      <div class="w-12 h-12 bg-slate-50 dark:bg-slate-900 rounded-full flex items-center justify-center mx-auto mb-4 border border-slate-100 dark:border-slate-800 text-slate-400">
        <Search size={22} />
      </div>
      <h3 class="text-lg font-bold text-slate-900 dark:text-white mb-1">No repair tickets found</h3>
      <p class="text-sm text-slate-500 dark:text-slate-400 max-w-sm mx-auto">Try modifying your search or filters, or add a new repair ticket to get started.</p>
    </div>
  {:else}
    <!-- Table -->
    <div class="bg-white dark:bg-slate-950 border border-slate-200 dark:border-slate-800 rounded-2xl overflow-hidden shadow-sm">
      <div class="overflow-x-auto">
        <table class="w-full text-left border-collapse">
          <thead>
            <tr class="bg-slate-50 dark:bg-slate-900/50 border-b border-slate-200 dark:border-slate-800 text-xs font-bold text-slate-500 dark:text-slate-400 uppercase tracking-wider">
              <th class="py-4 px-6">Device</th>
              <th class="py-4 px-6">Customer</th>
              <th class="py-4 px-6">Issue</th>
              <th class="py-4 px-6">Price</th>
              <th class="py-4 px-6">Status</th>
              <th class="py-4 px-6 text-right">Quick Action</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-200 dark:divide-slate-800 text-sm text-slate-900 dark:text-slate-200">
            {#each filteredTickets as ticket (ticket.id)}
              <tr
                onclick={() => openEditDrawer(ticket)}
                class="hover:bg-slate-50/70 dark:hover:bg-slate-900/30 transition-colors cursor-pointer group"
              >
                <td class="py-4 px-6">
                  <div class="font-semibold text-slate-900 dark:text-white group-hover:text-blue-600 transition-colors">{ticket.brand}</div>
                  <div class="text-xs text-slate-500 dark:text-slate-400 mt-0.5">{ticket.model}</div>
                </td>
                <td class="py-4 px-6">
                  <div class="font-medium text-slate-900 dark:text-white">{ticket.customer_name}</div>
                  <div class="text-xs text-slate-500 dark:text-slate-400 mt-0.5">{ticket.customer_gender}</div>
                </td>
                <td class="py-4 px-6">
                  <div class="line-clamp-1 max-w-xs">{ticket.issue}</div>
                </td>
                <td class="py-4 px-6 font-semibold">
                  {formatCurrency(Number(ticket.price))}
                </td>
                <td class="py-4 px-6">
                  <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-semibold border {getStatusBadgeClass(ticket.status)}">
                    {getStatusLabel(ticket.status)}
                  </span>
                </td>
                <td class="py-4 px-6 text-right" onclick={(e) => e.stopPropagation()}>
                  {#if ticket.status !== 'picked_up'}
                    <button
                      onclick={() => handleQuickStatusUpdate(ticket.id, ticket.status)}
                      disabled={isActionLoading[ticket.id]}
                      class="px-4 py-1.5 bg-slate-900 hover:bg-slate-800 dark:bg-slate-800 dark:hover:bg-slate-700 disabled:bg-slate-100 disabled:text-slate-400 text-white font-bold text-xs uppercase tracking-wider rounded-lg transition-all active:scale-95 inline-flex items-center gap-1.5 shadow-sm"
                    >
                      {#if isActionLoading[ticket.id]}
                        <Loader2 size={12} class="animate-spin" />
                      {/if}
                      {#if ticket.status === 'service_in'}
                        Start Process
                      {:else if ticket.status === 'on_process'}
                        Mark Fixed
                      {:else if ticket.status === 'fixed'}
                        Pickup & Pay
                      {/if}
                    </button>
                  {:else}
                    <span class="text-xs font-semibold text-slate-400">Completed</span>
                  {/if}
                </td>
              </tr>
            {/each}
          </tbody>
        </table>
      </div>
    </div>
  {/if}
</div>

<!-- Create Modal -->
{#if showCreateModal}
  <div class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-900/60 backdrop-blur-sm animate-fade-in">
    <div class="bg-white dark:bg-slate-950 w-full max-w-2xl rounded-2xl border border-slate-200 dark:border-slate-800 shadow-2xl overflow-hidden flex flex-col max-h-[90vh]">
      <!-- Header -->
      <div class="flex justify-between items-center px-6 py-4 border-b border-slate-200 dark:border-slate-800 bg-slate-50 dark:bg-slate-900/30">
        <h3 class="font-bold text-lg text-slate-900 dark:text-white flex items-center gap-2">
          <Smartphone size={20} class="text-blue-600" />
          Intake New Device
        </h3>
        <button
          onclick={() => showCreateModal = false}
          class="p-1.5 rounded-lg text-slate-400 hover:bg-slate-100 dark:hover:bg-slate-900 transition-colors"
        >
          <X size={18} />
        </button>
      </div>

      <!-- Body -->
      <form onsubmit={handleCreateTicket} class="flex-1 overflow-y-auto p-6 space-y-6">
        <!-- Customer Info -->
        <div>
          <h4 class="text-xs font-bold text-slate-400 uppercase tracking-widest mb-3 flex items-center gap-1.5">
            <User size={14} />
            Customer Information
          </h4>
          <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <div class="space-y-1.5">
              <label for="create-customer-name" class="text-xs font-bold text-slate-500 block">Name *</label>
              <input
                id="create-customer-name"
                type="text"
                bind:value={createForm.customer_name}
                required
                placeholder="Customer Name"
                class="w-full bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-xl py-2 px-3 text-sm text-slate-900 dark:text-white focus:outline-none focus:border-blue-600 transition-colors"
              />
            </div>
            <div class="space-y-1.5">
              <label for="create-customer-gender" class="text-xs font-bold text-slate-500 block">Gender *</label>
              <select
                id="create-customer-gender"
                bind:value={createForm.customer_gender}
                class="w-full bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-xl py-2 px-3 text-sm text-slate-900 dark:text-white focus:outline-none focus:border-blue-600 transition-colors cursor-pointer"
              >
                <option value="Male">Male</option>
                <option value="Female">Female</option>
                <option value="Other">Other</option>
              </select>
            </div>
          </div>
        </div>

        <!-- Device Info -->
        <div>
          <h4 class="text-xs font-bold text-slate-400 uppercase tracking-widest mb-3 flex items-center gap-1.5">
            <Wrench size={14} />
            Device Details
          </h4>
          <div class="grid grid-cols-1 sm:grid-cols-2 gap-4 mb-4">
            <div class="space-y-1.5">
              <label for="create-brand" class="text-xs font-bold text-slate-500 block">Brand *</label>
              <input
                id="create-brand"
                type="text"
                bind:value={createForm.brand}
                required
                placeholder="e.g. Apple, Samsung"
                class="w-full bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-xl py-2 px-3 text-sm text-slate-900 dark:text-white focus:outline-none focus:border-blue-600 transition-colors"
              />
            </div>
            <div class="space-y-1.5">
              <label for="create-model" class="text-xs font-bold text-slate-500 block">Model *</label>
              <input
                id="create-model"
                type="text"
                bind:value={createForm.model}
                required
                placeholder="e.g. iPhone 15 Pro, Galaxy S24"
                class="w-full bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-xl py-2 px-3 text-sm text-slate-900 dark:text-white focus:outline-none focus:border-blue-600 transition-colors"
              />
            </div>
          </div>

          <div class="space-y-4">
            <div class="space-y-1.5">
              <label for="create-issue" class="text-xs font-bold text-slate-500 block">Issue Description *</label>
              <input
                id="create-issue"
                type="text"
                bind:value={createForm.issue}
                required
                placeholder="e.g. Broken LCD screen, Battery drain"
                class="w-full bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-xl py-2 px-3 text-sm text-slate-900 dark:text-white focus:outline-none focus:border-blue-600 transition-colors"
              />
            </div>
            <div class="space-y-1.5">
              <label for="create-additional-description" class="text-xs font-bold text-slate-500 block">Additional Description</label>
              <textarea
                id="create-additional-description"
                bind:value={createForm.additional_description}
                rows="3"
                placeholder="Any cosmetic damage or technical background details..."
                class="w-full bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-xl py-2 px-3 text-sm text-slate-900 dark:text-white focus:outline-none focus:border-blue-600 transition-colors resize-none"
              ></textarea>
            </div>
            <div class="space-y-1.5">
              <label for="create-accessories" class="text-xs font-bold text-slate-500 block">Accessories Left Behind</label>
              <input
                id="create-accessories"
                type="text"
                bind:value={createForm.accessories}
                placeholder="e.g. Charger, Case, SIM card tray"
                class="w-full bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-xl py-2 px-3 text-sm text-slate-900 dark:text-white focus:outline-none focus:border-blue-600 transition-colors"
              />
            </div>
          </div>
        </div>

        <!-- Price & Warranty -->
        <div>
          <h4 class="text-xs font-bold text-slate-400 uppercase tracking-widest mb-3 flex items-center gap-1.5">
            <DollarSign size={14} />
            Pricing & Warranty
          </h4>
          <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <div class="space-y-1.5">
              <label for="create-price" class="text-xs font-bold text-slate-500 block">Estimated Price (IDR)</label>
              <input
                id="create-price"
                type="number"
                bind:value={createForm.price}
                min="0"
                class="w-full bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-xl py-2 px-3 text-sm text-slate-900 dark:text-white focus:outline-none focus:border-blue-600 transition-colors"
              />
            </div>
            <div class="space-y-1.5">
              <label for="create-warranty-days" class="text-xs font-bold text-slate-500 block">Warranty Period (Days)</label>
              <input
                id="create-warranty-days"
                type="number"
                bind:value={createForm.warranty_days}
                min="0"
                class="w-full bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-xl py-2 px-3 text-sm text-slate-900 dark:text-white focus:outline-none focus:border-blue-600 transition-colors"
              />
            </div>
          </div>
        </div>

        <!-- Footer -->
        <div class="pt-6 border-t border-slate-200 dark:border-slate-800 flex justify-end gap-3">
          <button
            type="button"
            onclick={() => showCreateModal = false}
            class="px-4 py-2 border border-slate-200 dark:border-slate-800 text-slate-700 dark:text-slate-300 font-bold text-xs uppercase tracking-wider rounded-xl hover:bg-slate-50 dark:hover:bg-slate-900 transition-colors"
          >
            Cancel
          </button>
          <button
            type="submit"
            class="px-6 py-2 bg-blue-600 hover:bg-blue-700 text-white font-bold text-xs uppercase tracking-wider rounded-xl transition-colors shadow-sm"
          >
            Submit Intake
          </button>
        </div>
      </form>
    </div>
  </div>
{/if}

<!-- Edit Drawer -->
{#if selectedTicket}
  <!-- svelte-ignore a11y_click_events_have_key_events -->
  <!-- svelte-ignore a11y_no_static_element_interactions -->
  <div class="fixed inset-0 z-50 flex justify-end bg-slate-900/60 backdrop-blur-sm animate-fade-in" onclick={closeEditDrawer}>
    <!-- svelte-ignore a11y_click_events_have_key_events -->
    <!-- svelte-ignore a11y_no_static_element_interactions -->
    <div
      class="bg-white dark:bg-slate-950 w-full max-w-lg border-l border-slate-200 dark:border-slate-800 shadow-2xl h-full flex flex-col animate-slide-in"
      onclick={(e) => e.stopPropagation()}
    >
      <!-- Header -->
      <div class="flex justify-between items-center px-6 py-4 border-b border-slate-200 dark:border-slate-800 bg-slate-50 dark:bg-slate-900/30">
        <div>
          <h3 class="font-bold text-lg text-slate-900 dark:text-white">Repair Details</h3>
          <p class="text-xs text-slate-500 dark:text-slate-400 mt-0.5 font-mono">ID: {selectedTicket.id}</p>
        </div>
        <button
          onclick={closeEditDrawer}
          class="p-1.5 rounded-lg text-slate-400 hover:bg-slate-100 dark:hover:bg-slate-900 transition-colors"
        >
          <X size={18} />
        </button>
      </div>

      <!-- Body -->
      <form bind:this={editFormElement} onsubmit={handleUpdateTicket} class="flex-1 overflow-y-auto p-6 space-y-6">
        <!-- Status Panel -->
        <div class="bg-slate-50 dark:bg-slate-900/40 p-4 rounded-xl border border-slate-200 dark:border-slate-850 space-y-3">
          <div class="flex justify-between items-center">
            <span class="text-xs font-bold text-slate-500">Current Status</span>
            <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-semibold border {getStatusBadgeClass(editForm.status)}">
              {getStatusLabel(editForm.status)}
            </span>
          </div>
          <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
            <div class="space-y-1">
              <span class="text-[10px] font-black uppercase text-slate-400">Entry Date</span>
              <p class="text-xs font-semibold text-slate-700 dark:text-slate-300">{formatDate(selectedTicket.entry_date)}</p>
            </div>
            <div class="space-y-1">
              <span class="text-[10px] font-black uppercase text-slate-400">Exit Date</span>
              <p class="text-xs font-semibold text-slate-700 dark:text-slate-300">{selectedTicket.exit_date ? formatDate(selectedTicket.exit_date) : '-'}</p>
            </div>
          </div>
          {#if selectedTicket.warranty_expiry_date}
            <div class="pt-2 border-t border-slate-200/50 dark:border-slate-800/50 flex justify-between items-center">
              <span class="text-xs font-bold text-slate-500">Warranty Expiry</span>
              <span class="text-xs font-bold text-indigo-600 dark:text-indigo-400">{formatDate(selectedTicket.warranty_expiry_date)}</span>
            </div>
          {/if}
        </div>

        <!-- Customer Form Fields -->
        <div>
          <h4 class="text-xs font-bold text-slate-400 uppercase tracking-widest mb-3">Customer Profile</h4>
          <div class="grid grid-cols-2 gap-4">
            <div class="space-y-1.5">
              <label for="edit-customer-name" class="text-xs font-bold text-slate-500">Customer Name</label>
              <input
                id="edit-customer-name"
                type="text"
                bind:value={editForm.customer_name}
                class="w-full bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-xl py-2 px-3 text-sm text-slate-900 dark:text-white focus:outline-none focus:border-blue-600 transition-colors"
              />
            </div>
            <div class="space-y-1.5">
              <label for="edit-customer-gender" class="text-xs font-bold text-slate-500">Customer Gender</label>
              <select
                id="edit-customer-gender"
                bind:value={editForm.customer_gender}
                class="w-full bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-xl py-2 px-3 text-sm text-slate-900 dark:text-white focus:outline-none focus:border-blue-600 transition-colors cursor-pointer"
              >
                <option value="Male">Male</option>
                <option value="Female">Female</option>
                <option value="Other">Other</option>
              </select>
            </div>
          </div>
        </div>

        <!-- Device Info -->
        <div>
          <h4 class="text-xs font-bold text-slate-400 uppercase tracking-widest mb-3">Repair Information</h4>
          <div class="grid grid-cols-2 gap-4 mb-4">
            <div class="space-y-1.5">
              <label for="edit-brand" class="text-xs font-bold text-slate-500">Brand</label>
              <input
                id="edit-brand"
                type="text"
                bind:value={editForm.brand}
                class="w-full bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-xl py-2 px-3 text-sm text-slate-900 dark:text-white focus:outline-none focus:border-blue-600 transition-colors"
              />
            </div>
            <div class="space-y-1.5">
              <label for="edit-model" class="text-xs font-bold text-slate-500">Model</label>
              <input
                id="edit-model"
                type="text"
                bind:value={editForm.model}
                class="w-full bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-xl py-2 px-3 text-sm text-slate-900 dark:text-white focus:outline-none focus:border-blue-600 transition-colors"
              />
            </div>
          </div>

          <div class="space-y-4">
            <div class="space-y-1.5">
              <label for="edit-issue" class="text-xs font-bold text-slate-500">Issue</label>
              <input
                id="edit-issue"
                type="text"
                bind:value={editForm.issue}
                class="w-full bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-xl py-2 px-3 text-sm text-slate-900 dark:text-white focus:outline-none focus:border-blue-600 transition-colors"
              />
            </div>
            <div class="space-y-1.5">
              <label for="edit-additional-description" class="text-xs font-bold text-slate-500">Additional Description</label>
              <textarea
                id="edit-additional-description"
                bind:value={editForm.additional_description}
                rows="3"
                class="w-full bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-xl py-2 px-3 text-sm text-slate-900 dark:text-white focus:outline-none focus:border-blue-600 transition-colors resize-none"
              ></textarea>
            </div>
            <div class="space-y-1.5">
              <label for="edit-accessories" class="text-xs font-bold text-slate-500">Accessories</label>
              <input
                id="edit-accessories"
                type="text"
                bind:value={editForm.accessories}
                class="w-full bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-xl py-2 px-3 text-sm text-slate-900 dark:text-white focus:outline-none focus:border-blue-600 transition-colors"
              />
            </div>
          </div>
        </div>

        <!-- Finances & Status overrides -->
        <div>
          <h4 class="text-xs font-bold text-slate-400 uppercase tracking-widest mb-3">Billing & Parameters</h4>
          <div class="grid grid-cols-2 gap-4 mb-4">
            <div class="space-y-1.5">
              <label for="edit-price" class="text-xs font-bold text-slate-500">Price (IDR)</label>
              <input
                id="edit-price"
                type="number"
                bind:value={editForm.price}
                class="w-full bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-xl py-2 px-3 text-sm text-slate-900 dark:text-white focus:outline-none focus:border-blue-600 transition-colors"
              />
            </div>
            <div class="space-y-1.5">
              <label for="edit-warranty-days" class="text-xs font-bold text-slate-500">Warranty Period (Days)</label>
              <input
                id="edit-warranty-days"
                type="number"
                bind:value={editForm.warranty_days}
                class="w-full bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-xl py-2 px-3 text-sm text-slate-900 dark:text-white focus:outline-none focus:border-blue-600 transition-colors"
              />
            </div>
          </div>

          <div class="grid grid-cols-2 gap-4">
            <div class="space-y-1.5">
              <label for="edit-status" class="text-xs font-bold text-slate-500">Status</label>
              <select
                id="edit-status"
                bind:value={editForm.status}
                class="w-full bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-xl py-2.5 px-3 text-sm text-slate-900 dark:text-white focus:outline-none focus:border-blue-600 transition-colors cursor-pointer"
              >
                <option value="service_in">Service In</option>
                <option value="on_process">On Process</option>
                <option value="fixed">Fixed</option>
                <option value="picked_up">Picked Up</option>
              </select>
            </div>
            <div class="space-y-1.5">
              <label for="edit-payment-status" class="text-xs font-bold text-slate-500">Payment Status</label>
              <select
                id="edit-payment-status"
                bind:value={editForm.payment_status}
                class="w-full bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-xl py-2.5 px-3 text-sm text-slate-900 dark:text-white focus:outline-none focus:border-blue-600 transition-colors cursor-pointer"
              >
                <option value="unpaid">Unpaid</option>
                <option value="paid">Paid</option>
              </select>
            </div>
          </div>
        </div>

        <!-- Footer Actions -->
        <div class="pt-6 border-t border-slate-200 dark:border-slate-800 flex justify-between items-center">
          <button
            type="button"
            onclick={deleteSelectedTicket}
            class="px-4 py-2 bg-red-50 hover:bg-red-100 text-red-600 font-bold text-xs uppercase tracking-wider rounded-xl transition-colors inline-flex items-center gap-1.5 border border-red-200"
          >
            <Trash2 size={14} />
            Delete Ticket
          </button>

          <div class="flex gap-3">
            <button
              type="button"
              onclick={closeEditDrawer}
              class="px-4 py-2 border border-slate-200 dark:border-slate-800 text-slate-700 dark:text-slate-300 font-bold text-xs uppercase tracking-wider rounded-xl hover:bg-slate-50 dark:hover:bg-slate-900 transition-colors"
            >
              Cancel
            </button>
            <button
              type="button"
              onclick={submitEditForm}
              class="px-6 py-2 bg-blue-600 hover:bg-blue-700 text-white font-bold text-xs uppercase tracking-wider rounded-xl transition-colors shadow-sm"
            >
              Save Changes
            </button>
          </div>
        </div>
      </form>
    </div>
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
