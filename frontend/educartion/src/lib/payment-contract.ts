import type { CartItem } from "./cart-contract";

export type DeliveryAddress = {
  name: string;
  address: string;
  city: string;
  postal: string;
};

export type PaymentDetails = {
  cardName: string;
  cardNumber: string;
  expiry: string;
  cvc: string;
};

export type CheckoutRequest = {
  deliveryAddress: DeliveryAddress;
  paymentDetails: PaymentDetails;
};

export type CheckoutResponse = {
  orderId: string;
  status: "success" | "error";
  message: string;
};

export type CheckoutPageData = {
  items: CartItem[];
  subtotal: number;
  deliveryFee: number;
  total: number;
  loaded: boolean;
};
