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

import { deleteApplicationById } from "@/lib/api/applications"
import { user } from "@/lib/userStore"

import { RemoveItemDialog } from "@/components/remove-item-dialog"

export const getColumns = (setApplications: any, applications: any) => {
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
      accessorKey: "name",
      header: ({ column }) => {
          return (
            <Button
              variant="ghost"
              onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
            >
              Name
              <ArrowUpDown className="ml-2 h-4 w-4" />
            </Button>
          )
        },
      cell: ({row, column}) => (
        <span className="pl-3 block">{row.getValue(column.id)}</span>
      ),
    },
    {
      accessorKey: "description",
      header: ({ column }) => {
          return (
            <Button
              variant="ghost"
              onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
            >
              Description
              <ArrowUpDown className="ml-2 h-4 w-4" />
            </Button>
          )
        },
      cell: ({row, column}) => (
        <span className="pl-3 block">{row.getValue(column.id)}</span>
      ),
    },
    {
      accessorKey: "type",
      header: ({ column }) => {
          return (
            <Button
              variant="ghost"
              onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
            >
              Application Type
              <ArrowUpDown className="ml-2 h-4 w-4" />
            </Button>
          )
        },
      cell: ({row, column}) => (
        <span className="pl-3 block">{row.getValue(column.id)}</span>
      ),
    },
    {
      accessorKey: "version",
      header: ({ column }) => {
          return (
            <Button
              variant="ghost"
              onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
            >
                Version
              <ArrowUpDown className="ml-2 h-4 w-4" />
            </Button>
          )
        },
      cell: ({row, column}) => (
        <span className="pl-3 block">{row.getValue(column.id)}</span>
      ),
    },
    {
      accessorKey: "target",
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
      accessorKey: "created_at",
      header: ({ column }) => {
          return (
            <Button
              variant="ghost"
              onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
            >
                Discovered
              <ArrowUpDown className="ml-2 h-4 w-4" />
            </Button>
          )
        },
      cell: ({row, column}) => (
        <span className="block pl-3">{new Date(row.getValue(column.id)).toLocaleString()}</span>
      ),
    },
    {
      id: "actions",
      cell: ({ row }) => {
        const application = row.original
  
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
                  navigator.clipboard.writeText(application.id)
                  toast("Copied application ID to clipboard")
                }}
              >
                Copy application ID
              </DropdownMenuItem>
              <DropdownMenuItem>
                <a target="_blank" href={`/applications/${application.id}`}>View details</a>
              </DropdownMenuItem>
              <DropdownMenuSeparator />
              <div className="hover:cursor-pointer text-red-500">
                <RemoveItemDialog
                handleSubmit={async () => {
                  const $user = user.get();
                  const result = await deleteApplicationById(application.id, $user.access_token);
                  if ("error" in result) {
                    toast("Failed to remove application");
                    return;
                  }
                  const tmpApplications = applications.slice().filter(t => t.id !== application.id);
                  setApplications(tmpApplications);
                  toast("Successfully removed application");
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
