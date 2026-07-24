"use client";

import Link from "next/link";
import type { StaticImageData } from "next/image";
import { useRouter } from "next/navigation";
import { useEffect, useMemo, useState } from "react";

import Categories from "@/components/Categories";
import electronicsImage from "@/images/pexels-indraprojectsofficial-18662969.jpg";
import homeImage from "@/images/pexels-peterdanthy-32760879.jpg";
import fashionImage from "@/images/pexels-rachel-claire-5531541.jpg";
import beautyImage from "@/images/pexels-deepa-nishad-1620501673-27544680.jpg";
import sportsImage from "@/images/pexels-mikhail-nilov-6740821.jpg";
import booksImage from "@/images/pexels-beyzaa-yurtkuran-279977530-16412993.jpg";
import toolsImage from "@/images/pexels-enginakyurt-1571736.jpg";
import gamingImage from "@/images/OIP1.jpg";
import officeImage from "@/images/pexels-cup-of-couple-7657384.jpg";
import groceriesImage from "@/images/pexels-macshamim-30689829.jpg";
import samsungImage from "@/images/imgi_167_galaxy-s26-ultra-reasons-buy.jpg"
import iPhoneImage from"@/images/OIP.jpg"
import tvImage from "@/images/cb179d4d-samsung-65-oled-s95f-4k-smart-tv-2025.jpg"
import Slider from "@/components/Slider";
import ThemeToggle from "@/components/theme-toggle";
import { handleLoadProducts } from "@/controllers/productController";
import type { Product } from "@/lib/product-contract";
import * as cartService from "@/services/cartService";

const productImageBasePath = "/product_images";

const heroSlides = [
  {
    id: 1,
    image: samsungImage,
    title: "Fresh arrivals",
    subtitle: "Shop curated products for your everyday needs.",
  },
  {
    id: 2,
    image: iPhoneImage,
    title: "New season picks",
    subtitle: "Explore top-rated items with fast delivery.",
  },
  {
    id: 3,
    image: tvImage,
    title: "Limited offers",
    subtitle: "Save on best-selling essentials today.",
  },
];

const categories: Array<{ id: number; name: string; image?: string | StaticImageData }> = [
  { id: 1, name: "Electronics", image: electronicsImage },
  { id: 2, name: "Home & Living", image: homeImage },
  { id: 3, name: "Fashion", image: fashionImage },
  { id: 4, name: "Beauty", image: beautyImage },
  { id: 5, name: "Sports", image: sportsImage },
  { id: 6, name: "Books", image: booksImage },
  { id: 7, name: "Tools & Machinery", image: toolsImage },
  { id: 8, name: "Gaming, Movies & Music", image: gamingImage },
  { id: 9, name: "Stationary & Office", image: officeImage },
  { id: 10, name: "Groceries & Household", image: groceriesImage },
];

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

