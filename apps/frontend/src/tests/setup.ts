import { vi } from 'vitest';
import '@testing-library/jest-dom';

// Mock SvelteKit runtime modules
vi.mock('$app/navigation', () => ({
  goto: vi.fn(),
  invalidate: vi.fn(),
  prefetch: vi.fn(),
  prefetchRoutes: vi.fn()
}));

vi.mock('$app/environment', () => ({
  browser: true,
  dev: true,
  building: false,
  version: 'mock-version'
}));

vi.mock('$app/paths', () => ({
  base: '',
  assets: ''
}));

vi.mock('$env/dynamic/public', () => ({
  env: {
    PUBLIC_API_URL: '',
    PUBLIC_MOCK_API: 'false'
  }
}));
