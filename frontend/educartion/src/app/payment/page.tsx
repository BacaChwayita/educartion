"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import React, { useEffect, useState } from "react";

type CartItem = { id: string | number; title: string; price: number; image?: string; description?: string; qty: number };

const DELIVERY_FEE = 5.0;
const formatMoney = (v: number) => `$${v.toFixed(2)}`;

export default function PaymentPage() {
  const router = useRouter();
  const [cart, setCart] = useState<CartItem[]>([]);
  const [loaded, setLoaded] = useState(false);

  const [name, setName] = useState("");
  const [address, setAddress] = useState("");
  const [city, setCity] = useState("");
  const [postal, setPostal] = useState("");

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
    if (!name || !address || !cardName || cardNumber.length < 12) {
      alert("Please complete your delivery and payment details.");
      return;
    }

    window.localStorage.removeItem("cart");
    setCart([]);
    alert("Order confirmed — thank you!");
    router.push("/orders");
  };

  return (
    <section className="relative min-h-screen overflow-hidden bg-[radial-gradient(circle_at_top_left,rgba(251,191,36,0.2),transparent_30%),radial-gradient(circle_at_bottom_right,rgba(14,165,233,0.18),transparent_28%),linear-gradient(180deg,#060816_0%,#0b1020_100%)] px-6 py-10 text-slate-100 sm:px-8 lg:px-10">
      <div className="absolute inset-0 bg-[linear-gradient(rgba(255,255,255,0.03)_1px,transparent_1px),linear-gradient(90deg,rgba(255,255,255,0.03)_1px,transparent_1px)] bg-size-[28px_28px] opacity-20" />
      <div className="relative mx-auto w-full max-w-4xl">
        <div className="w-full rounded-4xl border border-white/10 bg-white/5 p-6 shadow-2xl shadow-black/30 backdrop-blur-xl lg:p-8">
          <div className="rounded-3xl border border-white/10 bg-slate-950/80 p-6 sm:p-8">
            <div className="mb-8 space-y-2">
              <p className="text-sm font-medium uppercase tracking-[0.3em] text-amber-200">
                Payment
              </p>
              <h2 className="text-3xl font-semibold tracking-tight text-white">Checkout</h2>
              <p className="text-sm leading-6 text-slate-300">Complete your order with delivery and payment details.</p>
            </div>

            <div className="grid gap-6 lg:grid-cols-3">
              <div className="lg:col-span-2 space-y-6">
                <div className="rounded-2xl border border-white/10 bg-white/5 p-6">
                  <h3 className="text-sm font-semibold text-white mb-4 uppercase tracking-[0.2em] text-amber-200">Delivery Address</h3>
                  <div className="space-y-4">
                    <input
                      placeholder="Full name"
                      value={name}
                      onChange={(e) => setName(e.target.value)}
                      className="w-full rounded-2xl border border-white/10 bg-white/5 px-4 py-3 text-slate-100 outline-none transition placeholder:text-slate-500 focus:border-amber-300/60 focus:bg-white/8 focus:ring-2 focus:ring-amber-300/20"
                    />
                    <input
                      placeholder="Street address"
                      value={address}
                      onChange={(e) => setAddress(e.target.value)}
                      className="w-full rounded-2xl border border-white/10 bg-white/5 px-4 py-3 text-slate-100 outline-none transition placeholder:text-slate-500 focus:border-amber-300/60 focus:bg-white/8 focus:ring-2 focus:ring-amber-300/20"
                    />
                    <div className="grid grid-cols-3 gap-3">
                      <input
                        placeholder="City"
                        value={city}
                        onChange={(e) => setCity(e.target.value)}
                        className="col-span-2 rounded-2xl border border-white/10 bg-white/5 px-4 py-3 text-slate-100 outline-none transition placeholder:text-slate-500 focus:border-amber-300/60 focus:bg-white/8 focus:ring-2 focus:ring-amber-300/20"
                      />
                      <input
                        placeholder="Postal code"
                        value={postal}
                        onChange={(e) => setPostal(e.target.value)}
                        className="rounded-2xl border border-white/10 bg-white/5 px-4 py-3 text-slate-100 outline-none transition placeholder:text-slate-500 focus:border-amber-300/60 focus:bg-white/8 focus:ring-2 focus:ring-amber-300/20"
                      />
                    </div>
                  </div>
                </div>

                <div className="rounded-2xl border border-white/10 bg-white/5 p-6">
                  <h3 className="text-sm font-semibold text-white mb-4 uppercase tracking-[0.2em] text-amber-200">Payment Details</h3>
                  <div className="space-y-4">
                    <input
                      placeholder="Name on card"
                      value={cardName}
                      onChange={(e) => setCardName(e.target.value)}
                      className="w-full rounded-2xl border border-white/10 bg-white/5 px-4 py-3 text-slate-100 outline-none transition placeholder:text-slate-500 focus:border-amber-300/60 focus:bg-white/8 focus:ring-2 focus:ring-amber-300/20"
                    />
                    <input
                      placeholder="Card number"
                      value={cardNumber}
                      onChange={(e) => setCardNumber(e.target.value.replace(/[^0-9 ]/g, ""))}
                      className="w-full rounded-2xl border border-white/10 bg-white/5 px-4 py-3 text-slate-100 outline-none transition placeholder:text-slate-500 focus:border-amber-300/60 focus:bg-white/8 focus:ring-2 focus:ring-amber-300/20"
                    />
                    <div className="grid grid-cols-2 gap-3">
                      <input
                        placeholder="MM/YY"
                        value={expiry}
                        onChange={(e) => setExpiry(e.target.value)}
                        className="rounded-2xl border border-white/10 bg-white/5 px-4 py-3 text-slate-100 outline-none transition placeholder:text-slate-500 focus:border-amber-300/60 focus:bg-white/8 focus:ring-2 focus:ring-amber-300/20"
                      />
                      <input
                        placeholder="CVC"
                        value={cvc}
                        onChange={(e) => setCvc(e.target.value.replace(/[^0-9]/g, ""))}
                        className="rounded-2xl border border-white/10 bg-white/5 px-4 py-3 text-slate-100 outline-none transition placeholder:text-slate-500 focus:border-amber-300/60 focus:bg-white/8 focus:ring-2 focus:ring-amber-300/20"
                      />
                    </div>
                  </div>
                </div>
              </div>

              <div className="space-y-4">
                <div className="rounded-2xl border border-white/10 bg-white/5 px-4 py-4 text-sm text-slate-300">
                  <div>Need to edit your cart?{" "}
                    <Link href="/cart" className="text-amber-300 hover:text-amber-200 transition">
                      Go back to cart
                    </Link>
                  </div>
                </div>

                <div className="rounded-2xl border border-white/10 bg-white/5 p-6">
                  <h3 className="text-sm font-semibold text-white mb-4">Order Summary</h3>
                  <div className="space-y-3 text-sm">
                    {cart.length === 0 ? (
                      <div className="text-slate-400">Your cart is empty.</div>
                    ) : (
                      cart.map((it) => (
                        <div key={it.id} className="flex justify-between text-slate-300">
                          <div>
                            <div className="font-semibold">{it.title}</div>
                            <div className="text-xs text-slate-400">{it.qty} × {formatMoney(it.price)}</div>
                          </div>
                          <div className="font-semibold text-amber-400">{formatMoney(it.price * it.qty)}</div>
                        </div>
                      ))
                    )}

                    <div className="border-t border-white/10 pt-3 space-y-2 text-slate-300">
                      <div className="flex justify-between">
                        <span>Subtotal</span>
                        <span>{formatMoney(subtotal)}</span>
                      </div>
                      <div className="flex justify-between">
                        <span>Delivery fee</span>
                        <span>{formatMoney(delivery)}</span>
                      </div>
                      <div className="border-t border-white/10 pt-3 flex justify-between font-semibold text-amber-400 text-base">
                        <span>Total</span>
                        <span>{formatMoney(total)}</span>
                      </div>
                    </div>
                  </div>
                </div>

                <button
                  onClick={confirmOrder}
                  className="w-full flex items-center justify-center rounded-2xl bg-amber-400 px-5 py-3.5 text-sm font-semibold text-slate-950 transition hover:bg-amber-300"
                >
                  Confirm order
                </button>

                <div className="flex items-center gap-2 text-xs text-slate-400">
                  <span className="text-sm">🔒</span>
                  <span>Your payment information is secure and encrypted</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}
