"use client";

import Link from "next/link";
import React, { useEffect, useState } from "react";

import ThemeToggle from "@/components/theme-toggle";

type CartItem = {
  id: string | number;
  title: string;
  price: number;
  image?: string;
  description?: string;
  qty: number;
};

const DELIVERY_FEE = 5.0;

const formatMoney = (value: number) => `$${value.toFixed(2)}`;

const normalizeCart = (cart: any[]): CartItem[] => {
  const merged: Record<string, CartItem> = {};

  cart.forEach((item) => {
    const id = String(item.id ?? item.title ?? Math.random());
    const qty = Math.max(1, Number(item.qty || 1));
    const price = Number(item.price || 0);

    if (!merged[id]) {
      merged[id] = {
        id,
        title: String(item.title || "Untitled product"),
        price,
        image: item.image || "/images/product-placeholder.png",
        description: item.description,
        qty,
      };
    } else {
      merged[id].qty += qty;
    }
  });

  return Object.values(merged);
};

export default function CartPage() {
  const [items, setItems] = useState<CartItem[]>([]);
  const [isLoaded, setIsLoaded] = useState(false);

  useEffect(() => {
    if (typeof window === "undefined") return;
    const stored = window.localStorage.getItem("cart") || "[]";
    let parsed: any[] = [];

    try {
      parsed = JSON.parse(stored);
    } catch {
      parsed = [];
    }

    setItems(normalizeCart(parsed));
    setIsLoaded(true);
  }, []);

  useEffect(() => {
    if (!isLoaded) return;
    if (typeof window === "undefined") return;

    window.localStorage.setItem("cart", JSON.stringify(items));
  }, [items, isLoaded]);

  const updateItemQty = (id: string | number, qty: number) => {
    setItems((current) =>
      current.map((item) =>
        String(item.id) === String(id) ? { ...item, qty: Math.max(1, qty) } : item
      )
    );
  };

  const removeItem = (id: string | number) => {
    setItems((current) => current.filter((item) => String(item.id) !== String(id)));
  };

  const subtotal = items.reduce((sum, item) => sum + item.price * item.qty, 0);
  const delivery = items.length > 0 ? DELIVERY_FEE : 0;
  const total = subtotal + delivery;

  return (
    <section className="relative min-h-screen overflow-hidden bg-[radial-gradient(circle_at_top_left,var(--hero-accent),transparent_30%),radial-gradient(circle_at_bottom_right,var(--hero-accent-2),transparent_28%),linear-gradient(180deg,var(--hero-bg-top)_0%,var(--hero-bg-bottom)_100%)] px-6 py-10 text-slate-100 sm:px-8 lg:px-10">
      <div className="absolute inset-0 bg-[linear-gradient(var(--grid-line)_1px,transparent_1px),linear-gradient(90deg,var(--grid-line)_1px,transparent_1px)] bg-size-[28px_28px] opacity-20" />
      <div className="relative mx-auto w-full max-w-5xl">
        <div className="w-full rounded-4xl border border-white/10 bg-white/5 p-6 shadow-2xl shadow-black/30 backdrop-blur-xl lg:p-8">
          <div className="rounded-3xl border border-white/10 bg-slate-950/80 p-6 sm:p-8">
            <header className="mb-8 flex flex-wrap items-center justify-between gap-4 rounded-2xl border border-white/10 bg-white/5 p-4 sm:p-5">
              <div className="space-y-2">
                <p className="text-sm font-medium uppercase tracking-[0.3em] text-amber-200">
                  Shopping
                </p>
                <h2 className="text-3xl font-semibold tracking-tight text-white">Shopping Cart</h2>
                <p className="text-sm leading-6 text-slate-300">Review your items before you checkout.</p>
              </div>
              <div className="flex items-center gap-2">
                <Link
                  href="/products"
                  className="inline-flex items-center justify-center rounded-2xl border border-white/15 px-4 py-2 text-sm font-semibold text-white transition hover:bg-white/10"
                >
                  Continue shopping
                </Link>
                <Link
                  href="/cart"
                  className="inline-flex items-center gap-2 rounded-2xl border border-white/10 bg-white/5 px-3 py-2 text-sm font-medium text-slate-200 transition hover:bg-white/10"
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
            </header>

            {items.length === 0 ? (
              <div className="rounded-2xl border border-white/10 bg-white/5 px-6 py-8 text-center">
                <p className="text-base text-slate-300 mb-3">Your cart is empty.</p>
                <p className="text-sm text-slate-400 mb-6">Add products from the shop and return here to complete your order.</p>
                <Link
                  href="/products"
                  className="inline-flex items-center justify-center rounded-2xl bg-amber-400 px-6 py-3 text-sm font-semibold text-slate-950 transition hover:bg-amber-300"
                >
                  Browse products
                </Link>
              </div>
            ) : (
              <div className="grid gap-6 lg:grid-cols-3">
                <div className="lg:col-span-2 space-y-4">
                  {items.map((item) => (
                    <div
                      key={item.id}
                      className="rounded-2xl border border-white/10 bg-white/5 p-4 flex gap-4 items-start"
                    >
                      <div className="w-24 h-24 rounded-lg overflow-hidden flex-shrink-0 bg-slate-900/50">
                        <img
                          src={item.image || "/images/product-placeholder.png"}
                          alt={item.title}
                          className="w-full h-full object-cover"
                        />
                      </div>
                      <div className="flex-1 min-w-0">
                        <h3 className="text-sm font-semibold text-white truncate">{item.title}</h3>
                        {item.description && (
                          <p className="text-xs text-slate-400 line-clamp-2 mt-1">{item.description}</p>
                        )}
                        <p className="text-amber-400 font-semibold mt-2">{formatMoney(item.price)}</p>
                      </div>
                      <div className="flex items-center gap-2 bg-slate-900/50 rounded-lg p-1">
                        <button
                          type="button"
                          onClick={() => updateItemQty(item.id, item.qty - 1)}
                          className="w-6 h-6 rounded text-slate-200 hover:bg-white/10 transition"
                        >
                          −
                        </button>
                        <span className="w-6 text-center text-xs font-semibold text-slate-200">{item.qty}</span>
                        <button
                          type="button"
                          onClick={() => updateItemQty(item.id, item.qty + 1)}
                          className="w-6 h-6 rounded text-slate-200 hover:bg-white/10 transition"
                        >
                          +
                        </button>
                      </div>
                      <p className="text-amber-400 font-semibold text-sm text-right min-w-[60px]">
                        {formatMoney(item.price * item.qty)}
                      </p>
                      <button
                        type="button"
                        onClick={() => removeItem(item.id)}
                        className="text-rose-400 hover:text-rose-300 transition text-sm font-semibold"
                      >
                        Remove
                      </button>
                    </div>
                  ))}
                </div>

                <div className="space-y-4">
                  <div className="rounded-2xl border border-white/10 bg-white/5 px-4 py-5">
                    <h3 className="text-sm font-semibold text-white mb-4">Order Summary</h3>
                    <div className="space-y-3 text-sm">
                      <div className="flex justify-between text-slate-300">
                        <span>Subtotal</span>
                        <span>{formatMoney(subtotal)}</span>
                      </div>
                      <div className="flex justify-between text-slate-300">
                        <span>Delivery fee</span>
                        <span>{formatMoney(delivery)}</span>
                      </div>
                      <div className="border-t border-white/10 pt-3 flex justify-between font-semibold text-amber-400">
                        <span>Total</span>
                        <span>{formatMoney(total)}</span>
                      </div>
                    </div>
                  </div>
                  <Link
                    href="/payment"
                    className="w-full flex items-center justify-center rounded-2xl bg-amber-400 px-5 py-3.5 text-sm font-semibold text-slate-950 transition hover:bg-amber-300"
                  >
                    Proceed to checkout
                  </Link>
                  <Link
                    href="/"
                    className="w-full flex items-center justify-center rounded-2xl border border-white/15 px-5 py-3.5 text-sm font-semibold text-white transition hover:bg-white/10"
                  >
                    Continue shopping
                  </Link>
                </div>
              </div>
            )}
          </div>
        </div>
      </div>
    </section>
  );
}
