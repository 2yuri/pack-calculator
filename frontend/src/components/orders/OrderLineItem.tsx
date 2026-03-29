import { Trash2 } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import type { Product, Pack } from "@/types";

function formatPackSizes(packs: Pack[] | undefined): string {
  if (!packs || packs.length === 0) return "No packs configured";
  const sorted = [...packs].sort((a, b) => a.size - b.size);
  return `Packs: ${sorted.map((p) => p.size.toLocaleString()).join(", ")}`;
}

interface OrderLineItemProps {
  products: Product[];
  packsMap: Map<number, Pack[]>;
  productId: string;
  quantity: string;
  onProductChange: (productId: string) => void;
  onQuantityChange: (quantity: string) => void;
  onRemove: () => void;
  canRemove: boolean;
  disabledProductIds: number[];
  productError?: string;
  quantityError?: string;
}

export function OrderLineItem({
  products,
  packsMap,
  productId,
  quantity,
  onProductChange,
  onQuantityChange,
  onRemove,
  canRemove,
  disabledProductIds,
  productError,
  quantityError,
}: OrderLineItemProps) {
  const selectedProduct = productId
    ? products.find((p) => p.id === Number(productId))
    : null;

  return (
    <div className="flex flex-col gap-2 rounded-lg border p-4">
      <div className="flex flex-col gap-1">
        <Select value={productId || undefined} onValueChange={(value) => { if (value) onProductChange(value); }}>
          <SelectTrigger>
            <SelectValue placeholder="Select product">
              {selectedProduct?.name}
            </SelectValue>
          </SelectTrigger>
          <SelectContent>
            {products.map((p) => (
              <SelectItem
                key={p.id}
                value={String(p.id)}
                disabled={disabledProductIds.includes(p.id)}
              >
                {p.name}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
        {productError && (
          <p className="text-sm text-destructive">{productError}</p>
        )}
      </div>
      {selectedProduct && (
        <p className="text-xs text-muted-foreground">
          {formatPackSizes(packsMap.get(selectedProduct.id))}
        </p>
      )}
      <div className="flex items-start gap-2">
        <div className="flex flex-col gap-1">
          <Input
            type="number"
            min="1"
            value={quantity}
            onChange={(e) => onQuantityChange(e.target.value)}
            placeholder="Qty"
            className="w-28"
          />
          {quantityError && (
            <p className="text-sm text-destructive">{quantityError}</p>
          )}
        </div>
        {canRemove && (
          <Button
            variant="ghost"
            size="icon"
            onClick={onRemove}
            className="shrink-0"
          >
            <Trash2 />
          </Button>
        )}
      </div>
    </div>
  );
}
