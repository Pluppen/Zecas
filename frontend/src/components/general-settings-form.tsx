"use client"

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
import { Textarea } from "./ui/textarea"

import { activeProjectStore, activeProjectIdStore } from "@/lib/projectsStore"
import { useStore } from "@nanostores/react"

const FormSchema = z.object({
  name: z.string().min(2, {
    message: "Project must be at least 1 characters.",
  }),
  description: z.string().optional()
})

export default function InputForm() {
  const $activeProject = useStore(activeProjectStore);
  const $activeProjectId = useStore(activeProjectIdStore);
  const form = useForm<z.infer<typeof FormSchema>>({
    resolver: zodResolver(FormSchema),
    defaultValues: {
      name: $activeProject?.name,
      description: $activeProject?.description,
    },
  })


    function onSubmit(data: z.infer<typeof FormSchema>) {
      // TODO: Add func to edit general settings
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
        <Button type="submit">Save</Button>
      </form>
    </Form>
  )
}
