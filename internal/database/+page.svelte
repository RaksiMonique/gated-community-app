<script>
  import { onMount } from 'svelte';
  let requests = [];
  let loading = true;
  const TEST_COMMUNITY_ID = '550e8400-e29b-41d4-a716-446655440000';

  async function fetchRequests() {
    try {
      const token = localStorage.getItem('token');
      const res = await fetch('/api/v1/maintenance', {
        headers: { 
          'Authorization': `Bearer ${token}`,
          'X-Community-ID': TEST_COMMUNITY_ID 
        }
      });
      const json = await res.json();
      if (json.success) requests = json.data;
    } finally {
      loading = false;
    }
  }

  onMount(fetchRequests);
</script>

<div class="space-y-6">
  <h1 class="text-2xl font-black text-slate-900 uppercase tracking-tight">Maintenance Tickets</h1>

  <div class="space-y-4">
    {#if loading}
      <div class="p-20 text-center text-slate-400 font-black animate-pulse">Retrieving Tickets...</div>
    {:else}
      {#each requests as req}
        <div class="card flex flex-col md:flex-row justify-between items-start md:items-center gap-4 border-l-4 {req.priority === 'URGENT' ? 'border-l-rose-500' : 'border-l-amber-500'}">
          <div class="flex-1">
            <div class="flex items-center gap-3 mb-1">
              <span class="badge {req.status === 'OPEN' ? 'badge-danger' : 'badge-warning'}">{req.status}</span>
              <span class="text-[10px] font-black text-slate-400 uppercase tracking-widest">{req.category}</span>
            </div>
            <h3 class="text-lg font-black text-slate-900">{req.title}</h3>
            <p class="text-sm text-slate-500 mt-1">{req.description}</p>
          </div>
          
          <div class="flex items-center gap-6 w-full md:w-auto border-t md:border-t-0 pt-4 md:pt-0">
            <div class="text-right">
              <p class="text-[10px] font-black text-slate-400 uppercase">Priority</p>
              <p class="text-sm font-bold {req.priority === 'URGENT' ? 'text-rose-600' : 'text-slate-900'}">{req.priority}</p>
            </div>
            <div class="text-right">
              <p class="text-[10px] font-black text-slate-400 uppercase">Created</p>
              <p class="text-sm font-bold text-slate-900">{new Date(req.created_at).toLocaleDateString()}</p>
            </div>
            <button class="btn btn-secondary text-xs">Assign Vendor</button>
          </div>
        </div>
      {:else}
        <div class="card text-center py-20 text-slate-400 font-bold uppercase tracking-widest">
          No active maintenance requests.
        </div>
      {/each}
    {/if}
  </div>
</div>