"use client";

import Link from "next/link";
import { useParams } from "next/navigation";
import { useEffect, useMemo, useState } from "react";

import ThemeToggle from "@/components/theme-toggle";
import { handleLoadProducts } from "@/controllers/productController";
import type { CartItem } from "@/lib/cart-contract";
import type { Product } from "@/lib/product-contract";
import * as cartService from "@/services/cartService";

const productImageBasePath = "/product_images";

function formatCurrency(cents: number): string {
  return new Intl.NumberFormat("en-US", {
    style: "currency",
    currency: "USD",
  }).format(cents / 100);
}

function buildProductImagePath(imageValue?: string): string {
  if (typeof imageValue !== "string" || !imageValue.trim()) {
    return "/images/product-placeholder.png";
  }

  const trimmed = imageValue.trim().replace(/^\/+/, "");
  return `${productImageBasePath}/${trimmed}`;
}

function extractImageUrl(value: unknown): string | undefined {
  if (!value || typeof value !== "object") {
    return undefined;
  }

  const record = value as Record<string, unknown>;
  return (
    (typeof record.image_url === "string" && record.image_url) ||
    (typeof record.imageUrl === "string" && record.imageUrl) ||
    (typeof record.url === "string" && record.url) ||
    undefined
  );
}

function resolveImageFromCollection(images: unknown): string | undefined {
  if (!Array.isArray(images) || images.length === 0) {
    return undefined;
  }

  const primary = images.find((item) => {
    if (!item || typeof item !== "object") {
      return false;
    }
    const record = item as Record<string, unknown>;
    return record.is_primary === true || record.isPrimary === true;
  });

  return extractImageUrl(primary ?? images[0]);
}

function resolveProductId(product: Product, fallbackId: string): string {
  const rawId = product.id ?? product.product_id ?? product["product_id"] ?? fallbackId;
  return typeof rawId === "string" || typeof rawId === "number" ? String(rawId) : fallbackId;
}

function resolveSupplierLabel(product: Product): string {
  const supplierName =
    (typeof product.supplier_name === "string" && product.supplier_name.trim()) ||
    (typeof product["supplier_name"] === "string" && String(product["supplier_name"]).trim()) ||
    (typeof product["supplier"] === "string" && String(product["supplier"]).trim());

  if (supplierName) {
    return supplierName;
  }

  const supplierId = product.supplier_id ?? product["supplier_id"];
  if (typeof supplierId === "string" || typeof supplierId === "number") {
    return `Supplier #${supplierId}`;
  }

  return "Unknown supplier";
}

function resolveStockQuantity(product: Product): number | null {
  const candidates = [
    product.stock_quantity,
    product["stock_quantity"],
    product["stock"],
    product["quantity"],
    product["available"],
  ];

  for (const value of candidates) {
    if (typeof value === "number" && !Number.isNaN(value)) {
      return value;
    }
    if (typeof value === "string" && value.trim()) {
      const parsed = Number(value);
      if (!Number.isNaN(parsed)) {
        return parsed;
      }
    }
  }

  return null;
}

function resolveImageUrl(product: Product): string {
  const image =
    (typeof product.image === "string" && product.image) ||
    (typeof product.image_url === "string" && product.image_url) ||
    (typeof product["image_url"] === "string" && String(product["image_url"])) ||
    (typeof product["image"] === "string" && String(product["image"])) ||
    (typeof product["imageUrl"] === "string" && String(product["imageUrl"]));

  if (image) {
    return buildProductImagePath(image);
  }

  const imageFromCollection =
    resolveImageFromCollection(product["product_image"]) ??
    resolveImageFromCollection(product["product_images"]) ??
    resolveImageFromCollection(product["images"]);

  if (imageFromCollection) {
    return buildProductImagePath(imageFromCollection);
  }

  return "/images/product-placeholder.png";
}

