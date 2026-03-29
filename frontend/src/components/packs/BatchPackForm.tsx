import { useState } from "react";
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

interface BatchPackFormProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onSubmit: (sizes: number[]) => void;
  isPending: boolean;
  error?: string | null;
}

function BatchPackFormContent({
  onSubmit,
  isPending,
  error: externalError,
}: Omit<BatchPackFormProps, "open" | "onOpenChange">) {
  const [input, setInput] = useState("");
  const [error, setError] = useState("");

  function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    const parts = input
      .split(",")
      .map((s) => s.trim())
      .filter((s) => s !== "");

    if (parts.length === 0) {
      setError("Enter at least one pack size");
      return;
    }

    const sizes: number[] = [];
    for (const part of parts) {
      const n = parseInt(part, 10);
      if (isNaN(n) || n <= 0) {
        setError(`"${part}" is not a valid positive integer`);
        return;
      }
      sizes.push(n);
    }

    setError("");
    onSubmit(sizes);
  }

  const displayError = externalError || error;

  return (
    <form onSubmit={handleSubmit} className="flex flex-col gap-4 p-4 sm:p-0">
      <div className="flex flex-col gap-2">
        <Label htmlFor="batch-sizes">Pack Sizes</Label>
        <Input
          id="batch-sizes"
          value={input}
          onChange={(e) => setInput(e.target.value)}
          placeholder="e.g. 250, 500, 1000"
          autoFocus
        />
        <p className="text-sm text-muted-foreground">
          Comma-separated list of sizes
        </p>
        {displayError && (
          <p className="text-sm text-destructive">{displayError}</p>
        )}
      </div>
      <Button type="submit" disabled={isPending} className="w-full">
        {isPending ? "Adding..." : "Add Packs"}
      </Button>
    </form>
  );
}

export function BatchPackForm({
  open,
  onOpenChange,
  onSubmit,
  isPending,
  error,
}: BatchPackFormProps) {
  const isDesktop = useMediaQuery("(min-width: 640px)");

  if (isDesktop) {
    return (
      <Dialog open={open} onOpenChange={onOpenChange}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Add Multiple Packs</DialogTitle>
          </DialogHeader>
          <BatchPackFormContent
            onSubmit={onSubmit}
            isPending={isPending}
            error={error}
          />
        </DialogContent>
      </Dialog>
    );
  }

  return (
    <Drawer open={open} onOpenChange={onOpenChange}>
      <DrawerContent>
        <DrawerHeader>
          <DrawerTitle>Add Multiple Packs</DrawerTitle>
        </DrawerHeader>
        <BatchPackFormContent
          onSubmit={onSubmit}
          isPending={isPending}
          error={error}
        />
      </DrawerContent>
    </Drawer>
  );
}
