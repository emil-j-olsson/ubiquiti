# Ubiquiti Frontend

A modern web application built with SvelteKit that provides a monitoring dashboard for Ubiquiti network devices. The frontend communicates with the backend monitoring service to display real-time device diagnostics, health status, and network topology.

## Access

Access the web application by visiting `localhost:3000` using the Docker Compose setup. In development use `localhost:5173`.

## Technology Stack

- **Framework**: [SvelteKit](https://kit.svelte.dev/) v2.48.5 with Svelte 5
- **Styling**: [Tailwind CSS](https://tailwindcss.com/) v4.1.17 with custom glassmorphism design
- **Animations**: [GSAP](https://greensock.com/gsap/) v3.13.0

## Features

- Responsive glassmorphic UI with dark mode support.
- Real-time device monitoring dashboard.
- Server-side rendering (SSR) and client-side hydration.
- Optimized production builds with code splitting.

## Development

```bash
# install dependencies
npm install

# run development server
npm run dev
# or:
make dev

# build for production
npm run build

# preview production build
npm run preview
```
