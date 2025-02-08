import {
    Table,
    TableBody,
    TableHead,
    TableHeader,
    TableRow,
  } from "@/components/ui/table"

  export function TableComponent({
    headers,
    data,
    renderRow,
  }: {
    headers: string[];
    data: any[];
    renderRow: (item: any) => React.ReactNode;
  }) {
    return (
      <Table>
        <TableHeader>
          <TableRow>
            {headers.map((header) => (
                <TableHead key={header}>{header}</TableHead>
            ))}
          </TableRow>
        </TableHeader>
        <TableBody>
          {data.map((item, index) => (
            <TableRow key={index}>{renderRow(item)}</TableRow>
          ))}
        </TableBody>
      </Table>
    )
  }
  