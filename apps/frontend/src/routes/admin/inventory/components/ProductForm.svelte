<script lang="ts">
  import { Card, Button, Input } from '$lib';
  import { formatCurrencyInput, parseCurrencyInput } from '$lib/utils/format';

  let {
    editingId,
    prodName = $bindable(''),
    prodCategory = $bindable('retail'),
    prodCostPrice = $bindable(0),
    prodPrice = $bindable(0),
    prodStock = $bindable(0),
    prodMinStock = $bindable(5),
    onSubmit,
    onCancel
  }: {
    editingId: string | null;
    prodName: string;
    prodCategory: 'retail' | 'spare_part';
    prodCostPrice: number;
    prodPrice: number;
    prodStock: number;
    prodMinStock: number;
    onSubmit: (e: Event) => void;
    onCancel: () => void;
  } = $props();

  let profitMargin = $derived(prodPrice - prodCostPrice);
  let profitMarginPercent = $derived(prodPrice > 0 ? Math.round((profitMargin / prodPrice) * 100) : 0);

  // Dynamic formatted display values
  let displayCostPrice = $state(formatCurrencyInput(prodCostPrice));
  let displayPrice = $state(formatCurrencyInput(prodPrice));

  $effect(() => {
    displayCostPrice = formatCurrencyInput(prodCostPrice);
  });

  $effect(() => {
    displayPrice = formatCurrencyInput(prodPrice);
  });

  function handleCostInput(e: Event) {
    const target = e.target as HTMLInputElement;
    const numeric = parseCurrencyInput(target.value);
    prodCostPrice = numeric;
    displayCostPrice = formatCurrencyInput(numeric);
  }

  function handlePriceInput(e: Event) {
    const target = e.target as HTMLInputElement;
    const numeric = parseCurrencyInput(target.value);
    prodPrice = numeric;
    displayPrice = formatCurrencyInput(numeric);
  }
</script>

<Card bgColor="bg-white" class="border-4 border-neubrutalism-charcoal p-6 shadow-neubrutalism-md">
  <form onsubmit={onSubmit} class="flex flex-col gap-4">
    <h3 class="font-display font-bold text-sm uppercase text-zinc-700 border-b-2 border-dashed border-zinc-200 pb-2">
      {editingId ? 'Edit Catalog Product' : 'Add New Product / Component'}
    </h3>
    
    <div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-4">
      <!-- Name -->
      <div class="sm:col-span-2">
        <Input id="p-name" label="Item Name" type="text" placeholder="e.g. Battery Replacement iPhone 13" bind:value={prodName} />
      </div>

      <!-- Category -->
      <div class="flex flex-col gap-1.5">
        <label for="p-category" class="font-mono text-[10px] font-bold uppercase text-zinc-500">Category</label>
        <select 
          id="p-category" 
          bind:value={prodCategory}
          class="w-full bg-white border-4 border-neubrutalism-charcoal px-3 py-2 font-mono text-xs shadow-neubrutalism-sm focus:outline-none"
        >
          <option value="retail">RETAIL ACCESSORIES (GENERAL SALE)</option>
          <option value="spare_part">SPARE PARTS (FOR SERVICES/REPAIRS)</option>
        </select>
      </div>

      <!-- Cost Price -->
      <div class="flex flex-col gap-2 w-full">
        <label for="p-cost" class="font-display font-bold text-sm text-neubrutalism-charcoal uppercase tracking-wider">Capital Cost Price (IDR)</label>
        <input 
          id="p-cost" 
          type="text" 
          placeholder="0"
          value={displayCostPrice} 
          oninput={handleCostInput}
          class="w-full border-4 border-neubrutalism-charcoal bg-white p-3 font-sans text-neubrutalism-charcoal rounded-none transition-all duration-150 focus:outline-none focus:ring-4 focus:ring-neubrutalism-charcoal focus:bg-[#fefefe] focus:placeholder-transparent"
        />
      </div>

      <!-- Sell Price -->
      <div class="flex flex-col gap-2 w-full">
        <label for="p-price" class="font-display font-bold text-sm text-neubrutalism-charcoal uppercase tracking-wider">Selling Price (IDR)</label>
        <input 
          id="p-price" 
          type="text" 
          placeholder="0"
          value={displayPrice} 
          oninput={handlePriceInput}
          class="w-full border-4 border-neubrutalism-charcoal bg-white p-3 font-sans text-neubrutalism-charcoal rounded-none transition-all duration-150 focus:outline-none focus:ring-4 focus:ring-neubrutalism-charcoal focus:bg-[#fefefe] focus:placeholder-transparent"
        />
      </div>

      <!-- Margin calculator preview -->
      <div class="flex flex-col gap-1.5 justify-end pb-1 font-mono text-xs font-bold text-zinc-700">
        <span>Estimated Profit Margin:</span>
        <span class="text-neubrutalism-pink text-sm">
          IDR {profitMargin.toLocaleString('id-ID')} ({profitMarginPercent}%)
        </span>
      </div>

      <!-- Stock -->
      <div class="flex flex-col gap-2 w-full">
        <label for="p-stock" class="font-display font-bold text-sm text-neubrutalism-charcoal uppercase tracking-wider">Initial Stock Quantity</label>
        <input 
          id="p-stock" 
          type="number" 
          min="0"
          placeholder="0"
          bind:value={prodStock} 
          class="w-full border-4 border-neubrutalism-charcoal bg-white p-3 font-sans text-neubrutalism-charcoal rounded-none transition-all duration-150 focus:outline-none focus:ring-4 focus:ring-neubrutalism-charcoal focus:bg-[#fefefe] focus:placeholder-transparent"
        />
      </div>

      <!-- Min Stock Alert level -->
      <div class="flex flex-col gap-2 w-full">
        <label for="p-min" class="font-display font-bold text-sm text-neubrutalism-charcoal uppercase tracking-wider">Min Stock Warning Level</label>
        <input 
          id="p-min" 
          type="number" 
          min="0"
          placeholder="5"
          bind:value={prodMinStock} 
          class="w-full border-4 border-neubrutalism-charcoal bg-white p-3 font-sans text-neubrutalism-charcoal rounded-none transition-all duration-150 focus:outline-none focus:ring-4 focus:ring-neubrutalism-charcoal focus:bg-[#fefefe] focus:placeholder-transparent"
        />
      </div>
    </div>

    <div class="flex justify-end gap-3 mt-2">
      <Button 
        bgColor="bg-white" 
        type="button" 
        onclick={onCancel}
        class="py-2 px-5 font-bold shadow-neubrutalism-sm border-2 border-neubrutalism-charcoal"
      >
        CANCEL
      </Button>

      <Button 
        bgColor="bg-neubrutalism-green" 
        type="submit"
        class="py-2 px-5 font-bold shadow-neubrutalism-sm border-2 border-neubrutalism-charcoal"
      >
        <span>{editingId ? 'UPDATE PRODUCT' : 'CREATE PRODUCT'}</span>
      </Button>
    </div>
  </form>
</Card>
