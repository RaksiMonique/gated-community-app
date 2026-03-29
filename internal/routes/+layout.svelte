<script>
  import "./app.css";
  import NotificationDropdown from "../handlers/NotificationDropdown.svelte";

  // Mocking SvelteKit's $page store for a standard SPA setup
  let currentPath = window.location.pathname;

  let user = { role: 'ADMIN', name: 'Admin User' };

  const adminItems = [
    { name: 'Dashboard', href: '/admin/dashboard', icon: 'M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6' },
    { name: 'Units', href: '/units', icon: 'M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4' },
    { name: 'Visitors', href: '/visitors', icon: 'M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z' },
    { name: 'Maintenance', href: '/maintenance', icon: 'M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z' },
  ];

  const residentItems = [
    { name: 'My Home', href: '/resident/dashboard', icon: 'M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6' },
    { name: 'Visitors', href: '/resident/visitors', icon: 'M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z' },
    { name: 'Payments', href: '/resident/payments', icon: 'M3 10h18M7 15h1m4 0h1m-7 4h12a3 3 0 003-3V8a3 3 0 00-3-3H6a3 3 0 00-3 3v8a3 3 0 003 3z' },
    { name: 'Maintenance', href: '/resident/maintenance', icon: 'M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z' },
  ];

  const securityItems = [
    { name: 'Scanner', href: '/', icon: 'M10 12a2 2 0 100-4 2 2 0 000 4z M.458 10C1.732 5.943 5.522 3 10 3s8.268 2.943 9.542 7c-1.274 4.057-5.064 7-9.542 7S1.732 14.057.458 10z' },
    { name: 'Logs', href: '/security/history', icon: 'M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z' },
  ];

  $: navItems = user.role === 'ADMIN' 
    ? adminItems 
    : (user.role === 'SECURITY' ? securityItems : residentItems);
</script>

<div class="flex min-h-screen">
  <!-- Sidebar -->
  <aside class="w-64 bg-slate-900 text-white flex flex-col fixed h-full z-50">
    <div class="p-6 flex items-center gap-3">
      <div class="h-8 w-8 bg-indigo-500 rounded-lg flex items-center justify-center font-black text-xl italic">G</div>
      <span class="text-xl font-black tracking-tighter">GatedAPI</span>
    </div>

    <nav class="flex-1 px-4 space-y-1">
      <p class="text-[10px] font-black text-slate-500 uppercase tracking-widest px-2 mb-4">Main Menu</p>
      {#each navItems as item}
        <a href={item.href} class="flex items-center px-4 py-3 rounded-xl transition-all font-bold text-sm { currentPath === item.href ? 'bg-indigo-600 text-white' : 'text-slate-400 hover:bg-slate-800 hover:text-white' }">
          <svg class="w-5 h-5 mr-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d={item.icon} />
          </svg>
          {item.name}
        </a>
      {/each}
    </nav>

    <div class="p-4 mt-auto border-t border-slate-800">
      <button class="w-full flex items-center px-4 py-3 text-slate-400 font-bold text-sm hover:text-white">
        <svg class="w-5 h-5 mr-3 text-rose-500" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" /></svg>
        Logout
      </button>
    </div>
  </aside>

  <!-- Main Content Area -->
  <main class="flex-1 ml-64">
    <header class="h-20 bg-white border-b border-slate-100 flex items-center justify-between px-8 sticky top-0 z-40">
      <div>
        <span class="text-[10px] font-black text-slate-400 uppercase tracking-widest">System / {currentPath.split('/').pop() || 'Overview'}</span>
      </div>
      <div class="flex items-center gap-6">
        <NotificationDropdown />
        <div class="flex items-center gap-3 lg:border-l lg:border-slate-100 lg:pl-6">
          <div class="text-right">
            <p class="text-sm font-black text-slate-900 leading-none">{user.name}</p>
            <p class="text-[10px] font-bold text-indigo-600 uppercase tracking-widest mt-1">{user.role}</p>
          </div>
          <div class="h-10 w-10 rounded-xl bg-slate-100 flex items-center justify-center font-black text-slate-400 border border-slate-200">
            U
          </div>
        </div>
      </div>
    </header>

    <div class="p-8">
      <slot />
    </div>
  </main>
</div>