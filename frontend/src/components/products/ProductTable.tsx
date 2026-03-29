import { useNavigate } from "@tanstack/react-router";
import { Package, Pencil, Trash2 } from "lucide-react";
import { Button } from "@/components/ui/button";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Card, CardContent } from "@/components/ui/card";
import type { Product } from "@/types";

interface ProductTableProps {
  products: Product[];
  onEdit: (product: Product) => void;
  onDelete: (product: Product) => void;
}

function formatDate(dateStr: string) {
  return new Date(dateStr).toLocaleDateString();
}

const iconBtn = "cursor-pointer transition-colors hover:text-primary";
const iconBtnDanger = "cursor-pointer transition-colors hover:text-destructive";

function ActionButtons({
  product,
  onPacks,
  onEdit,
  onDelete,
}: {
  product: Product;
  onPacks: (product: Product) => void;
  onEdit: (product: Product) => void;
  onDelete: (product: Product) => void;
}) {
  return (
    <div className="flex gap-1">
      <Button
        variant="ghost"
        size="icon"
        className={iconBtn}
        onClick={(e) => {
          e.stopPropagation();
          onPacks(product);
        }}
      >
        <Package />
      </Button>
      <Button
        variant="ghost"
        size="icon"
        className={iconBtn}
        onClick={(e) => {
          e.stopPropagation();
          onEdit(product);
        }}
      >
        <Pencil />
      </Button>
      <Button
        variant="ghost"
        size="icon"
        className={iconBtnDanger}
        onClick={(e) => {
          e.stopPropagation();
          onDelete(product);
        }}
      >
        <Trash2 />
      </Button>
    </div>
  );
}

export function ProductTable({
  products,
  onEdit,
  onDelete,
}: ProductTableProps) {
  const navigate = useNavigate();

  function handleRowClick(product: Product) {
    navigate({ to: "/products/$productId/packs", params: { productId: String(product.id) } });
  }

  return (
    <>
      {/* Desktop table */}
      <div className="hidden sm:block">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Name</TableHead>
              <TableHead>Created</TableHead>
              <TableHead className="w-32">Actions</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {products.map((product) => (
              <TableRow
                key={product.id}
                className="cursor-pointer"
                onClick={() => handleRowClick(product)}
              >
                <TableCell className="font-medium">{product.name}</TableCell>
                <TableCell>{formatDate(product.created_at)}</TableCell>
                <TableCell>
                  <ActionButtons
                    product={product}
                    onPacks={handleRowClick}
                    onEdit={onEdit}
                    onDelete={onDelete}
                  />
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </div>

      {/* Mobile cards */}
      <div className="flex flex-col gap-3 sm:hidden">
        {products.map((product) => (
          <Card
            key={product.id}
            className="cursor-pointer"
            onClick={() => handleRowClick(product)}
          >
            <CardContent className="flex items-center justify-between p-4">
              <div>
                <p className="font-medium">{product.name}</p>
                <p className="text-sm text-muted-foreground">
                  {formatDate(product.created_at)}
                </p>
              </div>
              <ActionButtons
                product={product}
                onPacks={handleRowClick}
                onEdit={onEdit}
                onDelete={onDelete}
              />
            </CardContent>
          </Card>
        ))}
      </div>
    </>
  );
}
