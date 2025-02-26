import React from "react"
import { AppSidebar } from "@/components/app-sidebar"
import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator,
} from "@/components/ui/breadcrumb"
import { Separator } from "@/components/ui/separator"
import {
  SidebarInset,
  SidebarProvider,
  SidebarTrigger,
} from "@/components/ui/sidebar"

import { activeProjectStore, activeProjectIdStore, projects, type Project } from "@/lib/projectsStore"
import { useStore } from "@nanostores/react"

interface Props {
  children: React.ReactNode
  defaultSidebarOpen: boolean
  breadcrumbL1: string
  breadcrumbL2: string
  projectData: {
    activeProject: Project
    projects: Project[]
  }
}

export default function Page(props: Props) {
  const $activeProjectId = useStore(activeProjectIdStore);

  React.useEffect(() => {
    let activeProject = props.projectData.projects[0]
    if ($activeProjectId !== "null" && $activeProjectId) {
      activeProject = props.projectData.projects.filter(project => project.id == $activeProjectId)[0];
    }

    if (props.projectData.activeProject) {
      activeProject = props.projectData.activeProject;
    }

    activeProjectStore.set(activeProject);
    projects.set({projects: props.projectData.projects});
  }, [$activeProjectId])

  return (
    <SidebarProvider>
      <AppSidebar />
      <SidebarInset>
        <header className="flex h-16 shrink-0 items-center gap-2 transition-[width,height] ease-linear group-has-data-[collapsible=icon]/sidebar-wrapper:h-12">
          <div className="flex items-center gap-2 px-4">
            <SidebarTrigger className="-ml-1" />
            <Separator
              orientation="vertical"
              className="mr-2 data-[orientation=vertical]:h-4"
            />
            <Breadcrumb>
              <BreadcrumbList>
                <BreadcrumbItem className="hidden md:block">
                  <BreadcrumbLink href="#">
                    {props.breadcrumbL1}
                  </BreadcrumbLink>
                </BreadcrumbItem>
                <BreadcrumbSeparator className="hidden md:block" />
                <BreadcrumbItem>
                  <BreadcrumbPage>{props.breadcrumbL2}</BreadcrumbPage>
                </BreadcrumbItem>
              </BreadcrumbList>
            </Breadcrumb>
          </div>
        </header>
        <div className="flex flex-1 flex-col gap-4 p-4 pt-0">
          {props.children}
        </div>
      </SidebarInset>
    </SidebarProvider>
  )
}
