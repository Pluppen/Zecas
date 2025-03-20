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
import { type Service, ServiceSchema } from "@/lib/api/services";

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

const NewServiceSchema = ServiceSchema.extend(z.object({
    targets: z.array(z.string())
}))

export default function CreateServiceDialog ({setServices}: {setServices: Dispatch<SetStateAction<any>>}) {
  const $activeProjectId = useStore(activeProjectIdStore);
  const $user = useStore(user);

  const [open, setOpen] = useState(false);
  const [targets, setTargets] = useState([]);

  const form = useForm<z.infer<typeof NewServiceSchema>>({
    resolver: zodResolver(NewServiceSchema),
  })

  useEffect(() => {
    if ($activeProjectId) {
      getProjectTargets($activeProjectId, $user.access_token).then(targetsData => {
        setTargets(targetsData);
      });
    }
  }, [$activeProjectId, $user])

  const onSubmit = (data: z.infer<typeof NewServiceSchema>) => {
    if($activeProjectId) {
        //data.targets.forEach(target => {
        //    createFinding({
        //        title: data.title,
        //        description: data.description,
        //        finding_type: data.finding_type,
        //        target_id: target.value,
        //        severity: data.severity,
        //        manual: true
        //    }, $user.access_token).then((result) => {
        //        if ("error" in result) {
        //            toast(result.error);
        //            return
        //        }
        //        setFindings((prev) => [...prev, result]);
        //        toast("Added new finding successfully!");
        //    });
        //})
        setOpen(false);
    }
  }
  
  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild className="mt-4">
        <Button>Add service</Button>
      </DialogTrigger>
      <DialogContent className="sm:max-w-[525px]">
        <DialogHeader>
          <DialogTitle>Add service</DialogTitle>
          <DialogDescription>
            Add a new service manually
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
                    name="port"
                    render={({ field }) => (
                        <FormItem>
                            <FormLabel>
                                Port
                            </FormLabel>
                            <FormControl>
                                <Input type="number" min={1} max={65535} placeholder="Ex. 3389" {...field} />
                            </FormControl>
                            <FormMessage />
                        </FormItem>
                    )}
                />
                <FormField
                    control={form.control}
                    name="protocol"
                    render={({ field }) => (
                        <FormItem>
                            <FormLabel>
                                Protocol
                            </FormLabel>
                            <FormControl>
                                <Input placeholder="Protocol on endpoint, example TCP, UDP" {...field} />
                            </FormControl>
                            <FormMessage />
                        </FormItem>
                    )}
                />
                <FormField
                    control={form.control}
                    name="version"
                    render={({ field }) => (
                        <FormItem>
                            <FormLabel>
                                Version
                            </FormLabel>
                            <FormControl>
                                <Input placeholder="Verison of service" {...field} />
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
                <div>
                    {/* TODO File upload for details / raw data*/}
                </div>
                <DialogFooter>
                    <Button type="submit">
                        Add Service
                    </Button>
                </DialogFooter>
            </form>
        </Form>
      </DialogContent>
    </Dialog>
  )
}
