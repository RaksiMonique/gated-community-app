import { defineConfig } from 'vite';
import { svelte } from '@sveltejs/vite-plugin-svelte';
import path from 'path';
import { fileURLToPath } from 'url';

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

export default defineConfig({
  plugins: [svelte()],
  server: {
    host: true,
    port: 5173,
    proxy: {
      '/api': 'http://127.0.0.1:8080'
    }
  },
  resolve: {
    alias: {
      $lib: path.resolve(__dirname, './web/src/lib'),
      $internal: path.resolve(__dirname, './internal/routes'),
      $handlers: path.resolve(__dirname, './internal/handlers'),
    },
  },
});