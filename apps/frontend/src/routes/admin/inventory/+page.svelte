<script lang="ts">
  import { inventoryService, type Product } from '$lib/services/inventory';
  import { onMount } from 'svelte';
  import { ShieldAlert, Check } from 'lucide-svelte';

  import InventoryHeader from './components/InventoryHeader.svelte';
  import InventoryStats from './components/InventoryStats.svelte';
  import InventoryFilters from './components/InventoryFilters.svelte';
  import ProductForm from './components/ProductForm.svelte';
  import InventoryTable from './components/InventoryTable.svelte';

  let products = $state<Product[]>([]);
  let isLoading = $state(true);

  // Filter & Search states
  let searchQuery = $state('');
  let categoryFilter = $state<string>('all');
  let stockFilter = $state<string>('all'); // all, low_stock

  // Form states (Edit)
  let editingId = $state<string | null>(null);

  // Edit Product states
  let prodName = $state('');
  let prodCategory = $state<'retail' | 'spare_part'>('retail');
  let prodStock = $state<number>(0);
  let prodPrice = $state<number>(0);
  let prodCostPrice = $state<number>(0);
  let prodMinStock = $state<number>(5);

  let successMessage = $state('');
  let errorMessage = $state('');

  onMount(async () => {
    await loadInventory();
  });

  async function loadInventory() {
    try {
      products = await inventoryService.getInventory();
    } catch (err: any) {
      console.error('Error loading inventory:', err);
      errorMessage = 'Failed to load inventory: ' + err.message;
      products = [];
    } finally {
      isLoading = false;
    }
  }

  // Filter products reactively
  let filteredProducts = $derived(
    products.filter(p => {
      // 1. Search Query
      if (searchQuery.trim() !== '') {
        if (!p.name.toLowerCase().includes(searchQuery.toLowerCase())) return false;
      }
      // 2. Category Filter
      if (categoryFilter !== 'all' && p.category !== categoryFilter) return false;
      // 3. Stock Filter
      if (stockFilter === 'low_stock' && p.stock > p.min_stock) return false;

      return true;
    })
  );

  // Low stock counter
  let lowStockCount = $derived(products.filter(p => p.stock <= p.min_stock).length);

  function resetForm() {
    prodName = '';
    prodCategory = 'retail';
    prodStock = 0;
    prodPrice = 0;
    prodCostPrice = 0;
    prodMinStock = 5;
    editingId = null;
    errorMessage = '';
    successMessage = '';
  }

  async function handleEditProduct(e: Event) {
    e.preventDefault();
    errorMessage = '';
    successMessage = '';

    if (!editingId) return;

    if (!prodName.trim()) {
      errorMessage = 'Product name is required.';
      return;
    }
    if (prodStock < 0 || prodPrice < 0 || prodCostPrice < 0 || prodMinStock < 0) {
      errorMessage = 'Values cannot be negative.';
      return;
    }

    try {
      // Edit existing product
      await inventoryService.updateProduct(editingId, {
        name: prodName.trim(),
        category: prodCategory,
        stock: prodStock,
        price: prodPrice,
        cost_price: prodCostPrice,
        min_stock: prodMinStock
      });
      successMessage = 'Product successfully updated.';
      resetForm();
      await loadInventory();
    } catch (err: any) {
      errorMessage = err.message || 'Failed to save product.';
    }
  }

  function startEdit(p: Product) {
    editingId = p.id;
    prodName = p.name;
    prodCategory = p.category;
    prodStock = p.stock;
    prodPrice = p.price;
    prodCostPrice = p.cost_price;
    prodMinStock = p.min_stock;
    errorMessage = '';
    successMessage = '';
  }

  async function handleDelete(id: string) {
    if (!confirm('Are you sure you want to delete this product from the catalog?')) return;
    try {
      await inventoryService.deleteProduct(id);
      successMessage = 'Product successfully deleted.';
      await loadInventory();
    } catch (err: any) {
      errorMessage = err.message || 'Failed to delete product.';
    }
  }

  async function adjustStock(id: string, currentStock: number, amt: number) {
    try {
      const targetStock = Math.max(0, currentStock + amt);
      await inventoryService.updateProduct(id, { stock: targetStock });
      await loadInventory();
    } catch (err: any) {
      alert('Failed to adjust stock: ' + err.message);
    }
  }
</script>

<svelte:head>
  <title>Inventory / Stock - OpenBench Admin</title>
</svelte:head>

<div class="flex flex-col gap-6">

  <!-- Header Area -->
  <InventoryHeader />

  <!-- Messages -->
  {#if successMessage}
    <div class="bg-neubrutalism-green border-4 border-neubrutalism-charcoal p-4 font-mono text-xs flex items-center gap-2 shadow-neubrutalism-sm">
      <Check class="w-5 h-5 shrink-0" />
      <span>{successMessage}</span>
    </div>
  {/if}

  {#if errorMessage}
    <div class="bg-neubrutalism-pink text-white border-4 border-neubrutalism-charcoal p-4 font-mono text-xs flex items-center gap-2 shadow-neubrutalism-sm">
      <ShieldAlert class="w-5 h-5 shrink-0" />
      <span>{errorMessage}</span>
    </div>
  {/if}

  <!-- Edit Form Box -->
  {#if editingId !== null}
    <ProductForm
      {editingId}
      bind:prodName
      bind:prodCategory
      bind:prodCostPrice
      bind:prodPrice
      bind:prodStock
      bind:prodMinStock
      onSubmit={handleEditProduct}
      onCancel={resetForm}
    />
  {/if}

  <!-- Stats Indicators Quick Overview -->
  <InventoryStats {products} />

  <!-- Search & Filter Area -->
  <InventoryFilters 
    bind:searchQuery
    bind:categoryFilter
    bind:stockFilter
  />

  <!-- Inventory Catalog Table / Grid -->
  <InventoryTable
    {isLoading}
    {filteredProducts}
    onAdjustStock={adjustStock}
    onStartEdit={startEdit}
    onDelete={handleDelete}
  />

</div>
