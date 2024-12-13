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

  
  const shells = [
    {
      shellId: "#1230",
      shellStatus: "Alive",
      upTime: "12 Days",
      targetIp: "10.10.14.12",
    },
    {
        shellId: "#1231",
        shellStatus: "dead",
        upTime: "12 Days",
        targetIp: "10.10.14.12",
    },
    {
        shellId: "#1232",
        shellStatus: "Alive",
        upTime: "12 Days",
        targetIp: "10.10.14.12",
    },
    {
        shellId: "#1233",
        shellStatus: "Alive",
        upTime: "12 Days",
        targetIp: "10.10.14.12",
    },
  ]
  
  export function TableComponent() {
    return (
      <Table>
        <TableCaption>A list of your recent invoices.</TableCaption>
        <TableHeader>
          <TableRow>
            <TableHead className="w-[100px]">Target</TableHead>
            <TableHead>Status</TableHead>
            <TableHead>Up Time</TableHead>
            <TableHead>Options</TableHead>
            <TableHead className="text-right">Action</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {shells.map((shell) => (
            <TableRow key={shell.shellId}>
              <TableCell>{shell.targetIp}</TableCell>
              <TableCell className={`font-medium ${shell.shellStatus === "dead" ? "text-red-600" : "text-green-400"}`}>{shell.shellStatus}</TableCell>
              <TableCell>{shell.upTime}</TableCell>
              <TableCell>modify</TableCell>
              <TableCell className="text-right"><Button variant="outline">Interact</Button>
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    )
  }
  