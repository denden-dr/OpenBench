<script lang="ts">
  import { ticketService, type Ticket } from '$lib/services/ticket';
  import { page } from '$app/state';
  import { goto } from '$app/navigation';
  import { onMount } from 'svelte';
  import { ShieldAlert, ShieldCheck, ArrowLeft } from 'lucide-svelte';
  import EmergencyWarningBanner from './components/EmergencyWarningBanner.svelte';
  import EmergencyEditForm from './components/EmergencyEditForm.svelte';

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
      await ticketService.emergencyUpdateTicket(ticketId, {
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
    <EmergencyWarningBanner />

    <!-- Form layout -->
    <EmergencyEditForm
      bind:customerName
      bind:customerPhone
      bind:brandPhone
      bind:modelPhone
      bind:serialNumber
      bind:damageDescription
      bind:repairAction
      bind:status
      bind:devicePosition
      bind:warrantyDurationDays
      bind:cost
      bind:paymentStatus
      bind:paymentMethod
      {isSubmitting}
      onsubmit={handleSave}
      oncancel={handleCancel}
    />
  {/if}
</div>
