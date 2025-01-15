"use client"

import { TableComponent } from "@/components/util/navbar/table";
import * as React from "react"
import { useState } from "react";
import ConsoleWidget from "@/components/util/navbar/console";

import { Button } from "@/components/ui/button"
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
  CardDescription
} from "@/components/ui/card"

import { Label } from "@radix-ui/react-label";

import { TableRow, TableCell } from "@/components/ui/table";
import { CalendarIcon, Ellipsis } from "lucide-react";
import { Combobox } from "@/components/util/items/combobox";
import { Input } from "@/components/ui/input";
import { Select, SelectContent, SelectGroup, SelectItem, SelectLabel, SelectTrigger, SelectValue } from "@/components/ui/select";

import { cn } from "@/lib/utils"
import { Calendar } from "@/components/ui/calendar"
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover"

import { format } from "date-fns"
import { Switch } from "@/components/ui/switch";


export default function Listeners() {

  const [date, setDate] = useState();

  const listeners = [
    {
      name: "CIA Listener",
      listenerStatus: "Alive",
      listenerType: "HTTP",
      listenerBind: "0.0.0.0:80",
      startTime: "12/19/2024"
    }
  ]

  const agentHeaders = ["Name", "Status", "Listener Type", "Bind", "Start Time"];

  const renderRow = (listener: any) => (
    <>
      <TableCell>{listener.name}</TableCell>
      <TableCell
        className={`font-medium ${
          listener.listenerStatus === "Dead" ? "text-red-600" : "text-green-400"
        }`}
      >
        {listener.listenerStatus}
      </TableCell>
      <TableCell>{listener.listenerType}</TableCell>
      <TableCell>{listener.listenerBind}</TableCell>
      <TableCell>{listener.startTime}</TableCell>
      <TableCell className="text-right m-0 p-0">
        <Ellipsis />
      </TableCell>
    </>
  );

  return (
    // <div className="w-[calc(100vw-var(--sidebar-width))] h-[100%] flex flex-col gap-4 justify-items-center min-h-screen pb-4 p-0 m-6">
    //   <div className="flex flex-row justify-between p-1 bg-secondary w-[25%] rounded">
    //     <Button className="w-[10em]">Table</Button>
    //     <Button className="bg-transparent text-white w-[10em]">Add</Button>
    //   </div>
    //   <Card className="w-[calc(95vw-var(--sidebar-width))]">
    //     <CardContent className="m-0 p-0">
    //       <TableComponent headers={agentHeaders} data={listeners} renderRow={renderRow} />
    //     </CardContent>
    //   </Card>
    //   <div className="flex flex-row w-[calc(95vw-var(--sidebar-width))] justify-between">
    //     <p>1 total listener(s) found.</p>
    //     <div className="flex flex-row gap-2">
    //       <Button>Previous</Button>
    //       <Button>Next</Button>
    //     </div>
    //   </div>
    // </div>

    <div className="w-[calc(100vw-var(--sidebar-width))] h-[100%] flex flex-col gap-4 justify-items-center min-h-screen pb-4 p-0 m-6">
      <div className="flex flex-row justify-between p-1 bg-secondary w-[25%] rounded">
          <Button className="bg-transparent text-white w-[10em]">Table</Button>
          <Button className="w-[10em]">Add</Button>
      </div>
      <div className="max-w-[calc(100vw-var(--sidebar-width)-3em)] grid grid-cols-2 gap-0 items-start justify-items-end p-0 mb-0 mt-0 border-solid border-2 pt-5 rounded">
        <Card className="mx-auto max-w-sm border-none mb-4">
          <CardContent>
              <div className="space-y-4">
                <div className="space-y-2">
                    <Label htmlFor="email">Name</Label>
                    <Input id="email" type="email" placeholder="CIA Listener" required />
                </div>
                <div className="flex flex-row justify-center items-center gap-2">
                  <div>
                    <Label htmlFor="password">Bind IP</Label>
                    <Input id="password" type="password" placeholder="0.0.0.0" required />
                  </div>
                  <div>
                    <Label htmlFor="password">Bind Port</Label>
                    <Input id="password" type="password" placeholder="80" required />
                  </div>
                </div>
                <div className="space-y-2">
                    <Label htmlFor="email">Connection Port</Label>
                    <Input id="email" type="email" placeholder="8080" required />
                </div>
              </div>
          </CardContent>
        </Card>

        <Card className="mx-auto max-w-sm border-none mb-4">
          <CardContent>
              <div className="space-y-4">
              <div className="space-y-2">
                  <Label htmlFor="email">Rotation Type</Label>
                  {/* <Input id="email" type="email" placeholder="Round Robin" required /> */}
                  <Select>
                    <SelectTrigger>
                      <SelectValue placeholder="Select Rotation Type" />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectGroup>
                        <SelectItem value="apple">Round Robin</SelectItem>
                      </SelectGroup>
                    </SelectContent>
                  </Select>
              </div>
              <div className="flex flex-row justify-center items-center gap-2">
                <Popover>
                  <PopoverTrigger asChild>
                    <Button
                      variant={"outline"}
                      className={cn(
                        "w-[240px] justify-start text-left font-normal",
                        !date && "text-muted-foreground"
                      )}
                    >
                      <CalendarIcon />
                      {date ? format(date, "P") : <span>From</span>}
                    </Button>
                  </PopoverTrigger>
                  <PopoverContent className="w-auto p-0" align="start">
                    <Calendar
                      mode="single"
                      selected={date}
                      onSelect={setDate}
                      initialFocus
                    />
                  </PopoverContent>
                </Popover>
                <Popover>
                  <PopoverTrigger asChild>
                    <Button
                      variant={"outline"}
                      className={cn(
                        "w-[240px] justify-start text-left font-normal",
                        !date && "text-muted-foreground"
                      )}
                    >
                      <CalendarIcon />
                      {date ? format(date, "P") : <span>Until</span>}
                    </Button>
                  </PopoverTrigger>
                  <PopoverContent className="w-auto p-0" align="start">
                    <Calendar
                      mode="single"
                      selected={date}
                      onSelect={setDate}
                      initialFocus
                    />
                  </PopoverContent>
                </Popover>
              </div>
              <div className="space-y-2 flex flex-col">
                  <Label htmlFor="email">Kill Date</Label>
                  <Popover>
                  <PopoverTrigger asChild>
                    <Button
                      variant={"outline"}
                      className={cn(
                        "justify-start text-left font-normal",
                        !date && "text-muted-foreground"
                      )}
                    >
                      <CalendarIcon />
                      {date ? format(date, "PPPP") : <span>Pick a date</span>}
                    </Button>
                  </PopoverTrigger>
                  <PopoverContent className="w-auto p-0" align="start">
                    <Calendar
                      mode="single"
                      selected={date}
                      onSelect={setDate}
                      initialFocus
                    />
                  </PopoverContent>
                </Popover>
              </div>
              </div>
          </CardContent>
        </Card>

        <span className="border-solid border-0 border-b rounded-none h-[1px] w-[100%]" />
        <span className="border-solid border-0 border-b rounded-none h-[1px] w-[100%]" />

        <Card className="mx-auto max-w-sm border-none w-[100%] mt-4">
          <CardContent className="flex flex-col gap-5">
            <div className="flex items-center space-x-2">
              <Switch id="airplane-mode" />
              <Label htmlFor="enable-proxy">Enable Proxy</Label>
            </div>
            <div className="space-y-2">
                <Label htmlFor="email">Enable Proxy</Label>
                {/* <Input id="email" type="email" placeholder="Round Robin" required /> */}
                <Select>
                  <SelectTrigger>
                    <SelectValue placeholder="Select Proxy Type" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectGroup>
                      <SelectItem value="apple">RIO (Experimental)</SelectItem>
                    </SelectGroup>
                  </SelectContent>
                </Select>
            </div>
          </CardContent>
        </Card>

        <Card className="mx-auto max-w-sm border-none w-[100%] mt-4">
          <CardContent className="flex flex-col gap-5">
            <div className="flex items-center space-x-2">
              <Switch id="airplane-mode" />
              <Label htmlFor="enable-proxy">Enable Secure Connection</Label>
            </div>
            <div className="space-y-2">
                <Label htmlFor="email">Certificate File</Label>
                {/* <Input id="email" type="email" placeholder="Round Robin" required /> */}
                <Select>
                  <SelectTrigger>
                    <SelectValue placeholder="cert.cer" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectGroup>
                      <SelectItem value="apple">Cert</SelectItem>
                    </SelectGroup>
                  </SelectContent>
                </Select>
            </div>

            <div>
              <Label htmlFor="email">Key File</Label>
              {/* <Input id="email" type="email" placeholder="Round Robin" required /> */}
              <Select>
                <SelectTrigger>
                  <SelectValue placeholder="key.pem" />
                </SelectTrigger>
                <SelectContent>
                  <SelectGroup>
                    <SelectItem value="apple">Key</SelectItem>
                  </SelectGroup>
                </SelectContent>
              </Select>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>

  );
}