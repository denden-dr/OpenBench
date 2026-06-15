<script lang="ts">
  import { Card, Button } from '$lib';
  import { ticketService, type Ticket } from '$lib/services/ticket';
  import { onMount } from 'svelte';
  import { 
    Wrench, Plus, Search, Filter, 
    ArrowUpRight, Phone, Calendar, User 
  } from 'lucide-svelte';

  let tickets = $state<Ticket[]>([]);
  let isLoading = $state(true);
  let activeTab = $state<'active' | 'archive'>('active');
  let searchQuery = $state('');
  let statusFilter = $state<string>('all');

  let fetchId = 0;

  async function loadTickets() {
    const currentFetchId = ++fetchId;
    isLoading = true;
    try {
      const data = await ticketService.getTickets();
      if (currentFetchId === fetchId) {
        tickets = data;
        isLoading = false;
      }
    } catch (err) {
      console.error('Error fetching tickets:', err);
      if (currentFetchId === fetchId) {
        isLoading = false;
      }
    }
  }

  async function switchTab(tab: 'active' | 'archive') {
    activeTab = tab;
    statusFilter = 'all';
    await loadTickets();
  }

  onMount(() => {
    loadTickets();
  });

  // Filter tickets reactively using Svelte 5 derived states
  let filteredTickets = $derived(
    tickets.filter(ticket => {
      // 1. Tab filter
      const isArchived = ticket.status === 'picked_up' || ticket.status === 'cancelled';
      if (activeTab === 'active' && isArchived) return false;
      if (activeTab === 'archive' && !isArchived) return false;

      // 2. Status filter dropdown
      if (statusFilter !== 'all' && ticket.status !== statusFilter) return false;

      // 3. Search query
      if (searchQuery.trim() !== '') {
        const query = searchQuery.toLowerCase();
        const numMatch = ticket.ticket_number.toLowerCase().includes(query);
        const nameMatch = ticket.customer_name.toLowerCase().includes(query);
        const phoneMatch = ticket.customer_phone.includes(query);
        const deviceMatch = `${ticket.brand_phone} ${ticket.model_phone}`.toLowerCase().includes(query);
        return numMatch || nameMatch || phoneMatch || deviceMatch;
      }

      return true;
    })
  );

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'ready_for_pickup': return 'bg-neubrutalism-green';
      case 'in_repair': return 'bg-neubrutalism-yellow';
      case 'diagnosing': return 'bg-neubrutalism-pink text-white';
      case 'received': return 'bg-zinc-200';
      case 'picked_up': return 'bg-zinc-100 text-zinc-500';
      case 'cancelled': return 'bg-rose-100 text-rose-600 border-rose-300';
      default: return 'bg-white';
    }
  };

  const getStatusText = (status: string) => {
    return status.replace(/_/g, ' ');
  };

  const formatDate = (dateStr: string) => {
    const d = new Date(dateStr);
    return d.toLocaleDateString('en-US', { day: 'numeric', month: 'short', year: 'numeric' });
  };
</script>

<svelte:head>
  <title>Repair Tickets - OpenBench Admin</title>
</svelte:head>

