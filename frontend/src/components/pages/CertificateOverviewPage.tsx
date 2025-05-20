import {useState, useEffect} from "react"

import {user} from "@/lib/userStore";
import { activeProjectIdStore } from "@/lib/projectsStore";
import { useStore } from "@nanostores/react";

import { getProjectCertificates, getProjectTargets } from "@/lib/api/projects";
import { type Target } from "@/lib/api/targets";
import { type Certificate, deleteCertificateById } from "@/lib/api/certificates";

import { getColumns } from "@/components/certificates/data-table/columns";
import { DataTable } from "@/components/findings/data-table/data-table";

import { toast } from "sonner";


export default function CertificateOverviewPage() {
    const $activeProjectId = useStore(activeProjectIdStore);
    const [certificates, setCertificates] = useState<Certificate[]>([]);
    const [targetsMap, setTargetsMap] = useState<Record<string, Target>>({});
    const $user = useStore(user);

    useEffect(() => {
        if($activeProjectId && $user?.access_token) {
            getProjectCertificates($activeProjectId, $user.access_token).then(certificates => {
                if ("error" in certificates) {
                    return
                }
                setCertificates(certificates);
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
        const removedCertificateIds: string[] = []
        for(let row of selectedRows) {
            if (row.id && $user?.access_token != undefined) {
                const resp = await deleteCertificateById(row.id, $user?.access_token)
                if ("error" in resp) {
                    toast("Something went wrong")
                }
                removedCertificateIds.push(row.id)
            }
        }
        const newCertificates = certificates.slice().filter(t => !removedCertificateIds.includes(t.id ?? ""))
        setCertificates(newCertificates);
        toast(`Successfully removed ${selectedRows.length} certificates.`)
    }

    return (
        <div className="mt-4 container">
            <div className="py-10">
                <DataTable
                    columns={getColumns(setCertificates, certificates)}
                    data={certificates.map(certificate => ({...certificate, target: targetsMap[certificate.target_id]?.value ?? "Unknown"}))}
                    filterSettings={{placeholder: "Filter by target", filterKey: "target"}}
                    handleDelete={handleBulkDelete}
                />
            </div>
        </div>
    );
}
