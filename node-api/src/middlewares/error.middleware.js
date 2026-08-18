import { AppError } from '../errors/app-error.js';
import { logger } from '../config/logger.js';
import { failure } from '../utils/response.util.js';

function toAppError(err) {
  if (err instanceof AppError) {
    return err;
  }

  if (err.type === 'entity.too.large' || err.status === 413) {
    return AppError.payloadTooLarge();
  }

  if (err instanceof SyntaxError && err.status === 400 && 'body' in err) {
    return AppError.validation([
      { field: 'body', reason: 'malformed JSON' },
    ]);
  }

  return AppError.internal();
}

export function errorHandler(err, req, res, next) {
  if (res.headersSent) {
    next(err);
    return;
  }

  const appError = toAppError(err);

  if (appError.statusCode >= 500) {
    logger.error({ err, requestId: req.requestId }, 'unhandled error');
  }

  res.status(appError.statusCode).json(
    failure({
      code: appError.code,
      message: appError.message,
      details: appError.details,
      requestId: req.requestId,
    }),
  );
}
