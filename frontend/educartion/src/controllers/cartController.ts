import type { CartItem, UpdateCartItemRequest, RemoveCartItemRequest } from "@/lib/cart-contract";
import * as cartService from "@/services/cartService";

export async function handleLoadCart() {
  return cartService.loadCartFromStorage();
}

export async function handleUpdateItemQuantity(
  items: CartItem[],
  request: UpdateCartItemRequest
): Promise<CartItem[]> {
  return items.map((item) =>
    String(item.id) === String(request.id)
      ? { ...item, qty: Math.max(1, request.qty) }
      : item
  );
}

export async function handleRemoveItem(
  items: CartItem[],
  request: RemoveCartItemRequest
): Promise<CartItem[]> {
  return items.filter((item) => String(item.id) !== String(request.id));
}

export async function handleSaveCart(items: CartItem[]): Promise<void> {
  cartService.saveCartToStorage(items);
}

export async function handleGetCartSummary(items: CartItem[]) {
  return cartService.calculateCartSummary(items);
}
