import { NextResponse } from "next/server";

import { getBackendAuthUrl } from "@/lib/backend-auth";
import { readBackendError } from "@/lib/auth-errors";
import { type OrderDetail } from "@/lib/orders-contract";

export async function GET(
  _request: Request,
  context: { params: { id: string } } | { params: Promise<{ id: string }> },
) {
  const resolvedParams = await Promise.resolve(context.params);
  const orderId = typeof resolvedParams?.id === "string" ? resolvedParams.id.trim() : "";

  if (!orderId) {
    return NextResponse.json({ error: "Order id is required." }, { status: 400 });
  }

  const backendResponse = await fetch(getBackendAuthUrl(`/api/orders/${orderId}`), {
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

  const data = (await backendResponse.json()) as OrderDetail;
  return NextResponse.json(data, { status: 200 });
}
