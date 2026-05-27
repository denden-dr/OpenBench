<script lang="ts">
  import { X, Smartphone, User, Wrench, LoaderCircle, AlertTriangle, PhoneCall, Trash2 } from "lucide-svelte";
  import type { Ticket } from "$lib/types/ticket";
  import { formatCurrency, formatDate } from "$lib/utils/format";

  let {
    ticket,
    isUpdating,
    isDeleting,
    onClose,
    onSubmitUpdate,
    onDelete,
    copyTrackingLink,
    openQRModal,
    submitIssue,
    approveByCustomer,
    rejectByCustomer,
    getStatusBadgeClass,
    getStatusLabel,
  }: {
    ticket: Ticket | null;
    isUpdating: boolean;
    isDeleting: boolean;
    onClose: () => void;
    onSubmitUpdate: (updatedForm: any) => void;
    onDelete: () => void;
    copyTrackingLink: (ticketId: string) => void;
    openQRModal: (ticketId: string) => void;
    submitIssue: (issueType: string, newPrice: number, newDiagnosis: string) => void;
    approveByCustomer: () => void;
    rejectByCustomer: () => void;
    getStatusBadgeClass: (status: string, isWarranty: boolean) => string;
    getStatusLabel: (status: string, isWarranty: boolean) => string;
  } = $props();

  // local form states derived from ticket when it changes
  let editForm = $state({
    customer_name: "",
    customer_phone: "",
    customer_gender: "Male",
    brand: "",
    model: "",
    issue: "",
    additional_description: "",
    accessories: "",
    price: 0,
    warranty_days: 0,
    status: "",
    payment_status: "",
  });

  $effect(() => {
    if (ticket) {
      editForm.customer_name = ticket.customer_name;
      editForm.customer_phone = ticket.customer_phone || "";
      editForm.customer_gender = ticket.customer_gender;
      editForm.brand = ticket.brand;
      editForm.model = ticket.model;
      editForm.issue = ticket.issue;
      editForm.additional_description = ticket.additional_description || "";
      editForm.accessories = ticket.accessories || "";
      editForm.price = ticket.price;
      editForm.warranty_days = ticket.warranty_days;
      editForm.status = ticket.status;
      editForm.payment_status = ticket.payment_status;
    }
  });

  $effect(() => {
    if (editForm.status === "picked_up" && editForm.payment_status !== "paid") {
      editForm.payment_status = "paid";
    }
  });

  // issue states
  let issueType = $state("");
  let newPrice = $state(0);
  let newDiagnosis = $state("");

  function handleSubmit(e: SubmitEvent) {
    e.preventDefault();
    onSubmitUpdate(editForm);
  }

  function handleIssueSubmit() {
    submitIssue(issueType, newPrice, newDiagnosis);
    issueType = "";
    newPrice = 0;
    newDiagnosis = "";
  }
</script>

<svelte:window onkeydown={(e) => { if (e.key === 'Escape' && ticket) onClose(); }} />

