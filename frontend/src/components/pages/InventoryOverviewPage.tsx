import {useState, useEffect} from "react"

import {user} from "@/lib/userStore";
import { activeProjectIdStore } from "@/lib/projectsStore";
import { useStore } from "@nanostores/react";

import { getProjectServices, getProjectTargets, getProjectDNSRecords, getProjectApplications, getProjectCertificates } from "@/lib/api/projects";
import { type Target } from "@/lib/api/targets";
import { type Application } from "@/lib/api/applications";
import { type DNSRecord } from "@/lib/api/dns";
import { type Certificate } from "@/lib/api/certificates";
import { type Service } from "@/lib/api/services";

import { BasicCard } from "@/components/basic-card";

export default function ServicesOverviewPage() {
    const $activeProjectId = useStore(activeProjectIdStore);
    const [services, setServices] = useState<Service[]>([]);
    const [targets, setTargets] = useState<Target[]>([]);
    const [applications, setApplications] = useState<Application[]>([]);
    const [dnsRecords, setDNSRecords] = useState<DNSRecord[]>([]);
    const [certificates, setCertificates] = useState<Certificate[]>([]);
    const $user = useStore(user);

    useEffect(() => {
        if($activeProjectId && $user?.access_token) {
            getProjectServices($activeProjectId, $user.access_token).then(services => {
                if ("error" in services) {
                    return
                }
                setServices(services);
            })

            getProjectDNSRecords($activeProjectId, $user.access_token).then(dnsRecordsData => {
                if ("error" in dnsRecordsData) {
                    return
                }
                setDNSRecords(dnsRecordsData);
            })

            getProjectCertificates($activeProjectId, $user.access_token).then(certificatesData => {
                if ("error" in certificatesData) {
                    return
                }
                setCertificates(certificatesData)
            })

            getProjectApplications($activeProjectId, $user.access_token).then(applicationsData => {
                if ("error" in applicationsData) {
                    return
                }
                setApplications(applicationsData)
            })

            getProjectTargets($activeProjectId, $user.access_token).then(result => {
                if ("error" in result) {
                    return
                }
                setTargets(result);
            });
        }
    }, [$activeProjectId, $user])

    return (
        <div className="mt-4 container">
            <BasicCard cardTitle="Inventory Items" cardDescription="The count of items and different assets you have">
                <ul>
                    <li>Hosts: {targets.length}</li>
                    <li>DNS Records: {dnsRecords.length}</li>
                    <li>Certificates: {certificates.length}</li>
                    <li>Services: {services.length}</li>
                    <li>Applications: {applications.length}</li>
                </ul>
            </BasicCard>
        </div>
    );
}
