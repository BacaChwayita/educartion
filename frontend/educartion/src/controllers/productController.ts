import type { Product } from "@/lib/product-contract";
import * as productService from "@/services/productService";

export async function handleLoadProducts(): Promise<{
  ok: boolean;
  data?: Product[];
  error?: string;
}> {
  const result = await productService.getProducts();
  if (!result.ok) {
    return { ok: false, error: result.error ?? "Failed to load products." };
  }

  return { ok: true, data: result.data ?? [] };
}