{#if ticket}
  <!-- svelte-ignore a11y_click_events_have_key_events -->
  <!-- svelte-ignore a11y_no_static_element_interactions -->
  <div
    class="fixed inset-0 z-50 flex justify-end bg-slate-900/60 backdrop-blur-sm animate-fade-in font-sans"
    onclick={onClose}
  >
    <!-- svelte-ignore a11y_click_events_have_key_events -->
    <!-- svelte-ignore a11y_no_static_element_interactions -->
    <div
      class="bg-white dark:bg-slate-950 w-full max-w-lg border-l border-slate-200 dark:border-slate-800 shadow-2xl h-full flex flex-col animate-slide-in"
      onclick={(e) => e.stopPropagation()}
    >
      <!-- Header -->
      <div
        class="flex justify-between items-center px-6 py-4 border-b border-slate-200 dark:border-slate-800 bg-slate-50 dark:bg-slate-900/30"
      >
        <div class="text-left">
          <h3 class="font-bold text-lg text-slate-900 dark:text-white">
            Detail Perbaikan
          </h3>
          <p class="text-xs text-slate-500 dark:text-slate-400 mt-0.5 font-mono">
            ID: {ticket.id}
          </p>
        </div>
        <button
          onclick={onClose}
          disabled={isUpdating || isDeleting}
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
        <!-- Status Panel -->
        <div class="bg-slate-50 dark:bg-slate-900/40 p-4 rounded-xl border border-slate-200 dark:border-slate-850 space-y-3">
          <div class="flex justify-between items-center">
            <span class="text-xs font-bold text-slate-500">Status Saat Ini</span>
            <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-semibold border {getStatusBadgeClass(editForm.status, Boolean(ticket.is_warranty))}">
              {getStatusLabel(editForm.status, Boolean(ticket.is_warranty))}
            </span>
          </div>
          <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
            <div class="space-y-1">
              <span class="text-[10px] font-black uppercase text-slate-400">Tanggal Masuk</span>
              <p class="text-xs font-semibold text-slate-700 dark:text-slate-300">
                {formatDate(ticket.entry_date)}
              </p>
            </div>
            <div class="space-y-1">
              <span class="text-[10px] font-black uppercase text-slate-400">Tanggal Keluar</span>
              <p class="text-xs font-semibold text-slate-700 dark:text-slate-300">
                {ticket.exit_date ? formatDate(ticket.exit_date) : "-"}
              </p>
            </div>
          </div>
          {#if ticket.warranty_expiry_date}
            <div class="pt-2 border-t border-slate-200/50 dark:border-slate-800/50 flex justify-between items-center">
              <span class="text-xs font-bold text-slate-500">Batas Garansi</span>
              <span class="text-xs font-bold text-indigo-600 dark:text-indigo-400">
                {formatDate(ticket.warranty_expiry_date)}
              </span>
            </div>
          {/if}
        </div>

        <!-- Tracking Section -->
        <div class="bg-slate-50 dark:bg-slate-900/40 p-4 rounded-xl border border-slate-200 dark:border-slate-800 space-y-3">
          <span class="text-xs font-bold text-slate-500 block">Pelacakan Pelanggan</span>
          <div class="flex gap-2">
            <button
              type="button"
              onclick={() => copyTrackingLink(ticket.id)}
              class="flex-1 py-2 border border-slate-200 dark:border-slate-800 text-slate-700 dark:text-slate-300 text-xs font-bold rounded-xl hover:bg-slate-100 dark:hover:bg-slate-900 transition-colors cursor-pointer"
            >
              Salin Link
            </button>
            <button
              type="button"
              onclick={() => openQRModal(ticket.id)}
              class="flex-1 py-2 bg-blue-600 hover:bg-blue-700 text-white text-xs font-bold rounded-xl transition-all shadow-sm active:scale-95 cursor-pointer"
            >
              Tampilkan QR Code
            </button>
          </div>
        </div>

        <!-- Customer Form Fields -->
        <div>
          <h4 class="text-xs font-bold text-slate-400 uppercase tracking-widest mb-3">Profil Pelanggan</h4>
          <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <div class="space-y-1.5">
              <label for="edit-customer-name" class="text-xs font-bold text-slate-500">Nama Pelanggan</label>
              <input
                id="edit-customer-name"
                type="text"
                bind:value={editForm.customer_name}
                class="w-full bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-xl py-2 px-3 text-sm text-slate-900 dark:text-white focus:outline-none focus:border-blue-600 transition-colors"
              />
            </div>
            <div class="space-y-1.5">
              <label for="edit-customer-phone" class="text-xs font-bold text-slate-500">Nomor HP</label>
              <input
                id="edit-customer-phone"
                type="text"
                bind:value={editForm.customer_phone}
                placeholder="0812..."
                class="w-full bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-xl py-2 px-3 text-sm text-slate-900 dark:text-white focus:outline-none focus:border-blue-600 transition-colors"
              />
            </div>
            <div class="space-y-1.5">
              <label for="edit-customer-gender" class="text-xs font-bold text-slate-500">Jenis Kelamin</label>
              <select
                id="edit-customer-gender"
                bind:value={editForm.customer_gender}
                class="w-full bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-xl py-2 px-3 text-sm text-slate-900 dark:text-white focus:outline-none focus:border-blue-600 transition-colors cursor-pointer"
              >
                <option value="Male">Laki-laki</option>
                <option value="Female">Perempuan</option>
                <option value="Other">Other</option>
              </select>
            </div>
          </div>
        </div>

        <!-- Device Info -->
        <div>
          <h4 class="text-xs font-bold text-slate-400 uppercase tracking-widest mb-3">Informasi Servis</h4>
          <div class="grid grid-cols-2 gap-4 mb-4">
            <div class="space-y-1.5">
              <label for="edit-brand" class="text-xs font-bold text-slate-500">Merek</label>
              <input
                id="edit-brand"
                type="text"
                bind:value={editForm.brand}
                class="w-full bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-xl py-2 px-3 text-sm text-slate-900 dark:text-white focus:outline-none focus:border-blue-600 transition-colors"
              />
            </div>
            <div class="space-y-1.5">
              <label for="edit-model" class="text-xs font-bold text-slate-500">Model</label>
              <input
                id="edit-model"
                type="text"
                bind:value={editForm.model}
                class="w-full bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-xl py-2 px-3 text-sm text-slate-900 dark:text-white focus:outline-none focus:border-blue-600 transition-colors"
              />
            </div>
          </div>

          <div class="space-y-4">
            <div class="space-y-1.5">
              <label for="edit-issue" class="text-xs font-bold text-slate-500">Kerusakan</label>
              <input
                id="edit-issue"
                type="text"
                bind:value={editForm.issue}
                class="w-full bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-xl py-2 px-3 text-sm text-slate-900 dark:text-white focus:outline-none focus:border-blue-600 transition-colors"
              />
            </div>
            <div class="space-y-1.5">
              <label for="edit-additional-description" class="text-xs font-bold text-slate-500">Deskripsi Tambahan</label>
              <textarea
                id="edit-additional-description"
                bind:value={editForm.additional_description}
                rows="3"
                class="w-full bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-xl py-2 px-3 text-sm text-slate-900 dark:text-white focus:outline-none focus:border-blue-600 transition-colors resize-none"
              ></textarea>
            </div>
            <div class="space-y-1.5">
              <label for="edit-accessories" class="text-xs font-bold text-slate-500">Kelengkapan</label>
              <input
                id="edit-accessories"
                type="text"
                bind:value={editForm.accessories}
                class="w-full bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-xl py-2 px-3 text-sm text-slate-900 dark:text-white focus:outline-none focus:border-blue-600 transition-colors"
              />
            </div>
          </div>
        </div>

        <!-- Finances & Status overrides -->
        <div>
          <h4 class="text-xs font-bold text-slate-400 uppercase tracking-widest mb-3">Biaya & Parameter</h4>
          <div class="grid grid-cols-2 gap-4 mb-4">
            <div class="space-y-1.5">
              <label for="edit-price" class="text-xs font-bold text-slate-500">Harga (Rp)</label>
              <input
                id="edit-price"
                type="number"
                bind:value={editForm.price}
                class="w-full bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-xl py-2 px-3 text-sm text-slate-900 dark:text-white focus:outline-none focus:border-blue-600 transition-colors"
              />
            </div>
            <div class="space-y-1.5">
              <label for="edit-warranty-days" class="text-xs font-bold text-slate-500">Masa Garansi (Hari)</label>
              <input
                id="edit-warranty-days"
                type="number"
                bind:value={editForm.warranty_days}
                class="w-full bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-xl py-2 px-3 text-sm text-slate-900 dark:text-white focus:outline-none focus:border-blue-600 transition-colors"
              />
            </div>
          </div>

          <div class="grid grid-cols-2 gap-4">
            <div class="space-y-1.5">
              <label for="edit-status" class="text-xs font-bold text-slate-500">Status</label>
              <select
                id="edit-status"
                bind:value={editForm.status}
                class="w-full bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-xl py-2.5 px-3 text-sm text-slate-900 dark:text-white focus:outline-none focus:border-blue-600 transition-colors cursor-pointer"
              >
                <option value="service_in">Masuk</option>
                <option value="on_process">Sedang Diproses</option>
                <option value="waiting_confirmation">Menunggu Konfirmasi</option>
                <option value="cancelled">Dibatalkan</option>
                <option value="fixed">Selesai Diperbaiki</option>
                <option value="picked_up">Sudah Diambil</option>
              </select>
            </div>
            <div class="space-y-1.5">
              <label for="edit-payment-status" class="text-xs font-bold text-slate-500">Status Pembayaran</label>
              <select
                id="edit-payment-status"
                bind:value={editForm.payment_status}
                class="w-full bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-xl py-2.5 px-3 text-sm text-slate-900 dark:text-white focus:outline-none focus:border-blue-600 transition-colors cursor-pointer"
              >
                <option value="unpaid" disabled={editForm.status === "picked_up"}>Unpaid</option>
                <option value="paid">Lunas</option>
              </select>
            </div>
          </div>
        </div>

        {#if editForm.status === "on_process"}
          <div class="mt-6 border-t border-slate-200 dark:border-slate-800 pt-6">
            <h4 class="text-xs font-bold text-rose-500 uppercase tracking-widest mb-3 flex items-center gap-1.5">
              <AlertTriangle size={14} />
              Kendala Teknisi (Opsional)
            </h4>

            <div class="space-y-4 bg-rose-50/50 dark:bg-rose-900/10 p-4 rounded-xl border border-rose-100 dark:border-rose-800/50">
              <div class="space-y-1.5">
                <label for="issue-type" class="text-xs font-bold text-slate-500 block">Pilih Jenis Kendala</label>
                <select
                  id="issue-type"
                  bind:value={issueType}
                  class="w-full bg-white dark:bg-slate-950 border border-slate-200 dark:border-slate-800 rounded-xl py-2 px-3 text-sm text-slate-900 dark:text-white focus:outline-none focus:border-rose-500 transition-colors cursor-pointer"
                >
                  <option value="">-- Pilih Kendala --</option>
                  <option value="unrepairable">1. Tidak dapat dilanjutkan (Rusak Total)</option>
                  <option value="additional_damage">2. Kerusakan Tambahan</option>
                  <option value="wrong_diagnosis">3. Salah Diagnosa Awal</option>
                </select>
              </div>

              {#if issueType === "additional_damage" || issueType === "wrong_diagnosis"}
                <div class="space-y-1.5">
                  <label for="new-price" class="text-xs font-bold text-slate-500 block">Penyesuaian Harga Baru (IDR)</label>
                  <input
                    id="new-price"
                    type="number"
                    bind:value={newPrice}
                    class="w-full bg-white dark:bg-slate-950 border border-slate-200 dark:border-slate-800 rounded-xl py-2 px-3 text-sm text-slate-900 dark:text-white focus:outline-none focus:border-rose-500 transition-colors"
                  />
                </div>
                <div class="space-y-1.5">
                  <label for="new-diagnosis" class="text-xs font-bold text-slate-500 block">Detail Diagnosa Baru</label>
                  <textarea
                    id="new-diagnosis"
                    bind:value={newDiagnosis}
                    placeholder="Contoh: Ternyata IC Power ikut rusak..."
                    class="w-full bg-white dark:bg-slate-950 border border-slate-200 dark:border-slate-800 rounded-xl py-2 px-3 text-sm text-slate-900 dark:text-white focus:outline-none focus:border-rose-500 transition-colors resize-none"
                  ></textarea>
                </div>

                <button
                  type="button"
                  onclick={handleIssueSubmit}
                  disabled={isUpdating || isDeleting}
                  class="w-full py-2 bg-rose-600 hover:bg-rose-700 text-white font-bold text-xs rounded-xl transition-all shadow-sm active:scale-95 inline-flex items-center justify-center gap-1.5 disabled:opacity-50"
                >
                  {#if isUpdating}
                    <LoaderCircle class="animate-spin" size={14} />
                    Menyimpan...
                  {:else}
                    Kirim ke Status Menunggu Konfirmasi
                  {/if}
                </button>
              {/if}

              {#if issueType === "unrepairable"}
                <button
                  type="button"
                  onclick={handleIssueSubmit}
                  disabled={isUpdating || isDeleting}
                  class="w-full py-2 bg-red-600 hover:bg-red-700 text-white font-bold text-xs rounded-xl transition-all shadow-sm active:scale-95 inline-flex items-center justify-center gap-1.5 disabled:opacity-50"
                >
                  {#if isUpdating}
                    <LoaderCircle class="animate-spin" size={14} />
                    Menyimpan...
                  {:else}
                    Batalkan Perbaikan (Kembalikan ke Customer)
                  {/if}
                </button>
              {/if}
            </div>
          </div>
        {/if}

        {#if editForm.status === "waiting_confirmation"}
          <div class="mt-6 border-t border-slate-200 dark:border-slate-800 pt-6">
            <div class="bg-amber-100 dark:bg-amber-900/30 p-4 rounded-xl border border-amber-200 dark:border-amber-800 space-y-4">
              <div class="flex gap-3 items-start">
                <PhoneCall class="text-amber-600 dark:text-amber-500 mt-1" size={20} />
                <div>
                  <h4 class="font-bold text-amber-900 dark:text-amber-100">Tindakan Diperlukan</h4>
                  <p class="text-xs text-amber-700 dark:text-amber-300 mt-1">
                    Terdapat perubahan harga atau kerusakan. Harap hubungi customer untuk meminta persetujuan.
                  </p>

                  <div class="mt-3 p-3 bg-white/60 dark:bg-black/20 rounded-lg text-sm">
                    <p class="text-slate-700 dark:text-slate-300">
                      <strong>Harga Baru:</strong> {formatCurrency(editForm.price)}
                    </p>
                    <p class="text-slate-700 dark:text-slate-300 mt-1 whitespace-pre-wrap">
                      <strong>Alasan/Keterangan:</strong><br />{editForm.additional_description}
                    </p>
                  </div>
                </div>
              </div>

              <div class="flex gap-2 pt-2">
                <button
                  type="button"
                  onclick={approveByCustomer}
                  disabled={isUpdating || isDeleting}
                  class="flex-1 py-2 bg-emerald-600 hover:bg-emerald-700 text-white text-xs font-bold rounded-xl transition-all shadow-sm active:scale-95 inline-flex items-center justify-center gap-1.5 disabled:opacity-50"
                >
                  {#if isUpdating}
                    <LoaderCircle class="animate-spin" size={14} />
                  {/if}
                  Pelanggan Setuju (Lanjut)
                </button>
                <button
                  type="button"
                  onclick={rejectByCustomer}
                  disabled={isUpdating || isDeleting}
                  class="flex-1 py-2 bg-red-600 hover:bg-red-700 text-white text-xs font-bold rounded-xl transition-all shadow-sm active:scale-95 inline-flex items-center justify-center gap-1.5 disabled:opacity-50"
                >
                  {#if isUpdating}
                    <LoaderCircle class="animate-spin" size={14} />
                  {/if}
                  Pelanggan Menolak (Batal)
                </button>
              </div>
            </div>
          </div>
        {/if}

        <!-- Footer Actions -->
        <div class="pt-6 border-t border-slate-200 dark:border-slate-800 flex justify-between items-center">
          <button
            type="button"
            onclick={onDelete}
            disabled={isUpdating || isDeleting}
            class="px-4 py-2 bg-red-50 hover:bg-red-100 text-red-600 font-bold text-xs uppercase tracking-wider rounded-xl transition-colors inline-flex items-center gap-1.5 border border-red-200 disabled:opacity-50"
          >
            {#if isDeleting}
              <LoaderCircle class="animate-spin" size={14} />
              Deleting...
            {:else}
              <Trash2 size={14} />
              Delete Ticket
            {/if}
          </button>

          <div class="flex gap-3">
            <button
              type="button"
              onclick={onClose}
              disabled={isUpdating || isDeleting}
              class="px-4 py-2 border border-slate-200 dark:border-slate-800 text-slate-700 dark:text-slate-300 font-bold text-xs uppercase tracking-wider rounded-xl hover:bg-slate-50 dark:hover:bg-slate-900 transition-colors disabled:opacity-50"
            >
              Cancel
            </button>
            <button
              type="submit"
              disabled={isUpdating || isDeleting}
              class="px-6 py-2 bg-blue-600 hover:bg-blue-700 text-white font-bold text-xs uppercase tracking-wider rounded-xl transition-colors shadow-sm inline-flex items-center gap-1.5 disabled:opacity-50"
            >
              {#if isUpdating}
                <LoaderCircle class="animate-spin" size={14} />
                Menyimpan...
              {:else}
                Save Changes
              {/if}
            </button>
          </div>
        </div>
      </form>
    </div>
  </div>
{/if}
