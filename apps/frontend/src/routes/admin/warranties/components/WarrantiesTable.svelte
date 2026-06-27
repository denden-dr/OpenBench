<script lang="ts">
  import { User, ArrowUpRight } from 'lucide-svelte';
  import type { Warranty } from '$lib/services/warranty';

  interface Props {
    filteredWarranties: Warranty[];
    isLoading: boolean;
  }

  let { filteredWarranties, isLoading }: Props = $props();

  const getWarrantyStatus = (endDateStr: string): 'active' | 'expired' => {
    const end = new Date(endDateStr);
    return end > new Date() ? 'active' : 'expired';
  };

  const getRemainingDays = (endDateStr: string): number => {
    const end = new Date(endDateStr).getTime();
    const now = new Date().getTime();
    const diffTime = end - now;
    if (diffTime <= 0) return 0;
    return Math.ceil(diffTime / (1000 * 60 * 60 * 24));
  };

  const formatDate = (dateStr: string) => {
    const d = new Date(dateStr);
    return d.toLocaleDateString('en-US', { day: 'numeric', month: 'short', year: 'numeric' });
  };
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
          <th class="p-3 border-r-2 border-neubrutalism-charcoal">Ticket Number</th>
          <th class="p-3 border-r-2 border-neubrutalism-charcoal">Customer Name</th>
          <th class="p-3 border-r-2 border-neubrutalism-charcoal">Device Info</th>
          <th class="p-3 border-r-2 border-neubrutalism-charcoal">Guarantee Period</th>
          <th class="p-3 border-r-2 border-neubrutalism-charcoal text-center">Warranty Status</th>
          <th class="p-3 text-center">Actions</th>
        </tr>
      </thead>
      <tbody class="font-mono">
        {#if filteredWarranties.length === 0}
          <tr>
            <td colspan="6" class="p-8 text-center text-zinc-500 font-sans text-sm">
              No active or expired warranties match current search filters.
            </td>
          </tr>
        {:else}
          {#each filteredWarranties as w (w.id)}
            {@const status = getWarrantyStatus(w.end_date)}
            {@const daysLeft = getRemainingDays(w.end_date)}
            <tr class="border-b-2 border-neubrutalism-charcoal hover:bg-zinc-50 transition-colors">
              <!-- Ticket Number -->
              <td class="p-3 border-r-2 border-neubrutalism-charcoal font-bold">
                {w.ticket_number}
              </td>

              <!-- Customer -->
              <td class="p-3 border-r-2 border-neubrutalism-charcoal font-sans text-zinc-700">
                <div class="flex items-center gap-1">
                  <User class="w-3.5 h-3.5 text-zinc-400" />
                  <span class="font-bold">{w.customer_name}</span>
                </div>
              </td>

              <!-- Device Info -->
              <td class="p-3 border-r-2 border-neubrutalism-charcoal font-sans text-zinc-800">
                {w.device_info}
              </td>

              <!-- Period -->
              <td class="p-3 border-r-2 border-neubrutalism-charcoal font-sans text-zinc-600">
                <div class="flex flex-col gap-0.5">
                  <span>Starts: {formatDate(w.start_date)}</span>
                  <span class="font-bold">Expires: {formatDate(w.end_date)}</span>
                </div>
              </td>

              <!-- Status & Time Remaining -->
              <td class="p-3 border-r-2 border-neubrutalism-charcoal text-center">
                <div class="flex flex-col items-center gap-1">
                  {#if status === 'active'}
                    <span class="bg-neubrutalism-green border border-neubrutalism-charcoal py-0.5 px-2 font-mono text-[9px] font-bold shadow-neubrutalism-sm">
                      ACTIVE
                    </span>
                    <span class="text-[8px] text-emerald-600 font-bold uppercase">{daysLeft} days remaining</span>
                  {:else}
                    <span class="bg-zinc-200 border border-zinc-400 text-zinc-500 py-0.5 px-2 font-mono text-[9px] font-bold">
                      EXPIRED
                    </span>
                    <span class="text-[8px] text-zinc-400 font-bold uppercase">Guarantee Ended</span>
                  {/if}
                </div>
              </td>

              <!-- View Ticket Link -->
              <td class="p-3 text-center">
                <a href="/admin/tickets/{w.ticket_id}">
                  <button 
                    class="p-1.5 border-2 border-neubrutalism-charcoal bg-white hover:bg-zinc-100 shadow-neubrutalism-sm active:translate-y-0.5 active:shadow-none cursor-pointer"
                    title="View Associated Repair Ticket"
                  >
                    <ArrowUpRight class="w-4 h-4 text-neubrutalism-charcoal" />
                  </button>
                </a>
              </td>
            </tr>
          {/each}
        {/if}
      </tbody>
    </table>
  </div>
{/if}
