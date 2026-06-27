<script lang="ts">
  import { Button } from '$lib';
  import { ticketService, type Ticket } from '$lib/services/ticket';
  import { page } from '$app/state'; // SvelteKit v2 / Svelte 5 state-based routing
  import { goto } from '$app/navigation';
  import { onMount } from 'svelte';
  import { ShieldAlert, ShieldCheck, ArrowLeft, Save, RotateCcw } from 'lucide-svelte';

  import TicketHeader from './components/TicketHeader.svelte';
  import CustomerDeviceDetails from '../components/CustomerDeviceDetails.svelte';
  import DiagnosticsTechActions from '../components/DiagnosticsTechActions.svelte';
  import TicketStatusController from './components/TicketStatusController.svelte';
  import PaymentInfoController from './components/PaymentInfoController.svelte';
  import WarrantyDetailsCard from './components/WarrantyDetailsCard.svelte';
  import EmergencyConfirmModal from './components/EmergencyConfirmModal.svelte';

  let ticketId = $derived(page.params.id || '');
  let ticket = $state<Ticket | null>(null);
  let isLoading = $state(true);
  let isEditing = $state(false);
  let showConfirmModal = $state(false);

  // Edit form states
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
        // Seed states
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

    if (devicePosition === 'picked_up') {
      if (status !== 'completed' && status !== 'cancelled') {
        errorMessage = 'Service status must be COMPLETED or CANCELLED when device location status is PICKED UP.';
        isSubmitting = false;
        return;
      }
      if (status === 'completed') {
        if (paymentStatus !== 'paid') {
          errorMessage = 'Payment status must be set to PAID when status is COMPLETED.';
          isSubmitting = false;
          return;
        }
        if (!paymentMethod || (paymentMethod !== 'cash' && paymentMethod !== 'qris')) {
          errorMessage = 'Payment method must be set to CASH or QRIS when status is COMPLETED.';
          isSubmitting = false;
          return;
        }
      }
    }

    try {
      const updated = await ticketService.updateTicket(ticketId, {
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
        payment_method: paymentMethod,
        warranty_duration_days: warrantyDurationDays
      });

      ticket = updated;
      
      // Update form fields to match database trigger overrides (like picked_up automatic rules)
      devicePosition = updated.device_position;
      status = updated.status;
      paymentStatus = updated.payment_status;
      if (updated.payment_method) {
        paymentMethod = updated.payment_method;
      }

      successMessage = 'Repair ticket updated successfully.';
      isEditing = false;
      window.scrollTo({ top: 0, behavior: 'smooth' });
    } catch (err: any) {
      errorMessage = err.message || 'Failed to update ticket.';
    } finally {
      isSubmitting = false;
    }
  }

  function discardChanges() {
    if (ticket) {
      customerName = ticket.customer_name;
      customerPhone = ticket.customer_phone;
      brandPhone = ticket.brand_phone;
      modelPhone = ticket.model_phone;
      serialNumber = ticket.serial_number;
      damageDescription = ticket.damage_description;
      repairAction = ticket.repair_action || '';
      cost = ticket.cost;
      status = ticket.status;
      devicePosition = ticket.device_position;
      paymentStatus = ticket.payment_status;
      paymentMethod = ticket.payment_method;
      warrantyDurationDays = ticket.warranty_duration_days !== undefined ? ticket.warranty_duration_days : 30;
    }
    isEditing = false;
  }

  const getStatusColor = (statusVal: string) => {
    switch (statusVal) {
      case 'ready_for_pickup': return 'bg-neubrutalism-green';
      case 'in_repair': return 'bg-neubrutalism-yellow';
      case 'received': return 'bg-zinc-200';
      case 'completed': return 'bg-zinc-100 text-zinc-500';
      case 'cancelled': return 'bg-rose-100 text-rose-600 border-rose-300';
      default: return 'bg-white';
    }
  };

  const getStatusText = (statusVal: string) => {
    return statusVal.replace(/_/g, ' ');
  };

  const formatDate = (dateStr?: string) => {
    if (!dateStr) return 'N/A';
    const d = new Date(dateStr);
    return d.toLocaleDateString('en-US', { day: 'numeric', month: 'short', year: 'numeric', hour: '2-digit', minute: '2-digit' });
  };
</script>

