<script lang="ts">
  import { Calendar, Printer } from 'lucide-svelte';
  import type { Sale } from '$lib/services/sales';

  let {
    filteredSales,
    isLoading,
    onSelectSale,
    formatDate
  }: {
    filteredSales: Sale[];
    isLoading: boolean;
    onSelectSale: (s: Sale) => void;
    formatDate: (dateStr: string) => string;
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
          <th class="p-3 border-r-2 border-neubrutalism-charcoal">Invoice Info</th>
          <th class="p-3 border-r-2 border-neubrutalism-charcoal">Transaction Date</th>
          <th class="p-3 border-r-2 border-neubrutalism-charcoal">Purchased Items</th>
          <th class="p-3 border-r-2 border-neubrutalism-charcoal">Net Total</th>
          <th class="p-3 border-r-2 border-neubrutalism-charcoal">Payment</th>
          <th class="p-3 text-center">Receipt</th>
        </tr>
      </thead>
      <tbody class="font-mono">
        {#if filteredSales.length === 0}
          <tr>
            <td colspan="6" class="p-8 text-center text-zinc-500 font-sans text-sm">
              No transactions match current filters.
            </td>
          </tr>
        {:else}
          {#each filteredSales as s}
            <tr class="border-b-2 border-neubrutalism-charcoal hover:bg-zinc-50 transition-colors">
              <!-- Invoice Info -->
              <td class="p-3 border-r-2 border-neubrutalism-charcoal font-bold">
                {s.invoice_number}
              </td>

              <!-- Date -->
              <td class="p-3 border-r-2 border-neubrutalism-charcoal font-sans flex items-center gap-1.5 mt-2">
                <Calendar class="w-3.5 h-3.5 text-zinc-500" />
                <span>{formatDate(s.created_at)}</span>
              </td>

              <!-- Items list -->
              <td class="p-3 border-r-2 border-neubrutalism-charcoal font-sans text-zinc-700">
                <div class="flex flex-col gap-0.5">
                  {#each s.items as item}
                    <span>{item.name} <strong class="font-mono">x{item.qty}</strong></span>
                  {/each}
                </div>
              </td>

              <!-- Net Total -->
              <td class="p-3 border-r-2 border-neubrutalism-charcoal font-extrabold text-neubrutalism-pink text-xs sm:text-sm">
                IDR {s.total.toLocaleString('id-ID')}
              </td>

              <!-- Payment Method -->
              <td class="p-3 border-r-2 border-neubrutalism-charcoal uppercase text-[10px] font-bold">
                {#if s.payment_method === 'cash'}
                  <span class="bg-neubrutalism-yellow border border-neubrutalism-charcoal py-0.5 px-1.5 shadow-neubrutalism-sm">CASH</span>
                {:else}
                  <span class="bg-neubrutalism-green border border-neubrutalism-charcoal py-0.5 px-1.5 shadow-neubrutalism-sm">QRIS</span>
                {/if}
              </td>

              <!-- Actions (re-print) -->
              <td class="p-3 text-center">
                <button 
                  class="p-1.5 border-2 border-neubrutalism-charcoal bg-white hover:bg-zinc-100 shadow-neubrutalism-sm active:translate-y-0.5 active:shadow-none cursor-pointer"
                  onclick={() => onSelectSale(s)}
                  title="View & Re-print Receipt"
                >
                  <Printer class="w-4 h-4 text-neubrutalism-charcoal" />
                </button>
              </td>
            </tr>
          {/each}
        {/if}
      </tbody>
    </table>
  </div>
{/if}
