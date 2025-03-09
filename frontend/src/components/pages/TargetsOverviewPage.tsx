import {useState, useEffect} from "react"

import {user} from "@/lib/userStore";
import { activeProjectIdStore } from "@/lib/projectsStore";
import { useStore } from "@nanostores/react";

import { getProjectTargets } from "@/lib/api/projects";
import {getTargetFindings, getTargetServices} from "@/lib/api/targets"

import CreateFindingDialog from "@/components/findings/create-finding-dialog";

import { getColumns } from "@/components/targets/data-table/columns";
import { DataTable } from "@/components/findings/data-table/data-table";


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
                console.log(targets);
                let tmpTargets = [];
                for (const target of targets) {
                    const findings = await getTargetFindings(target.id, $user.access_token);
                    const services = await getTargetServices(target.id, $user.access_token);
                    tmpTargets.push({
                        ...target,
                        findings,
                        services
                    })
                };

                setTargets(tmpTargets);
            });
        }
    }, [$activeProjectId, $user])

    return (
        <div className="mt-4 container mx-auto ">
            <div className="py-10">
                <DataTable columns={getColumns(setTargets, targets)} data={targets} />
            </div>
        </div>
    );
}