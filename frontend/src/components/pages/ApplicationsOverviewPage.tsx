import {useState, useEffect} from "react"

import {user} from "@/lib/userStore";
import { activeProjectIdStore } from "@/lib/projectsStore";
import { useStore } from "@nanostores/react";

import { getProjectApplications, getProjectTargets } from "@/lib/api/projects";
import { type Target } from "@/lib/api/targets";
import { type Application, deleteApplicationById } from "@/lib/api/applications";

import { toast } from "sonner";

import { getColumns } from "@/components/applications/data-table/columns";
import { DataTable } from "@/components/findings/data-table/data-table";

import CreateApplicationDialog from "@/components/applications/create-application-dialog";


export default function ApplicationsOverviewPage() {
    const $activeProjectId = useStore(activeProjectIdStore);
    const [applications, setApplications] = useState<Application[]>([]);
    const [targetsMap, setTargetsMap] = useState<Record<string, Target>>({});
    const $user = useStore(user);

    useEffect(() => {
        if($activeProjectId && $user?.access_token) {
            getProjectApplications($activeProjectId, $user.access_token).then(applications => {
                if ("error" in applications) {
                    return
                }
                console.log(applications[0]);
                setApplications(applications);
            })

            getProjectTargets($activeProjectId, $user.access_token).then(result => {
                if ("error" in result) {
                    return
                }
                const targetsMapTmp: Record<string, Target> = {}
                result.forEach((target: Target) => {
                    if (!(target.id in targetsMapTmp)) {
                        targetsMapTmp[target.id] = {...target}
                    }
                });
                setTargetsMap(targetsMapTmp);
            });
        }
    }, [$activeProjectId, $user?.access_token])

    const handleBulkDelete = async (selectedRows: any) => {
        const removedApplicationIds: string[] = []
        for(let row of selectedRows) {
            if (row.id && $user?.access_token != undefined) {
                const resp = await deleteApplicationById(row.id, $user?.access_token)
                if ("error" in resp) {
                    toast("Something went wrong")
                }
                removedApplicationIds.push(row.id)
            }
        }
        const newApplications = applications.slice().filter(t => !removedApplicationIds.includes(t.id ?? ""))
        setApplications(newApplications);
        toast(`Successfully removed ${selectedRows.length} applications.`)
    }

    return (
        <div className="mt-4 container mx-auto ">
            <CreateApplicationDialog setApplications={setApplications} />
            <div className="py-10">
                <DataTable
                    columns={getColumns(setApplications, applications)}
                    data={applications.map(s => ({...s, target: targetsMap[s.host_target]?.value ?? ""}))}
                    filterSettings={{placeholder: "Filter by description...", filterKey: "description"}}
                    handleDelete={handleBulkDelete}
                />
            </div>
        </div>
    );
}