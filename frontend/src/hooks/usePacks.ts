import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import {
  getPacks,
  createPack,
  createPacksBatch,
  updatePack,
  deletePack,
} from "@/api/packs";

export function usePacks(productId: number) {
  return useQuery({
    queryKey: ["products", productId, "packs"],
    queryFn: () => getPacks(productId),
  });
}

export function useCreatePack(productId: number) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (size: number) => createPack(productId, size),
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: ["products", productId, "packs"],
      });
    },
  });
}

export function useCreatePacksBatch(productId: number) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (sizes: number[]) => createPacksBatch(productId, sizes),
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: ["products", productId, "packs"],
      });
    },
  });
}

export function useUpdatePack(productId: number) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ packId, size }: { packId: number; size: number }) =>
      updatePack(packId, size),
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: ["products", productId, "packs"],
      });
    },
  });
}

export function useDeletePack(productId: number) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (packId: number) => deletePack(packId),
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: ["products", productId, "packs"],
      });
    },
  });
}
