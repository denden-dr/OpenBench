<script lang="ts">
  import { Button } from '$lib';
  import { FileText, Printer, X } from 'lucide-svelte';
  import type { Sale } from '$lib/services/sales';

  let {
    selectedSale,
    onClose
  }: {
    selectedSale: Sale;
    onClose: () => void;
  } = $props();

  function printReceipt() {
    window.print();
  }
</script>

<div class="fixed inset-0 bg-black/60 z-50 flex items-center justify-center p-4">
  <div class="max-w-md w-full bg-white border-4 border-neubrutalism-charcoal shadow-neubrutalism-lg flex flex-col overflow-hidden max-h-[90vh]">
    
    <!-- Header -->
    <div class="p-4 bg-neubrutalism-green border-b-4 border-neubrutalism-charcoal flex justify-between items-center text-neubrutalism-charcoal font-display font-extrabold text-sm uppercase">
      <div class="flex items-center gap-2">
        <FileText class="w-5 h-5" />
        <span>Transaction Receipt</span>
      </div>
      <button onclick={onClose} class="p-1 hover:bg-white/20">
        <X class="w-5 h-5" />
      </button>
    </div>

    <!-- Receipt Body -->
    <div class="p-6 overflow-y-auto font-mono text-xs flex flex-col gap-4 text-neubrutalism-charcoal" id="print-area">
      <div class="text-center flex flex-col gap-1 border-b-2 border-dashed border-zinc-300 pb-3">
        <span class="font-display font-extrabold text-base tracking-widest uppercase">OPENBENCH REPAIRS</span>
        <span class="text-[10px] text-zinc-500 uppercase">Jalan Teknisi Raya No. 101, Bandung</span>
        <span class="text-[10px] text-zinc-500 uppercase">Telp: 0812-3456-7890</span>
      </div>

      <div class="flex flex-col gap-1 border-b border-dashed border-zinc-200 pb-2">
        <p>Invoice: {selectedSale.invoice_number}</p>
        <p>Date: {new Date(selectedSale.created_at).toLocaleString('id-ID')}</p>
        <p>Cashier: admin@openbench.dev</p>
        <p>Payment: {selectedSale.payment_method.toUpperCase()}</p>
      </div>

      <!-- Items list -->
      <div class="flex flex-col gap-2.5 border-b-2 border-dashed border-zinc-300 pb-3">
        {#each selectedSale.items as item}
          <div class="flex flex-col gap-0.5">
            <span class="font-sans font-semibold text-zinc-800">{item.name}</span>
            <div class="flex justify-between text-[10px] text-zinc-600">
              <span>{item.qty} x IDR {item.price.toLocaleString('id-ID')}</span>
              <span>IDR {(item.price * item.qty).toLocaleString('id-ID')}</span>
            </div>
          </div>
        {/each}
      </div>

      <!-- Totals -->
      <div class="flex flex-col gap-1 text-right font-bold pr-1 border-b border-dashed border-zinc-200 pb-2">
        <div class="flex justify-between">
          <span>SUBTOTAL:</span>
          <span>IDR {selectedSale.subtotal.toLocaleString('id-ID')}</span>
        </div>
        {#if selectedSale.discount > 0}
          <div class="flex justify-between text-rose-600">
            <span>DISCOUNT:</span>
            <span>- IDR {selectedSale.discount.toLocaleString('id-ID')}</span>
          </div>
        {/if}
        <div class="flex justify-between text-sm font-black text-neubrutalism-charcoal mt-1 pt-1 border-t border-zinc-200">
          <span>TOTAL PAID:</span>
          <span>IDR {selectedSale.total.toLocaleString('id-ID')}</span>
        </div>
      </div>

      <!-- Footer Msg -->
      <div class="text-center font-sans text-[9px] text-zinc-500 mt-2 flex flex-col gap-0.5">
        <span>THANK YOU FOR SHOPPING!</span>
        <span>Keep this receipt as valid proof of purchase.</span>
      </div>
    </div>

    <!-- Actions -->
    <div class="p-4 border-t-4 border-neubrutalism-charcoal bg-zinc-50 flex gap-3">
      <Button 
        bgColor="bg-white" 
        onclick={printReceipt}
        class="flex-1 py-2 font-bold shadow-neubrutalism-sm border-2 border-neubrutalism-charcoal flex items-center justify-center gap-1.5"
      >
        <Printer class="w-4 h-4" />
        <span>PRINT RECEIPT</span>
      </Button>

      <Button 
        bgColor="bg-neubrutalism-yellow" 
        onclick={onClose}
        class="flex-1 py-2 font-bold shadow-neubrutalism-sm border-2 border-neubrutalism-charcoal flex items-center justify-center"
      >
        <span>CLOSE</span>
      </Button>
    </div>

  </div>
</div>
