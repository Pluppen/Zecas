import {useState, useEffect} from "react"

import { activeProjectIdStore } from "@/lib/projectsStore";
import { user } from "@/lib/userStore";
import { useStore } from "@nanostores/react";

import { getProjectScans, getScanConfigs } from "@/lib/scans";

import SimpleTable from "@/components/simple-table";

export default function ScansOverviewPage() {
    const $activeProjectId = useStore(activeProjectIdStore);
    const $user = useStore(user);
    const [scans, setScans] = useState([]);
    const [scanConfigs, setScanConfigs] = useState({});

    useEffect(() => {
        if($activeProjectId) {
            getProjectScans($activeProjectId, $user.access_token).then(result => {
                if ("error" in result) {
                    return
                }
                setScans(result);
            })

            getScanConfigs($user.access_token).then(result => {
                if ("error" in result) {
                    return
                }
                const scanConfigHT = {}
                result.forEach(r => {
                    scanConfigHT[r.id] = r
                });
                setScanConfigs(scanConfigHT);
            })
        }
    }, [$activeProjectId, $user])

    return (
        <div className="mt-8">
            <h1>Scans Table</h1>
            <SimpleTable
                tableCaption="A list of this projects scans"
                headers={[
                    {
                        key: "short_id",
                        label: "Scan ID"
                    },
                    {
                        key: "status",
                        label: "Status"
                    },
                    {
                        key: "scan_config_name",
                        label: "Scan Config Name"
                    },
                    {
                        key: "linkComponent",
                        label: "Link"
                    }
                ]}
                tableRows={scans.map(s => ({
                    ...s, 
                    short_id: s.id.substring(0, 8),
                    scan_config_name: s.scan_config_id in scanConfigs ? scanConfigs[s.scan_config_id].name : "",
                    linkComponent: <a className="underline hover:cursor-pointer" href={`/scans/${s.id}`}>View Scan</a>
                }))}
            />
        </div>
    );
}