import { callAPI } from ".";

export const getService = async (serviceId: string, access_token: string) => {
    return await callAPI(`/api/v1/services/${serviceId}`, {
        method: "GET",
        expected_status: 200,
        access_token
    })
}