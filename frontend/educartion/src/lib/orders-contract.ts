export interface OrderSummary {
  order_id: number;
  account_id: number;
  order_number: string;
  status: string;
  subtotal_amount: number;
  discount_amount: number;
  total_amount: number;
  placed_at: string;
}
