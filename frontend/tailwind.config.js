/** @type {import('tailwindcss').Config} */
export default {
    content: ['./src/**/*.{html,js,svelte,ts}'],
    darkMode: 'class',
    theme: {
        extend: {
            fontFamily: { sans: ['Inter', 'sans-serif'] },
            colors: {
                ubiquiti: {
                    50: '#f0f7ff',
                    100: '#e0effe',
                    400: '#3b82f6', // Standard Blue
                    500: '#2563eb', // Brand Blue
                    600: '#0057e7', // Ubiquiti-ish Blue
                    900: '#1e3a8a',
                }
            }
        }
    },
    plugins: []
};
