<script lang="ts">
  import { formatCurrencyInput, parseCurrencyInput } from '$lib/utils/format';

  interface Props {
    cost: number;
    warrantyDurationDays: number;
    isSubmitting: boolean;
  }

  let {
    cost = $bindable(),
    warrantyDurationDays = $bindable(),
    isSubmitting
  }: Props = $props();

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

<div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
  <div class="flex flex-col gap-2 w-full">
    <label for="est-cost" class="font-display font-bold text-sm text-neubrutalism-charcoal uppercase tracking-wider">Estimated Cost (IDR) *</label>
    <input
      id="est-cost"
      type="text"
      placeholder="0"
      required
      value={displayCost}
      oninput={handleCostInput}
      disabled={isSubmitting}
      class="w-full border-4 border-neubrutalism-charcoal bg-white p-3 font-sans text-neubrutalism-charcoal rounded-none transition-all duration-150 focus:outline-none focus:ring-4 focus:ring-neubrutalism-charcoal focus:bg-[#fefefe] focus:placeholder-transparent disabled:opacity-50 disabled:cursor-not-allowed disabled:bg-zinc-100 disabled:border-dashed"
    />
  </div>

  <div class="flex flex-col gap-2 w-full">
    <label for="warranty-select" class="font-display font-bold text-sm text-neubrutalism-charcoal uppercase tracking-wider">Warranty Duration</label>
    <select
      id="warranty-select"
      bind:value={warrantyDurationDays}
      disabled={isSubmitting}
      class="w-full border-4 border-neubrutalism-charcoal bg-white p-3 font-sans text-sm focus:outline-none shadow-neubrutalism-sm rounded-none h-[52px]"
    >
      <option value={0}>No Warranty</option>
      <option value={7}>7 Days</option>
      <option value={14}>14 Days</option>
      <option value={30}>30 Days</option>
      <option value={90}>90 Days</option>
      <option value={180}>180 Days</option>
    </select>
  </div>
</div>
