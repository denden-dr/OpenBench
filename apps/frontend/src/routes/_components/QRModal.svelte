<script lang="ts">
  import QRCode from 'qrcode';

  let {
    isOpen,
    qrUrl,
    onClose,
    onCopyQRUrl,
  }: {
    isOpen: boolean;
    qrUrl: string;
    onClose: () => void;
    onCopyQRUrl: () => void;
  } = $props();

  let canvasElement = $state<HTMLCanvasElement | null>(null);

  $effect(() => {
    if (isOpen && qrUrl && canvasElement) {
      QRCode.toCanvas(canvasElement, qrUrl, { width: 200, margin: 1 }, (err) => {
        if (err) {
          console.error('Failed to generate QR code:', err);
        }
      });
    }
  });
</script>

<svelte:window onkeydown={(e) => { if (e.key === 'Escape' && isOpen) onClose(); }} />

{#if isOpen}
  <div class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-900/60 backdrop-blur-sm animate-fade-in font-sans">
    <div class="bg-white dark:bg-slate-950 w-full max-w-sm rounded-2xl border border-slate-200 dark:border-slate-800 shadow-2xl p-6 text-center space-y-6 animate-scale-in">
      <h3 class="font-bold text-lg text-slate-900 dark:text-white">QR Code Pelacakan</h3>
      <div class="bg-slate-50 dark:bg-slate-900 border border-slate-100 dark:border-slate-850 p-4 rounded-xl inline-block mx-auto">
        <canvas bind:this={canvasElement} class="mx-auto w-[200px] h-[200px]"></canvas>
      </div>

      <div class="space-y-2 text-left">
        <label for="qr-url-input" class="text-[10px] font-bold tracking-wider text-slate-400 dark:text-slate-500 uppercase">URL Pelacakan</label>
        <div class="flex gap-2">
          <input
            id="qr-url-input"
            type="text"
            bind:value={qrUrl}
            class="flex-1 px-3 py-2 text-xs border border-slate-200 dark:border-slate-800 rounded-xl bg-slate-50 dark:bg-slate-900 text-slate-850 dark:text-slate-150 focus:outline-none focus:ring-1 focus:ring-slate-450 dark:focus:ring-slate-700 transition-all font-mono"
            placeholder="http://localhost:5173/track/..."
          />
          <button
            type="button"
            onclick={onCopyQRUrl}
            class="px-3.5 py-2 bg-slate-100 hover:bg-slate-200 dark:bg-slate-800 dark:hover:bg-slate-700 text-slate-700 dark:text-slate-300 text-xs font-semibold rounded-xl cursor-pointer transition-all active:scale-95"
          >
            Salin
          </button>
        </div>
        <p class="text-[10px] text-slate-400 dark:text-slate-500 mt-1 leading-normal">
          Ubah URL jika Anda mengakses menggunakan port forwarding (misal. VS Code preview URL) agar HP pelanggan dapat memindai & mengakses halaman tersebut.
        </p>
      </div>

      <p class="text-xs text-slate-500 dark:text-slate-400">Scan QR Code ini dengan HP untuk langsung membuka halaman pelacakan status.</p>
      <button
        type="button"
        onclick={onClose}
        class="w-full py-2.5 bg-slate-900 hover:bg-slate-800 dark:bg-slate-800 dark:hover:bg-slate-700 text-white text-xs font-bold rounded-xl cursor-pointer transition-all"
      >
        Tutup
      </button>
    </div>
  </div>
{/if}
