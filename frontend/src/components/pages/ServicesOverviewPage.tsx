import {useState, useEffect} from "react"

import {user} from "@/lib/userStore";
import { activeProjectIdStore } from "@/lib/projectsStore";
import { useStore } from "@nanostores/react";

import { getProjectServices, getProjectTargets } from "@/lib/api/projects";
import { type Target } from "@/lib/api/targets";
import { type Service, deleteServiceById } from "@/lib/api/services";

import { getColumns } from "@/components/services/data-table/columns";
import { DataTable } from "@/components/findings/data-table/data-table";
import CreateServiceDialog from "../services/create-service-dialog";

import { toast } from "sonner";


export default function ServicesOverviewPage() {
    const $activeProjectId = useStore(activeProjectIdStore);
    const [services, setServices] = useState<Service[]>([]);
    const [targetsMap, setTargetsMap] = useState<Record<string, Target>>({});
    const $user = useStore(user);

    useEffect(() => {
        if($activeProjectId && $user?.access_token) {
            getProjectServices($activeProjectId, $user.access_token).then(services => {
                if ("error" in services) {
                    return
                }
                setServices(services);
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
    }, [$activeProjectId, $user])

    const handleBulkDelete = async (selectedRows: any) => {
        const removedServiceIds: string[] = []
        for(let row of selectedRows) {
            if (row.id && $user?.access_token != undefined) {
                const resp = await deleteServiceById(row.id, $user?.access_token)
                if ("error" in resp) {
                    toast("Something went wrong")
                }
                removedServiceIds.push(row.id)
            }
        }
        const newServices = services.slice().filter(t => !removedServiceIds.includes(t.id ?? ""))
        setServices(newServices);
        toast(`Successfully removed ${selectedRows.length} services.`)
    }

    return (
        <div className="mt-4 container mx-auto ">
            <CreateServiceDialog setServices={setServices} />
            <div className="py-10">
                <DataTable
                    columns={getColumns(setServices, services)}
                    data={services.map(s => ({...s, target: targetsMap[s.target_id]?.value ?? ""}))}
                    filterSettings={{placeholder: "Filter by description...", filterKey: "description"}}
                    handleDelete={handleBulkDelete}
                />
            </div>
        </div>
    );
}