<div class="flex flex-col gap-6">
  
  <!-- Page Header Area -->
  <div class="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
    <div>
      <h2 class="font-display font-extrabold text-2xl md:text-3xl uppercase tracking-tight">
        Repair Tickets Management
      </h2>
      <p class="font-sans text-xs sm:text-sm text-zinc-500 font-semibold mt-1">
        Manage hardware diagnostics, active repairs, component changes, and archives.
      </p>
    </div>
    
    <a href="/admin/tickets/new" class="w-full sm:w-auto">
      <Button bgColor="bg-neubrutalism-yellow" class="w-full sm:w-auto flex items-center justify-center gap-2 py-2 px-4 font-bold shadow-neubrutalism-sm">
        <Plus class="w-4 h-4" />
        <span>NEW TICKET</span>
      </Button>
    </a>
  </div>

  <!-- Tabs Switcher -->
  <div class="flex border-b-4 border-neubrutalism-charcoal font-display font-bold">
    <button 
      class="px-6 py-3 border-2 border-b-0 border-neubrutalism-charcoal transition-all text-sm uppercase tracking-wider {activeTab === 'active' ? 'bg-neubrutalism-green border-b-4 border-b-neubrutalism-green -mb-1 translate-y-0.5' : 'bg-white hover:bg-zinc-50'}"
      onclick={() => switchTab('active')}
    >
      Active Repairs ({tickets.filter(t => t.status !== 'picked_up' && t.status !== 'cancelled').length})
    </button>
    <button 
      class="px-6 py-3 border-2 border-l-0 border-b-0 border-neubrutalism-charcoal transition-all text-sm uppercase tracking-wider {activeTab === 'archive' ? 'bg-neubrutalism-pink text-white border-b-4 border-b-neubrutalism-pink -mb-1 translate-y-0.5' : 'bg-white hover:bg-zinc-50'}"
      onclick={() => switchTab('archive')}
    >
      Archive ({tickets.filter(t => t.status === 'picked_up' || t.status === 'cancelled').length})
    </button>
  </div>

  <!-- Search and Filtering Bar -->
  <div class="flex flex-col md:flex-row gap-4 items-center justify-between">
    <!-- Search Box -->
    <div class="relative w-full md:max-w-md">
      <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none text-zinc-500">
        <Search class="w-4 h-4" />
      </div>
      <input 
        type="text" 
        placeholder="Search ticket number, customer name, device..."
        bind:value={searchQuery}
        class="w-full pl-9 pr-4 py-2 border-4 border-neubrutalism-charcoal bg-white focus:outline-none focus:bg-zinc-50 font-mono text-sm shadow-neubrutalism-sm"
      />
    </div>

    <!-- Status Dropdown Filter -->
    <div class="flex items-center gap-2 w-full md:w-auto">
      <Filter class="w-4 h-4 text-zinc-600 shrink-0" />
      <span class="font-mono text-xs font-bold shrink-0">STATUS:</span>
      <select 
        bind:value={statusFilter}
        class="w-full md:w-48 bg-white border-4 border-neubrutalism-charcoal px-3 py-1.5 font-mono text-xs shadow-neubrutalism-sm focus:outline-none"
      >
        <option value="all">ALL STATUS</option>
        {#if activeTab === 'active'}
          <option value="received">RECEIVED</option>
          <option value="diagnosing">DIAGNOSING</option>
          <option value="in_repair">IN REPAIR</option>
          <option value="ready_for_pickup">READY FOR PICKUP</option>
        {:else}
          <option value="picked_up">PICKED UP (COMPLETED)</option>
          <option value="cancelled">CANCELLED</option>
        {/if}
      </select>
    </div>
  </div>

  <!-- Main Repairs List -->
  {#if isLoading}
    <div class="flex flex-col gap-4">
      {#each Array(4) as _}
        <div class="h-24 bg-zinc-200 border-4 border-neubrutalism-charcoal animate-pulse"></div>
      {/each}
    </div>
  {:else}
    <div class="flex flex-col gap-4">
      {#if filteredTickets.length === 0}
        <Card bgColor="bg-white" class="p-12 text-center font-mono text-sm text-zinc-500">
          NO TICKETS MATCH THE CURRENT SEARCH / FILTER.
        </Card>
      {:else}
        {#each filteredTickets as ticket}
          <Card bgColor="bg-white" class="hover:shadow-neubrutalism-lg transition-all border-4 border-neubrutalism-charcoal">
            <div class="flex flex-col md:flex-row md:items-center justify-between gap-4">
              
              <!-- Left side info -->
              <div class="flex-grow flex flex-col sm:flex-row sm:items-start gap-4">
                
                <div class="flex flex-col gap-1 shrink-0">
                  <span class="font-mono text-xs font-extrabold bg-zinc-200 px-2 py-0.5 border-2 border-neubrutalism-charcoal text-neubrutalism-charcoal shadow-neubrutalism-sm w-fit">
                    {ticket.ticket_number}
                  </span>
                  
                  <div class="flex items-center gap-1.5 font-mono text-[10px] text-zinc-500 mt-1">
                    <Calendar class="w-3.5 h-3.5" />
                    <span>{formatDate(ticket.created_at)}</span>
                  </div>
                </div>

                <div class="flex flex-col gap-1">
                  <h3 class="font-display font-bold text-lg leading-tight">
                    {ticket.brand_phone} {ticket.model_phone}
                  </h3>
                  
                  <div class="flex flex-wrap items-center gap-x-4 gap-y-1 font-sans text-xs text-zinc-700 mt-1">
                    <span class="flex items-center gap-1">
                      <User class="w-3.5 h-3.5 text-zinc-500" />
                      <strong>{ticket.customer_name}</strong>
                    </span>
                    <span class="flex items-center gap-1">
                      <Phone class="w-3.5 h-3.5 text-zinc-500" />
                      <span>{ticket.customer_phone}</span>
                    </span>
                  </div>

                  <p class="font-sans text-xs text-neubrutalism-charcoal mt-1 line-clamp-1">
                    <span class="font-semibold text-zinc-600">Damage:</span> {ticket.damage_description}
                  </p>
                </div>

              </div>

              <!-- Right side actions & status -->
              <div class="flex flex-row md:flex-col items-center md:items-end justify-between md:justify-center gap-3 shrink-0 border-t-2 md:border-t-0 border-dashed border-zinc-200 pt-3 md:pt-0">
                <div class="flex items-center gap-2">
                  <span class="font-mono text-xs font-bold py-1 px-3 border-2 border-neubrutalism-charcoal shadow-neubrutalism-sm uppercase tracking-wide {getStatusColor(ticket.status)}">
                    {getStatusText(ticket.status)}
                  </span>
                </div>

                <div class="flex items-center gap-3">
                  <span class="font-mono text-sm font-extrabold">
                    IDR {ticket.cost.toLocaleString('id-ID')}
                  </span>
                  
                  <a href="/admin/tickets/{ticket.id}">
                    <button class="p-2 border-2 border-neubrutalism-charcoal hover:bg-zinc-100 transition duration-150 cursor-pointer shadow-neubrutalism-sm active:translate-x-0.5 active:translate-y-0.5 active:shadow-none bg-white">
                      <ArrowUpRight class="w-4 h-4 text-neubrutalism-charcoal" />
                    </button>
                  </a>
                </div>
              </div>

            </div>
          </Card>
        {/each}
      {/if}
    </div>
  {/if}

</div>
