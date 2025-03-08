"use client"
import {useEffect, useState} from "react";
import { zodResolver } from "@hookform/resolvers/zod"
import { useForm } from "react-hook-form"
import { z } from "zod"

import { API_URL } from "@/config"
import { Button } from "@/components/ui/button"
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form"
import { Input } from "@/components/ui/input"
import {isValidCIDR, isValidDomain, isValidIP} from '@/lib/utils'
import { toast } from "sonner";

import { activeProjectIdStore } from "@/lib/projectsStore"
import { useStore } from "@nanostores/react"

import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"

import { createProjectTarget } from "@/lib/api/targets";
import { getProjectTargets } from "@/lib/api/projects";
import { user } from "@/lib/userStore";

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
  // Ensure at least one target type is specified
  return data.ipAddress || data.cidrRange || data.domain;
}, {
  message: "At least one target (IP, CIDR or domain) must be specified.",
  path: ["ipAddress"], // Show error on the first field
});

export default function InputForm() {
  const $activeProjectId = useStore(activeProjectIdStore);
  const $user = useStore(user);

  const form = useForm<z.infer<typeof FormSchema>>({
    resolver: zodResolver(FormSchema),
    defaultValues: {},
  })

  const [targets, setTargets] = useState([]);

  useEffect(() => {
    if ($activeProjectId) {
      getProjectTargets($activeProjectId, $user.access_token).then(targets => {
        setTargets(targets);
      });
    }
  }, [$activeProjectId])


  function onSubmit(data: z.infer<typeof FormSchema>) {
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
      let targetsCopy = targets.slice();
      targetsCopy.push(res);
      setTargets(targetsCopy);
      toast("Successfully added new target.")
    })
  }

  const targetType = form.watch("targetType");

  return (
    <div>
      <h2>Current Targets</h2>
      <div className="flex gap-2 mt-4 mb-8 w-2/3">
        {targets.map((target) => (
          <div key={"target-id-"+target.id} className="text-sm bg-primary text-primary-foreground shadow-xs hover:bg-primary/90 rounded-md font-medium px-4 py-2 w-max-content">
            <span>{target.value}</span>
          </div>
        ))}
      </div>
      <hr />
      <h2 className="mb-4 mt-4">Add new target</h2>
      <Form {...form}>
        <form onSubmit={form.handleSubmit(onSubmit)} className="w-2/3 space-y-6">
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
          <Button type="submit">Submit</Button>
        </form>
      </Form>
    </div>
  )
}
