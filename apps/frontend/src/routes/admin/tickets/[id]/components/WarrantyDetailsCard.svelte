<script lang="ts">
  import { Card } from '$lib';
  import { Award } from 'lucide-svelte';
  import type { Ticket } from '$lib/services/ticket';

  let {
    ticket,
    status,
    formatDate
  }: {
    ticket: Ticket;
    status: string;
    formatDate: (dateStr?: string) => string;
  } = $props();

  let hasWarranty = $derived(!!ticket.warranty_expiry_date);
</script>

<Card
  bgColor={hasWarranty ? "bg-neubrutalism-green/10" : "bg-neubrutalism-charcoal/5"}
  class={`border-4 p-6 flex flex-col gap-3 ${hasWarranty ? "border-neubrutalism-green/50" : "border-neubrutalism-charcoal/30"}`}
>
  <h3 class={`font-display font-bold text-sm uppercase border-b-2 border-dashed pb-2 flex items-center gap-1.5 ${hasWarranty ? "text-emerald-800 border-emerald-300" : "text-neubrutalism-charcoal border-neubrutalism-charcoal/20"}`}>
    <Award class={`w-4 h-4 ${hasWarranty ? "text-emerald-700" : "text-neubrutalism-charcoal/70"}`} />
    <span>Warranty Details</span>
  </h3>
  
  <div class={`font-sans text-xs flex flex-col gap-1 font-semibold ${hasWarranty ? "text-emerald-800" : "text-neubrutalism-charcoal"}`}>
    <p>
      Status:
      {#if hasWarranty}
        <span class="bg-emerald-600 text-white px-2 py-0.5 text-[9px] font-mono rounded font-bold">ACTIVE</span>
      {:else}
        <span class="bg-amber-100 text-amber-800 border border-amber-300 px-2 py-0.5 text-[9px] font-mono rounded font-bold">PENDING</span>
      {/if}
    </p>
    <p>Starts: {hasWarranty && ticket.picked_up_at ? formatDate(ticket.picked_up_at) : 'Not Activated'}</p>
    <p>Expires: {hasWarranty ? formatDate(ticket.warranty_expiry_date) : 'Not Activated'}</p>
  </div>
</Card>
