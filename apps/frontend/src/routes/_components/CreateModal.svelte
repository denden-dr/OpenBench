<script lang="ts">
  import { X, Smartphone, User, Wrench, LoaderCircle } from "lucide-svelte";

  let {
    isOpen,
    isCreating,
    onSubmit,
    onClose,
  }: {
    isOpen: boolean;
    isCreating: boolean;
    onSubmit: (form: any) => void;
    onClose: () => void;
  } = $props();

  let customer_name = $state("");
  let customer_phone = $state("");
  let customer_gender = $state("Male");
  let brand = $state("");
  let model = $state("");
  let issue = $state("");
  let additional_description = $state("");
  let accessories = $state("");
  let price = $state(0);
  let warranty_days = $state(0);

  function handleSubmit(e: SubmitEvent) {
    e.preventDefault();
    onSubmit({
      customer_name,
      customer_phone,
      customer_gender,
      brand,
      model,
      issue,
      additional_description,
      accessories,
      price,
      warranty_days,
    });
    // Reset form after submission
    customer_name = "";
    customer_phone = "";
    customer_gender = "Male";
    brand = "";
    model = "";
    issue = "";
    additional_description = "";
    accessories = "";
    price = 0;
    warranty_days = 0;
  }
</script>

<svelte:window onkeydown={(e) => { if (e.key === 'Escape' && isOpen) onClose(); }} />

