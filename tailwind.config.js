/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./**/*.svelte",
    "./**/*.js",
    "./**/*.ts",
    "!./node_modules/**",
  ],
  theme: {
    extend: {
      colors: {
        primary: {
          600: '#4f46e5',
          700: '#4338ca',
        }
      }
    },
  },
  plugins: [require('@tailwindcss/forms')],
}