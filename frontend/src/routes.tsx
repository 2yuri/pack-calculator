import {
  createRouter,
  createRoute,
  createRootRoute,
  redirect,
} from "@tanstack/react-router";
import { AppLayout } from "@/components/layout/AppLayout";
import { ProductsPage } from "@/pages/ProductsPage";
import { PacksPage } from "@/pages/PacksPage";
import { OrdersPage } from "@/pages/OrdersPage";

const rootRoute = createRootRoute({
  component: AppLayout,
});

const indexRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/",
  beforeLoad: () => {
    throw redirect({ to: "/products" });
  },
});

const productsRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/products",
  component: ProductsPage,
});

const packsRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/products/$productId/packs",
  component: PacksPage,
});

const ordersRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/orders",
  component: OrdersPage,
});

const routeTree = rootRoute.addChildren([
  indexRoute,
  productsRoute,
  packsRoute,
  ordersRoute,
]);

export const router = createRouter({ routeTree });

declare module "@tanstack/react-router" {
  interface Register {
    router: typeof router;
  }
}
