<script lang="ts">
  import { Card, Button, Input } from '$lib';
  import { onMount } from 'svelte';
  import { Settings, Save, ShieldCheck, ShieldAlert, Store, Clock, RefreshCw } from 'lucide-svelte';

  let shopName = $state('OpenBench Repairs');
  let shopAddress = $state('Jalan Teknisi Raya No. 101, Bandung');
  let shopPhone = $state('0812-3456-7890');
  let warrantyDays = $state(30);
  let taxRate = $state(0);

  let successMessage = $state('');
  let errorMessage = $state('');
  let isSaving = $state(false);

  onMount(() => {
    // Load config from localStorage
    const saved = localStorage.getItem('openbench_mock_settings');
    if (saved) {
      try {
        const config = JSON.parse(saved);
        shopName = config.shopName || shopName;
        shopAddress = config.shopAddress || shopAddress;
        shopPhone = config.shopPhone || shopPhone;
        warrantyDays = config.warrantyDays || warrantyDays;
        taxRate = config.taxRate || taxRate;
      } catch (e) {
        // use defaults
      }
    }
  });

  async function handleSave(e: Event) {
    e.preventDefault();
    successMessage = '';
    errorMessage = '';
    isSaving = true;

    if (!shopName.trim() || !shopPhone.trim()) {
      errorMessage = 'Shop name and contact number are required.';
      isSaving = false;
      return;
    }

    try {
      const config = {
        shopName: shopName.trim(),
        shopAddress: shopAddress.trim(),
        shopPhone: shopPhone.trim(),
        warrantyDays,
        taxRate
      };

      localStorage.setItem('openbench_mock_settings', JSON.stringify(config));
      
      // Simulate delay
      await new Promise(resolve => setTimeout(resolve, 500));
      successMessage = 'Shop settings successfully saved.';
    } catch (err: any) {
      errorMessage = err.message || 'Failed to save settings.';
    } finally {
      isSaving = false;
    }
  }

  function handleResetMockDB() {
    if (confirm('Are you sure you want to reset mock data? All added tickets, inventory, and new transactions will be deleted.')) {
      localStorage.removeItem('openbench_mock_db');
      successMessage = 'Mock database successfully reset to default data. Please reload the page.';
      setTimeout(() => window.location.reload(), 1500);
    }
  }
</script>

<svelte:head>
  <title>Settings - OpenBench Admin</title>
</svelte:head>

<div class="flex flex-col gap-6 max-w-2xl mx-auto">

  <!-- Header Area -->
  <div>
    <h2 class="font-display font-extrabold text-2xl md:text-3xl uppercase tracking-tight flex items-center gap-2">
      <Settings class="w-8 h-8 text-zinc-700" />
      <span>System Configuration</span>
    </h2>
    <p class="font-sans text-xs sm:text-sm text-zinc-500 font-semibold mt-1">
      Adjust business parameters, checkout policies, and environment flags.
    </p>
  </div>

  <!-- Messages -->
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

  <form onsubmit={handleSave} class="flex flex-col gap-6">
    <!-- Store Identity -->
    <Card bgColor="bg-white" class="border-4 border-neubrutalism-charcoal p-6 shadow-neubrutalism-md flex flex-col gap-4">
      <h3 class="font-display font-bold text-sm uppercase text-zinc-700 border-b-2 border-dashed border-zinc-200 pb-2 flex items-center gap-1.5">
        <Store class="w-4 h-4 text-neubrutalism-charcoal" />
        <span>Store Identity</span>
      </h3>

      <div class="flex flex-col gap-4">
        <Input id="shop-name" label="Shop Name" type="text" bind:value={shopName} disabled={isSaving} />
        <Input id="shop-address" label="Shop Address" type="text" bind:value={shopAddress} disabled={isSaving} />
        <Input id="shop-phone" label="Contact / WA Number" type="text" bind:value={shopPhone} disabled={isSaving} />
      </div>
    </Card>

    <!-- Policies -->
    <Card bgColor="bg-white" class="border-4 border-neubrutalism-charcoal p-6 shadow-neubrutalism-md flex flex-col gap-4">
      <h3 class="font-display font-bold text-sm uppercase text-zinc-700 border-b-2 border-dashed border-zinc-200 pb-2 flex items-center gap-1.5">
        <Clock class="w-4 h-4 text-neubrutalism-charcoal" />
        <span>Checkout & Warranty Rules</span>
      </h3>

      <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
        <div class="flex flex-col gap-2 w-full">
          <label for="warranty-days" class="font-display font-bold text-sm text-neubrutalism-charcoal uppercase tracking-wider">Default Warranty duration (Days)</label>
          <input 
            id="warranty-days" 
            type="number" 
            min="0"
            placeholder="30"
            bind:value={warrantyDays} 
            disabled={isSaving}
            class="w-full border-4 border-neubrutalism-charcoal bg-white p-3 font-sans text-neubrutalism-charcoal rounded-none transition-all duration-150 focus:outline-none focus:ring-4 focus:ring-neubrutalism-charcoal focus:bg-[#fefefe] focus:placeholder-transparent"
          />
        </div>

        <div class="flex flex-col gap-2 w-full">
          <label for="tax-rate" class="font-display font-bold text-sm text-neubrutalism-charcoal uppercase tracking-wider">Service Tax Rate (%)</label>
          <input 
            id="tax-rate" 
            type="number" 
            min="0"
            max="100"
            placeholder="11"
            bind:value={taxRate} 
            disabled={isSaving}
            class="w-full border-4 border-neubrutalism-charcoal bg-white p-3 font-sans text-neubrutalism-charcoal rounded-none transition-all duration-150 focus:outline-none focus:ring-4 focus:ring-neubrutalism-charcoal focus:bg-[#fefefe] focus:placeholder-transparent"
          />
        </div>
      </div>
    </Card>

    <!-- Sandbox / Developer configurations -->
    <Card bgColor="bg-zinc-50" class="border-4 border-dashed border-zinc-400 p-6 flex flex-col gap-4">
      <h3 class="font-display font-bold text-sm uppercase text-zinc-700 border-b-2 border-dashed border-zinc-200 pb-2">
        Sandbox / Developer Utilities
      </h3>
      
      <div class="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
        <div class="flex flex-col gap-1">
          <span class="font-mono text-[10px] font-bold uppercase text-zinc-600">Reset Local Sandbox Data</span>
          <p class="font-sans text-xs text-zinc-500">Reset the local client storage and return all mock tickets and items back to defaults.</p>
        </div>
        
        <Button 
          bgColor="bg-neubrutalism-pink" 
          type="button" 
          onclick={handleResetMockDB}
          class="w-full sm:w-auto py-2 px-4 text-xs font-bold text-white shadow-neubrutalism-sm flex items-center justify-center gap-1.5"
        >
          <RefreshCw class="w-3.5 h-3.5" />
          <span>RESET MOCK DB</span>
        </Button>
      </div>
    </Card>

    <!-- Save configuration actions -->
    <div class="flex justify-end mt-2">
      <Button 
        bgColor="bg-neubrutalism-yellow" 
        type="submit" 
        disabled={isSaving}
        class="py-3 px-8 font-bold shadow-neubrutalism-md border-2 border-neubrutalism-charcoal flex items-center justify-center gap-2"
      >
        <Save class="w-4 h-4" />
        <span>{isSaving ? 'SAVING...' : 'SAVE SETTINGS'}</span>
      </Button>
    </div>
  </form>

</div>
