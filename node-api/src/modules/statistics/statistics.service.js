import { env } from '../../config/env.js';
import {
  average,
  flattenMatrices,
  isDiagonal,
  max,
  min,
  sum,
} from '../../utils/matrix.util.js';
import { toStatisticsResponse } from './statistics.dto.js';

export function computeStatistics(q, r) {
  const values = flattenMatrices(q, r);
  const diagonalQ = isDiagonal(q, env.diagonalEpsilon);
  const diagonalR = isDiagonal(r, env.diagonalEpsilon);

  return toStatisticsResponse({
    statistics: {
      max: max(values),
      min: min(values),
      average: average(values),
      sum: sum(values),
    },
    diagonal: {
      q: diagonalQ,
      r: diagonalR,
      anyDiagonal: diagonalQ || diagonalR,
    },
  });
}
