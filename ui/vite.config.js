import { defineConfig } from "vite";
import { svelte } from "@sveltejs/vite-plugin-svelte";

export default defineConfig({
  plugins: [svelte()],
  define: {
    __API_URL__: JSON.stringify(
      process.env.MOOGIE_API_URL || "http://localhost:8080"
    ),
  },
  server: {
    host: true,
    port: 3000,
  },
});
