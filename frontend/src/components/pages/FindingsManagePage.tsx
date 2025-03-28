import {useState, useEffect} from "react"

import {user} from "@/lib/userStore";
import { activeProjectIdStore } from "@/lib/projectsStore";
import { useStore } from "@nanostores/react";

import { getProjectFindings, getProjectTargets } from "@/lib/api/projects";

import CreateFindingDialog from "@/components/findings/create-finding-dialog";

import { getColumns } from "../findings/data-table/columns";
import { DataTable } from "../findings/data-table/data-table";
import type { Target } from "@/lib/api/targets";
import { type Finding } from "@/lib/api/findings";


export default function FindingsManagePage() {
    const $activeProjectId = useStore(activeProjectIdStore);
    const [findings, setFindings] = useState<Finding[]>([]);
    const [targetsMap, setTargetsMap] = useState<Record<string, Target>>({});
    const $user = useStore(user);

    useEffect(() => {
        if($activeProjectId && $user?.access_token) {
            getProjectFindings($activeProjectId, $user.access_token).then(result => {
                if ("error" in result) {
                    return
                }
                setFindings(result);
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

    return (
        <div className="mt-8 container mx-auto ">
            <h1>Findings</h1>
            <CreateFindingDialog setFindings={setFindings} />
            <div className="py-10">
                <DataTable columns={getColumns(setFindings, findings)} data={findings.map(f => ({...f, target_value: targetsMap[f.target_id]?.value ?? ""}))} filterSettings={{placeholder: "Filter by title...", filterKey: "title"}} />
            </div>
        </div>
    );
}
