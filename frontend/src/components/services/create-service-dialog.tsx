import {useEffect, useState, type Dispatch, type SetStateAction} from "react";
import { zodResolver } from "@hookform/resolvers/zod"
import { z } from "zod"
import { useForm, type ErrorOption } from "react-hook-form"

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
import { createService } from "@/lib/api/services";

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
import { user } from "@/lib/userStore";
import {type Service, ServiceSchema} from "@/lib/api/services"

const NewServiceSchema = ServiceSchema.extend({
    targets: z.array(z.object({
        value: z.string().uuid(),
        label: z.string()
    }))
}).omit({target_id: true})
type NewServiceType = z.infer<typeof NewServiceSchema>

export default function CreateServiceDialog ({setServices}: {setServices: Dispatch<SetStateAction<any>>}) {
  const $activeProjectId = useStore(activeProjectIdStore);
  const $user = useStore(user);

  const [open, setOpen] = useState(false);
  const [targets, setTargets] = useState([]);

  const form = useForm<NewServiceType>({
    resolver: zodResolver(NewServiceSchema),
    defaultValues: {
        targets: []
    },
  })

  useEffect(() => {
    if($activeProjectId) {
        getProjectTargets($activeProjectId, $user.access_token).then(result => {
            console.log(result);
            setTargets(result);
        })
    }
  }, [$activeProjectId, $user])

  const onSubmit = (data: NewServiceType) => {
    if($activeProjectId) {
        let rawData;
        try {
            if (data.raw_info) {
                rawData = JSON.parse(data.raw_info);
            }
        } catch (e) {
            const errorOpt: ErrorOption = {
                "message": "Invalid JSON"
            }
            form.setError("raw_info", errorOpt);
            return
        }

        data.targets.forEach(target => {
            const newTarget = {
                title: data.title,
                description: data.description,
                port: data.port,
                protocol: data.protocol,
                service_name: data.service_name,
                version: data.version,
                banner: data.banner,
                raw_info: data.raw_info,
                target_id: target.value
            }
            createService(newTarget, $user.access_token).then((result) => {
                if ("error" in result) {
                    toast(result.error);
                    return
                }
                setServices((prev) => [...prev, result]);
                toast(`Added new service for ${target.label}`);
            });
        })
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
          <DialogTitle>Add service to target(s)</DialogTitle>
          <DialogDescription>
            Add a new service manually to target(s)
          </DialogDescription>
        </DialogHeader>
        <Form {...form}>
            <form onSubmit={form.handleSubmit(onSubmit, (aaaa) => {
                console.error(aaaa);
            })} className="space-y-4">
                <FormField
                    control={form.control}
                    name="title"
                    render={({ field }) => (
                        <FormItem>
                            <FormLabel>
                                Title
                            </FormLabel>
                            <FormControl>
                                <Input placeholder="Title of service" {...field} />
                            </FormControl>
                            <FormMessage />
                        </FormItem>
                    )}
                />
                <FormField
                    control={form.control}
                    name="service_name"
                    render={({ field }) => (
                        <FormItem>
                            <FormLabel>
                                Service name
                            </FormLabel>
                            <FormControl>
                                <Input placeholder="Service name..." {...field} />
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
                            <Textarea placeholder="Description of service..." {...field} />
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
                                <Input placeholder="Version of service..." {...field} />
                            </FormControl>
                            <FormMessage />
                        </FormItem>
                    )}
                />
                <FormField
                    control={form.control}
                    name="banner"
                    render={({ field }) => (
                        <FormItem>
                            <FormLabel>
                                Version
                            </FormLabel>
                            <FormControl>
                                <Input placeholder="Banner of service..." {...field} />
                            </FormControl>
                            <FormMessage />
                        </FormItem>
                    )}
                />
                <FormField
                    control={form.control}
                    name="port"
                    render={({ field }) => (
                        <FormItem>
                            <FormLabel>Port</FormLabel>
                            <FormControl>
                                <Input type="number" placeholder="Ex: 3389" {...field} min="1" max="65535" />
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
                                <Input placeholder="Protocol of service. Ex (TCP or UDP)" {...field} />
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
                    name="raw_info"
                    render={({field}) => (
                        <FormItem>
                            <Label htmlFor="scanner_type">
                              Raw Info (JSON)
                            </Label>
                            <FormControl>
                              <Textarea {...field} />
                            </FormControl>
                            <FormMessage />
                        </FormItem>
                    )}
                />
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
