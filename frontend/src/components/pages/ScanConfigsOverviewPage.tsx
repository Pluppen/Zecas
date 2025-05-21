import {useState, useEffect} from "react";

import CreateScanConfigDialog from '@/components/create-scan-config-dialog';
import EditScanConfigDialog from '@/components/edit-scan-config-dialog';
import { Card, CardDescription, CardFooter, CardHeader, CardTitle} from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { deleteScanConfigById, type ScanConfig } from "@/lib/api/scans";
import { Trash } from "lucide-react";
import { RemoveItemDialog } from "../remove-item-dialog";
import { user } from "@/lib/userStore";
import { useStore } from "@nanostores/react";
import {toast} from "sonner";

interface ScanConfigOverviewPageProps {
    scanConfigs: ScanConfig[]
}

export default function ScanConfigOverviewPage(props: ScanConfigOverviewPageProps) {
    const [scanConfigs, setScanConfigs] = useState(props.scanConfigs);
    const [scanConfig, setScanConfig] = useState<ScanConfig>();
    const [editDialogOpen, setEditDialogOpen] = useState(false);

    const $user = useStore(user);

    useEffect(() => {
        if(scanConfig !== undefined) {
            console.log("Ready to edit");
            setEditDialogOpen(true)
        } else {
            console.log("Not Ready to edit");
        }
    }, [scanConfig])

    useEffect(() => {
        if(!editDialogOpen) {
            setScanConfig(undefined);
        }
    }, [editDialogOpen])

    return (
        <div className="container">
            <h1 className="mb-4">Scan Configs</h1>
            <CreateScanConfigDialog setScanConfigs={setScanConfigs} />
            {scanConfig ? <EditScanConfigDialog
                setScanConfigs={setScanConfigs}
                scanConfigs={scanConfigs}
                scanConfig={scanConfig}
                open={editDialogOpen}
                setOpen={setEditDialogOpen}
            /> : null}
            <div className="grid grid-cols-1 lg:grid-cols-3 xl:grid-cols-4 gap-4 mt-8"> 
                {scanConfigs.map(scanConfig => (
                    <Card key={scanConfig.id} className="flex flex-column">
                        <CardHeader>
                            <CardTitle>{scanConfig.name}</CardTitle>
                            <CardDescription>
                                <Badge className="mt-2 uppercase" variant="outline">{scanConfig.scanner_type}</Badge>
                            </CardDescription>
                        </CardHeader>
                        <CardFooter className="flex-grow">
                            <RemoveItemDialog
                                handleSubmit={async () => {
                                    if ($user?.access_token && scanConfig.id) {
                                        const result = await deleteScanConfigById(scanConfig.id, $user.access_token)
                                        if ("error" in result) {
                                            toast(result.error);
                                            return
                                        }
                                        let tmpScanConfigs = [...scanConfigs].filter(s => s.id !== scanConfig.id);
                                        setScanConfigs(tmpScanConfigs);
                                        toast("Removed scan config successfully!");
                                    }
                                }}
                                button={
                                    <Trash color="red" />
                                }
                            />
                            <Button onClick={() => setScanConfig(scanConfig)} variant="secondary" className="ml-2">Edit</Button>
                        </CardFooter>
                    </Card>
                ))}
            </div>
        </div>
    );
}
