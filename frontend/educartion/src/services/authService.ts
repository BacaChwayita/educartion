import type {
  LoginClientRequest,
  RegisterClientRequest,
  LoginResponse,
} from "@/lib/auth-contract";
import { readBackendError } from "@/lib/auth-errors";

type ServiceResult<T> = {
  ok: boolean;
  data?: T;
  error?: string;
  status?: number;
};

export async function login(
  payload: LoginClientRequest,
): Promise<ServiceResult<LoginResponse>> {
  try {
    const res = await fetch("/api/auth/login", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(payload),
    });

    if (!res.ok) {
      return { ok: false, error: await readBackendError(res), status: res.status };
    }

    const data = (await res.json()) as LoginResponse;
    return { ok: true, data };
  } catch (err) {
    return { ok: false, error: "Network error while contacting auth endpoint." };
  }
}

export async function register(
  payload: RegisterClientRequest,
): Promise<ServiceResult<null>> {
  try {
    const res = await fetch("/api/auth/register", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(payload),
    });

    if (!res.ok) {
      return { ok: false, error: await readBackendError(res), status: res.status };
    }

    return { ok: true };
  } catch (err) {
    return { ok: false, error: "Network error while contacting register endpoint." };
  }
}
