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

export interface OrderItemProductImage {
  product_image_id?: number;
  product_id?: number;
  image_url?: string;
  alt_text?: string;
  is_primary?: boolean;
  sort_order?: number;
  created_at?: string;
  [key: string]: unknown;
}

export interface OrderItemSupplier {
  supplier_id: number;
  name: string;
}

export interface OrderItemProductDetails {
  product_id: number;
  name: string;
  description: string;
  price: number;
  discount_percent: number;
  stock_quantity: number;
  is_active: boolean;
  supplier_details?: OrderItemSupplier;
  product_image?: OrderItemProductImage[];
  product_images?: OrderItemProductImage[];
  [key: string]: unknown;
}

export interface OrderItemDetail {
  quantity: number;
  unit_price: number;
  discount_amount: number;
  product_details: OrderItemProductDetails;
}

export interface OrderDetail {
  order_id: number;
  order_number: string;
  status: string;
  subtotal_amount: number;
  discount_amount: number;
  total_amount: number;
  placed_at: string;
  order_item: OrderItemDetail[];
}
