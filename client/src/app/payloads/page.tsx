"use client"

import { TableComponent } from "@/components/util/navbar/table";
import * as React from "react"
import { useState } from "react";
import ConsoleWidget from "@/components/util/navbar/console";

import { MultiSelectPopover } from "@/components/appendable";
import { saveAs, encodeBase64 } from '@progress/kendo-file-saver';

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

import axios from "axios";

import { formatISO } from "date-fns";

import { Label } from "@radix-ui/react-label";
import { Download, Pause, Play, PlayIcon, PlaySquare, PlaySquareIcon, Trash2 } from "lucide-react";

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

import axiosInstance from "@/axiosInstance";


export default function Payloads() {

  const [date, setDate] = useState();
  const [isCreating, setIsCreating] = useState(false);

  const [payloadName, setPayloadName] = useState("");
  const [payloadArchitecture, setPayloadArchitecture] = useState("win-x64"); //default
  const [selectedListener, setSelectedListener] = useState("");
  const [sleepValue, setSleepValue] = useState(0);
  const [jitterValue, setJitterValue] = useState(0);

  const [selfContained, setSelfContained] = useState(false);
  const [singleFile, setSingleFile] = useState(true);

  interface Listener {
    name: string;
    listenerStatus: number;
    listenerType: string;
    listenerBind: string;
    startTime: string;
    id: string;
  }

  interface Payload {
    name: string;
    payloadArch: string;
    payloadListener: string;
    startTime: string;
    id: string;
  }

  const [listeners, setListeners] = useState<Listener[]>([]);
  const [payloads, setPayloads] = useState<Payload[]>([]);

  const [cookies, setCookie, removeCookie] = useCookies(["cookie-name"]);
  const { toast } = useToast()

  const agentHeaders = ["Name", "Architecture", "Listener Type", "Listener Name", "Start Time"];

  const renderRow = (payload: any) => (
    <>
      <TableCell>{payload.name}</TableCell>
      <TableCell>{payload.payloadArch}</TableCell>
      <TableCell>{getPayloadListenerType(payload.payloadListener)}</TableCell>
      <TableCell>{getPayloadListenerName(payload.payloadListener)}</TableCell>
      <TableCell>{format(payload.startTime, "Pp")}</TableCell>
      <TableCell className="text-right flex flex-row items-center justify-center gap-4">

        <Trash2 className="cursor-pointer" size={18} strokeWidth={1} onClick={() => deletePayload(payload)} />
        <Download className="cursor-pointer" size={18} strokeWidth={1} onClick={() => downloadPayload(payload)} />
      </TableCell>
    </>
  )

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


  const fetchPayloads = async () => {
    try {
        const url = "/payloads";
        const payloadsData = await axiosInstance.get(url,
          {
            headers: { authorization: `Bearer ${cookies.jwt}` }
          }
        );
        console.log(payloadsData.data);

        if(payloadsData.data.data){
          setPayloads((prevPayloads) => [
            ...prevPayloads,
            ...payloadsData.data.data.map((payload : any) => (
              {
              name: payload.payload_name || "Unknown Listener",
              payloadArch: payload.architecture || "Unknown",
              payloadListener: payload.listener_id,
              startTime: payload.created_at || "N/A",
              id: payload.payload_id
            })),
          ]);
        }
        
    } catch (error) {
        console.log(error);
    }
  };

  const getPayloadListenerType = (listenerId: string) => {
    let foundType = "";
    listeners.map((listener) => {
      if(listener.id === listenerId){
        foundType = listener.listenerType;
        return `${listener.listenerType}`;
      }
    })
    if(foundType){
      return foundType;
    } else{
      return "N/A";
    }
  }

  const getPayloadListenerName = (listenerId: string) => {
    let foundName = "";
    listeners.map((listener) => {
      if(listener.id === listenerId){
        foundName = listener.name;
      }
    })
    if(foundName){
      return foundName;
    } else{
      return "N/A";
    }
  }

  const createPayload = async () => {
    try{
        const url = '/payloads';
        const { data } = await axiosInstance.post(
            url,
            {
              "payload_name": payloadName,
              "listener_id": selectedListener,
              "architecture": payloadArchitecture,
              "self_contained": selfContained,
              // "publish_single_file": singleFile,
              "sleep": sleepValue,
              "jitter": jitterValue,
              // "start_date": "",
              // "kill_date": "",
              // "working_hours_start": 9,
              // "working_hours_end": 17
            },            
            {
                headers: { authorization: `Bearer ${cookies.jwt}` }
            }
        );
        toast({
          variant: "success",
          title: "Payload Creation Successful!",
          description: "You have successfully created a payload.",
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
          title: "Payload Creation Failed...",
          description: "Failed to generate payload. Please try again later.",
          action: (
            <ToastAction altText="undo">Close</ToastAction>
          ),
        })
    }
  }

  const downloadPayload = async (payload: any) => {
    try{
        console.log(payload)
        const url = `/payloads/download/${payload.id}`;
        const data : any = await axiosInstance.get(
            url,
            {
                headers: { authorization: `Bearer ${cookies.jwt}` }
            }
        ).then((data) => {
          const file = new Blob([data.data], {
            type: "application/octet-stream",
          });
          // saveAs(file, `${data.headers["Content-Disposition"]}`);
          saveAs(file, "payload.exe")
        });
        toast({
          variant: "success",
          title: "Your payload has been downloaded.",
          description: "A popup window should appear to download your payload.",
          action: (
            <ToastAction altText="undo">Close</ToastAction>
          ),
        })
    
    }
    catch(error){
        console.log(error);
        toast({
          variant: "destructive",
          title: "Unable to Download Payload...",
          description: "Please try again later.",
          action: (
            <ToastAction altText="undo">Close</ToastAction>
          ),
        })
    }
  }

  const deletePayload = async (payload: any) => {
    try{
      console.log(payload)
        const url = `/payloads/${payload.id}`;
        const { data } = await axiosInstance.delete(
            url,
            {
                headers: { authorization: `Bearer ${cookies.jwt}` }
            }
        );
        toast({
          variant: "success",
          title: "Your payload has been deleted.",
          description: "You payload is now deleted.",
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
          title: "Unable to Delete Payload...",
          description: "Please try again later.",
          action: (
            <ToastAction altText="undo">Close</ToastAction>
          ),
        })
    }
  }

  React.useEffect(() => {
    fetchListeners();
    fetchPayloads();
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
              <TableComponent headers={agentHeaders} data={payloads} renderRow={renderRow} />
            </CardContent>
          </Card>
          <div className="flex flex-row w-[calc(100vw-var(--sidebar-width)-6.5em)] justify-between">
            {/* <p>1 total listener(s) found.</p> */}
            <p>{payloads.length} total payload(s) found.</p>
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
          <Button className="w-[10em]" onClick={() => createPayload()}>+ Create</Button>
        </div>
        <div className="w-[calc(100vw-var(--sidebar-width)-6.5em)] grid grid-cols-3 gap-0 items-start justify-items-end p-0 mb-0 mt-0 border-solid border-[1px] pt-5 rounded-lg">
          <Card className="mx-auto border-none w-[100%] space-x-2">
            <CardContent>
              <div className="space-y-4">
                <div className="space-y-2">
                    <Label htmlFor="email">Name</Label>
                    <Input id="email" type="email" placeholder="Payload Name" value={payloadName} onChange={(e) => setPayloadName(e.target.value)} required />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="email">Select Listener</Label>
                  <Select onValueChange={setSelectedListener}>
                    <SelectTrigger>
                      <SelectValue placeholder="Select Listener By Name"/>
                    </SelectTrigger>
                    <SelectContent>
                      <SelectGroup>
                        {/* <SelectItem value="win-x64">win-x64</SelectItem>
                        <SelectItem value="win-x86">win-x86</SelectItem>
                        <SelectItem value="linux-x64">linux-x64</SelectItem>
                        <SelectItem value="linux-x86">linux-x86</SelectItem> */}
                        {listeners.map((listener) => {
                          return <SelectItem value={listener.id}>{listener.name}</SelectItem>
                        })

                        }
                      </SelectGroup>
                    </SelectContent>
                  </Select>
                </div>
                <div className="space-y-2">
                  <Label htmlFor="email">Payload Architecture</Label>
                  <Select onValueChange={setPayloadArchitecture}>
                    <SelectTrigger>
                      <SelectValue placeholder="Select Payload Architecture"/>
                    </SelectTrigger>
                    <SelectContent>
                      <SelectGroup>
                        <SelectItem value="win-x64">win-x64</SelectItem>
                        <SelectItem value="win-x86">win-x86</SelectItem>
                        <SelectItem value="linux-x64">linux-x64</SelectItem>
                        <SelectItem value="linux-x86">linux-x86</SelectItem>
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
                <div className="space-y-2">
                  <Label htmlFor="email">Sleep</Label>
                  <Input id="email" type="email" placeholder="5" value={sleepValue} onChange={(e) => setSleepValue(parseInt(e.target.value))} required />
                </div>

                <div className="space-y-2">
                  <Label htmlFor="email">Jitter</Label>
                  <Input id="email" type="email" placeholder="5" value={jitterValue} onChange={(e) => setJitterValue(parseInt(e.target.value))} required />
                </div>

                <div className="space-y-2">
                  <Label htmlFor="email">Generation Options</Label>
                </div>

                <div className="flex items-center space-x-2">
                  <Switch id="proxy-mode" checked={selfContained} onCheckedChange={setSelfContained} />
                  <Label htmlFor="enable-contained">Generate Self Contained</Label>
                </div>
                <div className="flex items-center space-x-2">
                  <Switch id="proxy-mode" checked={singleFile} onCheckedChange={setSingleFile} />
                  <Label htmlFor="single-file">Generate Single File</Label>
                </div>
              </div>
            </CardContent>
          </Card>

          <Card className="mx-auto border-none w-[100%] space-x-2">
            <CardContent className="flex flex-col gap-5">

              <div className="space-y-2" >
                <Label htmlFor="email">Start / Kill Date</Label>
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
            </CardContent>
          </Card>
        </div>
      </div>

    );
  }
}