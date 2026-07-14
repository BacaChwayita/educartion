"use client";
import React, { useRef } from "react";

type Category = {
  id: string | number;
  name: string;
  image?: string;
};

export default function Categories({ categories, onSelect }: { categories: Category[]; onSelect?: (id: string | number) => void }) {
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
          gap: 12,
          overflowX: "auto",
          scrollBehavior: "smooth",
          flex: 1,
          maxWidth: "calc(100% - 100px)",
          scrollbarWidth: "none",
          msOverflowStyle: "none",
        }}
      >
        <style>{`div::-webkit-scrollbar { display: none; }`}</style>
        {categories.map((c) => (
          <button
            key={c.id}
            onClick={() => onSelect && onSelect(c.id)}
            style={{
              flex: "0 0 160px",
              minWidth: 160,
              display: "flex",
              flexDirection: "column",
              alignItems: "center",
              padding: 12,
              borderRadius: 8,
              border: "1px solid #eee",
              background: "white",
              cursor: "pointer",
            }}
          >
            {c.image ? <img src={c.image} alt={c.name} style={{ width: 80, height: 80, objectFit: "cover", borderRadius: 8 }} /> : <div style={{ width: 80, height: 80, background: "#f3f3f3", borderRadius: 8 }} />}
            <div style={{ marginTop: 8, fontWeight: 600, fontSize: 14 }}>{c.name}</div>
          </button>
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
