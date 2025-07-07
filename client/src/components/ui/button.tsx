"use client"

import * as React from "react"
import { Slot } from "@radix-ui/react-slot"
import { cva, type VariantProps } from "class-variance-authority"

import { cn } from "@/lib/utils"

const buttonVariants = cva(
  "inline-flex items-center justify-center gap-2 whitespace-nowrap rounded-md text-sm font-medium ring-offset-background transition-all duration-200 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50 [&_svg]:pointer-events-none [&_svg]:size-4 [&_svg]:shrink-0 relative overflow-hidden ripple-container",
  {
    variants: {
      variant: {
        default: "bg-primary text-primary-foreground hover:bg-primary/90 hover:shadow-lg hover:-translate-y-0.5 active:translate-y-0 active:shadow-md transition-all duration-200 ease-out",
        destructive: "bg-gradient-destructive text-destructive-foreground hover:shadow-lg hover:shadow-red-500/25 hover:-translate-y-0.5 active:translate-y-0 transition-all duration-200",
        outline: "border border-input bg-background hover:bg-accent hover:text-accent-foreground hover:border-ring hover:shadow-md hover:-translate-y-0.5 active:translate-y-0 transition-all duration-200",
        secondary: "bg-secondary text-secondary-foreground hover:bg-secondary/80 hover:shadow-md hover:-translate-y-0.5 active:translate-y-0 transition-all duration-200",
        ghost: "hover:bg-accent hover:text-accent-foreground hover:shadow-sm hover:-translate-y-0.5 active:translate-y-0 transition-all duration-200",
        link: "text-primary underline-offset-4 hover:underline hover:text-primary/80 transition-colors duration-200",
        gradient: "bg-gradient-primary text-white hover:shadow-glow hover:-translate-y-0.5 active:translate-y-0 transition-all duration-200 shadow-md",
        glassmorphism: "glass backdrop-blur-sm text-foreground hover:shadow-lg hover:-translate-y-0.5 active:translate-y-0 transition-all duration-200",
        premium: "bg-gradient-secondary text-white hover:shadow-glow-strong hover:-translate-y-1 hover:scale-105 active:translate-y-0 active:scale-100 transition-all duration-300 ease-out shadow-lg animate-pulse-soft"
      },
      size: {
        default: "h-10 px-4 py-2",
        sm: "h-9 rounded-md px-3",
        lg: "h-11 rounded-md px-8",
        icon: "h-10 w-10",
        xl: "h-12 rounded-lg px-10 text-base"
      },
      animation: {
        none: "",
        bounce: "hover:animate-bounce-soft",
        glow: "animate-glow",
        float: "animate-float",
        wiggle: "hover:animate-wiggle",
        heartbeat: "animate-heart-beat"
      }
    },
    defaultVariants: {
      variant: "default",
      size: "default",
      animation: "none"
    },
  }
)

export interface ButtonProps
  extends React.ButtonHTMLAttributes<HTMLButtonElement>,
    VariantProps<typeof buttonVariants> {
  asChild?: boolean
  ripple?: boolean
}

const Button = React.forwardRef<HTMLButtonElement, ButtonProps>(
  ({ className, variant, size, animation, asChild = false, ripple = true, onClick, ...props }, ref) => {
    const Comp = asChild ? Slot : "button"
    
    const handleClick = (e: React.MouseEvent<HTMLButtonElement>) => {
      if (ripple && !asChild) {
        const button = e.currentTarget
        const rect = button.getBoundingClientRect()
        const size = Math.max(rect.width, rect.height)
        const x = e.clientX - rect.left - size / 2
        const y = e.clientY - rect.top - size / 2
        
        const rippleElement = document.createElement('span')
        rippleElement.className = 'ripple-effect'
        rippleElement.style.width = rippleElement.style.height = size + 'px'
        rippleElement.style.left = x + 'px'
        rippleElement.style.top = y + 'px'
        
        button.appendChild(rippleElement)
        
        setTimeout(() => {
          rippleElement.remove()
        }, 600)
      }
      
      if (onClick) {
        onClick(e)
      }
    }
    
    return (
      <Comp
        className={cn(buttonVariants({ variant, size, animation, className }))}
        ref={ref}
        onClick={handleClick}
        {...props}
      />
    )
  }
)
Button.displayName = "Button"

export { Button, buttonVariants }
