import type { CheckoutRequest } from "@/lib/payment-contract";
import type { CartItem } from "@/lib/cart-contract";
import { loadCartFromStorage, clearCart } from "./cartService";

export function loadCartForCheckout(): CartItem[] {
  return loadCartFromStorage();
}

export async function submitCheckout(
  request: CheckoutRequest
): Promise<{ ok: boolean; error?: string; orderId?: string }> {
  // Basic client-side validation
  const { deliveryAddress, paymentDetails } = request;

  if (
    !deliveryAddress.name ||
    !deliveryAddress.address ||
    !paymentDetails.cardName ||
    paymentDetails.cardNumber.length < 12
  ) {
    return { ok: false, error: "Please complete all required fields." };
  }

  try {
    // Simulate order submission
    // In a real app, this would call an API endpoint
    const orderId = `ORD-${Date.now()}`;

    // Clear the cart on successful submission
    clearCart();

    return { ok: true, orderId };
  } catch (err) {
    return { ok: false, error: "Failed to process order. Please try again." };
  }
}
