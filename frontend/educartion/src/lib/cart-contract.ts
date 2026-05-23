export type CartItem = {
  id: string | number;
  title: string;
  price: number;
  image?: string;
  description?: string;
  qty: number;
};

export type CartSummary = {
  items: CartItem[];
  subtotal: number;
  deliveryFee: number;
  total: number;
};

export type UpdateCartItemRequest = {
  id: string | number;
  qty: number;
};

export type RemoveCartItemRequest = {
  id: string | number;
};
