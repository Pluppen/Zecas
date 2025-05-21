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

import { deleteDNSRecordById, type DNSRecord } from "@/lib/api/dns"
import { user } from "@/lib/userStore"

import { RemoveItemDialog } from "@/components/remove-item-dialog"
import type { Dispatch, SetStateAction } from "react"

export const getColumns = (setDNSRecords: Dispatch<SetStateAction<DNSRecord[]>>, dnsRecords: DNSRecord[]) => {
    const columnDefinitions: ColumnDef<DNSRecord>[] = [
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
      accessorKey: "record_type",
      header: ({ column }) => {
          return (
            <Button
              variant="ghost"
              onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
            >
              Record Type
              <ArrowUpDown className="ml-2 h-4 w-4" />
            </Button>
          )
        },
      cell: ({row, column}) => (
        <span className="pl-3 block">{row.getValue(column.id)}</span>
      ),
    },
    {
      accessorKey: "record_value",
      header: ({ column }) => {
          return (
            <Button
              variant="ghost"
              onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
            >
              Record Value
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
        const dnsRecord = row.original
  
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
                  navigator.clipboard.writeText(dnsRecord.id ?? "Error occured")
                  toast("Copied dns record ID to clipboard")
                }}
              >
                Copy DNS Record ID
              </DropdownMenuItem>
              <DropdownMenuItem>
                <a target="_blank" href={`/dns-records/${dnsRecord.id}`}>View details</a>
              </DropdownMenuItem>
              <DropdownMenuSeparator />
              <div className="hover:cursor-pointer text-red-500">
                <RemoveItemDialog
                handleSubmit={async () => {
                  const $user = user.get();
                  if ($user?.access_token && dnsRecord.id){
                    const result = await deleteDNSRecordById(dnsRecord.id, $user.access_token);
                    if ("error" in result) {
                      toast("Failed to remove DNS record");
                      return;
                    }
                    const tmpServices = dnsRecords.slice().filter(t => t.id !== dnsRecord.id);
                    setDNSRecords(tmpServices);
                    toast("Successfully removed dns record");
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
