import { handleLoadProducts } from "@/controllers/productController";
import type { Product } from "@/lib/product-contract";

export default async function ProductsPage() {
  const result = await handleLoadProducts();
  const products: Product[] = result.ok ? result.data ?? [] : [];

  return (
    <main className="min-h-screen bg-zinc-50 text-zinc-950 px-6 py-10">
      <div className="mx-auto max-w-6xl">
        <h1 className="text-4xl font-semibold tracking-tight">Products</h1>
        <p className="mt-2 text-sm text-zinc-600">
          Loaded from the backend API.
        </p>

        {!result.ok ? (
          <div className="mt-10 rounded-3xl border border-dashed border-zinc-300 bg-white p-10 text-center shadow-sm">
            <p className="text-lg font-medium text-red-700">Failed to load products</p>
            <p className="mt-2 text-sm text-red-600">{result.error}</p>
          </div>
        ) : products.length === 0 ? (
          <div className="mt-10 rounded-3xl border border-dashed border-zinc-300 bg-white p-10 text-center shadow-sm">
            <p className="text-lg font-medium text-zinc-700">No products were found.</p>
            <p className="mt-2 text-sm text-zinc-500">
              Make sure the backend is running and that the <code className="rounded bg-white px-2 py-1 text-xs font-mono">/api/products</code> endpoint returns products.
            </p>
          </div>
        ) : (
          <div className="mt-10 grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
            {products.map((product, index) => (
              <article
                key={product.id ?? index}
                className="overflow-hidden rounded-3xl border border-zinc-200 bg-white p-6 shadow-sm transition hover:-translate-y-1 hover:shadow-md"
              >
                <h2 className="text-xl font-semibold text-zinc-950">{product.name ?? "Unnamed product"}</h2>
                {product.description ? (
                  <p className="mt-3 text-sm leading-6 text-zinc-600">{product.description}</p>
                ) : (
                  <p className="mt-3 text-sm leading-6 text-zinc-500">No description available.</p>
                )}
                {typeof product.price === "number" ? (
                  <p className="mt-5 text-lg font-semibold text-zinc-900">
                    ${product.price.toFixed(2)}
                  </p>
                ) : null}
              </article>
            ))}
          </div>
        )}
      </div>
    </main>
  );
}
