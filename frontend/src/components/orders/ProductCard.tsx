import { Minus, Plus } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import type { Product, Pack } from "@/types";

function formatPackSizes(packs: Pack[]): string {
  if (packs.length === 0) return "No packs configured";
  const sorted = [...packs].sort((a, b) => a.size - b.size);
  return sorted.map((p) => p.size.toLocaleString()).join(", ");
}

interface ProductCardProps {
  product: Product;
  packs: Pack[];
  quantity: number;
  onQuantityChange: (quantity: number) => void;
}

export function ProductCard({
  product,
  packs,
  quantity,
  onQuantityChange,
}: ProductCardProps) {
  const hasPacks = packs.length > 0;

  return (
    <Card className={quantity > 0 ? "border-primary" : ""}>
      <CardContent className="flex flex-col gap-3 p-4">
        <div className="flex flex-col gap-0.5">
          <span className="text-xs text-muted-foreground">#{product.id}</span>
          <span className="font-semibold">{product.name}</span>
          <span className="text-xs text-muted-foreground">
            {hasPacks ? formatPackSizes(packs) : "No packs configured"}
          </span>
        </div>

        <div className="flex items-center gap-2">
          <Button
            variant="outline"
            size="icon"
            className="cursor-pointer transition-colors hover:text-primary"
            disabled={quantity <= 0}
            onClick={() => onQuantityChange(Math.max(0, quantity - 1))}
          >
            <Minus />
          </Button>
          <input
            type="number"
            min="0"
            value={quantity || ""}
            onChange={(e) => {
              const val = parseInt(e.target.value, 10);
              onQuantityChange(isNaN(val) ? 0 : Math.max(0, val));
            }}
            placeholder="0"
            className="h-8 w-16 rounded-lg border border-input bg-transparent text-center text-sm outline-none focus-visible:border-ring focus-visible:ring-3 focus-visible:ring-ring/50"
          />
          <Button
            variant="outline"
            size="icon"
            className="cursor-pointer transition-colors hover:text-primary"
            onClick={() => onQuantityChange(quantity + 1)}
          >
            <Plus />
          </Button>
        </div>
      </CardContent>
    </Card>
  );
}