export default function ProductsPage() {
  const router = useRouter();
  const [products, setProducts] = useState<Product[]>([]);
  const [search, setSearch] = useState("");
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [cartMessage, setCartMessage] = useState<string | null>(null);
  const [visibleCount, setVisibleCount] = useState(4);

  useEffect(() => {
    let isMounted = true;

    async function loadProducts() {
      setIsLoading(true);
      setError(null);

      try {
        const result = await handleLoadProducts();

        if (!isMounted) {
          return;
        }

        if (result.ok) {
          setProducts(result.data ?? []);
        } else {
          setProducts([]);
          setError(result.error ?? "Failed to load products.");
        }
      } catch (fetchError) {
        if (isMounted) {
          setProducts([]);
          setError(fetchError instanceof Error ? fetchError.message : "Failed to load products.");
        }
      } finally {
        if (isMounted) {
          setIsLoading(false);
        }
      }
    }

    loadProducts();

    return () => {
      isMounted = false;
    };
  }, []);

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

  useEffect(() => {
    setVisibleCount(4);
  }, [search, products.length]);

  const visibleProducts = filteredProducts.slice(0, visibleCount);

  const handleAddToCart = (product: Product) => {
    const priceCents = typeof product.price === "number" ? product.price : 0;
    const rawId = product.id ?? product.product_id ?? product["product_id"];
    const productId = typeof rawId === "string" || typeof rawId === "number" ? String(rawId) : "";
    const imageUrl = resolveImageUrl(product);

    const newItem = {
      id: productId || `${product.name ?? "product"}-${Date.now()}`,
      title: product.name ?? "Untitled product",
      price: priceCents / 100,
      image: imageUrl,
      description: product.description,
      qty: 1,
    };

    const current = cartService.loadCartFromStorage();
    const existingIndex = current.findIndex((item) => String(item.id) === String(productId));
    let updated = [] as typeof current;

    if (existingIndex >= 0) {
      updated = current.map((item, index) =>
        index === existingIndex ? { ...item, qty: item.qty + 1 } : item,
      );
    } else {
      updated = [...current, newItem];
    }

    cartService.saveCartToStorage(updated);
    setCartMessage(`${product.name ?? "Item"} added to cart.`);
  };

  return (
    <div className="relative min-h-screen overflow-hidden bg-[radial-gradient(circle_at_top_left,var(--hero-accent),transparent_30%),radial-gradient(circle_at_bottom_right,var(--hero-accent-2),transparent_28%),linear-gradient(180deg,var(--hero-bg-top)_0%,var(--hero-bg-bottom)_100%)] px-6 py-10 text-slate-900 sm:px-8 lg:px-10 dark:text-slate-100">
      <div className="absolute inset-0 bg-[linear-gradient(var(--grid-line)_1px,transparent_1px),linear-gradient(90deg,var(--grid-line)_1px,transparent_1px)] bg-size-[28px_28px] opacity-20" />
      <header className="rounded-4xl border border-slate-200/70 bg-white/70 p-6 shadow-[0_20px_60px_-24px_var(--shadow-color)] backdrop-blur-xl dark:border-white/10 dark:bg-white/10 dark:shadow-2xl dark:shadow-black/20">
          <div className="flex flex-col gap-4 xl:flex-row xl:items-center xl:gap-6">
            <div className="flex min-w-[170px] items-center justify-start">
              <div className="inline-flex h-11 w-40 items-center justify-center rounded-xl border border-white/15 bg-slate-950/80 text-sm font-semibold uppercase tracking-[0.2em] text-amber-100">
                Logo
              </div>
            </div>

            <div className="flex flex-1 justify-center">
              <label className="relative w-full max-w-2xl">
                <input
                  type="text"
                  value={search}
                  onChange={(event) => setSearch(event.target.value)}
                  placeholder="Search products, categories..."
                  className="w-full rounded-2xl border border-white/10 bg-white/5 py-3 pl-4 pr-12 text-sm text-slate-100 outline-none transition placeholder:text-slate-500 focus:border-amber-300/60 focus:bg-white/8 focus:ring-2 focus:ring-amber-300/20"
                />
                <span className="pointer-events-none absolute inset-y-0 right-4 flex items-center text-slate-400">
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
                className="inline-flex items-center gap-2 rounded-2xl border border-white/10 bg-white/5 px-3 py-2 text-sm font-medium text-slate-200 transition hover:bg-white/10"
              >
                <svg viewBox="0 0 24 24" className="h-5 w-5" fill="none" stroke="currentColor" strokeWidth="2" aria-hidden="true">
                  <circle cx="12" cy="8" r="4" />
                  <path d="M4 20a8 8 0 0 1 16 0" />
                </svg>
                <span>Account</span>
              </Link>
              <button
                type="button"
                className="inline-flex items-center gap-2 rounded-2xl border border-white/10 bg-white/5 px-3 py-2 text-sm font-medium text-slate-200 transition hover:bg-white/10"
              >
                <svg viewBox="0 0 24 24" className="h-5 w-5" fill="none" stroke="currentColor" strokeWidth="2" aria-hidden="true">
                  <circle cx="9" cy="20" r="1.5" />
                  <circle cx="17" cy="20" r="1.5" />
                  <path d="M3 4h2l2.4 10.2a1 1 0 0 0 1 .8H18a1 1 0 0 0 1-.8L21 7H7" />
                </svg>
                <span>Cart</span>
              </button>
              <ThemeToggle />
            </div>
          </div>
        </header>
        <br/>
      <div className="relative mx-auto flex min-h-[calc(100vh-5rem)] w-full max-w-[1600px] flex-col gap-6">
        <div className="rounded-4xl border border-slate-200/80 bg-white/90 p-6 shadow-[0_20px_60px_-24px_var(--shadow-color)] backdrop-blur-xl dark:border-white/10 dark:bg-white/5 dark:shadow-2xl dark:shadow-black/30">
          <div className="rounded-3xl border border-slate-200/80 bg-slate-50/95 p-6 sm:p-8 dark:border-white/10 dark:bg-slate-950/80">
            <div className="mb-6 flex flex-col gap-4 lg:flex-row lg:items-center lg:justify-between">
              <div>
                <p className="text-sm font-medium uppercase tracking-[0.3em] text-amber-700 dark:text-amber-200">Featured</p>
                <h1 className="mt-2 text-3xl font-semibold tracking-tight text-slate-900 dark:text-white">Top picks of the week</h1>
              </div>
              <div className="flex flex-wrap items-center gap-3">
                <Link href="/products" className="rounded-2xl border border-slate-300/80 bg-white/90 px-4 py-2 text-sm font-semibold text-slate-900 shadow-sm transition hover:bg-slate-50 dark:border-white/10 dark:bg-white/5 dark:text-slate-200 dark:hover:bg-white/10">
                  Shop all products
                </Link>
              </div>
            </div>

            <div className="overflow-hidden rounded-3xl border border-slate-200/80 bg-slate-100/80 dark:border-white/10 dark:bg-slate-900/60">
              <Slider slides={heroSlides} />
            </div>

            <div className="mt-8 rounded-3xl border border-slate-200/80 bg-white/95 p-5 dark:border-white/10 dark:bg-white/5">
              <div className="mb-5 flex items-center justify-between gap-4">
                <div>
                  <p className="text-xs font-semibold uppercase tracking-[0.2em] text-amber-700 dark:text-amber-200">Categories</p>
                  <h2 className="text-xl font-semibold text-slate-900 dark:text-white">Browse by category</h2>
                </div>
                <Link href="/products" className="text-sm font-semibold text-slate-700 transition hover:text-slate-900 dark:text-slate-300 dark:hover:text-white">
                  View all
                </Link>
              </div>
              <Categories categories={categories} />
            </div>
          </div>
        </div>
      </div>
      <br/>
      <main className="min-w-0 flex-1 rounded-4xl border border-slate-200/80 bg-white/90 p-6 shadow-[0_20px_60px_-24px_var(--shadow-color)] backdrop-blur-xl dark:border-white/10 dark:bg-white/5 dark:shadow-2xl dark:shadow-black/30 lg:p-8">
        <div className="rounded-3xl border border-slate-200/80 bg-slate-50/95 p-6 sm:p-8 dark:border-white/10 dark:bg-slate-950/80">
          <div className="mb-8 flex flex-col gap-3 sm:flex-row sm:items-end sm:justify-between">
            <div className="space-y-2">
              <p className="text-sm font-medium uppercase tracking-[0.3em] text-amber-700 dark:text-amber-200">Browse</p>
              <h1 className="text-3xl font-semibold tracking-tight text-slate-900 dark:text-white">Products</h1>
              <p className="text-sm leading-6 text-slate-600 dark:text-slate-300">Discover our collection of available products.</p>
            </div>

            <div className="flex flex-wrap items-center gap-3">
              {cartMessage ? (
                <span className="rounded-full border border-emerald-400/40 bg-emerald-50 px-3 py-1 text-xs font-semibold uppercase tracking-[0.12em] text-emerald-700 dark:border-emerald-400/30 dark:bg-emerald-400/10 dark:text-emerald-200">
                  {cartMessage}
                </span>
              ) : null}
              {filteredProducts.length > 4 ? (
                <button
                  type="button"
                  onClick={() =>
                    setVisibleCount((current) =>
                      current >= filteredProducts.length ? 4 : Math.min(filteredProducts.length, current + 4),
                    )
                  }
                  className="inline-flex w-fit rounded-full border border-amber-500/40 bg-amber-500/90 px-3 py-1 text-xs font-semibold uppercase tracking-[0.12em] text-slate-950 shadow-sm transition hover:bg-amber-400"
                >
                  {visibleCount >= filteredProducts.length ? "Show less" : "View more"}
                </button>
              ) : null}
            </div>
          </div>

          {error ? (
            <div className="mb-4 rounded-2xl border border-rose-400/30 bg-rose-50 px-4 py-3 text-sm text-rose-700 dark:bg-rose-400/10 dark:text-rose-100">{error}</div>
          ) : null}

          {isLoading ? (
            <div className="grid gap-4 sm:grid-cols-2 xl:grid-cols-4">
              {Array.from({ length: 6 }).map((_, index) => (
                <div key={index} className="rounded-3xl border border-slate-200/80 bg-white/80 p-6 shadow-sm dark:border-white/10 dark:bg-white/5">
                  <div className="h-6 w-3/5 animate-pulse rounded bg-white/10" />
                  <div className="mt-3 h-4 w-full animate-pulse rounded bg-white/5" />
                  <div className="mt-2 h-4 w-4/5 animate-pulse rounded bg-white/5" />
                  <div className="mt-5 h-6 w-1/3 animate-pulse rounded bg-white/10" />
                </div>
              ))}
            </div>
          ) : filteredProducts.length === 0 ? (
            <div className="rounded-2xl border border-slate-200/80 bg-white/90 p-8 text-sm text-slate-700 shadow-sm dark:border-white/10 dark:bg-white/5 dark:text-slate-300">No products found.</div>
          ) : (
            <div className="grid gap-4 sm:grid-cols-2 xl:grid-cols-4">
              {visibleProducts.map((product, index) => {
                const rawId = product.id ?? product.product_id ?? product["product_id"] ?? index;
                const productId = typeof rawId === "string" || typeof rawId === "number" ? String(rawId) : String(index);
                const imageUrl = resolveImageUrl(product);

                return (
                  <div
                    key={productId}
                    onClick={() => router.push(`/products/${productId}`)}
                    className="group cursor-pointer overflow-hidden rounded-3xl border border-slate-200/80 bg-white/95 p-3 shadow-[0_16px_35px_-20px_var(--shadow-color)] transition hover:-translate-y-1 hover:border-amber-500/40 hover:bg-white dark:border-white/10 dark:bg-white/5 dark:shadow-2xl dark:shadow-black/30 dark:hover:bg-white/8 dark:hover:border-amber-300/40 backdrop-blur-xl"
                  >
                    <article>
                      <div className="relative mb-3 overflow-hidden rounded-2xl border border-slate-200/80 bg-slate-100/80 dark:border-white/10 dark:bg-slate-900/60">
                        <div className="relative h-28 w-full">
                          <img src={imageUrl} alt={product.name ?? "Product image"} className="h-full w-full object-cover" />
                        </div>
                        <button
                          type="button"
                          onClick={(event) => {
                            event.preventDefault();
                            event.stopPropagation();
                            handleAddToCart(product);
                          }}
                          className="absolute right-2 top-2 inline-flex items-center justify-center rounded-full border border-amber-500/40 bg-amber-500/95 p-2 text-slate-950 shadow-lg transition hover:bg-amber-400"
                          aria-label={`Add ${product.name ?? "product"} to cart`}
                        >
                          <svg viewBox="0 0 24 24" className="h-4 w-4" fill="none" stroke="currentColor" strokeWidth="2" aria-hidden="true">
                            <circle cx="9" cy="20" r="1.5" />
                            <circle cx="17" cy="20" r="1.5" />
                            <path d="M3 4h2l2.4 10.2a1 1 0 0 0 1 .8H18a1 1 0 0 0 1-.8L21 7H7" />
                          </svg>
                        </button>
                      </div>
                      <h2 className="text-base font-semibold text-slate-900 transition group-hover:text-amber-700 dark:text-white dark:group-hover:text-amber-200">
                        {product.name ?? "Unnamed product"}
                      </h2>
                      {product.description ? (
                        <p className="mt-2 text-sm leading-5 text-slate-700 dark:text-slate-300">{product.description}</p>
                      ) : (
                        <p className="mt-2 text-sm leading-6 text-slate-500 dark:text-slate-500">No description available.</p>
                      )}
                      {typeof product.price === "number" ? (
                        <p className="mt-2 text-sm font-semibold text-amber-700 dark:text-amber-300">{formatCurrency(product.price)}</p>
                      ) : null}
                    </article>
                  </div>
                );
              })}
            </div>
          )}
        </div>
      </main>
    </div>
  );
}
