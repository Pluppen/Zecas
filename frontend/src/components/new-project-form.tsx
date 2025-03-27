"use client"

import { zodResolver } from "@hookform/resolvers/zod"
import { useForm } from "react-hook-form"
import { z } from "zod"

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
import { Textarea } from "./ui/textarea"
import {isValidCIDR, isValidDomain, isValidIP} from '@/lib/utils'

import { activeProjectIdStore } from "@/lib/projectsStore"
import { user } from "@/lib/userStore"
import { createNewProject } from "@/lib/api/projects"
import { useStore } from "@nanostores/react"

const FormSchema = z.object({
  name: z.string().min(2, {
    message: "Project must be at least 1 characters.",
  }),
  description: z.string().optional(),
  ipAddresses: z.string().optional().refine(value => {
    if (!value) return true; // Allow empty
    return value.split('\n')
      .map(ip => ip.trim())
      .filter(ip => ip !== '')
      .every(isValidIP);
  }, {
    message: "Invalid IP address format. Use one IP per line (e.g., 192.168.1.1).",
  }),
  cidrRanges: z.string().optional().refine(value => {
    if (!value) return true; // Allow empty
    return value.split('\n')
      .map(cidr => cidr.trim())
      .filter(cidr => cidr !== '')
      .every(isValidCIDR);
  }, {
    message: "Invalid CIDR format. Use one CIDR range per line (e.g., 192.168.1.0/24).",
  }),
  domains: z.string().optional().refine(value => {
    if (!value) return true; // Allow empty
    return value.split('\n')
      .map(domain => domain.trim())
      .filter(domain => domain !== '')
      .every(isValidDomain);
  }, {
    message: "Invalid domain format. Use one domain per line (e.g., example.com).",
  }),
  keywords: z.string().optional(),
}).refine(data => {
  // Ensure at least one target type is specified
  return data.ipAddresses || data.cidrRanges || data.domains;
}, {
  message: "At least one target (IP, CIDR or domain) must be specified.",
  path: ["ipAddresses"], // Show error on the first field
});

function transformTargets(targets: string) {
  console.log(targets);
  if (targets.split("\n").length < 1) {
    return []
  }

  return targets.split("\n").map(t => t.trim()).filter(t => t != "")
}

export default function InputForm() {
  const $user = useStore(user);

  const form = useForm<z.infer<typeof FormSchema>>({
    resolver: zodResolver(FormSchema),
    defaultValues: {
      name: "",
      description: ""
    },
  })

  function onSubmit(data: z.infer<typeof FormSchema>) {
    if ($user?.access_token) {
      let body = {
        name: data.name,
        description: data.description,
        ip_ranges: data.ipAddresses ? transformTargets(data.ipAddresses) : [],
        cidr_ranges: data.cidrRanges ? transformTargets(data.cidrRanges) : [],
        domains: data.domains ? transformTargets(data.domains) : []
      }

      createNewProject(JSON.stringify(body), $user.access_token).then((result) => {
        console.log(result);
        activeProjectIdStore.set(result.id)
        window.location.href="/"
      });
    }
  }

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="w-2/3 space-y-6">
        <FormField
          control={form.control}
          name="name"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Project Name</FormLabel>
              <FormControl>
                <Input placeholder="shadcn" {...field} />
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
              <FormLabel>Project Description</FormLabel>
              <FormControl>
                <Textarea placeholder="Description of project..." {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />

        <FormField
        control={form.control}
        name="ipAddresses"
        render={({ field }) => (
            <FormItem>
            <FormLabel>IP Addresses</FormLabel>
            <FormControl>
                <Textarea 
                placeholder="192.168.1.1
192.168.1.2" 
                className="h-32 resize-y font-mono text-sm" 
                {...field} 
                />
            </FormControl>
            <FormDescription>
                Enter one IP address per line
            </FormDescription>
            <FormMessage />
            </FormItem>
        )}
        />
        
        <FormField
        control={form.control}
        name="cidrRanges"
        render={({ field }) => (
            <FormItem>
            <FormLabel>CIDR Ranges</FormLabel>
            <FormControl>
                <Textarea 
                placeholder="192.168.1.0/24
10.0.0.0/8" 
                className="h-32 resize-y font-mono text-sm" 
                {...field} 
                />
            </FormControl>
            <FormDescription>
                Enter one CIDR range per line
            </FormDescription>
            <FormMessage />
            </FormItem>
        )}
        />
    
        <FormField
        control={form.control}
        name="domains"
        render={({ field }) => (
            <FormItem>
            <FormLabel>Domains</FormLabel>
            <FormControl>
                <Textarea 
                placeholder="example.com
subdomain.example.com" 
                className="h-32 resize-y font-mono text-sm" 
                {...field} 
                />
            </FormControl>
            <FormDescription>
                Enter one domain per line
            </FormDescription>
            <FormMessage />
            </FormItem>
        )}
        />
        
        <FormField
        control={form.control}
        name="keywords"
        render={({ field }) => (
            <FormItem>
            <FormLabel>Keywords (Optional)</FormLabel>
            <FormControl>
                <Textarea 
                placeholder="keyword1
keyword2" 
                className="h-32 resize-y font-mono text-sm" 
                {...field} 
                />
            </FormControl>
            <FormDescription>
                Enter one keyword per line
            </FormDescription>
            <FormMessage />
            </FormItem>
        )}
        />
        <Button type="submit">Submit</Button>
      </form>
    </Form>
  )
}
