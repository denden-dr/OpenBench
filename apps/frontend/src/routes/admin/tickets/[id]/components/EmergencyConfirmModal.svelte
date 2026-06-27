<script lang="ts">
  import { Button } from '$lib';
  import { ShieldAlert, X } from 'lucide-svelte';
  import { goto } from '$app/navigation';

  interface Props {
    isOpen: boolean;
    ticketId: string;
  }

  let { isOpen = $bindable(), ticketId }: Props = $props();
</script>

{#if isOpen}
  <div class="fixed inset-0 bg-black/60 z-50 flex items-center justify-center p-4">
    <div class="max-w-md w-full bg-white border-4 border-neubrutalism-charcoal shadow-neubrutalism-lg flex flex-col overflow-hidden max-h-[90vh]">
      <!-- Header -->
      <div class="p-4 bg-neubrutalism-pink border-b-4 border-neubrutalism-charcoal flex justify-between items-center text-white font-display font-extrabold text-sm uppercase">
        <div class="flex items-center gap-2">
          <ShieldAlert class="w-5 h-5 text-white" />
          <span class="text-white">Emergency Edit Mode</span>
        </div>
        <button onclick={() => isOpen = false} class="p-1 hover:bg-white/20" aria-label="Close modal">
          <X class="w-5 h-5 text-white" />
        </button>
      </div>

      <!-- Content -->
      <div class="p-6 flex flex-col gap-4">
        <p class="font-sans text-sm text-neubrutalism-charcoal leading-relaxed font-semibold">
          Are you sure you want to enter Emergency Edit mode?
        </p>
        <p class="font-sans text-xs text-zinc-650 leading-relaxed">
          Manual changes to pricing, device data, or statuses bypass normal validation checks. This can affect customer invoices, payment statuses, and warranty validation dates.
        </p>

        <!-- Actions -->
        <div class="flex justify-end gap-3 mt-2">
          <Button 
            bgColor="bg-zinc-200" 
            onclick={() => isOpen = false}
            class="py-2 px-4 text-xs font-bold shadow-neubrutalism-sm"
          >
            CANCEL
          </Button>
          <Button 
            bgColor="bg-neubrutalism-pink" 
            onclick={() => {
              isOpen = false;
              goto(`/admin/tickets/${ticketId}/emergency`);
            }}
            class="py-2 px-4 text-xs font-bold text-white shadow-neubrutalism-sm"
          >
            <span class="text-white">CONFIRM EDIT</span>
          </Button>
        </div>
      </div>
    </div>
  </div>
{/if}
