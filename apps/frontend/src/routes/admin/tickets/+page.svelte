<script lang="ts">
  import { ticketService, type Ticket } from '$lib/services/ticket';
  import { onMount } from 'svelte';

  import TicketsHeader from './components/TicketsHeader.svelte';
  import TicketsSearchFilters from './components/TicketsSearchFilters.svelte';
  import TicketsTable from './components/TicketsTable.svelte';

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
        tickets = data ?? [];
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
      const displayStatus = getDisplayStatus(ticket);
      const isArchived = displayStatus === 'completed' || displayStatus === 'cancelled';
      if (activeTab === 'active' && isArchived) return false;
      if (activeTab === 'archive' && !isArchived) return false;

      // 2. Status filter dropdown
      if (statusFilter !== 'all' && displayStatus !== statusFilter) return false;

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

  const getDisplayStatus = (t: Ticket) => {
    return t.status === 'completed' && t.device_position === 'warehouse'
      ? 'ready_for_pickup'
      : t.status;
  };
</script>

<svelte:head>
  <title>Repair Tickets - OpenBench Admin</title>
</svelte:head>

<div class="flex flex-col gap-6">
  
  <!-- Page Header -->
  <TicketsHeader />

  <!-- Tabs Switcher -->
  <div class="flex border-b-4 border-neubrutalism-charcoal font-display font-bold">
    <button 
      class="px-6 py-3 border-2 border-b-0 border-neubrutalism-charcoal transition-all text-sm uppercase tracking-wider {activeTab === 'active' ? 'bg-neubrutalism-green border-b-4 border-b-neubrutalism-green -mb-1 translate-y-0.5' : 'bg-white hover:bg-zinc-50'}"
      onclick={() => switchTab('active')}
    >
      Active Repairs ({tickets.filter(t => getDisplayStatus(t) !== 'completed' && getDisplayStatus(t) !== 'cancelled').length})
    </button>
    <button 
      class="px-6 py-3 border-2 border-l-0 border-b-0 border-neubrutalism-charcoal transition-all text-sm uppercase tracking-wider {activeTab === 'archive' ? 'bg-neubrutalism-pink text-white border-b-4 border-b-neubrutalism-pink -mb-1 translate-y-0.5' : 'bg-white hover:bg-zinc-50'}"
      onclick={() => switchTab('archive')}
    >
      Archive ({tickets.filter(t => getDisplayStatus(t) === 'completed' || getDisplayStatus(t) === 'cancelled').length})
    </button>
  </div>

  <!-- Search and Filtering Bar -->
  <TicketsSearchFilters
    bind:searchQuery
    bind:statusFilter
    {activeTab}
  />

  <!-- Main Repairs List -->
  <TicketsTable 
    {filteredTickets} 
    {isLoading} 
  />

</div>
