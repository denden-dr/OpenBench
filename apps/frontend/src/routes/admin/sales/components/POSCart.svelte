<script lang="ts">
  import { Card, Button } from '$lib';
  import { Plus, Minus, ArrowRight } from 'lucide-svelte';
  import type { Product } from '$lib/services/inventory';
  import { formatCurrencyInput, parseCurrencyInput } from '$lib/utils/format';

  interface CartItem {
    product: Product;
    qty: number;
  }

  let {
    cart,
    discount = $bindable(0),
    paymentMethod = $bindable('cash'),
    cashPaid = $bindable(0),
    subtotal,
    finalTotal,
    changeAmount,
    onAdjustQty,
    onCheckout
  }: {
    cart: CartItem[];
    discount: number;
    paymentMethod: 'cash' | 'qris';
    cashPaid: number;
    subtotal: number;
    finalTotal: number;
    changeAmount: number;
    onAdjustQty: (id: string, amt: number) => void;
    onCheckout: () => void;
  } = $props();

  let totalItems = $derived(cart.reduce((sum, item) => sum + item.qty, 0));

  let displayDiscount = $state(formatCurrencyInput(discount));
  let displayCashPaid = $state(formatCurrencyInput(cashPaid));

  $effect(() => {
    displayDiscount = formatCurrencyInput(discount);
  });

  $effect(() => {
    displayCashPaid = formatCurrencyInput(cashPaid);
  });

  function handleDiscountInput(e: Event) {
    const target = e.target as HTMLInputElement;
    const numeric = parseCurrencyInput(target.value);
    discount = numeric;
    displayDiscount = formatCurrencyInput(numeric);
  }

  function handleCashPaidInput(e: Event) {
    const target = e.target as HTMLInputElement;
    const numeric = parseCurrencyInput(target.value);
    cashPaid = numeric;
    displayCashPaid = formatCurrencyInput(numeric);
  }
</script>

