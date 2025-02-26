import * as React from "react"
import { ChevronsUpDown, Plus } from "lucide-react"
import { useContext } from 'react';
import { ProjectContext } from "@/lib/ProjectContext";
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

import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuShortcut,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import {
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  useSidebar,
} from "@/components/ui/sidebar"
import { activeProjectStore, projects, activeProjectIdStore, type Project } from "@/lib/projectsStore";
import { useStore } from '@nanostores/react';

export function ProjectSwitcher() {
  const { isMobile } = useSidebar()

  const $projectData = useStore(projects);
  const $activeProject = useStore(activeProjectStore);

  const setActiveProjectId = (projectId: string) => {
    const activeProject = $projectData.projects.filter(project => project.id === projectId)[0];
    if (activeProject.id !== undefined) {
      activeProjectIdStore.set(activeProject.id);
      window.location.href = "/";
    }
  }

  return (
    <SidebarMenu>
      <SidebarMenuItem>
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <SidebarMenuButton
              size="lg"
              className="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"
            >
              <div className="bg-sidebar-primary text-sidebar-primary-foreground flex aspect-square size-8 items-center justify-center rounded-lg">
                <GalleryVerticalEnd className="size-4" />
              </div>
              <div className="grid flex-1 text-left text-sm leading-tight">
                <span className="truncate font-medium">{$activeProject?.name}</span>
                <span className="truncate text-xs">{}</span>
              </div>
              <ChevronsUpDown className="ml-auto" />
            </SidebarMenuButton>
          </DropdownMenuTrigger>
          <DropdownMenuContent
            className="w-(--radix-dropdown-menu-trigger-width) min-w-56 rounded-lg"
            align="start"
            side={isMobile ? "bottom" : "right"}
            sideOffset={4}
          >
            <DropdownMenuLabel className="text-muted-foreground text-xs">
              Projects
            </DropdownMenuLabel>
            {$projectData.projects.map((project: any, index: number) => (
              <DropdownMenuItem
                key={"project-id-"+project.id}
                onClick={() => setActiveProjectId(project.id)}
                className="gap-2 p-2"
              >
                <div className="flex size-6 items-center justify-center rounded-xs border">
                  <GalleryVerticalEnd className="size-4 shrink-0" />
                </div>
                {project.name}
                <DropdownMenuShortcut>⌘{index + 1}</DropdownMenuShortcut>
              </DropdownMenuItem>
            ))}
            <DropdownMenuSeparator />
            <DropdownMenuItem>
                <a href="/project/new" className="hover:cursor-pointer gap-2 p-2 focus:bg-accent focus:text-accent-foreground data-[variant=destructive]:text-destructive-foreground data-[variant=destructive]:focus:bg-destructive/10 dark:data-[variant=destructive]:focus:bg-destructive/40 data-[variant=destructive]:focus:text-destructive-foreground data-[variant=destructive]:*:[svg]:!text-destructive-foreground [&_svg:not([class*='text-'])]:text-muted-foreground relative flex cursor-default items-center gap-2 rounded-sm px-2 py-1.5 text-sm outline-hidden select-none data-[disabled]:pointer-events-none data-[disabled]:opacity-50 data-[inset]:pl-8 [&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*='size-'])]:size-4">
                  <div className="bg-background flex size-6 items-center justify-center rounded-md border">
                    <Plus className="size-4" />
                  </div>
                  <div className="text-muted-foreground font-medium">Add project</div>
                </a>
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </SidebarMenuItem>
    </SidebarMenu>
  )
}
