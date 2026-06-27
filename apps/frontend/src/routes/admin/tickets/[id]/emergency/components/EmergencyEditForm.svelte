<script lang="ts">
  import { Button } from '$lib';
  import { Save, RotateCcw } from 'lucide-svelte';

  interface Props {
    customerName: string;
    customerPhone: string;
    brandPhone: string;
    modelPhone: string;
    serialNumber: string;
    damageDescription: string;
    repairAction: string;
    status: 'received' | 'in_repair' | 'completed' | 'cancelled';
    devicePosition: 'warehouse' | 'picked_up';
    warrantyDurationDays: number;
    cost: number;
    paymentStatus: 'none' | 'requesting' | 'paid';
    paymentMethod: 'cash' | 'qris' | undefined;
    isSubmitting: boolean;
    onsubmit: (event: Event) => void;
    oncancel: () => void;
  }

  let {
    customerName = $bindable(),
    customerPhone = $bindable(),
    brandPhone = $bindable(),
    modelPhone = $bindable(),
    serialNumber = $bindable(),
    damageDescription = $bindable(),
    repairAction = $bindable(),
    status = $bindable(),
    devicePosition = $bindable(),
    warrantyDurationDays = $bindable(),
    cost = $bindable(),
    paymentStatus = $bindable(),
    paymentMethod = $bindable(),
    isSubmitting,
    onsubmit,
    oncancel
  }: Props = $props();

  import { formatCurrencyInput, parseCurrencyInput } from '$lib/utils/format';

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

