<script lang="ts">
  // Responsive Skeleton Loader for Tickets
</script>

<div class="w-full space-y-4 font-sans animate-pulse">
  <!-- Desktop Layout (Hidden on Mobile) -->
  <div class="hidden md:block bg-white dark:bg-slate-950 border border-slate-200 dark:border-slate-800 rounded-2xl overflow-hidden shadow-sm">
    <table class="w-full text-left border-collapse">
      <thead>
        <tr class="bg-slate-50 dark:bg-slate-900/50 border-b border-slate-200 dark:border-slate-800 text-xs font-bold uppercase tracking-wider">
          <th class="py-4 px-6 text-slate-400 dark:text-slate-600">Device</th>
          <th class="py-4 px-6 text-slate-400 dark:text-slate-600">Customer</th>
          <th class="py-4 px-6 text-slate-400 dark:text-slate-600">Kerusakan</th>
          <th class="py-4 px-6 text-slate-400 dark:text-slate-600">Price</th>
          <th class="py-4 px-6 text-slate-400 dark:text-slate-600">Status</th>
          <th class="py-4 px-6 text-right text-slate-400 dark:text-slate-600">Quick Action</th>
        </tr>
      </thead>
      <tbody class="divide-y divide-slate-200 dark:divide-slate-800">
        {#each Array(5) as _, i}
          <tr class="transition-colors">
            <!-- Device Column -->
            <td class="py-4 px-6 space-y-2">
              <div class="h-4 bg-slate-200 dark:bg-slate-800 rounded w-24"></div>
              <div class="h-3 bg-slate-100 dark:bg-slate-900 rounded w-16"></div>
            </td>
            <!-- Customer Column -->
            <td class="py-4 px-6 space-y-2">
              <div class="h-4 bg-slate-200 dark:bg-slate-800 rounded w-28"></div>
              <div class="h-3 bg-slate-100 dark:bg-slate-900 rounded w-12"></div>
            </td>
            <!-- Kerusakan Column -->
            <td class="py-4 px-6">
              <div class="h-4 bg-slate-200 dark:bg-slate-800 rounded w-36"></div>
            </td>
            <!-- Price Column -->
            <td class="py-4 px-6">
              <div class="h-4 bg-slate-200 dark:bg-slate-800 rounded w-20"></div>
            </td>
            <!-- Status Column -->
            <td class="py-4 px-6">
              <div class="h-5 bg-slate-200 dark:bg-slate-800 rounded-full w-20"></div>
            </td>
            <!-- Action Column -->
            <td class="py-4 px-6 text-right">
              <div class="h-8 bg-slate-200 dark:bg-slate-800 rounded-lg w-28 inline-block"></div>
            </td>
          </tr>
        {/each}
      </tbody>
    </table>
  </div>

  <!-- Mobile Layout (Hidden on Desktop) -->
  <div class="grid grid-cols-1 gap-4 md:hidden">
    {#each Array(3) as _}
      <div class="bg-white dark:bg-slate-950 border border-slate-200 dark:border-slate-800 rounded-2xl p-5 shadow-sm space-y-4">
        <div class="flex justify-between items-start">
          <div class="space-y-2">
            <div class="h-4 bg-slate-200 dark:bg-slate-800 rounded w-24"></div>
            <div class="h-3 bg-slate-100 dark:bg-slate-900 rounded w-32"></div>
          </div>
          <div class="h-5 bg-slate-200 dark:bg-slate-800 rounded-full w-16"></div>
        </div>
        <div class="border-t border-slate-100 dark:border-slate-900 pt-3 flex justify-between items-center">
          <div class="h-4 bg-slate-200 dark:bg-slate-800 rounded w-20"></div>
          <div class="h-8 bg-slate-200 dark:bg-slate-800 rounded-lg w-28"></div>
        </div>
      </div>
    {/each}
  </div>
</div>
