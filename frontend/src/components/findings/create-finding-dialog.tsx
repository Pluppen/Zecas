import {useEffect, useState, type Dispatch, type SetStateAction} from "react";
import { zodResolver } from "@hookform/resolvers/zod"
import { z } from "zod"
import { useForm } from "react-hook-form"

import { Button } from "@/components/ui/button"
import {
  Dialog,
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

import { getProjectTargets } from "@/lib/api/projects";
import { createFinding, type Finding } from "@/lib/api/findings";
import { type Target } from "@/lib/api/targets";

import { toast } from "sonner";
import { Input } from "@/components/ui/input";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form"
import { user } from "@/lib/userStore";

const SEVERITIES = [
    "unknown",
    "info",
    "low",
    "medium",
    "high",
    "critical"
]

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
    severity: z.string()
})

export default function CreateFindingDialog ({setFindings}: {setFindings: Dispatch<SetStateAction<Finding[]>>}) {
  const $activeProjectId = useStore(activeProjectIdStore);
  const $user = useStore(user);

  const [open, setOpen] = useState(false);
  const [targets, setTargets] = useState<Target[]>([]);

  const form = useForm<z.infer<typeof NewFindingSchema>>({
    resolver: zodResolver(NewFindingSchema),
    defaultValues: {
      description: "",
      targets: []
    },
  })

  useEffect(() => {
    if ($activeProjectId && $user?.access_token) {
      getProjectTargets($activeProjectId, $user.access_token).then(targetsData => {
        setTargets(targetsData);
      });
    }
  }, [$activeProjectId, $user?.access_token])

  const onSubmit = (data: z.infer<typeof NewFindingSchema>) => {
    if($activeProjectId && $user?.access_token) {
        data.targets.forEach(target => {
            createFinding({
                title: data.title,
                description: data.description,
                finding_type: data.finding_type,
                target_id: target.value,
                severity: data.severity,
                manual: true
            }, $user.access_token ?? "").then((result) => {
                if ("error" in result) {
                    toast(result.error);
                    return
                }
                setFindings((prev) => [...prev, result]);
                toast("Added new finding successfully!");
            });
        })
        setOpen(false);
    }
  }
  
  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild className="mt-4">
        <Button>Add finding</Button>
      </DialogTrigger>
      <DialogContent className="sm:max-w-[525px]">
        <DialogHeader>
          <DialogTitle>Add finding</DialogTitle>
          <DialogDescription>
            Add a new finding manually
          </DialogDescription>
        </DialogHeader>
        <Form {...form}>
            <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
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
                            <FormMessage />
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
                            <FormMessage />
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
                            <FormMessage />
                        </FormItem>
                    )}
                />
                <FormField
                    control={form.control}
                    name="targets"
                    render={({ field }) => (
                        <FormItem>
                            <FormLabel>Targets</FormLabel>
                            <FormControl>
                                <MultiSelect
                                    data={targets.map(t => ({value: t.id, label: t.value}))}
                                    placeholder="Select target(s)"
                                    selected={field.value}
                                    setSelected={field.onChange}
                                />
                            </FormControl>
                            <FormMessage />
                        </FormItem>
                    )}
                />
                <FormField
                    control={form.control}
                    name="severity"
                    render={({field}) => (
                        <FormItem>
                            <Label htmlFor="name">
                                Severity
                            </Label>
                            <Select onValueChange={field.onChange} value={field.value}>
                                <FormControl>
                                    <SelectTrigger className="w-[180px] capitalize">
                                        <SelectValue placeholder="Target Type" />
                                    </SelectTrigger>
                                </FormControl>
                                <FormMessage />
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
                        </FormItem>
                    )}
                />
                <div>
                    {/* TODO File upload for details / raw data*/}
                </div>
                <DialogFooter>
                    <Button type="submit">
                        Add Finding 
                    </Button>
                </DialogFooter>
            </form>
        </Form>
      </DialogContent>
    </Dialog>
  )
}
