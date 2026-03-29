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
import type { Pack } from "@/types";

interface PackFormProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  pack?: Pack | null;
  onSubmit: (size: number) => void;
  isPending: boolean;
  error?: string | null;
}

function PackFormContent({
  pack,
  onSubmit,
  isPending,
  error: externalError,
}: Omit<PackFormProps, "open" | "onOpenChange">) {
  const [size, setSize] = useState(pack?.size?.toString() ?? "");
  const [error, setError] = useState("");

  useEffect(() => {
    setSize(pack?.size?.toString() ?? "");
    setError("");
  }, [pack]);

  function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    const parsed = parseInt(size, 10);
    if (!size || isNaN(parsed) || parsed <= 0) {
      setError("Size must be a positive integer");
      return;
    }
    setError("");
    onSubmit(parsed);
  }

  const displayError = externalError || error;

  return (
    <form onSubmit={handleSubmit} className="flex flex-col gap-4 p-4 sm:p-0">
      <div className="flex flex-col gap-2">
        <Label htmlFor="pack-size">Pack Size</Label>
        <Input
          id="pack-size"
          type="number"
          min="1"
          value={size}
          onChange={(e) => setSize(e.target.value)}
          placeholder="e.g. 250"
          autoFocus
        />
        {displayError && (
          <p className="text-sm text-destructive">{displayError}</p>
        )}
      </div>
      <Button type="submit" disabled={isPending} className="w-full">
        {isPending ? "Saving..." : pack ? "Update" : "Add Pack"}
      </Button>
    </form>
  );
}

export function PackForm({
  open,
  onOpenChange,
  pack,
  onSubmit,
  isPending,
  error,
}: PackFormProps) {
  const isDesktop = useMediaQuery("(min-width: 640px)");
  const title = pack ? "Edit Pack" : "Add Pack";

  if (isDesktop) {
    return (
      <Dialog open={open} onOpenChange={onOpenChange}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>{title}</DialogTitle>
          </DialogHeader>
          <PackFormContent
            pack={pack}
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
          <DrawerTitle>{title}</DrawerTitle>
        </DrawerHeader>
        <PackFormContent
          pack={pack}
          onSubmit={onSubmit}
          isPending={isPending}
          error={error}
        />
      </DrawerContent>
    </Drawer>
  );
}
