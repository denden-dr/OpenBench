<script lang="ts">
  import { Button } from '$lib';
  import { ticketService, type Ticket } from '$lib/services/ticket';
  import { page } from '$app/state';
  import { goto } from '$app/navigation';
  import { onMount } from 'svelte';
  import { ShieldAlert, ShieldCheck, ArrowLeft, Save, RotateCcw } from 'lucide-svelte';

  let ticketId = $derived(page.params.id || '');
  let ticket = $state<Ticket | null>(null);
  let isLoading = $state(true);

  // Form states
  let customerName = $state('');
  let customerPhone = $state('');
  let brandPhone = $state('');
  let modelPhone = $state('');
  let serialNumber = $state('');
  let damageDescription = $state('');
  let repairAction = $state('');
  let cost = $state<number>(0);
  let status = $state<'received' | 'in_repair' | 'completed' | 'cancelled'>('received');
  let devicePosition = $state<'warehouse' | 'picked_up'>('warehouse');
  let paymentStatus = $state<'none' | 'requesting' | 'paid'>('none');
  let paymentMethod = $state<'cash' | 'qris' | undefined>(undefined);
  let warrantyDurationDays = $state(30);

  let successMessage = $state('');
  let errorMessage = $state('');
  let isSubmitting = $state(false);

  onMount(async () => {
    try {
      const data = await ticketService.getTicket(ticketId);
      if (data) {
        ticket = data;
        customerName = data.customer_name;
        customerPhone = data.customer_phone;
        brandPhone = data.brand_phone;
        modelPhone = data.model_phone;
        serialNumber = data.serial_number;
        damageDescription = data.damage_description;
        repairAction = data.repair_action || '';
        cost = data.cost;
        status = data.status;
        devicePosition = data.device_position;
        paymentStatus = data.payment_status;
        paymentMethod = data.payment_method;
        warrantyDurationDays = data.warranty_duration_days !== undefined ? data.warranty_duration_days : 30;
      } else {
        errorMessage = 'Repair ticket not found.';
      }
    } catch (err: any) {
      errorMessage = err.message || 'Failed to load ticket details.';
    } finally {
      isLoading = false;
    }
  });

  async function handleSave(event: Event) {
    event.preventDefault();
    successMessage = '';
    errorMessage = '';
    isSubmitting = true;

    // Validate inputs
    if (!customerName.trim() || !brandPhone.trim() || !modelPhone.trim()) {
      errorMessage = 'Customer name, brand, and phone model are required.';
      isSubmitting = false;
      return;
    }

    try {
      const updated = await ticketService.emergencyUpdateTicket(ticketId, {
        customer_name: customerName.trim(),
        customer_phone: customerPhone.trim(),
        brand_phone: brandPhone.trim(),
        model_phone: modelPhone.trim(),
        serial_number: serialNumber.trim(),
        damage_description: damageDescription.trim(),
        repair_action: repairAction.trim(),
        cost: cost,
        status: status,
        device_position: devicePosition,
        payment_status: paymentStatus,
        payment_method: paymentMethod || null as any,
        warranty_duration_days: warrantyDurationDays
      });

      successMessage = 'Ticket updated successfully under Emergency Override. Redirecting to details...';
      setTimeout(() => {
        goto(`/admin/tickets/${ticketId}`);
      }, 1500);
    } catch (err: any) {
      errorMessage = err.message || 'Failed to save emergency update.';
      isSubmitting = false;
    }
  }

  function handleCancel() {
    goto(`/admin/tickets/${ticketId}`);
  }
</script>

<svelte:head>
  <title>Emergency Edit - OpenBench Admin</title>
</svelte:head>

<div class="flex flex-col gap-6 max-w-4xl mx-auto">
  <!-- Back Action -->
  <div class="flex items-center justify-between w-full">
    <button onclick={handleCancel} class="inline-flex items-center gap-1.5 font-mono text-xs font-bold uppercase hover:underline">
      <ArrowLeft class="w-4 h-4" />
      <span>BACK TO TICKET DETAILS</span>
    </button>
    <div class="font-mono text-[10px] font-extrabold bg-neubrutalism-pink text-white border-2 border-neubrutalism-charcoal px-3 py-1.5 shadow-neubrutalism-sm uppercase tracking-wide">
      EMERGENCY OVERRIDE ACTIVE
    </div>
  </div>

  {#if isLoading}
    <div class="h-64 bg-zinc-200 border-4 border-neubrutalism-charcoal animate-pulse flex items-center justify-center font-mono text-sm uppercase">
      LOADING TICKET DETAILS...
    </div>
  {:else if ticket === null}
    <div class="bg-neubrutalism-pink text-white border-4 border-neubrutalism-charcoal p-6 font-mono text-sm flex items-center gap-3 shadow-neubrutalism-sm">
      <ShieldAlert class="w-6 h-6 shrink-0" />
      <span>{errorMessage || 'TICKET NOT FOUND.'}</span>
    </div>
  {:else}
    <!-- Alert Messages -->
    {#if successMessage}
      <div class="bg-neubrutalism-green border-4 border-neubrutalism-charcoal p-4 font-mono text-xs flex items-center gap-2 shadow-neubrutalism-sm">
        <ShieldCheck class="w-5 h-5 shrink-0" />
        <span>{successMessage}</span>
      </div>
    {/if}

    {#if errorMessage}
      <div class="bg-neubrutalism-pink text-white border-4 border-neubrutalism-charcoal p-4 font-mono text-xs flex items-center gap-2 shadow-neubrutalism-sm">
        <ShieldAlert class="w-5 h-5 shrink-0" />
        <span>{errorMessage}</span>
      </div>
    {/if}

    <!-- Warning Banner -->
    <div class="bg-amber-100 border-4 border-neubrutalism-charcoal p-6 flex flex-col gap-2 shadow-neubrutalism-sm">
      <div class="flex items-center gap-2 text-amber-800 font-display font-extrabold text-sm uppercase">
        <ShieldAlert class="w-5 h-5" />
        <span>WARNING: Emergency Bypass Operation</span>
      </div>
      <p class="font-sans text-xs text-amber-900 leading-relaxed font-semibold">
        You are making modifications that bypass standard business rules. If you revert the status from "Picked Up" to "In Warehouse/Store", any active warranty records associated with this ticket will be automatically voided and deleted from the system.
      </p>
    </div>

    <!-- Form layout -->
    <form onsubmit={handleSave} class="grid grid-cols-1 lg:grid-cols-3 gap-8">
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
              <input id="cost-input" type="number" bind:value={cost} min="0" class="w-full bg-white border-4 border-neubrutalism-charcoal px-3 py-2 font-mono text-xs shadow-neubrutalism-sm focus:outline-none" />
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
            onclick={handleCancel}
            class="w-full py-2.5 px-6 font-bold shadow-neubrutalism-sm border-2 border-neubrutalism-charcoal flex items-center justify-center gap-2"
          >
            <RotateCcw class="w-4 h-4" />
            <span>CANCEL / GO BACK</span>
          </Button>
        </div>
      </div>
    </form>
  {/if}
</div>
