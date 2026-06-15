<script lang="ts">
  import { Card } from '$lib';
  import { Package, AlertTriangle } from 'lucide-svelte';
  import type { Product } from '$lib/services/inventory';

  let { products }: { products: Product[] } = $props();

  let totalItems = $derived(products.length);
  let lowStockCount = $derived(products.filter(p => p.stock <= p.min_stock).length);
</script>

<div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-4">
  <Card bgColor="bg-white" class="p-4 flex items-center justify-between border-4 border-neubrutalism-charcoal shadow-neubrutalism-sm">
    <div class="flex flex-col">
      <span class="font-mono text-[10px] font-bold text-zinc-500 uppercase">Total Items</span>
      <span class="font-mono text-2xl font-extrabold">{totalItems}</span>
    </div>
    <Package class="w-8 h-8 text-zinc-400" />
  </Card>

  <Card bgColor="bg-white" class="p-4 flex items-center justify-between border-4 border-neubrutalism-charcoal shadow-neubrutalism-sm">
    <div class="flex flex-col">
      <span class="font-mono text-[10px] font-bold text-rose-600 uppercase font-semibold">Low Stock Warnings</span>
      <span class="font-mono text-2xl font-extrabold text-rose-600">{lowStockCount}</span>
    </div>
    <AlertTriangle class="w-8 h-8 text-rose-500 animate-pulse" />
  </Card>
</div>
