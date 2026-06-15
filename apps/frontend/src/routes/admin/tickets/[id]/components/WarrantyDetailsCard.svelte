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

  let expiryDate = $derived(ticket.warranty_expiry_date || new Date(Date.now() + 30*24*60*60*1000).toISOString());
</script>

<Card bgColor="bg-neubrutalism-green/10" class="border-4 border-neubrutalism-green/50 p-6 flex flex-col gap-3">
  <h3 class="font-display font-bold text-sm uppercase text-emerald-800 border-b-2 border-dashed border-emerald-300 pb-2 flex items-center gap-1.5">
    <Award class="w-4 h-4 text-emerald-700" />
    <span>Warranty Details</span>
  </h3>
  
  <div class="font-sans text-xs text-emerald-800 flex flex-col gap-1 font-semibold">
    <p>Status: <span class="bg-emerald-600 text-white px-2 py-0.5 text-[9px] font-mono rounded">ACTIVE</span></p>
    <p>Starts: {formatDate(ticket.created_at)}</p>
    <p>Expires: {formatDate(expiryDate)}</p>
  </div>
</Card>
