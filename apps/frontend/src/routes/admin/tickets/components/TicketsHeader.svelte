<script lang="ts">
  import { Button } from '$lib';
  import { Plus } from 'lucide-svelte';
</script>

<div class="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
  <div>
    <h2 class="font-display font-extrabold text-2xl md:text-3xl uppercase tracking-tight">
      Repair Tickets Management
    </h2>
    <p class="font-sans text-xs sm:text-sm text-zinc-500 font-semibold mt-1">
      Manage hardware diagnostics, active repairs, component changes, and archives.
    </p>
  </div>
  
  <a href="/admin/tickets/new" class="w-full sm:w-auto">
    <Button bgColor="bg-neubrutalism-yellow" class="w-full sm:w-auto flex items-center justify-center gap-2 py-2 px-4 font-bold shadow-neubrutalism-sm">
      <Plus class="w-4 h-4" />
      <span>NEW TICKET</span>
    </Button>
  </a>
</div>
