import { useState } from "react";
import { useQueries } from "@tanstack/react-query";
import { toast } from "sonner";
import { ShoppingCart } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Skeleton } from "@/components/ui/skeleton";
import { ProductCard } from "@/components/orders/ProductCard";
import { OrderResults } from "@/components/orders/OrderResults";
import { useProducts } from "@/hooks/useProducts";
import { useCalculateOrder } from "@/hooks/useOrders";
import { getPacks } from "@/api/packs";
import { ApiRequestError } from "@/api/client";
import type { OrderResponse, Pack } from "@/types";

export function OrdersPage() {
  const { data: products, isLoading } = useProducts();
  const calculateOrder = useCalculateOrder();

  const packsQueries = useQueries({
    queries: (products ?? []).map((p) => ({
      queryKey: ["products", p.id, "packs"],
      queryFn: () => getPacks(p.id),
    })),
  });

  const packsMap = new Map<number, Pack[]>();
  (products ?? []).forEach((p, i) => {
    if (packsQueries[i]?.data) {
      packsMap.set(p.id, packsQueries[i].data);
    }
  });

  const [quantities, setQuantities] = useState<Map<number, number>>(new Map());
  const [results, setResults] = useState<OrderResponse | null>(null);

  function handleQuantityChange(productId: number, quantity: number) {
    setQuantities((prev) => {
      const next = new Map(prev);
      if (quantity <= 0) {
        next.delete(productId);
      } else {
        next.set(productId, quantity);
      }
      return next;
    });
  }

  function handleCalculate() {
    const items = Array.from(quantities.entries())
      .filter(([, qty]) => qty > 0)
      .map(([product_id, quantity]) => ({ product_id, quantity }));

    if (items.length === 0) {
      toast.error("Select at least one product");
      return;
    }

    calculateOrder.mutate(
      { items },
      {
        onSuccess: (data) => setResults(data),
        onError: (err) => {
          toast.error(
            err instanceof ApiRequestError
              ? err.message
              : "Failed to calculate order",
          );
        },
      },
    );
  }

  const totalItems = Array.from(quantities.values()).reduce((a, b) => a + b, 0);

  if (isLoading) {
    return (
      <div className="flex flex-col gap-4">
        <Skeleton className="h-8 w-48" />
        <div className="grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
          {Array.from({ length: 6 }).map((_, i) => (
            <Skeleton key={i} className="h-32 w-full" />
          ))}
        </div>
      </div>
    );
  }

  return (
    <div className="flex flex-col gap-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold">Order Calculator</h1>
        <Button
          onClick={handleCalculate}
          disabled={totalItems === 0 || calculateOrder.isPending}
        >
          <ShoppingCart className="mr-2" />
          {calculateOrder.isPending
            ? "Calculating..."
            : `Calculate (${totalItems})`}
        </Button>
      </div>

      {(!products || products.length === 0) ? (
        <p className="text-sm text-muted-foreground">
          No products yet. Create products first.
        </p>
      ) : (
        <div className="grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
          {products.map((product) => (
            <ProductCard
              key={product.id}
              product={product}
              packs={packsMap.get(product.id) ?? []}
              quantity={quantities.get(product.id) ?? 0}
              onQuantityChange={(qty) => handleQuantityChange(product.id, qty)}
            />
          ))}
        </div>
      )}

      {results && (
        <OrderResults results={results.results} products={products ?? []} />
      )}
    </div>
  );
}
