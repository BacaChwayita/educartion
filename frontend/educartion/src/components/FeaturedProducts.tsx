"use client";
import React, { useRef } from "react";

type Product = {
  id: string | number;
  title: string;
  price: number;
  image?: string;
  description?: string;
};

export default function FeaturedProducts({ products, onAddToCart }: { products: Product[]; onAddToCart?: (p: Product) => void }) {
  const scrollRef = useRef<HTMLDivElement>(null);

  const scroll = (direction: "left" | "right") => {
    if (scrollRef.current) {
      const scrollAmount = 400;
      scrollRef.current.scrollBy({
        left: direction === "left" ? -scrollAmount : scrollAmount,
        behavior: "smooth",
      });
    }
  };

  return (
    <div style={{ display: "flex", alignItems: "center", gap: 12 }}>
      <button
        onClick={() => scroll("left")}
        style={{
          flex: "0 0 40px",
          height: 40,
          background: "#f0f0f0",
          border: "1px solid #ddd",
          borderRadius: 6,
          cursor: "pointer",
          fontSize: 18,
          display: "flex",
          alignItems: "center",
          justifyContent: "center",
        }}
      >
        ‹
      </button>

      <div
        ref={scrollRef}
        style={{
          display: "flex",
          gap: 16,
          overflowX: "auto",
          scrollBehavior: "smooth",
          flex: 1,
          maxWidth: "calc(100% - 100px)",
          scrollbarWidth: "none",
          msOverflowStyle: "none",
        }}
      >
        <style>{`div::-webkit-scrollbar { display: none; }`}</style>
        {products.map((p) => (
          <div key={p.id} style={{ flex: "0 0 220px", minWidth: 220, border: "1px solid #eee", borderRadius: 8, padding: 12, background: "white", display: "flex", flexDirection: "column" }}>
            <div style={{ width: "100%", height: 160, marginBottom: 8 }}>
              <img src={p.image || "/placeholder.png"} alt={p.title} style={{ width: "100%", height: "100%", objectFit: "cover", borderRadius: 6 }} />
            </div>
            <div style={{ fontWeight: 700 }}>{p.title}</div>
            {p.description && <div style={{ color: "#666", fontSize: 13, marginTop: 6 }}>{p.description}</div>}
            <div style={{ marginTop: "auto", display: "flex", justifyContent: "space-between", alignItems: "center" }}>
              <div style={{ fontWeight: 700 }}>${p.price.toFixed(2)}</div>
              <button onClick={() => onAddToCart && onAddToCart(p)} style={{ background: "#0070f3", color: "white", border: "none", padding: "8px 12px", borderRadius: 6, cursor: "pointer" }}>
                Add to cart
              </button>
            </div>
          </div>
        ))}
      </div>

      <button
        onClick={() => scroll("right")}
        style={{
          flex: "0 0 40px",
          height: 40,
          background: "#f0f0f0",
          border: "1px solid #ddd",
          borderRadius: 6,
          cursor: "pointer",
          fontSize: 18,
          display: "flex",
          alignItems: "center",
          justifyContent: "center",
        }}
      >
        ›
      </button>
    </div>
  );
}
