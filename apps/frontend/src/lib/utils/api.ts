export function checkSuccess(res: Response, payload: any): boolean {
  return payload && (
    payload.success === true || 
    (res.ok && (payload.code === undefined || (payload.code >= 200 && payload.code < 300)))
  );
}

export function getErrorMessage(payload: any, fallback: string): string {
  return payload?.detail || payload?.title || payload?.message || payload?.error || fallback;
}
