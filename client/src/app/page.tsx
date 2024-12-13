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


export default function Home() {
  return (
    <div className="items-center justify-items-center min-h-screen pb-4 p-0 mb-0 mt-10 font-[family-name:var(--font-geist-sans)]">
      <Navbar />
      <div className="grid grid-cols-2 gap-4 items-center justify-items-center min-h-screen p-0 mb-0 mt-10 font-[family-name:var(--font-geist-sans)]">
        <Card className="">
          <CardHeader>
            <CardTitle>Agent Connections</CardTitle>
            <CardDescription>Interact with your active agents.</CardDescription>
          </CardHeader>
          <CardContent className="w-[45vw] h-[20em] overflow-y-scroll">
            <form>
              <div className="grid w-full items-center gap-4">
                <div className="flex flex-col space-y-1.5">
                  <Label htmlFor="framework">Filter</Label>
                  <Select>
                    <SelectTrigger id="framework">
                      <SelectValue placeholder="Select" />
                    </SelectTrigger>
                    <SelectContent position="popper">
                      <SelectItem value="next">N/A</SelectItem>
                      <SelectItem value="sveltekit">Active</SelectItem>
                      <SelectItem value="astro">Dead</SelectItem>
                    </SelectContent>
                  </Select>
                </div>
              </div>
            </form>
            <TableComponent />
          </CardContent>
        </Card>

        <Card className="">
          <CardHeader>
            <CardTitle>Agent Connections</CardTitle>
            <CardDescription>Deploy your new project in one-click.</CardDescription>
          </CardHeader>
          <CardContent className="w-[45vw] h-[20em] overflow-y-scroll">
            <form>
              <div className="grid w-full items-center gap-4">
                <div className="flex flex-col space-y-1.5">
                  <Label htmlFor="framework">Filter</Label>
                  <Select>
                    <SelectTrigger id="framework">
                      <SelectValue placeholder="Select" />
                    </SelectTrigger>
                    <SelectContent position="popper">
                      <SelectItem value="next">N/A</SelectItem>
                      <SelectItem value="sveltekit">Active</SelectItem>
                      <SelectItem value="astro">Dead</SelectItem>
                    </SelectContent>
                  </Select>
                </div>
              </div>
            </form>
            <TableComponent />
          </CardContent>
        </Card>

        <Card className="">
          <CardHeader>
            <CardTitle>Teamchat</CardTitle>
            <CardDescription>Communicate with your team using teamchat.</CardDescription>
          </CardHeader>
          <CardContent className="w-[45vw] h-[20em] overflow-y-scroll pb-0 mb-0">
            <ChatRoom />
          </CardContent>
        </Card>

        <Card className="">
          <CardHeader>
            <CardTitle>Terminal Window</CardTitle>
            <CardDescription>Deploy your new project in one-click.</CardDescription>
          </CardHeader>
          <CardContent className="w-[45vw] h-[20em] overflow-y-scroll pb-0 mb-0">
            <ConsoleWidget />
          </CardContent>
        </Card>
      </div>
    </div>
  );
}