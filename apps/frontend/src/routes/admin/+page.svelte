<script lang="ts">
  import { Card, Button } from '$lib';
  import { ticketService, type Ticket } from '$lib/services/ticket';
  import { onMount } from 'svelte';
  import { 
    Wrench, CheckCircle2, Clock, ClipboardList, 
    AlertTriangle, ArrowUpRight, ArrowRight 
  } from 'lucide-svelte';

  let tickets = $state<Ticket[]>([]);
  let isLoading = $state(true);

  // Compute stats reactively using Svelte 5 derived states
  let activeTicketsCount = $derived(tickets.filter(t => t.status !== 'picked_up' && t.status !== 'cancelled').length);
  let completedTodayCount = $derived(tickets.filter(t => t.status === 'picked_up').length);
  let diagnosingCount = $derived(tickets.filter(t => t.status === 'diagnosing').length);

  onMount(async () => {
    tickets = await ticketService.getTickets();
    isLoading = false;
  });

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'ready_for_pickup': return 'bg-neubrutalism-green';
      case 'in_repair': return 'bg-neubrutalism-yellow';
      case 'diagnosing': return 'bg-neubrutalism-pink text-white';
      case 'received': return 'bg-zinc-200';
      case 'picked_up': return 'bg-zinc-100 text-zinc-500 line-through';
      default: return 'bg-white';
    }
  };

  const getStatusText = (status: string) => {
    return status.replace(/_/g, ' ');
  };
</script>

<svelte:head>
  <title>Overview - OpenBench Admin</title>
</svelte:head>

