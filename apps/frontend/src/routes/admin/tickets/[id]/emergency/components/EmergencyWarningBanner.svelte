<script lang="ts">
  import { ShieldAlert } from 'lucide-svelte';
</script>

<div class="bg-amber-100 border-4 border-neubrutalism-charcoal p-6 flex flex-col gap-2 shadow-neubrutalism-sm">
  <div class="flex items-center gap-2 text-amber-800 font-display font-extrabold text-sm uppercase">
    <ShieldAlert class="w-5 h-5" />
    <span>WARNING: Emergency Bypass Operation</span>
  </div>
  <p class="font-sans text-xs text-amber-900 leading-relaxed font-semibold">
    You are making modifications that bypass standard business rules. If you revert the status from "Picked Up" to "In Warehouse/Store", any active warranty records associated with this ticket will be automatically voided and deleted from the system.
  </p>
</div>
