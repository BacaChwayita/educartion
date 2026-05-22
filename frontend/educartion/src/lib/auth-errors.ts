import type { BackendErrorResponse } from "@/lib/auth-contract";

export async function readBackendError(response: Response): Promise<string> {
  const rawBody = await response.text();

  if (!rawBody) {
    return `Request failed with status ${response.status}`;
  }

  try {
    const parsed = JSON.parse(rawBody) as BackendErrorResponse | { message?: string } | string;

    if (typeof parsed === "string") {
      return parsed;
    }

    if (parsed && typeof parsed === "object") {
      if (typeof parsed.error === "string" && parsed.error.trim()) {
        return parsed.error;
      }

      if (typeof parsed.message === "string" && parsed.message.trim()) {
        return parsed.message;
      }
    }
  } catch {
    // Fall through to the raw text below.
  }

  return rawBody.trim() || `Request failed with status ${response.status}`;
}