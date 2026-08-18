import {
  average,
  flattenMatrices,
  flattenMatrix,
  isDiagonal,
  max,
  min,
  sum,
} from '../../src/utils/matrix.util.js';

describe('matrix.util', () => {
  describe('flatten', () => {
    it('flattens a single matrix in row-major order', () => {
      expect(
        flattenMatrix([
          [1, 2],
          [3, 4],
        ]),
      ).toEqual([1, 2, 3, 4]);
    });

    it('flattens multiple matrices into one list of values', () => {
      const q = [
        [1, 0],
        [0, 1],
      ];
      const r = [
        [2, 3],
        [0, 4],
      ];

      expect(flattenMatrices(q, r)).toEqual([1, 0, 0, 1, 2, 3, 0, 4]);
    });
  });

  describe('max', () => {
    it('returns the largest value', () => {
      expect(max([1, 8, 3, 0])).toBe(8);
    });

    it('supports negative numbers', () => {
      expect(max([-10, -2, -7])).toBe(-2);
    });

    it('supports decimals', () => {
      expect(max([0.1, 0.25, 0.2])).toBeCloseTo(0.25);
    });
  });

  describe('min', () => {
    it('returns the smallest value', () => {
      expect(min([1, 8, 3, 0])).toBe(0);
    });

    it('supports negative numbers', () => {
      expect(min([-10, -2, -7])).toBe(-10);
    });

    it('supports decimals', () => {
      expect(min([0.1, 0.25, 0.2])).toBeCloseTo(0.1);
    });
  });

  describe('sum', () => {
    it('adds every value', () => {
      expect(sum([1, 2, 3, 4])).toBe(10);
    });

    it('supports negative numbers', () => {
      expect(sum([-5, 2, -1])).toBe(-4);
    });

    it('supports decimals', () => {
      expect(sum([0.1, 0.2, 0.3])).toBeCloseTo(0.6);
    });
  });

  describe('average', () => {
    it('divides the sum by the number of values', () => {
      expect(average([1, 2, 3])).toBe(2);
    });

    it('keeps a fractional result', () => {
      expect(average([1, 2])).toBeCloseTo(1.5);
    });

    it('supports negative numbers and decimals', () => {
      expect(average([-1.5, 0.5, 1])).toBeCloseTo(0);
    });
  });

  describe('isDiagonal', () => {
    const epsilon = 1e-10;

    it('returns true for a diagonal matrix', () => {
      expect(
        isDiagonal(
          [
            [2, 0, 0],
            [0, -3, 0],
            [0, 0, 4.5],
          ],
          epsilon,
        ),
      ).toBe(true);
    });

    it('returns false for a non-diagonal matrix', () => {
      expect(
        isDiagonal(
          [
            [1, 2],
            [0, 1],
          ],
          epsilon,
        ),
      ).toBe(false);
    });

    it('treats values within epsilon as zero', () => {
      expect(
        isDiagonal(
          [
            [1, 1e-12],
            [-1e-12, 1],
          ],
          epsilon,
        ),
      ).toBe(true);
    });

    it('treats a 1x1 matrix as diagonal', () => {
      expect(isDiagonal([[-7]], epsilon)).toBe(true);
    });

    it('allows a rectangular diagonal matrix', () => {
      expect(
        isDiagonal(
          [
            [5, 0, 0],
            [0, 6, 0],
          ],
          epsilon,
        ),
      ).toBe(true);
    });
  });
});
