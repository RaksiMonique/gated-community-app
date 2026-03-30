<script>
  export let title = '';
  export let className = '';
</script>

<div class="rounded-xl border bg-white text-gray-950 shadow-sm {className}">
  {#if title}
    <div class="flex flex-col space-y-1.5 p-6 pb-2">
      <h3 class="font-semibold leading-none tracking-tight">{title}</h3>
    </div>
  {/if}
  <div class="p-6 pt-2">
    <slot />
  </div>
</div>