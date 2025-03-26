import { callAPI } from "@/lib/api";
import { z } from "zod";

export const FindingSchema = z.object({
    id: z.string(),
    scan_id: z.string().uuid().optional(),
    target_id: z.string().uuid(),
    service_id: z.string().uuid().optional(),
    application_id: z.string().uuid().optional(),
    title: z.string(),
    description: z.string().optional(),
    severity: z.string(),
    finding_type: z.string(),
    details: z.string().optional(),
    verified: z.boolean().optional(),
    fixed: z.boolean().optional(),
    manual: z.boolean().optional()
});

export type Finding = z.infer<typeof FindingSchema>;

export const getFindings = async (access_token: string) => {
    return await callAPI(`/api/v1/findings`, {
        method: "GET",
        access_token,
        expected_status: 200
    })
}

export const getFindingById = async (findingId: string, access_token: string) => {
    return await callAPI(`/api/v1/findings/${findingId}`, {
        method: 'GET',
        expected_status: 200,
        access_token
    })
}

interface FindingParam {
    id?: string,
    target_id: string,
    title: string,
    description?: string,
    severity: string,
    finding_type: string,
    details?: string,
    manual?: boolean,
    scan_id?: string
}

export const createFinding = async (finding: FindingParam, access_token: string) => {
    const body: FindingParam = {
        target_id: finding.target_id,
        title: finding.title,
        severity: finding.severity,
        finding_type: finding.finding_type,
    }

    if (finding.description) {
        body.description = finding.description
    }

    if (finding.scan_id) {
        body.scan_id= finding.scan_id
    }

    if (finding.details) {
        body.details = finding.details
    }

    if (finding.manual) {
        body.manual = finding.manual
    }

    return await callAPI("/api/v1/findings", {
        method: 'POST',
        access_token,
        body: JSON.stringify(body),
        expected_status: 201
    })
}

export const updateFinding = async (finding: FindingParam, access_token: string) => {
    const body: FindingParam = {
        target_id: finding.target_id,
        title: finding.title,
        severity: finding.severity,
        finding_type: finding.finding_type,
    }

    if (finding.description) {
        body.description = finding.description
    }

    if (finding.details) {
        body.details = finding.details
    }

    if (finding.manual) {
        body.manual = finding.manual
    }

    return await callAPI(`/api/v1/findings/${finding.id}`, {
        method: 'PUT',
        access_token,
        body: JSON.stringify(body),
        expected_status: 200
    })
}

export const removeFinding = async (findingId: string, access_token: string) => {
    return await callAPI(`/api/v1/findings/${findingId}`, {
        method: 'DELETE',
        access_token,
        expected_status: 200
    });
}