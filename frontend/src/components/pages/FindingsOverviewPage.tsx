import {useState, useEffect} from "react"

import { activeProjectIdStore } from "@/lib/projectsStore";
import { useStore } from "@nanostores/react";
import {user} from "@/lib/userStore";

import { getProjectFindings } from "@/lib/api/projects";

import SimpleTable from "@/components/simple-table";
import PieChart from "@/components/pie-chart";

export default function FindingsOverviewPage() {
    const $activeProjectId = useStore(activeProjectIdStore);
    const $user = useStore(user);

    const [findings, setFindings] = useState([]);
    const [severityCounts, setSeverityCounts] = useState({})

    useEffect(() => {
        if($activeProjectId) {
            getProjectFindings($activeProjectId, $user.access_token).then(result => {
                if ("error" in result) {
                    return
                }
                setFindings(result);
                let severityCountsTmp = {...severityCounts}
                result.forEach(f => {
                    if (f.severity in severityCountsTmp) {
                        severityCountsTmp[f.severity] += 1
                    } else {
                        severityCountsTmp[f.severity] = 1
                    }
                });
                setSeverityCounts(severityCountsTmp)
            })
        }
    }, [$activeProjectId, $user])

    const chartData = [
        { severity: "critical", count: severityCounts["critical"] ?? 0, fill: "var(--color-critical)" },
        { severity: "high", count: severityCounts["high"] ?? 0, fill: "var(--color-high)" },
        { severity: "medium", count: severityCounts["medium"] ?? 0, fill: "var(--color-medium)" },
        { severity: "low", count: severityCounts["low"] ?? 0, fill: "var(--color-low)" },
        { severity: "unknown", count: severityCounts["unknown"] ?? 0, fill: "var(--color-unknown)" },
    ]

    const chartConfig = {
        count: {
            label: "Count",
        },
        unknown: {
            label: "Unknown",
            color: "hsl(var(--chart-1))",
        },
        low: {
            label: "Low",
            color: "hsl(var(--chart-2))",
        },
        medium: {
            label: "Medium",
            color: "hsl(var(--chart-3))",
        },
        high: {
            label: "High",
            color: "hsl(var(--chart-4))",
        },
        critical: {
            label: "Critical",
            color: "hsl(var(--chart-5))",
        },
    } satisfies ChartConfig

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
                <div className="">
                    <PieChart chartConfig={chartConfig} chartData={chartData} />
                </div>
            </div>
        </div>
    );
}