<script>
  import { page } from '$app/stores';
  import { user, logout } from '$lib/stores/auth';
  
  let isSidebarOpen = false;

  const menuItems = [
    { label: 'Dashboard', href: '/dashboard', icon: 'LayoutDashboard' },
    { label: 'Residents', href: '/residents', icon: 'Users' },
    { label: 'Gate Access', href: '/access', icon: 'QrCode' },
    { label: 'Security Logs', href: '/logs', icon: 'ShieldAlert' },
  ];

  function toggleSidebar() {
    isSidebarOpen = !isSidebarOpen;
  }
</script>

<div class="min-h-screen bg-gray-50 flex">
  <!-- Mobile Sidebar Backdrop -->
  {#if isSidebarOpen}
    <div class="fixed inset-0 bg-gray-600 bg-opacity-75 z-20 lg:hidden" on:click={toggleSidebar}></div>
  {/if}

  <!-- Sidebar -->
  <div class="{isSidebarOpen ? 'translate-x-0' : '-translate-x-full'} lg:translate-x-0 fixed inset-y-0 left-0 z-30 w-64 bg-white border-r border-gray-200 transition-transform duration-300 ease-in-out lg:static lg:inset-auto lg:flex lg:flex-col">
    <div class="flex h-16 shrink-0 items-center px-6 border-b border-gray-200">
      <span class="text-xl font-bold text-blue-600">GateKeeper</span>
    </div>
    
    <nav class="flex-1 overflow-y-auto px-4 py-4 space-y-1">
      {#each menuItems as item}
        <a 
          href={item.href} 
          class="group flex items-center px-2 py-2 text-sm font-medium rounded-md {$page.url.pathname.startsWith(item.href) ? 'bg-blue-50 text-blue-700' : 'text-gray-700 hover:bg-gray-50 hover:text-gray-900'}"
        >
          {item.label}
        </a>
      {/each}
    </nav>

    <div class="border-t border-gray-200 p-4">
      <button on:click={logout} class="w-full flex items-center px-2 py-2 text-sm font-medium text-red-600 rounded-md hover:bg-red-50">
        Sign Out
      </button>
    </div>
  </div>

  <!-- Main Content -->
  <div class="flex-1 flex flex-col min-w-0 overflow-hidden">
    <!-- Top Header -->
    <header class="bg-white border-b border-gray-200 flex h-16 items-center justify-between px-4 sm:px-6 lg:px-8">
      <button class="lg:hidden text-gray-500 focus:outline-none" on:click={toggleSidebar}>
        <svg class="h-6 w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16"></path></svg>
      </button>
      
      <div class="flex items-center gap-4">
        <span class="text-sm text-gray-700">{$user?.name || 'User'}</span>
        <div class="h-8 w-8 rounded-full bg-blue-100 flex items-center justify-center text-blue-600 font-bold">
          {($user?.name || 'U')[0]}
        </div>
      </div>
    </header>

    <!-- Page Content -->
    <main class="flex-1 overflow-y-auto p-4 sm:p-6 lg:p-8">
      <slot />
    </main>
  </div>
</div>
```

### 5. Example Pages

**Login Page**
```diff
<script>
  import Input from '$lib/components/ui/Input.svelte';
  import Button from '$lib/components/ui/Button.svelte';
  import { login, isLoading } from '$lib/stores/auth';
  
  let email = '';
  let password = '';

  const handleLogin = async () => {
    await login(email, password);
    // Redirect handled by logic or protected route guard
  };
</script>

<div class="bg-white py-8 px-4 shadow sm:rounded-lg sm:px-10">
  <div class="mb-6 text-center">
    <h2 class="text-2xl font-bold text-gray-900">Sign in to your account</h2>
  </div>
  <form class="space-y-6" on:submit|preventDefault={handleLogin}>
    <Input label="Email address" type="email" bind:value={email} required />
    <Input label="Password" type="password" bind:value={password} required />
    <Button type="submit" className="w-full" disabled={$isLoading}>
      {$isLoading ? 'Signing in...' : 'Sign in'}
    </Button>
  </form>
</div>
```

**Dashboard Page**
```diff
<script>
  import Card from '$lib/components/ui/Card.svelte';
  import Badge from '$lib/components/ui/Badge.svelte';
</script>

<div class="mb-8">
  <h1 class="text-2xl font-bold text-gray-900">Dashboard</h1>
  <p class="mt-1 text-sm text-gray-500">Overview of your community status.</p>
</div>

<div class="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3">
  <Card title="Active Residents">
    <div class="text-3xl font-bold text-gray-900">1,240</div>
    <div class="mt-2 text-sm text-green-600 flex items-center">
      <span class="mr-1">↑</span> 12% from last month
    </div>
  </Card>

  <Card title="Recent Access Logs">
    <ul class="space-y-3 mt-2">
      <li class="flex justify-between items-center text-sm">
        <span>Visitor: John Smith</span>
        <Badge variant="success">Allowed</Badge>
      </li>
      <li class="flex justify-between items-center text-sm">
        <span>Visitor: Delivery Van</span>
        <Badge variant="warning">Pending</Badge>
      </li>
    </ul>
  </Card>
</div>
```

<!--
[PROMPT_SUGGESTION]Integrate the Svelte frontend Login page with the Go backend authentication endpoint[/PROMPT_SUGGESTION]
[PROMPT_SUGGESTION]Create the 'Visitor Access Control' module with QR code generation in Svelte[/PROMPT_SUGGESTION]
