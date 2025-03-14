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

import { getScanConfigs, startNewScan } from "@/lib/api/scans";
import { getProjectTargets } from "@/lib/api/projects";
import { toast } from "sonner";
import { user } from "@/lib/userStore";

export default function CreateScanConfigDialog () {
  const $activeProjectId = useStore(activeProjectIdStore);
  const $user = useStore(user);

  const [targets, setTargets] = useState([]);
  const [scanConfigs, setScanConfigs] = useState([]);
  const [selectedTargets, setSelectedTargets] = useState([]);
  const [selectedScanConfig, setSelectedScanConfig] = useState("");

  useEffect(() => {
    getScanConfigs($user.access_token).then(result => {
        setScanConfigs(result);
    });
  }, [$user])

  useEffect(() => {
    if ($activeProjectId) {
      getProjectTargets($activeProjectId, $user.access_token).then(targetsData => {
        setTargets(targetsData);
      });
    }
  }, [$activeProjectId, $user])

  const handleSubmit = () => {
    const targetList = selectedTargets.map(t => t.value);

    if($activeProjectId) {
        startNewScan($activeProjectId, selectedScanConfig, targetList, $user.access_token).then(res => {
            if ("error" in res) {
                toast(res.error);
                return;
            }
            toast(res.message);
        })
    }
  }

  // TODO: What data is needed for this?
  // TODO: What params are needed
  // TODO: Possible to run dynamic form validation/parameters json validation based on scanner type?

  return (
    <Dialog>
      <DialogTrigger asChild>
        <Button>Create scan config</Button>
      </DialogTrigger>
      <DialogContent className="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>Create new scan config</DialogTitle>
          <DialogDescription>
            Create a new scan configuration that can be used for scans
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
                            value={scanConfig.id}
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
