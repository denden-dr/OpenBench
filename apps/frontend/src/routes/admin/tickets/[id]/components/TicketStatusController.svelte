<script lang="ts">
  import { Card } from '$lib';
  import { Tag } from 'lucide-svelte';

  let {
    status = $bindable('received'),
    devicePosition = $bindable('warehouse'),
    isSubmitting,
    isEditing = false
  }: {
    status: 'received' | 'diagnosing' | 'in_repair' | 'ready_for_pickup' | 'picked_up' | 'cancelled';
    devicePosition: 'warehouse' | 'picked_up';
    isSubmitting: boolean;
    isEditing?: boolean;
  } = $props();
</script>

<Card bgColor="bg-white" class="border-4 border-neubrutalism-charcoal p-6 shadow-neubrutalism-md flex flex-col gap-4">
  <h3 class="font-display font-bold text-sm uppercase text-zinc-700 border-b-2 border-dashed border-zinc-200 pb-2 flex items-center gap-1.5">
    <Tag class="w-4 h-4 text-neubrutalism-charcoal" />
    <span>Repair Status</span>
  </h3>

  <div class="flex flex-col gap-1.5">
    <label for="status-select" class="font-mono text-[10px] font-bold uppercase text-zinc-500">Service Status</label>
    <select 
      id="status-select" 
      bind:value={status}
      disabled={!isEditing || isSubmitting}
      class="w-full bg-white border-4 border-neubrutalism-charcoal px-3 py-2 font-mono text-xs shadow-neubrutalism-sm focus:outline-none"
    >
      <option value="received">RECEIVED</option>
      <option value="diagnosing">DIAGNOSING</option>
      <option value="in_repair">IN REPAIR</option>
      <option value="ready_for_pickup">READY FOR PICKUP</option>
      <option value="picked_up">PICKED UP (COMPLETED)</option>
      <option value="cancelled">CANCELLED</option>
    </select>
  </div>

  {#if status === 'picked_up'}
    <div class="bg-amber-100 border-2 border-amber-400 p-3 font-mono text-[10px] text-amber-800 leading-snug">
      <strong>NOTE:</strong> Setting status to PICKED UP will automatically trigger active warranty registration (30 days) and update payment status to PAID.
    </div>
  {/if}

  <div class="flex flex-col gap-1.5 mt-2">
    <label for="pos-select" class="font-mono text-[10px] font-bold uppercase text-zinc-500">Device Location Status</label>
    <select 
      id="pos-select" 
      bind:value={devicePosition}
      disabled={!isEditing || isSubmitting || status === 'picked_up'}
      class="w-full bg-white border-4 border-neubrutalism-charcoal px-3 py-2 font-mono text-xs shadow-neubrutalism-sm focus:outline-none"
    >
      <option value="warehouse">IN WAREHOUSE / STORE</option>
      <option value="picked_up">PICKED UP / TAKEN HOME</option>
    </select>
  </div>
</Card>