<svelte:head>
  <title>Ticket Details - OpenBench Admin</title>
</svelte:head>

<div class="flex flex-col gap-6 max-w-4xl mx-auto">
  
  <!-- Back Action -->
  <div class="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4 w-full">
    <a href="/admin/tickets" class="inline-flex items-center gap-1.5 font-mono text-xs font-bold uppercase hover:underline">
      <ArrowLeft class="w-4 h-4" />
      <span>BACK TO TICKETS</span>
    </a>

    {#if !isEditing}
      <div class="flex gap-2">
        {#if devicePosition !== 'picked_up'}
          <Button
            bgColor="bg-neubrutalism-yellow"
            onclick={() => isEditing = true}
            class="py-1.5 px-3 text-xs font-bold flex items-center gap-1.5 shadow-neubrutalism-sm border-2 border-neubrutalism-charcoal"
          >
            <span>EDIT TICKET</span>
          </Button>
        {/if}
        <Button
          bgColor="bg-neubrutalism-pink"
          onclick={() => showConfirmModal = true}
          class="py-1.5 px-3 text-xs font-bold flex items-center gap-1.5 shadow-neubrutalism-sm text-white border-2 border-neubrutalism-charcoal"
        >
          <span class="text-white">EMERGENCY EDIT</span>
        </Button>
      </div>
    {:else}
      <div class="font-mono text-[10px] font-extrabold bg-neubrutalism-yellow text-neubrutalism-charcoal border-2 border-neubrutalism-charcoal px-3 py-1.5 shadow-neubrutalism-sm uppercase tracking-wide animate-pulse">
        EDITING ACTIVE
      </div>
    {/if}
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
    <!-- Heading -->
    <TicketHeader
      {ticket}
      {getStatusColor}
      {getStatusText}
      {formatDate}
    />

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

    <!-- Detail layout -->
    <form onsubmit={handleSave} class="grid grid-cols-1 lg:grid-cols-3 gap-8">
      
      <!-- Left columns: Edit Fields -->
      <div class="lg:col-span-2 flex flex-col gap-6">
        <CustomerDeviceDetails
          bind:customerName
          bind:customerPhone
          bind:brandPhone
          bind:modelPhone
          bind:serialNumber
          {isSubmitting}
          {isEditing}
        />

        <DiagnosticsTechActions
          bind:damageDescription
          bind:repairAction
          {isSubmitting}
          {isEditing}
        />
      </div>

      <!-- Right column: Status & Actions -->
      <div class="flex flex-col gap-6">
        <TicketStatusController
          bind:status
          bind:devicePosition
          bind:warrantyDurationDays
          {isSubmitting}
          {isEditing}
        />

        <PaymentInfoController
          bind:cost
          bind:paymentStatus
          bind:paymentMethod
          {isSubmitting}
          {status}
          {isEditing}
        />

        <!-- Warranty display (Read-Only) -->
        {#if ticket.warranty_expiry_date || devicePosition === 'picked_up'}
          <WarrantyDetailsCard
            {ticket}
            {status}
            {formatDate}
          />
        {/if}

        <!-- Save / Cancel Buttons -->
        {#if isEditing}
          <div class="flex flex-col gap-3">
            <Button 
              bgColor="bg-neubrutalism-yellow" 
              type="submit"
              disabled={isSubmitting}
              class="w-full py-3 px-6 font-bold shadow-neubrutalism-md border-2 border-neubrutalism-charcoal flex items-center justify-center gap-2"
            >
              <Save class="w-4 h-4" />
              <span>{isSubmitting ? 'SAVING CHANGES...' : 'SAVE CHANGES'}</span>
            </Button>

            <Button 
              bgColor="bg-zinc-200" 
              type="button"
              disabled={isSubmitting}
              onclick={discardChanges}
              class="w-full py-2.5 px-6 font-bold shadow-neubrutalism-sm border-2 border-neubrutalism-charcoal flex items-center justify-center gap-2"
            >
              <RotateCcw class="w-4 h-4" />
              <span>CANCEL</span>
            </Button>
          </div>
        {/if}

      </div>

    </form>
  {/if}

  <!-- Emergency Edit Confirmation Modal -->
  <EmergencyConfirmModal bind:isOpen={showConfirmModal} {ticketId} />

</div>
