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
  return (
    <div className="w-[calc(100vw-var(--sidebar-width))] h-[100%] flex flex-col gap-4 justify-items-center min-h-screen pb-4 p-0 m-6 font-[family-name:var(--font-geist-sans)]">
      Hello!
    </div>
  );
}