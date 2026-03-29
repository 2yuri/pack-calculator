import { useState } from "react";
import { Link, useRouterState } from "@tanstack/react-router";
import { Menu, Package, ShoppingCart } from "lucide-react";
import { Button } from "@/components/ui/button";
import {
  Sheet,
  SheetContent,
  SheetHeader,
  SheetTitle,
  SheetTrigger,
} from "@/components/ui/sheet";

const navLinks = [
  { to: "/products" as const, label: "Products", icon: Package },
  { to: "/orders" as const, label: "Orders", icon: ShoppingCart },
];

export function Navbar() {
  const [open, setOpen] = useState(false);
  const routerState = useRouterState();
  const currentPath = routerState.location.pathname;

  return (
    <header className="border-b bg-background">
      <div className="mx-auto flex h-14 max-w-5xl items-center justify-between px-4">
        <Link to="/products" className="text-lg font-semibold">
          Pack Calculator
        </Link>

        {/* Desktop nav */}
        <nav className="hidden gap-1 sm:flex">
          {navLinks.map(({ to, label }) => (
            <Button
              key={to}
              variant={currentPath.startsWith(to) ? "secondary" : "ghost"}
              size="sm"
              render={<Link to={to} />}
            >
              {label}
            </Button>
          ))}
        </nav>

        {/* Mobile hamburger */}
        <Sheet open={open} onOpenChange={setOpen}>
          <SheetTrigger
            render={
              <Button variant="ghost" size="icon" className="sm:hidden" />
            }
          >
            <Menu />
          </SheetTrigger>
          <SheetContent side="left">
            <SheetHeader>
              <SheetTitle>Pack Calculator</SheetTitle>
            </SheetHeader>
            <nav className="mt-4 flex flex-col gap-2">
              {navLinks.map(({ to, label, icon: Icon }) => (
                <Button
                  key={to}
                  variant={currentPath.startsWith(to) ? "secondary" : "ghost"}
                  className="justify-start"
                  render={<Link to={to} />}
                  onClick={() => setOpen(false)}
                >
                  <Icon className="mr-2" />
                  {label}
                </Button>
              ))}
            </nav>
          </SheetContent>
        </Sheet>
      </div>
    </header>
  );
}
