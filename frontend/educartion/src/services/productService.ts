import type { Product, ProductsResponse } from "@/lib/product-contract";

type ServiceResult<T> = {
  ok: boolean;
  data?: T;
  error?: string;
  status?: number;
};

export async function getProducts(): Promise<ServiceResult<Product[]>> {
  try {
    const res = await fetch("/api/products", {
      cache: "no-store",
    });

    if (!res.ok) {
      const body = (await res.json().catch(() => null)) as { error?: string } | null;
      return {
        ok: false,
        error: body?.error ?? `Unable to load products: ${res.status}`,
        status: res.status,
      };
    }

    const data = (await res.json()) as ProductsResponse;
    const products = normalizeProducts(data);
    return { ok: true, data: products };
  } catch (err) {
    const errorMessage = err instanceof Error ? err.message : "Unknown error";
    return {
      ok: false,
      error: `Failed to fetch products: ${errorMessage}`,
    };
  }
}

function normalizeProducts(payload: unknown): Product[] {
  if (!payload) {
    return [];
  }

  if (Array.isArray(payload)) {
    return payload as Product[];
  }

  if (typeof payload === "object") {
    const obj = payload as Record<string, unknown>;
    // If it has a products property, use that
    if (obj.products && Array.isArray(obj.products)) {
      return obj.products as Product[];
    }
    // Otherwise, treat the object values as products
    return Object.values(obj).filter((item) => item && typeof item === "object") as Product[];
  }

  return [];
}
