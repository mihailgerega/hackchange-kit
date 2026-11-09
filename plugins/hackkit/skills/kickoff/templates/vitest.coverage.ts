// Merge this `test` block into frontend/vite.config.ts (or vitest.config.ts).
// `coverage.include` is required: without it Vitest only counts files that tests import,
// which hides untested code and overstates coverage.
import { defineConfig } from 'vitest/config'

export default defineConfig({
  test: {
    environment: 'jsdom',
    setupFiles: ['./src/test/setup.ts'], // @testing-library/jest-dom + MSW server lifecycle
    coverage: {
      provider: 'v8',
      include: ['src/**/*.{ts,tsx}'],
      exclude: [
        'src/**/*.test.{ts,tsx}',
        'src/**/*.d.ts',
        'src/test/**',
        'src/**/generated/**',
        'src/main.tsx',
      ],
      thresholds: { lines: 80 },
    },
  },
})
