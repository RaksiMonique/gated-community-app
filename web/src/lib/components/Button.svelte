<script>
  export let variant = 'primary'; // primary, secondary, outline, ghost, danger
  export let size = 'md'; // sm, md, lg
  export let type = 'button';
  export let className = '';
  export let disabled = false;

  const base = "inline-flex items-center justify-center rounded-lg font-medium transition-all duration-200 active:scale-[0.98] focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-offset-2 disabled:opacity-50 disabled:pointer-events-none ring-offset-white shadow-sm";
  
  const variants = {
    primary: "bg-blue-600 text-white hover:bg-blue-700 focus-visible:ring-blue-600",
    secondary: "bg-gray-100 text-gray-900 hover:bg-gray-200 focus-visible:ring-gray-500",
    outline: "border border-gray-300 bg-transparent hover:bg-gray-100 focus-visible:ring-gray-500",
    ghost: "hover:bg-gray-100 hover:text-gray-900 focus-visible:ring-gray-500",
    danger: "bg-red-600 text-white hover:bg-red-700 focus-visible:ring-red-600"
  };
  
  const sizes = {
    sm: "h-8 px-3 text-xs",
    md: "h-10 py-2 px-4 text-sm",
    lg: "h-12 px-8 text-base"
  };
</script>

<button {type} {disabled} class="{base} {variants[variant]} {sizes[size]} {className}" on:click {...$$restProps}>
  <slot />
</button>