export function toStatisticsResponse({ statistics, diagonal }) {
  return {
    statistics: {
      max: statistics.max,
      min: statistics.min,
      average: statistics.average,
      sum: statistics.sum,
    },
    diagonal: {
      q: diagonal.q,
      r: diagonal.r,
      anyDiagonal: diagonal.anyDiagonal,
    },
  };
}
