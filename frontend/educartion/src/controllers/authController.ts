import type { LoginClientRequest, RegisterClientRequest } from "@/lib/auth-contract";
import * as authService from "@/services/authService";

export async function handleLogin(payload: LoginClientRequest) {
  // Basic client-side validation
  if (!payload.email || !payload.password) {
    return { ok: false, error: "Email and password are required." };
  }

  const res = await authService.login(payload);
  if (!res.ok) {
    return { ok: false, error: res.error ?? "Login failed." };
  }

  return { ok: true, data: res.data };
}

export async function handleRegister(payload: RegisterClientRequest) {
  if (!payload.fullName || !payload.email || !payload.password) {
    return { ok: false, error: "Full name, email, and password are required." };
  }

  const res = await authService.register(payload);
  if (!res.ok) {
    return { ok: false, error: res.error ?? "Register failed." };
  }

  return { ok: true };
}
