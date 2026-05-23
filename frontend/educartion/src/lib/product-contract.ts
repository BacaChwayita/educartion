export interface Product {
  id?: string | number;
  name?: string;
  description?: string;
  price?: number;
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
