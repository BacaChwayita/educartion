import { NextResponse } from "next/server";

import { readBackendError } from "@/lib/auth-errors";
import { getBackendAuthUrl } from "@/lib/backend-auth";

export async function POST(request: Request) {
  const requestBody = await request.json();

  const backendResponse = await fetch(getBackendAuthUrl("/api/payments"), {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      ...(request.headers.get("authorization")
        ? { Authorization: request.headers.get("authorization")! }
        : {}),
    },
    body: JSON.stringify(requestBody),
    cache: "no-store",
  });

  if (!backendResponse.ok) {
    return NextResponse.json(
      { error: await readBackendError(backendResponse) },
      { status: backendResponse.status },
    );
  }

  const data = await backendResponse.json();
  return NextResponse.json(data, { status: backendResponse.status });
}
