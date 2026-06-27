<script lang="ts">
  import { Card, Button } from '$lib';
  import { ticketService } from '$lib/services/ticket';
  import { goto } from '$app/navigation';
  import { ShieldAlert, ArrowLeft } from 'lucide-svelte';
  
  import CustomerDeviceDetails from '../components/CustomerDeviceDetails.svelte';
  import DiagnosticsTechActions from '../components/DiagnosticsTechActions.svelte';
  import NewTicketPriceWarranty from './components/NewTicketPriceWarranty.svelte';

  // Form states using Svelte 5 state
  let customerName = $state('');
  let customerPhone = $state('');
  let brandPhone = $state('');
  let modelPhone = $state('');
  let serialNumber = $state('');
  let damageDescription = $state('');
  let repairAction = $state('');
  let cost = $state<number>(0);
  let warrantyDurationDays = $state(30);

  let errorMessage = $state('');
  let isSubmitting = $state(false);

  async function handleSubmit(event: Event) {
    event.preventDefault();
    errorMessage = '';
    isSubmitting = true;

    // Basic Client-Side Validations
    if (!customerName.trim()) {
      errorMessage = 'Customer name is required.';
      isSubmitting = false;
      return;
    }
    if (!customerPhone.trim() || !customerPhone.startsWith('08') || customerPhone.length < 9) {
      errorMessage = 'Invalid phone number (must start with 08 and be at least 9 digits).';
      isSubmitting = false;
      return;
    }
    if (!brandPhone.trim()) {
      errorMessage = 'Phone brand is required.';
      isSubmitting = false;
      return;
    }
    if (!modelPhone.trim()) {
      errorMessage = 'Phone model is required.';
      isSubmitting = false;
      return;
    }
    if (cost < 0) {
      errorMessage = 'Estimated cost cannot be negative.';
      isSubmitting = false;
      return;
    }

    try {
      await ticketService.createTicket({
        customer_name: customerName.trim(),
        customer_phone: customerPhone.trim(),
        brand_phone: brandPhone.trim(),
        model_phone: modelPhone.trim(),
        serial_number: serialNumber.trim() || undefined,
        damage_description: damageDescription.trim(),
        repair_action: repairAction.trim() || undefined,
        cost: cost,
        warranty_duration_days: warrantyDurationDays
      });

      // Redirect back to tickets index
      await goto('/admin/tickets');
    } catch (err: any) {
      errorMessage = err.message || 'Failed to save new ticket.';
    } finally {
      isSubmitting = false;
    }
  }
</script>

<svelte:head>
  <title>New Ticket - OpenBench Admin</title>
</svelte:head>

<div class="flex flex-col gap-6 max-w-2xl mx-auto">
  
  <!-- Back Action -->
  <div>
    <a href="/admin/tickets" class="inline-flex items-center gap-1.5 font-mono text-xs font-bold uppercase hover:underline">
      <ArrowLeft class="w-4 h-4" />
      <span>BACK TO TICKETS</span>
    </a>
  </div>

  <!-- Heading -->
  <div>
    <h2 class="font-display font-extrabold text-2xl md:text-3xl uppercase tracking-tight">
      Create New Repair Ticket
    </h2>
    <p class="font-sans text-xs sm:text-sm text-zinc-500 font-semibold mt-1">
      Register a customer device, diagnose details, and estimate repair cost.
    </p>
  </div>

  {#if errorMessage}
    <div class="bg-neubrutalism-pink text-white border-4 border-neubrutalism-charcoal p-4 font-mono text-xs flex items-center gap-2 shadow-neubrutalism-sm">
      <ShieldAlert class="w-5 h-5 shrink-0" />
      <span>{errorMessage}</span>
    </div>
  {/if}

  <!-- Form Card -->
  <Card bgColor="bg-white" class="border-4 border-neubrutalism-charcoal p-6 shadow-neubrutalism-md">
    <form onsubmit={handleSubmit} class="flex flex-col gap-5">
      
      <!-- Section: Customer & Device Details (Reusable Components) -->
      <CustomerDeviceDetails
        bind:customerName
        bind:customerPhone
        bind:brandPhone
        bind:modelPhone
        bind:serialNumber
        {isSubmitting}
        isEditing={true}
      />

      <DiagnosticsTechActions
        bind:damageDescription
        bind:repairAction
        {isSubmitting}
        isEditing={true}
      />

      <!-- Section: Pricing & Warranty (Modular Component) -->
      <div class="border-t-2 border-dashed border-zinc-200 pt-4">
        <h3 class="font-display font-bold text-sm uppercase text-zinc-700 mb-3 flex items-center gap-1.5">
          <span>REPAIR ESTIMATE & WARRANTY</span>
        </h3>
        
        <NewTicketPriceWarranty
          bind:cost
          bind:warrantyDurationDays
          {isSubmitting}
        />
      </div>

      <!-- Submit Actions -->
      <div class="flex items-center justify-end gap-3 mt-4 border-t-4 border-neubrutalism-charcoal pt-4">
        <a href="/admin/tickets" class="w-full sm:w-auto">
          <Button 
            bgColor="bg-white" 
            type="button"
            disabled={isSubmitting}
            class="w-full sm:w-auto py-2.5 px-6 font-bold shadow-neubrutalism-sm border-2 border-neubrutalism-charcoal"
          >
            CANCEL
          </Button>
        </a>

        <Button 
          bgColor="bg-neubrutalism-green" 
          type="submit"
          disabled={isSubmitting}
          class="w-full sm:w-auto py-2.5 px-6 font-bold shadow-neubrutalism-sm border-2 border-neubrutalism-charcoal flex items-center justify-center gap-2"
        >
          <span>{isSubmitting ? 'CREATING...' : 'CREATE TICKET'}</span>
        </Button>
      </div>

    </form>
  </Card>

</div>
