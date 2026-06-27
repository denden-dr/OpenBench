<script lang="ts">
  import { Card } from '$lib';
  import { ShieldCheck, Award, Clock } from 'lucide-svelte';

  interface Props {
    activeWarrantiesCount: number;
    expiredWarrantiesCount: number;
  }

  let { activeWarrantiesCount, expiredWarrantiesCount }: Props = $props();
</script>

<!-- Header Area -->
<div class="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
  <div>
    <h2 class="font-display font-extrabold text-2xl md:text-3xl uppercase tracking-tight flex items-center gap-2">
      <ShieldCheck class="w-8 h-8 text-neubrutalism-green" />
      <span>Warranty Guarantee Monitor</span>
    </h2>
    <p class="font-sans text-xs sm:text-sm text-zinc-500 font-semibold mt-1">
      Monitor active device warranties registered upon successful repair pickups.
    </p>
  </div>
</div>

<!-- Stats Grid -->
<div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-4">
  <!-- Active Warranties -->
  <Card bgColor="bg-white" class="p-4 flex items-center justify-between border-4 border-neubrutalism-charcoal shadow-neubrutalism-sm">
    <div class="flex flex-col">
      <span class="font-mono text-[10px] font-bold text-zinc-500 uppercase">Active Guarantees</span>
      <span class="font-mono text-2xl font-extrabold text-neubrutalism-green">{activeWarrantiesCount}</span>
    </div>
    <Award class="w-8 h-8 text-neubrutalism-green" />
  </Card>

  <!-- Expired Warranties -->
  <Card bgColor="bg-white" class="p-4 flex items-center justify-between border-4 border-neubrutalism-charcoal shadow-neubrutalism-sm">
    <div class="flex flex-col">
      <span class="font-mono text-[10px] font-bold text-zinc-500 uppercase font-semibold">Expired Guarantees</span>
      <span class="font-mono text-2xl font-extrabold text-zinc-400">{expiredWarrantiesCount}</span>
    </div>
    <Clock class="w-8 h-8 text-zinc-300" />
  </Card>
</div>
