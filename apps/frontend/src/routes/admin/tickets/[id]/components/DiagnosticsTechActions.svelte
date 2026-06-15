<script lang="ts">
  import { Card } from '$lib';
  import { Wrench } from 'lucide-svelte';

  let {
    damageDescription = $bindable(''),
    repairAction = $bindable(''),
    isSubmitting,
    isEditing = false
  }: {
    damageDescription: string;
    repairAction: string;
    isSubmitting: boolean;
    isEditing?: boolean;
  } = $props();
</script>

<Card bgColor="bg-white" class="border-4 border-neubrutalism-charcoal p-6 shadow-neubrutalism-md flex flex-col gap-4">
  <h3 class="font-display font-bold text-sm uppercase text-zinc-700 border-b-2 border-dashed border-zinc-200 pb-2 flex items-center gap-1.5">
    <Wrench class="w-4 h-4 text-neubrutalism-charcoal" />
    <span>Diagnostics & Technical Actions</span>
  </h3>

  <div class="flex flex-col gap-4">
    <div class="flex flex-col gap-1">
      <label for="damage" class="font-display font-bold text-sm text-neubrutalism-charcoal uppercase tracking-wider">Customer Damage Complaint</label>
      <textarea 
        id="damage" 
        bind:value={damageDescription} 
        disabled={!isEditing || isSubmitting}
        rows="3"
        class="w-full border-4 border-neubrutalism-charcoal bg-white p-3 font-sans text-sm focus:outline-none shadow-neubrutalism-sm"
      ></textarea>
    </div>

    <div class="flex flex-col gap-1">
      <label for="action" class="font-display font-bold text-sm text-neubrutalism-charcoal uppercase tracking-wider">Repair Action (Technician Log)</label>
      <textarea 
        id="action" 
        placeholder="Log the repair operations performed (e.g. replaced LCD, soldered charging board)..."
        bind:value={repairAction} 
        disabled={!isEditing || isSubmitting}
        rows="3"
        class="w-full border-4 border-neubrutalism-charcoal bg-white p-3 font-sans text-sm focus:outline-none shadow-neubrutalism-sm"
      ></textarea>
    </div>
  </div>
</Card>
