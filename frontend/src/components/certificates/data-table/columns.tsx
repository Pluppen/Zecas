import { type ColumnDef } from "@tanstack/react-table"
import { MoreHorizontal, ArrowUpDown, Trash } from "lucide-react"
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

import { deleteCertificateById, type Certificate } from "@/lib/api/certificates"
import { user } from "@/lib/userStore"

import { RemoveItemDialog } from "@/components/remove-item-dialog"
import type { Dispatch, SetStateAction } from "react"

export const getColumns = (setCertificates: Dispatch<SetStateAction<Certificate[]>>, certificates: Certificate[]) => {
    const columnDefinitions: ColumnDef<Certificate>[] = [
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
        <span className="pl-3 block">{row.getValue(column.id)}</span>
      ),
    },
    {
      accessorKey: "domain",
      header: ({ column }) => {
          return (
            <Button
              variant="ghost"
              onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
            >
              Domain
              <ArrowUpDown className="ml-2 h-4 w-4" />
            </Button>
          )
        },
      cell: ({row, column}) => (
        <span className="pl-3 block">{row.getValue(column.id)}</span>
      ),
    },
    {
      accessorKey: "issuer",
      header: ({ column }) => {
          return (
            <Button
              variant="ghost"
              onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
            >
                Issuer
              <ArrowUpDown className="ml-2 h-4 w-4" />
            </Button>
          )
        },
      cell: ({row, column}) => (
        <span className="pl-3 block">{row.getValue(column.id)}</span>
      ),
    },
    {
      accessorKey: "issued_at",
      header: ({ column }) => {
          return (
            <Button
              variant="ghost"
              onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
            >
              Issued
              <ArrowUpDown className="ml-2 h-4 w-4" />
            </Button>
          )
        },
      cell: ({row, column}) => (
        <span className="pl-3 block">{row.getValue(column.id)}</span>
      ),
    },
    {
      accessorKey: "expires_at",
      header: ({ column }) => {
          return (
            <Button
              variant="ghost"
              onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
            >
              Expires
              <ArrowUpDown className="ml-2 h-4 w-4" />
            </Button>
          )
        },
      cell: ({row, column}) => (
        <span className="pl-3 block">{row.getValue(column.id)}</span>
      ),
    },
    {
      id: "actions",
      cell: ({ row }) => {
        const certificate = row.original
  
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
                  navigator.clipboard.writeText(certificate.id ?? "Error occured")
                  toast("Copied certificate ID to clipboard")
                }}
              >
                Copy Certificate ID
              </DropdownMenuItem>
              <DropdownMenuItem>
                <a target="_blank" href={`/certificates/${certificate.id}`}>View details</a>
              </DropdownMenuItem>
              <DropdownMenuSeparator />
              <div className="hover:cursor-pointer text-red-500">
                <RemoveItemDialog
                handleSubmit={async () => {
                  const $user = user.get();
                  if ($user?.access_token && certificate.id){
                    const result = await deleteCertificateById(certificate.id, $user.access_token);
                    if ("error" in result) {
                      toast("Failed to remove certificate");
                      return;
                    }
                    const tmpServices = certificates.slice().filter(t => t.id !== certificate.id);
                    setCertificates(tmpServices);
                    toast("Successfully removed certificate");
                  }
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
  return columnDefinitions
}
