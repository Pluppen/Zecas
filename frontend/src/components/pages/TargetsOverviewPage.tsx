import {useState, useEffect} from "react"

import {user} from "@/lib/userStore";
import { activeProjectIdStore } from "@/lib/projectsStore";
import { useStore } from "@nanostores/react";

import { getProjectTargets } from "@/lib/api/projects";
import { type Target, deleteTargetById } from "@/lib/api/targets";
import { toast } from "sonner"

import { getColumns } from "@/components/targets/data-table/columns";
import { DataTable } from "@/components/findings/data-table/data-table";

import CreateTargetDialog from "@/components/targets/create-target-dialog";


export default function FindingsManagePage() {
    const $activeProjectId = useStore(activeProjectIdStore);
    const [targets, setTargets] = useState<Target[]>([]);
    const $user = useStore(user);

    useEffect(() => {
        if($activeProjectId && $user?.access_token) {
            getProjectTargets($activeProjectId, $user.access_token).then(async targets => {
                if ("error" in targets) {
                    return
                }
                setTargets(targets);
            });
        }
    }, [$activeProjectId, $user?.access_token])

    const handleBulkDelete = async (selectedRows: any) => {
        const removedTargetIds: string[] = []
        for(let row of selectedRows) {
            if (row.id && $user?.access_token != undefined) {
                const resp = await deleteTargetById(row.id, $user?.access_token)
                if ("error" in resp) {
                    toast("Something went wrong")
                }
                removedTargetIds.push(row.id)
            }
        }
        const newTargets = targets.slice().filter(t => !removedTargetIds.includes(t.id))
        setTargets(newTargets);
        toast(`Successfully removed ${selectedRows.length} targets.`)
    }

    return (
        <div className="mt-4 container mx-auto ">
            <CreateTargetDialog setTargets={setTargets} />
            <div className="py-10">
                <DataTable handleDelete={handleBulkDelete} columns={getColumns(setTargets, targets)} data={targets} filterSettings={{placeholder: "Filter by target", filterKey: "value"}} />
            </div>
        </div>
    );
}