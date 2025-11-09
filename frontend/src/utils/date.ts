import { format, parse, isValid } from 'date-fns';

/**
 * Format date to YYYY-MM-DD
 */
export const formatDateISO = (date: Date | string): string => {
  const d = typeof date === 'string' ? new Date(date) : date;
  return format(d, 'yyyy-MM-dd');
};

/**
 * Format date to readable format
 */
export const formatDateReadable = (date: Date | string): string => {
  const d = typeof date === 'string' ? new Date(date) : date;
  return format(d, 'MMMM d, yyyy');
};

/**
 * Format time to HH:mm:ss
 */
export const formatTime = (date: Date | string): string => {
  const d = typeof date === 'string' ? new Date(date) : date;
  return format(d, 'HH:mm:ss');
};

/**
 * Parse date string
 */
export const parseDate = (dateString: string, formatString: string = 'yyyy-MM-dd'): Date | null => {
  const parsed = parse(dateString, formatString, new Date());
  return isValid(parsed) ? parsed : null;
};

/**
 * Get today's date in ISO format
 */
export const getTodayISO = (): string => {
  return formatDateISO(new Date());
};

/**
 * Get date range for current month
 */
export const getCurrentMonthRange = (): { from: string; to: string } => {
  const now = new Date();
  const from = new Date(now.getFullYear(), now.getMonth(), 1);
  const to = new Date(now.getFullYear(), now.getMonth() + 1, 0);

  return {
    from: formatDateISO(from),
    to: formatDateISO(to),
  };
};
