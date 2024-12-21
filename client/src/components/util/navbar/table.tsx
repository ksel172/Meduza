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
import { Button } from "@/components/ui/button"
import { Ellipsis } from "lucide-react"

  
  const shells = [
    {
      name: "Tim Cooks Laptop",
      shellStatus: "Alive",
      compileType: "x64-win",
      targetIp: "10.10.14.12",
    },
    {
      name: "NSA's Systems",
      shellStatus: "Dead",
      compileType: "x64-win",
      targetIp: "10.10.14.13",
    },
    {
      name: "CIA Backdoor",
      shellStatus: "Alive",
      compileType: "x64-win",
      targetIp: "10.10.14.14",
    },
    {
      name: "Another Backdoor",
      shellStatus: "Alive",
      compileType: "x64-win",
      targetIp: "10.10.14.15",
    },
  ]
  
  export function TableComponent() {
    return (
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>Name</TableHead>
            <TableHead>Status</TableHead>
            <TableHead>Compile Type</TableHead>
            <TableHead>Target</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {shells.map((shell) => (
            <TableRow key={shell.targetIp}>
              <TableCell>{shell.name}</TableCell>
              <TableCell className={`font-medium ${shell.shellStatus === "Dead" ? "text-red-600" : "text-green-400"}`}>{shell.shellStatus}</TableCell>
              <TableCell>{shell.compileType}</TableCell>
              <TableCell>{shell.targetIp}</TableCell>
              <TableCell className="text-right m-0 p-0"><Ellipsis /></TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    )
  }
  