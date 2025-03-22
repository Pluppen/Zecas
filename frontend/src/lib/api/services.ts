import { callAPI } from ".";
import { z } from "zod";

export const ServiceSchema = z.object({
    id: z.string().uuid().optional(),
    target_id: z.string().uuid(),
    port: z.coerce.number().min(1).max(65535),
    protocol: z.string(),
    service_name: z.string().optional(),
    version: z.string().optional(),
    title: z.string().optional(),
    description: z.string().optional(),
    banner: z.string().optional(),
    raw_info: z.string().optional(),
})

export type Service = z.infer<typeof ServiceSchema>

export const getService = async (serviceId: string, access_token: string) => {
    return await callAPI(`/api/v1/services/${serviceId}`, {
        method: "GET",
        expected_status: 200,
        access_token
    })
}

export const createService = async (service: Service, access_token: string) => {
    return await callAPI(`/api/v1/services`, {
        method: "POST",
        expected_status: 201,
        body: JSON.stringify(service),
        access_token
    })
}

export const deleteServiceById = async (serviceId: string, access_token: string) => {
    return await callAPI(`/api/v1/services/${serviceId}`, {
        method: "DELETE",
        expected_status: 200,
        access_token
    })
}