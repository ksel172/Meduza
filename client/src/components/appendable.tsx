"use client"

import * as React from "react"
import { Check, ChevronsUpDown, X } from 'lucide-react'

import { cn } from "@/lib/utils"
import { Button } from "@/components/ui/button"
import {
  Command,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
} from "@/components/ui/command"
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover"
import { Label } from "@/components/ui/label"
import { Badge } from "@/components/ui/badge"

// const initialFrameworks = [
//   { value: "next.js", label: "Next.js" },
//   { value: "sveltekit", label: "SvelteKit" },
//   { value: "nuxt.js", label: "Nuxt.js" },
// ]

export function MultiSelectPopover({
    initialFrameworks,
    selectPlaceholder,
    addPlaceholder,
  }: {
    initialFrameworks: { value: string, label: string }[]
    selectPlaceholder: string;
    addPlaceholder: string;
  }) {
  const [open, setOpen] = React.useState(false)
  const [selectedValues, setSelectedValues] = React.useState<string[]>([])
  const [frameworks, setFrameworks] = React.useState(initialFrameworks)
  const [inputValue, setInputValue] = React.useState("")

  const handleInputChange = (value: string) => {
    setInputValue(value)
  }

  const handleKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === 'Enter' && inputValue.trim() !== '') {
      e.preventDefault()
      const newFramework = {
        value: inputValue.toLowerCase().replace(/\s+/g, '-'),
        label: inputValue.trim(),
      }
      setFrameworks([...frameworks, newFramework])
      setInputValue('')
    }
  }

  const handleSelect = (currentValue: string) => {
    setSelectedValues((prev) =>
      prev.includes(currentValue)
        ? prev.filter((value) => value !== currentValue)
        : [...prev, currentValue]
    )
  }

  const removeValue = (valueToRemove: string) => {
    setSelectedValues((prev) => prev.filter((value) => value !== valueToRemove))
  }

  const selectedBadges = selectedValues.slice(0, 1).map((value) => (
    <Badge key={value} variant="secondary" className="mr-1">
      {frameworks.find((framework) => framework.value === value)?.label}
      <button
        className="ml-1 ring-offset-background rounded-full outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2"
        onKeyDown={(e) => {
          if (e.key === "Enter") {
            removeValue(value)
          }
        }}
        onMouseDown={(e) => {
          e.preventDefault()
          e.stopPropagation()
        }}
        onClick={() => removeValue(value)}
      >
        <X className="h-3 w-3 text-muted-foreground hover:text-foreground" />
      </button>
    </Badge>
  ))

  const additionalCount = selectedValues.length - 1
  const overflowBadge = additionalCount > 0 && (
    <Badge variant="secondary">+{additionalCount} more</Badge>
  )

  return (
    <Popover open={open} onOpenChange={setOpen}>
    <PopoverTrigger asChild>
        <Button
        variant="outline"
        role="combobox"
        aria-expanded={open}
        className="w-full justify-between"
        >
        <div className="flex flex-wrap items-center gap-1">
            {selectedValues.length > 0 ? (
            <>
                {selectedBadges}
                {overflowBadge}
            </>
            ) : (
            // "Select frameworks..."
            selectPlaceholder
            )}
        </div>
        <ChevronsUpDown className="ml-2 h-4 w-4 shrink-0 opacity-50" />
        </Button>
    </PopoverTrigger>
    <PopoverContent className="w-full p-0">
        <Command>
        <CommandInput 
            placeholder={addPlaceholder}
            value={inputValue}
            onValueChange={handleInputChange}
            onKeyDown={handleKeyDown}
        />
        <CommandList>
            <CommandEmpty>Hit enter to add.</CommandEmpty>
            <CommandGroup>
            {frameworks.map((framework) => (
                <CommandItem
                key={framework.value}
                value={framework.value}
                onSelect={() => handleSelect(framework.value)}
                >
                <Check
                    className={cn(
                    "mr-2 h-4 w-4",
                    selectedValues.includes(framework.value) ? "opacity-100" : "opacity-0"
                    )}
                />
                {framework.label}
                </CommandItem>
            ))}
            </CommandGroup>
        </CommandList>
        </Command>
    </PopoverContent>
    </Popover>
  )
}