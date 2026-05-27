<script lang="ts">
  import { LoaderCircle, AlertTriangle } from "lucide-svelte";
  import type { Ticket } from "$lib/types/ticket";

  let {
    isOpen,
    ticket,
    isDeleting,
    onConfirm,
    onClose,
  }: {
    isOpen: boolean;
    ticket: Ticket | null;
    isDeleting: boolean;
    onConfirm: () => void;
    onClose: () => void;
  } = $props();
</script>

<svelte:window onkeydown={(e) => { if (e.key === 'Escape' && isOpen) onClose(); }} />

{#if isOpen}
  <div class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-900/60 backdrop-blur-sm animate-fade-in font-sans">
    <div class="bg-white dark:bg-slate-950 w-full max-w-md rounded-2xl border border-slate-200 dark:border-slate-800 shadow-2xl p-6 text-center space-y-6 animate-scale-in">
      
      <!-- Warning Icon Header -->
      <div class="mx-auto w-12 h-12 bg-red-50 dark:bg-red-950/30 rounded-full flex items-center justify-center border border-red-100 dark:border-red-900 text-red-600 dark:text-red-400">
        <AlertTriangle size={24} />
      </div>

      <div class="space-y-2">
        <h3 class="font-bold text-lg text-slate-900 dark:text-white">Hapus Tiket Permanen</h3>
        <p class="text-xs text-slate-500 dark:text-slate-400 leading-relaxed">
          Apakah Anda yakin ingin menghapus tiket ini secara permanen?
        </p>
      </div>

      <!-- Ticket Summary Card -->
      {#if ticket}
        <div class="bg-slate-50 dark:bg-slate-900/50 border border-slate-100 dark:border-slate-850 rounded-xl p-4 text-left space-y-2">
          <div class="flex justify-between items-center text-xs">
            <span class="font-mono font-bold text-slate-700 dark:text-slate-300">
              {ticket.id.substring(0, 8)}
            </span>
            <span class="text-slate-400 dark:text-slate-500">
              {ticket.brand} - {ticket.model}
            </span>
          </div>
          <div class="text-xs text-slate-600 dark:text-slate-400">
            <span class="font-semibold text-slate-800 dark:text-slate-200">Pelanggan:</span> {ticket.customer_name}
          </div>
          <div class="text-xs text-slate-600 dark:text-slate-400">
            <span class="font-semibold text-slate-800 dark:text-slate-200">Kerusakan:</span> {ticket.issue}
          </div>
        </div>
      {/if}

      <div class="text-xs text-rose-600 dark:text-rose-400 font-semibold bg-rose-50 dark:bg-rose-950/20 py-2.5 px-3 rounded-lg border border-rose-100 dark:border-rose-900/30">
        Peringatan: Tindakan ini tidak dapat dibatalkan!
      </div>

      <!-- Action Buttons -->
      <div class="flex gap-3">
        <button
          type="button"
          disabled={isDeleting}
          onclick={onClose}
          class="flex-1 py-2.5 bg-slate-100 hover:bg-slate-200 dark:bg-slate-900 dark:hover:bg-slate-800 text-slate-700 dark:text-slate-300 text-xs font-bold rounded-xl cursor-pointer transition-all active:scale-95 disabled:opacity-50 disabled:cursor-not-allowed"
        >
          Batal
        </button>
        <button
          type="button"
          disabled={isDeleting}
          onclick={onConfirm}
          class="flex-1 py-2.5 bg-red-600 hover:bg-red-700 text-white text-xs font-bold rounded-xl cursor-pointer transition-all active:scale-95 flex items-center justify-center gap-2 disabled:bg-red-500 disabled:opacity-75 disabled:cursor-not-allowed"
        >
          {#if isDeleting}
            <LoaderCircle class="animate-spin" size={14} />
            <span>Menghapus...</span>
          {:else}
            <span>Ya, Hapus</span>
          {/if}
        </button>
      </div>
    </div>
  </div>
{/if}

<style>
  @keyframes fade-in {
    from { opacity: 0; }
    to { opacity: 1; }
  }
  @keyframes scale-in {
    from { transform: scale(0.95); opacity: 0; }
    to { transform: scale(1); opacity: 1; }
  }
  .animate-fade-in {
    animation: fade-in 0.2s ease-out forwards;
  }
  .animate-scale-in {
    animation: scale-in 0.2s cubic-bezier(0.16, 1, 0.3, 1) forwards;
  }
</style>
