import * as React from "react"
import { Button as ShadcnButton, buttonVariants } from "~/components/ui/button"
import { type VariantProps } from "class-variance-authority"
import { cn } from "~/lib/utils"

export interface ButtonProps
  extends React.ComponentProps<"button">,
  VariantProps<typeof buttonVariants> {
  isLoading?: boolean
  asChild?: boolean
}

const Button = React.forwardRef<HTMLButtonElement, ButtonProps>(
  ({ className, isLoading, children, disabled, ...props }, ref) => {
    return (
      <ShadcnButton
        className={cn(className, "cursor-pointer")}
        disabled={isLoading || disabled}
        ref={ref}
        {...props}
      >
        {isLoading ? (
          <>
            <span className="me-2 h-4 w-4 animate-spin rounded-full border-2 border-current border-t-transparent" />
            {children}
          </>
        ) : (
          children
        )}
      </ShadcnButton>
    )
  }
)
Button.displayName = "Button"

export { Button }