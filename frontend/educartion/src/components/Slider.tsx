"use client";

import type { StaticImageData } from "next/image";
import { useEffect, useState } from "react";

type Slide = {
  id: string | number;
  image: string | StaticImageData;
  title?: string;
  subtitle?: string;
};

export default function Slider({ slides, interval = 4500 }: { slides: Slide[]; interval?: number }) {
  const [index, setIndex] = useState(0);

  useEffect(() => {
    if (!slides || slides.length === 0) return;
    const timer = setInterval(() => setIndex((current) => (current + 1) % slides.length), interval);
    return () => clearInterval(timer);
  }, [slides, interval]);

  if (!slides || slides.length === 0) return null;

  const prev = () => setIndex((current) => (current - 1 + slides.length) % slides.length);
  const next = () => setIndex((current) => (current + 1) % slides.length);

  return (
    <div className="relative overflow-hidden rounded-[2rem] border border-slate-200/80 bg-slate-950/90 shadow-sm dark:border-white/10 dark:bg-slate-900/80">
      {slides.map((slide, idx) => {
        const imageSrc = typeof slide.image === "string" ? slide.image : slide.image.src;

        return (
          <div
            key={slide.id}
            className={`relative h-[360px] w-full transition duration-700 ${idx === index ? "opacity-100" : "pointer-events-none absolute inset-0 opacity-0"}`}
          >
            <img
              src={imageSrc}
              alt={slide.title ?? `slide-${idx}`}
              className="h-full w-full object-cover"
            />
            <div className="absolute inset-0 bg-gradient-to-t from-slate-950/80 via-slate-950/10 to-transparent" />
            <div className="absolute left-6 bottom-6 right-6 text-white sm:left-10 sm:bottom-10">
              {slide.title && <h2 className="text-3xl font-semibold tracking-tight sm:text-4xl">{slide.title}</h2>}
              {slide.subtitle && <p className="mt-3 max-w-2xl text-sm sm:text-base text-slate-100/90">{slide.subtitle}</p>}
            </div>
          </div>
        );
      })}

      <button
        type="button"
        onClick={prev}
        aria-label="Previous slide"
        className="absolute left-4 top-1/2 z-10 flex h-11 w-11 -translate-y-1/2 items-center justify-center rounded-full border border-white/20 bg-slate-950/70 text-white shadow-lg transition hover:bg-slate-950/90 sm:left-6"
      >
        ‹
      </button>
      <button
        type="button"
        onClick={next}
        aria-label="Next slide"
        className="absolute right-4 top-1/2 z-10 flex h-11 w-11 -translate-y-1/2 items-center justify-center rounded-full border border-white/20 bg-slate-950/70 text-white shadow-lg transition hover:bg-slate-950/90 sm:right-6"
      >
        ›
      </button>

      <div className="absolute bottom-4 left-1/2 flex -translate-x-1/2 gap-2">
        {slides.map((_, dotIndex) => (
          <button
            key={dotIndex}
            type="button"
            onClick={() => setIndex(dotIndex)}
            className={`h-2.5 w-2.5 rounded-full transition ${dotIndex === index ? "bg-amber-400" : "bg-white/40 hover:bg-white/70"}`}
            aria-label={`Go to slide ${dotIndex + 1}`}
          />
        ))}
      </div>
    </div>
  );
}
