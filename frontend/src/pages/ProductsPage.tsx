import { useState } from "react";
import { Plus, Package } from "lucide-react";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "@/components/ui/alert-dialog";
import { Skeleton } from "@/components/ui/skeleton";
import { ProductTable } from "@/components/products/ProductTable";
import { ProductForm } from "@/components/products/ProductForm";
import {
  useProducts,
  useCreateProduct,
  useUpdateProduct,
  useDeleteProduct,
} from "@/hooks/useProducts";
import { ApiRequestError } from "@/api/client";
import type { Product } from "@/types";

export function ProductsPage() {
  const { data: products, isLoading } = useProducts();
  const createProduct = useCreateProduct();
  const updateProduct = useUpdateProduct();
  const deleteProduct = useDeleteProduct();

  const [formOpen, setFormOpen] = useState(false);
  const [editingProduct, setEditingProduct] = useState<Product | null>(null);
  const [deletingProduct, setDeletingProduct] = useState<Product | null>(null);

  function handleCreate(name: string) {
    createProduct.mutate(name, {
      onSuccess: () => {
        setFormOpen(false);
        toast.success("Product created");
      },
      onError: (err) => {
        toast.error(err instanceof ApiRequestError ? err.message : "Failed to create product");
      },
    });
  }

  function handleUpdate(name: string) {
    if (!editingProduct) return;
    updateProduct.mutate(
      { id: editingProduct.id, name },
      {
        onSuccess: () => {
          setEditingProduct(null);
          toast.success("Product updated");
        },
        onError: (err) => {
          toast.error(err instanceof ApiRequestError ? err.message : "Failed to update product");
        },
      },
    );
  }

  function handleDelete() {
    if (!deletingProduct) return;
    deleteProduct.mutate(deletingProduct.id, {
      onSuccess: () => {
        setDeletingProduct(null);
        toast.success("Product deleted");
      },
      onError: (err) => {
        toast.error(err instanceof ApiRequestError ? err.message : "Failed to delete product");
      },
    });
  }

  if (isLoading) {
    return (
      <div className="flex flex-col gap-4">
        <div className="flex items-center justify-between">
          <Skeleton className="h-8 w-32" />
          <Skeleton className="h-9 w-28" />
        </div>
        <div className="flex flex-col gap-2">
          {Array.from({ length: 3 }).map((_, i) => (
            <Skeleton key={i} className="h-12 w-full" />
          ))}
        </div>
      </div>
    );
  }

  const isEmpty = !products || products.length === 0;

  return (
    <div className="flex flex-col gap-4">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold">Products</h1>
        <Button onClick={() => setFormOpen(true)}>
          <Plus className="mr-2" />
          New Product
        </Button>
      </div>

      {isEmpty ? (
        <div className="flex flex-col items-center justify-center gap-4 rounded-lg border border-dashed py-12">
          <Package className="size-12 text-muted-foreground" />
          <div className="text-center">
            <p className="text-lg font-medium">No products yet</p>
            <p className="text-sm text-muted-foreground">
              Create your first product to get started.
            </p>
          </div>
          <Button onClick={() => setFormOpen(true)}>
            <Plus className="mr-2" />
            Create Product
          </Button>
        </div>
      ) : (
        <ProductTable
          products={products}
          onEdit={(p) => setEditingProduct(p)}
          onDelete={(p) => setDeletingProduct(p)}
        />
      )}

      {/* Create dialog */}
      <ProductForm
        open={formOpen}
        onOpenChange={setFormOpen}
        onSubmit={handleCreate}
        isPending={createProduct.isPending}
      />

      {/* Edit dialog */}
      <ProductForm
        open={!!editingProduct}
        onOpenChange={(open) => {
          if (!open) setEditingProduct(null);
        }}
        product={editingProduct}
        onSubmit={handleUpdate}
        isPending={updateProduct.isPending}
      />

      {/* Delete confirmation */}
      <AlertDialog
        open={!!deletingProduct}
        onOpenChange={(open) => {
          if (!open) setDeletingProduct(null);
        }}
      >
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Delete product?</AlertDialogTitle>
            <AlertDialogDescription>
              This will delete &ldquo;{deletingProduct?.name}&rdquo;. This
              action cannot be undone.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>Cancel</AlertDialogCancel>
            <AlertDialogAction
              onClick={handleDelete}
              disabled={deleteProduct.isPending}
            >
              {deleteProduct.isPending ? "Deleting..." : "Delete"}
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </div>
  );
}
