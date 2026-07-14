import { NextResponse } from "next/server";
import { getBackendAuthUrl } from "@/lib/backend-auth";
import { readBackendError } from "@/lib/auth-errors";
import { type OrderSummary } from "@/lib/orders-contract";

export async function GET() {
  const backendResponse = await fetch(getBackendAuthUrl("/api/orders/"), {
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

  const data = (await backendResponse.json()) as OrderSummary[];
  return NextResponse.json(data, { status: 200 });
}
