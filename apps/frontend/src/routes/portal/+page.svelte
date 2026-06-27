<script lang="ts">
  import { onMount } from 'svelte';
  import { Button } from '$lib';
  import { authService, type UserSession } from '$lib/services/auth';
  import { ticketService, type Ticket } from '$lib/services/ticket';
  import { goto } from '$app/navigation';
  import { AlertCircle } from 'lucide-svelte';

  import PortalHeader from './components/PortalHeader.svelte';
  import PortalStats from './components/PortalStats.svelte';
  import PortalTicketsList from './components/PortalTicketsList.svelte';

  let session = $state<UserSession | null>(null);
  let tickets = $state<Ticket[]>([]);
  let loading = $state(true);
  let error = $state('');

  onMount(async () => {
    try {
      // Fetch or verify session
      const currentSession = await authService.checkSession();
      if (!currentSession) {
        await goto('/auth/signin');
        return;
      }
      
      // Redirect to onboarding profile completion if details are missing
      if (!currentSession.username || !currentSession.full_name) {
        await goto('/portal/setup');
        return;
      }

      session = currentSession;

      // Fetch user's tickets
      tickets = await ticketService.getMyTickets();
    } catch (err: any) {
      error = err.message || 'Failed to load dashboard.';
    } finally {
      loading = false;
    }
  });

  async function handleSignOut() {
    try {
      await authService.signOut();
      await goto('/auth/signin');
    } catch (err: any) {
      error = 'Failed to sign out.';
    }
  }

  // Derived stats
  const activeTicketsCount = $derived(
    tickets.filter(t => t.status !== 'completed' && t.status !== 'cancelled').length
  );
  const completedTicketsCount = $derived(
    tickets.filter(t => t.status === 'completed').length
  );
</script>

<svelte:head>
  <title>Customer Portal - OpenBench</title>
</svelte:head>

<main class="min-h-screen bg-neubrutalism-bg p-4 sm:p-8">
  <div class="max-w-6xl mx-auto flex flex-col gap-8">
    
    <!-- Header -->
    <PortalHeader {session} onsignout={handleSignOut} />

    {#if loading}
      <!-- Loading Skeleton -->
      <div class="grid grid-cols-1 md:grid-cols-3 gap-6 animate-pulse">
        <div class="h-28 bg-zinc-200 border-4 border-neubrutalism-charcoal"></div>
        <div class="h-28 bg-zinc-200 border-4 border-neubrutalism-charcoal"></div>
        <div class="h-28 bg-zinc-200 border-4 border-neubrutalism-charcoal"></div>
      </div>
      <div class="h-64 bg-zinc-200 border-4 border-neubrutalism-charcoal animate-pulse"></div>
    {:else if error}
      <!-- Error Display -->
      <div class="border-4 border-neubrutalism-charcoal bg-rose-100 p-6 shadow-neubrutalism-md text-center">
        <AlertCircle class="w-12 h-12 text-neubrutalism-pink mx-auto mb-4" />
        <h3 class="font-display font-bold text-lg text-neubrutalism-charcoal uppercase">Error Loading Portal</h3>
        <p class="font-sans text-sm mt-1 text-neubrutalism-charcoal opacity-80">{error}</p>
        <Button bgColor="bg-white" class="mt-4 border-2 border-neubrutalism-charcoal shadow-neubrutalism-sm" onclick={() => window.location.reload()}>
          RETRY
        </Button>
      </div>
    {:else}
      <!-- Stats Panel -->
      <PortalStats 
        {activeTicketsCount} 
        {completedTicketsCount} 
        totalTicketsCount={tickets.length} 
      />

      <!-- Tickets Grid/List -->
      <PortalTicketsList {tickets} />
    {/if}
  </div>
</main>
