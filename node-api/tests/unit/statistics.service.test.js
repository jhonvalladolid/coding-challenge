import { computeStatistics } from '../../src/modules/statistics/statistics.service.js';

describe('statistics.service', () => {
  it('computes statistics over Q and R together', () => {
    const q = [
      [1, 0],
      [0, 1],
    ];
    const r = [
      [2, 3],
      [0, 4],
    ];

    const result = computeStatistics(q, r);

    expect(result.statistics.max).toBe(4);
    expect(result.statistics.min).toBe(0);
    expect(result.statistics.sum).toBe(11);
    expect(result.statistics.average).toBeCloseTo(1.375);
    expect(result.diagonal).toEqual({
      q: true,
      r: false,
      anyDiagonal: true,
    });
  });

  it('sets anyDiagonal when only R is diagonal', () => {
    const q = [
      [1, 2],
      [3, 4],
    ];
    const r = [
      [-1, 0],
      [0, -2],
    ];

    const result = computeStatistics(q, r);

    expect(result.diagonal.q).toBe(false);
    expect(result.diagonal.r).toBe(true);
    expect(result.diagonal.anyDiagonal).toBe(true);
  });
});
