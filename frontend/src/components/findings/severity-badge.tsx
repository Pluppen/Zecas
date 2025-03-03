import { Badge } from "@/components/ui/badge"

export const SEVERITY_COLOR = {
    "unknown": "bg-gray-400 text-white",
    "info": "bg-blue-600 text-white",
    "low": "bg-green-600 text-white",
    "medium": "bg-yellow-600 text-white",
    "high": "bg-red-600 text-white",
    "critical": "bg-purple-800 text-white",
}

export type Severity = "unknown" | "info" | "low" | "medium" | "high" | "critical"

export interface SeverityBadgeProps {
    severity: Severity,
    className: string
}

export default function SeverityBadge ({severity, className}: SeverityBadgeProps) {
    return (
        <Badge className={`${className} ${SEVERITY_COLOR[severity]}`}>{severity}</Badge>
    );
}