<div class="flex flex-col gap-8">
  <!-- Welcome banner -->
  <div class="border-4 border-neubrutalism-charcoal bg-neubrutalism-yellow p-6 shadow-neubrutalism-md">
    <h1 class="font-display font-bold text-2xl md:text-3xl uppercase leading-none mb-2">
      Welcome to the Workbench, Admin!
    </h1>
    <p class="font-sans text-xs md:text-sm max-w-2xl opacity-90 font-semibold">
      This console allows you to manage repair tickets, Point of Sales transaction logs, inventory stock, and warranty periods. Currently running in Sandbox Mock Mode.
    </p>
  </div>

  {#if isLoading}
    <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
      {#each Array(3) as _}
        <div class="h-32 bg-zinc-200 border-4 border-neubrutalism-charcoal animate-pulse"></div>
      {/each}
    </div>
  {:else}
    <!-- Metrics Cards Grid -->
    <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
      
      <!-- Metric 1 -->
      <Card bgColor="bg-white" class="flex flex-col gap-2 relative">
        <div class="absolute top-4 right-4 bg-neubrutalism-charcoal text-white p-2 border-2 border-neubrutalism-charcoal">
          <Wrench class="w-5 h-5" />
        </div>
        <span class="font-mono text-xs font-bold text-neubrutalism-charcoal opacity-70 uppercase tracking-widest">Active Tickets</span>
        <span class="font-mono text-5xl font-extrabold text-neubrutalism-charcoal mt-2">{activeTicketsCount}</span>
        <span class="font-sans text-xs text-neubrutalism-charcoal opacity-60 mt-2 flex items-center gap-1">
          <AlertTriangle class="w-3.5 h-3.5 text-amber-500" />
          Active repairs in the queue
        </span>
      </Card>

      <!-- Metric 2 -->
      <Card bgColor="bg-white" class="flex flex-col gap-2 relative">
        <div class="absolute top-4 right-4 bg-neubrutalism-charcoal text-white p-2 border-2 border-neubrutalism-charcoal">
          <CheckCircle2 class="w-5 h-5" />
        </div>
        <span class="font-mono text-xs font-bold text-neubrutalism-charcoal opacity-70 uppercase tracking-widest">Completed Repairs</span>
        <span class="font-mono text-5xl font-extrabold text-neubrutalism-charcoal mt-2">{completedTodayCount}</span>
        <span class="font-sans text-xs text-neubrutalism-charcoal opacity-60 mt-2 flex items-center gap-1">
          <ArrowUpRight class="w-3.5 h-3.5 text-emerald-500" />
          Devices picked up by customers
        </span>
      </Card>

      <!-- Metric 3 -->
      <Card bgColor="bg-white" class="flex flex-col gap-2 relative">
        <div class="absolute top-4 right-4 bg-neubrutalism-charcoal text-white p-2 border-2 border-neubrutalism-charcoal">
          <Clock class="w-5 h-5" />
        </div>
        <span class="font-mono text-xs font-bold text-neubrutalism-charcoal opacity-70 uppercase tracking-widest">In Diagnosis</span>
        <span class="font-mono text-5xl font-extrabold text-neubrutalism-charcoal mt-2">{diagnosingCount}</span>
        <span class="font-sans text-xs text-neubrutalism-charcoal opacity-60 mt-2 flex items-center gap-1">
          <ClipboardList class="w-3.5 h-3.5 text-zinc-500" />
          Awaiting technical diagnostics
        </span>
      </Card>
      
    </div>

    <!-- Active Repairs Sandbox -->
    <div class="flex flex-col gap-6 col-span-2">
      <div class="flex justify-between items-center">
        <h2 class="font-display font-bold text-xl md:text-2xl text-neubrutalism-charcoal uppercase tracking-tight">
          Recent Active Repairs
        </h2>
        <a href="/admin/tickets" class="inline-block">
          <Button bgColor="bg-neubrutalism-yellow" class="py-1.5 px-3 text-xs font-bold flex items-center gap-1.5 shadow-neubrutalism-sm">
            <span>ALL TICKETS</span>
            <ArrowRight class="w-3.5 h-3.5" />
          </Button>
        </a>
      </div>
      
      <div class="flex flex-col gap-4">
        {#if tickets.filter(t => t.status !== 'picked_up').length === 0}
          <Card bgColor="bg-white" class="p-8 text-center font-mono text-sm text-zinc-500">
            NO ACTIVE REPAIR TICKETS IN QUEUE
          </Card>
        {:else}
          {#each tickets.filter(t => t.status !== 'picked_up').slice(0, 5) as ticket}
            <Card bgColor="bg-white" class="hover:shadow-neubrutalism-lg transition-all">
              <div class="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
                <div class="flex flex-col gap-1">
                  <div class="flex items-center gap-2">
                    <span class="font-mono text-xs font-bold bg-zinc-200 px-2 py-0.5 border-2 border-neubrutalism-charcoal text-neubrutalism-charcoal shadow-neubrutalism-sm">
                      {ticket.ticket_number}
                    </span>
                    <span class="font-display font-bold text-base md:text-lg">{ticket.brand_phone} {ticket.model_phone}</span>
                  </div>
                  <p class="font-sans text-xs sm:text-sm text-neubrutalism-charcoal mt-1">
                    <span class="font-semibold">Issue:</span> {ticket.damage_description}
                  </p>
                  <p class="font-sans text-xs text-neubrutalism-charcoal opacity-60">
                    Customer: {ticket.customer_name} &bull; Serial: {ticket.serial_number} &bull; Cost: IDR {ticket.cost.toLocaleString('id-ID')}
                  </p>
                </div>
                
                <div class="flex items-center gap-3 w-full sm:w-auto justify-between sm:justify-end">
                  <span class="font-mono text-[10px] sm:text-xs font-bold py-1 px-2.5 border-2 border-neubrutalism-charcoal shadow-neubrutalism-sm uppercase tracking-wide {getStatusColor(ticket.status)}">
                    {getStatusText(ticket.status)}
                  </span>
                  
                  <a href="/admin/tickets/{ticket.id}" class="inline-block">
                    <button class="p-2 border-2 border-neubrutalism-charcoal hover:bg-zinc-100 transition duration-150 cursor-pointer shadow-neubrutalism-sm active:translate-x-0.5 active:translate-y-0.5 active:shadow-none bg-white">
                      <ArrowUpRight class="w-4 h-4 text-neubrutalism-charcoal" />
                    </button>
                  </a>
                </div>
              </div>
            </Card>
          {/each}
        {/if}
      </div>
    </div>
  {/if}
</div>
