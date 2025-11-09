import { Middleware } from '@reduxjs/toolkit';

/**
 * Logger middleware for development
 * Logs all actions and state changes
 */
const loggerMiddleware: Middleware = (store) => (next) => (action) => {
  if (process.env.NODE_ENV === 'development') {
    console.group(action.type);
    console.info('dispatching', action);
    const result = next(action);
    console.log('next state', store.getState());
    console.groupEnd();
    return result;
  }

  return next(action);
};

export default loggerMiddleware;
