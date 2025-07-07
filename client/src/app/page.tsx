import * as React from "react"
import ChatRoom from "@/components/util/navbar/chat";
import ConsoleWidget from "@/components/util/navbar/console";

import { Button } from "@/components/ui/button"
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"

import { Ear, Radio, Skull} from "lucide-react";
import { Infocard } from "@/components/util/items/infocard";

export default function Home() {
  return (
    <div className="w-[calc(100vw-var(--sidebar-width))] h-[100%] flex flex-col gap-4 justify-items-center min-h-screen pb-4 p-0 m-6 overflow-x-hidden animate-fade-in">
      {/* Enhanced Button Showcase */}
      <Card variant="interactive" animate className="mb-6">
        <CardHeader>
          <CardTitle className="animate-slide-down">UI Enhancement Demo</CardTitle>
          <CardDescription className="animate-slide-down" style={{animationDelay: '0.1s'}}>
            Showcasing impressive button animations and modern design elements
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div className="flex flex-wrap gap-4 mb-4">
            <Button variant="default" size="default">
              Default Button
            </Button>
            <Button variant="gradient" size="default" animation="glow">
              Gradient Button
            </Button>
            <Button variant="glassmorphism" size="default">
              Glass Button
            </Button>
            <Button variant="premium" size="lg" animation="float">
              Premium Button
            </Button>
            <Button variant="destructive" size="default" animation="heartbeat">
              Destructive
            </Button>
            <Button variant="outline" size="default" animation="bounce">
              Outline Bounce
            </Button>
          </div>
          <div className="flex gap-4">
            <Button variant="gradient" size="xl" ripple={true}>
              Large Ripple Effect
            </Button>
            <Button variant="premium" size="default" animation="wiggle">
              Wiggle Animation
            </Button>
          </div>
        </CardContent>
      </Card>

      <div className="flex flex-row gap-[3.2vw] animate-slide-up">
        <div className="animate-slide-right" style={{animationDelay: '0.1s'}}>
          <Infocard title={"Active Listeners"} value={"1"} icon={<Ear size={20}/>}/>
        </div>
        <div className="animate-slide-right" style={{animationDelay: '0.2s'}}>
          <Infocard title={"Live Agents"} value={"2"} icon={<Radio size={20}/>}/>
        </div>
        <div className="animate-slide-right" style={{animationDelay: '0.3s'}}>
          <Infocard title={"Dead Agents"} value={"0"} icon={<Skull size={20}/>}/>
        </div>
        <div className="animate-slide-right" style={{animationDelay: '0.4s'}}>
          <Infocard title={"Total Agents"} value={"2"} icon={<Skull size={20}/>}/>
        </div>
      </div>
      <div className="flex flex-row gap-4 w-[calc(95vw-var(--sidebar-width))] animate-scale-in" style={{animationDelay: '0.5s'}}>
        <Card variant="glow" className="w-[60%]">
          <CardContent className="m-0 p-0 h-[70vh]">
            <ConsoleWidget />
          </CardContent>
        </Card>
        <Card variant="glass" className="w-[40%] h-[70vh] overflow-y-scroll">
          <CardHeader className="sticky top-0 bg-black/80 backdrop-blur-sm z-10 border-b border-white/10">
            <CardTitle className="animate-glow">Teamchat</CardTitle>
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