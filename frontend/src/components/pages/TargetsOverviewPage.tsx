import {useState, useEffect} from "react"

import {user} from "@/lib/userStore";
import { activeProjectIdStore } from "@/lib/projectsStore";
import { useStore } from "@nanostores/react";

import { getProjectTargets } from "@/lib/api/projects";

import { getColumns } from "@/components/targets/data-table/columns";
import { DataTable } from "@/components/findings/data-table/data-table";

import CreateTargetDialog from "@/components/targets/create-target-dialog";


export default function FindingsManagePage() {
    const $activeProjectId = useStore(activeProjectIdStore);
    const [targets, setTargets] = useState({});
    const $user = useStore(user);

    useEffect(() => {
        if($activeProjectId) {
            getProjectTargets($activeProjectId, $user.access_token).then(async targets => {
                if ("error" in targets) {
                    return
                }
                setTargets(targets);
            });
        }
    }, [$activeProjectId, $user])

    return (
        <div className="mt-4 container mx-auto ">
            <CreateTargetDialog setTargets={setTargets} />
            <div className="py-10">
                <DataTable columns={getColumns(setTargets, targets)} data={targets} filterSettings={{placeholder: "Filter by target", filterKey: "value"}} />
            </div>
        </div>
    );
}