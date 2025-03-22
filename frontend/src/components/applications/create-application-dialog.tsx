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
import { Label } from "@/components/ui/label"
import MultiSelect from "@/components/multi-select"
import { Textarea } from "@/components/ui/textarea";

import { activeProjectIdStore } from "@/lib/projectsStore";
import { useStore } from "@nanostores/react";

import { getProjectServices, getProjectTargets } from "@/lib/api/projects";
import { createApplication } from "@/lib/api/applications";

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
import {ApplicationSchema} from "@/lib/api/applications"

const NewApplicationSchema = ApplicationSchema.extend({
    targets: z.array(z.object({
        value: z.string().uuid(),
        label: z.string()
    })),
    services: z.array(z.object({
        value: z.string().uuid(),
        label: z.string()
    }))
}).omit({host_target: true, service_id: true, project_id: true})

type NewApplicationType = z.infer<typeof NewApplicationSchema>

export default function CreateApplicationDialog ({setApplications}: {setApplications: Dispatch<SetStateAction<any>>}) {
  const $activeProjectId = useStore(activeProjectIdStore);
  const $user = useStore(user);

  const [open, setOpen] = useState(false);
  const [targets, setTargets] = useState([]);
  const [services, setServices] = useState([]);

  const form = useForm<NewApplicationType>({
    resolver: zodResolver(NewApplicationSchema),
    defaultValues: {
        targets: [],
        services: []
    },
  })

  useEffect(() => {
    if($activeProjectId) {
        getProjectTargets($activeProjectId, $user.access_token).then(result => {
            setTargets(result);
        })

        getProjectServices($activeProjectId, $user.access_token).then(result => {
            setServices(result);
        })
    }
  }, [$activeProjectId, $user])

  const onSubmit = (data: NewApplicationType) => {
    if($activeProjectId) {
        let rawData;
        try {
            if (data.metadata) {
                rawData = JSON.parse(data.metadata);
            }
        } catch (e) {
            const errorOpt: ErrorOption = {
                "message": "Invalid JSON"
            }
            form.setError("metadata", errorOpt);
            return
        }

        data.targets.forEach(target => {
            data.services.forEach(service => {
                const newApplication = {
                    project_id: $activeProjectId,
                    name: data.name,
                    type: data.type,
                    version: data.version,
                    description: data.description,
                    url: data.url,
                    host_target: target.value,
                    service_id: service.value
                }
                createApplication(newApplication, $user.access_token).then((result) => {
                    if ("error" in result) {
                        toast(result.error);
                        return
                    }
                    setApplications((prev) => [...prev, result]);
                    toast(`Added new application for ${target.label}`);
                });
            })
        })
        setOpen(false);
    }
  }
  
  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild className="mt-4">
        <Button>Add application</Button>
      </DialogTrigger>
      <DialogContent className="sm:max-w-[525px]">
        <DialogHeader>
          <DialogTitle>Add application to target(s)</DialogTitle>
          <DialogDescription>
            Add a new application manually to target(s)
          </DialogDescription>
        </DialogHeader>
        <Form {...form}>
            <form onSubmit={form.handleSubmit(onSubmit, (aaaa) => {
                console.error(aaaa);
            })} className="space-y-4">
                <FormField
                    control={form.control}
                    name="name"
                    render={({ field }) => (
                        <FormItem>
                            <FormLabel>
                                Name
                            </FormLabel>
                            <FormControl>
                                <Input placeholder="Name of application" {...field} />
                            </FormControl>
                            <FormMessage />
                        </FormItem>
                    )}
                />
                <FormField
                    control={form.control}
                    name="type"
                    render={({ field }) => (
                        <FormItem>
                            <FormLabel>
                                Type
                            </FormLabel>
                            <FormControl>
                                <Input placeholder="Application type (ex: CMS)..." {...field} />
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
                                <Input placeholder="Application version..." {...field} />
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
                            <Textarea placeholder="Description of application..." {...field} />
                            <FormMessage />
                        </FormItem>
                    )}
                />
                <FormField
                    control={form.control}
                    name="url"
                    render={({ field }) => (
                        <FormItem>
                            <FormLabel>
                                URL
                            </FormLabel>
                            <FormControl>
                                <Input placeholder="URL of application..." {...field} />
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
                    name="services"
                    render={({ field }) => (
                        <FormItem>
                            <FormLabel>Services</FormLabel>
                            <FormControl>
                                <MultiSelect
                                    data={services.map(t => ({value: t.id, label: `${t.banner ?? t.port} (${t.target_id})`}))}
                                    placeholder="Select service(s)"
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
                    name="metadata"
                    render={({field}) => (
                        <FormItem>
                            <Label>
                              Metadata (JSON)
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
                        Add Application 
                    </Button>
                </DialogFooter>
            </form>
        </Form>
      </DialogContent>
    </Dialog>
  )
}
