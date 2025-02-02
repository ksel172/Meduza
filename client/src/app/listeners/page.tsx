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

  const [date, setDate] = useState<Date | undefined>(undefined);
  const [value, setValue] = useState();
  const [open, setOpen] = useState();
  const [isCreating, setIsCreating] = useState(false);
  
  const [isSecureConnEnabled, setIsSecureConnEnabled] = useState(false);
  const handleSwitchToggle = (checked: boolean) => {
    setIsSecureConnEnabled(checked);
  };

  const [isProxyEnabled, setIsProxyEnabled] = useState(false);
  const handleProxySwitchToggle = (checked: boolean) => {
    setIsProxyEnabled(checked);
  };

  const [isWhitelistedEnabled, setIsWhitelistedEnabled] = useState(false);
  const handleWhitelistedSwitchToggle = (checked: boolean) => {
    setIsWhitelistedEnabled(checked);
  };

  const [isBlacklistedEnabled, setIsBlacklistedEnabled] = useState(false);
  const handleBlacklistedSwitchToggle = (checked: boolean) => {
    setIsBlacklistedEnabled(checked);
  };
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
          <div className="flex flex-row justify-between p-1 bg-secondary w-[19em] border- rounded">
              <Button className="bg-transparent text-white w-[10em]" onClick={() => setIsCreating(false)}>Table</Button>
              <Button className="w-[10em]" onClick={() => setIsCreating(true)}>Add</Button>
          </div>
          <Button className="w-[10em]">Create</Button>
        </div>
        <div className="w-[calc(100vw-var(--sidebar-width)-6.5em)] grid grid-cols-3 gap-0 items-start justify-items-end p-0 mb-0 mt-0 border-solid border-[1px] pt-5 rounded-lg">
  <Card className="mx-auto max-w-sm border-none mb-4 w-[100%] mt-4">
    <CardContent>
      <div className="space-y-4">
        {/* Name Input */}
        <div className="space-y-2">
          <Label htmlFor="name">Name</Label>
          <div className="relative">
            <Input
              id="name"
              type="text"
              placeholder=""
              required
              className="pr-8"
            />
            <span className="absolute right-2 top-1/2 transform -translate-y-1/2 text-white text-xl">*</span>
          </div>
        </div>

        {/* Description Input */}
        <div className="space-y-2">
          <Label htmlFor="description">Description</Label>
          <div className="relative">
            <Input
              id="description"
              type="text"
              placeholder=""
              required
              className="pr-8"
            />
          </div>
        </div>

        {/* Listener Type Select */}
        <div className="space-y-2">
          <Label htmlFor="listener-type">Listener Type</Label>
          <div className="relative">
            <Select>
              <SelectTrigger>
                <SelectValue placeholder="" />
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
        </div>
      </div>
    </CardContent>
  </Card>

  <Card className="mx-auto max-w-sm border-none mb-4 w-[100%] mt-4">
    <CardContent>
      <div className="space-y-4">
        <div className="space-y-2">
          {/* Labels Row */}
          <div className="flex justify-between">
            <div className="w-[200px]">
              <Label htmlFor="bind-ip">Bind IP</Label>
            </div>
            <div className="w-[120px] text-right">
              <Label htmlFor="bind-port">Bind Port</Label>
            </div>
          </div>
          
          {/* Inputs Row */}
          <div className="flex items-center gap-2">
            <div className="w-[250px] relative">
              <Input
                id="bind-ip"
                type="text"
                placeholder=""
                required
                className="pr-8"
              />
              <span className="absolute right-2 top-1/2 transform -translate-y-1/2 text-white text-xl">*</span>
            </div>
            <span className="text-white text-xl">:</span>
            <div className="w-[100px] relative">
              <Input
                id="bind-port"
                type="text"
                placeholder=""
                required
                className="pr-8"
              />
              <span className="absolute right-2 top-1/2 transform -translate-y-1/2 text-white text-xl">*</span>
            </div>
          </div>
        </div>

        {/* Connection Port */}
        <div className="space-y-2">
          <Label htmlFor="connection-port">Connection Port</Label>
          <div className="relative">
            <Input
              id="connection-port"
              type="text"
              placeholder=""
              required
              className="pr-8"
            />
          </div>
        </div>

        {/* Working Hours */}
        <div className="space-y-2">
          <Label htmlFor="working-hours">Working Hours</Label>
          <div className="flex flex-row justify-center items-center gap-2">
            <Popover>
              <PopoverTrigger asChild>
                <Button
                  variant={"outline"}
                  className="w-[240px] justify-start text-left font-normal"
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
                  className="w-[240px] justify-start text-left font-normal"
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


  <Card className="mx-auto max-w-sm border-none mb-4 w-[100%] mt-4">
    <CardContent className="flex flex-col gap-5">
      <div className="space-y-2">
        <Label htmlFor="email">Hosts Selection</Label>
        <MultiSelectPopover initialFrameworks={defaultHosts} selectPlaceholder="" addPlaceholder="" />
      </div>
      <div className="space-y-2">
        <Label htmlFor="email">Host Rotation</Label>
        <Select>
          <SelectTrigger>
            <SelectValue placeholder="" />
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
        <Input id="email" type="email" /*placeholder="Mozilla/5.0 (Windows NT 6.1; WOW64; rv:40.0) Gecko/20100101 Firefox/40.0"*/ />
      </div>
      <div className="space-y-2">
        <Label htmlFor="email">Custom Headers</Label>
        <MultiSelectPopover initialFrameworks={defaultHeaders} selectPlaceholder="" addPlaceholder="Add Header" />
      </div>
    </CardContent>
  </Card>

          <span className="border-solid border-0 border-b rounded-none h-[1px] w-[100%]" />
          <span className="border-solid border-0 border-b rounded-none h-[1px] w-[100%]" />
          <span className="border-solid border-0 border-b rounded-none h-[1px] w-[100%]" />

          <Card className="mx-auto max-w-sm border-none w-[100%] mt-4">
      <CardContent className="flex flex-col gap-5">
        {/* Secure Connection Switch */}
        <div className="flex items-center space-x-2">
          <Switch
            id="enable-secure-connection"
            checked={isSecureConnEnabled}
            onCheckedChange={handleSwitchToggle}
          />
          <Label htmlFor="enable-secure-connection">Enable Secure Connection</Label>
        </div>

        {/* Certificate File Select */}
        <div className="space-y-2">
          <Label htmlFor="certificate-file" className={isSecureConnEnabled ? '' : 'opacity-50'}>
            Certificate File
          </Label>
          <Select disabled={!isSecureConnEnabled}>
            <SelectTrigger className={isSecureConnEnabled ? '' : 'opacity-50 pointer-events-none'}>
              <SelectValue placeholder="(.cer)" />
            </SelectTrigger>
            <SelectContent>
              <SelectGroup>
                <SelectItem value="apple">Cert</SelectItem>
              </SelectGroup>
            </SelectContent>
          </Select>
        </div>

        {/* Key File Select */}
        <div className="space-y-2">
          <Label htmlFor="key-file" className={isSecureConnEnabled ? '' : 'opacity-50'}>
            Key File
          </Label>
          <Select disabled={!isSecureConnEnabled}>
            <SelectTrigger className={isSecureConnEnabled ? '' : 'opacity-50 pointer-events-none'}>
              <SelectValue placeholder="(.pem)" />
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


    <Card className="mx-auto max-w-sm border-none w-[100%] mt-4">
  <CardContent className="flex flex-col gap-5">
    {/* Proxy Switch */}
    <div className="flex items-center space-x-2">
      <Switch
        id="enable-proxy-switch"
        checked={isProxyEnabled}
        onCheckedChange={handleProxySwitchToggle}
      />
      <Label htmlFor="enable-proxy-switch">Enable Proxy</Label>
    </div>

    
    <div className="flex flex-row space-x-2">
      {/* Proxy Type Select */}
      <div className="w-full space-y-2">
        <Label htmlFor="proxy-type" className={isProxyEnabled ? '' : 'opacity-50'}>
          Proxy Type
        </Label>
        <Select disabled={!isProxyEnabled}>
          <SelectTrigger
            className={`${isProxyEnabled ? '' : 'opacity-50 pointer-events-none'}`}
          >
            <SelectValue placeholder="" />
          </SelectTrigger>
          <SelectContent>
            <SelectGroup>
              <SelectItem value="apple">RIO (Experimental)</SelectItem>
            </SelectGroup>
          </SelectContent>
        </Select>
      </div>

      {/* Proxy Port Input */}
      <div className="w-[150px] space-y-2">
        <Label htmlFor="proxy-port" className={isProxyEnabled ? '' : 'opacity-50'}>
          Proxy Port
        </Label>
        <Input
          id="proxy-port"
          type="number"
          placeholder=""
          required
          disabled={!isProxyEnabled}
          className={`${isProxyEnabled ? '' : 'opacity-50 pointer-events-none'}`}
        />
      </div>
    </div>

    {/* Username and Password Inputs */}
    <div className="flex flex-row justify-center items-center gap-2">
      <div className="space-y-2">
        <Label htmlFor="username" className={isProxyEnabled ? '' : 'opacity-50'}>
          Username
        </Label>
        <Input
          id="username"
          type="text"
          placeholder=""
          required
          disabled={!isProxyEnabled}
          className={`${isProxyEnabled ? '' : 'opacity-50 pointer-events-none'}`}
        />
      </div>
      <div className="space-y-2">
        <Label htmlFor="password" className={isProxyEnabled ? '' : 'opacity-50'}>
          Password
        </Label>
        <Input
          id="password"
          type="password"
          placeholder=""
          required
          disabled={!isProxyEnabled}
          className={`${isProxyEnabled ? '' : 'opacity-50 pointer-events-none'}`}
        />
      </div>
    </div>
  </CardContent>
