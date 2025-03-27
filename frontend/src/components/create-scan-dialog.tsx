import {useEffect, useState} from "react";

import { Button } from "@/components/ui/button"
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"
import { Label } from "@/components/ui/label"
import MultiSelect from "@/components/multi-select"

import { activeProjectIdStore } from "@/lib/projectsStore";
import { useStore } from "@nanostores/react";

import { getScanConfigs, startNewScan, type ScanConfig } from "@/lib/api/scans";
import {type Target} from "@/lib/api/targets";
import { getProjectTargets } from "@/lib/api/projects";
import { toast } from "sonner";
import { user } from "@/lib/userStore";

export default function CreateScanDialog () {
  const $activeProjectId = useStore(activeProjectIdStore);
  const $user = useStore(user);

  const [targets, setTargets] = useState<Target[]>([]);
  const [scanConfigs, setScanConfigs] = useState<ScanConfig[]>([]);
  const [selectedTargets, setSelectedTargets] = useState<Record<"label" | "value", string>[]>([]);
  const [selectedScanConfig, setSelectedScanConfig] = useState("");

  useEffect(() => {
    if ($user?.access_token) {
      getScanConfigs($user.access_token).then(result => {
          setScanConfigs(result);
      });
    }
  }, [$user])

  useEffect(() => {
    if ($activeProjectId && $user?.access_token) {
      getProjectTargets($activeProjectId, $user.access_token).then(targetsData => {
        setTargets(targetsData);
      });
    }
  }, [$activeProjectId, $user])

  const handleSubmit = () => {
    const targetList = selectedTargets.map(t => t.value);

    if($activeProjectId && $user?.access_token) {
        startNewScan($activeProjectId, selectedScanConfig, targetList, $user.access_token).then(res => {
            if ("error" in res) {
                toast(res.error);
                return;
            }
            toast(res.message);
        })
    }
  }

  return (
    <Dialog>
      <DialogTrigger asChild>
        <Button>Start new scan</Button>
      </DialogTrigger>
      <DialogContent className="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>Start new scan</DialogTitle>
          <DialogDescription>
            Start a new scan towards chosen targets using scan config.
          </DialogDescription>
        </DialogHeader>
        <div className="">
            <Label htmlFor="name" className="text-right">
                Scan Config
            </Label>
            <Select onValueChange={(value) => setSelectedScanConfig(value)} value={selectedScanConfig}>
                <SelectTrigger className="w-[180px]">
                    <SelectValue placeholder="Target Type" />
                </SelectTrigger>
                <SelectContent>
                    {scanConfigs.map(scanConfig => (
                        <SelectItem
                            key={"scan-config-"+scanConfig.id}
                            value={scanConfig.id ?? ""}
                        >
                            {scanConfig.name}
                        </SelectItem>
                    ))}
                </SelectContent>
            </Select>
        </div>
        <div>
            <Label>
                Targets
            </Label>
            <MultiSelect
                data={targets.map(t => ({value: t.id, label: t.value}))}
                placeholder="Select target(s)"
                selected={selectedTargets}
                setSelected={setSelectedTargets}
            />
        </div>
        <DialogFooter>
            <DialogClose>
                <Button onClick={handleSubmit} type="submit">Start scan</Button>
            </DialogClose>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
