import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
  scenarios: {
    get_reviews: {
      executor: 'constant-arrival-rate',
      rate: 5,
      timeUnit: '1s',
      duration: '2m',
      preAllocatedVUs: 5,
      maxVUs: 10,
    },
  },
};

const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';
const USER_TOKEN = __ENV.USER_TOKEN || 'user:u2';
const USER_ID = __ENV.USER_ID || 'u2';

export default function () {
  const url = `${BASE_URL}/users/getReview?user_id=${USER_ID}`;

  const params = {
    headers: {
      'Authorization': `Bearer ${USER_TOKEN}`,
    },
  };

  const res = http.get(url, params);

  check(res, {
    'status is 200': (r) => r.status === 200,
  });

  sleep(0.1);
}
