<script lang="ts">
  import { Search, Wrench, ShieldCheck, Smartphone, ArrowRight, AlertCircle, X } from 'lucide-svelte';
  import { fade, fly } from 'svelte/transition';

  // Mock Data
  type RepairService = { 
    name: string; 
    price: number; 
    grade: 'Original' | 'ODM';
    type?: 'OLED' | 'IPS' | 'Standard';
  };
  type DevicePricing = { id: string; brand: string; model: string; services: RepairService[] };

  const devices: DevicePricing[] = [
    // Apple
    { id: 'a1', brand: 'Apple', model: 'iPhone 14 Pro Max', services: [{ name: 'Screen Replacement', grade: 'Original', type: 'OLED', price: 6500000 }, { name: 'Screen Replacement', grade: 'ODM', type: 'OLED', price: 3800000 }, { name: 'Battery', grade: 'Original', price: 1200000 }, { name: 'Back Door', grade: 'Original', price: 2200000 }, { name: 'Rear Camera', grade: 'Original', price: 2800000 }, { name: 'Charging Port', grade: 'Original', price: 950000 }] },
    { id: 'a2', brand: 'Apple', model: 'iPhone 14 Pro', services: [{ name: 'Screen Replacement', grade: 'Original', type: 'OLED', price: 5800000 }, { name: 'Screen Replacement', grade: 'ODM', type: 'OLED', price: 3200000 }, { name: 'Battery', grade: 'Original', price: 1100000 }, { name: 'Back Door', grade: 'Original', price: 2100000 }, { name: 'Rear Camera', grade: 'Original', price: 2500000 }, { name: 'Charging Port', grade: 'Original', price: 950000 }] },
    { id: 'a3', brand: 'Apple', model: 'iPhone 14', services: [{ name: 'Screen Replacement', grade: 'Original', type: 'OLED', price: 4200000 }, { name: 'Screen Replacement', grade: 'ODM', type: 'OLED', price: 2200000 }, { name: 'Battery', grade: 'Original', price: 950000 }, { name: 'Battery', grade: 'ODM', price: 450000 }, { name: 'Back Door', grade: 'Original', price: 1500000 }] },
    { id: 'a4', brand: 'Apple', model: 'iPhone 13 Pro Max', services: [{ name: 'Screen Replacement', grade: 'Original', type: 'OLED', price: 5200000 }, { name: 'Screen Replacement', grade: 'ODM', type: 'OLED', price: 2800000 }, { name: 'Battery', grade: 'Original', price: 1000000 }, { name: 'Back Door', grade: 'Original', price: 1800000 }] },
    { id: 'a5', brand: 'Apple', model: 'iPhone 13 Pro', services: [{ name: 'Screen Replacement', grade: 'Original', type: 'OLED', price: 4500000 }, { name: 'Screen Replacement', grade: 'ODM', type: 'OLED', price: 2400000 }, { name: 'Battery', grade: 'Original', price: 950000 }, { name: 'Back Door', grade: 'Original', price: 1600000 }] },
    { id: 'a6', brand: 'Apple', model: 'iPhone 13', services: [{ name: 'Screen Replacement', grade: 'Original', type: 'OLED', price: 3200000 }, { name: 'Screen Replacement', grade: 'ODM', type: 'OLED', price: 1800000 }, { name: 'Battery', grade: 'Original', price: 850000 }, { name: 'Battery', grade: 'ODM', price: 350000 }] },
    { id: 'a7', brand: 'Apple', model: 'iPhone 12 Pro Max', services: [{ name: 'Screen Replacement', grade: 'Original', type: 'OLED', price: 4200000 }, { name: 'Screen Replacement', grade: 'ODM', type: 'OLED', price: 1950000 }, { name: 'Battery', grade: 'Original', price: 850000 }, { name: 'Battery', grade: 'ODM', price: 350000 }] },
    { id: 'a8', brand: 'Apple', model: 'iPhone 12', services: [{ name: 'Screen Replacement', grade: 'Original', type: 'OLED', price: 2800000 }, { name: 'Screen Replacement', grade: 'ODM', type: 'OLED', price: 1450000 }, { name: 'Battery', grade: 'Original', price: 750000 }, { name: 'Battery', grade: 'ODM', price: 350000 }] },
    { id: 'a9', brand: 'Apple', model: 'iPhone 11 Pro Max', services: [{ name: 'Screen Replacement', grade: 'Original', type: 'OLED', price: 3200000 }, { name: 'Screen Replacement', grade: 'ODM', type: 'OLED', price: 1550000 }, { name: 'Battery', grade: 'Original', price: 650000 }, { name: 'Battery', grade: 'ODM', price: 350000 }] },
    { id: 'a10', brand: 'Apple', model: 'iPhone 11', services: [{ name: 'Screen Replacement', grade: 'Original', type: 'IPS', price: 1800000 }, { name: 'Screen Replacement', grade: 'ODM', type: 'IPS', price: 850000 }, { name: 'Battery', grade: 'Original', price: 550000 }, { name: 'Battery', grade: 'ODM', price: 250000 }] },

    // Samsung
    { id: 's1', brand: 'Samsung', model: 'Galaxy S22 Ultra', services: [{ name: 'Screen Replacement', grade: 'Original', type: 'OLED', price: 5200000 }, { name: 'Screen Replacement', grade: 'ODM', price: 2800000 }, { name: 'Battery', grade: 'Original', price: 850000 }, { name: 'Back Door', grade: 'Original', price: 950000 }, { name: 'Charging Port', grade: 'Original', price: 750000 }] },
    { id: 's2', brand: 'Samsung', model: 'Galaxy S22+', services: [{ name: 'Screen Replacement', grade: 'Original', type: 'OLED', price: 3800000 }, { name: 'Screen Replacement', grade: 'ODM', type: 'OLED', price: 1950000 }, { name: 'Battery', grade: 'Original', price: 750000 }, { name: 'Battery', grade: 'ODM', price: 350000 }, { name: 'Back Door', grade: 'Original', price: 850000 }] },
    { id: 's3', brand: 'Samsung', model: 'Galaxy S22', services: [{ name: 'Screen Replacement', grade: 'Original', type: 'OLED', price: 3200000 }, { name: 'Screen Replacement', grade: 'ODM', type: 'OLED', price: 1650000 }, { name: 'Battery', grade: 'Original', price: 650000 }, { name: 'Battery', grade: 'ODM', price: 250000 }] },
    { id: 's4', brand: 'Samsung', model: 'Galaxy S21 Ultra', services: [{ name: 'Screen Replacement', grade: 'Original', type: 'OLED', price: 4500000 }, { name: 'Screen Replacement', grade: 'ODM', type: 'OLED', price: 2200000 }, { name: 'Battery', grade: 'Original', price: 750000 }, { name: 'Battery', grade: 'ODM', price: 350000 }, { name: 'Back Door', grade: 'Original', price: 850000 }] },
    { id: 's5', brand: 'Samsung', model: 'Galaxy S21', services: [{ name: 'Screen Replacement', grade: 'Original', type: 'OLED', price: 2800000 }, { name: 'Screen Replacement', grade: 'ODM', type: 'OLED', price: 1250000 }, { name: 'Battery', grade: 'Original', price: 650000 }, { name: 'Battery', grade: 'ODM', price: 250000 }] },
    { id: 's6', brand: 'Samsung', model: 'Galaxy A53 5G', services: [{ name: 'Screen Replacement', grade: 'Original', price: 1200000 }, { name: 'Battery', grade: 'Original', price: 450000 }, { name: 'Back Door', grade: 'Original', price: 350000 }] },
    { id: 's7', brand: 'Samsung', model: 'Galaxy A33 5G', services: [{ name: 'Screen Replacement', grade: 'Original', price: 950000 }, { name: 'Battery', grade: 'Original', price: 350000 }] },
    { id: 's8', brand: 'Samsung', model: 'Galaxy A73 5G', services: [{ name: 'Screen Replacement', grade: 'Original', price: 1800000 }, { name: 'Battery', grade: 'Original', price: 550000 }] },
    { id: 's9', brand: 'Samsung', model: 'Galaxy Note 20 Ultra', services: [{ name: 'Screen Replacement', grade: 'Original', type: 'OLED', price: 4800000 }, { name: 'Battery', grade: 'Original', price: 850000 }] },
    { id: 's10', brand: 'Samsung', model: 'Galaxy Z Fold 4', services: [{ name: 'Inner Screen', grade: 'Original', type: 'OLED', price: 8500000 }, { name: 'Battery', grade: 'Original', price: 1500000 }] },

    // Google
    { id: 'g1', brand: 'Google', model: 'Pixel 7 Pro', services: [{ name: 'Screen Replacement', grade: 'Original', type: 'OLED', price: 3800000 }, { name: 'Battery', grade: 'Original', price: 750000 }, { name: 'Charging Port', grade: 'Original', price: 850000 }] },
    { id: 'g2', brand: 'Google', model: 'Pixel 7', services: [{ name: 'Screen Replacement', grade: 'Original', type: 'OLED', price: 2800000 }, { name: 'Battery', grade: 'Original', price: 650000 }] },
    { id: 'g3', brand: 'Google', model: 'Pixel 6 Pro', services: [{ name: 'Screen Replacement', grade: 'Original', type: 'OLED', price: 3200000 }, { name: 'Battery', grade: 'Original', price: 750000 }] },
    { id: 'g4', brand: 'Google', model: 'Pixel 6', services: [{ name: 'Screen Replacement', grade: 'Original', type: 'OLED', price: 2200000 }, { name: 'Battery', grade: 'Original', price: 650000 }] },
    { id: 'g5', brand: 'Google', model: 'Pixel 6a', services: [{ name: 'Screen Replacement', grade: 'Original', price: 1500000 }, { name: 'Battery', grade: 'Original', price: 550000 }] },
    { id: 'g6', brand: 'Google', model: 'Pixel 5', services: [{ name: 'Screen Replacement', grade: 'Original', price: 1800000 }, { name: 'Battery', grade: 'Original', price: 550000 }] },
    { id: 'g7', brand: 'Google', model: 'Pixel 5a', services: [{ name: 'Screen Replacement', grade: 'Original', price: 1400000 }, { name: 'Battery', grade: 'Original', price: 450000 }] },
    { id: 'g8', brand: 'Google', model: 'Pixel 4 XL', services: [{ name: 'Screen Replacement', grade: 'Original', price: 1600000 }, { name: 'Battery', grade: 'Original', price: 550000 }] },
    { id: 'g9', brand: 'Google', model: 'Pixel 4', services: [{ name: 'Screen Replacement', grade: 'Original', price: 1200000 }, { name: 'Battery', grade: 'Original', price: 450000 }] },
    { id: 'g10', brand: 'Google', model: 'Pixel 4a', services: [{ name: 'Screen Replacement', grade: 'Original', price: 1100000 }, { name: 'Battery', grade: 'Original', price: 450000 }] },

    // Xiaomi
    { id: 'x1', brand: 'Xiaomi', model: '12 Pro', services: [{ name: 'Screen Replacement', grade: 'Original', type: 'OLED', price: 3500000 }, { name: 'Screen Replacement', grade: 'ODM', type: 'OLED', price: 1850000 }, { name: 'Battery', grade: 'Original', price: 750000 }, { name: 'Charging Port', grade: 'Original', price: 450000 }] },
    { id: 'x2', brand: 'Xiaomi', model: '12', services: [{ name: 'Screen Replacement', grade: 'Original', type: 'OLED', price: 2800000 }, { name: 'Screen Replacement', grade: 'ODM', type: 'OLED', price: 1450000 }, { name: 'Battery', grade: 'Original', price: 650000 }] },
    { id: 'x3', brand: 'Xiaomi', model: '12 Lite', services: [{ name: 'Screen Replacement', grade: 'Original', price: 1500000 }, { name: 'Battery', grade: 'Original', price: 550000 }] },
    { id: 'x4', brand: 'Xiaomi', model: 'Redmi Note 11 Pro', services: [{ name: 'Screen Replacement', grade: 'Original', price: 950000 }, { name: 'Battery', grade: 'Original', price: 450000 }] },
    { id: 'x5', brand: 'Xiaomi', model: 'Redmi Note 11', services: [{ name: 'Screen Replacement', grade: 'Original', price: 650000 }, { name: 'Battery', grade: 'Original', price: 350000 }] },
    { id: 'x6', brand: 'Xiaomi', model: 'Poco F4 GT', services: [{ name: 'Screen Replacement', grade: 'Original', price: 2200000 }, { name: 'Battery', grade: 'Original', price: 650000 }] },
    { id: 'x7', brand: 'Xiaomi', model: 'Poco X4 Pro', services: [{ name: 'Screen Replacement', grade: 'Original', price: 1100000 }, { name: 'Battery', grade: 'Original', price: 550000 }] },
    { id: 'x8', brand: 'Xiaomi', model: 'Mi 11 Ultra', services: [{ name: 'Screen Replacement', grade: 'Original', type: 'OLED', price: 4200000 }, { name: 'Battery', grade: 'Original', price: 850000 }] },
    { id: 'x9', brand: 'Xiaomi', model: 'Mi 11', services: [{ name: 'Screen Replacement', grade: 'Original', type: 'OLED', price: 3200000 }, { name: 'Battery', grade: 'Original', price: 750000 }] },
    { id: 'x10', brand: 'Xiaomi', model: 'Redmi 10', services: [{ name: 'Screen Replacement', grade: 'Original', price: 550000 }, { name: 'Battery', grade: 'Original', price: 250000 }] },

    // Oppo
    { id: 'o1', brand: 'Oppo', model: 'Find X5 Pro', services: [{ name: 'Screen Replacement', grade: 'Original', type: 'OLED', price: 5500000 }, { name: 'Screen Replacement', grade: 'ODM', type: 'OLED', price: 2850000 }, { name: 'Battery', grade: 'Original', price: 950000 }, { name: 'Back Door', grade: 'Original', price: 1200000 }] },
    { id: 'o2', brand: 'Oppo', model: 'Find X5', services: [{ name: 'Screen Replacement', grade: 'Original', type: 'OLED', price: 3800000 }, { name: 'Screen Replacement', grade: 'ODM', type: 'OLED', price: 1950000 }, { name: 'Battery', grade: 'Original', price: 750000 }] },
    { id: 'o3', brand: 'Oppo', model: 'Reno 8 Pro', services: [{ name: 'Screen Replacement', grade: 'Original', type: 'OLED', price: 3200000 }, { name: 'Battery', grade: 'Original', price: 750000 }] },
    { id: 'o4', brand: 'Oppo', model: 'Reno 8', services: [{ name: 'Screen Replacement', grade: 'Original', price: 1800000 }, { name: 'Battery', grade: 'Original', price: 650000 }] },
    { id: 'o5', brand: 'Oppo', model: 'Reno 7', services: [{ name: 'Screen Replacement', grade: 'Original', price: 1100000 }, { name: 'Battery', grade: 'Original', price: 400000 }] },
    { id: 'o6', brand: 'Oppo', model: 'A96', services: [{ name: 'Screen Replacement', grade: 'Original', price: 950000 }, { name: 'Battery', grade: 'Original', price: 450000 }] },
    { id: 'o7', brand: 'Oppo', model: 'A76', services: [{ name: 'Screen Replacement', grade: 'Original', price: 750000 }, { name: 'Battery', grade: 'Original', price: 350000 }] },
    { id: 'o8', brand: 'Oppo', model: 'A16', services: [{ name: 'Screen Replacement', grade: 'Original', price: 550000 }, { name: 'Battery', grade: 'Original', price: 250000 }] },
    { id: 'o9', brand: 'Oppo', model: 'Reno 6 Pro', services: [{ name: 'Screen Replacement', grade: 'Original', price: 2800000 }, { name: 'Battery', grade: 'Original', price: 650000 }] },
    { id: 'o10', brand: 'Oppo', model: 'Find N', services: [{ name: 'Inner Screen', grade: 'Original', type: 'OLED', price: 7500000 }, { name: 'Battery', grade: 'Original', price: 1200000 }] },

    // Asus
    { id: 'as1', brand: 'Asus', model: 'ROG Phone 6 Pro', services: [{ name: 'Screen Replacement', grade: 'Original', type: 'OLED', price: 4800000 }, { name: 'Battery', grade: 'Original', price: 1200000 }, { name: 'Side Port', grade: 'Original', price: 850000 }] },
    { id: 'as2', brand: 'Asus', model: 'ROG Phone 6', services: [{ name: 'Screen Replacement', grade: 'Original', type: 'OLED', price: 4200000 }, { name: 'Battery', grade: 'Original', price: 950000 }] },
    { id: 'as3', brand: 'Asus', model: 'ROG Phone 5s', services: [{ name: 'Screen Replacement', grade: 'Original', price: 3800000 }, { name: 'Battery', grade: 'Original', price: 850000 }] },
    { id: 'as4', brand: 'Asus', model: 'ROG Phone 5', services: [{ name: 'Screen Replacement', grade: 'Original', price: 3500000 }, { name: 'Battery', grade: 'Original', price: 850000 }] },
    { id: 'as5', brand: 'Asus', model: 'Zenfone 9', services: [{ name: 'Screen Replacement', grade: 'Original', price: 2800000 }, { name: 'Battery', grade: 'Original', price: 650000 }] },
    { id: 'as6', brand: 'Asus', model: 'Zenfone 8', services: [{ name: 'Screen Replacement', grade: 'Original', price: 2200000 }, { name: 'Battery', grade: 'Original', price: 650000 }] },
    { id: 'as7', brand: 'Asus', model: 'Zenfone 8 Flip', services: [{ name: 'Screen Replacement', grade: 'Original', price: 2500000 }, { name: 'Battery', grade: 'Original', price: 650000 }] },
    { id: 'as8', brand: 'Asus', model: 'ROG Phone 3', services: [{ name: 'Screen Replacement', grade: 'Original', price: 2200000 }, { name: 'Battery', grade: 'Original', price: 750000 }] },
    { id: 'as9', brand: 'Asus', model: 'Zenfone 7 Pro', services: [{ name: 'Screen Replacement', grade: 'Original', price: 2200000 }, { name: 'Battery', grade: 'Original', price: 650000 }] },
    { id: 'as10', brand: 'Asus', model: 'Zenfone 6', services: [{ name: 'Screen Replacement', grade: 'Original', price: 1200000 }, { name: 'Battery', grade: 'Original', price: 550000 }] },

    // Vivo
    { id: 'v1', brand: 'Vivo', model: 'X80 Pro', services: [{ name: 'Screen Replacement', grade: 'Original', type: 'OLED', price: 5200000 }, { name: 'Battery', grade: 'Original', price: 950000 }, { name: 'Rear Camera', grade: 'Original', price: 1800000 }] },
    { id: 'v2', brand: 'Vivo', model: 'X80', services: [{ name: 'Screen Replacement', grade: 'Original', type: 'OLED', price: 3500000 }, { name: 'Battery', grade: 'Original', price: 850000 }] },
    { id: 'v3', brand: 'Vivo', model: 'V23 5G', services: [{ name: 'Screen Replacement', grade: 'Original', price: 1300000 }, { name: 'Battery', grade: 'Original', price: 450000 }, { name: 'Front Camera', grade: 'Original', price: 750000 }] },
    { id: 'v4', brand: 'Vivo', model: 'V23e', services: [{ name: 'Screen Replacement', grade: 'Original', price: 950000 }, { name: 'Battery', grade: 'Original', price: 350000 }] },
    { id: 'v5', brand: 'Vivo', model: 'T1 Pro', services: [{ name: 'Screen Replacement', grade: 'Original', price: 1100000 }, { name: 'Battery', grade: 'Original', price: 450000 }] },
    { id: 'v6', brand: 'Vivo', model: 'Y35', services: [{ name: 'Screen Replacement', grade: 'Original', price: 650000 }, { name: 'Battery', grade: 'Original', price: 350000 }] },
    { id: 'v7', brand: 'Vivo', model: 'Y22', services: [{ name: 'Screen Replacement', grade: 'Original', price: 550000 }, { name: 'Battery', grade: 'Original', price: 250000 }] },
    { id: 'v8', brand: 'Vivo', model: 'Y16', services: [{ name: 'Screen Replacement', grade: 'Original', price: 450000 }, { name: 'Battery', grade: 'Original', price: 250000 }] },
    { id: 'v9', brand: 'Vivo', model: 'V21', services: [{ name: 'Screen Replacement', grade: 'Original', price: 1100000 }, { name: 'Battery', grade: 'Original', price: 450000 }] },
    { id: 'v10', brand: 'Vivo', model: 'X70 Pro', services: [{ name: 'Screen Replacement', grade: 'Original', type: 'OLED', price: 3200000 }, { name: 'Battery', grade: 'Original', price: 750000 }] }
  ];

  let selectedBrand = $state('All');
  let searchQuery = $state('');
  let selectedDevice: DevicePricing | null = $state(null);
  let selectedGrade: 'Original' | 'ODM' = $state('Original');

  const originalServices = $derived<RepairService[]>(
    selectedDevice ? (selectedDevice as DevicePricing).services.filter((s: RepairService) => s.grade === 'Original') : []
  );
  const odmServices = $derived<RepairService[]>(
    selectedDevice ? (selectedDevice as DevicePricing).services.filter((s: RepairService) => s.grade === 'ODM') : []
  );
  const currentServices = $derived<RepairService[]>(
    selectedGrade === 'Original' ? originalServices : odmServices
  );

  const brands = $derived(['All', ...new Set(devices.map(d => d.brand))]);

  let filteredDevices = $derived(
    devices.filter(d => {
      const matchesBrand = selectedBrand === 'All' || d.brand === selectedBrand;
      const matchesSearch = searchQuery.trim() === '' || 
        d.model.toLowerCase().includes(searchQuery.toLowerCase());
      return matchesBrand && matchesSearch;
    })
  );

  let visibleLimit = $state(12);
  let displayedDevices = $derived(filteredDevices.slice(0, visibleLimit));

  // Reset limit when filtering changes
  $effect(() => {
    selectedBrand;
    searchQuery;
    visibleLimit = 12;
  });

  const formatIDR = (amount: number) => {
    return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(amount);
  };

  const handleKeydown = (e: KeyboardEvent) => {
    if (e.key === 'Escape') selectedDevice = null;
  };

  $effect(() => {
    if (selectedDevice) {
      document.body.style.overflow = 'hidden';
    } else {
      document.body.style.overflow = '';
    }
  });
