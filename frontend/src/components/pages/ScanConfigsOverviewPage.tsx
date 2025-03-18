import {useState} from "react";

import CreateScanDialog from '@/components/create-scan-dialog';
import CreateScanConfigDialog from '@/components/create-scan-config-dialog';

import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle} from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { type ScanConfig } from "@/lib/api/scans";

interface ScanConfigOverviewPageProps {
    scanConfigs: ScanConfig[]
}

export default function ScanConfigOverviewPage(props: ScanConfigOverviewPageProps) {
    const [scanConfigs, setScanConfigs] = useState(props.scanConfigs);


    return (
        <div className="w-2/3">
            <h1 className="mb-4">Scan Configs</h1>
            <CreateScanConfigDialog setScanConfigs={setScanConfigs} />
            <div className="grid grid-cols-1 lg:grid-cols-3 xl:grid-cols-4 gap-4 mt-8"> 
                {scanConfigs.map(scanConfig => (
                    <Card key={scanConfig.id}>
                        <CardHeader>
                            <CardTitle>{scanConfig.name}</CardTitle>
                            <CardDescription>
                                <Badge className="mt-2 uppercase" variant="outline">{scanConfig.scanner_type}</Badge>
                            </CardDescription>
                        </CardHeader>
                        <CardFooter>
                            <Button variant="secondary">Edit</Button>
                        </CardFooter>
                    </Card>
                ))}
            </div>
        </div>
    );
}