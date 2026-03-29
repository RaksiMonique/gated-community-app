<script>
  import { onMount } from 'svelte';
  
  let inviteCode = '';
  let validationResult = null; 
  let loading = false;
  let recentLogs = [];

  const TEST_COMMUNITY_ID = '550e8400-e29b-41d4-a716-446655440000';

  async function validatePass(actionType) {
    loading = true;
    try {
      const token = localStorage.getItem('token');
      const res = await fetch('/api/v1/visitors/validate', {
        method: 'POST',
        headers: { 
          'Content-Type': 'application/json',
          'X-Community-ID': TEST_COMMUNITY_ID,
          'Authorization': `Bearer ${token}`
        },
        body: JSON.stringify({
          invite_code: inviteCode.toUpperCase(),
          action_type: actionType
        })
      });
      const json = await res.json();
      validationResult = json;
      if (json.success) {
        inviteCode = '';
        fetchRecentLogs();
      }
    } catch (e) {
      validationResult = { success: false, message: 'Network error' };
    } finally {
      loading = false;
    }
  }

  async function fetchRecentLogs() {
    try {
      const token = localStorage.getItem('token');
      const res = await fetch('/api/v1/security/logs?limit=10', {
        headers: { 
          'X-Community-ID': TEST_COMMUNITY_ID,
          'Authorization': `Bearer ${token}`
        }
      });
      const json = await res.json();
      if (json.success) recentLogs = json.data;
    } catch (e) {
      console.error("Log fetch failed:", e);
    }
  }

  onMount(fetchRecentLogs);
</script>

