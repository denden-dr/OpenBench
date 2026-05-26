<script lang="ts">
  import { Search, LoaderCircle, ArrowRight, Smartphone, ShieldAlert } from 'lucide-svelte';

  let ticketCode = $state('');
  let phone = $state('');
  let isLoading = $state(false);
  let errorMsg = $state('');
  let showPhoneVerification = $state(false);

  function isFullUUID(value: string) {
      return /^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$/.test(value.trim());
  }

  function isShortID(value: string) {
      return /^[0-9a-fA-F]{8}$/.test(value.trim());
  }

  async function handleSearch(event: SubmitEvent) {
      event.preventDefault();
      isLoading = true;
      errorMsg = '';
      try {
          const code = ticketCode.trim();
          if (isFullUUID(code)) {
              window.location.href = `/track/${code}`;
              return;
          }
          if (!isShortID(code)) {
              errorMsg = 'Masukkan 8 karakter kode tiket atau UUID lengkap.';
              return;
          }
          showPhoneVerification = true;
      } catch (err) {
          errorMsg = 'Koneksi gagal. Silakan coba lagi.';
      } finally {
          isLoading = false;
      }
  }

  async function handleVerify(event: SubmitEvent) {
      event.preventDefault();
      isLoading = true;
      errorMsg = '';
      try {
          const res = await fetch('/api/v1/public/track', {
              method: 'POST',
              headers: { 'Content-Type': 'application/json' },
              body: JSON.stringify({ short_id: ticketCode.trim(), phone })
          });
          const payload = await res.json();
          if (!res.ok) {
              errorMsg = payload.error || 'Verifikasi gagal. Periksa nomor HP Anda.';
          } else {
              window.location.href = `/track/${payload.ticket_id}`;
          }
      } catch (err) {
          errorMsg = 'Koneksi gagal. Silakan coba lagi.';
      } finally {
          isLoading = false;
      }
  }
</script>

<div class="min-h-screen bg-gradient-to-br from-slate-900 via-slate-950 to-blue-950 flex flex-col justify-center items-center p-4">
    <div class="w-full max-w-md bg-slate-900/80 backdrop-blur-md border border-slate-700/50 rounded-2xl p-8 shadow-2xl space-y-6">
        <div class="text-center">
            <h1 class="text-2xl font-black text-white tracking-wide">OPENBENCH</h1>
            <p class="text-xs text-slate-400 mt-1 uppercase tracking-widest">Repair Tracking Portal</p>
        </div>

        {#if errorMsg}
            <div class="bg-rose-950/30 border border-rose-500/30 rounded-xl p-4 flex gap-2 text-rose-300 text-xs">
                <ShieldAlert size={16} class="shrink-0" />
                <span>{errorMsg}</span>
            </div>
        {/if}

        {#if !showPhoneVerification}
            <!-- Search ticket code -->
            <form onsubmit={handleSearch} class="space-y-4">
                <div class="space-y-1.5">
                    <label for="track-ticket-code" class="text-xs font-bold text-slate-400 block uppercase tracking-wider">Masukkan Kode Tiket</label>
                    <input
                        id="track-ticket-code"
                        type="text"
                        bind:value={ticketCode}
                        required
                        placeholder="8 karakter atau UUID lengkap"
                        class="w-full bg-slate-800/50 border border-slate-700 text-white rounded-lg px-4 py-3 focus:outline-none focus:ring-2 focus:ring-blue-500/50 text-sm font-semibold transition-all"
                    />
                </div>
                <button
                    type="submit"
                    disabled={isLoading}
                    class="w-full py-3 bg-blue-600 hover:bg-blue-700 text-white font-bold text-xs uppercase tracking-widest rounded-xl transition-all shadow-lg active:scale-95 inline-flex items-center justify-center gap-2"
                >
                    {#if isLoading}
                        <LoaderCircle size={14} class="animate-spin" />
                        Mencari...
                    {:else}
                        <Search size={14} />
                        Cari Status
                    {/if}
                </button>
            </form>
        {:else}
            <!-- Verify Phone Number -->
            <form onsubmit={handleVerify} class="space-y-4 animate-fade-in">
                <div class="space-y-1.5">
                    <label for="track-phone" class="text-xs font-bold text-slate-400 block uppercase tracking-wider">Verifikasi Nomor HP</label>
                    <input
                        id="track-phone"
                        type="text"
                        bind:value={phone}
                        required
                        placeholder="Masukkan No HP terdaftar (misal: 0812...)"
                        class="w-full bg-slate-800/50 border border-slate-700 text-white rounded-lg px-4 py-3 focus:outline-none focus:ring-2 focus:ring-blue-500/50 text-sm font-semibold transition-all"
                    />
                </div>
                <div class="flex gap-2">
                    <button
                        type="button"
                        onclick={() => showPhoneVerification = false}
                        class="w-1/3 py-3 border border-slate-700 text-slate-300 font-bold text-xs uppercase tracking-widest rounded-xl transition-all hover:bg-slate-800/50"
                    >
                        Kembali
                    </button>
                    <button
                        type="submit"
                        disabled={isLoading}
                        class="w-2/3 py-3 bg-blue-600 hover:bg-blue-700 text-white font-bold text-xs uppercase tracking-widest rounded-xl transition-all shadow-lg active:scale-95 inline-flex items-center justify-center gap-2"
                    >
                        {#if isLoading}
                            <LoaderCircle size={14} class="animate-spin" />
                        {:else}
                            <ArrowRight size={14} />
                        {/if}
                        Verifikasi
                    </button>
                </div>
            </form>
        {/if}
    </div>
</div>
