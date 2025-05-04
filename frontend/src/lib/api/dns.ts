import { callAPI } from ".";
import { z } from "zod";

export const DNSRecordSchema = z.object({
    id: z.string().uuid().optional(),
    project_id: z.string().uuid(),
    target_id: z.string().uuid(),
    scan_id: z.string().uuid(),
    record_type: z.string(),
    record_value: z.string().optional(),
    discovered_at: z.string().optional(),
    details: z.string().optional(),
})

export type DNSRecord = z.infer<typeof DNSRecordSchema>

export const getDNSRecord = async (dnsRecordId: string, access_token: string) => {
    return await callAPI(`/api/v1/dns-records/${dnsRecordId}`, {
        method: "GET",
        expected_status: 200,
        access_token
    })
}

export const createDNSRecord = async (dnsRecord: DNSRecord, access_token: string) => {
    return await callAPI(`/api/v1/dns-records`, {
        method: "POST",
        expected_status: 201,
        body: JSON.stringify(dnsRecord),
        access_token
    })
}

export const deleteDNSRecordById = async (dnsRecordId: string, access_token: string) => {
    return await callAPI(`/api/v1/dns-records/${dnsRecordId}`, {
        method: "DELETE",
        expected_status: 200,
        access_token
    })
}
