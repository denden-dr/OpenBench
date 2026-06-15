<script lang="ts">
  import { toastService } from '$lib/services/toast.svelte';
  import { fly } from 'svelte/transition';
  import { ShieldAlert, CheckCircle2, AlertTriangle, Info } from 'lucide-svelte';
</script>

<div class="fixed bottom-6 right-6 z-50 flex flex-col gap-3 max-w-sm w-full pointer-events-none">
  {#each toastService.messages as toast (toast.id)}
    <div 
      transition:fly={{ y: 20, duration: 250 }}
      class="pointer-events-auto border-4 border-neubrutalism-charcoal p-4 font-mono text-xs font-bold uppercase shadow-neubrutalism-md flex items-start gap-3
        {toast.type === 'error' || toast.type === 'warning' 
          ? 'bg-neubrutalism-pink text-white' 
          : toast.type === 'success' 
            ? 'bg-neubrutalism-green text-neubrutalism-charcoal' 
            : 'bg-neubrutalism-yellow text-neubrutalism-charcoal'}"
    >
      {#if toast.type === 'error' || toast.type === 'warning'}
        <ShieldAlert class="w-4 h-4 shrink-0 mt-0.5" />
      {:else if toast.type === 'success'}
        <CheckCircle2 class="w-4 h-4 shrink-0 mt-0.5" />
      {:else}
        <Info class="w-4 h-4 shrink-0 mt-0.5" />
      {/if}
      
      <span class="flex-grow leading-tight break-words">{toast.message}</span>
      
      <button 
        onclick={() => toastService.dismiss(toast.id)}
        class="font-sans font-bold hover:opacity-80 shrink-0 select-none text-sm leading-none focus:outline-none"
      >
        ✕
      </button>
    </div>
  {/each}
</div>
