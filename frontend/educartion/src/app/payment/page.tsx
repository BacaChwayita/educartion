"use client";

import React, { useEffect, useState } from "react";
import Link from "next/link";

type CartItem = { id: string | number; title: string; price: number; image?: string; description?: string; qty: number };

const DELIVERY_FEE = 5.0;
const formatMoney = (v: number) => `$${v.toFixed(2)}`;

export default function PaymentPage() {
  const [cart, setCart] = useState<CartItem[]>([]);
  const [loaded, setLoaded] = useState(false);

  // Delivery address fields
  const [name, setName] = useState("");
  const [address, setAddress] = useState("");
  const [city, setCity] = useState("");
  const [postal, setPostal] = useState("");

  // Payment fields
  const [cardName, setCardName] = useState("");
  const [cardNumber, setCardNumber] = useState("");
  const [expiry, setExpiry] = useState("");
  const [cvc, setCvc] = useState("");

  useEffect(() => {
    if (typeof window === "undefined") return;
    const stored = window.localStorage.getItem("cart") || "[]";
    try {
      const parsed = JSON.parse(stored || "[]");
      setCart(parsed.map((it: any) => ({ id: it.id ?? it.title, title: it.title ?? "Untitled", price: Number(it.price || 0), image: it.image, description: it.description, qty: Number(it.qty || 1) })));
    } catch {
      setCart([]);
    }
    setLoaded(true);
  }, []);

  const subtotal = cart.reduce((s, it) => s + it.price * it.qty, 0);
  const delivery = cart.length > 0 ? DELIVERY_FEE : 0;
  const total = subtotal + delivery;

  const confirmOrder = () => {
    // Basic client-side simulation: validate minimal fields
    if (!name || !address || !cardName || cardNumber.length < 12) {
      alert("Please complete your delivery and payment details.");
      return;
    }

    // Simulate order submission
    window.localStorage.removeItem("cart");
    setCart([]);
    alert("Order confirmed — thank you!");
  };

  return (
    <main style={{ width: "100%", padding: 24, display: "grid", justifyItems: "center" }}>
      <div style={{ width: "100%", maxWidth: 900, display: "grid", gap: 20 }}>
        <h1 style={{ margin: 0 }}>Checkout</h1>

        <section style={{ display: "grid", gap: 18, gridTemplateColumns: "1fr 360px" }}>
          <div style={{ display: "grid", gap: 12 }}>
            <div style={{ padding: 16, border: "1px solid #e5e7eb", borderRadius: 12, background: "white" }}>
              <h2 style={{ margin: "0 0 8px" }}>Delivery address</h2>
              <div style={{ display: "grid", gap: 8 }}>
                <input placeholder="Full name" value={name} onChange={(e) => setName(e.target.value)} style={{ padding: 10, borderRadius: 8, border: "1px solid #d1d5db" }} />
                <input placeholder="Street address" value={address} onChange={(e) => setAddress(e.target.value)} style={{ padding: 10, borderRadius: 8, border: "1px solid #d1d5db" }} />
                <div style={{ display: "flex", gap: 8 }}>
                  <input placeholder="City" value={city} onChange={(e) => setCity(e.target.value)} style={{ padding: 10, borderRadius: 8, border: "1px solid #d1d5db", flex: 1 }} />
                  <input placeholder="Postal code" value={postal} onChange={(e) => setPostal(e.target.value)} style={{ padding: 10, borderRadius: 8, border: "1px solid #d1d5db", width: 120 }} />
                </div>
              </div>
            </div>

            <div style={{ padding: 16, border: "1px solid #e5e7eb", borderRadius: 12, background: "white" }}>
              <h2 style={{ margin: "0 0 8px" }}>Payment details</h2>
              <div style={{ display: "grid", gap: 8 }}>
                <input placeholder="Name on card" value={cardName} onChange={(e) => setCardName(e.target.value)} style={{ padding: 10, borderRadius: 8, border: "1px solid #d1d5db" }} />
                <input placeholder="Card number" value={cardNumber} onChange={(e) => setCardNumber(e.target.value.replace(/[^0-9 ]/g, ""))} style={{ padding: 10, borderRadius: 8, border: "1px solid #d1d5db" }} />
                <div style={{ display: "flex", gap: 8 }}>
                  <input placeholder="MM/YY" value={expiry} onChange={(e) => setExpiry(e.target.value)} style={{ padding: 10, borderRadius: 8, border: "1px solid #d1d5db", width: 120 }} />
                  <input placeholder="CVC" value={cvc} onChange={(e) => setCvc(e.target.value.replace(/[^0-9]/g, ""))} style={{ padding: 10, borderRadius: 8, border: "1px solid #d1d5db", width: 120 }} />
                </div>
              </div>
            </div>
          </div>

          <aside style={{ display: "grid", gap: 12, alignSelf: "start" }}>
            <div style={{ padding: 12, borderRadius: 12, border: "1px solid #e5e7eb", background: "white", color: "#6b7280" }}>
              <div style={{ fontSize: 14 }}>Need to edit your cart? <Link href="/cart">Go back to cart</Link></div>
            </div>

            <div style={{ padding: 16, border: "1px solid #e5e7eb", borderRadius: 12, background: "white" }}>
              <h3 style={{ margin: 0, marginBottom: 8 }}>Order summary</h3>
              <div style={{ display: "grid", gap: 8 }}>
                {cart.length === 0 ? <div style={{ color: "#6b7280" }}>Your cart is empty.</div> : cart.map((it) => (
                  <div key={it.id} style={{ display: "flex", justifyContent: "space-between", color: "#374151" }}>
                    <div>
                      <div style={{ fontWeight: 700 }}>{it.title}</div>
                      <div style={{ color: "#6b7280", fontSize: 13 }}>{it.qty} × {formatMoney(it.price)}</div>
                    </div>
                    <div style={{ fontWeight: 700 }}>{formatMoney(it.price * it.qty)}</div>
                  </div>
                ))}

                <div style={{ height: 1, background: "#e5e7eb", margin: "8px 0" }} />
                <div style={{ display: "flex", justifyContent: "space-between", color: "#4b5563" }}>
                  <span>Subtotal</span>
                  <span>{formatMoney(subtotal)}</span>
                </div>
                <div style={{ display: "flex", justifyContent: "space-between", color: "#4b5563" }}>
                  <span>Delivery fee</span>
                  <span>{formatMoney(delivery)}</span>
                </div>
                <div style={{ display: "flex", justifyContent: "space-between", fontWeight: 700, fontSize: 18 }}>
                  <span>Total</span>
                  <span>{formatMoney(total)}</span>
                </div>
              </div>
            </div>

            <button onClick={confirmOrder} style={{ width: "100%", background: "#0070f3", color: "white", padding: "14px 18px", borderRadius: 10, border: "none", fontSize: 16, fontWeight: 700 }}>
              Confirm order
            </button>

            <div style={{ display: "flex", alignItems: "center", gap: 8, color: "#6b7280", fontSize: 14 }}>
              <span style={{ fontSize: 16 }}>🔒</span>
              <span>Your payment information is secure and encrypted</span>
            </div>
          </aside>
        </section>
      </div>
    </main>
  );
}
