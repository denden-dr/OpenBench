<script lang="ts">
  import { Button, Card } from '$lib';
  import { User, LogOut, Plus } from 'lucide-svelte';
  import { goto } from '$app/navigation';
  import type { UserSession } from '$lib/services/auth';

  interface Props {
    session: UserSession | null;
    onsignout: () => void;
  }

  let { session, onsignout }: Props = $props();
</script>

<Card class="p-6 bg-white flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
  <div>
    <h1 class="font-display font-black text-3xl text-neubrutalism-charcoal uppercase tracking-tight">
      CUSTOMER PORTAL
    </h1>
    {#if session}
      <div class="flex items-center gap-2 mt-2 font-mono text-sm text-neubrutalism-charcoal">
        <User class="w-4 h-4 text-neubrutalism-pink" />
        <span>Welcome back, <strong class="text-neubrutalism-pink uppercase">{session.full_name}</strong> (@{session.username})</span>
      </div>
    {/if}
  </div>

  <div class="flex gap-3 w-full sm:w-auto">
    <Button 
      bgColor="bg-neubrutalism-yellow" 
      class="flex items-center gap-2 uppercase font-bold text-sm shadow-neubrutalism-sm"
      onclick={() => goto('/repair-request')}
    >
      <Plus class="w-4 h-4" />
      <span>Book Repair</span>
    </Button>

    <Button 
      bgColor="bg-zinc-100" 
      class="flex items-center gap-2 uppercase font-bold text-sm shadow-neubrutalism-sm hover:bg-rose-100 hover:text-rose-700 transition-colors"
      onclick={onsignout}
    >
      <LogOut class="w-4 h-4" />
      <span>Sign Out</span>
    </Button>
  </div>
</Card>
