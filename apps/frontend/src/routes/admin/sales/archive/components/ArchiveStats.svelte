<script lang="ts">
  import { Card } from '$lib';
  import { DollarSign, FileText } from 'lucide-svelte';

  let {
    totalRevenue,
    salesCount
  }: {
    totalRevenue: number;
    salesCount: number;
  } = $props();
</script>

<div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-4">
  <!-- Gross Revenue Card -->
  <Card bgColor="bg-white" class="p-4 flex items-center justify-between border-4 border-neubrutalism-charcoal shadow-neubrutalism-sm">
    <div class="flex flex-col">
      <span class="font-mono text-[10px] font-bold text-zinc-500 uppercase">Gross POS Revenue</span>
      <span class="font-mono text-2xl font-extrabold text-neubrutalism-pink">IDR {totalRevenue.toLocaleString('id-ID')}</span>
    </div>
    <DollarSign class="w-8 h-8 text-neubrutalism-pink" />
  </Card>

  <!-- Transaction Count Card -->
  <Card bgColor="bg-white" class="p-4 flex items-center justify-between border-4 border-neubrutalism-charcoal shadow-neubrutalism-sm">
    <div class="flex flex-col">
      <span class="font-mono text-[10px] font-bold text-zinc-500 uppercase font-semibold">Total Completed Sales</span>
      <span class="font-mono text-2xl font-extrabold">{salesCount}</span>
    </div>
    <FileText class="w-8 h-8 text-zinc-400" />
  </Card>
</div>
