import { defineConfig, loadEnv, createLogger } from "vite";
import { gsx, devFallback } from "@gsxhq/vite-plugin-gsx";

export default defineConfig(({ command, mode }) => {
  const env = loadEnv(mode, process.cwd(), "");
  const goPort = env.GO_PORT || "7777";
  const vitePort = parseInt(env.VITE_PORT || "5173", 10);
  // Single-computed upstream: `gsx dev` injects GSX_DEV_UPSTREAM into the real
  // process env, injected by gsx dev at spawn; everything else
  // (standalone `vite`, `vite build`) falls back to GO_PORT. Passed explicitly
  // to devFallback AND the proxy below so evaluation order can't matter — a
  // no-args devFallback call throws when neither is set, which would break
  // `vite build` (devFallback is inert there, but the zero-arg constructor
  // still throws) and standalone `vite serve`.
  const upstream = process.env.GSX_DEV_UPSTREAM ?? `http://localhost:${goPort}`;
  // Serve a self-recovering interstitial while the Go server is down/restarting
  // (instead of a raw proxy error).
  const fallback = devFallback({ target: upstream });

  // While the Go server is down/restarting, the dev-fallback interstitial already
  // shows it — so drop Vite's redundant "http proxy error … ECONNREFUSED" spam.
  const logger = createLogger();
  const baseError = logger.error;
  logger.error = (msg, opts) => {
    if (typeof msg === "string" && msg.includes("http proxy error")) return;
    baseError(msg, opts);
  };

  return {
    clearScreen: false,
    // Dev serves all Vite assets under /__vite/ (matches gsxhq/vite DevBase);
    // prod uses the default base ("/") since hashed assets are served from
    // /static via gsxhq/vite.
    base: command === "serve" ? "/__vite/" : "/",
    publicDir: false,
    customLogger: logger,
    plugins: [
      gsx(),
      fallback.plugin,
      {
        // Log the app front door ("/"); Vite's own line shows the empty /__vite/ base.
        name: "gsx-dev-url",
        configureServer(server) {
          server.printUrls = () => {
            server.config.logger.info(
              `\n  \x1b[32m➜\x1b[0m  Open \x1b[36mhttp://localhost:${vitePort}/\x1b[0m to view your app\n`,
            );
          };
        },
      },
    ],
    server: {
      port: vitePort,
      // `gsx dev` chooses an available VITE_PORT and matching VITE_DEV_URL before
      // starting Vite. Keep Vite strict so proxy, overlay, and printed URL agree.
      strictPort: true,
      proxy: {
        // Everything except /__vite/ (Vite's own dev-asset namespace) and
        // /__dev/ (fallback plugin status) goes to the Go server.
        // Single exclusion — no source-dir denylist, no app-route collisions.
        // No `ws: true` — the Go server has no WebSocket; proxying ws would
        // capture Vite's HMR socket.
        "^(?!/__vite/|/__dev).*": {
          target: upstream,
          changeOrigin: true,
          configure: fallback.configureProxy,
        },
      },
    },
    build: {
      manifest: true,
      outDir: "dist",
      rollupOptions: { input: "web/main.js" },
    },
  };
});
