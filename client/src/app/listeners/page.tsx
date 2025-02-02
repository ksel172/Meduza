"use client"

import { TableComponent } from "@/components/util/navbar/table";
import * as React from "react"
import { useState } from "react";
import ConsoleWidget from "@/components/util/navbar/console";

import { MultiSelectPopover } from "@/components/appendable";

import { Button } from "@/components/ui/button"
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
  CardDescription
} from "@/components/ui/card"

import { Label } from "@radix-ui/react-label";

import { Command, CommandInput, CommandList, CommandEmpty, CommandItem, CommandGroup } from "@/components/ui/command";
import { Check, ChevronsUpDown } from "lucide-react"

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
  const [value, setValue] = useState();
  const [open, setOpen] = useState();
  const [isCreating, setIsCreating] = useState(false);
  const [selectedHostValues, setSelectedHostValues] = React.useState<string[]>([]);
  const [selectedHeaderValues, setSelectedHeaderValues] = React.useState<string[]>([]);
  const [selectedWhitelistValues, setSelectedWhitelistValues] = React.useState<string[]>([]);
  const [selectedBlacklistValues, setSelectedBlacklistValues] = React.useState<string[]>([]);
  const [defaultAgents, setDefaultAgents] = useState([
    {
      value: "Firefox",
      label: "Mozilla/5.0 (Windows NT 6.1; WOW64; Trident/7.0; rv:11.0) like Gecko",
    },
    {
      value: "Firefox 2",
      label: "Mozilla/5.0 (Windows NT 5.1; rv:37.0) Gecko/20100101 Firefox/37.0",
    },
    {
      value: "Firefox 3",
      label: "Mozilla/5.0 (Windows NT 6.1; rv:32.0) Gecko/20100101 Firefox/32.0",
    },
    {
      value: "Intel Mac",
      label: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.10; rv:31.0) Gecko/20100101 Firefox/31.0",
    },
    {
      value: "Linux",
      label: "Mozilla/5.0 (X11; Linux i686; rv:40.0) Gecko/20100101 Firefox/40.0",
    },
  ])

  const [defaultHeaders, setDefaultHeaders] = useState([
    {
      value: "Content Type JSON",
      label: '{"key": "Content-Type", "value": "application/json"}',
    },
    {
      value: "Connection Alive",
      label: '{"key": "Connection", "value": "keep-alive"}',
    },
    {
      value: "Cache Control",
      label: '{"key": "Accept", "value": "*/*"}',
    },
  ])

  const [whitelistedIPs, setWhitelistedIPs] = useState([
    {
      value: "localhost",
      label: 'localhost',
    }
  ])

  const [blacklistedIPs, setBlacklistedIPs] = useState([
    {
      value: "localhost",
      label: 'localhost',
    }
  ])

  const [defaultHosts, setDefaultHosts] = useState([
    {
      value: "localhost",
      label: 'http://localhost:8080',
    }
  ])

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

  React.useEffect(() => {
    console.log("Hosts Array: ", selectedHostValues)
    console.log("Header Array: ", selectedHeaderValues)
    console.log("Whitelist Array: ", selectedWhitelistValues)
    console.log("Blacklist Array: ", selectedBlacklistValues)
  }, [selectedBlacklistValues, selectedHeaderValues, selectedHostValues, selectedWhitelistValues])

  if(isCreating === false){
    return (
      <div className="w-[calc(100vw-var(--sidebar-width))] h-[100%] gap-4 justify-items-center items-center flex flex-col pb-0 mb-0 p-0">
        <div className="w-[calc(100vw-var(--sidebar-width)-6.5em)] h-[100%] flex flex-col gap-4 justify-items-center min-h-screen pb-4 p-0 m-6">
          <div className="flex flex-row justify-between p-1 bg-secondary w-[19em] rounded">
            <Button className="w-[10em]">Table</Button>
            <Button className="bg-transparent text-white w-[10em]" onClick={() => setIsCreating(true)}>Add</Button>
          </div>
          <Card className="w-[calc(100vw-var(--sidebar-width)-6.5em)]">
            <CardContent className="m-0 p-0">
              <TableComponent headers={agentHeaders} data={listeners} renderRow={renderRow} />
            </CardContent>
          </Card>
          <div className="flex flex-row w-[calc(100vw-var(--sidebar-width)-6.5em)] justify-between">
            <p>1 total listener(s) found.</p>
            <div className="flex flex-row gap-2">
              <Button>Previous</Button>
              <Button>Next</Button>
            </div>
          </div>
        </div>
      </div>
    )
  } else{
    return (

      <div className="w-[calc(100vw-var(--sidebar-width))] h-[100%] gap-4 justify-items-center items-center flex flex-col pb-0 mb-0 p-0 mt-6">
        <div className="flex flex-row w-[calc(100vw-var(--sidebar-width)-6.5em)] justify-between">
          <div className="flex flex-row justify-between p-1 bg-secondary w-[19em] rounded">
            <Button className="bg-transparent text-white w-[10em]" onClick={() => setIsCreating(false)}>Table</Button>
            <Button className="w-[10em]">Add</Button>
          </div>
          <Button className="w-[10em]">Create</Button>
        </div>
        <div className="w-[calc(100vw-var(--sidebar-width)-6.5em)] grid grid-cols-3 gap-0 items-start justify-items-end p-0 mb-0 mt-0 border-solid border-2 pt-5 rounded">
          <Card className="mx-auto border-none w-[100%] space-x-2">
            <CardContent>
              <div className="space-y-4">
                <div className="space-y-2">
                    <Label htmlFor="email">Name</Label>
                    <Input id="email" type="email" placeholder="CIA Listener" required />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="email">Description</Label>
                  <Input id="email" type="email" placeholder="Definitely not spyware..." required />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="email">Listener Type</Label>
                  <Select>
                    <SelectTrigger>
                      <SelectValue placeholder="Select Listener Type" />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectGroup>
                        <SelectItem value="http">HTTP/S</SelectItem>
                        <SelectItem value="tcp">TCP</SelectItem>
                        <SelectItem value="udp">UDP</SelectItem>
                        <SelectItem value="smb">SMB</SelectItem>
                        <SelectItem value="winrm">WINRM</SelectItem>
                      </SelectGroup>
                    </SelectContent>
                  </Select>
                </div>
                  {/* <div className="flex items-center space-x-2">
                    <Switch id="ssl-mode" />
                    <Label htmlFor="enable-proxy">Enable Secure Connection</Label>
                  </div> */}
              </div>
            </CardContent>
          </Card>

          <Card className="mx-auto border-none w-[100%] space-x-2">
            <CardContent>
              <div className="space-y-4">
                <div className="flex flex-row justify-between items-center gap-2">
                  <div className="space-y-2 w-[100%]">
                    <Label htmlFor="password">Bind IP</Label>
                    <Input id="password" type="password" placeholder="0.0.0.0" required />
                  </div>
                  <div className="space-y-2">
                    <Label htmlFor="password">Bind Port</Label>
                    <Input id="password" type="password" placeholder="80" required />
                  </div>
                </div>

                <div className="space-y-2">
                  <Label htmlFor="email">Connection Port</Label>
                  <Input id="email" type="email" placeholder="8080" required />
                </div>

                <div className="space-y-2" >
                  <Label htmlFor="email">Working Hours</Label>
                  <div className="flex flex-row justify-center items-center gap-2">
                    <Popover>
                      <PopoverTrigger asChild>
                        <Button variant={"outline"} className={cn("w-[240px] justify-start text-left font-normal", !date && "text-muted-foreground")} >
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
                </div>
              </div>
            </CardContent>
          </Card>

          <Card className="mx-auto border-none w-[100%] space-x-2">
            <CardContent className="flex flex-col gap-5">
              <div className="space-y-2 w-[100%]">
                  <Label htmlFor="email">Hosts Selection</Label>
                  {/* <Input id="email" type="email" placeholder="Round Robin" required /> */}
                  <MultiSelectPopover initialFrameworks={defaultHosts} selectPlaceholder="Select Hosts..." addPlaceholder="Add Hosts" selectedValues={selectedHostValues} setSelectedValues={setSelectedHostValues}/>
              </div>
              <div className="space-y-2 w-[100%]">
                <Label htmlFor="email">Host Rotation</Label>
                <Select>
                  <SelectTrigger>
                    <SelectValue placeholder="Select Rotation Method" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectGroup>
                      <SelectItem value="apple">Fallback</SelectItem>
                    </SelectGroup>
                  </SelectContent>
                </Select>
              </div>
              <div className="space-y-2">
                  <Label htmlFor="email">User Agent</Label>
                  {/* <MultiSelectPopover initialFrameworks={defaultAgents} selectPlaceholder="Select Agents..." addPlaceholder="Add Agent"/> */}
                  <Input id="email" type="email" placeholder="Mozilla/5.0 (Windows NT 6.1; WOW64; rv:40.0) Gecko/20100101 Firefox/40.0" />
              </div>
              <div className="space-y-2">
                  <Label htmlFor="email">Custom Headers</Label>
                  {/* <Input id="email" type="email" placeholder='{"key": "X-Custom-Header", "value": "CustomValue"}' /> */}
                  <MultiSelectPopover initialFrameworks={defaultHeaders} selectPlaceholder="Select Headers..." addPlaceholder="Add Header" selectedValues={selectedHeaderValues} setSelectedValues={setSelectedHeaderValues}/>
              </div>
            </CardContent>
          </Card>

          <span className="border-solid border-0 border-b rounded-none h-[1px] w-[100%] mt-6 mb-6" />
          <span className="border-solid border-0 border-b rounded-none h-[1px] w-[100%] mt-6 mb-6" />
          <span className="border-solid border-0 border-b rounded-none h-[1px] w-[100%] mt-6 mb-6" />

          <Card className="mx-auto border-none w-[100%] space-x-2">
            <CardContent className="flex flex-col gap-5">
              <div className="flex items-center space-x-2">
                <Switch id="airplane-mode" />
                <Label htmlFor="enable-proxy">Enable Secure Connection</Label>
              </div>
              <div className="space-y-2">
                <Label htmlFor="email">Certificate File</Label>
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

          <Card className="mx-auto border-none w-[100%] space-x-2">
            <CardContent className="flex flex-col gap-5">
              <div className="flex items-center space-x-2">
                <Switch id="airplane-mode" />
                <Label htmlFor="enable-proxy">Enable Proxy</Label>
              </div>
              <div className="flex flex-row items-center justify-between w-[100%] gap-2">
                <div className="space-y-2 w-[100%]">
                  <Label htmlFor="email">Proxy Type</Label>
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
                <div className="space-y-2">
                  <Label htmlFor="email">Proxy Port</Label>
                  <Input id="email" type="email" placeholder="1234" required />
                </div>
              </div>
              <div className="flex flex-row justify-between items-center gap-2">
                <div className="w-[100%]">
                  <Label htmlFor="password">Username</Label>
                  <Input id="password" type="password" placeholder="Batman" required />
                </div>
                <div className="w-[100%]">
                  <Label htmlFor="password">Password</Label>
                  <Input id="password" type="password" placeholder="****" required />
                </div>
              </div>
            </CardContent>
          </Card>

          <Card className="mx-auto border-none w-[100%] space-x-2">
            <CardContent className="flex flex-col gap-5">
              <div className="flex items-center space-x-2">
                <Switch id="airplane-mode" />
                <Label htmlFor="enable-proxy">Enable Whitelisted IP's</Label>
              </div>
              <div className="space-y-2">
                  <Label htmlFor="email">Select IP's</Label>
                  {/* <Input id="email" type="email" placeholder="Round Robin" required /> */}
                  <MultiSelectPopover initialFrameworks={whitelistedIPs} selectPlaceholder="Select IPs..." addPlaceholder="Add IP" selectedValues={selectedWhitelistValues} setSelectedValues={setSelectedWhitelistValues}/>
              </div>

              <div className="flex items-center space-x-2">
                <Switch id="airplane-mode" />
                <Label htmlFor="enable-proxy">Enable Blacklisted IP's</Label>
              </div>
              <div className="space-y-2">
                  <Label htmlFor="email">Select IP's</Label>
                  {/* <Input id="email" type="email" placeholder="Round Robin" required /> */}
                  <MultiSelectPopover initialFrameworks={whitelistedIPs} selectPlaceholder="Select IPs..." addPlaceholder="Add IP" selectedValues={selectedBlacklistValues} setSelectedValues={setSelectedBlacklistValues}/>
              </div>
            </CardContent>
          </Card>
        </div>
      </div>

    );
  }
}