import { handleLoadProducts } from "@/controllers/productController";
import type { Product } from "@/lib/product-contract";

export default async function ProductsPage() {
  const result = await handleLoadProducts();
  const products: Product[] = result.ok ? result.data ?? [] : [];

  return (
    <main className="relative min-h-screen overflow-hidden bg-[radial-gradient(circle_at_top_left,rgba(251,191,36,0.2),transparent_30%),radial-gradient(circle_at_bottom_right,rgba(14,165,233,0.18),transparent_28%),linear-gradient(180deg,#060816_0%,#0b1020_100%)] px-6 py-10 text-slate-100 sm:px-8 lg:px-10">
      <div className="absolute inset-0 bg-[linear-gradient(rgba(255,255,255,0.03)_1px,transparent_1px),linear-gradient(90deg,rgba(255,255,255,0.03)_1px,transparent_1px)] bg-size-[28px_28px] opacity-20" />
      
      <div className="relative mx-auto max-w-6xl">
        <div className="mb-8 space-y-3">
          <p className="text-sm font-medium uppercase tracking-[0.3em] text-amber-200">
            Browse
          </p>
          <h1 className="text-4xl font-semibold tracking-tight text-white">Products</h1>
          <p className="text-sm leading-6 text-slate-300">
            Discover our collection of available products.
          </p>
        </div>

        {!result.ok ? (
          <div className="mt-10 rounded-3xl border border-rose-400/30 bg-rose-400/10 p-8 text-center shadow-2xl shadow-black/30 backdrop-blur-xl">
            <p className="text-lg font-semibold text-rose-100">Failed to load products</p>
            <p className="mt-2 text-sm text-rose-200">{result.error}</p>
          </div>
        ) : products.length === 0 ? (
          <div className="mt-10 rounded-3xl border border-white/10 bg-white/5 p-8 text-center shadow-2xl shadow-black/30 backdrop-blur-xl">
            <p className="text-lg font-semibold text-slate-100">No products were found.</p>
            <p className="mt-2 text-sm text-slate-400">
              Make sure the backend is running and that the <code className="rounded bg-slate-950/50 px-2 py-1 text-xs font-mono text-slate-200">/api/products</code> endpoint returns products.
            </p>
          </div>
        ) : (
          <div className="mt-10 grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
            {products.map((product, index) => (
              <article
                key={product.id ?? index}
                className="group overflow-hidden rounded-3xl border border-white/10 bg-white/5 p-6 shadow-2xl shadow-black/30 transition hover:-translate-y-1 hover:bg-white/8 hover:border-amber-300/40 backdrop-blur-xl"
              >
                <h2 className="text-xl font-semibold text-white group-hover:text-amber-200 transition">{product.name ?? "Unnamed product"}</h2>
                {product.description ? (
                  <p className="mt-3 text-sm leading-6 text-slate-300">{product.description}</p>
                ) : (
                  <p className="mt-3 text-sm leading-6 text-slate-500">No description available.</p>
                )}
                {typeof product.price === "number" ? (
                  <p className="mt-5 text-lg font-semibold text-amber-300">
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
