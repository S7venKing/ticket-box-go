# Ticket Box Web

React/Vite frontend for the Ticket Box Go API gateway.

## Run

```bash
npm install
npm run dev
```

Set `VITE_API_URL` when the gateway is not running at the default
`http://localhost:8081`.

## Structure

- `src/components`: reusable UI and layout components
- `src/hooks`: reusable Redux and async hooks
- `src/modules`: Redux slices grouped by business module
- `src/pages`: route-level screens
- `src/lib`: API and infrastructure helpers
