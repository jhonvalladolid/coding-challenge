import request from 'supertest';
import app from '../../src/app.js';

const validPayload = {
  matrices: {
    q: [
      [1, 0],
      [0, 1],
    ],
    r: [
      [2, 3],
      [0, 4],
    ],
  },
};

function postStatistics(body) {
  return request(app).post('/api/v1/statistics').send(body);
}

describe('POST /api/v1/statistics', () => {
  it('returns statistics and diagonal flags for a valid payload', async () => {
    const response = await postStatistics(validPayload);

    expect(response.status).toBe(200);
    expect(response.headers['x-request-id']).toEqual(expect.any(String));
    expect(response.body.success).toBe(true);
    expect(response.body.data.statistics.max).toBe(4);
    expect(response.body.data.statistics.min).toBe(0);
    expect(response.body.data.statistics.sum).toBe(11);
    expect(response.body.data.statistics.average).toBeCloseTo(1.375);
    expect(response.body.data.diagonal).toEqual({
      q: true,
      r: false,
      anyDiagonal: true,
    });
  });

  it('rejects an empty body', async () => {
    const response = await postStatistics({});

    expect(response.status).toBe(400);
    expect(response.body.success).toBe(false);
    expect(response.body.error.code).toBe('VALIDATION_ERROR');
    expect(response.body.error.requestId).toEqual(expect.any(String));
    expect(response.body.error.details).toEqual(
      expect.arrayContaining([
        expect.objectContaining({ field: 'matrices' }),
      ]),
    );
  });

  it('rejects a payload without q', async () => {
    const response = await postStatistics({
      matrices: {
        r: [[1]],
      },
    });

    expect(response.status).toBe(400);
    expect(response.body.error.code).toBe('VALIDATION_ERROR');
    expect(response.body.error.details).toEqual(
      expect.arrayContaining([
        expect.objectContaining({ field: 'matrices.q' }),
      ]),
    );
  });

  it('rejects a payload without r', async () => {
    const response = await postStatistics({
      matrices: {
        q: [[1]],
      },
    });

    expect(response.status).toBe(400);
    expect(response.body.error.code).toBe('VALIDATION_ERROR');
    expect(response.body.error.details).toEqual(
      expect.arrayContaining([
        expect.objectContaining({ field: 'matrices.r' }),
      ]),
    );
  });

  it('rejects an empty matrix', async () => {
    const response = await postStatistics({
      matrices: {
        q: [],
        r: [[1]],
      },
    });

    expect(response.status).toBe(400);
    expect(response.body.error.code).toBe('VALIDATION_ERROR');
    expect(response.body.error.details).toEqual(
      expect.arrayContaining([
        expect.objectContaining({
          field: 'matrices.q',
          reason: 'matrices.q must not be empty',
        }),
      ]),
    );
  });

  it('rejects rows with different dimensions', async () => {
    const response = await postStatistics({
      matrices: {
        q: [
          [1, 2],
          [3],
        ],
        r: [[1]],
      },
    });

    expect(response.status).toBe(400);
    expect(response.body.error.code).toBe('VALIDATION_ERROR');
    expect(response.body.error.details).toEqual(
      expect.arrayContaining([
        expect.objectContaining({
          field: 'matrices.q',
          reason: 'matrices.q rows must have equal length',
        }),
      ]),
    );
  });

  it('rejects a non-numeric value', async () => {
    const response = await postStatistics({
      matrices: {
        q: [[1, 'x']],
        r: [[1]],
      },
    });

    expect(response.status).toBe(400);
    expect(response.body.error.code).toBe('VALIDATION_ERROR');
    expect(response.body.error.details).toEqual(
      expect.arrayContaining([
        expect.objectContaining({
          reason: 'matrices.q values must be finite numbers',
        }),
      ]),
    );
  });
});
