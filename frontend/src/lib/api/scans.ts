import { callAPI } from "@/lib/api";
import { z } from "zod";

export const ScanSchema = z.object({
    id: z.string().optional(),
    project_id: z.string().uuid(),
    scan_config_id: z.string().uuid(),
    status: z.string().optional(),
    raw_results: z.string().optional(),
    error: z.string().optional(),
});

export type Scan = z.infer<typeof ScanSchema>;

export type ScannerType =
    | "nmap"
    | "dns"
    | "subdomain"
    | "nuclei"
    | "httpx"
    | "testSSL";

export const ScannerTypeEnum = z.enum([
    "nmap",
    "dns",
    "subdomain",
    "nuclei",
    "httpx",
    "testSSL",
]);

const NmapParametersSchema = z.object({
    scan_type: z
        .enum(["quick", "comprehensive", "service", "all_ports", "basic"], {
            message:
                "Parameter scan_type can only be quick, comprehensive, service, all_ports or basic",
        })
        .optional(),
    port_range: z
        .string({
            message: "Parameter port_range is required to be a valid string",
        })
        .optional(),
    timing: z
        .enum(["0", "1", "2", "3", "4", "5"], {
            message: "Parameter timing can only be a number between 0-5",
        })
        .optional(),
});

const allowedRecordTypes = ["A", "AAAA", "CNAME", "MX", "TXT", "NS"] as const;
const DNSParametersSchema = z.object({
    record_types: z
        .array(
            z.enum(allowedRecordTypes, {
    rate_limit: z.number({message: "Parameter rate_limit needs to be a valid number"}).optional(),
    bulk_size: z.number({message: "Parameter bulk_size needs to be a valid number"}).optional(),
    templates_dir: z.string({message: "Parameter templates_dir needs to be a valid string"}).optional(),
    headless: z.boolean({message: "Parameter headless needs to be either true or false"}).optional(),
    include_all: z.boolean({message: "Parameter include_all needs to be either true or false"}).optional(),
                message:
                    "Parameter record_types needs to be an array of valid DNS records",
            })
        )
        .optional(),
});

const SubdomainParametersSchema = z.object({
    recursive: z
        .boolean({
            message: "Parameter recursive needs to be either true or false",
        })
        .optional(),
    resolve_ip: z
        .boolean({
            message: "Parameter resolve_ip needs to be either true or false",
        })
        .optional(),
    wordlist: z
        .string({ message: "Parameter wordlist needs to be a valid string" })
        .optional(),
    timeout: z
        .number({ message: "Parameter timeout needs to be a valid number" })
        .optional(),
});

const NucleiParametersSchema = z.object({
    template_tags: z
        .array(
            z.string({
                message: "Parameter template_tags needs to be a list of strings",
            })
        )
        .optional(),
    template_paths: z
        .array(
            z.string({
                message: "Parameter template_paths needs to be a list of strings",
            })
        )
        .optional(),
    template_exclude: z
        .array(
            z.string({
                message: "Parameter template_exclude needs to be a list of strings",
            })
        )
        .optional(),
    severity: z
        .string({ message: "Parameter severity needs to be a valid string" })
        .optional(),
    timeout: z
        .number({ message: "Parameter timeout needs to be a valid number" })
        .optional(),
    rate_limit: z
        .number({ message: "Parameter rate_limit needs to be a valid number" })
        .optional(),
    bulk_size: z
        .number({ message: "Parameter bulk_size needs to be a valid number" })
        .optional(),
    templates_dir: z
        .string({ message: "Parameter templates_dir needs to be a valid string" })
        .optional(),
    headless: z
        .boolean({ message: "Parameter headless needs to be either true or false" })
        .optional(),
    include_all: z
        .boolean({
            message: "Parameter include_all needs to be either true or false",
        })
        .optional(),
});

const TestSSLParametersSchema = z.object({
    severity: z
        .string({ message: "Parameter severity needs to be a valid string" })
        .optional(),
    timeout: z
        .number({ message: "Parameter timeout needs to be a valid number" })
        .optional(),
});