export default function ProductDetailsPage() {
  const params = useParams<{ id: string | string[] }>();
  const routeId = Array.isArray(params.id) ? params.id[0] : params.id;
  const [product, setProduct] = useState<Product | null>(null);
  const [products, setProducts] = useState<Product[]>([]);
  const [quantity, setQuantity] = useState(1);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [successMessage, setSuccessMessage] = useState<string | null>(null);
  const [search, setSearch] = useState("");

  useEffect(() => {
    let isMounted = true;

    async function loadProduct() {
      setIsLoading(true);
      setError(null);

      try {
        const result = await handleLoadProducts();

        if (!isMounted) {
          return;
        }

        if (!result.ok) {
          setError(result.error ?? "Failed to load product.");
          setProduct(null);
          return;
        }

        const found = (result.data ?? []).find((item, index) => {
          const itemId = resolveProductId(item, String(index));
          return itemId === routeId;
        });

        if (!found) {
          setError("Product not found.");
          setProduct(null);
          return;
        }

        setProduct(found);
      } catch (fetchError) {
        if (isMounted) {
          setError(fetchError instanceof Error ? fetchError.message : "Failed to load product.");
          setProduct(null);
        }
      } finally {
        if (isMounted) {
          setIsLoading(false);
        }
      }
    }

    loadProduct();

    return () => {
      isMounted = false;
    };
  }, [routeId]);

  const stockQuantity = useMemo(() => (product ? resolveStockQuantity(product) : null), [product]);
  const normalizedStock = typeof stockQuantity === "number" ? Math.max(0, stockQuantity) : 0;
  const isInStock = normalizedStock > 0;
  const supplierLabel = product ? resolveSupplierLabel(product) : "";
  const imageUrl = product ? resolveImageUrl(product) : "/images/product-placeholder.png";
  const maxQuantity = isInStock ? normalizedStock : 1;

  const updateQuantity = (value: number) => {
    setSuccessMessage(null);
    const safeValue = Number.isFinite(value) ? value : 1;
    const clamped = Math.min(Math.max(1, safeValue), Math.max(1, maxQuantity));
    setQuantity(clamped);
  };

  const handleAddToCart = () => {
    if (!product || !isInStock) {
      return;
    }

    const priceCents = typeof product.price === "number" ? product.price : 0;
    const productId = resolveProductId(product, routeId);
    const newItem: CartItem = {
      id: productId,
      title: product.name ?? "Untitled product",
      price: priceCents / 100,
      image: imageUrl,
      description: product.description,
      qty: quantity,
    };

    const current = cartService.loadCartFromStorage();
    const existingIndex = current.findIndex((item) => String(item.id) === String(productId));
    let updated: CartItem[] = [];

    if (existingIndex >= 0) {
      updated = current.map((item, index) =>
        index === existingIndex ? { ...item, qty: item.qty + quantity } : item,
      );
    } else {
      updated = [...current, newItem];
    }

    cartService.saveCartToStorage(updated);

    setSuccessMessage(`${quantity} item(s) added to cart.`);
  };

  const filteredProducts = useMemo(() => {
    const query = search.trim().toLowerCase();

    if (!query) {
      return products;
    }

    return products.filter((product) => {
      const name = typeof product.name === "string" ? product.name.toLowerCase() : "";
      const description = typeof product.description === "string" ? product.description.toLowerCase() : "";
      const rawId = product.id ?? product["product_id"];
      const idText =
        typeof rawId === "string" || typeof rawId === "number" ? String(rawId).toLowerCase() : "";

      return name.includes(query) || description.includes(query) || idText.includes(query);
    });
  }, [products, search]);

  return (
    <section className="relative min-h-screen overflow-hidden bg-[radial-gradient(circle_at_top_left,var(--hero-accent),transparent_30%),radial-gradient(circle_at_bottom_right,var(--hero-accent-2),transparent_28%),linear-gradient(180deg,var(--hero-bg-top)_0%,var(--hero-bg-bottom)_100%)] px-6 py-10 text-slate-900 sm:px-8 lg:px-10 dark:text-slate-100">
      <div className="absolute inset-0 bg-[linear-gradient(var(--grid-line)_1px,transparent_1px),linear-gradient(90deg,var(--grid-line)_1px,transparent_1px)] bg-size-[28px_28px] opacity-20" />
      <div className="relative mx-auto w-full max-w-[1400px]">
        <header className="rounded-4xl border border-slate-200/80 bg-white/90 p-6 shadow-[0_20px_60px_-24px_var(--shadow-color)] backdrop-blur-xl dark:border-white/10 dark:bg-white/5 dark:shadow-2xl dark:shadow-black/30">
          <div className="flex flex-col gap-4 xl:flex-row xl:items-center xl:gap-6">
            <div className="flex min-w-[170px] items-center justify-start">
              <div className="inline-flex h-11 w-40 items-center justify-center rounded-xl border border-slate-300/80 bg-white/90 text-sm font-semibold uppercase tracking-[0.2em] text-slate-900 shadow-sm dark:border-white/15 dark:bg-slate-950/80 dark:text-amber-100">
                Logo
              </div>
            </div>

            <div className="flex flex-1 justify-center">
              <label className="relative w-full max-w-2xl">
                <input
                  type="text"
                  value={search}
                  onChange={(event) => setSearch(event.target.value)}
                  placeholder="Search products"
                  className="w-full rounded-2xl border border-slate-300/80 bg-white/90 py-3 pl-4 pr-12 text-sm text-slate-900 outline-none transition placeholder:text-slate-500 focus:border-amber-500/60 focus:bg-white focus:ring-2 focus:ring-amber-500/20 dark:border-white/10 dark:bg-white/5 dark:text-slate-100 dark:placeholder:text-slate-500 dark:focus:border-amber-300/60 dark:focus:bg-white/8 dark:focus:ring-amber-300/20"
                />
                <span className="pointer-events-none absolute inset-y-0 right-4 flex items-center text-slate-500 dark:text-slate-400">
                  <svg viewBox="0 0 24 24" className="h-5 w-5" fill="none" stroke="currentColor" strokeWidth="2" aria-hidden="true">
                    <circle cx="11" cy="11" r="7" />
                    <path d="M20 20L16.65 16.65" />
                  </svg>
                </span>
              </label>
            </div>

            <div className="flex min-w-[170px] items-center justify-end gap-2 md:gap-3">
              <Link
                href="/orders"
                className="inline-flex items-center gap-2 rounded-2xl border border-slate-300/80 bg-white/90 px-3 py-2 text-sm font-medium text-slate-800 shadow-sm transition hover:bg-slate-50 dark:border-white/10 dark:bg-white/5 dark:text-slate-200 dark:hover:bg-white/10"
              >
                <svg viewBox="0 0 24 24" className="h-5 w-5" fill="none" stroke="currentColor" strokeWidth="2" aria-hidden="true">
                  <circle cx="12" cy="8" r="4" />
                  <path d="M4 20a8 8 0 0 1 16 0" />
                </svg>
                <span>Account</span>
              </Link>
              <Link
                href="/cart"
                className="inline-flex items-center gap-2 rounded-2xl border border-slate-300/80 bg-white/90 px-3 py-2 text-sm font-medium text-slate-800 shadow-sm transition hover:bg-slate-50 dark:border-white/10 dark:bg-white/5 dark:text-slate-200 dark:hover:bg-white/10"
              >
                <svg viewBox="0 0 24 24" className="h-5 w-5" fill="none" stroke="currentColor" strokeWidth="2" aria-hidden="true">
                  <circle cx="9" cy="20" r="1.5" />
                  <circle cx="17" cy="20" r="1.5" />
                  <path d="M3 4h2l2.4 10.2a1 1 0 0 0 1 .8H18a1 1 0 0 0 1-.8L21 7H7" />
                </svg>
                <span>Cart</span>
              </Link>
              <ThemeToggle />
            </div>
          </div>
        </header>
        <br></br>
        <div className="rounded-4xl border border-slate-200/80 bg-white/90 p-6 shadow-[0_20px_60px_-24px_var(--shadow-color)] backdrop-blur-xl dark:border-white/10 dark:bg-white/5 dark:shadow-2xl dark:shadow-black/30 lg:p-8">
          <div className="rounded-3xl border border-slate-200/80 bg-slate-50/95 p-6 sm:p-8 dark:border-white/10 dark:bg-slate-950/80">
            <div className="mb-8 flex flex-col gap-4 sm:flex-row sm:items-start sm:justify-between">
              <div className="space-y-2">
                <p className="text-sm font-medium uppercase tracking-[0.3em] text-amber-700 dark:text-amber-200">Product</p>
                <h1 className="text-3xl font-semibold tracking-tight text-slate-900 dark:text-white">
                  {product?.name ?? "Product details"}
                </h1>
                <p className="text-sm leading-6 text-slate-600 dark:text-slate-300">
                  Review the details before adding this item to your cart.
                </p>
              </div>
              <Link
                href="/products"
                className="inline-flex items-center gap-2 rounded-2xl border border-slate-300/80 bg-white/90 px-4 py-2 text-sm font-semibold text-slate-800 shadow-sm transition hover:bg-slate-50 dark:border-white/10 dark:bg-white/5 dark:text-slate-200 dark:hover:bg-white/10"
              >
                <span>Back to products</span>
              </Link>
            </div>

            {error ? (
              <div className="mb-6 rounded-2xl border border-rose-400/30 bg-rose-50 px-4 py-3 text-sm text-rose-700 dark:bg-rose-400/10 dark:text-rose-100">
                {error}
              </div>
            ) : null}

            {isLoading ? (
              <div className="grid gap-6 lg:grid-cols-[minmax(0,380px)_minmax(0,1fr)]">
                <div className="h-[420px] rounded-3xl border border-slate-200/80 bg-slate-100/80 animate-pulse dark:border-white/10 dark:bg-white/5" />
                <div className="space-y-4">
                  <div className="h-7 w-3/5 animate-pulse rounded bg-white/10" />
                  <div className="h-4 w-full animate-pulse rounded bg-white/5" />
                  <div className="h-4 w-4/5 animate-pulse rounded bg-white/5" />
                  <div className="h-6 w-1/3 animate-pulse rounded bg-white/10" />
                </div>
              </div>
            ) : product ? (
              <div className="grid gap-8 lg:grid-cols-[minmax(0,380px)_minmax(0,1fr)]">
                <div className="space-y-4">
                  <div className="overflow-hidden rounded-3xl border border-slate-200/80 bg-slate-100/80 dark:border-white/10 dark:bg-slate-900/60">
                    <div className="relative aspect-[4/5] w-full">
                      <img src={imageUrl} alt={product.name ?? "Product image"} className="h-full w-full object-cover" />
                    </div>
                  </div>
                  <div className="rounded-2xl border border-slate-200/80 bg-white/95 px-4 py-3 text-sm text-slate-700 shadow-sm dark:border-white/10 dark:bg-white/5 dark:text-slate-300">
                    <span className="text-slate-500 dark:text-slate-400">Supplier</span>
                    <p className="mt-1 text-base font-semibold text-slate-900 dark:text-white">{supplierLabel}</p>
                  </div>
                </div>

                <div className="space-y-6">
                  <div>
                    <p className="text-sm font-medium uppercase tracking-[0.2em] text-slate-500 dark:text-slate-400">Price</p>
                    <p className="mt-2 text-3xl font-semibold text-amber-700 dark:text-amber-300">
                      {typeof product.price === "number" ? formatCurrency(product.price) : "$0.00"}
                    </p>
                  </div>

                  <div className="rounded-2xl border border-slate-200/80 bg-white/95 px-4 py-4 shadow-sm dark:border-white/10 dark:bg-white/5">
                    <p className="text-sm font-semibold text-slate-900 dark:text-white">Availability</p>
                    <div className="mt-2 flex flex-wrap items-center gap-2">
                      <span className={isInStock ? "text-emerald-600 font-semibold dark:text-emerald-300" : "text-rose-600 font-semibold dark:text-rose-300"}>
                        {isInStock ? "In stock" : "Out of stock"}
                      </span>
                      <span className="text-sm text-slate-600 dark:text-slate-300">
                        {normalizedStock} available
                      </span>
                    </div>
                  </div>

                  <div>
                    <p className="text-sm font-medium uppercase tracking-[0.2em] text-slate-500 dark:text-slate-400">Description</p>
                    <p className="mt-2 text-sm leading-6 text-slate-700 dark:text-slate-300">
                      {product.description ?? "No description available for this product."}
                    </p>
                  </div>

                  <div className="rounded-2xl border border-slate-200/80 bg-white/95 px-4 py-4 shadow-sm dark:border-white/10 dark:bg-white/5">
                    <p className="text-sm font-semibold text-slate-900 mb-3 dark:text-white">Quantity</p>
                    <div className="flex flex-wrap items-center gap-4">
                      <div className="flex items-center gap-2 rounded-2xl border border-slate-200/80 bg-slate-100/80 px-3 py-2 dark:border-white/10 dark:bg-slate-900/50">
                        <button
                          type="button"
                          onClick={() => updateQuantity(quantity - 1)}
                          className="h-8 w-8 rounded-lg text-slate-700 transition hover:bg-slate-200 dark:text-slate-200 dark:hover:bg-white/10"
                          aria-label="Decrease quantity"
                        >
                          −
                        </button>
                        <input
                          type="number"
                          min={1}
                          max={Math.max(1, maxQuantity)}
                          value={quantity}
                          onChange={(event) => updateQuantity(Number(event.target.value))}
                          className="h-8 w-16 rounded-lg border border-slate-300/80 bg-white text-center text-sm text-slate-900 outline-none focus:border-amber-500/60 focus:ring-2 focus:ring-amber-500/20 dark:border-white/10 dark:bg-slate-950/80 dark:text-white dark:focus:border-amber-300/60 dark:focus:ring-amber-300/20"
                        />
                        <button
                          type="button"
                          onClick={() => updateQuantity(quantity + 1)}
                          className="h-8 w-8 rounded-lg text-slate-700 transition hover:bg-slate-200 dark:text-slate-200 dark:hover:bg-white/10"
                          aria-label="Increase quantity"
                        >
                          +
                        </button>
                      </div>

                      <button
                        type="button"
                        onClick={handleAddToCart}
                        disabled={!isInStock}
                        className={`inline-flex items-center justify-center rounded-2xl px-6 py-3 text-sm font-semibold transition ${isInStock
                          ? "bg-amber-400 text-slate-950 hover:bg-amber-300"
                          : "cursor-not-allowed bg-white/10 text-slate-500"
                          }`}
                      >
                        Add to cart
                      </button>
                    </div>
                    {successMessage && (
                      <div>
                        <br />
                        <div className="mb-4 rounded-2xl border border-emerald-400/30 bg-emerald-50 px-4 py-3 text-sm text-emerald-700 dark:bg-emerald-400/10 dark:text-emerald-100">
                          {successMessage}
                        </div>
                      </div>
                    )}
                  </div>
                </div>
              </div>
            ) : null}
          </div>
        </div>
      </div>
    </section>
  );
}
