import { useMutation } from "@tanstack/react-query";
import { calculateOrder } from "@/api/orders";
import type { OrderRequest } from "@/types";

export function useCalculateOrder() {
  return useMutation({
    mutationFn: (request: OrderRequest) => calculateOrder(request),
  });
}
