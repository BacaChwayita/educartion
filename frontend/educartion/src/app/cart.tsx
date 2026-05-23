"use client";

import Link from "next/link";
import React, { useEffect, useState } from "react";

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
    <main style={{ width: "100%", padding: "24px 20px", display: "grid", justifyItems: "center" }}>
      <div style={{ width: "100%", maxWidth: 1200, display: "grid", gap: 24 }}>
        <header style={{ display: "flex", justifyContent: "space-between", alignItems: "center", gap: 16, flexWrap: "wrap" }}>
          <div>
            <h1 style={{ margin: 0, fontSize: 32 }}>Shopping Cart</h1>
            <p style={{ margin: "8px 0 0", color: "#555" }}>Review your items before you checkout.</p>
          </div>

          <div style={{ display: "flex", gap: 12, flexWrap: "wrap" }}>
            <Link href="/" style={{ display: "inline-flex", alignItems: "center", justifyContent: "center", padding: "12px 18px", borderRadius: 8, border: "1px solid #ccc", background: "white", color: "#111", textDecoration: "none", minWidth: 170 }}>
              Continue shopping
            </Link>
            <Link href="/payment" style={{ display: "inline-flex", alignItems: "center", justifyContent: "center", padding: "12px 18px", borderRadius: 8, border: "none", background: "#0070f3", color: "white", textDecoration: "none", minWidth: 170 }}>
              Proceed to checkout
            </Link>
          </div>
        </header>

        {items.length === 0 ? (
          <section style={{ padding: 28, border: "1px solid #e5e7eb", borderRadius: 16, background: "#fafafa", textAlign: "center" }}>
            <p style={{ margin: 0, fontSize: 18, color: "#333" }}>Your cart is empty.</p>
            <p style={{ margin: "12px 0 0", color: "#666" }}>Add products from the shop and return here to complete your order.</p>
            <Link href="/" style={{ marginTop: 18, display: "inline-flex", padding: "12px 20px", borderRadius: 8, border: "none", background: "#0070f3", color: "white", textDecoration: "none" }}>
              Browse products
            </Link>
          </section>
        ) : (
          <section style={{ display: "grid", gap: 24, gridTemplateColumns: "1.7fr 0.9fr" }}>
            <div style={{ display: "grid", gap: 16 }}>
              <div style={{ padding: 16, borderRadius: 16, border: "1px solid #e5e7eb", background: "white" }}>
                <div style={{ display: "grid", gridTemplateColumns: "2fr 1fr 1fr 1fr auto", gap: 12, padding: "12px 0", fontSize: 14, fontWeight: 700, color: "#555" }}>
                  <div>Product</div>
                  <div style={{ textAlign: "right" }}>Price</div>
                  <div style={{ textAlign: "center" }}>Quantity</div>
                  <div style={{ textAlign: "right" }}>Total</div>
                  <div />
                </div>
              </div>

              {items.map((item) => (
                <article key={item.id} style={{ display: "grid", gridTemplateColumns: "2fr 1fr 1fr 1fr auto", gap: 12, padding: 18, borderRadius: 16, border: "1px solid #e5e7eb", background: "white", alignItems: "center" }}>
                  <div style={{ display: "flex", gap: 14, alignItems: "center" }}>
                    <div style={{ width: 100, height: 100, borderRadius: 16, overflow: "hidden", background: "#f3f4f6", flexShrink: 0 }}>
                      <img src={item.image || "/images/product-placeholder.png"} alt={item.title} style={{ width: "100%", height: "100%", objectFit: "cover" }} />
                    </div>
                    <div>
                      <div style={{ fontWeight: 700, fontSize: 16 }}>{item.title}</div>
                      {item.description && <div style={{ marginTop: 4, color: "#6b7280", fontSize: 13 }}>{item.description}</div>}
                    </div>
                  </div>

                  <div style={{ textAlign: "right", fontWeight: 600 }}>{formatMoney(item.price)}</div>

                  <div style={{ display: "inline-flex", justifyContent: "center", alignItems: "center", gap: 8, background: "#f8fafc", borderRadius: 999, padding: "6px 10px" }}>
                    <button
                      type="button"
                      onClick={() => updateItemQty(item.id, item.qty - 1)}
                      style={{ width: 30, height: 30, borderRadius: 999, border: "1px solid #d1d5db", background: "white", cursor: "pointer", fontWeight: 700 }}
                      aria-label={`Decrease quantity for ${item.title}`}
                    >
                      –
                    </button>
                    <span style={{ minWidth: 24, textAlign: "center", fontWeight: 700 }}>{item.qty}</span>
                    <button
                      type="button"
                      onClick={() => updateItemQty(item.id, item.qty + 1)}
                      style={{ width: 30, height: 30, borderRadius: 999, border: "1px solid #d1d5db", background: "white", cursor: "pointer", fontWeight: 700 }}
                      aria-label={`Increase quantity for ${item.title}`}
                    >
                      +
                    </button>
                  </div>

                  <div style={{ textAlign: "right", fontWeight: 700 }}>{formatMoney(item.price * item.qty)}</div>

                  <button
                    type="button"
                    onClick={() => removeItem(item.id)}
                    style={{ width: 40, height: 40, borderRadius: 12, border: "1px solid #e5e7eb", background: "#f8fafc", cursor: "pointer", color: "#ef4444" }}
                    aria-label={`Remove ${item.title} from cart`}
                  >
                    🗑️
                  </button>
                </article>
              ))}
            </div>

            <aside style={{ display: "grid", gap: 16, alignSelf: "start" }}>
              <div style={{ padding: 20, borderRadius: 18, border: "1px solid #e5e7eb", background: "white" }}>
                <h2 style={{ margin: 0, marginBottom: 14, fontSize: 20 }}>Order summary</h2>
                <div style={{ display: "grid", gap: 12 }}>
                  <div style={{ display: "flex", justifyContent: "space-between", color: "#4b5563" }}>
                    <span>Subtotal</span>
                    <span>{formatMoney(subtotal)}</span>
                  </div>
                  <div style={{ display: "flex", justifyContent: "space-between", color: "#4b5563" }}>
                    <span>Delivery fee</span>
                    <span>{formatMoney(delivery)}</span>
                  </div>
                  <div style={{ height: 1, background: "#e5e7eb", margin: "8px 0" }} />
                  <div style={{ display: "flex", justifyContent: "space-between", fontWeight: 700, fontSize: 18 }}>
                    <span>Total</span>
                    <span>{formatMoney(total)}</span>
                  </div>
                </div>
              </div>
              <div style={{ padding: 20, borderRadius: 18, border: "1px solid #e5e7eb", background: "#f9fafb" }}>
                <p style={{ margin: 0, color: "#374151", fontSize: 14 }}>Need help? Your order summary updates automatically as you change quantity or remove items.</p>
              </div>
            </aside>
          </section>
        )}
      </div>
    </main>
  );
}
