import { NextResponse } from "next/server";

import { readBackendError } from "@/lib/auth-errors";
import { getBackendAuthUrl } from "@/lib/backend-auth";
import type { Product } from "@/lib/product-contract";

export async function GET() {
  const backendResponse = await fetch(getBackendAuthUrl("/api/products"), {
    method: "GET",
    headers: {
      "Content-Type": "application/json",
    },
    cache: "no-store",
  });

  if (!backendResponse.ok) {
    return NextResponse.json(
      { error: await readBackendError(backendResponse) },
      { status: backendResponse.status },
    );
  }

  const data = (await backendResponse.json()) as Product[];
  return NextResponse.json(data, { status: 200 });
}
