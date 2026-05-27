export function checkWarrantyExpiry(exitDateStr: string | undefined, warrantyDays: number) {
  if (!exitDateStr) {
    return {
      expiryDate: new Date(),
      isValid: false,
      remainingDays: 0,
      formattedExpiry: '-'
    };
  }
  const exitDate = new Date(exitDateStr);
  const expiryDate = new Date(exitDate.getTime() + warrantyDays * 24 * 60 * 60 * 1000);
  const today = new Date();
  const remainingDays = Math.ceil((expiryDate.getTime() - today.getTime()) / (24 * 60 * 60 * 1000));
  
  return {
    expiryDate,
    isValid: remainingDays >= 0,
    remainingDays: remainingDays >= 0 ? remainingDays : 0,
    formattedExpiry: expiryDate.toLocaleDateString('id-ID', { day: 'numeric', month: 'long', year: 'numeric' })
  };
}

export function getWarrantyExpiryDate(exitDateStr: string | undefined, days: number): string {
  if (!exitDateStr) return '-';
  const d = new Date(exitDateStr);
  d.setDate(d.getDate() + days);
  return d.toLocaleDateString('id-ID', {
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  });
}
