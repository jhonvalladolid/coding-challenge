import { AppError } from '../errors/app-error.js';

export function notFound(req, res, next) {
  next(AppError.notFound());
}

export function methodNotAllowed(req, res, next) {
  next(AppError.methodNotAllowed());
}
