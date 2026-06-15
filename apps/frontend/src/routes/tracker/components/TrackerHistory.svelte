<script lang="ts">
  import { Card } from '$lib';
  import { History, Trash2 } from 'lucide-svelte';

  interface HistoryItem {
    id: string;
    label: string;
    timestamp: string;
  }

  let {
    searchHistory,
    onSelect,
    onClear
  }: {
    searchHistory: HistoryItem[];
    onSelect: (id: string) => void;
    onClear: () => void;
  } = $props();
</script>

<Card bgColor="bg-white" class="border-4 border-neubrutalism-charcoal p-6 shadow-neubrutalism-sm flex flex-col gap-3">
  <h3 class="font-display font-bold text-sm uppercase text-zinc-700 flex items-center justify-between border-b border-dashed border-zinc-200 pb-2">
    <span class="flex items-center gap-1.5">
      <History class="w-4 h-4 text-zinc-500" />
      <span>Recent Tracking History</span>
    </span>
    <button 
      onclick={onClear}
      class="text-[9px] text-rose-500 font-mono font-bold uppercase hover:underline flex items-center gap-1"
    >
      <Trash2 class="w-3.5 h-3.5" />
      <span>Clear All</span>
    </button>
  </h3>

  <div class="flex flex-col gap-2">
    {#each searchHistory as item}
      <button 
        onclick={() => onSelect(item.id)}
        class="w-full text-left bg-zinc-50 border-2 border-neubrutalism-charcoal p-3 hover:bg-zinc-100 transition-colors flex items-center justify-between gap-4 cursor-pointer active:translate-y-0.5 shadow-neubrutalism-sm active:shadow-none"
      >
        <div class="flex flex-col">
          <span class="font-sans font-bold text-xs text-neubrutalism-charcoal">{item.label}</span>
          <span class="font-mono text-[9px] text-zinc-500 mt-1">ID: {item.id}</span>
        </div>
        <span class="font-mono text-[9px] text-zinc-400 font-semibold">{item.timestamp}</span>
      </button>
    {/each}
  </div>
</Card>
