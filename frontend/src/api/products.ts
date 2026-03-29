import type { Product } from "@/types";
import { apiRequest } from "./client";

export function getProducts(): Promise<Product[]> {
  return apiRequest<Product[]>("/products");
}

export function getProduct(id: number): Promise<Product> {
  return apiRequest<Product>(`/products/${id}`);
}

export function createProduct(name: string): Promise<Product> {
  return apiRequest<Product>("/products", {
    method: "POST",
    body: JSON.stringify({ name }),
  });
}

export function updateProduct(id: number, name: string): Promise<Product> {
  return apiRequest<Product>(`/products/${id}`, {
    method: "PUT",
    body: JSON.stringify({ name }),
  });
}

export function deleteProduct(id: number): Promise<void> {
  return apiRequest<void>(`/products/${id}`, {
    method: "DELETE",
  });
}
