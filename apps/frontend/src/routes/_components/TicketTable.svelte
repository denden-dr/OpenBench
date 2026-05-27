<script lang="ts">
  import { RefreshCw, ShieldCheck, LoaderCircle } from "lucide-svelte";
  import type { Ticket } from "$lib/types/ticket";
  import { formatCurrency } from "$lib/utils/format";
  import { checkWarrantyExpiry } from "$lib/utils/warranty";

  let {
    tickets,
    isActionLoading,
    onSelectTicket,
    onQuickAction,
    getStatusBadgeClass,
    getStatusLabel,
  }: {
    tickets: Ticket[];
    isActionLoading: Record<string, boolean>;
    onSelectTicket: (ticket: Ticket) => void;
    onQuickAction: (ticketId: string, currentStatus: string) => void;
    getStatusBadgeClass: (status: string, isWarranty: boolean) => string;
    getStatusLabel: (status: string, isWarranty: boolean) => string;
  } = $props();
</script>

<div class="hidden md:block bg-white dark:bg-slate-950 border border-slate-200 dark:border-slate-800 rounded-2xl overflow-hidden shadow-sm">
  <div class="overflow-x-auto">
    <table class="w-full text-left border-collapse">
      <thead>
        <tr class="bg-slate-50 dark:bg-slate-900/50 border-b border-slate-200 dark:border-slate-800 text-xs font-bold text-slate-500 dark:text-slate-400 uppercase tracking-wider">
          <th class="py-4 px-6">Device</th>
          <th class="py-4 px-6">Customer</th>
          <th class="py-4 px-6">Kerusakan</th>
          <th class="py-4 px-6">Price</th>
          <th class="py-4 px-6">Status</th>
          <th class="py-4 px-6 text-right">Quick Action</th>
        </tr>
      </thead>
      <tbody class="divide-y divide-slate-200 dark:divide-slate-800 text-sm text-slate-900 dark:text-slate-200">
        {#each tickets as ticket (ticket.id)}
          <tr
            onclick={() => onSelectTicket(ticket)}
            class="hover:bg-slate-50/70 dark:hover:bg-slate-900/30 transition-colors cursor-pointer group"
          >
            <td class="py-4 px-6">
              <div class="font-semibold text-slate-900 dark:text-white group-hover:text-blue-600 transition-colors flex items-center gap-1.5">
                {ticket.brand}
                {#if ticket.is_warranty}
                  <span class="inline-flex items-center gap-0.5 px-1.5 py-0.5 rounded text-[10px] font-bold bg-purple-50 text-purple-700 border border-purple-100 dark:bg-purple-950/20 dark:text-purple-400 dark:border-purple-900/30">
                    <RefreshCw size={10} />
                    Klaim Garansi
                  </span>
                {/if}
              </div>
              <div class="text-xs text-slate-500 dark:text-slate-400 mt-0.5 flex flex-wrap items-center gap-1.5">
                <span>{ticket.model}</span>
                {#if !ticket.is_warranty && ticket.status === 'picked_up'}
                  {@const warranty = checkWarrantyExpiry(ticket.exit_date, ticket.warranty_days)}
                  {#if warranty && warranty.isValid}
                    <span class="inline-flex items-center gap-0.5 px-1.5 py-0.5 rounded text-[10px] font-bold bg-indigo-50 text-indigo-700 border border-indigo-100 dark:bg-indigo-950/20 dark:text-indigo-400 dark:border-indigo-900/30">
                      <ShieldCheck size={10} />
                      Garansi Aktif ({warranty.remainingDays} H)
                    </span>
                  {/if}
                {/if}
              </div>
            </td>
            <td class="py-4 px-6">
              <div class="font-medium text-slate-900 dark:text-white">
                {ticket.customer_name}
              </div>
              <div class="text-xs text-slate-500 dark:text-slate-400 mt-0.5">
                {ticket.customer_gender}
              </div>
            </td>
            <td class="py-4 px-6">
              <div class="line-clamp-1 max-w-xs">{ticket.issue}</div>
            </td>
            <td class="py-4 px-6 font-semibold">
              {formatCurrency(Number(ticket.price))}
            </td>
            <td class="py-4 px-6">
              <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-semibold border {getStatusBadgeClass(ticket.status, Boolean(ticket.is_warranty))}">
                {getStatusLabel(ticket.status, Boolean(ticket.is_warranty))}
              </span>
            </td>
            <td class="py-4 px-6 text-right" onclick={(e) => e.stopPropagation()}>
              {#if ticket.status !== "picked_up"}
                <button
                  onclick={() => {
                    onQuickAction(ticket.id, ticket.status);
                  }}
                  disabled={isActionLoading[ticket.id]}
                  class="px-4 py-1.5 bg-slate-900 hover:bg-slate-800 dark:bg-slate-800 dark:hover:bg-slate-700 disabled:bg-slate-100 disabled:text-slate-400 text-white font-bold text-xs uppercase tracking-wider rounded-lg transition-all active:scale-95 inline-flex items-center gap-1.5 shadow-sm"
                >
                  {#if isActionLoading[ticket.id]}
                    <LoaderCircle size={12} class="animate-spin" />
                  {/if}
                  {#if ticket.status === "service_in"}
                    Mulai Proses
                  {:else if ticket.status === "on_process"}
                    Tandai Selesai
                  {:else if ticket.status === "waiting_confirmation"}
                    Hubungi Pelanggan
                  {:else if ticket.status === "fixed"}
                    Ambil & Bayar
                  {/if}
                </button>
              {:else}
                <span class="text-xs font-semibold text-slate-400">Selesai</span>
              {/if}
            </td>
          </tr>
        {/each}
      </tbody>
    </table>
  </div>
</div>
