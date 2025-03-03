import {useState} from "react";

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

interface RemoveItemDialogProps {
    handleSubmit: any
    button: any
}

export function RemoveItemDialog({handleSubmit, button}: RemoveItemDialogProps) {
    const [open, setOpen] = useState(false);

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button variant="ghost">{button}</Button>
      </DialogTrigger>
      <DialogContent className="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>Remove item</DialogTitle>
          <DialogDescription>
            Are you sure you want to remove this item?
          </DialogDescription>
        </DialogHeader>
        <DialogFooter>
          <Button onClick={() => setOpen(false)} type="submit">Cancel</Button>
          <Button onClick={() => {
            handleSubmit();
            setOpen(false);
          }} variant="destructive" type="submit">Remove</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
