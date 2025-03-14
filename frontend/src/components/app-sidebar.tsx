import * as React from "react"
import {
  AudioWaveform,
  BookOpen,
  Bot,
  Command,
  ComputerIcon,
  Frame,
  GalleryVerticalEnd,
  LayoutDashboard,
  Map,
  PieChart,
  Radar,
  Settings2,
  ShieldAlert,
  SquareTerminal,
  Target,
  Warehouse,
} from "lucide-react"

import { NavMain } from "@/components/nav-main"
import { NavUser } from "@/components/nav-user"
import { ProjectSwitcher } from "@/components/project-switcher"
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarHeader,
  SidebarRail,
} from "@/components/ui/sidebar"

import { useStore } from "@nanostores/react"
import { user } from "@/lib/userStore"

import { ThemeModeToggle } from "./theme-mode-toggle"

// This is sample data.
const data = {
  user: {
    name: "shadcn",
    email: "m@example.com",
    avatar: "/avatars/shadcn.jpg",
  },
  navMain: [
    {
      title: "Scans",
      url: "/project/scans/overview",
      icon: Radar,
      isActive: true,
      items: [
        {
          title: "Overview",
          url: "/project/scans/overview",
        },
      ],
    },
    {
      title: "Assets",
      url: "/project/targets/overview",
      icon: Target,
      isActive: true,
      items: [
        {
          title: "Targets",
          url: "/project/targets/overview",
        },
        {
          title: "Services",
          url: "/project/targets/services",
        },
        {
          title: "Applications",
          url: "/project/targets/applications",
        },
      ],
    },
    {
      title: "Findings",
      url: "#",
      icon: ShieldAlert,
      isActive: true,
      items: [
        {
          title: "Overview",
          url: "/project/findings/overview"
        },
        {
          title: "Manage",
          url: "/project/findings/manage",
        },
      ],
    },
    {
      title: "Report",
      url: "#",
      icon: BookOpen,
      isActive: true,
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
      title: "Inventory",
      url: "#",
      icon: Warehouse,
      isActive: false,
      items: [
        {
          title: "Overview",
          url: "#",
        },
        {
          title: "Ports",
          url: "#",
        },
        {
          title: "Certificates",
          url: "#",
        },
        {
          title: "Services",
          url: "#",
        },
        {
          title: "Applications",
          url: "#",
        },
        {
          title: "Assets",
          url: "#",
        },
      ],
    },
    {
      title: "Settings",
      url: "/project/settings/general",
      icon: Settings2,
      isActive: false,
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
  const $activeUser = useStore(user);

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
        <NavUser user={$activeUser} />
      </SidebarFooter>
      <SidebarRail />
    </Sidebar>
  )
}
