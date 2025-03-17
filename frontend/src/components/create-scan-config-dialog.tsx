import {useEffect, useState} from "react";
import { zodResolver } from "@hookform/resolvers/zod"
import { z } from "zod"
import { useForm } from "react-hook-form"

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

import {type ScannerType, ScanConfigSchema, createScanConfig} from "@/lib/api/scans";

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

const SCANNER_TYPES: ScannerType[] = [
    "nmap",
    "dns",
    "subdomain",
    "nuclei",
    "httpx"
]

export default function CreateFindingDialog ({setScanConfigs}: {setScanConfigs: SetStateAction}) {
  const $user = useStore(user);
  const [open, setOpen] = useState(false);

  const form = useForm<z.infer<typeof ScanConfigSchema>>({
    resolver: zodResolver(ScanConfigSchema),
    defaultValues: {
      name: "",
    },
  })

  const onSubmit = (data: z.infer<typeof ScanConfigSchema>) => {
    if ($user.access_token) {
        createScanConfig({
            name: data.name,
            active: true,
            scanner_type: data.scanner_type,
            parameters: data.parameters
        }, $user.access_token).then((result) => {
            if ("error" in result) {
                toast(result.error);
                return
            }
            setScanConfigs((prev) => [...prev, result]);
            toast("Added new finding successfully!");
        });
        setOpen(false);
    }
  }
  
  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild className="mt-4">
        <Button>Create scan config</Button>
      </DialogTrigger>
      <DialogContent className="sm:max-w-[525px]">
        <DialogHeader>
          <DialogTitle>Create scan config</DialogTitle>
          <DialogDescription>
            Create a new scan config
          </DialogDescription>
        </DialogHeader>
        <Form {...form}>
            <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
                <FormField
                    control={form.control}
                    name="name"
                    render={({ field }) => (
                        <FormItem>
                            <FormLabel>
                                Name
                            </FormLabel>
                            <FormControl>
                                <Input placeholder="Name of scan config" {...field} />
                            </FormControl>
                            <FormMessage />
                        </FormItem>
                    )}
                />
                <FormField
                    control={form.control}
                    name="scanner_type"
                    render={({field}) => (
                        <FormItem>
                            <Label htmlFor="scanner_type">
                                Scanner Type
                            </Label>
                            <Select onValueChange={field.onChange} value={field.value}>
                                <FormControl>
                                    <SelectTrigger className="w-[180px] capitalize">
                                        <SelectValue placeholder="Target Type" />
                                    </SelectTrigger>
                                </FormControl>
                                <FormMessage />
                                <SelectContent>
                                    {SCANNER_TYPES.map(scannerType => (
                                        <SelectItem
                                            key={"scannerType-item-"+scannerType}
                                            value={scannerType}
                                            className="capitalize"
                                        >
                                            {scannerType}
                                        </SelectItem>
                                    ))}
                                </SelectContent>
                            </Select>
                        </FormItem>
                    )}
                />
                <div>
                    {/* TODO JSON editor for parameters */}
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
