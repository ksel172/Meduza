/**
 * v0 by Vercel.
 * @see https://v0.dev/t/kpODUatqMES
 * Documentation: https://v0.dev/docs#integrating-generated-code-into-your-nextjs-app
 */
"use client"

import { useState } from "react"
import { Input } from "@/components/ui/input"
import { Avatar, AvatarImage, AvatarFallback } from "@/components/ui/avatar"
import { Button } from "@/components/ui/button"

export default function ChatRoom() {
  const [searchTerm, setSearchTerm] = useState("")

  return (
    <div className="flex flex-col h-[100%]">
    <div className="flex-1 flex flex-col">
        <div className="flex-col overflow-y-auto p-1 grid gap-4">
            <div className="flex items-start gap-4">
                <Avatar className="w-8 h-8 border">
                    <AvatarImage src="/placeholder-user.jpg" alt="John Doe" />
                    <AvatarFallback>JD</AvatarFallback>
                </Avatar>
                <div className="bg-muted p-3 rounded-lg max-w-[80%]">
                    <p>Hey, did you see the latest update?</p>
                    <p className="text-xs text-muted-foreground mt-1">3:45 PM</p>
                </div>
            </div>
            <div className="flex items-start gap-4">
                <Avatar className="w-8 h-8 border">
                    <AvatarImage src="/placeholder-user.jpg" alt="John Doe" />
                    <AvatarFallback>JD</AvatarFallback>
                </Avatar>
                <div className="bg-muted p-3 rounded-lg max-w-[80%]">
                    <p>Hey, did you see the latest update?</p>
                    <p className="text-xs text-muted-foreground mt-1">3:45 PM</p>
                </div>
            </div>
            <div className="flex items-start gap-4 justify-end">
                <div className="bg-green-900 p-3 rounded-lg max-w-[80%]">
                    <p>Hey, did you see the latest update?</p>
                    <p className="text-xs text-muted-foreground mt-1">3:45 PM</p>
                </div>
                <Avatar className="w-8 h-8 border">
                    <AvatarImage src="/placeholder-user.jpg" alt="Jane Doe" />
                    <AvatarFallback>JD</AvatarFallback>
                </Avatar>
            </div>
            <div className="flex items-start gap-4 justify-end">
                <div className="bg-green-900 p-3 rounded-lg max-w-[80%]">
                    <p>Hey, did you see the latest update?</p>
                    <p className="text-xs text-muted-foreground mt-1">3:45 PM</p>
                </div>
                <Avatar className="w-8 h-8 border">
                    <AvatarImage src="/placeholder-user.jpg" alt="Jane Doe" />
                    <AvatarFallback>JD</AvatarFallback>
                </Avatar>
            </div>
            <div className="flex items-start gap-4">
                <Avatar className="w-8 h-8 border">
                    <AvatarImage src="/placeholder-user.jpg" alt="John Doe" />
                    <AvatarFallback>JD</AvatarFallback>
                </Avatar>
                <div className="bg-muted p-3 rounded-lg max-w-[80%]">
                    <p>Hey, did you see the latest update?</p>
                    <p className="text-xs text-muted-foreground mt-1">3:45 PM</p>
                </div>
            </div>
        </div>
        <div className="sticky bottom-0 bg-background p-4 flex items-center gap-2">
            <Input
                type="text"
                placeholder="Type your message..."
                className="flex-1 bg-background rounded-full px-4 py-2 focus:outline-none focus:ring-2 focus:ring-primary focus:ring-offset-2"
            />
            <Button variant="ghost" size="icon" className="rounded-full">
                <SendIcon className="w-5 h-5" />
            </Button>
        </div>
    </div>
    </div>
  )
}


function SendIcon(props) {
  return (
    <svg
      {...props}
      xmlns="http://www.w3.org/2000/svg"
      width="24"
      height="24"
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      strokeWidth="2"
      strokeLinecap="round"
      strokeLinejoin="round"
    >
      <path d="m22 2-7 20-4-9-9-4Z" />
      <path d="M22 2 11 13" />
    </svg>
  )
}