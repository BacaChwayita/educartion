import type { CartItem, CartSummary } from "@/lib/cart-contract";

const DELIVERY_FEE = 5.0;
const STORAGE_KEY = "cart";

export function normalizeCart(cart: any[]): CartItem[] {
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
}

export function loadCartFromStorage(): CartItem[] {
  if (typeof window === "undefined") return [];
  
  const stored = window.localStorage.getItem(STORAGE_KEY) || "[]";
  try {
    const parsed = JSON.parse(stored);
    return normalizeCart(parsed);
  } catch {
    return [];
  }
}

export function saveCartToStorage(items: CartItem[]): void {
  if (typeof window === "undefined") return;
  window.localStorage.setItem(STORAGE_KEY, JSON.stringify(items));
}

export function clearCart(): void {
  if (typeof window === "undefined") return;
  window.localStorage.removeItem(STORAGE_KEY);
}

export function calculateCartSummary(items: CartItem[]): CartSummary {
  const subtotal = items.reduce((sum, item) => sum + item.price * item.qty, 0);
  const deliveryFee = items.length > 0 ? DELIVERY_FEE : 0;
  const total = subtotal + deliveryFee;

  return { items, subtotal, deliveryFee, total };
}
