import preact from '@preact/preset-vite'
import { defineConfig } from 'vite'

// https://vite.dev/config/
export default defineConfig({
  plugins: [preact()],
  base: "./",
  build: {
    outDir: "../bonefire/webview/output",
    assetsDir: "assets",
    sourcemap: false,
    minify: true,
    rollupOptions: {
      input: {
        main: "./index.html",
      },
    },
  },
})
