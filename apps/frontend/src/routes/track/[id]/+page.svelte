<script lang="ts">
  import { page } from '$app/stores';
  import { onMount } from 'svelte';
  import { Check, Smartphone, Calendar, ShieldCheck, DollarSign, LoaderCircle, ArrowLeft } from 'lucide-svelte';
  import type { Ticket } from '$lib/types/ticket';
  import { formatCurrency, formatDate } from '$lib/utils/format';
  import { getWarrantyExpiryDate } from '$lib/utils/warranty';

  let ticket = $state<Ticket | null>(null);
  let isLoading = $state(true);
  let errorMsg = $state('');

  const steps = [
      { status: 'service_in', label: 'Diterima', desc: 'HP masuk ke antrean servis' },
      { status: 'on_process', label: 'Diperbaiki', desc: 'Teknisi sedang melakukan reparasi' },
      { status: 'waiting_confirmation', label: 'Menunggu', desc: 'Awaiting customer confirmation' },
      { status: 'fixed', label: 'Siap Ambil', desc: 'Perbaikan selesai dilakukan' },
      { status: 'picked_up', label: 'Sudah Diambil', desc: 'Perangkat diambil & garansi aktif' }
  ];

  let id = $derived($page.params.id || '');

  $effect(() => {
      if (id) {
          fetchTicket(id);
      }
  });

  async function fetchTicket(ticketId: string) {
      isLoading = true;
      const uuidRegex = /^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$/;
      
      if (!uuidRegex.test(ticketId)) {
          errorMsg = 'Tiket tidak ditemukan atau format ID salah.';
          isLoading = false;
          return;
      }

      try {
          const res = await fetch(`/api/v1/public/tickets/${ticketId}`);
          const payload = await res.json();
          if (!res.ok) {
              errorMsg = payload.error || 'Gagal memuat detail perbaikan';
          } else {
              ticket = payload.data;
              errorMsg = '';
          }
      } catch (err) {
          errorMsg = 'Koneksi gagal. Silakan coba lagi.';
      } finally {
          isLoading = false;
      }
  }

  function getStepIndex(currentStatus: string): number {
      if (currentStatus === 'cancelled') return -1;
      return steps.findIndex(s => s.status === currentStatus);
  }
</script>

