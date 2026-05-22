import { NextResponse } from "next/server";
import { type LoginClientRequest, type LoginResponse } from "@/lib/auth-contract";
import { readBackendError } from "@/lib/auth-errors";
import { getBackendAuthUrl } from "@/lib/backend-auth";

const SESSION_COOKIE_NAME = "educartion_auth_token";

export async function POST(request: Request) {
  let payload: LoginClientRequest;

  try {
    payload = (await request.json()) as LoginClientRequest;
  } catch {
    return NextResponse.json({ error: "Invalid request body." }, { status: 400 });
  }

  const email = payload.email?.trim();
  const password = payload.password?.trim();

  if (!email || !password) {
    return NextResponse.json({ error: "Email and password are required." }, { status: 400 });
  }

  const backendResponse = await fetch(getBackendAuthUrl("/api/auth/login"), {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      email,
      password_hash: password, // TODO: should be password_text - issue on backend
    }),
  });

  if (!backendResponse.ok) {
    return NextResponse.json(
      { error: await readBackendError(backendResponse) },
      { status: backendResponse.status },
    );
  }

  const data = (await backendResponse.json()) as LoginResponse;
  const response = NextResponse.json(data, { status: backendResponse.status });
  const cookieOptions = {
    httpOnly: true,
    sameSite: "lax" as const,
    secure: process.env.NODE_ENV === "production",
    path: "/",
    ...(payload.rememberMe ? { maxAge: 60 * 60 * 24 * 30 } : {}),
  };

  response.cookies.set(SESSION_COOKIE_NAME, data.token, cookieOptions);

  return response;
}
