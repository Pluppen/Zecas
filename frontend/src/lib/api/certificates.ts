import { callAPI } from ".";
import { z } from "zod";

export const CertificateSchema = z.object({
    id: z.string().uuid().optional(),
    project_id: z.string().uuid(),
    target_id: z.string().uuid(),
    service_id: z.string().uuid().optional(),
    application_id: z.string().uuid().optional(),
    scan_id: z.string().uuid(),
    domain: z.string(),
    issuer: z.string(),
    expires_at: z.string().optional(),
    issued_at: z.string().optional(),
    details: z.string().optional(),
    discovered_at: z.string().optional(),
})

export type Certificate = z.infer<typeof CertificateSchema>

export const getCertificate = async (certificateId: string, access_token: string) => {
    return await callAPI(`/api/v1/certificates/${certificateId}`, {
        method: "GET",
        expected_status: 200,
        access_token
    })
}

export const createCertificate = async (certificate: Certificate, access_token: string) => {
    return await callAPI(`/api/v1/certificates`, {
        method: "POST",
        expected_status: 201,
        body: JSON.stringify(certificate),
        access_token
    })
}

export const deleteCertificateById = async (certificateId: string, access_token: string) => {
    return await callAPI(`/api/v1/certificates/${certificateId}`, {
        method: "DELETE",
        expected_status: 200,
        access_token
    })
}
