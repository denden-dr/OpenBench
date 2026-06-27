<script lang="ts">
  import { Card, Button } from '$lib';
  import { Search, ShieldAlert } from 'lucide-svelte';

  let {
    ticketIdInput = $bindable(''),
    isLoading,
    errorMessage,
    onTrack
  }: {
    ticketIdInput: string;
    isLoading: boolean;
    errorMessage: string;
    onTrack: (e?: Event) => void;
  } = $props();
</script>

<Card bgColor="bg-white" class="border-4 border-neubrutalism-charcoal p-6 shadow-neubrutalism-md flex flex-col gap-4">
  <form onsubmit={onTrack} class="flex flex-col sm:flex-row gap-3">
    <div class="relative flex-grow">
      <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none text-zinc-400">
        <Search class="w-4 h-4" />
      </div>
      <input 
        type="text" 
        aria-label="Tracking Number"
        placeholder="Enter Tracking Number (Example: OB-202606-0001-A9X2)"
        bind:value={ticketIdInput}
        disabled={isLoading}
        class="w-full pl-9 pr-4 py-3 border-4 border-neubrutalism-charcoal bg-white focus:outline-none focus:bg-zinc-50 font-mono text-xs shadow-neubrutalism-sm"
      />
    </div>
    <Button 
      bgColor="bg-neubrutalism-yellow" 
      type="submit"
      disabled={isLoading}
      class="py-3 px-6 font-display font-bold uppercase border-2 border-neubrutalism-charcoal shadow-neubrutalism-sm flex items-center justify-center gap-2 shrink-0"
    >
      <span>{isLoading ? 'TRACKING...' : 'TRACK PROGRESS'}</span>
    </Button>
  </form>

  {#if errorMessage}
    <div class="bg-neubrutalism-pink text-white border-4 border-neubrutalism-charcoal p-4 font-mono text-xs flex items-center gap-2">
      <ShieldAlert class="w-5 h-5 shrink-0" />
      <span>{errorMessage}</span>
    </div>
  {/if}
</Card>