{#if isOpen}
  <div class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-900/60 backdrop-blur-sm animate-fade-in">
    <div class="bg-white dark:bg-slate-950 w-full max-w-2xl rounded-2xl border border-slate-200 dark:border-slate-800 shadow-2xl overflow-hidden flex flex-col max-h-[90vh]">
      <!-- Header -->
      <div class="flex justify-between items-center px-6 py-4 border-b border-slate-200 dark:border-slate-800 bg-slate-50 dark:bg-slate-900/30">
        <h3 class="font-bold text-lg text-slate-900 dark:text-white flex items-center gap-2">
          <Smartphone size={20} class="text-blue-600" />
          Pendaftaran Servis
        </h3>
        <button
          onclick={onClose}
          disabled={isCreating}
          class="p-1.5 rounded-lg text-slate-400 hover:bg-slate-100 dark:hover:bg-slate-900 transition-colors disabled:opacity-50"
        >
          <X size={18} />
        </button>
      </div>

      <!-- Body -->
      <form
        onsubmit={handleSubmit}
        class="flex-1 overflow-y-auto p-6 space-y-6 text-left"
      >
        <!-- Customer Info -->
        <div>
          <h4 class="text-xs font-bold text-slate-400 uppercase tracking-widest mb-3 flex items-center gap-1.5">
            <User size={14} />
            Informasi Pelanggan
          </h4>
          <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <div class="space-y-1.5">
              <label for="create-customer-name" class="text-xs font-bold text-slate-500 block">Nama Lengkap *</label>
              <input
                id="create-customer-name"
                type="text"
                bind:value={customer_name}
                required
                placeholder="Customer Name"
                class="w-full bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-xl py-2 px-3 text-sm text-slate-900 dark:text-white focus:outline-none focus:border-blue-600 transition-colors"
              />
            </div>
            <div class="space-y-1.5">
              <label for="create-customer-phone" class="text-xs font-bold text-slate-500 block">Nomor HP</label>
              <input
                id="create-customer-phone"
                type="text"
                bind:value={customer_phone}
                placeholder="0812..."
                class="w-full bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-xl py-2 px-3 text-sm text-slate-900 dark:text-white focus:outline-none focus:border-blue-600 transition-colors"
              />
            </div>
            <div class="space-y-1.5">
              <label for="create-customer-gender" class="text-xs font-bold text-slate-500 block">Jenis Kelamin *</label>
              <select
                id="create-customer-gender"
                bind:value={customer_gender}
                class="w-full bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-xl py-2 px-3 text-sm text-slate-900 dark:text-white focus:outline-none focus:border-blue-600 transition-colors cursor-pointer"
              >
                <option value="Male">Laki-laki</option>
                <option value="Female">Perempuan</option>
                <option value="Other">Other</option>
              </select>
            </div>
          </div>
        </div>

        <!-- Detail Perangkat -->
        <div>
          <h4 class="text-xs font-bold text-slate-400 uppercase tracking-widest mb-3 flex items-center gap-1.5">
            <Wrench size={14} />
            Detail Perangkat
          </h4>
          <div class="grid grid-cols-1 sm:grid-cols-2 gap-4 mb-4">
            <div class="space-y-1.5">
              <label for="create-brand" class="text-xs font-bold text-slate-500 block">Merek *</label>
              <input
                id="create-brand"
                type="text"
                bind:value={brand}
                required
                placeholder="e.g. Apple, Samsung"
                class="w-full bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-xl py-2 px-3 text-sm text-slate-900 dark:text-white focus:outline-none focus:border-blue-600 transition-colors"
              />
            </div>
            <div class="space-y-1.5">
              <label for="create-model" class="text-xs font-bold text-slate-500 block">Tipe / Model *</label>
              <input
                id="create-model"
                type="text"
                bind:value={model}
                required
                placeholder="e.g. iPhone 15 Pro, Galaxy S24"
                class="w-full bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-xl py-2 px-3 text-sm text-slate-900 dark:text-white focus:outline-none focus:border-blue-600 transition-colors"
              />
            </div>
          </div>

          <div class="space-y-4">
            <div class="space-y-1.5">
              <label for="create-issue" class="text-xs font-bold text-slate-500 block">Kerusakan *</label>
              <input
                id="create-issue"
                type="text"
                bind:value={issue}
                required
                placeholder="e.g. Broken LCD screen, Battery drain"
                class="w-full bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-xl py-2 px-3 text-sm text-slate-900 dark:text-white focus:outline-none focus:border-blue-600 transition-colors"
              />
            </div>
            <div class="space-y-1.5">
              <label for="create-additional-description" class="text-xs font-bold text-slate-500 block">Deskripsi Tambahan</label>
              <textarea
                id="create-additional-description"
                bind:value={additional_description}
                rows="3"
                placeholder="Any cosmetic damage or technical background details..."
                class="w-full bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-xl py-2 px-3 text-sm text-slate-900 dark:text-white focus:outline-none focus:border-blue-600 transition-colors resize-none"
              ></textarea>
            </div>
            <div class="space-y-1.5">
              <label for="create-accessories" class="text-xs font-bold text-slate-500 block">Kelengkapan Bawaan</label>
              <input
                id="create-accessories"
                type="text"
                bind:value={accessories}
                placeholder="e.g. Charger, Case, SIM card tray"
                class="w-full bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-xl py-2 px-3 text-sm text-slate-900 dark:text-white focus:outline-none focus:border-blue-600 transition-colors"
              />
            </div>
          </div>
        </div>

        <!-- Harga & Garansi -->
        <div>
          <h4 class="text-xs font-bold text-slate-400 uppercase tracking-widest mb-3 flex items-center gap-1.5">
            <span class="text-[10px] font-black bg-slate-100 dark:bg-slate-800 px-1.5 py-0.5 rounded text-slate-500 dark:text-slate-400 mr-1 select-none">RP</span>
            Harga & Garansi
          </h4>
          <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <div class="space-y-1.5">
              <label for="create-price" class="text-xs font-bold text-slate-500 block">Estimasi Harga (Rp)</label>
              <input
                id="create-price"
                type="number"
                bind:value={price}
                min="0"
                class="w-full bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-xl py-2 px-3 text-sm text-slate-900 dark:text-white focus:outline-none focus:border-blue-600 transition-colors"
              />
            </div>
            <div class="space-y-1.5">
              <label for="create-warranty-days" class="text-xs font-bold text-slate-500 block">Masa Garansi (Hari)</label>
              <input
                id="create-warranty-days"
                type="number"
                bind:value={warranty_days}
                min="0"
                class="w-full bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-xl py-2 px-3 text-sm text-slate-900 dark:text-white focus:outline-none focus:border-blue-600 transition-colors"
              />
            </div>
          </div>
        </div>

        <!-- Footer -->
        <div class="pt-6 border-t border-slate-200 dark:border-slate-800 flex justify-end gap-3">
          <button
            type="button"
            onclick={onClose}
            disabled={isCreating}
            class="px-4 py-2 border border-slate-200 dark:border-slate-800 text-slate-700 dark:text-slate-300 font-bold text-xs uppercase tracking-wider rounded-xl hover:bg-slate-50 dark:hover:bg-slate-900 transition-colors disabled:opacity-50"
          >
            Cancel
          </button>
          <button
            type="submit"
            disabled={isCreating}
            class="px-6 py-2 bg-blue-600 hover:bg-blue-700 text-white font-bold text-xs uppercase tracking-wider rounded-xl transition-colors shadow-sm inline-flex items-center gap-1.5 disabled:opacity-50"
          >
            {#if isCreating}
              <LoaderCircle class="animate-spin" size={14} />
              Menyimpan...
            {:else}
              Submit Intake
            {/if}
          </button>
        </div>
      </form>
    </div>
  </div>
{/if}
