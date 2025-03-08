"use client"
import {useState, useEffect} from "react";

import { zodResolver } from "@hookform/resolvers/zod"
import { useForm } from "react-hook-form"
import { z } from "zod"

import { Button } from "@/components/ui/button"
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form"
import { Input } from "@/components/ui/input"
import { Textarea } from "./ui/textarea"

import { activeProjectIdStore, projects, activeProjectStore } from "@/lib/projectsStore"
import { user } from "@/lib/userStore"
import { useStore } from "@nanostores/react"
import { getProjectById, updateProject } from "@/lib/api/projects";

import { toast } from "sonner";

const FormSchema = z.object({
  name: z.string().min(2, {
    message: "Project must be at least 1 characters.",
  }),
  description: z.string().optional()
})

export default function InputForm() {
  const $user = useStore(user);
  const $activeProjectId = useStore(activeProjectIdStore);
  const $projects = useStore(projects);

  const form = useForm<z.infer<typeof FormSchema>>({
    resolver: zodResolver(FormSchema),
    defaultValues: async () => getProjectById($activeProjectId, $user.access_token).then(result => {
        if ("error" in result) {
          console.error(result.error);
          return {
            name: "",
            description: ""
          }
        }
        return result
      })
  })

  useEffect(() => {
    if ($activeProjectId) {
      getProjectById($activeProjectId, $user.access_token).then(result => {
        if ("error" in result) {
          console.error(result.error);
          form.reset({
            name: "",
            description: ""
          });
          return
        }
        form.reset({
          name: result.name,
          description: result.description
        });
      });
    }
  }, [$activeProjectId])

  function onSubmit(data: z.infer<typeof FormSchema>) {
    updateProject($activeProjectId, data.name, data.description, $user.access_token).then(result => {
      if ("error" in result) {
        console.error(result.error);
        toast("Something went wrong.")
        return
      }
      const newProjectsList = $projects.projects.map(p => {
        if (p.id == result.id) {
          return result
        }
        return p
      });
      projects.set({
        projects: newProjectsList,
      });
      activeProjectStore.set(result);
      toast("Successfully updated project settings")
    });
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
