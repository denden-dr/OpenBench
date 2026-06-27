<script lang="ts">
  import { Card } from '$lib';
  import { Wrench, ShieldCheck } from 'lucide-svelte';
  import type { PublicTrackerTicket } from '$lib/services/ticket';

  let {
    searchedTicket,
    getStatusColor,
    getStatusText
  }: {
    searchedTicket: PublicTrackerTicket;
    getStatusColor: (statusVal: string) => string;
    getStatusText: (statusVal: string) => string;
  } = $props();

  let displayStatus = $derived(
    searchedTicket.status === 'completed' && !searchedTicket.picked_up_at
      ? 'ready_for_pickup'
      : searchedTicket.status
  );
</script>

<Card bgColor="bg-white" class="border-4 border-neubrutalism-charcoal p-6 shadow-neubrutalism-md flex flex-col gap-5">
  <!-- Progress Header -->
  <div class="flex justify-between items-start border-b-2 border-dashed border-zinc-200 pb-4">
    <div>
      <span class="font-mono text-[9px] font-bold bg-zinc-100 text-zinc-500 border border-zinc-300 rounded px-1.5 py-0.5 uppercase">
        TICKET FOUND
      </span>
      <h2 class="font-display font-black text-xl text-neubrutalism-charcoal mt-1.5 uppercase leading-tight">
        {searchedTicket.brand_phone} {searchedTicket.model_phone}
      </h2>
    </div>

    <div class="flex flex-col items-end gap-1.5">
      <span class="font-mono text-xs font-extrabold border-2 border-neubrutalism-charcoal py-1 px-3 shadow-neubrutalism-sm uppercase {getStatusColor(displayStatus)}">
        {getStatusText(displayStatus)}
      </span>
    </div>
  </div>

  <!-- Progress Steps Tracker -->
  <div class="grid grid-cols-2 md:grid-cols-4 gap-3 text-center">
    <div class="border-2 border-neubrutalism-charcoal p-2.5 flex flex-col gap-1 rounded {['received','in_repair','ready_for_pickup','completed'].includes(displayStatus) ? 'bg-zinc-100' : 'opacity-40'}">
      <span class="font-mono text-[8px] font-bold text-zinc-500">STEP 1</span>
      <span class="font-display font-extrabold text-[10px] uppercase">RECEIVED</span>
    </div>

    <div class="border-2 border-neubrutalism-charcoal p-2.5 flex flex-col gap-1 rounded {['in_repair','ready_for_pickup','completed'].includes(displayStatus) ? 'bg-neubrutalism-yellow text-neubrutalism-charcoal font-bold' : 'opacity-40'}">
      <span class="font-mono text-[8px] {['in_repair','ready_for_pickup','completed'].includes(displayStatus) ? 'text-zinc-700' : 'text-zinc-500'} font-bold">STEP 2</span>
      <span class="font-display font-extrabold text-[10px] uppercase">IN REPAIR</span>
    </div>

    <div class="border-2 border-neubrutalism-charcoal p-2.5 flex flex-col gap-1 rounded {['ready_for_pickup','completed'].includes(displayStatus) ? 'bg-neubrutalism-green' : 'opacity-40'}">
      <span class="font-mono text-[8px] font-bold text-zinc-650">STEP 3</span>
      <span class="font-display font-extrabold text-[10px] uppercase">READY</span>
    </div>

    <div class="border-2 border-neubrutalism-charcoal p-2.5 flex flex-col gap-1 rounded {['completed'].includes(displayStatus) ? 'bg-zinc-100' : 'opacity-40'}">
      <span class="font-mono text-[8px] font-bold text-zinc-650">STEP 4</span>
      <span class="font-display font-extrabold text-[10px] uppercase">COMPLETED</span>
    </div>
  </div>

  <!-- Technical Description / Log -->
  <div class="bg-zinc-50 border-2 border-neubrutalism-charcoal p-4 font-mono text-xs flex flex-col gap-2">
    <div class="flex items-center gap-1 text-zinc-500 border-b border-dashed border-zinc-200 pb-1.5">
      <Wrench class="w-3.5 h-3.5" />
      <span class="font-bold uppercase text-[10px]">Repair Action / Technical Diagnosis:</span>
    </div>
    <p class="text-zinc-700 whitespace-pre-line leading-relaxed font-sans">
      {searchedTicket.repair_action || 'Technician is performing a thorough check on the circuit modules and supporting devices.'}
    </p>
  </div>

  <!-- Warranty Alert (If completed) -->
  {#if searchedTicket.status === 'completed'}
    <div class="bg-emerald-50 border-2 border-emerald-400 p-4 font-mono text-xs flex items-start gap-2.5 text-emerald-800">
      <ShieldCheck class="w-5 h-5 text-emerald-600 shrink-0 mt-0.5" />
      <div class="flex flex-col gap-0.5">
        <span class="font-bold uppercase text-[10px]">Device Warranty Active</span>
        <p class="font-sans text-[11px] mt-0.5 leading-snug">
          Your repair warranty is valid until <strong>{new Date(searchedTicket.warranty_expiry_date || '').toLocaleDateString('en-US')}</strong>. Present this electronic receipt if you need to make a warranty claim.
        </p>
      </div>
    </div>
  {/if}

  <!-- General Footer Metadata -->
  <div class="flex flex-wrap items-center justify-between gap-3 font-mono text-[10px] text-zinc-500">
    <span>Tracking Number: {searchedTicket.ticket_number}</span>
    <span>Date Received: {new Date(searchedTicket.created_at).toLocaleDateString('en-US')}</span>
  </div>
</Card>
