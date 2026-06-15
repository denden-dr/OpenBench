<script lang="ts">
  import { Card } from '$lib';
  import { CreditCard } from 'lucide-svelte';
  import { formatCurrencyInput, parseCurrencyInput } from '$lib/utils/format';

  let {
    cost = $bindable(0),
    paymentStatus = $bindable('none'),
    paymentMethod = $bindable(undefined),
    isSubmitting,
    status,
    isEditing = false
  }: {
    cost: number;
    paymentStatus: 'none' | 'requesting' | 'paid';
    paymentMethod: 'cash' | 'qris' | undefined;
    isSubmitting: boolean;
    status: string;
    isEditing?: boolean;
  } = $props();

  let displayCost = $state(formatCurrencyInput(cost));

  $effect(() => {
    displayCost = formatCurrencyInput(cost);
  });

  function handleCostInput(e: Event) {
    const target = e.target as HTMLInputElement;
    const numeric = parseCurrencyInput(target.value);
    cost = numeric;
    displayCost = formatCurrencyInput(numeric);
  }
</script>

<Card bgColor="bg-white" class="border-4 border-neubrutalism-charcoal p-6 shadow-neubrutalism-md flex flex-col gap-4">
  <h3 class="font-display font-bold text-sm uppercase text-zinc-700 border-b-2 border-dashed border-zinc-200 pb-2 flex items-center gap-1.5">
    <CreditCard class="w-4 h-4 text-neubrutalism-charcoal" />
    <span>Payment Info</span>
  </h3>

  <div class="flex flex-col gap-2 w-full">
    <label for="cost-input" class="font-display font-bold text-sm text-neubrutalism-charcoal uppercase tracking-wider">Cost (IDR)</label>
    <input 
      id="cost-input" 
      type="text" 
      placeholder="0"
      value={displayCost} 
      oninput={handleCostInput}
      disabled={!isEditing || isSubmitting} 
      class="w-full border-4 border-neubrutalism-charcoal bg-white p-3 font-sans text-neubrutalism-charcoal rounded-none transition-all duration-150 focus:outline-none focus:ring-4 focus:ring-neubrutalism-charcoal focus:bg-[#fefefe] focus:placeholder-transparent disabled:opacity-50 disabled:cursor-not-allowed disabled:bg-zinc-100 disabled:border-dashed"
    />
  </div>

  <div class="flex flex-col gap-1.5">
    <label for="pay-status" class="font-mono text-[10px] font-bold uppercase text-zinc-500">Payment Status</label>
    <select 
      id="pay-status" 
      bind:value={paymentStatus}
      disabled={!isEditing || isSubmitting || status === 'picked_up'}
      class="w-full bg-white border-4 border-neubrutalism-charcoal px-3 py-2 font-mono text-xs shadow-neubrutalism-sm focus:outline-none"
    >
      <option value="none">UNPAID</option>
      <option value="requesting">BILL REQUESTED</option>
      <option value="paid">PAID</option>
    </select>
  </div>

  <div class="flex flex-col gap-1.5">
    <label for="pay-method" class="font-mono text-[10px] font-bold uppercase text-zinc-500">Payment Method</label>
    <select 
      id="pay-method" 
      bind:value={paymentMethod}
      disabled={!isEditing || isSubmitting}
      class="w-full bg-white border-4 border-neubrutalism-charcoal px-3 py-2 font-mono text-xs shadow-neubrutalism-sm focus:outline-none"
    >
      <option value={undefined}>NOT SELECTED</option>
      <option value="cash">CASH</option>
      <option value="qris">QRIS DIGITAL</option>
    </select>
  </div>
</Card>
