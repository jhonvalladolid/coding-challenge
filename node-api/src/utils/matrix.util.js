export function flattenMatrix(matrix) {
  return matrix.flat();
}

export function flattenMatrices(...matrices) {
  return matrices.flatMap((matrix) => flattenMatrix(matrix));
}

export function max(values) {
  if (values.length === 0) {
    throw new RangeError('max requires a non-empty list');
  }

  return values.reduce((current, value) => (value > current ? value : current));
}

export function min(values) {
  if (values.length === 0) {
    throw new RangeError('min requires a non-empty list');
  }

  return values.reduce((current, value) => (value < current ? value : current));
}

export function sum(values) {
  return values.reduce((total, value) => total + value, 0);
}

export function average(values) {
  if (values.length === 0) {
    throw new RangeError('average requires a non-empty list');
  }

  return sum(values) / values.length;
}

export function isDiagonal(matrix, epsilon) {
  for (let i = 0; i < matrix.length; i += 1) {
    const row = matrix[i];
    for (let j = 0; j < row.length; j += 1) {
      if (i !== j && Math.abs(row[j]) > epsilon) {
        return false;
      }
    }
  }

  return true;
}
