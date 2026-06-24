<script lang="ts">
  import { inventoryService, type Product } from '$lib/services/inventory';
  import { saleService, type Sale } from '$lib/services/sales';
  import { onMount } from 'svelte';
  import { toastService } from '$lib/services/toast.svelte';

  import POSHeader from './components/POSHeader.svelte';
  import POSCatalog from './components/POSCatalog.svelte';
  import POSCart from './components/POSCart.svelte';
  import ReceiptModal from './components/ReceiptModal.svelte';

  let products = $state<Product[]>([]);
  let isLoading = $state(true);
  let searchQuery = $state('');

  // Cart state
  interface CartItem {
    product: Product;
    qty: number;
  }
  let cart = $state<CartItem[]>([]);
  let discount = $state<number>(0);
  let paymentMethod = $state<'cash' | 'qris'>('cash');
  let cashPaid = $state<number>(0);

  // Success Transaction Modal state
  let completedSale = $state<Sale | null>(null);

  // Compute cart calculations reactively
  let subtotal = $derived(cart.reduce((sum, item) => sum + (item.product.price * item.qty), 0));
  let finalTotal = $derived(Math.max(0, subtotal - discount));
  let changeAmount = $derived(paymentMethod === 'cash' ? Math.max(0, cashPaid - finalTotal) : 0);

  onMount(async () => {
    await loadProducts();
  });

  async function loadProducts() {
    try {
      // Load products and show retail accessories primarily, but also allow spare parts
      const data = await inventoryService.getInventory();
      products = data;
    } catch (err: any) {
      console.error('Error loading inventory:', err);
      toastService.error('Failed to load inventory: ' + err.message);
      products = [];
    } finally {
      isLoading = false;
    }
  }

  // Filter products for catalog search
  let filteredProducts = $derived(
    products.filter(p => {
      const matchSearch = p.name.toLowerCase().includes(searchQuery.toLowerCase());
      return matchSearch;
    })
  );

  function addToCart(p: Product) {
    if (p.stock <= 0) {
      toastService.warning(`Stock of ${p.name} is empty!`);
      return;
    }
    const existingIdx = cart.findIndex(item => item.product.id === p.id);
    if (existingIdx !== -1) {
      if (cart[existingIdx].qty >= p.stock) {
        toastService.warning(`Cannot exceed available stock (${p.stock} pcs)`);
        return;
      }
      cart[existingIdx].qty += 1;
    } else {
      cart.push({ product: p, qty: 1 });
    }
  }

  function removeFromCart(pId: string) {
    cart = cart.filter(item => item.product.id !== pId);
  }

  function adjustQty(pId: string, amt: number) {
    const idx = cart.findIndex(item => item.product.id === pId);
    if (idx === -1) return;

    const newQty = cart[idx].qty + amt;
    if (newQty <= 0) {
      removeFromCart(pId);
    } else {
      // Check stock limit
      if (newQty > cart[idx].product.stock) {
        toastService.warning(`Stock limited! Only ${cart[idx].product.stock} units available.`);
        return;
      }
      cart[idx].qty = newQty;
    }
  }

  function clearCart() {
    cart = [];
    discount = 0;
    cashPaid = 0;
    paymentMethod = 'cash';
  }

  async function handleCheckout() {
    if (cart.length === 0) {
      toastService.warning('Shopping cart is empty!');
      return;
    }

    if (paymentMethod === 'cash' && cashPaid < finalTotal) {
      toastService.warning('Cash paid is less than the final total!');
      return;
    }

    try {
      const saleItems = cart.map(item => ({
        product_id: item.product.id,
        qty: item.qty
      }));

      const newSale = await saleService.createSale({
        items: saleItems,
        discount: discount,
        payment_method: paymentMethod
      });

      completedSale = newSale;
      toastService.success('Transaction completed successfully!');
      clearCart();
      await loadProducts(); // reload products to reflect depleted stocks
    } catch (err: any) {
      toastService.error('Failed to complete transaction: ' + err.message);
    }
  }

  function closeSuccessModal() {
    completedSale = null;
  }
</script>

<svelte:head>
  <title>Point of Sales - OpenBench Admin</title>
</svelte:head>

<div class="flex flex-col gap-6">

  <!-- Header Area -->
  <POSHeader />

  <!-- Checkout Grid Layout: Left (Catalog) & Right (Cart) -->
  <div class="grid grid-cols-1 lg:grid-cols-5 gap-8">
    
    <!-- Left: Catalog List (3 cols) -->
    <div class="lg:col-span-3">
      <POSCatalog
        bind:searchQuery
        {isLoading}
        {filteredProducts}
        onAddToCart={addToCart}
      />
    </div>

    <!-- Right: Shopping Cart Checkout console (2 cols) -->
    <div class="lg:col-span-2">
      <POSCart
        {cart}
        bind:discount
        bind:paymentMethod
        bind:cashPaid
        {subtotal}
        {finalTotal}
        {changeAmount}
        onAdjustQty={adjustQty}
        onCheckout={handleCheckout}
      />
    </div>

  </div>

  <!-- Success Transaction Receipt Modal Overlay -->
  {#if completedSale !== null}
    <ReceiptModal
      {completedSale}
      onClose={closeSuccessModal}
    />
  {/if}

</div>
