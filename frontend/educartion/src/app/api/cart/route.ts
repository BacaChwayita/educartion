import { NextResponse } from "next/server";

import { readBackendError } from "@/lib/auth-errors";
import { getBackendAuthUrl } from "@/lib/backend-auth";
import type { CartItem } from "@/lib/cart-contract";

export async function GET(request: Request) {
  const backendResponse = await fetch(getBackendAuthUrl("/api/cart/items"), {
    method: "GET",
    headers: {
      "Content-Type": "application/json",
      ...(request.headers.get("authorization")
        ? { Authorization: request.headers.get("authorization")! }
        : {}),
    },
    cache: "no-store",
  });

  if (!backendResponse.ok) {
    return NextResponse.json(
      { error: await readBackendError(backendResponse) },
      { status: backendResponse.status },
    );
  }

  const data = (await backendResponse.json()) as CartItem[];
  return NextResponse.json(data, { status: 200 });
}
