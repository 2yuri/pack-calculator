export interface Product {
  id: number;
  name: string;
  created_at: string;
  updated_at: string;
}

export interface Pack {
  id: number;
  product_id: number;
  size: number;
  created_at: string;
  updated_at: string;
}

export interface OrderItem {
  product_id: number;
  quantity: number;
}

export interface OrderRequest {
  items: OrderItem[];
}

export interface CalculationResultPack {
  pack: Pack;
  quantity: number;
}

export interface OrderResult {
  product_id: number;
  quantity: number;
  total_items: number;
  total_packs: number;
  packs: CalculationResultPack[];
}

export interface OrderResponse {
  results: OrderResult[];
}

export interface ApiError {
  error: string;
}
