<script lang="ts">
  import { ticketService, type Ticket, type PublicTrackerTicket } from '$lib/services/ticket';
  import { onMount } from 'svelte';

  import TrackerHeader from './components/TrackerHeader.svelte';
  import TrackerInputForm from './components/TrackerInputForm.svelte';
  import TrackerProgressCard from './components/TrackerProgressCard.svelte';
  import TrackerHistory from './components/TrackerHistory.svelte';

  let ticketIdInput = $state('');
  let searchedTicket = $state<PublicTrackerTicket | null>(null);
  let isLoading = $state(false);
  let errorMessage = $state('');
  let hydrated = $state(false);
  let abortController: AbortController | null = null;

  // Search history state
  interface HistoryItem {
    id: string;
    label: string;
    timestamp: string;
  }
  let searchHistory = $state<HistoryItem[]>([]);

  onMount(() => {
    loadHistory();
    hydrated = true;
  });

  function loadHistory() {
    const saved = localStorage.getItem('openbench_tracked_history');
    if (saved) {
      try {
        searchHistory = JSON.parse(saved);
      } catch (e) {
        searchHistory = [];
      }
    }
  }

  function saveToHistory(ticket: PublicTrackerTicket) {
    // Remove if already exists to move to top
    let history = searchHistory.filter(item => item.id !== ticket.id);
    
    // Add to top of list
    const label = `${ticket.brand_phone} ${ticket.model_phone} (${ticket.id.slice(0, 8)})`;
    history.unshift({
      id: ticket.id,
      label,
      timestamp: new Date().toLocaleDateString('en-US', { hour: '2-digit', minute: '2-digit' })
    });

    // Cap at 5 items
    history = history.slice(0, 5);
    searchHistory = history;
    localStorage.setItem('openbench_tracked_history', JSON.stringify(history));
  }

  function clearHistory() {
    searchHistory = [];
    localStorage.removeItem('openbench_tracked_history');
  }

  // UUID regex pattern helper
  const uuidRegex = /^[0-9a-f]{8}-[0-9a-f]{4}-[1-5][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$/i;

  async function handleTrack(e?: Event) {
    if (e) e.preventDefault();
    errorMessage = '';
    searchedTicket = null;

    const trimmedInput = ticketIdInput.trim();
    if (!trimmedInput) {
      errorMessage = 'Please enter a Ticket ID first.';
      return;
    }

    // Security check: Reject Ticket Number (starts with OB-)
    if (trimmedInput.toUpperCase().startsWith('OB-')) {
      errorMessage = 'Public tracking is only allowed using Ticket ID (UUID), not Ticket Number.';
      return;
    }

    // Format validation: Validate UUID format
    if (!uuidRegex.test(trimmedInput)) {
      errorMessage = 'Invalid Ticket ID format. Make sure you entered the complete UUID code provided.';
      return;
    }

    if (abortController) {
      abortController.abort();
    }
    abortController = new AbortController();

    isLoading = true;
    try {
      const ticket = await ticketService.getPublicTrackerTicket(trimmedInput, abortController.signal);
      if (ticket) {
        searchedTicket = ticket;
        saveToHistory(ticket);
      } else {
        errorMessage = 'Ticket not found. Please double-check your Ticket ID.';
      }
    } catch (err: any) {
      if (err.name === 'AbortError') return;
      errorMessage = err.message || 'Failed to retrieve tracking status data.';
    } finally {
      isLoading = false;
    }
  }

  function handleSelectHistory(id: string) {
    ticketIdInput = id;
    handleTrack();
  }

  const getStatusColor = (statusVal: string) => {
    switch (statusVal) {
      case 'ready_for_pickup': return 'bg-neubrutalism-green border-neubrutalism-charcoal';
      case 'in_repair': return 'bg-neubrutalism-yellow border-neubrutalism-charcoal';
      case 'received': return 'bg-zinc-200';
      case 'completed': return 'bg-zinc-100 text-zinc-500 border-zinc-300';
      case 'cancelled': return 'bg-rose-100 text-rose-600 border-rose-300';
      default: return 'bg-white';
    }
  };

  const getStatusText = (statusVal: string) => {
    return statusVal.replace(/_/g, ' ');
  };
</script>

<svelte:head>
  <title>Track Ticket Status - OpenBench Tracker</title>
</svelte:head>

<main class="min-h-screen bg-neubrutalism-bg flex flex-col font-sans p-4 md:p-8" data-hydrated={hydrated}>
  <div class="max-w-2xl w-full mx-auto flex flex-col gap-6">
    
    <!-- Home Navigation & Title -->
    <TrackerHeader />

    <!-- Input Form Card -->
    <TrackerInputForm
      bind:ticketIdInput
      {isLoading}
      {errorMessage}
      onTrack={handleTrack}
    />

    <!-- Searched Ticket Progress -->
    {#if searchedTicket}
      <TrackerProgressCard
        {searchedTicket}
        {getStatusColor}
        {getStatusText}
      />
    {/if}

    <!-- Search History Card -->
    {#if searchHistory.length > 0}
      <TrackerHistory
        {searchHistory}
        onSelect={handleSelectHistory}
        onClear={clearHistory}
      />
    {/if}

  </div>
</main>
