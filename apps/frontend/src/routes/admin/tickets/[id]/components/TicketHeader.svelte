<script lang="ts">
  import { Calendar } from 'lucide-svelte';
  import type { Ticket } from '$lib/services/ticket';

  let {
    ticket,
    getStatusColor,
    getStatusText,
    formatDate
  }: {
    ticket: Ticket;
    getStatusColor: (statusVal: string) => string;
    getStatusText: (statusVal: string) => string;
    formatDate: (dateStr?: string) => string;
  } = $props();

  let displayStatus = $derived(
    ticket.status === 'completed' && ticket.device_position === 'warehouse'
      ? 'ready_for_pickup'
      : ticket.status
  );
</script>

<div class="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
  <div>
    <div class="flex items-center gap-3">
      <span class="font-mono text-sm font-extrabold bg-zinc-200 px-2.5 py-1 border-2 border-neubrutalism-charcoal shadow-neubrutalism-sm">
        {ticket.ticket_number}
      </span>
      <span class="font-mono text-xs font-bold py-1 px-3 border-2 border-neubrutalism-charcoal shadow-neubrutalism-sm uppercase {getStatusColor(displayStatus)}">
        {getStatusText(displayStatus)}
      </span>
    </div>
    <h2 class="font-display font-extrabold text-2xl md:text-3xl uppercase tracking-tight mt-2">
      {ticket.brand_phone} {ticket.model_phone}
    </h2>
  </div>
  
  <div class="font-mono text-[10px] text-zinc-500 text-left sm:text-right shrink-0">
    <div class="flex items-center sm:justify-end gap-1">
      <Calendar class="w-3.5 h-3.5" />
      <span>Registered: {formatDate(ticket.created_at)}</span>
    </div>
  </div>
</div>
