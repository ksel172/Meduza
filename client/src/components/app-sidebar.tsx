"use client"

import * as React from "react"
import {
  AudioWaveform,
  BookOpen,
  Bot,
  Circle,
  Command,
  EllipsisIcon,
  Frame,
  GalleryVerticalEnd,
  Map,
  PieChart,
  Settings2,
  SquareTerminal,
  LayoutDashboard,
  MessageSquareText,
  BookOpenText
} from "lucide-react"

import { NavMain } from "@/components/nav-main"
import { NavProjects } from "./nav-projects"
import { NavUser } from "@/components/nav-user"
import { TeamSwitcher } from "@/components/team-switcher"
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarHeader,
  SidebarRail,
} from "@/components/ui/sidebar"

// This is sample data.
const data = {
  user: {
    name: "shadcn",
    email: "m@example.com",
    avatar: "/avatars/shadcn.jpg",
  },
  teams: [
    {
      name: "Cyberaxis",
      logo: GalleryVerticalEnd
    }
  ],
  navMain: [
    {
      title: "Dashboard",
      url: "#",
      icon: LayoutDashboard,
      isActive: true,
    },
    {
      title: "Deployment",
      url: "#",
      icon: SquareTerminal,
      items: [
        {
          title: "Listeners",
          url: "#",
        },
        {
          title: "Agents",
          url: "#",
        },
        {
          title: "Payloads",
          url: "#",
        },
        {
            title: "Modules",
            url: "#",
        },
      ],
    },
    {
      title: "Options",
      url: "#",
      icon: Settings2,
    },
    {
      title: "Team Chat",
      url: "#",
      icon: MessageSquareText,
    },
    {
      title: "Documentation",
      url: "#",
      icon: BookOpenText,
    },
  ],
  projects: [
    {
        name: "10.10.14.12",
        url: "#",
        icon: Circle,
        status: "alive"
    },
    {
        name: "10.10.14.18",
        url: "#",
        icon: Circle,
        status: "dead"
    },
    {
        name: "10.10.14.11",
        url: "#",
        icon: Circle,
        status: "alive"
    }
  ],
}

export function AppSidebar({ ...props }: React.ComponentProps<typeof Sidebar>) {
  return (
    <Sidebar collapsible="icon" {...props}>
      <SidebarHeader>
        <TeamSwitcher teams={data.teams} />
      </SidebarHeader>
      <SidebarContent>
        <NavMain items={data.navMain} />
        <NavProjects projects={data.projects} />
      </SidebarContent>
      <SidebarFooter>
        <NavUser user={data.user} />
      </SidebarFooter>
      <SidebarRail />
    </Sidebar>
  )
}
