import type { CheckoutRequest } from "@/lib/payment-contract";
import * as paymentService from "@/services/paymentService";

export async function handleLoadCheckoutCart() {
  return paymentService.loadCartForCheckout();
}

export async function handleSubmitCheckout(request: CheckoutRequest) {
  return await paymentService.submitCheckout(request);
}
