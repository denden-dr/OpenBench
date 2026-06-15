<script lang="ts">
  import { inventoryService } from '$lib/services/inventory';
  import { goto } from '$app/navigation';
  import { ArrowLeft, ShieldAlert } from 'lucide-svelte';
  import ProductForm from '../components/ProductForm.svelte';

  // Form states using Svelte 5 state
  let prodName = $state('');
  let prodCategory = $state<'retail' | 'spare_part'>('retail');
  let prodCostPrice = $state<number>(0);
  let prodPrice = $state<number>(0);
  let prodStock = $state<number>(0);
  let prodMinStock = $state<number>(5);

  let errorMessage = $state('');
  let isSubmitting = $state(false);

  async function handleSubmit(event: Event) {
    event.preventDefault();
    errorMessage = '';
    isSubmitting = true;

    // Basic Client-Side Validations
    if (!prodName.trim()) {
      errorMessage = 'Product name is required.';
      isSubmitting = false;
      return;
    }
    if (prodStock < 0 || prodPrice < 0 || prodCostPrice < 0 || prodMinStock < 0) {
      errorMessage = 'Values cannot be negative.';
      isSubmitting = false;
      return;
    }

    try {
      await inventoryService.createProduct({
        name: prodName.trim(),
        category: prodCategory,
        stock: prodStock,
        price: prodPrice,
        cost_price: prodCostPrice,
        min_stock: prodMinStock
      });

      // Redirect back to inventory index
      await goto('/admin/inventory');
    } catch (err: any) {
      errorMessage = err.message || 'Failed to save new product.';
    } finally {
      isSubmitting = false;
    }
  }

  function handleCancel() {
    goto('/admin/inventory');
  }
</script>

<svelte:head>
  <title>New Product - OpenBench Admin</title>
</svelte:head>

<div class="flex flex-col gap-6 max-w-2xl mx-auto">
  
  <!-- Back Action -->
  <div>
    <a href="/admin/inventory" class="inline-flex items-center gap-1.5 font-mono text-xs font-bold uppercase hover:underline">
      <ArrowLeft class="w-4 h-4" />
      <span>BACK TO INVENTORY</span>
    </a>
  </div>

  <!-- Heading -->
  <div>
    <h2 class="font-display font-extrabold text-2xl md:text-3xl uppercase tracking-tight">
      Add New Product / Component
    </h2>
    <p class="font-sans text-xs sm:text-sm text-zinc-500 font-semibold mt-1">
      Add a sales accessory or repair component to the shop catalog.
    </p>
  </div>

  {#if errorMessage}
    <div class="bg-neubrutalism-pink text-white border-4 border-neubrutalism-charcoal p-4 font-mono text-xs flex items-center gap-2 shadow-neubrutalism-sm">
      <ShieldAlert class="w-5 h-5 shrink-0" />
      <span>{errorMessage}</span>
    </div>
  {/if}

  <!-- Form Component -->
  <ProductForm
    editingId={null}
    bind:prodName
    bind:prodCategory
    bind:prodCostPrice
    bind:prodPrice
    bind:prodStock
    bind:prodMinStock
    onSubmit={handleSubmit}
    onCancel={handleCancel}
  />

</div>
