import { Middleware } from '@reduxjs/toolkit';

/**
 * Generate a UUID v4
 */
function generateUUID(): string {
  return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function (c) {
    const r = (Math.random() * 16) | 0;
    const v = c === 'x' ? r : (r & 0x3) | 0x8;
    return v.toString(16);
  });
}

/**
 * Middleware to attach idempotency keys to outgoing mutations
 * This ensures that retried requests are safe and don't create duplicates
 */
const idempotencyMiddleware: Middleware = (_store) => (next) => (action) => {
  // Check if this is a mutation action that should have an idempotency key
  if (
    action.type &&
    (action.type.includes('/pending') ||
      action.type.includes('checkIn') ||
      action.type.includes('checkOut'))
  ) {
    // Attach idempotency key to action meta
    if (!action.meta) {
      action.meta = {};
    }

    if (!action.meta.idempotencyKey) {
      action.meta.idempotencyKey = generateUUID();
    }
  }

  return next(action);
};

export default idempotencyMiddleware;
