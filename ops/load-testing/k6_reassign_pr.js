import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
  scenarios: {
    reassign_pr: {
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
const PR_ID = __ENV.PR_ID || 'pr-load-1';
const OLD_USER_ID = __ENV.OLD_USER_ID || 'u2';

export default function () {
  const url = `${BASE_URL}/pullRequest/reassign`;
  const payload = JSON.stringify({
    pull_request_id: PR_ID,
    old_user_id: OLD_USER_ID,
  });

  const params = {
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${ADMIN_TOKEN}`,
    },
  };

  const res = http.post(url, payload, params);

  check(res, {
    'status is 200 or 400': (r) => r.status === 200 || r.status === 400,
  });

  sleep(0.1);
}
