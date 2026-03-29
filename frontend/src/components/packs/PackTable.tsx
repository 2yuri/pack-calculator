import { Pencil, Trash2 } from "lucide-react";
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
import type { Pack } from "@/types";

interface PackTableProps {
  packs: Pack[];
  onEdit: (pack: Pack) => void;
  onDelete: (pack: Pack) => void;
}

const iconBtn = "cursor-pointer transition-colors hover:text-primary";
const iconBtnDanger = "cursor-pointer transition-colors hover:text-destructive";

export function PackTable({ packs, onEdit, onDelete }: PackTableProps) {
  const sorted = [...packs].sort((a, b) => a.size - b.size);

  return (
    <>
      {/* Desktop table */}
      <div className="hidden sm:block">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Pack Size</TableHead>
              <TableHead className="w-24">Actions</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {sorted.map((pack) => (
              <TableRow key={pack.id}>
                <TableCell className="font-medium">
                  {pack.size.toLocaleString()}
                </TableCell>
                <TableCell>
                  <div className="flex gap-1">
                    <Button
                      variant="ghost"
                      size="icon"
                      className={iconBtn}
                      onClick={() => onEdit(pack)}
                    >
                      <Pencil />
                    </Button>
                    <Button
                      variant="ghost"
                      size="icon"
                      className={iconBtnDanger}
                      onClick={() => onDelete(pack)}
                    >
                      <Trash2 />
                    </Button>
                  </div>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </div>

      {/* Mobile cards */}
      <div className="flex flex-col gap-3 sm:hidden">
        {sorted.map((pack) => (
          <Card key={pack.id}>
            <CardContent className="flex items-center justify-between p-4">
              <p className="font-medium">{pack.size.toLocaleString()}</p>
              <div className="flex gap-1">
                <Button
                  variant="ghost"
                  size="icon"
                  className={iconBtn}
                  onClick={() => onEdit(pack)}
                >
                  <Pencil />
                </Button>
                <Button
                  variant="ghost"
                  size="icon"
                  className={iconBtnDanger}
                  onClick={() => onDelete(pack)}
                >
                  <Trash2 />
                </Button>
              </div>
            </CardContent>
          </Card>
        ))}
      </div>
    </>
  );
}
