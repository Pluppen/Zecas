import {useState, useEffect} from "react"

import { activeProjectIdStore } from "@/lib/projectsStore";
import { useStore } from "@nanostores/react";

import { getProjectFindings } from "@/lib/findings";
import { getProjectTargets } from "@/lib/targets";

import CreateFindingDialog from "@/components/findings/create-finding-dialog";

import { getColumns } from "../findings/data-table/columns";
import { DataTable } from "../findings/data-table/data-table";

export default function FindingsOverviewPage() {
    const $activeProjectId = useStore(activeProjectIdStore);
    const [findings, setFindings] = useState([]);
    const [targetsMap, setTargetsMap] = useState({});

    useEffect(() => {
        if($activeProjectId) {
            getProjectFindings($activeProjectId).then(result => {
                if ("error" in result) {
                    return
                }
                console.log(result);
                setFindings(result);
            })

            getProjectTargets($activeProjectId).then(result => {
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
    }, [$activeProjectId])

    return (
        <div className="mt-8 container mx-auto ">
            <h1>Findings</h1>
            <CreateFindingDialog setFindings={setFindings} />
            <div className="py-10">
                <DataTable columns={getColumns(setFindings, findings)} data={findings.map(f => ({...f, target_value: targetsMap[f.target_id]?.value ?? ""}))} />
            </div>
        </div>
    );
}