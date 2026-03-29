import http from 'k6/http';
import { uuidv4 } from 'https://jslib.k6.io/k6-utils/1.4.0/index.js';
import { check } from 'k6';

export const options = {
  stages: [
    { duration: '30s', target: 200 },   // aquecimento
    { duration: '1m',  target: 1000 },  // carga normal
    { duration: '2m',  target: 2000 },  // pico
    { duration: '30s', target: 0 },     // cooldown
  ],
  thresholds: {
    http_req_failed:   ['rate<0.01'],
    http_req_duration: ['p(95)<800'],
    http_req_duration: ['p(99)<2000'],
  },
};

export default function () {
  const res = http.post(
    'http://localhost:8080/pix',
    JSON.stringify({
      sender:   '1',
      receiver: '2',
      amount:   0.01,
    }),
    {
      headers: {
        'Content-Type':    'application/json',
        'Idempotency-Key': uuidv4(),
      },
      timeout: '10s',
    }
  );

  check(res, {
    'status 200 ou 201': (r) => r.status === 200 || r.status === 201,
    'tempo < 800ms':     (r) => r.timings.duration < 800,
  });
}