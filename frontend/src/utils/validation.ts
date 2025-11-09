/**
 * Validate email format
 */
export const isValidEmail = (email: string): boolean => {
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
  return emailRegex.test(email);
};

/**
 * Validate username (alphanumeric and underscore, 3-20 chars)
 */
export const isValidUsername = (username: string): boolean => {
  const usernameRegex = /^[a-zA-Z0-9_]{3,20}$/;
  return usernameRegex.test(username);
};

/**
 * Validate password (min 8 chars, at least one uppercase, one lowercase, one number)
 */
export const isValidPassword = (password: string): boolean => {
  const passwordRegex = /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d).{8,}$/;
  return passwordRegex.test(password);
};

/**
 * Validate phone number (basic international format)
 */
export const isValidPhoneNumber = (phone: string): boolean => {
  const phoneRegex = /^\+?[\d\s\-()]{10,}$/;
  return phoneRegex.test(phone);
};

/**
 * Validate date format YYYY-MM-DD
 */
export const isValidDateFormat = (date: string): boolean => {
  const dateRegex = /^\d{4}-\d{2}-\d{2}$/;
  if (!dateRegex.test(date)) return false;

  const [year, month, day] = date.split('-').map(Number);
  const d = new Date(year, month - 1, day);

  return d.getFullYear() === year && d.getMonth() === month - 1 && d.getDate() === day;
};

/**
 * Validate latitude and longitude
 */
export const isValidCoordinates = (
  latitude: number,
  longitude: number
): boolean => {
  return latitude >= -90 && latitude <= 90 && longitude >= -180 && longitude <= 180;
};

/**
 * Trim and validate string is not empty
 */
export const isNotEmpty = (value: string): boolean => {
  return value.trim().length > 0;
};

/**
 * Validate string length
 */
export const isValidLength = (value: string, min: number, max: number): boolean => {
  return value.length >= min && value.length <= max;
};
