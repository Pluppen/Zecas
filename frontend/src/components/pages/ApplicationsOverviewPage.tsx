import {useState, useEffect} from "react"

import {user} from "@/lib/userStore";
import { activeProjectIdStore } from "@/lib/projectsStore";
import { useStore } from "@nanostores/react";

import { getProjectApplications, getProjectTargets } from "@/lib/api/projects";

import { getColumns } from "@/components/applications/data-table/columns";
import { DataTable } from "@/components/findings/data-table/data-table";

import CreateApplicationDialog from "@/components/applications/create-application-dialog";


export default function ApplicationsOverviewPage() {
    const $activeProjectId = useStore(activeProjectIdStore);
    const [applications, setApplications] = useState([]);
    const [targetsMap, setTargetsMap] = useState({});
    const $user = useStore(user);

    useEffect(() => {
        if($activeProjectId) {
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
        <div className="mt-4 container mx-auto ">
            <CreateApplicationDialog setApplications={setApplications} />
            <div className="py-10">
                <DataTable
                    columns={getColumns(setApplications, applications)}
                    data={applications.map(s => ({...s, target: targetsMap[s.host_target]?.value ?? ""}))}
                    filterSettings={{placeholder: "Filter by description...", filterKey: "description"}}
                />
            </div>
        </div>
    );
}