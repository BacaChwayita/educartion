import {
  normalizeCart,
  loadCartFromStorage,
  saveCartToStorage,
  clearCart,
  calculateCartSummary,
} from "@/services/cartService";

const STORAGE_KEY = "cart";

describe("cartService", () => {
  beforeEach(() => {
    localStorage.clear();
    jest.clearAllMocks();
  });

  describe("normalizeCart", () => {
    it("should merge items with same id and sum qty", () => {
      const input = [
        { id: "1", title: "Apple", price: 10, qty: 1 },
        { id: "1", title: "Apple", price: 10, qty: 2 },
      ];

      const result = normalizeCart(input);

      expect(result).toHaveLength(1);
      expect(result[0].qty).toBe(3);
    });

    it("should fallback defaults for missing fields", () => {
      const result = normalizeCart([{}]);

      expect(result[0].id).toBeDefined();
      expect(result[0].title).toBe("Untitled product");
      expect(result[0].price).toBe(0);
      expect(result[0].image).toBe("/images/product-placeholder.png");
      expect(result[0].qty).toBe(1);
    });

    it("should handle string qty and price conversion", () => {
      const result = normalizeCart([
        { id: "1", title: "Test", price: "5", qty: "2" },
      ]);

      expect(result[0].price).toBe(5);
      expect(result[0].qty).toBe(2);
    });
  });

  describe("loadCartFromStorage", () => {
    it("returns empty array on SSR (no window)", () => {
      const originalWindow = global.window;

      // @ts-ignore
      global.window = undefined;

      expect(loadCartFromStorage()).toEqual([]);

      global.window = originalWindow;
    });

    it("returns normalized cart from localStorage", () => {
      const data = [
        { id: "1", title: "A", price: 10, qty: 1 },
      ];

      localStorage.setItem(STORAGE_KEY, JSON.stringify(data));

      const result = loadCartFromStorage();
      expect(result).toHaveLength(1);
      expect(result[0].title).toBe("A");
    });

    it("returns empty array on invalid JSON", () => {
      localStorage.setItem(STORAGE_KEY, "invalid-json");

      expect(loadCartFromStorage()).toEqual([]);
    });
  });

  describe("saveCartToStorage", () => {
    it("saves items to localStorage", () => {
      const items = [
        { id: "1", title: "A", price: 10, qty: 1, image: "", description: "" },
      ];

      saveCartToStorage(items);

      const stored = localStorage.getItem(STORAGE_KEY);
      expect(JSON.parse(stored!)).toEqual(items);
    });
  });

  describe("clearCart", () => {
    it("removes cart from localStorage", () => {
      localStorage.setItem(STORAGE_KEY, "[]");

      clearCart();

      expect(localStorage.getItem(STORAGE_KEY)).toBeNull();
    });
  });

  describe("calculateCartSummary", () => {
    it("calculates subtotal, delivery fee, and total", () => {
      const items = [
        { id: "1", title: "A", price: 10, qty: 2, image: "", description: "" },
        { id: "2", title: "B", price: 5, qty: 1, image: "", description: "" },
      ];

      const summary = calculateCartSummary(items);

      expect(summary.subtotal).toBe(25); // 20 + 5
      expect(summary.deliveryFee).toBe(5);
      expect(summary.total).toBe(30);
    });

    it("sets delivery fee to 0 when cart is empty", () => {
      const summary = calculateCartSummary([]);

      expect(summary.subtotal).toBe(0);
      expect(summary.deliveryFee).toBe(0);
      expect(summary.total).toBe(0);
    });
  });
});