</script>

<svelte:window onkeydown={handleKeydown} />

<div class="flex flex-col min-h-screen bg-slate-50 dark:bg-slate-950 pt-24 pb-12">
  <!-- Header Section -->
  <section class="container mx-auto px-4 mb-12 text-center">
    <div class="inline-flex items-center gap-2 px-4 py-1.5 rounded-full bg-blue-100 dark:bg-blue-900/40 text-sm font-bold text-blue-700 dark:text-blue-300 mb-6">
      <AlertCircle size={16} />
      Flat Rp 50.000 Diagnosis Fee
    </div>
    
    <h1 class="text-4xl lg:text-5xl font-extrabold text-slate-900 dark:text-white mb-4">
      Transparent Pricing. <br class="hidden sm:block"/> Expert Repairs.
    </h1>
    <p class="text-lg text-slate-600 dark:text-slate-400 max-w-2xl mx-auto mb-10">
      Search for your device below to see estimated repair costs. No hidden fees. If you proceed with the repair, the parts and labor are exactly as quoted.
    </p>

    <!-- Brand Tabs -->
    <div class="flex flex-wrap justify-center gap-2 mb-8">
      {#each brands as brand}
        <button 
          onclick={() => selectedBrand = brand}
          class="px-6 py-2.5 rounded-full text-sm font-bold transition-all border-2 
            {selectedBrand === brand 
              ? 'bg-blue-600 border-blue-600 text-white shadow-lg shadow-blue-600/20 scale-105' 
              : 'bg-white dark:bg-slate-900 border-slate-200 dark:border-slate-800 text-slate-600 dark:text-slate-400 hover:border-blue-600/30'}"
        >
          {brand}
        </button>
      {/each}
    </div>

    <!-- Search Input -->
    <div class="max-w-xl mx-auto relative group">
      <div class="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none text-slate-400 group-focus-within:text-blue-600 transition-colors">
        <Search size={20} />
      </div>
      <input 
        bind:value={searchQuery}
        type="text" 
        placeholder={selectedBrand === 'All' ? 'Search all models...' : `Search ${selectedBrand} models...`}
        class="w-full bg-white dark:bg-slate-900 border-2 border-slate-200 dark:border-slate-800 rounded-2xl py-4 pl-12 pr-4 text-lg text-slate-900 dark:text-white shadow-sm focus:outline-none focus:border-blue-600 dark:focus:border-blue-500 transition-all"
      />
    </div>
  </section>

  <!-- Results Section -->
  <section class="container mx-auto px-4 flex-1">
    {#if filteredDevices.length > 0}
      <div class="grid md:grid-cols-2 lg:grid-cols-3 gap-6">
        {#each displayedDevices as device (device.id)}
          <button 
            onclick={() => {
              selectedDevice = device;
              selectedGrade = 'Original';
            }}
            class="group bg-slate-50/50 dark:bg-slate-900/50 rounded-2xl border border-slate-200 dark:border-slate-800 p-6 flex items-center gap-4 text-left hover:border-blue-600 dark:hover:border-blue-500 hover:shadow-xl hover:shadow-blue-600/5 transition-all active:scale-95"
          >
            <div class="w-14 h-14 rounded-2xl bg-slate-50 dark:bg-slate-800 flex items-center justify-center text-slate-600 dark:text-slate-400 group-hover:bg-blue-50 dark:group-hover:bg-blue-900/30 group-hover:text-blue-600 dark:group-hover:text-blue-400 transition-colors">
              <Smartphone size={28} />
            </div>
            <div class="flex-1">
              <p class="text-xs font-bold text-slate-400 uppercase tracking-wider">{device.brand}</p>
              <h3 class="text-xl font-bold text-slate-900 dark:text-white mb-1">{device.model}</h3>
              <p class="text-sm font-medium text-blue-600 flex items-center gap-1 opacity-0 group-hover:opacity-100 transition-opacity">
                View {device.services.length} Repairs <ArrowRight size={14}/>
              </p>
            </div>
          </button>
        {/each}
      </div>

      <!-- Show More Button -->
      {#if visibleLimit < filteredDevices.length}
        <div class="mt-12 text-center">
          <button 
            onclick={() => visibleLimit += 12}
            class="px-8 py-4 bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-xl font-bold text-slate-600 dark:text-slate-300 hover:border-blue-600 hover:text-blue-600 transition-all active:scale-95"
          >
            Show More Models ({filteredDevices.length - visibleLimit} remaining)
          </button>
        </div>
      {/if}
    {:else}
      <!-- Empty State -->
      <div class="max-w-lg mx-auto text-center py-16 px-4 bg-white dark:bg-slate-900 rounded-3xl border border-slate-200 dark:border-slate-800 shadow-sm">
        <div class="w-20 h-20 mx-auto bg-blue-50 dark:bg-blue-900/30 text-blue-600 dark:text-blue-400 rounded-2xl flex items-center justify-center mb-6">
          <Wrench size={40} />
        </div>
        <h3 class="text-2xl font-bold text-slate-900 dark:text-white mb-3">Device not listed?</h3>
        <p class="text-slate-600 dark:text-slate-400 mb-8 leading-relaxed">
          Don't worry, we fix almost everything! Bring your device in for a professional assessment. We only charge a flat <strong>Rp 50.000</strong> diagnosis fee, and we'll give you a custom quote before doing any work.
        </p>
        <a href="/book" class="inline-flex items-center gap-2 px-8 py-4 bg-blue-600 text-white font-bold rounded-xl shadow-lg shadow-blue-600/20 hover:bg-blue-700 transition-all active:scale-95">
          Book Custom Diagnosis
        </a>
      </div>
    {/if}
  </section>

  <!-- Guarantee Footer -->
  <section class="container mx-auto px-4 mt-20">
    <div class="bg-blue-600 rounded-3xl p-8 sm:p-12 text-white flex flex-col md:flex-row items-center gap-8 md:gap-12">
      <div class="w-16 h-16 shrink-0 bg-white/20 rounded-2xl flex items-center justify-center">
        <ShieldCheck size={32} class="text-white" />
      </div>
      <div>
        <h3 class="text-2xl font-bold mb-3">The OpenBench Guarantee</h3>
        <p class="text-blue-100 leading-relaxed max-w-3xl">
          Transparency is our policy. The prices above include parts and professional installation. 
          <strong>Original Parts</strong> come with a strict 30-day warranty against defects. 
          <strong>Premium ODM Parts</strong> come with a 7-day guarantee. If we can't fix it, you only pay the Rp 50.000 diagnosis fee.
        </p>
      </div>
    </div>
  </section>

  <!-- Bottom Sheet / Modal -->
  {#if selectedDevice}
    <!-- Backdrop -->
    <!-- svelte-ignore a11y_click_events_have_key_events -->
    <!-- svelte-ignore a11y_no_static_element_interactions -->
    <div 
      transition:fade={{ duration: 200 }}
      onclick={() => selectedDevice = null}
      class="fixed inset-0 bg-slate-900/60 backdrop-blur-sm z-[60]"
    ></div>

    <!-- Sheet -->
    <div 
      transition:fly={{ y: 500, duration: 400, opacity: 1 }}
      class="fixed bottom-0 left-0 right-0 max-h-[85vh] bg-white dark:bg-slate-950 rounded-t-[40px] z-[70] shadow-2xl border-t border-slate-200 dark:border-slate-800 overflow-hidden flex flex-col"
    >
      <!-- Sheet Header -->
      <div class="p-8 pb-4 flex items-center justify-between">
        <div class="flex items-center gap-4">
          <div class="w-16 h-16 rounded-2xl bg-blue-50 dark:bg-blue-900/30 text-blue-600 dark:text-blue-400 flex items-center justify-center">
            <Smartphone size={32} />
          </div>
          <div>
            <p class="text-sm font-bold text-blue-600 uppercase tracking-widest">{selectedDevice.brand}</p>
            <h2 class="text-3xl font-extrabold text-slate-900 dark:text-white">{selectedDevice.model}</h2>
          </div>
        </div>
        <button 
          onclick={() => selectedDevice = null}
          class="w-12 h-12 rounded-full bg-slate-100 dark:bg-slate-800 text-slate-500 flex items-center justify-center hover:bg-slate-200 dark:hover:bg-slate-700 transition-colors"
        >
          <X size={24} />
        </button>
      </div>

      <!-- Sheet Content (Scrollable) -->
      <div class="flex-1 overflow-y-auto p-8 pt-4">
        <div class="flex items-center gap-3 text-slate-500 mb-8 bg-slate-50 dark:bg-slate-900/50 p-4 rounded-2xl">
          <AlertCircle size={18} />
          <span class="text-sm font-medium italic">Fixed Rp 50.000 diagnosis fee applies to all bookings.</span>
        </div>

        <!-- Grade Toggle Switch -->
        <div class="bg-slate-100 dark:bg-slate-900 p-1.5 rounded-2xl mb-10 flex gap-1">
          <button 
            onclick={() => selectedGrade = 'Original'}
            class="flex-1 py-3 px-4 rounded-xl text-sm font-bold transition-all
              {selectedGrade === 'Original' 
                ? 'bg-blue-600 text-white shadow-lg shadow-blue-600/20' 
                : 'text-slate-500 hover:text-slate-700 dark:hover:text-slate-300'}"
          >
            Original Standards
          </button>
          <button 
            onclick={() => selectedGrade = 'ODM'}
            class="flex-1 py-3 px-4 rounded-xl text-sm font-bold transition-all
              {selectedGrade === 'ODM' 
                ? 'bg-amber-500 text-white shadow-lg shadow-amber-500/20' 
                : 'text-slate-500 hover:text-slate-700 dark:hover:text-slate-300'}"
          >
            Premium ODM
          </button>
        </div>

        <div class="space-y-6">
          <div class="flex items-center justify-between mb-4 px-2">
            <h3 class="text-sm font-black text-slate-400 uppercase tracking-widest">
              {selectedGrade} Price List
            </h3>
            <span class="text-xs font-bold px-2 py-1 {selectedGrade === 'Original' ? 'bg-blue-100 text-blue-600' : 'bg-amber-100 text-amber-600'} rounded-lg">
              {selectedGrade === 'Original' ? '30-DAY WARRANTY' : '7-DAY GUARANTEE'}
            </span>
          </div>

          {#if currentServices.length > 0}
            <div class="grid gap-4">
              {#each currentServices as service}
                <div class="bg-white dark:bg-slate-900 p-6 rounded-2xl border border-slate-200 dark:border-slate-800 flex flex-col sm:flex-row sm:items-center justify-between gap-6 hover:border-blue-600/30 transition-colors">
                  <div>
                    <h4 class="font-bold text-slate-900 dark:text-white mb-1">
                      {service.name} 
                      {#if service.type && service.type !== 'Standard'}
                        <span class="ml-2 text-xs font-black px-1.5 py-0.5 bg-slate-100 dark:bg-slate-800 rounded">{service.type}</span>
                      {/if}
                    </h4>
                    <p class="text-sm text-slate-500">
                      {selectedGrade === 'Original' ? 'Genuine factory-spec component.' : 'High-grade verified aftermarket part.'}
                    </p>
                  </div>
                  <div class="flex items-center justify-between sm:justify-end gap-6 w-full sm:w-auto border-t sm:border-t-0 pt-4 sm:pt-0 border-slate-100 dark:border-slate-800">
                    <span class="text-xl font-black text-slate-900 dark:text-white">{formatIDR(service.price)}</span>
                    <a href="/book" class="px-6 py-3 {selectedGrade === 'Original' ? 'bg-blue-600' : 'bg-amber-500'} text-white font-bold rounded-xl hover:opacity-90 transition-all active:scale-95 shadow-lg">
                      Book
                    </a>
                  </div>
                </div>
              {/each}
            </div>
          {:else}
            <div class="text-center py-12 px-6 bg-slate-50 dark:bg-slate-900/50 rounded-3xl border-2 border-dashed border-slate-200 dark:border-slate-800">
              <AlertCircle class="mx-auto text-slate-300 mb-4" size={48} />
              <p class="text-slate-500 font-medium">No {selectedGrade} parts available for this model yet.</p>
              <button onclick={() => selectedGrade = selectedGrade === 'Original' ? 'ODM' : 'Original'} class="mt-4 text-blue-600 font-bold hover:underline">
                View {selectedGrade === 'Original' ? 'ODM' : 'Original'} prices instead
              </button>
            </div>
          {/if}
        </div>

        <div class="mt-12 flex items-start gap-4 p-6 bg-blue-50/50 dark:bg-blue-900/10 rounded-2xl border border-blue-100/50 dark:border-blue-900/30">
          <ShieldCheck class="text-blue-600 shrink-0 mt-1" size={24} />
          <p class="text-sm text-slate-600 dark:text-slate-400 leading-relaxed">
            <strong>Transparency Note:</strong> {selectedGrade === 'Original' 
              ? 'Original parts are identical to those installed at the factory. Ideal for maintaining resale value and peak performance.' 
              : 'ODM parts are high-quality third-party components that pass strict quality checks. Perfect for cost-effective repairs.'}
          </p>
        </div>
      </div>
      
      <!-- Safe Area Spacer for Mobile -->
      <div class="h-8"></div>
    </div>
  {/if}
</div>
