export type ApiOptions = {
    access_token: string,
    expected_status: number,
    method: string,
    body?: string
}

const CLIENT_API_URL = import.meta.env.PUBLIC_API_URL
const SSR_API_URL = import.meta.env.API_URL

export const callAPI = async (path: string, options: ApiOptions) => {
    let API_URL = CLIENT_API_URL
    if (import.meta.env.SSR) {
        API_URL = SSR_API_URL
    }

    const response = await fetch(`${API_URL}${path}`, {
        method: options.method,
        headers: {
            "Authorization": `Bearer ${options.access_token}`
        },
        body: options.body ?? undefined
    })
    if (response.status != options.expected_status) {
        console.log(response);
        return {error: "Something went wrong."};
    }
    return response.json();
}