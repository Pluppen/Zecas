import {useState, useEffect} from "react"

import {user} from "@/lib/userStore";
import { activeProjectIdStore } from "@/lib/projectsStore";
import { useStore } from "@nanostores/react";

import { getProjectDNSRecords, getProjectTargets } from "@/lib/api/projects";
import { type Target } from "@/lib/api/targets";
import { type DNSRecord, deleteDNSRecordById } from "@/lib/api/dns";

import { getColumns } from "@/components/dnsRecords/data-table/columns";
import { DataTable } from "@/components/findings/data-table/data-table";

import { toast } from "sonner";


export default function DNSOverviewPage() {
    const $activeProjectId = useStore(activeProjectIdStore);
    const [dnsRecords, setDNSRecords] = useState<DNSRecord[]>([]);
    const [targetsMap, setTargetsMap] = useState<Record<string, Target>>({});
    const $user = useStore(user);

    useEffect(() => {
        if($activeProjectId && $user?.access_token) {
            getProjectDNSRecords($activeProjectId, $user.access_token).then(dnsRecords => {
                if ("error" in dnsRecords) {
                    return
                }
                setDNSRecords(dnsRecords);
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
    }, [$activeProjectId, $user])

    const handleBulkDelete = async (selectedRows: any) => {
        const removedDNSRecordIds: string[] = []
        for(let row of selectedRows) {
            if (row.id && $user?.access_token != undefined) {
                const resp = await deleteDNSRecordById(row.id, $user?.access_token)
                if ("error" in resp) {
                    toast("Something went wrong")
                }
                removedDNSRecordIds.push(row.id)
            }
        }
        const newDNSRecords = dnsRecords.slice().filter(t => !removedDNSRecordIds.includes(t.id ?? ""))
        setDNSRecords(newDNSRecords);
        toast(`Successfully removed ${selectedRows.length} DNS Records.`)
    }

    return (
        <div className="mt-4 container">
            <div className="py-10">
                <DataTable
                    columns={getColumns(setDNSRecords, dnsRecords)}
                    data={dnsRecords.map(dnsRecord => ({...dnsRecord, target: targetsMap[dnsRecord.target_id]?.value ?? "Unknown"}))}
                    filterSettings={{placeholder: "Filter by target", filterKey: "target"}}
                    handleDelete={handleBulkDelete}
                />
            </div>
        </div>
    );
}
