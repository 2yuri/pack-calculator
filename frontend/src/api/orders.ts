import type { OrderRequest, OrderResponse } from "@/types";
import { apiRequest } from "./client";

export function calculateOrder(request: OrderRequest): Promise<OrderResponse> {
  return apiRequest<OrderResponse>("/orders", {
    method: "POST",
    body: JSON.stringify(request),
  });
}
