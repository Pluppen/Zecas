import { Button } from "@/components/ui/button"
import {useState, type Dispatch, type SetStateAction} from "react";
import { zodResolver } from "@hookform/resolvers/zod"
import { useForm } from "react-hook-form"
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog"
import { z } from "zod"
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"
import {isValidCIDR, isValidDomain, isValidIP} from '@/lib/utils'
import { Input } from "@/components/ui/input"
import { toast } from "sonner";
import { activeProjectIdStore } from "@/lib/projectsStore"
import { useStore } from "@nanostores/react"
import { user } from "@/lib/userStore";
import { createProjectTarget, type Target } from "@/lib/api/targets";

const FormSchema = z.object({
  targetType: z.string(),
  ipAddress: z.string().optional().refine(value => {
    if (!value) return true; // Allow empty
    return isValidIP(value)
  }, {
    message: "Invalid IP address format (e.g., 192.168.1.1).",
  }),
  cidrRange: z.string().optional().refine(value => {
    if (!value) return true; // Allow empty
    return isValidCIDR(value);
  }, {
    message: "Invalid CIDR format (e.g., 192.168.1.0/24).",
  }),
  domain: z.string().optional().refine(value => {
    if (!value) return true; // Allow empty
    return isValidDomain(value);
  }, {
    message: "Invalid domain format (e.g., example.com).",
  }),
}).refine(data => {
  return data.ipAddress || data.cidrRange || data.domain;
}, {
  message: "At least one target (IP, CIDR or domain) must be specified.",
  path: ["ipAddress"], // Show error on the first field
});

interface CreateTargetDialogProps {
    setTargets: Dispatch<SetStateAction<Target[]>>
}

export default function CreateTargetDialog({setTargets}: CreateTargetDialogProps) {
  const [open, setOpen] = useState(false);
  const $activeProjectId = useStore(activeProjectIdStore);
  const $user = useStore(user);

  const form = useForm<z.infer<typeof FormSchema>>({
    resolver: zodResolver(FormSchema),
    defaultValues: {},
  })

  function onSubmit(data: z.infer<typeof FormSchema>) {
    if ($user?.access_token) {
      let targetValue = ""

      if (targetType == "ip" && data.ipAddress) {
        targetValue = data.ipAddress;
      } else if (targetType == "cidr" && data.cidrRange) {
        targetValue = data.cidrRange;
      } else if (targetType == "domain" && data.domain) {
        targetValue = data.domain;
      } else {
        throw new Error("Something went wrong");
      }

      createProjectTarget($activeProjectId ?? "", data.targetType, targetValue, $user.access_token).then(res => {
        if ("err" in res) {
          toast("Something went wrong.")
          return
        }
        setTargets((prev) => [...prev, res]);
        toast("Successfully added new target.")
        setOpen(false);
      })
    }
  }

  const targetType = form.watch("targetType");
    return (
      <Dialog open={open} onOpenChange={() => setOpen(!open)}>
        <DialogTrigger>
            <Button>Add new target</Button>
        </DialogTrigger>
        <DialogContent>
            <DialogHeader>
                <DialogTitle>Add new target</DialogTitle>
                <DialogDescription>
                    Add a new target
                </DialogDescription>
            </DialogHeader>
            <Form {...form}>
                <form onSubmit={form.handleSubmit(onSubmit)}>
                    <FormField
                        control={form.control}
                        name="targetType"
                        render={({ field }) => (
                            <FormItem>
                            <FormLabel>Type</FormLabel>
                                <Select onValueChange={field.onChange} value={field.value}>
                                <FormControl>
                                    <SelectTrigger className="w-[180px]">
                                    <SelectValue placeholder="Target Type" />
                                    </SelectTrigger>
                                </FormControl>
                                <SelectContent>
                                    <SelectItem value="ip">IP</SelectItem>
                                    <SelectItem value="cidr">CIDR</SelectItem>
                                    <SelectItem value="domain">Domain</SelectItem>
                                </SelectContent>
                                </Select>
                            <FormMessage />
                            </FormItem>
                        )}
                    />

                    <FormField
                        control={form.control}
                        name="ipAddress"
                        render={({ field }) => (
                            <FormItem className={`ease-in-out delay-150 duration-300 ${targetType != "ip" ? "absolute opacity-0 translate-y-2" : "relative"}`}>
                            <FormLabel>IP Addresses {}</FormLabel>
                            <FormControl>
                                <Input
                                placeholder="192.168.1.1"
                                {...field} 
                                />
                            </FormControl>
                            <FormMessage />
                            </FormItem>
                        )}
                    />

                    <FormField
                        control={form.control}
                        name="cidrRange"
                        render={({ field }) => (
                            <FormItem className={`ease-in-out delay-150 duration-300 ${targetType != "cidr" ? "absolute opacity-0 translate-y-2" : "relative"}`}>
                            <FormLabel>CIDR</FormLabel>
                            <FormControl>
                                <Input
                                placeholder="8.8.8.8/31"
                                {...field} 
                                />
                            </FormControl>
                            <FormMessage />
                            </FormItem>
                        )}
                    />

                    <FormField
                        control={form.control}
                        name="domain"
                        render={({ field }) => (
                            <FormItem className={`ease-in-out delay-150 duration-300 ${targetType != "domain" ? "absolute opacity-0 translate-y-2" : "relative"}`}>
                            <FormLabel>Domain</FormLabel>
                            <FormControl>
                                <Input
                                placeholder="www.google.com"
                                {...field} 
                                />
                            </FormControl>
                            <FormMessage />
                            </FormItem>
                        )}
                    />
                    <DialogFooter>
                        <Button type="submit">Submit</Button>
                    </DialogFooter>
                </form>
            </Form>
        </DialogContent>
    </Dialog>
    )
}