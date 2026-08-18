import request from 'supertest';
import app from '../../src/app.js';

describe('GET /health', () => {
  it('returns the service liveness payload', async () => {
    const response = await request(app).get('/health');

    expect(response.status).toBe(200);
    expect(response.headers['x-request-id']).toEqual(expect.any(String));
    expect(response.body).toEqual({
      success: true,
      data: {
        status: 'ok',
        service: 'node-api',
      },
    });
  });

  it('reuses the incoming X-Request-ID', async () => {
    const response = await request(app)
      .get('/health')
      .set('X-Request-ID', 'health-trace-1');

    expect(response.status).toBe(200);
    expect(response.headers['x-request-id']).toBe('health-trace-1');
  });
});
