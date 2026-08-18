export function success(data) {
  return {
    success: true,
    data,
  };
}

export function failure({ code, message, details, requestId }) {
  const error = {
    code,
    message,
    requestId,
  };

  if (details && details.length > 0) {
    error.details = details;
  }

  return {
    success: false,
    error,
  };
}
