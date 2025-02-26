import {useState, useEffect} from "react"

import { activeProjectIdStore } from "@/lib/projectsStore";
import { useStore } from "@nanostores/react";

import { getProjectScans, getScanConfigs } from "@/lib/scans";

import SimpleTable from "@/components/simple-table";

export default function ScansOverviewPage() {
    const $activeProjectId = useStore(activeProjectIdStore);
    const [scans, setScans] = useState([]);
    const [scanConfigs, setScanConfigs] = useState({});

    useEffect(() => {
        if($activeProjectId) {
            getProjectScans($activeProjectId).then(result => {
                if ("error" in result) {
                    return
                }
                setScans(result);
            })

            getScanConfigs().then(result => {
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
    }, [$activeProjectId])

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
                    }
                ]}
                tableRows={scans.map(s => ({...s, short_id: s.id.substring(0, 8),scan_config_name: scanConfigs[s.scan_config_id].name}))}
            />
        </div>
    );
}