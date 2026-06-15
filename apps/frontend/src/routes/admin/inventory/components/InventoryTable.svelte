<script lang="ts">
  import { Trash2, Edit2 } from 'lucide-svelte';
  import type { Product } from '$lib/services/inventory';

  let {
    isLoading,
    filteredProducts,
    onAdjustStock,
    onStartEdit,
    onDelete
  }: {
    isLoading: boolean;
    filteredProducts: Product[];
    onAdjustStock: (id: string, currentStock: number, amt: number) => void;
    onStartEdit: (p: Product) => void;
    onDelete: (id: string) => void;
  } = $props();
</script>

{#if isLoading}
  <div class="flex flex-col gap-3">
    {#each Array(4) as _}
      <div class="h-14 bg-zinc-200 border-4 border-neubrutalism-charcoal animate-pulse"></div>
    {/each}
  </div>
{:else}
  <div class="overflow-x-auto border-4 border-neubrutalism-charcoal shadow-neubrutalism-md">
    <table class="w-full bg-white text-left font-sans text-xs border-collapse">
      <thead class="bg-zinc-100 font-display font-extrabold uppercase border-b-4 border-neubrutalism-charcoal text-[10px] sm:text-xs">
        <tr>
          <th class="p-3 border-r-2 border-neubrutalism-charcoal">Item Details</th>
          <th class="p-3 border-r-2 border-neubrutalism-charcoal">Category</th>
          <th class="p-3 border-r-2 border-neubrutalism-charcoal text-center">Current Stock</th>
          <th class="p-3 border-r-2 border-neubrutalism-charcoal">Capital (Cost)</th>
          <th class="p-3 border-r-2 border-neubrutalism-charcoal">Sell Price</th>
          <th class="p-3 border-r-2 border-neubrutalism-charcoal">Margin Profit</th>
          <th class="p-3 text-center">Actions</th>
        </tr>
      </thead>
      <tbody class="font-mono">
        {#if filteredProducts.length === 0}
          <tr>
            <td colspan="7" class="p-8 text-center text-zinc-500 font-sans text-sm">
              No catalog items match current filter/search.
            </td>
          </tr>
        {:else}
          {#each filteredProducts as p (p.id)}
            <tr class="border-b-2 border-neubrutalism-charcoal hover:bg-zinc-50 transition-colors">
              <!-- Details -->
              <td class="p-3 border-r-2 border-neubrutalism-charcoal font-sans">
                <div class="flex flex-col">
                  <span class="font-bold text-sm text-neubrutalism-charcoal">{p.name}</span>
                  <span class="text-[10px] text-zinc-500 font-mono mt-0.5">ID: {p.id}</span>
                </div>
              </td>
              
              <!-- Category -->
              <td class="p-3 border-r-2 border-neubrutalism-charcoal font-mono text-[10px] uppercase font-bold">
                {#if p.category === 'retail'}
                  <span class="bg-neubrutalism-yellow border border-neubrutalism-charcoal py-0.5 px-1.5 shadow-neubrutalism-sm">RETAIL</span>
                {:else}
                  <span class="bg-neubrutalism-green border border-neubrutalism-charcoal py-0.5 px-1.5 shadow-neubrutalism-sm">SPARE PART</span>
                {/if}
              </td>

              <!-- Stock Counter -->
              <td class="p-3 border-r-2 border-neubrutalism-charcoal text-center">
                <div class="flex items-center justify-center gap-2">
                  <button 
                    class="px-1.5 py-0.5 border border-neubrutalism-charcoal bg-white font-bold hover:bg-zinc-100 active:translate-y-0.5 shadow-neubrutalism-sm text-[10px]"
                    onclick={() => onAdjustStock(p.id, p.stock, -1)}
                    title="Decrease stock by 1"
                  >-</button>
                  
                  <span class="font-bold text-sm w-8 inline-block {p.stock <= p.min_stock ? 'text-rose-600 bg-rose-50 border-2 border-rose-300 py-0.5 px-1 font-bold animate-pulse' : ''}">
                    {p.stock}
                  </span>

                  <button 
                    class="px-1.5 py-0.5 border border-neubrutalism-charcoal bg-white font-bold hover:bg-zinc-100 active:translate-y-0.5 shadow-neubrutalism-sm text-[10px]"
                    onclick={() => onAdjustStock(p.id, p.stock, 1)}
                    title="Increase stock by 1"
                  >+</button>
                </div>
                {#if p.stock <= p.min_stock}
                  <span class="text-[8px] text-rose-500 font-bold block mt-1 uppercase">Low (Min: {p.min_stock})</span>
                {/if}
              </td>

              <!-- Cost Price -->
              <td class="p-3 border-r-2 border-neubrutalism-charcoal">
                IDR {p.cost_price.toLocaleString('id-ID')}
              </td>

              <!-- Sell Price -->
              <td class="p-3 border-r-2 border-neubrutalism-charcoal">
                IDR {p.price.toLocaleString('id-ID')}
              </td>

              <!-- Margin -->
              <td class="p-3 border-r-2 border-neubrutalism-charcoal text-neubrutalism-pink font-bold">
                IDR {(p.price - p.cost_price).toLocaleString('id-ID')}
              </td>

              <!-- Actions -->
              <td class="p-3 text-center">
                <div class="flex items-center justify-center gap-2">
                  <button 
                    class="p-1.5 border-2 border-neubrutalism-charcoal bg-white hover:bg-zinc-100 shadow-neubrutalism-sm active:translate-y-0.5 active:shadow-none"
                    onclick={() => onStartEdit(p)}
                    title="Edit Item Details"
                  >
                    <Edit2 class="w-3.5 h-3.5" />
                  </button>
                  <button 
                    class="p-1.5 border-2 border-neubrutalism-charcoal bg-white hover:bg-rose-100 text-rose-600 shadow-neubrutalism-sm active:translate-y-0.5 active:shadow-none"
                    onclick={() => onDelete(p.id)}
                    title="Delete Item"
                  >
                    <Trash2 class="w-3.5 h-3.5" />
                  </button>
                </div>
              </td>
            </tr>
          {/each}
        {/if}
      </tbody>
    </table>
  </div>
{/if}
