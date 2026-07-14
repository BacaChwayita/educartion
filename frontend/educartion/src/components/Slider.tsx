"use client";
import React, { useEffect, useState } from "react";

type Slide = {
  id: string | number;
  image: string;
  title?: string;
  subtitle?: string;
};

export default function Slider({ slides, interval = 4000 }: { slides: Slide[]; interval?: number }) {
  const [index, setIndex] = useState(0);

  useEffect(() => {
    if (!slides || slides.length === 0) return;
    const t = setInterval(() => setIndex((i) => (i + 1) % slides.length), interval);
    return () => clearInterval(t);
  }, [slides, interval]);

  if (!slides || slides.length === 0) return null;

  const prev = () => setIndex((i) => (i - 1 + slides.length) % slides.length);
  const next = () => setIndex((i) => (i + 1) % slides.length);

  return (
    <div style={{ position: "relative", width: "100%", overflow: "hidden", borderRadius: 8 }}>
      {slides.map((s, i) => (
        <div
          key={s.id}
          style={{
            display: i === index ? "block" : "none",
            width: "100%",
            transition: "opacity .5s ease-in-out",
          }}
        >
          <img src={s.image} alt={s.title || `slide-${i}`} style={{ width: "100%", height: 320, objectFit: "cover" }} />
          {(s.title || s.subtitle) && (
            <div style={{ position: "absolute", left: 20, bottom: 20, color: "white", textShadow: "0 1px 4px rgba(0,0,0,.6)" }}>
              {s.title && <div style={{ fontSize: 22, fontWeight: 700 }}>{s.title}</div>}
              {s.subtitle && <div style={{ fontSize: 14 }}>{s.subtitle}</div>}
            </div>
          )}
        </div>
      ))}

      <button onClick={prev} aria-label="Previous" style={{ position: "absolute", left: 10, top: "50%", transform: "translateY(-50%)", background: "rgba(0,0,0,.4)", color: "white", border: "none", padding: 8, borderRadius: 4 }}>
        ‹
      </button>
      <button onClick={next} aria-label="Next" style={{ position: "absolute", right: 10, top: "50%", transform: "translateY(-50%)", background: "rgba(0,0,0,.4)", color: "white", border: "none", padding: 8, borderRadius: 4 }}>
        ›
      </button>
    </div>
  );
}
