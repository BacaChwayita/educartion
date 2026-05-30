export const API_AUTH_BASE_URL =
  process.env.API_URL ?? "http://localhost:8080";

export interface LoginClientRequest {
  email: string;
  password: string;
  rememberMe: boolean;
}

export interface RegisterClientRequest {
  fullName: string;
  email: string;
  password: string;
}

export interface LoginRequest {
  email: string;
  password_hash: string;
}

export interface RegisterRequest {
  full_name: string;
  email: string;
  password_hash: string;
}

export interface LoginResponse {
  account_id: number;
  full_name: string;
  email: string;
  role: string;
  token: string;
}

export interface BackendErrorResponse {
  error?: string;
}
