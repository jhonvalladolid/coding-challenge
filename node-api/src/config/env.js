import dotenv from 'dotenv';
import path from 'node:path';
import { fileURLToPath } from 'node:url';

const rootDir = path.resolve(path.dirname(fileURLToPath(import.meta.url)), '../..');
dotenv.config({ path: path.join(rootDir, '.env'), quiet: true });

function parsePositiveInt(name, fallback) {
  const raw = process.env[name];
  if (raw === undefined || raw === '') {
    return fallback;
  }

  const value = Number(raw);
  if (!Number.isInteger(value) || value <= 0) {
    throw new Error(`${name} must be a positive integer`);
  }

  return value;
}

function parsePositiveNumber(name, fallback) {
  const raw = process.env[name];
  if (raw === undefined || raw === '') {
    return fallback;
  }

  const value = Number(raw);
  if (!Number.isFinite(value) || value <= 0) {
    throw new Error(`${name} must be a positive finite number`);
  }

  return value;
}

export const env = {
  port: parsePositiveInt('PORT', 3000),
  appEnv: process.env.APP_ENV || process.env.NODE_ENV || 'development',
  logLevel: process.env.LOG_LEVEL || 'info',
  maxMatrixDim: parsePositiveInt('MAX_MATRIX_DIM', 200),
  jsonBodyLimit: process.env.JSON_BODY_LIMIT || '1mb',
  diagonalEpsilon: parsePositiveNumber('DIAGONAL_EPSILON', 1e-10),
};
