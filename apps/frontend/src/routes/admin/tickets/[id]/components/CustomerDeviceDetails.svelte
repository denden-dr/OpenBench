<script lang="ts">
  import { Card, Input } from '$lib';
  import { User } from 'lucide-svelte';

  let {
    customerName = $bindable(''),
    customerPhone = $bindable(''),
    brandPhone = $bindable(''),
    modelPhone = $bindable(''),
    serialNumber = $bindable(''),
    isSubmitting,
    isEditing = false
  }: {
    customerName: string;
    customerPhone: string;
    brandPhone: string;
    modelPhone: string;
    serialNumber: string;
    isSubmitting: boolean;
    isEditing?: boolean;
  } = $props();
</script>

<Card bgColor="bg-white" class="border-4 border-neubrutalism-charcoal p-6 shadow-neubrutalism-md flex flex-col gap-4">
  <h3 class="font-display font-bold text-sm uppercase text-zinc-700 border-b-2 border-dashed border-zinc-200 pb-2 flex items-center gap-1.5">
    <User class="w-4 h-4 text-neubrutalism-charcoal" />
    <span>Customer & Device Details</span>
  </h3>

  <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
    <Input 
      id="cust-name" 
      label="Customer Name"
      type="text" 
      bind:value={customerName} 
      disabled={!isEditing || isSubmitting} 
    />

    <Input 
      id="cust-phone" 
      label="Phone Number"
      type="text" 
      bind:value={customerPhone} 
      disabled={!isEditing || isSubmitting} 
    />

    <Input 
      id="dev-brand" 
      label="Brand"
      type="text" 
      bind:value={brandPhone} 
      disabled={!isEditing || isSubmitting} 
    />

    <Input 
      id="dev-model" 
      label="Model"
      type="text" 
      bind:value={modelPhone} 
      disabled={!isEditing || isSubmitting} 
    />

    <div class="sm:col-span-2">
      <Input 
        id="dev-serial" 
        label="IMEI / Serial Number"
        type="text" 
        bind:value={serialNumber} 
        disabled={!isEditing || isSubmitting} 
      />
    </div>
  </div>
</Card>