</Card>


<Card className="mx-auto max-w-sm border-none w-[100%] mt-4">
  <CardContent className="flex flex-col gap-5">
    {/* Whitelisted IPs Switch */}
    <div className="flex items-center space-x-2">
      <Switch
        id="whitelisted-switch"
        checked={isWhitelistedEnabled}
        onCheckedChange={handleWhitelistedSwitchToggle}
      />
      <Label htmlFor="whitelisted-switch">Enable Whitelisted IP's</Label>
    </div>
    <div className="space-y-2">
      <Label
        htmlFor="whitelisted-ips"
        className={isWhitelistedEnabled ? '' : 'opacity-50'}
      >
        Select IP's
      </Label>
      <div
        className={`${!isWhitelistedEnabled ? 'opacity-50 pointer-events-none' : ''}`}
      >
        <MultiSelectPopover
          initialFrameworks={whitelistedIPs}
          selectPlaceholder=""
          addPlaceholder="Add IP"
        />
      </div>
    </div>

    {/* Blacklisted IPs Switch */}
    <div className="flex items-center space-x-2">
      <Switch
        id="blacklisted-switch"
        checked={isBlacklistedEnabled}
        onCheckedChange={handleBlacklistedSwitchToggle}
      />
      <Label htmlFor="blacklisted-switch">Enable Blacklisted IP's</Label>
    </div>
    <div className="space-y-2">
      <Label
        htmlFor="blacklisted-ips"
        className={isBlacklistedEnabled ? '' : 'opacity-50'}
      >
        Select IP's
      </Label>
      <div
        className={`${!isBlacklistedEnabled ? 'opacity-50 pointer-events-none' : ''}`}
      >
        <MultiSelectPopover
          initialFrameworks={blacklistedIPs}
          selectPlaceholder=""
          addPlaceholder="Add IP"
        />
      </div>
    </div>
  </CardContent>
</Card>

        </div>
      </div>

    );
  }
}