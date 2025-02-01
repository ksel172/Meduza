"use client"

import * as React from "react"
import { Card, CardHeader, CardDescription, CardContent, CardTitle } from "@/components/ui/card"
import { Ear } from "lucide-react"

export function Infocard({
    title,
    value,
    icon,
  }: {
    title: string
    value: string
    icon: React.ReactNode
  }) {

  return (
    <Card className="w-[100%]">
        <CardHeader className="flex flex-row justify-between items-center">
            <CardDescription>{title}</CardDescription>
            {/* <Ear size={20}/> */}
            {icon}
        </CardHeader>
        <CardContent>
            <CardTitle>{value}</CardTitle>
        </CardContent>
    </Card>
  )
}
