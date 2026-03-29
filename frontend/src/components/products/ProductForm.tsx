import { useState, useEffect } from "react";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import {
  Drawer,
  DrawerContent,
  DrawerHeader,
  DrawerTitle,
} from "@/components/ui/drawer";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { useMediaQuery } from "@/hooks/useMediaQuery";
import type { Product } from "@/types";

interface ProductFormProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  product?: Product | null;
  onSubmit: (name: string) => void;
  isPending: boolean;
}

function ProductFormContent({
  product,
  onSubmit,
  isPending,
}: Omit<ProductFormProps, "open" | "onOpenChange">) {
  const [name, setName] = useState(product?.name ?? "");
  const [error, setError] = useState("");

  useEffect(() => {
    setName(product?.name ?? "");
    setError("");
  }, [product]);

  function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    const trimmed = name.trim();
    if (!trimmed) {
      setError("Product name is required");
      return;
    }
    setError("");
    onSubmit(trimmed);
  }

  return (
    <form onSubmit={handleSubmit} className="flex flex-col gap-4 p-4 sm:p-0">
      <div className="flex flex-col gap-2">
        <Label htmlFor="product-name">Name</Label>
        <Input
          id="product-name"
          value={name}
          onChange={(e) => setName(e.target.value)}
          placeholder="Enter product name"
          autoFocus
        />
        {error && <p className="text-sm text-destructive">{error}</p>}
      </div>
      <Button type="submit" disabled={isPending} className="w-full">
        {isPending ? "Saving..." : product ? "Update" : "Create"}
      </Button>
    </form>
  );
}

export function ProductForm({
  open,
  onOpenChange,
  product,
  onSubmit,
  isPending,
}: ProductFormProps) {
  const isDesktop = useMediaQuery("(min-width: 640px)");
  const title = product ? "Edit Product" : "New Product";

  if (isDesktop) {
    return (
      <Dialog open={open} onOpenChange={onOpenChange}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>{title}</DialogTitle>
          </DialogHeader>
          <ProductFormContent
            product={product}
            onSubmit={onSubmit}
            isPending={isPending}
          />
        </DialogContent>
      </Dialog>
    );
  }

  return (
    <Drawer open={open} onOpenChange={onOpenChange}>
      <DrawerContent>
        <DrawerHeader>
          <DrawerTitle>{title}</DrawerTitle>
        </DrawerHeader>
        <ProductFormContent
          product={product}
          onSubmit={onSubmit}
          isPending={isPending}
        />
      </DrawerContent>
    </Drawer>
  );
}
