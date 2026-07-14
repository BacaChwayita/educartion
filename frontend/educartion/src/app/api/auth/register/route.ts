import { NextResponse } from "next/server";
import { type RegisterClientRequest } from "@/lib/auth-contract";
import { readBackendError } from "@/lib/auth-errors";
import { getBackendAuthUrl } from "@/lib/backend-auth";

export async function POST(request: Request) {
  let payload: RegisterClientRequest;

  try {
    payload = (await request.json()) as RegisterClientRequest;
  } catch {
    return NextResponse.json({ error: "Invalid request body." }, { status: 400 });
  }

  const fullName = payload.fullName?.trim();
  const email = payload.email?.trim();
  const password = payload.password?.trim();

  if (!fullName || !email || !password) {
    return NextResponse.json(
      { error: "Full name, email, and password are required." },
      { status: 400 },
    );
  }

  const backendResponse = await fetch(getBackendAuthUrl("/api/auth/register"), {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      full_name: fullName,
      email,
      password_hash: password,
    }),
  });

  if (!backendResponse.ok) {
    return NextResponse.json(
      { error: await readBackendError(backendResponse) },
      { status: backendResponse.status },
    );
  }

  return NextResponse.json({ ok: true }, { status: 201 });
}