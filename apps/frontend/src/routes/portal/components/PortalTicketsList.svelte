<script lang="ts">
  import { Card, Button } from '$lib';
  import { Wrench, Plus, Calendar, ChevronRight } from 'lucide-svelte';
  import { goto } from '$app/navigation';
  import type { Ticket } from '$lib/services/ticket';

  interface Props {
    tickets: Ticket[];
  }

  let { tickets }: Props = $props();

  function formatDate(dateStr: string): string {
    try {
      const date = new Date(dateStr);
      return date.toLocaleDateString('en-US', {
        year: 'numeric',
        month: 'short',
        day: 'numeric'
      });
    } catch {
      return dateStr;
    }
  }

  function getStatusStyle(status: string) {
    switch (status) {
      case 'received':
        return 'bg-neubrutalism-yellow text-neubrutalism-charcoal';
      case 'in_repair':
        return 'bg-amber-300 text-neubrutalism-charcoal';
      case 'ready_for_pickup':
      case 'completed':
        return 'bg-emerald-300 text-neubrutalism-charcoal';
      case 'cancelled':
        return 'bg-rose-300 text-neubrutalism-charcoal';
      default:
        return 'bg-zinc-200 text-neubrutalism-charcoal';
    }
  }

  function getStatusLabel(status: string) {
    switch (status) {
      case 'received':
        return 'RECEIVED';
      case 'in_repair':
        return 'IN REPAIR';
      case 'ready_for_pickup':
        return 'READY FOR PICKUP';
      case 'completed':
        return 'COMPLETED';
      case 'cancelled':
        return 'CANCELLED';
      default:
        return status.toUpperCase();
    }
  }
</script>

<div class="flex flex-col gap-6">
  <h2 class="font-display font-black text-2xl text-neubrutalism-charcoal uppercase tracking-tight">
    My Repair Requests
  </h2>

  {#if tickets.length === 0}
    <!-- Empty State -->
    <Card class="p-12 bg-white text-center flex flex-col items-center justify-center border-4 border-neubrutalism-charcoal shadow-neubrutalism-md">
      <Wrench class="w-16 h-16 text-zinc-300 mb-4" />
      <h3 class="font-display font-bold text-xl text-neubrutalism-charcoal uppercase">No repair tickets found</h3>
      <p class="font-sans text-sm text-zinc-500 mt-2 max-w-sm">
        You haven't submitted any repair requests yet. Create a request for any device issues you are facing.
      </p>
      <Button 
        bgColor="bg-neubrutalism-pink" 
        class="mt-6 flex items-center gap-2 uppercase font-bold text-sm shadow-neubrutalism-sm text-white"
        onclick={() => goto('/repair-request')}
      >
        <Plus class="w-4 h-4" />
        <span>Submit First Request</span>
      </Button>
    </Card>
  {:else}
    <!-- Tickets Grid -->
    <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
      {#each tickets as ticket (ticket.id)}
        <div class="bg-white border-4 border-neubrutalism-charcoal p-6 shadow-neubrutalism-sm hover:shadow-neubrutalism-md hover:-translate-y-1 transition-all flex flex-col justify-between gap-6">
          <!-- Top details -->
          <div class="flex flex-col gap-4">
            <div class="flex justify-between items-start gap-4">
              <div>
                <span class="font-mono text-xs font-bold text-zinc-400 block uppercase tracking-wide">Ticket Number</span>
                <span class="font-mono text-sm font-black text-neubrutalism-charcoal">{ticket.ticket_number}</span>
              </div>
              <span class="inline-block border-2 border-neubrutalism-charcoal font-mono text-[10px] font-bold px-2 py-0.5 rounded-none uppercase {getStatusStyle(ticket.status)}">
                {getStatusLabel(ticket.status)}
              </span>
            </div>

            <div>
              <h3 class="font-display font-black text-xl text-neubrutalism-charcoal uppercase tracking-tight leading-tight">
                {ticket.brand_phone} {ticket.model_phone}
              </h3>
              <p class="font-sans text-xs text-zinc-500 mt-1 line-clamp-2 italic">
                "{ticket.damage_description}"
              </p>
            </div>
          </div>

          <!-- Bottom stats & CTA -->
          <div class="flex justify-between items-center border-t-2 border-dashed border-zinc-150 pt-4 mt-auto">
            <div class="flex items-center gap-1.5 font-mono text-xs text-zinc-500">
              <Calendar class="w-3.5 h-3.5" />
              <span>{formatDate(ticket.created_at)}</span>
            </div>

            <Button 
              bgColor="bg-white" 
              class="border-2 border-neubrutalism-charcoal text-xs font-bold py-1 px-3 flex items-center gap-1 hover:bg-neubrutalism-pink hover:text-white transition-all shadow-[2px_2px_0px_0px_#1e1e24]"
              onclick={() => goto(`/tracker?id=${ticket.id}`)}
            >
              <span>Track</span>
              <ChevronRight class="w-3.5 h-3.5" />
            </Button>
          </div>
        </div>
      {/each}
    </div>
  {/if}
</div>
