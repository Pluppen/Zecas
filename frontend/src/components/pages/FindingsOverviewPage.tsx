import {useState, useEffect} from "react"

import { activeProjectIdStore } from "@/lib/projectsStore";
import { useStore } from "@nanostores/react";

import { getProjectFindings } from "@/lib/findings";

import SimpleTable from "@/components/simple-table";
import PieChart from "@/components/pie-chart";

export default function FindingsOverviewPage() {
    const $activeProjectId = useStore(activeProjectIdStore);
    const [findings, setFindings] = useState([]);
    const [scanConfigs, setScanConfigs] = useState({});

    useEffect(() => {
        if($activeProjectId) {
            getProjectFindings($activeProjectId).then(result => {
                if ("error" in result) {
                    return
                }
                console.log(result);
                setFindings(result);
            })
        }
    }, [$activeProjectId])

    return (
        <div className="mt-8 w-full">
            <h1>Findings Table</h1>
            <div className="grid auto-rows-min gap-4 md:grid-cols-3">
                <div className="col-span-3">
                    <SimpleTable
                        tableCaption="A list of this projects findings"
                        headers={[
                            {
                                key: "short_id",
                                label: "Finding ID"
                            },
                            {
                                key: "severity",
                                label: "Severity"
                            },
                            {
                                key: "title",
                                label: "Title"
                            },
                            {
                                key: "short_target_id",
                                label: "Target ID"
                            },
                            {
                                key: "finding_type",
                                label: "Finding Type"
                            },
                            {
                                key: "discovered_at",
                                label: "Time Generated"
                            },
                            {
                                key: "verified",
                                label: "Has been verified"
                            },
                            {
                                key: "fixed",
                                label: "Has been fixed"
                            }
                        ]}
                        tableRows={findings.map(s => ({
                            ...s,
                            short_id: s.id.substring(0,8),
                            short_target_id: s.id.substring(0, 8),
                            fixed: s.fixed ? "Yes" : "No",
                            verified: s.verified ? "Yes" : "No",
                        }))}
                    />
                </div>
                <div className="">
                    <PieChart />
                </div>
            </div>
        </div>
    );
}