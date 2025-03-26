import {
    Table,
    TableBody,
    TableCaption,
    TableCell,
    TableFooter,
    TableHead,
    TableHeader,
    TableRow,
  } from "@/components/ui/table"


  interface SimpleTableProps {
    headers: Record<"key" | "label", string>[]
    tableRows: {
      [key: string]: any
    }[]
    tableCaption?: string
  }
  
  export default function SimpleTable({headers, tableRows, tableCaption}: SimpleTableProps) {
    return (
      <Table>
        <TableCaption>{tableCaption ?? ""}</TableCaption>
        <TableHeader>
          <TableRow>
            {headers.map((h, i) => (
              <TableHead key={"table-head-"+i}>{h.label}</TableHead>
            ))}
          </TableRow>
        </TableHeader>
        <TableBody>
          {tableRows.map((row, i) => (
            <TableRow key={`table-row-${i}`}>
              {headers.map((h, j) => (
                <TableCell key={`table-cell-row-${i}-col-${j}`}>{row[h.key]}</TableCell>
              ))}
            </TableRow>
          ))}
        </TableBody>
        <TableFooter>
        </TableFooter>
      </Table>
    )
  }
  