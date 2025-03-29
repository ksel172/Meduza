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

import { useToast } from "@/hooks/use-toast"
import { ToastAction } from "@/components/ui/toast"

import axiosInstance from "@/axiosInstance";
import { formatISO } from "date-fns";

import { Label } from "@radix-ui/react-label";
import { Pause, Play, PlayIcon, PlaySquare, PlaySquareIcon } from "lucide-react";

import { Command, CommandInput, CommandList, CommandEmpty, CommandItem, CommandGroup } from "@/components/ui/command";
import { Check, ChevronsUpDown, MoreHorizontal } from "lucide-react"

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
import { useCookies } from "react-cookie";
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuLabel, DropdownMenuSeparator, DropdownMenuTrigger } from "@/components/ui/dropdown-menu";


export default function Listeners() {

  const [date, setDate] = useState();
  const [value, setValue] = useState();
  const [open, setOpen] = useState();
  const [isCreating, setIsCreating] = useState(false);

  const [selectedHostValues, setSelectedHostValues] = React.useState<string[]>([]);
  const [selectedHeaderValues, setSelectedHeaderValues] = React.useState<string[]>([]);
  const [selectedWhitelistValues, setSelectedWhitelistValues] = React.useState<string[]>([]);
  const [selectedBlacklistValues, setSelectedBlacklistValues] = React.useState<string[]>([]);

  const [listenerName, setListenerName] = useState("");
  const [listenerDescription, setListenerDescription] = useState("");
  const [listenerType, setListenerType] = useState("http"); //default

  // HTTP LISTENER OPTIONS
  const [hostIP, setHostIP] = useState("");
  const [hostPort, setHostPort] = useState("");
  const [connectionPort, setConnectionPort] = useState("");
  const [workingHours, setWorkingHours] = useState("9-5");
  const [hostRotation, setHostRotation] = useState("");
  const [userAgent, setUserAgent] = useState("");
  const [sslEnabled, setSslEnabled] = useState(false);
  const [proxyEnabled, setProxyEnabled] = useState(false);
  const [proxyType, setProxyType] = useState("");
  const [proxyPort, setProxyPort] = useState("");
  const [proxyUsername, setProxyUsername] = useState("");
  const [proxyPassword, setProxyPassword] = useState("");
  const [whitelistEnabled, setWhitelistEnabled] = useState(false);
  const [blacklistEnabled, setBlacklistEnabled] = useState(false);


  const [defaultAgents, setDefaultAgents] = useState([
    {
      value: "Mozilla/5.0 (Windows NT 6.1; WOW64; Trident/7.0; rv:11.0) like Gecko",
      label: "Mozilla/5.0 (Windows NT 6.1; WOW64; Trident/7.0; rv:11.0) like Gecko",
    },
    {
      value: "Mozilla/5.0 (Windows NT 5.1; rv:37.0) Gecko/20100101 Firefox/37.0",
      label: "Mozilla/5.0 (Windows NT 5.1; rv:37.0) Gecko/20100101 Firefox/37.0",
    },
    {
      value: "Mozilla/5.0 (Windows NT 6.1; rv:32.0) Gecko/20100101 Firefox/32.0",
      label: "Mozilla/5.0 (Windows NT 6.1; rv:32.0) Gecko/20100101 Firefox/32.0",
    },
    {
      value: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.10; rv:31.0) Gecko/20100101 Firefox/31.0",
      label: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.10; rv:31.0) Gecko/20100101 Firefox/31.0",
    },
    {
      value: "Mozilla/5.0 (X11; Linux i686; rv:40.0) Gecko/20100101 Firefox/40.0",
      label: "Mozilla/5.0 (X11; Linux i686; rv:40.0) Gecko/20100101 Firefox/40.0",
    },
  ])

  const [defaultHeaders, setDefaultHeaders] = useState([
    {
      value: '{"key": "Content-Type", "value": "application/json"}',
      label: '{"key": "Content-Type", "value": "application/json"}',
    },
    {
      value: '{"key": "Connection", "value": "keep-alive"}',
      label: '{"key": "Connection", "value": "keep-alive"}',
    },
    {
      value: '{"key": "Accept", "value": "*/*"}',
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

  const [defaultHosts, setDefaultHosts] = useState([])

  interface Listener {
    name: string;
    listenerStatus: number;
    listenerType: string;
    listenerBind: string;
    startTime: string;
  }

  const [listeners, setListeners] = useState<Listener[]>([]);

  const [cookies, setCookie, removeCookie] = useCookies<any>(["cookie-name"]);
  const { toast } = useToast()

  const agentHeaders = ["Name", "Status", "Listener Type", "Bind", "Start Time"];

  const renderRow = (listener: any) => (
    <>
      <TableCell>{listener.name}</TableCell>
      <TableCell
        className={`font-medium ${
          listener.listenerStatus === "Stopped" || listener.listenerStatus === "Error" ? "text-red-600" : (listener.listenerStatus === "Paused" || listener.listenerStatus === "Processing" ? "text-orange-400" : "text-green-400")
        }`}
      >
        {listener.listenerStatus}
      </TableCell>
      <TableCell>{listener.listenerType}</TableCell>
      <TableCell>{listener.listenerBind}</TableCell>
      <TableCell>{format(listener.startTime, "Pp")}</TableCell>
      {console.log(listener.listenerStatus)}
      <TableCell className="text-right flex flex-row items-center justify-center">

        {(listener.listenerStatus === "Stopped" || listener.listenerStatus === "Paused") ?
          <Play onClick={() => startListener(listener)} size={18} color="#ffffff" strokeWidth={1} fill="#ffffff"  />
          :
          (listener.listenerStatus === "Running" ?
            <Pause onClick={() => stopListener(listener)} size={18} color="#ffffff" strokeWidth={1} fill="#ffffff" />
            :
            null
          )
          
        }
        
        {/* {(listener.listenerStatus === "Stopped" || listener.listenerStatus === "Paused") ?
          console.log("STOPPED OR PAUSED")
          :
          console.log("ASD")
        } */}
        
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button variant="ghost" className="h-8 w-8 p-0">
              <Ellipsis className="h-4 w-4" />
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end">
            <DropdownMenuItem onClick={() => handleEdit(listener)}>Edit</DropdownMenuItem>
            <DropdownMenuItem onClick={() => handleDetails(listener)}>Details</DropdownMenuItem>
            <DropdownMenuSeparator />
            <DropdownMenuItem onClick={() => deleteListener(listener)}>Delete</DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </TableCell>
    </>
  )

  const handleDetails = (listener: any) => {
    console.log('Details', listener)
  }
  
  const handleEdit = (listener: any) => {
    console.log('Edit', listener)
  }

  // const axiosInstance = axios.create({
  //   baseURL: 'http://localhost:8080/api/v1', // Ensure this matches your API base URL
  //   headers: {
  //     'Content-Type': 'application/json', // Set default headers if required
  //   },
  // });

  const fetchListeners = async () => {
    try {
        const url = "/listeners";
        const listenersData = await axiosInstance.get(url,
          {
            headers: { authorization: `Bearer ${cookies.access_token}` }
          }
        );
        console.log(listenersData.data.data);

        if(listenersData.data.data){
          setListeners((prevListeners) => [
            ...prevListeners,
            ...listenersData.data.data.map((listener : any) => (
              {
              name: listener.name || "Unknown Listener",
              listenerStatus: (listener.status === 0 ? "Stopped" : (listener.status === 1 ? "Running" : (listener.status === 2 ? "Paused" : (listener.status === 3 ? "Processing" : "Error")))),
              listenerType: listener.type || "Unknown",
              listenerBind: listener.config.host_bind + ":" + listener.config.port_bind || "N/A",
              startTime: listener.created_at || "N/A",
              id: listener.id
            })),
          ]);
        }
        
    } catch (error) {
        console.log(error);
    }
  };

  const createListener = async () => {
    try{
        const url = '/listeners';
        const { data } = await axiosInstance.post(
            url,
            {
              "type": listenerType,
              "name": listenerName,
              "description": listenerDescription,
              "config": {
                "working_hours": "9-5",
                "hosts": selectedHostValues,
                "host_bind": hostIP,
                "host_rotation": hostRotation,
                "port_bind": hostPort,
                "port_conn": connectionPort,
                "secure": sslEnabled,
                "user_agent": userAgent,
                "headers": selectedHeaderValues,
                "certificate": {
                  "cert_path": "",
                  "key_path": ""
                },
                "whitelist_enabled": whitelistEnabled,
                "whitelist": selectedWhitelistValues,
                "blacklist_enabled": blacklistEnabled,
                "blacklist": selectedBlacklistValues,
                "proxy_settings": {
                  "enabled": proxyEnabled,
                  "type": proxyType,
                  "port": proxyPort,
                  "username": proxyUsername,
                  "password": proxyPassword
                }
              }
            },
            {
                headers: { authorization: `Bearer ${cookies.access_token}` }
            }
        );
        toast({
          variant: "success",
          title: "Listener Creation Successful!",
          description: "You have successfully created a listener.",
          action: (
            <ToastAction altText="undo">Close</ToastAction>
          ),
        })
        location.reload()
    }
    catch(error: any){
        console.log(error);
        toast({
          variant: "destructive",
          title: "Listener Creation Failed...",
          description: ( error.response.data.error ? error.response.data.error : "Ensure the port is within the same range as specified in the docker compose."),
          action: (
            <ToastAction altText="undo">Close</ToastAction>
          ),
        })
    }
  }

  const startListener = async (listener: any) => {
    try{
      console.log(listener)
        const url = `/listeners/${listener.id}/start`;
        const { data } = await axiosInstance.post(
            url,
            {
            },
            {
                headers: { authorization: `Bearer ${cookies.access_token}` }
            }
        );
        toast({
          variant: "success",
          title: "Listener Start Successful!",
          description: "You listener is now listening.",
          action: (
            <ToastAction altText="undo">Close</ToastAction>
          ),
        })
        location.reload()
    }
    catch(error){
        console.log(error);
        toast({
          variant: "destructive",
          title: "Unable to Start Listener...",
          description: "Ensure the port is within the same range as specified in the docker compose.",
          action: (
            <ToastAction altText="undo">Close</ToastAction>
          ),
        })
    }
  }

  const stopListener = async (listener: any) => {
    try{
      console.log(listener)
        const url = `/listeners/${listener.id}/stop`;
        const { data } = await axiosInstance.post(
            url,
            {
            },
            {
                headers: { authorization: `Bearer ${cookies.access_token}` }
            }
        );
        toast({
          variant: "success",
          title: "Listener Start Successful!",
          description: "You listener is now listening.",
          action: (
            <ToastAction altText="undo">Close</ToastAction>
          ),
        })
        location.reload()
    }
    catch(error){
        console.log(error);
        toast({
          variant: "destructive",
          title: "Unable to Start Listener...",
          description: "Ensure the port is within the same range as specified in the docker compose.",
          action: (
            <ToastAction altText="undo">Close</ToastAction>
          ),
        })
    }
  }

  const deleteListener = async (listener: any) => {
    try{
      console.log(listener)
        const url = `/listeners/${listener.id}`;
        const { data } = await axiosInstance.delete(
            url,
            {
                headers: { authorization: `Bearer ${cookies.access_token}` }
            }
        );
        toast({
          variant: "success",
          title: "Your listener has been deleted.",
          description: "You listener is now deleted.",
          action: (
            <ToastAction altText="undo">Close</ToastAction>
          ),
        })
        location.reload()
    }
    catch(error){
        console.log(error);
        toast({
          variant: "destructive",
          title: "Unable to Delete Listener...",
          description: "Please try again later.",
          action: (
            <ToastAction altText="undo">Close</ToastAction>
          ),
        })
    }
  }

  React.useEffect(() => {
    console.log("Hosts Array: ", selectedHostValues)
    console.log("Header Array: ", selectedHeaderValues)
    console.log("Whitelist Array: ", selectedWhitelistValues)
    console.log("Blacklist Array: ", selectedBlacklistValues)
  }, [selectedBlacklistValues, selectedHeaderValues, selectedHostValues, selectedWhitelistValues])

  React.useEffect(() => {
    fetchListeners();
  }, [])

  if(isCreating === false){
    return (
      <div className="w-[calc(100vw-var(--sidebar-width))] h-[100%] gap-4 justify-items-center items-center flex flex-col pb-0 mb-0 p-0">
        <div className="w-[calc(100vw-var(--sidebar-width)-6.5em)] h-[100%] flex flex-col gap-4 justify-items-center min-h-screen pb-4 p-0 m-6">
          <div className="flex flex-row justify-between p-1 bg-secondary w-[19em] h-[2.5em] rounded">
            <Button className="w-[10em] h-[100%]">Table</Button>
            <Button className="bg-transparent text-white w-[10em] h-[100%] hover:bg-[#3d3d3d]" onClick={() => setIsCreating(true)}>Add</Button>
          </div>
          <Card className="w-[calc(100vw-var(--sidebar-width)-6.5em)]">
            <CardContent className="m-0 p-0">
              <TableComponent headers={agentHeaders} data={listeners} renderRow={renderRow} />
            </CardContent>
          </Card>
          <div className="flex flex-row w-[calc(100vw-var(--sidebar-width)-6.5em)] justify-between">
            {/* <p>1 total listener(s) found.</p> */}
            <p>{listeners.length} total listener(s) found.</p>
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
          <div className="flex flex-row justify-between p-1 bg-secondary w-[19em] h-[2.5em] rounded">
            <Button className="bg-transparent text-white w-[10em] h-[100%] hover:bg-[#3d3d3d]" onClick={() => setIsCreating(false)}>Table</Button>
            <Button className="w-[10em] h-[100%]">Add</Button>
          </div>
          <Button className="w-[10em]" onClick={() => createListener()}>+ Create</Button>
        </div>
        <div className="w-[calc(100vw-var(--sidebar-width)-6.5em)] grid grid-cols-3 gap-0 items-start justify-items-end p-0 mb-0 mt-0 border-solid border-[1px] pt-5 rounded-lg">
          <Card className="mx-auto border-none w-[100%] space-x-2">
            <CardContent>
              <div className="space-y-4">
                <div className="space-y-2">
                    <Label htmlFor="email">Name</Label>
                    <Input id="email" type="email" placeholder="Listener Name" value={listenerName} onChange={(e) => setListenerName(e.target.value)} required />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="email">Description</Label>
                  <Input id="email" type="email" placeholder="Listener Description" value={listenerDescription} onChange={(e) => setListenerDescription(e.target.value)} required />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="email">Listener Type</Label>
                  <Select onValueChange={setListenerType}>
                    <SelectTrigger>
                      <SelectValue placeholder="Select Listener Type"/>
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
            </CardContent>
          </Card>

          <Card className="mx-auto border-none w-[100%] space-x-2">
            <CardContent>
              <div className="space-y-4">
                <div className="flex flex-row justify-between items-center gap-2">
                  <div className="space-y-2 w-[100%]">
                    <Label htmlFor="password">Bind IP</Label>
                    <Input placeholder="0.0.0.0" value={hostIP} onChange={(e) => setHostIP(e.target.value)} required />
                  </div>
                  <div className="space-y-2">
                    <Label htmlFor="password">Bind Port</Label>
                    <Input placeholder="80" value={hostPort} onChange={(e) => setHostPort(e.target.value)} required />
                  </div>
                </div>

                <div className="space-y-2">
                  <Label htmlFor="email">Connection Port</Label>
                  <Input id="email" type="email" placeholder="8080" value={connectionPort} onChange={(e) => setConnectionPort(e.target.value)} required />
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
                <Select onValueChange={setHostRotation}>
                  <SelectTrigger>
                    <SelectValue placeholder="Select Rotation Method" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectGroup>
                      <SelectItem value="fallback">Fallback</SelectItem>
                    </SelectGroup>
                  </SelectContent>
                </Select>
              </div>
              <div className="space-y-2">
                  <Label htmlFor="email">User Agent</Label>
                  {/* <MultiSelectPopover initialFrameworks={defaultAgents} selectPlaceholder="Select Agents..." addPlaceholder="Add Agent"/> */}
                  <Input id="email" type="email" placeholder="Mozilla/5.0 (Windows NT 6.1; WOW64; rv:40.0) Gecko/20100101 Firefox/40.0" value={userAgent} onChange={(e) => setUserAgent(e.target.value)} />
              </div>
              <div className="space-y-2">
                  <Label htmlFor="email">Custom Headers</Label>
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
                <Switch id="secure-connection" checked={sslEnabled} onCheckedChange={setSslEnabled} />
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
                      <SelectItem value="cert">Cert</SelectItem>
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
                        <SelectItem value="key">Key</SelectItem>
                      </SelectGroup>
                    </SelectContent>
                  </Select>
              </div>
            </CardContent>
          </Card>

          <Card className="mx-auto border-none w-[100%] space-x-2">
            <CardContent className="flex flex-col gap-5">
              <div className="flex items-center space-x-2">
                <Switch id="proxy-mode" checked={proxyEnabled} onCheckedChange={setProxyEnabled} />
                <Label htmlFor="enable-proxy">Enable Proxy</Label>
              </div>
              <div className="flex flex-row items-center justify-between w-[100%] gap-2">
                <div className="space-y-2 w-[100%]">
                  <Label htmlFor="email">Proxy Type</Label>
                  <Select onValueChange={setProxyType}>
                    <SelectTrigger>
                      <SelectValue placeholder="Select Proxy Type" />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectGroup>
                        <SelectItem value="rio">RIO (Experimental)</SelectItem>
                      </SelectGroup>
                    </SelectContent>
                  </Select>
                </div>
                <div className="space-y-2">
                  <Label htmlFor="email">Proxy Port</Label>
                  <Input id="email" type="email" placeholder="1234" value={proxyPort} onChange={(e) => setProxyPort(e.target.value)} required />
                </div>
              </div>
              <div className="flex flex-row justify-between items-center gap-2">
                <div className="w-[100%]">
                  <Label htmlFor="password">Username</Label>
                  <Input id="password" type="password" placeholder="Batman" value={proxyUsername} onChange={(e) => setProxyUsername(e.target.value)} required />
                </div>
                <div className="w-[100%]">
                  <Label htmlFor="password">Password</Label>
                  <Input id="password" type="password" placeholder="****" value={proxyPassword} onChange={(e) => setProxyPassword(e.target.value)} required />
                </div>
              </div>
            </CardContent>
          </Card>

          <Card className="mx-auto border-none w-[100%] space-x-2">
            <CardContent className="flex flex-col gap-5">
              <div className="flex items-center space-x-2">
                <Switch id="whitelist-option" checked={whitelistEnabled} onCheckedChange={setWhitelistEnabled} />
                <Label htmlFor="enable-proxy">Enable Whitelisted IP's</Label>
              </div>
              <div className="space-y-2">
                  <Label htmlFor="email">Select IP's</Label>
                  {/* <Input id="email" type="email" placeholder="Round Robin" required /> */}
                  <MultiSelectPopover initialFrameworks={whitelistedIPs} selectPlaceholder="Select IPs..." addPlaceholder="Add IP" selectedValues={selectedWhitelistValues} setSelectedValues={setSelectedWhitelistValues}/>
              </div>

              <div className="flex items-center space-x-2">
                <Switch id="blacklist-option" checked={blacklistEnabled} onCheckedChange={setBlacklistEnabled} />
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