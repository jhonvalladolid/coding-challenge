export class AppError extends Error {
  constructor({ code, message, statusCode, details }) {
    super(message);
    this.name = 'AppError';
    this.code = code;
    this.statusCode = statusCode;
    this.details = details;
  }

  static validation(details, message = 'The request body is invalid') {
    return new AppError({
      code: 'VALIDATION_ERROR',
      message,
      statusCode: 400,
      details,
    });
  }

  static notFound(message = 'The requested resource was not found') {
    return new AppError({
      code: 'NOT_FOUND',
      message,
      statusCode: 404,
    });
  }

  static methodNotAllowed(message = 'The HTTP method is not allowed for this resource') {
    return new AppError({
      code: 'METHOD_NOT_ALLOWED',
      message,
      statusCode: 405,
    });
  }

  static payloadTooLarge(message = 'The request payload exceeds the allowed size') {
    return new AppError({
      code: 'PAYLOAD_TOO_LARGE',
      message,
      statusCode: 413,
    });
  }

  static unsupportedMediaType(message = 'Content-Type must be application/json') {
    return new AppError({
      code: 'UNSUPPORTED_MEDIA_TYPE',
      message,
      statusCode: 415,
    });
  }

  static internal(message = 'An unexpected error occurred') {
    return new AppError({
      code: 'INTERNAL_ERROR',
      message,
      statusCode: 500,
    });
  }
}