<form {onsubmit} class="grid grid-cols-1 lg:grid-cols-3 gap-8">
  <!-- Left side: Customer & Device details -->
  <div class="lg:col-span-2 flex flex-col gap-6">
    <!-- Customer & Device -->
    <div class="bg-white border-4 border-neubrutalism-charcoal p-6 shadow-neubrutalism-md flex flex-col gap-4">
      <h3 class="font-display font-bold text-sm uppercase text-zinc-700 border-b-2 border-dashed border-zinc-200 pb-2">
        Customer & Device Information
      </h3>
      <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
        <div class="flex flex-col gap-1.5">
          <label for="cust-name" class="font-mono text-[10px] font-bold uppercase text-zinc-500">Customer Name</label>
          <input id="cust-name" type="text" bind:value={customerName} class="w-full bg-white border-4 border-neubrutalism-charcoal px-3 py-2 font-mono text-xs shadow-neubrutalism-sm focus:outline-none" required />
        </div>
        <div class="flex flex-col gap-1.5">
          <label for="cust-phone" class="font-mono text-[10px] font-bold uppercase text-zinc-500">Phone Number</label>
          <input id="cust-phone" type="text" bind:value={customerPhone} class="w-full bg-white border-4 border-neubrutalism-charcoal px-3 py-2 font-mono text-xs shadow-neubrutalism-sm focus:outline-none" required />
        </div>
        <div class="flex flex-col gap-1.5">
          <label for="dev-brand" class="font-mono text-[10px] font-bold uppercase text-zinc-500">Brand</label>
          <input id="dev-brand" type="text" bind:value={brandPhone} class="w-full bg-white border-4 border-neubrutalism-charcoal px-3 py-2 font-mono text-xs shadow-neubrutalism-sm focus:outline-none" required />
        </div>
        <div class="flex flex-col gap-1.5">
          <label for="dev-model" class="font-mono text-[10px] font-bold uppercase text-zinc-500">Model</label>
          <input id="dev-model" type="text" bind:value={modelPhone} class="w-full bg-white border-4 border-neubrutalism-charcoal px-3 py-2 font-mono text-xs shadow-neubrutalism-sm focus:outline-none" required />
        </div>
        <div class="sm:col-span-2 flex flex-col gap-1.5">
          <label for="dev-serial" class="font-mono text-[10px] font-bold uppercase text-zinc-500">IMEI / Serial Number</label>
          <input id="dev-serial" type="text" bind:value={serialNumber} class="w-full bg-white border-4 border-neubrutalism-charcoal px-3 py-2 font-mono text-xs shadow-neubrutalism-sm focus:outline-none" />
        </div>
      </div>
    </div>

    <!-- Diagnostics -->
    <div class="bg-white border-4 border-neubrutalism-charcoal p-6 shadow-neubrutalism-md flex flex-col gap-4">
      <h3 class="font-display font-bold text-sm uppercase text-zinc-700 border-b-2 border-dashed border-zinc-200 pb-2">
        Diagnostics & Repair Info
      </h3>
      <div class="flex flex-col gap-4">
        <div class="flex flex-col gap-1.5">
          <label for="damage" class="font-mono text-[10px] font-bold uppercase text-zinc-500">Damage Description</label>
          <textarea id="damage" bind:value={damageDescription} rows="3" class="w-full bg-white border-4 border-neubrutalism-charcoal px-3 py-2 font-mono text-xs shadow-neubrutalism-sm focus:outline-none" required></textarea>
        </div>
        <div class="flex flex-col gap-1.5">
          <label for="repair" class="font-mono text-[10px] font-bold uppercase text-zinc-500">Repair Action</label>
          <textarea id="repair" bind:value={repairAction} rows="3" class="w-full bg-white border-4 border-neubrutalism-charcoal px-3 py-2 font-mono text-xs shadow-neubrutalism-sm focus:outline-none"></textarea>
        </div>
      </div>
    </div>
  </div>

  <!-- Right side: Status and Payment overrides -->
  <div class="flex flex-col gap-6">
    <!-- Status Card -->
    <div class="bg-white border-4 border-neubrutalism-charcoal p-6 shadow-neubrutalism-md flex flex-col gap-4">
      <h3 class="font-display font-bold text-sm uppercase text-zinc-700 border-b-2 border-dashed border-zinc-200 pb-2">
        Status Override
      </h3>
      <div class="flex flex-col gap-3">
        <div class="flex flex-col gap-1.5">
          <label for="status-select" class="font-mono text-[10px] font-bold uppercase text-zinc-500">Service Status</label>
          <select id="status-select" bind:value={status} class="w-full bg-white border-4 border-neubrutalism-charcoal px-3 py-2 font-mono text-xs shadow-neubrutalism-sm focus:outline-none">
            <option value="received">RECEIVED</option>
            <option value="in_repair">IN REPAIR</option>
            <option value="completed">COMPLETED</option>
            <option value="cancelled">CANCELLED</option>
          </select>
        </div>

        <div class="flex flex-col gap-1.5">
          <label for="pos-select" class="font-mono text-[10px] font-bold uppercase text-zinc-500">Device Location Status</label>
          <select id="pos-select" bind:value={devicePosition} class="w-full bg-white border-4 border-neubrutalism-charcoal px-3 py-2 font-mono text-xs shadow-neubrutalism-sm focus:outline-none">
            <option value="warehouse">IN WAREHOUSE / STORE</option>
            <option value="picked_up">PICKED UP / TAKEN HOME</option>
          </select>
        </div>

        <div class="flex flex-col gap-1.5">
          <label for="warranty-select" class="font-mono text-[10px] font-bold uppercase text-zinc-500">Warranty Duration</label>
          <select id="warranty-select" bind:value={warrantyDurationDays} class="w-full bg-white border-4 border-neubrutalism-charcoal px-3 py-2 font-mono text-xs shadow-neubrutalism-sm focus:outline-none">
            <option value={0}>NO WARRANTY</option>
            <option value={7}>7 DAYS</option>
            <option value={14}>14 DAYS</option>
            <option value={30}>30 DAYS</option>
            <option value={90}>90 DAYS</option>
            <option value={180}>180 DAYS</option>
          </select>
        </div>
      </div>
    </div>

    <!-- Payment Card -->
    <div class="bg-white border-4 border-neubrutalism-charcoal p-6 shadow-neubrutalism-md flex flex-col gap-4">
      <h3 class="font-display font-bold text-sm uppercase text-zinc-700 border-b-2 border-dashed border-zinc-200 pb-2">
        Payment Override
      </h3>
      <div class="flex flex-col gap-3">
        <div class="flex flex-col gap-1.5">
          <label for="cost-input" class="font-mono text-[10px] font-bold uppercase text-zinc-500">Cost (IDR)</label>
          <input
            id="cost-input"
            type="text"
            value={displayCost}
            oninput={handleCostInput}
            disabled={isSubmitting}
            class="w-full bg-white border-4 border-neubrutalism-charcoal px-3 py-2 font-mono text-xs shadow-neubrutalism-sm focus:outline-none"
          />
        </div>

        <div class="flex flex-col gap-1.5">
          <label for="pay-status" class="font-mono text-[10px] font-bold uppercase text-zinc-500">Payment Status</label>
          <select id="pay-status" bind:value={paymentStatus} class="w-full bg-white border-4 border-neubrutalism-charcoal px-3 py-2 font-mono text-xs shadow-neubrutalism-sm focus:outline-none">
            <option value="none">UNPAID</option>
            <option value="requesting">BILLING INVOICE SENT</option>
            <option value="paid">PAID</option>
          </select>
        </div>

        <div class="flex flex-col gap-1.5">
          <label for="pay-method" class="font-mono text-[10px] font-bold uppercase text-zinc-500">Payment Method</label>
          <select id="pay-method" bind:value={paymentMethod} class="w-full bg-white border-4 border-neubrutalism-charcoal px-3 py-2 font-mono text-xs shadow-neubrutalism-sm focus:outline-none">
            <option value={undefined}>NONE</option>
            <option value="cash">CASH</option>
            <option value="qris">QRIS</option>
          </select>
        </div>
      </div>
    </div>

    <!-- Buttons -->
    <div class="flex flex-col gap-3">
      <Button 
        bgColor="bg-neubrutalism-pink" 
        type="submit"
        disabled={isSubmitting}
        class="w-full py-3 px-6 font-bold text-white shadow-neubrutalism-md border-2 border-neubrutalism-charcoal flex items-center justify-center gap-2"
      >
        <Save class="w-4 h-4 text-white" />
        <span class="text-white">{isSubmitting ? 'SAVING CHANGES...' : 'SAVE EMERGENCY OVERRIDES'}</span>
      </Button>

      <Button 
        bgColor="bg-zinc-200" 
        type="button"
        disabled={isSubmitting}
        onclick={oncancel}
        class="w-full py-2.5 px-6 font-bold shadow-neubrutalism-sm border-2 border-neubrutalism-charcoal flex items-center justify-center gap-2"
      >
        <RotateCcw class="w-4 h-4" />
        <span>CANCEL / GO BACK</span>
      </Button>
    </div>
  </div>
</form>
