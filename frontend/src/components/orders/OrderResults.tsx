import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import type { OrderResult, Product } from "@/types";

interface OrderResultsProps {
  results: OrderResult[];
  products: Product[];
}

export function OrderResults({ results, products }: OrderResultsProps) {
  function getProductName(productId: number) {
    return (
      products.find((p) => p.id === productId)?.name ?? `Product #${productId}`
    );
  }

  return (
    <div className="flex flex-col gap-4">
      <h2 className="text-xl font-semibold">Results</h2>
      <div className="grid gap-4 sm:grid-cols-2">
        {results.map((result) => (
          <Card key={result.product_id}>
            <CardHeader className="pb-3">
              <CardTitle className="text-lg">
                {getProductName(result.product_id)}
              </CardTitle>
            </CardHeader>
            <CardContent className="flex flex-col gap-3">
              <div className="grid grid-cols-3 gap-2 text-sm">
                <div>
                  <p className="text-muted-foreground">Ordered</p>
                  <p className="font-medium">
                    {result.quantity.toLocaleString()}
                  </p>
                </div>
                <div>
                  <p className="text-muted-foreground">Shipping</p>
                  <p className="font-medium">
                    {result.total_items.toLocaleString()}
                  </p>
                </div>
                <div>
                  <p className="text-muted-foreground">Packs</p>
                  <p className="font-medium">
                    {result.total_packs.toLocaleString()}
                  </p>
                </div>
              </div>

              {result.packs.length > 0 && (
                <Table>
                  <TableHeader>
                    <TableRow>
                      <TableHead>Pack Size</TableHead>
                      <TableHead className="text-right">Qty</TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    {result.packs.map((p) => (
                      <TableRow key={p.pack.id}>
                        <TableCell>{p.pack.size.toLocaleString()}</TableCell>
                        <TableCell className="text-right">
                          x{p.quantity}
                        </TableCell>
                      </TableRow>
                    ))}
                  </TableBody>
                </Table>
              )}
            </CardContent>
          </Card>
        ))}
      </div>
    </div>
  );
}
