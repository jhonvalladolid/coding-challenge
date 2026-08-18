import Joi from 'joi';
import { env } from '../../config/env.js';
import { AppError } from '../../errors/app-error.js';

function matrixSchema(field) {
  return Joi.array()
    .required()
    .min(1)
    .items(
      Joi.array()
        .min(1)
        .items(
          Joi.number()
            .strict()
            .required()
            .custom((value, helpers) => {
              if (!Number.isFinite(value)) {
                return helpers.error('number.finite');
              }
              return value;
            }),
        )
        .required(),
    )
    .custom((matrix, helpers) => {
      const columnCount = matrix[0].length;
      if (matrix.some((row) => row.length !== columnCount)) {
        return helpers.error('matrix.irregular');
      }

      if (matrix.length > env.maxMatrixDim || columnCount > env.maxMatrixDim) {
        return helpers.error('matrix.tooLarge');
      }

      return matrix;
    })
    .messages({
      'any.required': `${field} is required`,
      'array.base': `${field} must be a matrix`,
      'array.min': `${field} must not be empty`,
      'array.includesRequiredUnknowns': `${field} values must be finite numbers`,
      'number.base': `${field} values must be finite numbers`,
      'number.finite': `${field} values must be finite numbers`,
      'matrix.irregular': `${field} rows must have equal length`,
      'matrix.tooLarge': `${field} dimensions exceed the maximum allowed`,
    });
}

const statisticsBodySchema = Joi.object({
  matrices: Joi.object({
    q: matrixSchema('matrices.q'),
    r: matrixSchema('matrices.r'),
  })
    .required()
    .messages({
      'any.required': 'matrices is required',
      'object.base': 'matrices must be an object',
    }),
})
  .required()
  .unknown(false);

export function assertJsonContentType(req, res, next) {
  const contentType = req.header('content-type');
  if (!contentType || !contentType.toLowerCase().includes('application/json')) {
    next(AppError.unsupportedMediaType());
    return;
  }

  next();
}

export function validateStatisticsBody(req, res, next) {
  const { error, value } = statisticsBodySchema.validate(req.body, {
    abortEarly: false,
    convert: false,
  });

  if (error) {
    const isTooLarge = error.details.some((detail) => detail.type === 'matrix.tooLarge');
    if (isTooLarge) {
      next(AppError.payloadTooLarge());
      return;
    }

    next(
      AppError.validation(
        error.details.map((detail) => ({
          field: detail.path.join('.') || 'body',
          reason: detail.message,
        })),
      ),
    );
    return;
  }

  req.body = value;
  next();
}
