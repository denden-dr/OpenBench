<script lang="ts">
  interface Props {
    value: string;
    label: string;
    type?: string;
    id: string;
    placeholder?: string;
    error?: string;
    required?: boolean;
    disabled?: boolean;
  }

  let {
    value = $bindable(),
    label,
    type = 'text',
    id,
    placeholder = '',
    error = '',
    required = false,
    disabled = false
  }: Props = $props();
</script>

<div class="flex flex-col gap-2 w-full">
  <label for={id} class="font-display font-bold text-sm text-neubrutalism-charcoal uppercase tracking-wider">
    {label} {required ? '*' : ''}
  </label>

  <input
    {type}
    {id}
    {placeholder}
    {required}
    {disabled}
    bind:value={value}
    aria-invalid={error ? 'true' : 'false'}
    aria-describedby={error ? `${id}-error` : undefined}
    class="w-full border-4 border-neubrutalism-charcoal bg-white p-3 font-sans text-neubrutalism-charcoal rounded-none transition-all duration-150 focus:outline-none focus:ring-4 focus:ring-neubrutalism-charcoal focus:bg-[#fefefe]
    {error ? 'border-neubrutalism-pink bg-rose-50' : ''}
    {disabled ? 'opacity-50 cursor-not-allowed bg-zinc-100 border-dashed' : ''}"
  />

  {#if error}
    <p id={`${id}-error`} class="font-sans text-sm font-bold text-neubrutalism-pink mt-1" role="alert">
      {error}
    </p>
  {/if}
</div>
