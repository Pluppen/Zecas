import {useState, useEffect} from "react"

import { activeProjectIdStore } from "@/lib/projectsStore";
import { useStore } from "@nanostores/react";
import {user} from "@/lib/userStore";

import { getProjectFindings } from "@/lib/api/projects";

import { type Finding } from "@/lib/api/findings"

import SimpleTable from "@/components/simple-table";

export default function FindingsOverviewPage() {
    const $activeProjectId = useStore(activeProjectIdStore);
    const $user = useStore(user);

    const [findings, setFindings] = useState<Finding[]>([]);
    const [severityCounts, setSeverityCounts] = useState<Record<string, number>>({})

    useEffect(() => {
        if($activeProjectId && $user?.access_token) {
            getProjectFindings($activeProjectId, $user.access_token).then(result => {
                if ("error" in result) {
                    return
                }
                setFindings(result);
                let severityCountsTmp = {...severityCounts}
                result.forEach((f: Finding) => {
                    if (f.severity in severityCountsTmp) {
                        severityCountsTmp[f.severity] += 1
                    } else {
                        severityCountsTmp[f.severity] = 1
                    }
                });
                setSeverityCounts(severityCountsTmp)
            })
        }
    }, [$activeProjectId, $user?.access_token])

    return (
        <div className="mt-8 w-full">
            <h1>Findings</h1>
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
            </div>
        </div>
    );
}