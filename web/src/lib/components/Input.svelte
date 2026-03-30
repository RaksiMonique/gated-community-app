<script>
  export let label = '';
  export let value = '';
  export let type = 'text';
  export let placeholder = '';
  export let error = '';
  export let id = Math.random().toString(36).substr(2, 9);
</script>

<div class="w-full space-y-1">
  {#if label}
    <label for={id} class="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70">{label}</label>
  {/if}
  <input
    {id}
    {type}
    {placeholder}
    value={value}
    on:input={(e) => value = e.currentTarget.value}
    class="flex h-10 w-full rounded-md border border-gray-300 bg-white px-3 py-2 text-sm placeholder:text-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-600 focus:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50 {error ? 'border-red-500 focus:ring-red-500' : ''}"
    {...$$restProps}
  />
  {#if error}
    <p class="text-xs text-red-500 font-medium">{error}</p>
  {/if}
</div>