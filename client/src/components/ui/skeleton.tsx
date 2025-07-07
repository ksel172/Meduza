import { cn } from "@/lib/utils"

function Skeleton({
  className,
  variant = "default",
  ...props
}: React.HTMLAttributes<HTMLDivElement> & {
  variant?: "default" | "shimmer" | "wave" | "glow"
}) {
  const variants = {
    default: "animate-pulse rounded-md bg-muted",
    shimmer: "animate-shimmer rounded-md bg-gradient-to-r from-muted via-muted-foreground/20 to-muted bg-[length:200%_100%]",
    wave: "animate-pulse rounded-md bg-muted animate-float",
    glow: "animate-glow rounded-md bg-muted shadow-glow"
  }

  return (
    <div
      className={cn(variants[variant], className)}
      {...props}
    />
  )
}

export { Skeleton }
