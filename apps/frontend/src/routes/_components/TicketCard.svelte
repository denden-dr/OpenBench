<script lang="ts">
  import { RefreshCw, ShieldCheck, LoaderCircle } from "lucide-svelte";
  import type { Ticket } from "$lib/types/ticket";
  import { formatCurrency } from "$lib/utils/format";
  import { checkWarrantyExpiry } from "$lib/utils/warranty";

  let {
    ticket,
    isActionLoading,
    onSelectTicket,
    onQuickAction,
    getStatusBadgeClass,
    getStatusLabel,
  }: {
    ticket: Ticket;
    isActionLoading: boolean;
    onSelectTicket: () => void;
    onQuickAction: () => void;
    getStatusBadgeClass: (status: string, isWarranty: boolean) => string;
    getStatusLabel: (status: string, isWarranty: boolean) => string;
  } = $props();
</script>

<div
  role="button"
  tabindex="0"
  onclick={onSelectTicket}
  onkeydown={(e) => {
    if (e.key === 'Enter' || e.key === ' ') {
      e.preventDefault();
      onSelectTicket();
    }
  }}
  class="bg-white dark:bg-slate-950 border border-slate-200 dark:border-slate-800 rounded-xl p-4 shadow-sm active:scale-[0.99] transition-all cursor-pointer space-y-3.5 relative overflow-hidden text-left font-sans"
>
  <div class="absolute left-0 top-0 bottom-0 w-1 {ticket.status === 'cancelled' ? 'bg-rose-500' : 'bg-blue-600'}"></div>
  
  <div class="flex justify-between items-start pl-2">
    <div>
      <h4 class="font-bold text-slate-900 dark:text-white flex flex-wrap items-center gap-1.5 leading-tight">
        {ticket.brand}
        {#if ticket.is_warranty}
          <span class="inline-flex items-center gap-0.5 px-1.5 py-0.5 rounded text-[9px] font-bold bg-purple-50 text-purple-700 border border-purple-100 dark:bg-purple-950/20 dark:text-purple-400 dark:border-purple-900/30">
            <RefreshCw size={8} />
            Klaim
          </span>
        {/if}
      </h4>
      <div class="text-xs text-slate-500 dark:text-slate-400 mt-0.5 flex flex-wrap items-center gap-1.5">
        <span>{ticket.model}</span>
        {#if !ticket.is_warranty && ticket.status === 'picked_up'}
          {@const warranty = checkWarrantyExpiry(ticket.exit_date, ticket.warranty_days)}
          {#if warranty && warranty.isValid}
            <span class="inline-flex items-center gap-0.5 px-1.5 py-0.5 rounded text-[9px] font-bold bg-indigo-50 text-indigo-700 border border-indigo-100 dark:bg-indigo-950/20 dark:text-indigo-400 dark:border-indigo-900/30">
              <ShieldCheck size={8} />
              {warranty.remainingDays} H
            </span>
          {/if}
        {/if}
      </div>
    </div>
    
    <span class="inline-flex items-center px-2 py-0.5 rounded-full text-[10px] font-bold border {getStatusBadgeClass(ticket.status, Boolean(ticket.is_warranty))}">
      {getStatusLabel(ticket.status, Boolean(ticket.is_warranty))}
    </span>
  </div>

  <div class="pl-2 grid grid-cols-2 gap-2 text-xs">
    <div>
      <span class="text-slate-400 block text-[10px] uppercase font-bold tracking-wider mb-0.5">Customer</span>
      <span class="text-slate-900 dark:text-slate-200 font-semibold">{ticket.customer_name}</span>
      <span class="text-slate-500 block text-[10px] mt-0.5">{ticket.customer_gender}</span>
    </div>
    <div>
      <span class="text-slate-400 block text-[10px] uppercase font-bold tracking-wider mb-0.5">Kerusakan</span>
      <span class="text-slate-800 dark:text-slate-200 line-clamp-2">{ticket.issue}</span>
    </div>
  </div>

  <div class="pl-2 pt-3 border-t border-slate-100 dark:border-slate-800/50 flex justify-between items-center">
    <div>
      <span class="text-slate-400 block text-[10px] uppercase font-bold tracking-wider mb-0.5">Price</span>
      <span class="text-sm font-extrabold text-slate-950 dark:text-white">{formatCurrency(Number(ticket.price))}</span>
    </div>
    
    {#if ticket.status !== "picked_up"}
      <button
        onclick={(e) => {
          e.stopPropagation();
          onQuickAction();
        }}
        disabled={isActionLoading}
        class="px-3 py-1.5 bg-slate-900 hover:bg-slate-800 dark:bg-slate-800 dark:hover:bg-slate-700 disabled:bg-slate-100 disabled:text-slate-400 text-white font-bold text-[10px] uppercase tracking-wider rounded-lg transition-all inline-flex items-center gap-1 shadow-sm"
      >
        {#if isActionLoading}
          <LoaderCircle size={10} class="animate-spin" />
        {/if}
        {#if ticket.status === "service_in"}
          Mulai
        {:else if ticket.status === "on_process"}
          Selesai
        {:else if ticket.status === "waiting_confirmation"}
          Hubungi
        {:else if ticket.status === "fixed"}
          Ambil & Bayar
        {/if}
      </button>
    {:else}
      <span class="text-xs font-semibold text-slate-400">Selesai</span>
    {/if}
  </div>
</div>
