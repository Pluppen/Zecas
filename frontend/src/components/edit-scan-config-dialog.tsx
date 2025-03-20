import {useState, type Dispatch, type SetStateAction} from "react";
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
} from "@/components/ui/dialog"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"
import { Label } from "@/components/ui/label"
import { Textarea } from "@/components/ui/textarea";

import { useStore } from "@nanostores/react";

import {type ScanConfig, type ScannerType, ScanConfigSchema, ScannerTypeEnum, updateScanConfigById} from "@/lib/api/scans";

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

const SCANNER_TYPES: ScannerType[] = [
    "nmap",
    "dns",
    "subdomain",
    "nuclei",
    "httpx"
]

const scanConfigFormSchema = z.object({
  name: z.string(),
  scanner_type: ScannerTypeEnum,
  parameters: z.string()
})

interface EditScanConfigDialogProps {
    setScanConfigs: Dispatch<SetStateAction<ScanConfig[]>>
    scanConfigs: ScanConfig[]
    scanConfig?: ScanConfig
    open: boolean
    setOpen: Dispatch<SetStateAction<boolean>>
}

export default function EditScanConfigDialog ({setScanConfigs, scanConfigs, scanConfig, open, setOpen}: EditScanConfigDialogProps) {
  const $user = useStore(user);

  const form = useForm<z.infer<typeof scanConfigFormSchema>>({
    resolver: zodResolver(scanConfigFormSchema),
    defaultValues: {
      name: scanConfig?.name,
      scanner_type: scanConfig?.scanner_type,
      parameters: JSON.stringify(scanConfig?.parameters, null, 2)
    },
  })

  const onSubmit = (data: z.infer<typeof scanConfigFormSchema>) => {
    if ($user.access_token) {
        let parameters;
        try {
          // TODO: Make sure this parsing gets done safely
          parameters = JSON.parse(data.parameters);
        } catch (e) {
          const errorOpt: ErrorOption = {
            "message": "Invalid JSON"
          }
          form.setError("parameters", errorOpt)
          return
        }

        const scanConfigBody = {
          ...data,
          active: true,
          parameters,
          id: scanConfig?.id
        }
        const dataResult = ScanConfigSchema.safeParse(scanConfigBody);

        if (!dataResult.success)  {
          const errorOpt: ErrorOption = {
            "message": ""
          }
          dataResult.error.errors.forEach(e => {
            errorOpt["message"] += ` ${e.message}`
          });
          form.setError("parameters", errorOpt)
          return
        }

        updateScanConfigById(scanConfigBody, $user.access_token).then((result) => {
            if ("error" in result) {
                toast(result.error);
                return
            }
            let tmpScanConfigs = [...scanConfigs].map(s => s.id == result.id ? result : s);
            setScanConfigs(tmpScanConfigs);
            toast("Edited finding successfully!");
        });
        setOpen(false);
    }
  }
  
  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogContent className="sm:max-w-[525px]">
        <DialogHeader>
          <DialogTitle>Edit scan config</DialogTitle>
          <DialogDescription>
            Edit scan config
          </DialogDescription>
        </DialogHeader>
        <Form {...form}>
            <form onSubmit={form.handleSubmit(onSubmit, (data) => {
              console.log("Invalido Data");
              console.log(data);
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
                                        <SelectValue placeholder="Scanner Type" />
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
                <FormField
                    control={form.control}
                    name="parameters"
                    render={({field}) => (
                        <FormItem>
                            <Label htmlFor="scanner_type">
                              Parameters
                            </Label>
                            <FormControl>
                              <Textarea {...field} />
                            </FormControl>
                            <FormMessage />
                        </FormItem>
                    )}
                />
                <div>
                    {/* TODO JSON editor for parameters */}
                </div>
                <DialogFooter>
                    <Button type="submit">
                        Edit Scan Config
                    </Button>
                </DialogFooter>
            </form>
        </Form>
      </DialogContent>
    </Dialog>
  )
}
