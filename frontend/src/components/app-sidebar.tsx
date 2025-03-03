import * as React from "react"
import {
  AudioWaveform,
  BookOpen,
  Bot,
  Command,
  Frame,
  GalleryVerticalEnd,
  Map,
  PieChart,
  Settings2,
  SquareTerminal,
} from "lucide-react"

import { NavMain } from "@/components/nav-main"
import { NavProjects } from "@/components/nav-projects"
import { NavUser } from "@/components/nav-user"
import { ProjectSwitcher } from "@/components/project-switcher"
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarHeader,
  SidebarRail,
} from "@/components/ui/sidebar"

import { ThemeModeToggle } from "./theme-mode-toggle"

// This is sample data.
const data = {
  user: {
    name: "shadcn",
    email: "m@example.com",
    avatar: "/avatars/shadcn.jpg",
  },
  teams: [
    {
      name: "Pentest",
      logo: GalleryVerticalEnd,
      plan: "Enterprise",
    },
    {
      name: "Attack Surface Assessment",
      logo: AudioWaveform,
      plan: "Startup",
    },
    {
      name: "Iver Sverige",
      logo: Command,
      plan: "Free",
    },
  ],
  navMain: [
    {
      title: "Scans",
      url: "/project/scans/overview",
      icon: SquareTerminal,
      isActive: true,
      items: [
        {
          title: "Overview",
          url: "/project/scans/overview",
        },
        {
          title: "Results #TODO",
          url: "#",
        },
        {
          title: "Running #TODO",
          url: "#",
        },
        {
          title: "Settings #TODO",
          url: "#",
        },
      ],
    },
    {
      title: "Findings",
      url: "#",
      icon: Bot,
      items: [
        {
          title: "Overview # TODO",
          url: "/project/findings/overview"
        },
        {
          title: "Manage",
          url: "/project/findings/manage",
        },
        {
          title: "Vulnerabilites",
          url: "#",
        },
      ],
    },
    {
      title: "Report",
      url: "#",
      icon: BookOpen,
      items: [
        {
          title: "Edit Report",
          url: "#",
        },
        {
          title: "Generated Reports",
          url: "#",
        },
      ],
    },
    {
      title: "Settings",
      url: "/project/settings/general",
      icon: Settings2,
      items: [
        {
          title: "General",
          url: "/project/settings/general",
        },
        {
          title: "Targets",
          url: "/project/settings/targets",
        },
        {
          title: "Scanners",
          url: "#",
        },
      ],
    },
  ],
  projects: [
    {
      name: "LSR",
      url: "#",
      icon: Frame,
    },
    {
      name: "Mornington Hotel",
      url: "#",
      icon: PieChart,
    },
  ],
}

export function AppSidebar({ ...props }: React.ComponentProps<typeof Sidebar>) {
  return (
    <Sidebar collapsible="icon" {...props}>
      <SidebarHeader>
        <ProjectSwitcher />
      </SidebarHeader>
      <SidebarContent>
        <NavMain items={data.navMain} />
      </SidebarContent>
      <SidebarFooter>
        <ThemeModeToggle />
        <NavUser user={data.user} />
      </SidebarFooter>
      <SidebarRail />
    </Sidebar>
  )
}
