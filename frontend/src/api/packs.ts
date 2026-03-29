import type { Pack } from "@/types";
import { apiRequest } from "./client";

export function getPacks(productId: number): Promise<Pack[]> {
  return apiRequest<Pack[]>(`/products/${productId}/packs`);
}

export function createPack(productId: number, size: number): Promise<Pack> {
  return apiRequest<Pack>(`/products/${productId}/packs`, {
    method: "POST",
    body: JSON.stringify({ size }),
  });
}

export function createPacksBatch(productId: number, sizes: number[]): Promise<Pack[]> {
  return apiRequest<Pack[]>(`/products/${productId}/packs/batch`, {
    method: "POST",
    body: JSON.stringify({ sizes }),
  });
}

export function updatePack(packId: number, size: number): Promise<Pack> {
  return apiRequest<Pack>(`/packs/${packId}`, {
    method: "PUT",
    body: JSON.stringify({ size }),
  });
}

export function deletePack(packId: number): Promise<void> {
  return apiRequest<void>(`/packs/${packId}`, {
    method: "DELETE",
  });
}
