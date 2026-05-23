type Product = {
  id?: string;
  name?: string;
  description?: string;
  price?: number;
  [key: string]: unknown;
};

async function getProducts() {
  const response = await fetch("/api/products", {
    cache: "no-store",
  });

  if (!response.ok) {
    throw new Error(`Unable to load products: ${response.status}`);
  }

  return response.json();
}

function normalizeProducts(payload: unknown): Product[] {
  if (!payload) {
    return [];
  }

  if (Array.isArray(payload)) {
    return payload;
  }

  if (typeof payload === "object") {
    return Object.values(payload as Record<string, Product>);
  }

  return [];
}

export default async function ProductsPage() {
  const data = await getProducts();
  const products = normalizeProducts((data as { products?: unknown }).products ?? data);

  return (
    <main className="min-h-screen bg-zinc-50 text-zinc-950 px-6 py-10">
      <div className="mx-auto max-w-6xl">
        <h1 className="text-4xl font-semibold tracking-tight">Products</h1>
        <p className="mt-2 text-sm text-zinc-600">
          Loaded from <code className="rounded bg-white px-2 py-1 text-xs font-mono">/api/products</code>.
        </p>

        {products.length === 0 ? (
          <div className="mt-10 rounded-3xl border border-dashed border-zinc-300 bg-white p-10 text-center shadow-sm">
            <p className="text-lg font-medium text-zinc-700">No products were found.</p>
            <p className="mt-2 text-sm text-zinc-500">
              Make sure the backend is running and that <code>/api/products</code> returns a products object.
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
