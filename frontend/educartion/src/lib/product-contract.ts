export interface Product {
  id?: string | number;
  product_id?: string | number;
  supplier_id?: string | number;
  name?: string;
  description?: string;
  price?: number;
  discount_percent?: number;
  stock_quantity?: number;
  image?: string;
  image_url?: string;
  supplier_name?: string;
  is_active?: boolean;
  [key: string]: unknown;
}

export interface ProductsResponse {
  products?: Product[];
  [key: string]: unknown;
}

export interface BackendProductsResponse {
  error?: string;
  [key: string]: unknown;
}
