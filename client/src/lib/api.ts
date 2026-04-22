export interface HealthData {
    version: string;
}

export interface APIResponse<T> {
    message: string;
    status: number;
    data: T;
}

export async function fetchHealth(): Promise<APIResponse<HealthData>> {
    const response = await fetch('/api/health');
    if (!response.ok) {
        throw new Error(`API error: ${response.statusText}`);
    }
    return response.json();
}
