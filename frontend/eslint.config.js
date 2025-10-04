import antfu from '@antfu/eslint-config'
import { tanstackConfig } from '@tanstack/eslint-config'

export default antfu(
  {
    formatters: true,
    react: true,
  },
  {
    files: ['**/*.{ts,tsx}'],
    ...tanstackConfig[0],
  },
  {
    ignores: ['src/routeTree.gen.ts'],
  },
)
