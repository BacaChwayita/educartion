import { submitCheckout, loadCartForCheckout } from "@/services/paymentService";
import * as cartService from "@/services/cartService";

jest.mock("@/services/cartService", () => ({
  loadCartFromStorage: jest.fn(),
  clearCart: jest.fn(),
}));

describe("paymentService", () => {
  afterEach(() => {
    jest.clearAllMocks();
  });

  describe("loadCartForCheckout", () => {
    it("should return cart items from storage", () => {
      const mockCart = [
        { id: 1, title: "Item 1", price: 100 },
        { id: 2, title: "Item 2", price: 200 },
      ];

      (cartService.loadCartFromStorage as jest.Mock).mockReturnValue(mockCart);

      const result = loadCartForCheckout();

      expect(result).toEqual(mockCart);
      expect(cartService.loadCartFromStorage).toHaveBeenCalled();
    });
  });

  describe("submitCheckout", () => {
    it("should fail when required delivery fields are missing", async () => {
      const result = await submitCheckout({
        deliveryAddress: {
          name: "",
          address: "",
        },
        paymentDetails: {
          cardName: "",
          cardNumber: "123",
        },
      } as any);

      expect(result.ok).toBe(false);
      expect(result.error).toBe("Please complete all required fields.");
    });

    it("should fail when payment details are invalid", async () => {
      const result = await submitCheckout({
        deliveryAddress: {
          name: "John",
          address: "Street 1",
        },
        paymentDetails: {
          cardName: "John",
          cardNumber: "123", // too short
        },
      } as any);

      expect(result.ok).toBe(false);
      expect(result.error).toBe("Please complete all required fields.");
    });

    it("should successfully submit checkout and clear cart", async () => {
      const result = await submitCheckout({
        deliveryAddress: {
          name: "John",
          address: "Street 1",
        },
        paymentDetails: {
          cardName: "John",
          cardNumber: "123456789012",
        },
      } as any);

      expect(result.ok).toBe(true);
      expect(result.orderId).toMatch(/^ORD-/);
      expect(cartService.clearCart).toHaveBeenCalled();
    });
  });
});
