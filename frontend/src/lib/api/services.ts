import { z } from "zod";
import { callAPI } from ".";

export const ServiceSchema = z.object({
    id: z.string().uuid().optional(),
    target_id: z.string().uuid(),
    port: z.number().min(1).max(65535),
    protocol: z.string(),
    service_name: z.string().optional(),
    version: z.string().optional(),
    title: z.string().optional(),
    description: z.string().optional(),
    banner: z.string().optional(),
    raw_info: z.string().optional(),
});

export type Service = z.infer<typeof ServiceSchema>

export const getService = async (serviceId: string, access_token: string) => {
    return await callAPI(`/api/v1/services/${serviceId}`, {
        method: "GET",
        expected_status: 200,
        access_token
    })
}