<div class="min-h-screen bg-gradient-to-br from-slate-900 via-slate-950 to-blue-950 flex justify-center items-center p-4">
    {#if isLoading}
        <div class="text-center text-slate-400 flex flex-col items-center gap-2">
            <LoaderCircle size={32} class="animate-spin text-blue-500" />
            <span>Memuat status perbaikan...</span>
        </div>
    {:else if errorMsg}
        <div class="w-full max-w-md bg-slate-900/80 border border-rose-500/30 rounded-2xl p-8 text-center space-y-4">
            <h2 class="text-xl font-bold text-rose-400">Terjadi Kesalahan</h2>
            <p class="text-slate-400 text-sm">{errorMsg}</p>
            <a href="/track" class="inline-flex items-center gap-1.5 text-xs text-blue-400 font-bold uppercase tracking-wider hover:underline">
                <ArrowLeft size={14} /> Kembali ke halaman pencarian
            </a>
        </div>
    {:else if ticket}
        <div class="w-full max-w-2xl bg-slate-900/80 backdrop-blur-md border border-slate-700/50 rounded-2xl p-6 md:p-8 shadow-2xl space-y-8">
            <!-- Header -->
            <div class="flex justify-between items-start">
                <div>
                    <a href="/track" class="inline-flex items-center gap-1 text-xs text-slate-400 hover:text-white transition-colors mb-2">
                        <ArrowLeft size={12} /> Kembali
                    </a>
                    <h2 class="text-xl md:text-2xl font-black text-white">{ticket.brand} {ticket.model}</h2>
                    <p class="text-xs text-slate-400 font-mono mt-0.5">Kode Tiket: {ticket.id.slice(0,8).toUpperCase()}</p>
                </div>
                <span class="inline-flex items-center px-3 py-1 rounded-full text-xs font-bold uppercase tracking-wider
                    {ticket.status === 'cancelled' ? 'bg-rose-500/20 text-rose-400 border border-rose-500/30' : 'bg-emerald-500/20 text-emerald-400 border border-emerald-500/30'}">
                    {ticket.status === 'service_in' ? 'Masuk' :
                     ticket.status === 'on_process' ? 'Diproses' :
                     ticket.status === 'waiting_confirmation' ? 'Menunggu Konfirmasi' :
                     ticket.status === 'fixed' ? 'Siap Diambil' :
                     ticket.status === 'picked_up' ? 'Sudah Diambil' : ticket.status}
                </span>
            </div>

            <!-- Status Stepper Timeline -->
            {#if ticket.status !== 'cancelled'}
                {@const activeIndex = getStepIndex(ticket.status)}
                <div class="overflow-x-auto pb-4 -mx-6 px-6 scrollbar-none">
                    <div class="flex flex-row justify-between items-start gap-4 min-w-[550px] relative py-2">
                        <!-- Connecting Line Behind -->
                        <div class="absolute top-[18px] left-[10%] right-[10%] h-[2px] bg-slate-800 -z-0"></div>
                        <div class="absolute top-[18px] left-[10%] h-[2px] bg-blue-500 transition-all duration-500 -z-0" style="width: {activeIndex * 20}%"></div>

                        {#each steps as step, idx}
                            <div class="flex flex-col items-center gap-2 flex-1 text-center relative z-10">
                                <!-- Step Circle -->
                                <div class="w-8 h-8 rounded-full border-2 flex items-center justify-center font-bold text-xs shrink-0 transition-all duration-300
                                    {idx <= activeIndex ? 'bg-blue-600 border-blue-500 text-white shadow-lg shadow-blue-500/20' : 'bg-slate-800 border-slate-700 text-slate-500'}">
                                    {#if idx < activeIndex}
                                        <Check size={14} />
                                    {:else}
                                        {idx + 1}
                                    {/if}
                                </div>
                                <!-- Step Label & Description -->
                                <div>
                                    <h4 class="text-[11px] font-bold tracking-wide {idx <= activeIndex ? 'text-white' : 'text-slate-500'}">{step.label}</h4>
                                    <p class="text-[9px] text-slate-500 mt-0.5 px-1 leading-tight">{step.desc}</p>
                                </div>
                            </div>
                        {/each}
                    </div>
                </div>
            {:else}
                <div class="bg-rose-950/20 border border-rose-500/20 rounded-xl p-4 text-center">
                    <p class="text-sm font-bold text-rose-400">Status Perbaikan Dibatalkan</p>
                    <p class="text-xs text-rose-300/80 mt-1">Perangkat tidak dilanjutkan untuk diperbaiki. Silakan hubungi toko untuk informasi pengembalian.</p>
                </div>
            {/if}

            <!-- Details Cards Grid -->
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                <!-- Device Spec -->
                <div class="bg-slate-800/30 border border-slate-700/50 rounded-xl p-4 space-y-3">
                    <h3 class="text-xs font-black text-slate-400 uppercase tracking-widest flex items-center gap-1.5"><Smartphone size={14} /> Detail Servis</h3>
                    <div class="space-y-1.5 text-xs">
                        <p class="text-slate-400">Keluhan: <span class="text-white font-semibold">{ticket.issue}</span></p>
                        <p class="text-slate-400">Deskripsi: <span class="text-white font-semibold">{ticket.additional_description || '-'}</span></p>
                        <p class="text-slate-400">Aksesoris: <span class="text-white font-semibold">{ticket.accessories || '-'}</span></p>
                        <p class="text-slate-400">Tanggal Masuk: <span class="text-white font-semibold">{formatDate(ticket.entry_date)}</span></p>
                        {#if ticket.exit_date}
                            <p class="text-slate-400">Tanggal Selesai: <span class="text-white font-semibold">{formatDate(ticket.exit_date)}</span></p>
                        {/if}
                    </div>
                </div>

                <!-- Price & Warranty -->
                <div class="bg-slate-800/30 border border-slate-700/50 rounded-xl p-4 space-y-3">
                    <h3 class="text-xs font-black text-slate-400 uppercase tracking-widest flex items-center gap-1.5"><span class="text-[10px] font-black bg-slate-800 px-1.5 py-0.5 rounded text-slate-400 mr-1 select-none">RP</span> Estimasi Biaya & Garansi</h3>
                    <div class="space-y-1.5 text-xs">
                        <p class="text-slate-400">Biaya Perbaikan: <span class="text-white font-bold text-sm">{formatCurrency(Number(ticket.price))}</span></p>
                        <p class="text-slate-400">Status Pembayaran: 
                            <span class="font-bold uppercase px-1.5 py-0.5 rounded text-[10px] ml-1 {ticket.payment_status === 'paid' ? 'bg-emerald-500/10 text-emerald-400' : 'bg-amber-500/10 text-amber-400'}">
                                {ticket.payment_status}
                            </span>
                        </p>
                        <p class="text-slate-400">Durasi Garansi: <span class="text-white font-semibold">{ticket.warranty_days} Hari</span></p>
                        {#if ticket.status === 'picked_up'}
                            <p class="text-slate-400 flex items-center gap-1">Status Garansi: 
                                <span class="text-emerald-400 font-bold flex items-center gap-0.5"><ShieldCheck size={12} /> Aktif</span>
                            </p>
                            <p class="text-slate-400">Tanggal Berakhir: <span class="text-white font-bold">{getWarrantyExpiryDate(ticket.exit_date, ticket.warranty_days)}</span></p>
                        {:else}
                            <p class="text-slate-400 flex items-center gap-1">Status Garansi: 
                                <span class="text-slate-300 font-semibold">Aktif setelah HP diambil</span>
                            </p>
                        {/if}
                    </div>
                </div>
            </div>
        </div>
    {/if}
</div>

<style>
  .scrollbar-none::-webkit-scrollbar {
      display: none;
  }
  .scrollbar-none {
      -ms-overflow-style: none;
      scrollbar-width: none;
  }
</style>
