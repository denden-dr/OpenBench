<script lang="ts">
  import { Search } from 'lucide-svelte';
  import type { Product } from '$lib/services/inventory';

  let {
    searchQuery = $bindable(''),
    isLoading,
    filteredProducts,
    onAddToCart
  }: {
    searchQuery: string;
    isLoading: boolean;
    filteredProducts: Product[];
    onAddToCart: (p: Product) => void;
  } = $props();
</script>

<div class="flex flex-col gap-4">
  <div class="flex items-center gap-4">
    <!-- Search bar -->
    <div class="relative w-full">
      <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none text-zinc-500">
        <Search class="w-4 h-4" />
      </div>
      <input 
        type="text" 
        placeholder="Search retail items or spare parts..."
        bind:value={searchQuery}
        class="w-full pl-9 pr-4 py-2 border-4 border-neubrutalism-charcoal bg-white focus:outline-none focus:bg-zinc-50 font-mono text-sm shadow-neubrutalism-sm"
      />
    </div>
  </div>

  {#if isLoading}
    <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
      {#each Array(4) as _}
        <div class="h-28 bg-zinc-200 border-4 border-neubrutalism-charcoal animate-pulse"></div>
      {/each}
    </div>
  {:else}
    <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
      {#if filteredProducts.length === 0}
        <div class="col-span-2 p-12 text-center text-zinc-500 font-mono text-xs border-4 border-dashed border-zinc-300">
          PRODUCT NOT FOUND.
        </div>
      {:else}
        {#each filteredProducts as p (p.id)}
          <button 
            class="text-left bg-white border-4 border-neubrutalism-charcoal p-4 hover:shadow-neubrutalism-md transition-all active:translate-x-0.5 active:translate-y-0.5 active:shadow-none flex flex-col justify-between h-32 group cursor-pointer relative"
            onclick={() => onAddToCart(p)}
            disabled={p.stock <= 0}
          >
            <!-- Stock Badge -->
            <span class="absolute top-2 right-2 font-mono text-[9px] font-bold border-2 border-neubrutalism-charcoal px-1.5 py-0.5 bg-zinc-100 uppercase {p.stock <= p.min_stock ? 'bg-rose-100 text-rose-600 font-bold' : ''}">
              Stock: {p.stock}
            </span>

            <div class="flex flex-col gap-1 pr-14">
              <!-- Category Badge -->
              <span class="font-mono text-[8px] font-bold uppercase w-fit tracking-wider px-1 bg-zinc-150 border border-zinc-300 rounded text-zinc-600">
                {p.category.replace('_', ' ')}
              </span>
              
              <h4 class="font-display font-extrabold text-sm text-neubrutalism-charcoal leading-tight line-clamp-2 mt-1">
                {p.name}
              </h4>
            </div>

            <div class="flex justify-between items-center w-full border-t border-dashed border-zinc-200 pt-2 mt-2">
              <span class="font-mono text-xs font-extrabold text-neubrutalism-pink">
                IDR {p.price.toLocaleString('id-ID')}
              </span>
              
              {#if p.stock <= 0}
                <span class="font-mono text-[9px] font-bold text-rose-500 uppercase">SOLDOUT</span>
              {:else}
                <span class="font-mono text-[9px] font-bold bg-neubrutalism-green border border-neubrutalism-charcoal py-0.5 px-1.5 shadow-neubrutalism-sm group-hover:bg-neubrutalism-yellow transition-colors">
                  + ADD TO CART
                </span>
              {/if}
            </div>
          </button>
        {/each}
      {/if}
    </div>
  {/if}
</div>
