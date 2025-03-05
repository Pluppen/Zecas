import { API_URL } from "@/config"

export type ApiOptions = {
    access_token: string,
    expected_status: number,
    method: string,
    body?: string
}

export const callAPI = async (path: string, options: ApiOptions) => {
    const response = await fetch(`${API_URL}${path}`, {
        method: options.method,
        headers: {
            "Authorization": `Bearer ${options.access_token}`
        },
        body: options.body ?? undefined
    })
    if (response.status != options.expected_status) {
        return {error: "Something went wrong fetching scans for project"};
    }
    return response.json();
}