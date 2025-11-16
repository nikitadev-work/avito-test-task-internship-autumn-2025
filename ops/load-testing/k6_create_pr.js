import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
  scenarios: {
    create_pr: {
      executor: 'constant-arrival-rate',
      rate: 5,
      timeUnit: '1s',
      duration: '2m',
      preAllocatedVUs: 10,
      maxVUs: 20,
    },
  },
};

const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';
const ADMIN_TOKEN = __ENV.ADMIN_TOKEN || 'admin:u1';

export default function () {
  const url = `${BASE_URL}/pullRequest/create`;
  const payload = JSON.stringify({
    pull_request_id: `pr-${__VU}-${Date.now()}`,
    pull_request_name: 'load-test-pr',
    author_id: 'u1',
  });

  const params = {
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${ADMIN_TOKEN}`,
    },
  };

  const res = http.post(url, payload, params);

  check(res, {
    'status is 201 or 400': (r) => r.status === 201 || r.status === 400,
  });

  sleep(0.1);
}
