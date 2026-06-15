<script lang="ts">
  import { Card, Button, Input } from '$lib';
  import { ticketService } from '$lib/services/ticket';
  import { goto } from '$app/navigation';
  import { formatCurrencyInput, parseCurrencyInput } from '$lib/utils/format';
  import { Wrench, User, Phone, Tag, ShieldAlert, ArrowLeft } from 'lucide-svelte';

  // Form states using Svelte 5 state
  let customerName = $state('');
  let customerPhone = $state('');
  let brandPhone = $state('');
  let modelPhone = $state('');
  let serialNumber = $state('');
  let damageDescription = $state('');
  let cost = $state<number>(0);
  let displayCost = $state(formatCurrencyInput(0));

  function handleCostInput(e: Event) {
    const target = e.target as HTMLInputElement;
    const numeric = parseCurrencyInput(target.value);
    cost = numeric;
    displayCost = formatCurrencyInput(numeric);
  }

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
        serial_number: serialNumber.trim() || 'N/A',
        damage_description: damageDescription.trim(),
        repair_action: '',
        cost: cost,
        status: 'received',
        device_position: 'warehouse',
        payment_status: 'none'
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
      
      <!-- Section: Customer Details -->
      <div class="border-b-2 border-dashed border-zinc-200 pb-4">
        <h3 class="font-display font-bold text-sm uppercase text-zinc-700 mb-3 flex items-center gap-1.5">
          <User class="w-4 h-4" />
          <span>CUSTOMER INFORMATION</span>
        </h3>
        
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <Input 
            id="cust-name" 
            label="Customer Name"
            required={true}
            type="text" 
            placeholder="e.g. John Doe" 
            bind:value={customerName} 
            disabled={isSubmitting}
          />

          <Input 
            id="cust-phone" 
            label="Phone Number"
            required={true}
            type="text" 
            placeholder="e.g. 081234567890" 
            bind:value={customerPhone} 
            disabled={isSubmitting}
          />
        </div>
      </div>

      <!-- Section: Device Details -->
      <div class="border-b-2 border-dashed border-zinc-200 pb-4">
        <h3 class="font-display font-bold text-sm uppercase text-zinc-700 mb-3 flex items-center gap-1.5">
          <Wrench class="w-4 h-4" />
          <span>DEVICE INFORMATION</span>
        </h3>
        
        <div class="grid grid-cols-1 sm:grid-cols-3 gap-4">
          <Input 
            id="dev-brand" 
            label="Brand"
            required={true}
            type="text" 
            placeholder="e.g. Apple" 
            bind:value={brandPhone} 
            disabled={isSubmitting}
          />

          <Input 
            id="dev-model" 
            label="Model"
            required={true}
            type="text" 
            placeholder="e.g. iPhone 14 Pro" 
            bind:value={modelPhone} 
            disabled={isSubmitting}
          />

          <Input 
            id="dev-serial" 
            label="IMEI / Serial Number"
            type="text" 
            placeholder="e.g. SN-XYZ..." 
            bind:value={serialNumber} 
            disabled={isSubmitting}
          />
        </div>
      </div>

      <!-- Section: Diagnostics & Costs -->
      <div>
        <h3 class="font-display font-bold text-sm uppercase text-zinc-700 mb-3 flex items-center gap-1.5">
          <Tag class="w-4 h-4" />
          <span>REPAIR DETAILS & ESTIMATE</span>
        </h3>
        
        <div class="flex flex-col gap-4">
          <div class="flex flex-col gap-1.5">
            <label for="damage" class="font-display font-bold text-sm text-neubrutalism-charcoal uppercase tracking-wider">Damage Description</label>
            <textarea 
              id="damage" 
              placeholder="Describe what is wrong with the device..."
              bind:value={damageDescription}
              disabled={isSubmitting}
              rows="3"
              class="w-full border-4 border-neubrutalism-charcoal bg-white p-3 font-sans text-sm focus:outline-none focus:bg-zinc-50 focus:placeholder-transparent shadow-neubrutalism-sm"
            ></textarea>
          </div>

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
        </div>
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
