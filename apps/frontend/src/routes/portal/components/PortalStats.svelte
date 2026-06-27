<script lang="ts">
  import { TrendingUp, CheckCircle2, Wrench } from 'lucide-svelte';

  interface Props {
    activeTicketsCount: number;
    completedTicketsCount: number;
    totalTicketsCount: number;
  }

  let { activeTicketsCount, completedTicketsCount, totalTicketsCount }: Props = $props();
</script>

<div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-6">
  <!-- Stat: Active -->
  <div class="bg-amber-100 border-4 border-neubrutalism-charcoal p-5 shadow-neubrutalism-sm relative overflow-hidden flex flex-col justify-between min-h-[110px]">
    <div>
      <span class="font-mono text-xs font-bold text-zinc-500 uppercase tracking-wider block">Active Repairs</span>
      <span class="font-display font-black text-3xl text-neubrutalism-charcoal mt-1 block">{activeTicketsCount}</span>
    </div>
    <TrendingUp class="absolute right-4 bottom-4 w-10 h-10 text-amber-500 opacity-20" />
  </div>

  <!-- Stat: Completed -->
  <div class="bg-emerald-100 border-4 border-neubrutalism-charcoal p-5 shadow-neubrutalism-sm relative overflow-hidden flex flex-col justify-between min-h-[110px]">
    <div>
      <span class="font-mono text-xs font-bold text-zinc-500 uppercase tracking-wider block">Completed Repairs</span>
      <span class="font-display font-black text-3xl text-neubrutalism-charcoal mt-1 block">{completedTicketsCount}</span>
    </div>
    <CheckCircle2 class="absolute right-4 bottom-4 w-10 h-10 text-emerald-500 opacity-20" />
  </div>

  <!-- Stat: Total Requests -->
  <div class="bg-purple-100 border-4 border-neubrutalism-charcoal p-5 shadow-neubrutalism-sm relative overflow-hidden flex flex-col justify-between min-h-[110px] sm:col-span-2 md:col-span-1">
    <div>
      <span class="font-mono text-xs font-bold text-zinc-500 uppercase tracking-wider block">Total Requests</span>
      <span class="font-display font-black text-3xl text-neubrutalism-charcoal mt-1 block">{totalTicketsCount}</span>
    </div>
    <Wrench class="absolute right-4 bottom-4 w-10 h-10 text-purple-500 opacity-20" />
  </div>
</div>
