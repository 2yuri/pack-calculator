const API_BASE_URL =
  import.meta.env.VITE_API_BASE_URL || "http://localhost:8080/api/v1";

export class ApiRequestError extends Error {
  status: number;

  constructor(status: number, message: string) {
    super(message);
    this.name = "ApiRequestError";
    this.status = status;
  }
}

export async function apiRequest<T>(
  path: string,
  options: RequestInit = {},
): Promise<T> {
  const response = await fetch(`${API_BASE_URL}${path}`, {
    ...options,
    headers: {
      "Content-Type": "application/json",
      ...options.headers,
    },
  });

  if (response.status === 204) {
    return [] as unknown as T;
  }

  if (!response.ok) {
    let message = "Something went wrong";
    try {
      const body = await response.json();
      if (body.error) {
        message = body.error;
      }
    } catch {
      // use default message
    }
    throw new ApiRequestError(response.status, message);
  }

  return response.json();
}
