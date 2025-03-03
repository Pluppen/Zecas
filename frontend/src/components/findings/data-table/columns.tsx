import { type ColumnDef } from "@tanstack/react-table"
import { MoreHorizontal, ArrowUpDown, Trash, Edit } from "lucide-react"
import { Button } from "@/components/ui/button"
import { toast } from "sonner"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import { Checkbox } from "@/components/ui/checkbox"
import SeverityBadge from "../severity-badge"

import { RemoveItemDialog } from "@/components/remove-item-dialog"
import { removeFinding } from "@/lib/findings"

export const getColumns = (setFindings: any, findings: any) => {
    return [
      {
          id: "select",
          header: ({ table }) => (
            <Checkbox
              checked={
                table.getIsAllPageRowsSelected() ||
                (table.getIsSomePageRowsSelected() && "indeterminate")
              }
              onCheckedChange={(value) => table.toggleAllPageRowsSelected(!!value)}
              aria-label="Select all"
            />
          ),
          cell: ({ row }) => (
            <Checkbox
              checked={row.getIsSelected()}
              onCheckedChange={(value) => row.toggleSelected(!!value)}
              aria-label="Select row"
            />
          ),
          enableSorting: false,
          enableHiding: false,
        },
    {
      accessorKey: "severity",
      cell: ({row, column}) => (
        <SeverityBadge severity={row.getValue(column.id)} className="capitalize ml-3" />
      ),
      header: ({ column }) => {
          return (
            <Button
              variant="ghost"
              onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
            >
              Severity
              <ArrowUpDown className="ml-2 h-4 w-4" />
            </Button>
          )
        },
    },
    {
      accessorKey: "title",
      header: ({ column }) => {
          return (
            <Button
              variant="ghost"
              onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
            >
              Title
              <ArrowUpDown className="ml-2 h-4 w-4" />
            </Button>
          )
        },
      cell: ({row, column}) => (
        <span className="pl-3">{row.getValue(column.id)}</span>
      ),
    },
    {
      accessorKey: "target_value",
      header: ({ column }) => {
          return (
            <Button
              variant="ghost"
              onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
            >
              Target
              <ArrowUpDown className="ml-2 h-4 w-4" />
            </Button>
          )
        },
      cell: ({row, column}) => (
        <span className="pl-3">{row.getValue(column.id)}</span>
      ),
    },
    {
      accessorKey: "finding_type",
      header: ({ column }) => {
          return (
            <Button
              variant="ghost"
              onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
            >
              Finding Type
              <ArrowUpDown className="ml-2 h-4 w-4" />
            </Button>
          )
        },
      cell: ({row, column}) => (
        <span className="capitalize pl-3">{row.getValue(column.id)}</span>
      ),
    },
    {
      accessorKey: "discovered_at",
      header: ({ column }) => {
          return (
            <Button
              variant="ghost"
              onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
            >
              Time Generated
              <ArrowUpDown className="ml-2 h-4 w-4" />
            </Button>
          )
        },
      cell: ({row, column}) => (
        <span className="pl-3">{new Date(row.getValue(column.id)).toLocaleString()}</span>
      ),
    },
    {
      accessorKey: "verified",
      header: ({ column }) => {
          return (
            <Button
              variant="ghost"
              onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
            >
              Verified
              <ArrowUpDown className="ml-2 h-4 w-4" />
            </Button>
          )
        },
      cell: ({row, column}) => (
        <span className="pl-3">{row.getValue(column.id) ? "Yes" : "No"}</span>
      ),
    },
    {
      id: "actions",
      cell: ({ row }) => {
        const finding = row.original
  
        return (
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="ghost" className="h-8 w-8 p-0">
                <span className="sr-only">Open menu</span>
                <MoreHorizontal className="h-4 w-4" />
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end">
              <DropdownMenuLabel>Actions</DropdownMenuLabel>
              <DropdownMenuItem
                className="hover:cursor-pointer"
                onClick={() => {
                  navigator.clipboard.writeText(finding.id)
                  toast("Copied finding ID to clipboard")
                }}
              >
                Copy finding ID
              </DropdownMenuItem>
              <DropdownMenuItem>
                <a target="_blank" href={`/findings/${finding.id}`}>View details</a>
              </DropdownMenuItem>
              <DropdownMenuSeparator />
              <DropdownMenuItem
                className="hover:cursor-pointer"
              ><Edit />Edit</DropdownMenuItem>
              <div className="hover:cursor-pointer text-red-500">
                <RemoveItemDialog
                handleSubmit={() => {
                  removeFinding(finding.id).then(result => {
                    if ("error" in result) {
                      toast(result.error);
                    }
                    const findingsTmp = [...findings].filter(f => f.id !== finding.id);
                    setFindings(findingsTmp);
                    toast(result.message);
                  })
                }}
                  button={
                    <>
                      <Trash  color="red" /> Remove
                    </>
                  }
                />
              </div>
            </DropdownMenuContent>
          </DropdownMenu>
        )
      },
      enableSorting: false,
      enableHiding: false,
    },
  ]
}
