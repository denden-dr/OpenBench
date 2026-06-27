<script lang="ts">
  import { Card } from '$lib';
  import { Calendar, User, Phone, ArrowUpRight } from 'lucide-svelte';
  import type { Ticket } from '$lib/services/ticket';

  interface Props {
    filteredTickets: Ticket[];
    isLoading: boolean;
  }

  let { filteredTickets, isLoading }: Props = $props();

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'ready_for_pickup': return 'bg-neubrutalism-green';
      case 'in_repair': return 'bg-neubrutalism-yellow';
      case 'received': return 'bg-zinc-200';
      case 'completed': return 'bg-zinc-100 text-zinc-500';
      case 'cancelled': return 'bg-rose-100 text-rose-600 border-rose-300';
      default: return 'bg-white';
    }
  };

  const getStatusText = (status: string) => {
    return status.replace(/_/g, ' ');
  };

  const formatDate = (dateStr: string) => {
    const d = new Date(dateStr);
    return d.toLocaleDateString('en-US', { day: 'numeric', month: 'short', year: 'numeric' });
  };

  const getDisplayStatus = (t: Ticket) => {
    return t.status === 'completed' && t.device_position === 'warehouse'
      ? 'ready_for_pickup'
      : t.status;
  };
</script>

{#if isLoading}
  <div class="flex flex-col gap-4">
    {#each Array(4) as _}
      <div class="h-24 bg-zinc-200 border-4 border-neubrutalism-charcoal animate-pulse"></div>
    {/each}
  </div>
{:else}
  <div class="flex flex-col gap-4">
    {#if filteredTickets.length === 0}
      <Card bgColor="bg-white" class="p-12 text-center font-mono text-sm text-zinc-500">
        NO TICKETS MATCH THE CURRENT SEARCH / FILTER.
      </Card>
    {:else}
      {#each filteredTickets as ticket (ticket.id)}
        <Card bgColor="bg-white" class="hover:shadow-neubrutalism-lg transition-all border-4 border-neubrutalism-charcoal">
          <div class="flex flex-col md:flex-row md:items-center justify-between gap-4">
            
            <!-- Left side info -->
            <div class="flex-grow flex flex-col sm:flex-row sm:items-start gap-4">
              
              <div class="flex flex-col gap-1 shrink-0">
                <span class="font-mono text-xs font-extrabold bg-zinc-200 px-2 py-0.5 border-2 border-neubrutalism-charcoal text-neubrutalism-charcoal shadow-neubrutalism-sm w-fit">
                  {ticket.ticket_number}
                </span>
                
                <div class="flex items-center gap-1.5 font-mono text-[10px] text-zinc-500 mt-1">
                  <Calendar class="w-3.5 h-3.5" />
                  <span>{formatDate(ticket.created_at)}</span>
                </div>
              </div>

              <div class="flex flex-col gap-1">
                <h3 class="font-display font-bold text-lg leading-tight">
                  {ticket.brand_phone} {ticket.model_phone}
                </h3>
                
                <div class="flex flex-wrap items-center gap-x-4 gap-y-1 font-sans text-xs text-zinc-700 mt-1">
                  <span class="flex items-center gap-1">
                    <User class="w-3.5 h-3.5 text-zinc-500" />
                    <strong>{ticket.customer_name}</strong>
                  </span>
                  <span class="flex items-center gap-1">
                    <Phone class="w-3.5 h-3.5 text-zinc-500" />
                    <span>{ticket.customer_phone}</span>
                  </span>
                </div>

                <p class="font-sans text-xs text-neubrutalism-charcoal mt-1 line-clamp-1">
                  <span class="font-semibold text-zinc-600">Damage:</span> {ticket.damage_description}
                </p>
              </div>

            </div>

            <!-- Right side actions & status -->
            <div class="flex flex-row md:flex-col items-center md:items-end justify-between md:justify-center gap-3 shrink-0 border-t-2 md:border-t-0 border-dashed border-zinc-200 pt-3 md:pt-0">
              <div class="flex items-center gap-2">
                <span class="font-mono text-xs font-bold py-1 px-3 border-2 border-neubrutalism-charcoal shadow-neubrutalism-sm uppercase tracking-wide {getStatusColor(getDisplayStatus(ticket))}">
                  {getStatusText(getDisplayStatus(ticket))}
                </span>
              </div>

              <div class="flex items-center gap-3">
                <span class="font-mono text-sm font-extrabold">
                  IDR {ticket.cost.toLocaleString('id-ID')}
                </span>
                
                <a href="/admin/tickets/{ticket.id}">
                  <button class="p-2 border-2 border-neubrutalism-charcoal hover:bg-zinc-100 transition duration-150 cursor-pointer shadow-neubrutalism-sm active:translate-x-0.5 active:translate-y-0.5 active:shadow-none bg-white">
                    <ArrowUpRight class="w-4 h-4 text-neubrutalism-charcoal" />
                  </button>
                </a>
              </div>
            </div>

          </div>
        </Card>
      {/each}
    {/if}
  </div>
{/if}
