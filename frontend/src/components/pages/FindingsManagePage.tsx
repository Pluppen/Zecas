import {useState, useEffect} from "react"

import {user} from "@/lib/userStore";
import { activeProjectIdStore } from "@/lib/projectsStore";
import { useStore } from "@nanostores/react";

import { getProjectFindings, getProjectTargets } from "@/lib/api/projects";

import CreateFindingDialog from "@/components/findings/create-finding-dialog";

import { getColumns } from "../findings/data-table/columns";
import { DataTable } from "../findings/data-table/data-table";


export default function FindingsManagePage() {
    const $activeProjectId = useStore(activeProjectIdStore);
    const [findings, setFindings] = useState([]);
    const [targetsMap, setTargetsMap] = useState({});
    const $user = useStore(user);

    useEffect(() => {
        if($activeProjectId) {
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
                const targetsMapTmp = {}
                result.forEach(target => {
                    if (!(target.id in targetsMapTmp)) {
                        targetsMapTmp[target.id] = {...target}
                    }
                });
                setTargetsMap(targetsMapTmp);
            });
        }
    }, [$activeProjectId, $user])

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