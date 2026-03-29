<script>
  export let variant = 'default'; // default, success, warning, danger, outline
  
  const variants = {
    default: "border-transparent bg-slate-900 text-white",
    success: "border-transparent bg-emerald-50 text-emerald-700 border border-emerald-200",
    warning: "border-transparent bg-amber-50 text-amber-700 border border-amber-200",
    danger: "border-transparent bg-rose-50 text-rose-700 border border-rose-200",
    outline: "text-slate-600 border-slate-200 bg-white"
  };
</script>

<div class="inline-flex items-center rounded-md px-2 py-0.5 text-[10px] font-bold uppercase tracking-wider transition-colors focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2 {variants[variant]}">
  <slot />
</div>