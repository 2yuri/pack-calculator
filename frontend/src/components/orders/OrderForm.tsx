import { Plus } from "lucide-react";
import { Button } from "@/components/ui/button";
import { OrderLineItem } from "./OrderLineItem";
import type { Product, Pack } from "@/types";

export interface LineItem {
  productId: string;
  quantity: string;
}

export interface LineItemErrors {
  productError?: string;
  quantityError?: string;
}

interface OrderFormProps {
  products: Product[];
  packsMap: Map<number, Pack[]>;
  items: LineItem[];
  errors: LineItemErrors[];
  onItemChange: (index: number, field: keyof LineItem, value: string) => void;
  onAddItem: () => void;
  onRemoveItem: (index: number) => void;
  onSubmit: () => void;
  isPending: boolean;
}

export function OrderForm({
  products,
  packsMap,
  items,
  errors,
  onItemChange,
  onAddItem,
  onRemoveItem,
  onSubmit,
  isPending,
}: OrderFormProps) {
  const selectedProductIds = items
    .map((item) => parseInt(item.productId, 10))
    .filter((id) => !isNaN(id));

  return (
    <div className="flex flex-col gap-4">
      {items.map((item, index) => (
        <OrderLineItem
          key={index}
          products={products}
          packsMap={packsMap}
          productId={item.productId}
          quantity={item.quantity}
          onProductChange={(v) => onItemChange(index, "productId", v)}
          onQuantityChange={(v) => onItemChange(index, "quantity", v)}
          onRemove={() => onRemoveItem(index)}
          canRemove={items.length > 1}
          disabledProductIds={selectedProductIds.filter(
            (id) => id !== parseInt(item.productId, 10),
          )}
          productError={errors[index]?.productError}
          quantityError={errors[index]?.quantityError}
        />
      ))}

      <div className="flex gap-2">
        <Button variant="outline" onClick={onAddItem}>
          <Plus className="mr-2" />
          Add Item
        </Button>
        <Button onClick={onSubmit} disabled={isPending}>
          {isPending ? "Calculating..." : "Calculate"}
        </Button>
      </div>
    </div>
  );
}
