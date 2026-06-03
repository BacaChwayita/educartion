"use client";

import Link from "next/link";
import { useEffect, useMemo, useState } from "react";
import { type OrderSummary } from "@/lib/orders-contract";

const sidebarItems = [
  "Dashboard",
  "My Orders",
  "Profile",
  "Addresses",
  "Payment Methods",
  "Wishlist",
  "Logout",
];

function formatCurrency(cents: number): string {
  return new Intl.NumberFormat("en-US", {
    style: "currency",
    currency: "USD",
  }).format(cents / 100);
}

function formatDate(value: string): string {
  const date = new Date(value);

  if (Number.isNaN(date.getTime())) {
    return "Unknown date";
  }

  return new Intl.DateTimeFormat("en-US", {
    year: "numeric",
    month: "short",
    day: "numeric",
  }).format(date);
}

export default function OrdersPage() {
  const [orders, setOrders] = useState<OrderSummary[]>([]);
  const [search, setSearch] = useState("");
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    let isMounted = true;

    async function loadOrders() {
      setIsLoading(true);
      setError(null);

      try {
        const response = await fetch("/api/orders/", { cache: "no-store" });

        if (!response.ok) {
          const body = (await response.json().catch(() => null)) as { error?: string } | null;
          throw new Error(body?.error ?? "Failed to load orders.");
        }

        const data = (await response.json()) as OrderSummary[];

        if (isMounted) {
          setOrders(Array.isArray(data) ? data : []);
        }
      } catch (fetchError) {
        if (isMounted) {
          setError(fetchError instanceof Error ? fetchError.message : "Failed to load orders.");
          setOrders([]);
        }
      } finally {
        if (isMounted) {
          setIsLoading(false);
        }
      }
    }

    loadOrders();

    return () => {
      isMounted = false;
    };
  }, []);

  const filteredOrders = useMemo(() => {
    const query = search.trim().toLowerCase();

    if (!query) {
      return orders;
    }

    return orders.filter((order) => {
      return (
        order.order_number.toLowerCase().includes(query) ||
        order.status.toLowerCase().includes(query) ||
        String(order.order_id).includes(query)
      );
    });
  }, [orders, search]);

  return (
    <div className="relative min-h-screen overflow-hidden bg-[radial-gradient(circle_at_top_left,rgba(251,191,36,0.2),transparent_30%),radial-gradient(circle_at_bottom_right,rgba(14,165,233,0.18),transparent_28%),linear-gradient(180deg,#060816_0%,#0b1020_100%)] px-6 py-10 text-slate-100 sm:px-8 lg:px-10">
      <div className="absolute inset-0 bg-[linear-gradient(rgba(255,255,255,0.03)_1px,transparent_1px),linear-gradient(90deg,rgba(255,255,255,0.03)_1px,transparent_1px)] bg-size-[28px_28px] opacity-20" />

      <div className="relative mx-auto flex min-h-[calc(100vh-5rem)] w-full max-w-[1600px] flex-col gap-6">
        <header className="rounded-4xl border border-white/10 bg-white/5 p-6 shadow-2xl shadow-black/30 backdrop-blur-xl">
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
                  placeholder="Search orders"
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
            </div>
          </div>
        </header>

        <div className="flex flex-1 flex-col gap-6 xl:flex-row">

          <main className="min-w-0 flex-1 rounded-4xl border border-white/10 bg-white/5 p-6 shadow-2xl shadow-black/30 backdrop-blur-xl lg:p-8">
            <div className="rounded-3xl border border-white/10 bg-slate-950/80 p-6 sm:p-8">
              <div className="mb-8 flex flex-col gap-3 sm:flex-row sm:items-end sm:justify-between">
                <div className="space-y-2">
                  <p className="text-sm font-medium uppercase tracking-[0.3em] text-amber-200">Orders</p>
                  <h1 className="text-3xl font-semibold tracking-tight text-white">My Orders</h1>
                </div>

                <span className="inline-flex w-fit rounded-full border border-white/10 bg-white/5 px-3 py-1 text-xs font-semibold uppercase tracking-[0.12em] text-slate-300">
                  {filteredOrders.length} Results
                </span>
              </div>

              {error ? (
                <div className="mb-4 rounded-2xl border border-rose-400/30 bg-rose-400/10 px-4 py-3 text-sm text-rose-100">
                  {error}
                </div>
              ) : null}

              {isLoading ? (
                <div className="space-y-3">
                  {Array.from({ length: 3 }).map((_, index) => (
                    <div key={index} className="rounded-2xl border border-white/10 bg-white/5 p-5">
                      <div className="h-5 w-1/3 animate-pulse rounded bg-white/10" />
                      <div className="mt-3 h-4 w-2/3 animate-pulse rounded bg-white/8" />
                    </div>
                  ))}
                </div>
              ) : filteredOrders.length === 0 ? (
                <div className="rounded-2xl border border-white/10 bg-white/5 p-8 text-sm text-slate-300">
                  No orders found.
                </div>
              ) : (
                <div className="space-y-3">
                  {filteredOrders.map((order) => {
                    return (
                      <Link
                        key={order.order_id}
                        href={`/orders/${order.order_id}`}
                        aria-label={`View order ${order.order_number}`}
                        className="block rounded-2xl border border-white/10 bg-white/5 p-5 transition hover:border-amber-300/40 hover:bg-white/8"
                      >
                        <article>
                          <div className="flex flex-wrap items-start justify-between gap-3">
                            <div>
                              <p className="text-sm font-semibold uppercase tracking-[0.18em] text-amber-200">
                                Order #{order.order_number}
                              </p>
                              <h2 className="mt-1 text-lg font-semibold text-white">Order ID: {order.order_id}</h2>
                              <p className="mt-1 text-sm text-slate-300">Placed on {formatDate(order.placed_at)}</p>
                            </div>
                            <span className="rounded-full border border-amber-300/30 bg-amber-400/10 px-3 py-1 text-xs font-semibold uppercase tracking-[0.12em] text-amber-100">
                              {order.status}
                            </span>
                          </div>

                          <div className="mt-4 grid gap-2 text-sm text-slate-200 sm:grid-cols-3">
                            <p>
                              <span className="font-semibold text-white">Subtotal:</span>{" "}
                              {formatCurrency(order.subtotal_amount)}
                            </p>
                            <p>
                              <span className="font-semibold text-white">Discount:</span>{" "}
                              {formatCurrency(order.discount_amount)}
                            </p>
                            <p>
                              <span className="font-semibold text-white">Total:</span> {formatCurrency(order.total_amount)}
                            </p>
                          </div>
                        </article>
                      </Link>
                    );
                  })}
                </div>
              )}
            </div>
          </main>
        </div>
      </div>
    </div>
  );
}
