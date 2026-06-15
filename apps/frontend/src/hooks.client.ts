import { initMockNetwork } from '$lib/services/mocks/network';
import { isMockEnabled } from '$lib/services/auth';

if (isMockEnabled()) {
  initMockNetwork();
}

