"use client";

import { useRef } from "react";

type Category = {
  id: string | number;
  name: string;
  image?: string;
};

export default function Categories({ categories, onSelect }: { categories: Category[]; onSelect?: (id: string | number) => void }) {
  const scrollRef = useRef<HTMLDivElement>(null);

  const scroll = (direction: "left" | "right") => {
    if (scrollRef.current) {
      const scrollAmount = 340;
      scrollRef.current.scrollBy({
        left: direction === "left" ? -scrollAmount : scrollAmount,
        behavior: "smooth",
      });
    }
  };

  return (
    <div className="flex items-center gap-3 overflow-hidden">
      <button
        type="button"
        onClick={() => scroll("left")}
        className="flex h-11 w-11 items-center justify-center rounded-full border border-slate-200/80 bg-white/90 text-lg font-semibold text-slate-900 shadow-sm transition hover:bg-slate-100 dark:border-white/10 dark:bg-white/5 dark:text-slate-200 dark:hover:bg-white/10"
        aria-label="Scroll categories left"
      >
        ‹
      </button>

      <div
        ref={scrollRef}
        className="flex w-full gap-4 overflow-x-auto pb-1 scroll-smooth scrollbar-none"
      >
        {categories.map((category) => (
          <button
            key={category.id}
            type="button"
            onClick={() => onSelect?.(category.id)}
            className="group flex min-w-[180px] flex-col overflow-hidden rounded-[1.5rem] border border-slate-200/80 bg-white/90 shadow-sm transition hover:-translate-y-0.5 hover:shadow-md dark:border-white/10 dark:bg-white/5 dark:shadow-black/10"
          >
            <div className="relative h-32 overflow-hidden bg-slate-100 dark:bg-slate-900">
              {category.image ? (
                <img src={category.image} alt={category.name} className="h-full w-full object-cover transition-transform duration-500 group-hover:scale-105" />
              ) : (
                <div className="h-full w-full bg-slate-200 dark:bg-slate-800" />
              )}
            </div>
            <div className="px-4 py-3 text-left">
              <p className="text-sm font-semibold text-slate-900 transition group-hover:text-amber-700 dark:text-slate-100 dark:group-hover:text-amber-300">
                {category.name}
              </p>
            </div>
          </button>
        ))}
      </div>

      <button
        type="button"
        onClick={() => scroll("right")}
        className="flex h-11 w-11 items-center justify-center rounded-full border border-slate-200/80 bg-white/90 text-lg font-semibold text-slate-900 shadow-sm transition hover:bg-slate-100 dark:border-white/10 dark:bg-white/5 dark:text-slate-200 dark:hover:bg-white/10"
        aria-label="Scroll categories right"
      >
        ›
      </button>
    </div>
  );
}
