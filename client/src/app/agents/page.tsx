"use client"

import Image from "next/image";
import { Navbar } from "@/components/util/navbar/navbar";
import { TableComponent } from "@/components/util/navbar/table";
import * as React from "react"
import ChatRoom from "@/components/util/navbar/chat";
import ConsoleWidget from "@/components/util/navbar/console";
import axiosInstance from "@/axiosInstance";

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
import { useCookies } from "react-cookie";

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

  const [cookies, setCookie, removeCookie] = useCookies<any>(["cookie-name"]);

  // {
  //   "agent_id": "03c00218-044d-0523-bd06-980700080009",
  //   "name": "G9N9SG",
  //   "note": "",
  //   "status": 4,
  //   "first_callback": "2025-03-29T15:29:13.706401Z",
  //   "last_callback": "2025-03-29T15:31:40Z",
  //   "modified_at": "2025-03-29T15:29:13.706401Z"
  // }

  interface Agent {
    agent_id: string;
    name: string;
    note: string;
    status: string;
    // config_id: string;
    first_callback: string;
    last_callback: string;
    modified_at: string;
  }

  const [agents, setAgents] = React.useState<Agent[]>([]);

  const agentHeaders = ["Name", "Status", "Note", "Last Callback"];

  const getAgentStage = (value: number) => {
    // Uninitialized,
    // Stage0,
    // Stage1,
    // Stage2,
    // Active,
    // Lost,
    // Exited,
    // Disconnected,
    // Hidden
    if(value === 0){
      return "Uninitialized";
    } else if (value === 1){
      return "Stage 0";
    } else if (value === 2){
      return "Stage 1";
    } else if (value === 3){
      return "Stage 2";
    } else if (value === 4){
      return "Active";
    } else if (value === 5){
      return "Lost";
    } else if (value === 6){
      return "Exited";
    } else if (value === 7){
      return "Disconnected";
    } else{
      return "Hidden";
    }
  }

  const fetchAgents = async () => {
    try {
        const url = "/agents";
        const agentsData = await axiosInstance.get(url,
          {
            headers: { authorization: `Bearer ${cookies.access_token}` }
          }
        );
        console.log(agentsData.data.data);

        if(agentsData.data.data){
          setAgents((prevAgents) => [
            ...prevAgents,
            ...agentsData.data.data.map((agent : any) => (
              {
                agent_id: agent.agent_id,
                name: agent.name || "Unknown Agent",
                note: agent.note || "N/A",
                // status: (agent.status === 0 ? "Stopped" : (agent.status === 1 ? "Running" : (agent.status === 2 ? "Paused" : (agent.status === 3 ? "Processing" : "Error")))),
                status: getAgentStage(agent.status),
                first_callback: agent.first_callback || "N/A",
                last_callback: agent.last_callback || "N/A",
                modified_at: agent.modified_at || "N/A"
            })),
          ]);
        }
        
    } catch (error) {
        console.log(error);
    }
  };

  React.useEffect(() => {
    console.log("AGENTS UPDATED: ", agents)
  }, [agents])

  React.useEffect(() => {
      fetchAgents();
    }, [])

  const renderRow = (agent: any) => (
    <>
      <TableCell>{agent.name}</TableCell>
      <TableCell
        className={`font-medium ${
          agent.status === "Exited" || agent.status === "Lost" || agent.status === "Disconnected" ? "text-red-600" : ( agent.status === "Active"  ? "text-green-400" : "text-grey")
        }`}
      >
        {agent.status}
      </TableCell>
      <TableCell>{agent.note}</TableCell>
      <TableCell>{agent.last_callback}</TableCell>
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
           <TableComponent headers={agentHeaders} data={agents} renderRow={renderRow} />
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