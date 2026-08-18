import { randomUUID } from 'node:crypto';

const REQUEST_ID_HEADER = 'X-Request-ID';

export function requestId(req, res, next) {
  const headerValue = req.header(REQUEST_ID_HEADER);
  const id = headerValue && headerValue.trim() !== '' ? headerValue.trim() : randomUUID();

  req.requestId = id;
  res.setHeader(REQUEST_ID_HEADER, id);
  next();
}
