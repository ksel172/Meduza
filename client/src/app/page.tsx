import Image from "next/image";
import { Navbar } from "@/components/util/navbar/navbar";
import { TableComponent } from "@/components/util/navbar/table";
import * as React from "react"
import ChatRoom from "@/components/util/navbar/chat";
import ConsoleWidget from "@/components/util/navbar/console";
import axios from 'axios';

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

import { Ear, Radio, Skull} from "lucide-react";
import { Infocard } from "@/components/util/items/infocard";

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

  const axiosInstance = axios.create({
    baseURL: "http://localhost:2000",
  })

  const fetchRequest = async () => {
    try {
        const url = "/v1/chat";
        const blogsData = await axiosInstance.get(url);
    } catch (error) {
        console.log(error);
    }
  };

  const postRequest = async (param: string) => {
    try{
        const url = '/v1/post/endpoint';
        const { data } = await axiosInstance.post(
            url,
            // {
            //   param: param,
            // },
            // {
            //     headers: { authorization: `Bearer ${cookies.jwt}` }
            // }
        );
    }
    catch(error){
        console.log(error);
    }
  }

  return (
    <div className="w-[calc(100vw-var(--sidebar-width))] h-[100%] flex flex-col gap-4 justify-items-center min-h-screen pb-4 p-0 m-6 overflow-x-hidden">
      <div className="flex flex-row gap-[3.2vw]">
        <Infocard title={"Active Listeners"} value={"1"} icon={<Ear size={20}/>}/>
        <Infocard title={"Live Agents"} value={"2"} icon={<Radio size={20}/>}/>
        <Infocard title={"Dead Agents"} value={"0"} icon={<Skull size={20}/>}/>
        <Infocard title={"Dead Agents"} value={"0"} icon={<Skull size={20}/>}/>
      </div>
      <div className="flex flex-row gap-4 w-[calc(95vw-var(--sidebar-width))]">
        <Card className="w-[60%]">
          <CardContent className="m-0 p-0 h-[70vh]">
            <ConsoleWidget />
          </CardContent>
        </Card>
        <Card className="w-[40%] h-[70vh] overflow-y-scroll">
          <CardHeader className="sticky top-0 bg-black z-10">
            <CardTitle>Teamchat</CardTitle>
            <CardDescription>Communicate with your team using teamchat.</CardDescription>
          </CardHeader>
          <CardContent className="m-0 p-0 h-[70vh]">
            <ChatRoom />
          </CardContent>
        </Card>
      </div>
    </div>
  );
}