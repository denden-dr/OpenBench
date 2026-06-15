/**
 * Utility functions for formatting numbers and currencies.
 */

/**
 * Formats a number to IDR standard string (e.g. 300000 -> "300.000").
 * Strips any non-digit characters.
 */
export function formatCurrencyInput(value: number | string): string {
  if (value === undefined || value === null || value === '') return '';
  
  // Strip non-digits
  const clean = String(value).replace(/\D/g, '');
  if (!clean) return '';
  
  return Number(clean).toLocaleString('id-ID');
}

/**
 * Parses a formatted currency string back to a number (e.g. "300.000" -> 300000).
 */
export function parseCurrencyInput(value: string): number {
  if (!value) return 0;
  const clean = value.replace(/\D/g, '');
  if (!clean) return 0;
  return Number(clean);
}