<div class="min-h-screen bg-gradient-to-b from-slate-900 via-slate-900 to-black text-white pb-24 font-sans">
  <!-- Header -->
  <header class="p-5 bg-slate-800/50 backdrop-blur-md border-b border-white/5 flex justify-between items-center sticky top-0 z-40">
    <h1 class="text-xl font-black tracking-tighter bg-gradient-to-r from-indigo-400 to-cyan-400 bg-clip-text text-transparent">SECURITY GATE 1</h1>
    <div class="flex items-center gap-2">
      <span class="h-2 w-2 rounded-full bg-emerald-500 shadow-[0_0_8px_rgba(16,185,129,0.6)] animate-pulse"></span>
      <span class="text-[10px] font-black tracking-widest text-slate-400 uppercase">Live</span>
    </div>
  </header>

  <main class="p-4 space-y-6">
    <!-- Validation Console -->
    <section class="bg-slate-800/80 backdrop-blur-sm p-6 rounded-[2.5rem] shadow-2xl border border-white/5">
      <label for="code" class="block text-sm font-bold text-slate-400 mb-2 uppercase tracking-widest text-center">Enter Invite Code</label>
      <input 
        id="code"
        type="text" 
        bind:value={inviteCode}
        placeholder="CODE"
        class="w-full bg-slate-900 border-2 border-slate-700 rounded-2xl p-4 text-3xl font-mono text-center focus:border-indigo-500 focus:ring-0 transition-all uppercase mb-6 text-indigo-400 tracking-widest"
      />
      
      <div class="grid grid-cols-2 gap-4">
        <button 
          on:click={() => validatePass('ENTRY')}
          disabled={!inviteCode || loading}
          class="bg-indigo-600 hover:bg-indigo-500 disabled:opacity-50 py-6 rounded-3xl font-black text-xl shadow-lg shadow-indigo-500/20 active:scale-95 transition-all"
        >
          ENTRY
        </button>
        <button 
          on:click={() => validatePass('EXIT')}
          disabled={!inviteCode || loading}
          class="bg-slate-700 hover:bg-slate-600 disabled:opacity-50 py-6 rounded-3xl font-black text-xl shadow-lg active:scale-95 transition-all"
        >
          EXIT
        </button>
      </div>
    </section>

    <!-- Activity Feed -->
    <section>
      <div class="flex justify-between items-center mb-4 px-2">
        <h2 class="text-xs font-black text-slate-500 uppercase tracking-widest">Recent Activity</h2>
        <button on:click={fetchRecentLogs} class="text-xs font-bold text-indigo-400">Refresh</button>
      </div>
      <div class="space-y-3">
        {#each recentLogs as log}
          <div class="bg-slate-800/40 p-4 rounded-3xl flex justify-between items-center border border-white/5">
            <div class="flex items-center gap-4">
              <div class="h-10 w-10 rounded-2xl {log.time_out ? 'bg-slate-700' : 'bg-emerald-500/20'} flex items-center justify-center">
                <svg class="w-5 h-5 {log.time_out ? 'text-slate-400' : 'text-emerald-500'}" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d={log.time_out ? "M17 16l4-4m0 0l-4-4m4 4H7" : "M11 16l-4-4m0 0l4-4m-4 4h14"} />
                </svg>
              </div>
              <div>
                <p class="font-bold text-slate-200">{log.person_name}</p>
                <p class="text-[10px] text-slate-500 uppercase font-black tracking-tight">{log.unit_number} • {log.person_type}</p>
              </div>
            </div>
            <div class="text-right">
              <p class="text-xs font-mono font-bold text-indigo-400">
                {new Date(log.time_in).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}
              </p>
              <p class="text-[9px] font-black text-slate-600 uppercase mt-1">{log.time_out ? 'Departed' : 'Inside'}</p>
            </div>
          </div>
        {/each}
      </div>
    </section>
  </main>

  <!-- Full-screen Result Overlay -->
  {#if validationResult}
    <div class="fixed inset-0 z-50 flex items-center justify-center p-6 bg-slate-900/95 backdrop-blur-md">
      <div class="w-full max-w-sm rounded-[3rem] p-10 text-center shadow-2xl {validationResult.success ? 'bg-emerald-600 shadow-emerald-500/20' : 'bg-rose-600 shadow-rose-500/20'}">
        <div class="mb-6 scale-125">
          {#if validationResult.success}
            <svg class="w-16 h-16 mx-auto" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="3" d="M5 13l4 4L19 7" /></svg>
          {:else}
            <svg class="w-16 h-16 mx-auto" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="3" d="M6 18L18 6M6 6l12 12" /></svg>
          {/if}
        </div>
        
        <h2 class="text-4xl font-black mb-2 uppercase tracking-tighter italic">
          {validationResult.success ? 'Success' : 'Denied'}
        </h2>
        <p class="text-white/90 font-bold mb-10 leading-tight">
          {validationResult.message}
        </p>

        <button 
          on:click={() => validationResult = null}
          class="w-full bg-white text-black py-5 rounded-[2rem] font-black uppercase tracking-widest shadow-xl active:scale-95 transition-transform"
        >
          Next Visitor
        </button>
      </div>
    </div>
  {/if}

  <!-- Security Navigation -->
  <nav class="fixed bottom-0 left-0 right-0 bg-slate-900/80 backdrop-blur-xl border-t border-white/5 flex justify-around p-4 pb-10 z-40">
    <button class="text-indigo-400 flex flex-col items-center gap-1">
      <svg class="w-6 h-6" fill="currentColor" viewBox="0 0 20 20"><path d="M10 12a2 2 0 100-4 2 2 0 000 4z"/><path fill-rule="evenodd" d="M.458 10C1.732 5.943 5.522 3 10 3s8.268 2.943 9.542 7c-1.274 4.057-5.064 7-9.542 7S1.732 14.057.458 10zM14 10a4 4 0 11-8 0 4 4 0 018 0z" clip-rule="evenodd"/></svg>
      <span class="text-[9px] font-black uppercase tracking-widest">Scanner</span>
    </button>
    <button class="text-slate-500 flex flex-col items-center gap-1 opacity-50">
      <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"/></svg>
      <span class="text-[9px] font-black uppercase tracking-widest">Logs</span>
    </button>
  </nav>
</div>