import http from 'k6/http';

export const options = {
  vus: 2000,
  duration: '30s',
};

export default function () {
  http.post('http://localhost:8080/pix', JSON.stringify({
    sender: "1",
    receiver: "2",
    amount: 10
  }), {
    headers: {
      'Content-Type': 'application/json',
      'Idempotency-Key': Math.random().toString()
    }
  });
}