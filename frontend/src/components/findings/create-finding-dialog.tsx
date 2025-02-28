import {useEffect, useState} from "react";
import { zodResolver } from "@hookform/resolvers/zod"
import { z } from "zod"
import { useForm, useWatch } from "react-hook-form"

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
import { Textarea } from "@/components/ui/textarea";

import { activeProjectIdStore } from "@/lib/projectsStore";
import { useStore } from "@nanostores/react";

import { getScanConfigs, startNewScan } from "@/lib/scans";
import { getProjectTargets } from "@/lib/targets";
import { toast } from "sonner";
import { Input } from "@/components/ui/input";
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form"

const SEVERITIES = [
    "unknown",
    "informational",
    "low",
    "medium",
    "high",
    "critical"
]

const SeverityEnumZod = z.enum(["unknown", "informational", "low", "medium", "high", "critical"])

const NewFindingSchema = z.object({
    title: z.string(),
    description: z.string().optional(),
    finding_type: z.string(),
    targets: z.array(z.object({
        value: z.string(),
        label: z.string()
    })).min(1, {
        message: "At least 1 target must be selected"
    }),
    severity: SeverityEnumZod
})

export default function CreateFindingDialog () {
  const $activeProjectId = useStore(activeProjectIdStore);

  const [targets, setTargets] = useState([]);
  const [selectedTargets, setSelectedTargets] = useState([]);
  const [selectedScanConfig, setSelectedScanConfig] = useState("");
  const [selectedSeverity, setSelectedSeverity] = useState("unknown");

  const form = useForm<z.infer<typeof NewFindingSchema>>({
    resolver: zodResolver(NewFindingSchema),
    defaultValues: {
      description: "",
      targets: []
    },
  })

  useEffect(() => {
    if ($activeProjectId) {
      getProjectTargets($activeProjectId).then(targetsData => {
        setTargets(targetsData);
      });
    }
  }, [$activeProjectId])

  const onSubmit = (data: z.infer<typeof NewFindingSchema>) => {
    console.log(data);
    if($activeProjectId) {
    }
  }
  
  const formTargets = form.watch("targets");
  console.log(formTargets);

  return (
    <Dialog>
      <DialogTrigger asChild className="mt-4">
        <Button>Add new finding</Button>
      </DialogTrigger>
      <DialogContent className="sm:max-w-[525px]">
        <DialogHeader>
          <DialogTitle>Add new finding</DialogTitle>
          <DialogDescription>
            Add a new finding manually
          </DialogDescription>
        </DialogHeader>
        <Form {...form}>
            <form onSubmit={(event) => {
                event.preventDefault();
                form.handleSubmit(onSubmit)()
                console.log("ABC");
                }}>
                <FormField
                    control={form.control}
                    name="title"
                    render={({ field }) => (
                        <FormItem>
                            <FormLabel>
                                Title
                            </FormLabel>
                            <FormControl>
                                <Input placeholder="Title of finding" {...field} />
                            </FormControl>
                        </FormItem>
                    )}
                />
                <FormField
                    control={form.control}
                    name="description"
                    render={({ field }) => (
                        <FormItem>
                            <FormLabel>
                                Description
                            </FormLabel>
                            <Textarea placeholder="Description of finding..." {...field} />
                        </FormItem>
                    )}
                />
                <FormField
                    control={form.control}
                    name="finding_type"
                    render={({ field }) => (
                        <FormItem>
                            <FormLabel>
                                Finding Type
                            </FormLabel>
                            <FormControl>
                                <Input placeholder="Type of finding, e.g Vulnerable Service" {...field} />
                            </FormControl>
                        </FormItem>
                    )}
                />
                <FormField
                    control={form.control}
                    name="targets"
                    render={({ field }) => (
                        <FormItem>
                            <FormLabel>Targets</FormLabel>
                            <MultiSelect
                                data={targets.map(t => ({value: t.id, label: t.value}))}
                                placeholder="Select target(s)"
                                selected={field.value}
                                setSelected={field.onChange}
                            />
                        </FormItem>
                    )}
                />
                <FormField>
                    <Label htmlFor="name">
                        Severity
                    </Label>
                    <Select onValueChange={(value) => setSelectedSeverity(value)} value={selectedSeverity}>
                        <SelectTrigger className="w-[180px] capitalize">
                            <SelectValue placeholder="Target Type" />
                        </SelectTrigger>
                        <SelectContent>
                            {SEVERITIES.map(severity => (
                                <SelectItem
                                    key={"severity-item-"+severity}
                                    value={severity}
                                    className="capitalize"
                                >
                                    {severity}
                                </SelectItem>
                            ))}
                        </SelectContent>
                    </Select>
                </div>
                <div>
                    {/* TODO File upload for details / raw data*/}
                </div>
                <DialogFooter>
                    <Button type="submit">Start scan</Button>
                </DialogFooter>
            </form>
        </Form>
      </DialogContent>
    </Dialog>
  )
}
