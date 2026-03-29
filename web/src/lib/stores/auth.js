import { writable } from 'svelte/store';
import { browser } from '$app/environment';

// Initialize from local storage if available
const getStoredUser = () => {
    if (!browser) return null;
    try {
        const item = localStorage.getItem('user');
        return item ? JSON.parse(item) : null;
    } catch (e) {
        console.error('Failed to parse user from localStorage:', e);
        return null;
    }
};

const storedUser = getStoredUser();

export const user = writable(storedUser);
export const isLoading = writable(false);

export const login = async (email, password) => {
    isLoading.set(true);
    try {
        const response = await fetch('/api/v1/auth/login', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ email, password })
        });
        
        const result = await response.json();
        if (!result.success) throw new Error(result.message);

        const token = result.data.token;
        if (browser) localStorage.setItem('token', token);
        
        // Fetch user details
        const userRes = await fetch('/api/v1/auth/me', {
            headers: { 'Authorization': `Bearer ${token}` }
        });
        const userResult = await userRes.json();
        
        user.set(userResult.data);
        if (browser) localStorage.setItem('user', JSON.stringify(userResult.data));
    } finally {
        isLoading.set(false);
    }
};

export const logout = () => {
    user.set(null);
    if (browser) {
        localStorage.removeItem('user');
        localStorage.removeItem('token');
    }
};