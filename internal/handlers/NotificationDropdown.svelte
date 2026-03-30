<script>
  import { onMount } from 'svelte';
  let notifications = [];
  let unreadCount = 0;
  let isOpen = false;

  const TEST_COMMUNITY_ID = '550e8400-e29b-41d4-a716-446655440000';

  async function loadNotifications() {
    try {
      const token = localStorage.getItem('token');
      if (!token) {
        console.log("No auth token found, skipping notification fetch.");
        return;
      }

      const headers = { 
        'X-Community-ID': TEST_COMMUNITY_ID,
        'Authorization': `Bearer ${token}`
      };

      const res = await fetch('/api/v1/notifications', { headers });
      const json = await res.json();
      if (json.success) notifications = json.data;
      
      const countRes = await fetch('/api/v1/notifications/unread-count', { headers });
      const countJson = await countRes.json();
      if (countJson.success) unreadCount = countJson.data;
    } catch (e) {
      console.error("Failed to load notifications:", e);
    }
  }

  async function handleNotificationClick(id, link) {
    await fetch(`/api/v1/notifications/${id}/read`, { method: 'PATCH' });
    await loadNotifications();
    if (link) {
      window.location.href = link;
    }
  }

  onMount(loadNotifications);
</script>

<div class="relative inline-block text-left">
  <button on:click={() => isOpen = !isOpen} class="relative p-2 text-gray-600 hover:bg-gray-100 rounded-full transition-colors">
    <svg class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9" />
    </svg>
    {#if unreadCount > 0}
      <span class="absolute top-0 right-0 block h-4 w-4 rounded-full bg-red-500 text-[10px] text-white font-bold flex items-center justify-center ring-2 ring-white">
        {unreadCount}
      </span>
    {/if}
  </button>

  {#if isOpen}
    <div class="absolute right-0 mt-2 w-80 bg-white rounded-lg shadow-xl border border-gray-200 z-50 overflow-hidden">
      <div class="px-4 py-2 border-b border-gray-100 bg-gray-50 flex justify-between items-center">
        <span class="text-sm font-bold text-gray-700">Notifications</span>
        <button on:click={() => isOpen = false} class="text-xs text-indigo-600 hover:underline">Close</button>
      </div>
      <div class="max-h-96 overflow-y-auto">
        {#each notifications as note}
          <button 
            type="button"
            class={`w-full text-left p-4 border-b border-gray-50 hover:bg-gray-50 transition-colors ${!note.is_read ? 'bg-indigo-50/30' : ''}`}
            on:click={() => handleNotificationClick(note.id, note.link_to_resource)}
          >
            <div class="flex justify-between items-start">
              <h4 class="text-sm font-semibold text-gray-900">{note.title}</h4>
              {#if !note.is_read}
                <span class="h-2 w-2 bg-indigo-600 rounded-full"></span>
              {/if}
            </div>
            <p class="text-xs text-gray-600 mt-1 line-clamp-2">{note.message}</p>
            <span class="text-[10px] text-gray-400 mt-2 block">
              {new Date(note.created_at).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}
            </span>
          </button>
        {:else}
          <div class="p-8 text-center text-gray-400 text-sm">
            No notifications yet
          </div>
        {/each}
      </div>
      <a href="/notifications" class="block py-2 text-center text-xs font-medium text-indigo-600 bg-gray-50 hover:bg-gray-100">
        View All Notifications
      </a>
    </div>
    
    <!-- Click away backdrop -->
    <button 
      type="button" 
      class="fixed inset-0 z-40 w-full h-full cursor-default bg-transparent" 
      on:click={() => isOpen = false} 
      aria-label="Close notifications"
    ></button>
  {/if}
</div>

<style>
  .line-clamp-2 {
    display: -webkit-box; -webkit-line-clamp: 2; -webkit-box-orient: vertical; overflow: hidden;
  }
</style>