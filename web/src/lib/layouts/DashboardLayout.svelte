<script>
  import { onMount } from 'svelte';
  import { user, logout } from '$lib/stores/auth';

  let { children } = $props();
  let isMobileMenuOpen = $state(false);
  let currentPath = $state(typeof window !== 'undefined' ? window.location.pathname : '/');

  const navItems = [
    { name: 'Dashboard', href: '/dashboard', icon: 'home' },
    { name: 'Units', href: '/units', icon: 'building' },
    { name: 'Visitors', href: '/visitors', icon: 'users' },
    { name: 'Maintenance', href: '/maintenance', icon: 'tool' },
  ];

  onMount(() => {
    const handlePopState = () => {
      currentPath = window.location.pathname;
    };
    window.addEventListener('popstate', handlePopState);
    return () => window.removeEventListener('popstate', handlePopState);
  });

  function navigate(e, href) {
    if (typeof window !== 'undefined') {
      e.preventDefault();
      window.history.pushState({}, '', href);
      currentPath = href;
    }
  }
</script>

<div class="min-h-screen bg-slate-50 flex flex-col lg:flex-row">
  <!-- Sidebar (Desktop) -->
  <aside class="hidden lg:flex flex-col w-64 bg-slate-900 text-slate-300 border-r border-slate-800 space-y-8 fixed h-full z-50">
    <div class="p-6">
      <div class="text-2xl font-bold text-white tracking-tight flex items-center gap-2">
        <div class="w-8 h-8 bg-indigo-500 rounded-lg flex items-center justify-center font-black text-xl italic text-white">G</div>
        GatedAPI
      </div>
    </div>
    
    <nav class="flex-1 px-4 space-y-1">
      <p class="text-[10px] font-black text-slate-500 uppercase tracking-widest px-2 mb-4">Main Menu</p>
      {#each navItems as item}
        <a 
          href={item.href}
          onclick={(e) => navigate(e, item.href)}
          class="flex items-center gap-3 px-4 py-3 rounded-xl transition-all font-bold text-sm {currentPath === item.href ? 'bg-slate-800 text-white' : 'hover:bg-slate-800 hover:text-white'}"
        >
          <span class="font-medium">{item.name}</span>
        </a>
      {/each}
    </nav>

    <div class="p-4 mt-auto border-t border-slate-800">
      <button 
        onclick={logout} 
        class="w-full flex items-center gap-3 px-4 py-3 text-slate-400 font-bold text-sm hover:text-white transition-colors">
        <svg class="w-5 h-5 text-rose-500" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" /></svg>
        <span>Logout</span>
      </button>
    </div>
  </aside>

  <!-- Mobile Header -->
  <header class="lg:hidden bg-white border-b p-4 flex justify-between items-center sticky top-0 z-10">
    <div class="flex items-center gap-2">
      <div class="w-6 h-6 bg-indigo-500 rounded flex items-center justify-center font-black text-xs italic text-white">G</div>
      <span class="font-black tracking-tighter text-slate-900 uppercase text-sm">GatedAPI</span>
    </div>
    <button onclick={() => isMobileMenuOpen = !isMobileMenuOpen} class="p-2 text-slate-600">
      <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16"></path></svg>
    </button>
  </header>

  <!-- Mobile Nav Overlay -->
  {#if isMobileMenuOpen}
    <div class="fixed inset-0 bg-slate-900/50 z-20 lg:hidden" onclick={() => isMobileMenuOpen = false}>
      <div class="bg-white w-3/4 h-full p-6 shadow-xl" onclick={(e) => e.stopPropagation()}>
        <nav class="flex flex-col gap-4">
          {#each navItems as item}
            <a 
              href={item.href} 
              class="text-lg font-medium text-slate-800" 
              onclick={(e) => {
                navigate(e, item.href);
                isMobileMenuOpen = false;
              }}
            >{item.name}</a>
          {/each}
        </nav>
      </div>
    </div>
  {/if}

  <div class="flex-1 flex flex-col min-w-0 lg:ml-64">
    <!-- Top Bar -->
    <header class="hidden lg:flex h-20 bg-white border-b border-slate-100 items-center justify-between px-8 sticky top-0 z-40">
      <div class="flex items-center gap-2">
        <span class="text-[10px] font-black text-slate-400 uppercase tracking-widest">System / </span>
        <span class="text-[10px] font-black text-indigo-600 uppercase tracking-widest">{currentPath.split('/')[1] || 'Overview'}</span>
      </div>
      <div class="flex items-center gap-6">
        <div class="flex items-center gap-3 border-l border-slate-100 pl-6">
          <div class="text-right">
            <p class="text-sm font-black text-slate-900 leading-none">{$user?.name || 'Admin User'}</p>
            <p class="text-[10px] font-bold text-indigo-600 uppercase tracking-widest mt-1">Administrator</p>
          </div>
          <div class="w-10 h-10 rounded-xl bg-slate-100 border border-slate-200 flex items-center justify-center text-xs font-black text-slate-400">
            {$user?.name?.charAt(0) || 'A'}
          </div>
        </div>
      </div>
    </header>

    <!-- Main Content -->
    <main class="flex-1 p-4 md:p-8 overflow-y-auto">
      <!-- The dashboard content injected here will now inherit the card and button styles -->
      {@render children()}
    </main>
  </div>
</div>