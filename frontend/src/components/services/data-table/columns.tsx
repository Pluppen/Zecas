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

import { deleteServiceById } from "@/lib/api/services"
import { user } from "@/lib/userStore"

import { RemoveItemDialog } from "@/components/remove-item-dialog"

export const getColumns = (setServices: any, services: any) => {
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
      accessorKey: "banner",
      header: ({ column }) => {
          return (
            <Button
              variant="ghost"
              onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
            >
              Banner
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
      accessorKey: "port",
      header: ({ column }) => {
          return (
            <Button
              variant="ghost"
              onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
            >
                Port
              <ArrowUpDown className="ml-2 h-4 w-4" />
            </Button>
          )
        },
      cell: ({row, column}) => (
        <span className="pl-3 uppercase">{row.getValue(column.id)}</span>
      ),
    },
    {
      accessorKey: "protocol",
      header: ({ column }) => {
          return (
            <Button
              variant="ghost"
              onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
            >
                Protocol
              <ArrowUpDown className="ml-2 h-4 w-4" />
            </Button>
          )
        },
      cell: ({row, column}) => (
        <span className="pl-3 uppercase">{row.getValue(column.id)}</span>
      ),
    },
    {
      accessorKey: "service_name",
      header: ({ column }) => {
          return (
            <Button
              variant="ghost"
              onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
            >
                Service
              <ArrowUpDown className="ml-2 h-4 w-4" />
            </Button>
          )
        },
      cell: ({row, column}) => (
        <span className="pl-3 uppercase">{row.getValue(column.id)}</span>
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
        const service = row.original
  
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
                  navigator.clipboard.writeText(service.id)
                  toast("Copied service ID to clipboard")
                }}
              >
                Copy service ID
              </DropdownMenuItem>
              <DropdownMenuItem>
                <a target="_blank" href={`/services/${service.id}`}>View details</a>
              </DropdownMenuItem>
              <DropdownMenuSeparator />
              <div className="hover:cursor-pointer text-red-500">
                <RemoveItemDialog
                handleSubmit={async () => {
                  const $user = user.get();
                  const result = await deleteServiceById(service.id, $user.access_token);
                  if ("error" in result) {
                    toast("Failed to remove service");
                    return;
                  }
                  const tmpServices = services.slice().filter(t => t.id !== service.id);
                  setServices(tmpServices);
                  toast("Successfully removed service");
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
