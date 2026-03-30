<script>
    import { login, isLoading, user } from '$lib/stores/auth';
    import { goto } from '$app/navigation';
    import Button from '$lib/components/Button.svelte';
    import Input from '$lib/components/Input.svelte';

    let email = '';
    let password = '';
    let error = '';

    async function handleSubmit() {
        error = '';
        try {
            await login(email, password);
            goto('/');
        } catch (e) {
            error = e.message;
        }
    }

    $: if ($user) goto('/');
</script>

<div class="min-h-screen flex items-center justify-center bg-slate-50 px-4">
    <div class="max-w-md w-full space-y-8 bg-white p-8 rounded-2xl shadow-xl border border-slate-100">
        <div class="text-center">
            <div class="inline-flex items-center justify-center w-16 h-16 bg-blue-600 rounded-xl mb-4 shadow-lg shadow-blue-200">
                <span class="text-2xl font-bold text-white">G</span>
            </div>
            <h2 class="text-3xl font-extrabold text-slate-900">Welcome Back</h2>
            <p class="mt-2 text-sm text-slate-500">Sign in to manage your community</p>
        </div>

        <form class="mt-8 space-y-6" on:submit|preventDefault={handleSubmit}>
            <div class="space-y-4">
                <Input label="Email Address" type="email" bind:value={email} placeholder="you@example.com" required />
                <Input label="Password" type="password" bind:value={password} placeholder="••••••••" required />
            </div>

            {#if error}
                <div class="p-3 text-sm text-red-600 bg-red-50 border border-red-100 rounded-lg font-medium">
                    {error}
                </div>
            {/if}

            <Button type="submit" variant="primary" className="w-full h-12 text-base" disabled={$isLoading}>
                {$isLoading ? 'Signing in...' : 'Sign In'}
            </Button>
        </form>
    </div>
</div>