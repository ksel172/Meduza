import Image from "next/image";
import { Navbar } from "@/components/util/navbar/navbar";
import { TableComponent } from "@/components/util/navbar/table";
import * as React from "react"
import ChatRoom from "@/components/util/navbar/chat";
import ConsoleWidget from "@/components/util/navbar/console";

import { Button } from "@/components/ui/button"
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"

import { ScrollArea } from "@/components/ui/scroll-area"
import { Combobox } from "@/components/util/items/combobox";

import { TableCell } from "@/components/ui/table";
import { Ellipsis } from "lucide-react";

const comboboxOptions = [
  {
    value: "dead",
    label: "Dead",
  },
  {
    value: "alive",
    label: "Alive",
  },
]

export default function Home() {

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

  const agentHeaders = ["Name", "Status", "Compile Type", "Target"];

  const renderRow = (shell: any) => (
    <>
      <TableCell>{shell.name}</TableCell>
      <TableCell
        className={`font-medium ${
          shell.shellStatus === "Dead" ? "text-red-600" : "text-green-400"
        }`}
      >
        {shell.shellStatus}
      </TableCell>
      <TableCell>{shell.compileType}</TableCell>
      <TableCell>{shell.targetIp}</TableCell>
      <TableCell className="text-right m-0 p-0">
        <Ellipsis />
      </TableCell>
    </>
  );

  return (
    <div className="w-[calc(100vw-var(--sidebar-width))] h-[100%] flex flex-col gap-4 justify-items-center items-center min-h-screen pb-4 p-0 mt-6">
      <div className="flex flex-row gap-2 w-[calc(100vw-var(--sidebar-width)-6.5em)]">
        <Combobox options={comboboxOptions} deafultLabel="Select Filter" />
        <Input className="w-[50%]" type="email" placeholder="Search..." />
      </div>
      <Card className="w-[calc(100vw-var(--sidebar-width)-6.5em)]">
        <CardContent className="m-0 p-0">
           <TableComponent headers={agentHeaders} data={shells} renderRow={renderRow} />
        </CardContent>
      </Card>
      <div className="flex flex-row w-[calc(100vw-var(--sidebar-width)-6.5em)] justify-between">
        <p>4 total agent(s) found.</p>
        <div className="flex flex-row gap-2">
          <Button>Previous</Button>
          <Button>Next</Button>
        </div>
      </div>
      <Card className="w-[calc(100vw-var(--sidebar-width)-6.5em)]">
        <CardContent className="m-0 p-0 h-[60vh]">
          <ConsoleWidget />
        </CardContent>
      </Card>
    </div>
  );
}