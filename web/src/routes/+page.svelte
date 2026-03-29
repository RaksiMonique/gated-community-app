<script>
  import Button from '$lib/components/Button.svelte';
  import Input from '$lib/components/Input.svelte';
  import Badge from '$lib/components/Badge.svelte';
  import DashboardLayout from '$lib/layouts/DashboardLayout.svelte';
</script>

<DashboardLayout>
  <div class="space-y-6">
    <header class="flex flex-col md:flex-row md:items-center justify-between gap-4">
      <div>
        <h1 class="text-2xl font-bold text-slate-900">Community Overview</h1>
        <p class="text-slate-500">Welcome back! Here is what's happening today.</p>
      </div>
      <div class="flex items-center gap-2">
        <Button variant="outline" size="sm">Download Report</Button>
        <Button variant="primary" size="sm">Invite Visitor</Button>
      </div>
    </header>

    <!-- Stats Grid -->
    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
      {#each [
        { label: 'Total Units', value: '124', change: '+2', color: 'primary' },
        { label: 'Active Residents', value: '450', change: '+12', color: 'success' },
        { label: 'Expected Visitors', value: '18', change: '-3', color: 'warning' },
        { label: 'Pending Requests', value: '5', change: '0', color: 'danger' }
      ] as stat}
        <div class="bg-white p-6 rounded-xl border border-slate-200 shadow-sm">
          <p class="text-sm font-medium text-slate-500">{stat.label}</p>
          <div class="flex items-end justify-between mt-2">
            <h3 class="text-3xl font-bold text-slate-900">{stat.value}</h3>
            <Badge variant={stat.color === 'success' ? 'success' : stat.color === 'danger' ? 'danger' : 'default'}>
              {stat.change}
            </Badge>
          </div>
        </div>
      {/each}
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <!-- Recent Visitors Table -->
      <div class="lg:col-span-2 bg-white rounded-xl border border-slate-200 shadow-sm overflow-hidden">
        <div class="p-4 border-b border-slate-100 flex items-center justify-between">
          <h2 class="font-bold text-slate-900">Recent Visitors</h2>
          <a href="/visitors" class="text-sm text-blue-600 hover:underline">View all</a>
        </div>
        <div class="overflow-x-auto">
          <table class="w-full text-left text-sm">
            <thead class="bg-slate-50 text-slate-500 uppercase text-xs">
              <tr>
                <th class="px-4 py-3 font-semibold">Visitor</th>
                <th class="px-4 py-3 font-semibold">Unit</th>
                <th class="px-4 py-3 font-semibold">Status</th>
                <th class="px-4 py-3 font-semibold text-right">Time</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-100">
              {#each [
                { name: 'Sarah Connor', unit: 'B-204', status: 'In', time: '10:45 AM' },
                { name: 'John Wick', unit: 'A-101', status: 'Expected', time: '12:00 PM' },
                { name: 'Ellen Ripley', unit: 'C-305', status: 'Out', time: '09:30 AM' }
              ] as visitor}
                <tr class="hover:bg-slate-50 transition-colors">
                  <td class="px-4 py-3 font-medium text-slate-900">{visitor.name}</td>
                  <td class="px-4 py-3 text-slate-600">{visitor.unit}</td>
                  <td class="px-4 py-3">
                    <Badge variant={visitor.status === 'In' ? 'success' : visitor.status === 'Expected' ? 'warning' : 'outline'}>
                      {visitor.status}
                    </Badge>
                  </td>
                  <td class="px-4 py-3 text-right text-slate-500">{visitor.time}</td>
                </tr>
              {/each}
            </tbody>
          </table>
        </div>
      </div>

      <!-- Quick Actions / Search -->
      <div class="bg-white p-6 rounded-xl border border-slate-200 shadow-sm space-y-4">
        <h2 class="font-bold text-slate-900">Quick Search</h2>
        <Input placeholder="Search unit or resident..." />
        <div class="pt-4 border-t border-slate-100">
          <Button variant="outline" className="w-full justify-start">Log Maintenance Issue</Button>
        </div>
      </div>
    </div>
  </div>
</DashboardLayout>