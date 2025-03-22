import { callAPI } from "@/lib/api";
import {z} from "zod";

export const ApplicationSchema = z.object({
    project_id: z.string().uuid(),
    name: z.string(),
    type: z.string(),
    version: z.string().optional(),
    description: z.string().optional(),
    url: z.string().optional(),
    host_target: z.string().uuid(),
    service_id: z.string().uuid(),
    metadata: z.string().optional()

});

export type Application = z.infer<typeof ApplicationSchema>

export const createApplication = async (application: Application, access_token: string) => {
    return await callAPI(`/api/v1/applications`, {
        method: "POST",
        access_token,
        body: JSON.stringify(application),
        expected_status: 201
    })
}

export const deleteApplicationById = async (application_id: string, access_token: string) => {
    return await callAPI(`/api/v1/applications/${application_id}`, {
        method: "DELETE",
        access_token,
        expected_status: 200
    })
}