<Card bgColor="bg-white" class="border-4 border-neubrutalism-charcoal p-6 shadow-neubrutalism-md sticky top-24 flex flex-col gap-4">
  <h3 class="font-display font-extrabold text-lg uppercase tracking-tight flex items-center gap-1.5 border-b-2 border-dashed border-zinc-300 pb-2">
    <span>Shopping Cart</span>
    <span class="text-xs font-mono font-bold bg-neubrutalism-charcoal text-white rounded px-2 py-0.5 ml-auto">
      {totalItems} Items
    </span>
  </h3>

  <!-- Cart Items list -->
  <div class="flex-grow overflow-y-auto max-h-72 flex flex-col gap-3 pr-1">
    {#if cart.length === 0}
      <div class="p-8 text-center text-zinc-400 font-mono text-xs flex flex-col items-center justify-center gap-2 h-48 border-2 border-dashed border-zinc-200">
        <span>Shopping cart is empty.</span>
        <span class="text-[9px] opacity-75">Click products on the left panel to add them.</span>
      </div>
    {:else}
      {#each cart as item (item.product.id)}
        <div class="flex items-start justify-between gap-3 border-b border-dashed border-zinc-200 pb-2">
          <div class="flex-grow flex flex-col">
            <span class="font-display font-bold text-xs leading-tight line-clamp-2">{item.product.name}</span>
            <span class="font-mono text-[10px] text-zinc-500 mt-1">
              IDR {item.product.price.toLocaleString('id-ID')} / pc
            </span>
          </div>

          <div class="flex flex-col items-end gap-1 shrink-0">
            <!-- Qty adjusts -->
            <div class="flex items-center gap-1.5 border-2 border-neubrutalism-charcoal px-1 bg-zinc-100 shadow-neubrutalism-sm text-[10px]">
              <button type="button" class="p-0.5 hover:bg-zinc-200" onclick={() => onAdjustQty(item.product.id, -1)}>
                <Minus class="w-3 h-3" />
              </button>
              <span class="font-mono font-bold w-4 text-center">{item.qty}</span>
              <button type="button" class="p-0.5 hover:bg-zinc-200" onclick={() => onAdjustQty(item.product.id, 1)}>
                <Plus class="w-3 h-3" />
              </button>
            </div>

            <span class="font-mono text-xs font-bold mt-1 text-neubrutalism-charcoal">
              IDR {(item.product.price * item.qty).toLocaleString('id-ID')}
            </span>
          </div>
        </div>
      {/each}
    {/if}
  </div>

  <!-- Calculations -->
  {#if cart.length > 0}
    <div class="border-t-4 border-neubrutalism-charcoal pt-4 flex flex-col gap-2.5 font-mono text-xs">
      <!-- Subtotal -->
      <div class="flex justify-between font-bold">
        <span>SUBTOTAL:</span>
        <span>IDR {subtotal.toLocaleString('id-ID')}</span>
      </div>

      <!-- Discount -->
      <div class="flex items-center justify-between gap-4">
        <span>DISCOUNT (IDR):</span>
        <input 
          type="text" 
          placeholder="0"
          value={displayDiscount}
          oninput={handleDiscountInput}
          class="w-32 text-right border-2 border-neubrutalism-charcoal bg-white px-2 py-0.5 focus:outline-none focus:placeholder-transparent font-mono text-xs"
        />
      </div>

      <!-- Final Total -->
      <div class="flex justify-between font-display font-black text-base md:text-lg border-t-2 border-dashed border-zinc-300 pt-2">
        <span>TOTAL BILL:</span>
        <span class="text-neubrutalism-pink font-bold">IDR {finalTotal.toLocaleString('id-ID')}</span>
      </div>

      <!-- Payment Method -->
      <div class="flex items-center justify-between gap-4 border-t-2 border-dashed border-zinc-300 pt-2.5">
        <span>PAYMENT METHOD:</span>
        <div class="flex gap-2">
          <button 
            type="button"
            class="px-3 py-1 text-[10px] font-bold border-2 border-neubrutalism-charcoal shadow-neubrutalism-sm hover:bg-zinc-50 active:translate-y-0.5 active:shadow-none {paymentMethod === 'cash' ? 'bg-neubrutalism-yellow' : 'bg-white'}"
            onclick={() => { paymentMethod = 'cash'; }}
          >
            CASH
          </button>
          <button 
            type="button"
            class="px-3 py-1 text-[10px] font-bold border-2 border-neubrutalism-charcoal shadow-neubrutalism-sm hover:bg-zinc-50 active:translate-y-0.5 active:shadow-none {paymentMethod === 'qris' ? 'bg-neubrutalism-green' : 'bg-white'}"
            onclick={() => { paymentMethod = 'qris'; }}
          >
            QRIS Digital
          </button>
        </div>
      </div>

      <!-- Cash Paid input (Only cash) -->
      {#if paymentMethod === 'cash'}
        <div class="flex items-center justify-between gap-4">
          <span>CASH PAID (IDR):</span>
          <input 
            type="text" 
            placeholder="0"
            value={displayCashPaid}
            oninput={handleCashPaidInput}
            class="w-32 text-right border-2 border-neubrutalism-charcoal bg-white px-2 py-0.5 focus:outline-none focus:placeholder-transparent font-mono text-xs"
          />
        </div>

        <!-- Change return -->
        <div class="flex justify-between font-bold text-zinc-700 bg-zinc-100 p-2 border border-zinc-200">
          <span>CHANGE RETURN:</span>
          <span class="text-emerald-600 font-bold">IDR {changeAmount.toLocaleString('id-ID')}</span>
        </div>
      {/if}

      <!-- Checkout / Process Action -->
      <Button 
        bgColor="bg-neubrutalism-green" 
        onclick={onCheckout}
        class="w-full py-3 mt-2 font-display font-extrabold uppercase border-2 border-neubrutalism-charcoal shadow-neubrutalism-sm flex items-center justify-center gap-2"
      >
        <span>PROCESS CHECKOUT</span>
        <ArrowRight class="w-4 h-4" />
      </Button>
    </div>
  {/if}
</Card>
