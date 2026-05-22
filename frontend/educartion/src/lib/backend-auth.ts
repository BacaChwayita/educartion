import { BACKEND_AUTH_BASE_URL } from "@/lib/auth-contract";

export function getBackendAuthUrl(path: string): string {
  return `${BACKEND_AUTH_BASE_URL}${path}`;
}