const HttpxParametersSchema = z.object({
    timeout: z
        .number({ message: "Timeout needs to be a valid number" })
        .optional(),
    threads: z
        .number({ message: "Timeout needs to be a valid number" })
        .optional(),
    follow_redirects: z
        .boolean({
            message: "Parameter include_all needs to be either true or false",
        })
        .optional(),
    tech_detect: z
        .boolean({
            message: "Parameter tech_detect needs to be either true or false",
        })
        .optional(),
    status_code: z
        .boolean({
            message: "Parameter status_code needs to be either true or false",
        })
        .optional(),
    title: z
        .boolean({ message: "Parameter title needs to be either true or false" })
        .optional(),
    web_server: z
        .boolean({
            message: "Parameter web_server needs to be either true or false",
        })
        .optional(),
    content_type: z
        .boolean({
            message: "Parameter content_type needs to be either true or false",
        })
        .optional(),
    tls: z
        .boolean({ message: "Parameter tls needs to be either true or false" })
        .optional(),
    favicon: z
        .boolean({ message: "Parameter favicon needs to be either true or false" })
        .optional(),
    jarm: z
        .boolean({ message: "Parameter jarm needs to be either true or false" })
        .optional(),
    probe: z
        .boolean({ message: "Parameter probe needs to be either true or false" })
        .optional(),
    ports: z
        .string({ message: "Parameter ports needs to be a valid string" })
        .optional(),
    http2: z
        .boolean({ message: "Parameter http2 needs to be either true or false" })
        .optional(),
    security_headers: z
        .boolean({
            message: "Parameter security_headers needs to be either true or false",
        })
        .optional(),
    extract_cname: z
        .boolean({
            message: "Parameter extract_cname needs to be either true or false",
        })
        .optional(),
});

const ScanConfigBaseSchema = z.object({
    id: z.string().optional(),
    name: z.string(),
    active: z.boolean(),
});

export const ScanConfigSchema = z.discriminatedUnion("scanner_type", [
    ScanConfigBaseSchema.extend({
        scanner_type: z.literal("nmap"),
        parameters: NmapParametersSchema,
    }),
    ScanConfigBaseSchema.extend({
        scanner_type: z.literal("dns"),
        parameters: DNSParametersSchema,
    }),
    ScanConfigBaseSchema.extend({
        scanner_type: z.literal("subdomain"),
        parameters: SubdomainParametersSchema,
    }),
    ScanConfigBaseSchema.extend({
        scanner_type: z.literal("nuclei"),
        parameters: NucleiParametersSchema,
    }),
    ScanConfigBaseSchema.extend({
        scanner_type: z.literal("testSSL"),
        parameters: TestSSLParametersSchema,
    }),
    ScanConfigBaseSchema.extend({
        scanner_type: z.literal("httpx"),
        parameters: HttpxParametersSchema,
    }),
]);
export type ScanConfig = z.infer<typeof ScanConfigSchema>;

export const getScanConfigs = async (access_token: string) => {
    return await callAPI("/api/v1/scan-configs", {
        method: "GET",
        expected_status: 200,
        access_token,
    });
};

export const createScanConfig = async (
    scanConfig: ScanConfig,
    access_token: string
) => {
    const body = {
        name: scanConfig.name,
        scanner_type: scanConfig.scanner_type,
        parameters: scanConfig.parameters,
        active: scanConfig.active,
    };

    return await callAPI(`/api/v1/scan-configs`, {
        method: "POST",
        access_token,
        body: JSON.stringify(body),
        expected_status: 201,
    });
};

export const getScanConfigById = async (
    scanConfigId: string,
    access_token: string
) => {
    return await callAPI(`/api/v1/scan-configs/${scanConfigId}`, {
        method: "GET",
        access_token,
        expected_status: 200,
    });
};

export const updateScanConfigById = async (
    scanConfig: ScanConfig,
    access_token: string
) => {
    const body = {
        name: scanConfig.name,
        scanner_type: scanConfig.scanner_type,
        parameters: scanConfig.parameters,
        active: scanConfig.active,
    };

    return await callAPI(`/api/v1/scan-configs/${scanConfig.id}`, {
        method: "PUT",
        access_token,
        body: JSON.stringify(body),
        expected_status: 200,
    });
};

export const deleteScanConfigById = async (
    scanId: string,
    access_token: string
) => {
    return await callAPI(`/api/v1/scan-configs/${scanId}`, {
        method: "DELETE",
        access_token,
        expected_status: 200,
    });
};

// SCANS

export const startNewScan = async (
    projectId: string,
    scanConfigId: string,
    targetIds: string[],
    access_token: string
) => {
    const body = {
        project_id: projectId,
        scan_config_id: scanConfigId,
        target_ids: targetIds,
    };
    return await callAPI("/api/v1/scans", {
        method: "POST",
        expected_status: 202,
        body: JSON.stringify(body),
        access_token,
    });
};

export const getScanById = async (scanId: string, access_token: string) => {
    return await callAPI(`/api/v1/scans/${scanId}`, {
        method: "GET",
        access_token,
        expected_status: 200,
    });
};

export const getScans = async (access_token: string) => {
    return await callAPI(`/api/v1/scan/`, {
        method: "GET",
        access_token,
        expected_status: 200,
    });
};

export const cancelScanById = async (scanId: string, access_token: string) => {
    return await callAPI(`/api/v1/scan-configs/${scanId}/cancel`, {
        method: "POST",
        access_token,
        expected_status: 200,
    });
};

export const getFindingsByScanId = async (
    scanId: string,
    access_token: string
) => {
    return await callAPI(`/api/v1/scan-configs/${scanId}/findings`, {
        method: "POST",
        access_token,
        expected_status: 200,
    });
};
