import { useState } from "react";
import { Link, useParams } from "@tanstack/react-router";
import { ArrowLeft, Plus, Layers, Package } from "lucide-react";
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
import { PackTable } from "@/components/packs/PackTable";
import { PackForm } from "@/components/packs/PackForm";
import { BatchPackForm } from "@/components/packs/BatchPackForm";
import { useProduct } from "@/hooks/useProducts";
import {
  usePacks,
  useCreatePack,
  useCreatePacksBatch,
  useUpdatePack,
  useDeletePack,
} from "@/hooks/usePacks";
import { ApiRequestError } from "@/api/client";
import type { Pack } from "@/types";

export function PacksPage() {
  const { productId } = useParams({ from: "/products/$productId/packs" });
  const numericId = Number(productId);

  const { data: product, isLoading: productLoading } = useProduct(numericId);
  const { data: packs, isLoading: packsLoading } = usePacks(numericId);
  const createPack = useCreatePack(numericId);
  const createPacksBatch = useCreatePacksBatch(numericId);
  const updatePack = useUpdatePack(numericId);
  const deletePack = useDeletePack(numericId);

  const [formOpen, setFormOpen] = useState(false);
  const [batchFormOpen, setBatchFormOpen] = useState(false);
  const [editingPack, setEditingPack] = useState<Pack | null>(null);
  const [deletingPack, setDeletingPack] = useState<Pack | null>(null);
  const [formError, setFormError] = useState<string | null>(null);
  const [batchFormError, setBatchFormError] = useState<string | null>(null);

  function handleCreate(size: number) {
    setFormError(null);
    createPack.mutate(size, {
      onSuccess: () => {
        setFormOpen(false);
        toast.success("Pack added");
      },
      onError: (err) => {
        if (err instanceof ApiRequestError && err.status === 409) {
          setFormError(err.message);
        } else {
          toast.error(err instanceof ApiRequestError ? err.message : "Failed to add pack");
        }
      },
    });
  }

  function handleBatchCreate(sizes: number[]) {
    setBatchFormError(null);
    createPacksBatch.mutate(sizes, {
      onSuccess: () => {
        setBatchFormOpen(false);
        toast.success(`${sizes.length} packs added`);
      },
      onError: (err) => {
        if (err instanceof ApiRequestError && err.status === 409) {
          setBatchFormError(err.message);
        } else {
          toast.error(err instanceof ApiRequestError ? err.message : "Failed to add packs");
        }
      },
    });
  }

  function handleUpdate(size: number) {
    if (!editingPack) return;
    setFormError(null);
    updatePack.mutate(
      { packId: editingPack.id, size },
      {
        onSuccess: () => {
          setEditingPack(null);
          toast.success("Pack updated");
        },
        onError: (err) => {
          if (err instanceof ApiRequestError && err.status === 409) {
            setFormError(err.message);
          } else {
            toast.error(err instanceof ApiRequestError ? err.message : "Failed to update pack");
          }
        },
      },
    );
  }

  function handleDelete() {
    if (!deletingPack) return;
    deletePack.mutate(deletingPack.id, {
      onSuccess: () => {
        setDeletingPack(null);
        toast.success("Pack deleted");
      },
      onError: (err) => {
        toast.error(err instanceof ApiRequestError ? err.message : "Failed to delete pack");
      },
    });
  }

  if (productLoading || packsLoading) {
    return (
      <div className="flex flex-col gap-4">
        <Skeleton className="h-6 w-24" />
        <Skeleton className="h-8 w-48" />
        <div className="flex flex-col gap-2">
          {Array.from({ length: 3 }).map((_, i) => (
            <Skeleton key={i} className="h-12 w-full" />
          ))}
        </div>
      </div>
    );
  }

  const isEmpty = !packs || packs.length === 0;

  return (
    <div className="flex flex-col gap-4">
      <Button variant="ghost" size="sm" render={<Link to="/products" />}>
        <ArrowLeft className="mr-2" />
        Back to Products
      </Button>

      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold">
          {product?.name ?? "Product"} — Packs
        </h1>
        <div className="flex gap-2">
          <Button variant="outline" onClick={() => setBatchFormOpen(true)}>
            <Layers className="mr-2" />
            <span className="hidden sm:inline">Add Multiple</span>
          </Button>
          <Button onClick={() => setFormOpen(true)}>
            <Plus className="mr-2" />
            <span className="hidden sm:inline">Add Pack</span>
          </Button>
        </div>
      </div>

      {isEmpty ? (
        <div className="flex flex-col items-center justify-center gap-4 rounded-lg border border-dashed py-12">
          <Package className="size-12 text-muted-foreground" />
          <div className="text-center">
            <p className="text-lg font-medium">No packs configured</p>
            <p className="text-sm text-muted-foreground">
              Add pack sizes to use this product in orders.
            </p>
          </div>
          <Button onClick={() => setFormOpen(true)}>
            <Plus className="mr-2" />
            Add Pack
          </Button>
        </div>
      ) : (
        <PackTable
          packs={packs}
          onEdit={(p) => {
            setFormError(null);
            setEditingPack(p);
          }}
          onDelete={(p) => setDeletingPack(p)}
        />
      )}

      <PackForm
        open={formOpen}
        onOpenChange={(open) => {
          setFormOpen(open);
          if (!open) setFormError(null);
        }}
        onSubmit={handleCreate}
        isPending={createPack.isPending}
        error={formError}
      />

      <BatchPackForm
        open={batchFormOpen}
        onOpenChange={(open) => {
          setBatchFormOpen(open);
          if (!open) setBatchFormError(null);
        }}
        onSubmit={handleBatchCreate}
        isPending={createPacksBatch.isPending}
        error={batchFormError}
      />

      <PackForm
        open={!!editingPack}
        onOpenChange={(open) => {
          if (!open) {
            setEditingPack(null);
            setFormError(null);
          }
        }}
        pack={editingPack}
        onSubmit={handleUpdate}
        isPending={updatePack.isPending}
        error={formError}
      />

      <AlertDialog
        open={!!deletingPack}
        onOpenChange={(open) => {
          if (!open) setDeletingPack(null);
        }}
      >
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Delete pack?</AlertDialogTitle>
            <AlertDialogDescription>
              This will delete pack size {deletingPack?.size.toLocaleString()}.
              This action cannot be undone.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>Cancel</AlertDialogCancel>
            <AlertDialogAction
              onClick={handleDelete}
              disabled={deletePack.isPending}
            >
              {deletePack.isPending ? "Deleting..." : "Delete"}
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </div>
  );
}
