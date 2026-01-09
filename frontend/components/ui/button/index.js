
export { default as Button } from "./Button.vue";

import { cva } from "class-variance-authority";

export const buttonVariants = cva(
  "inline-flex items-center justify-center rounded-md font-medium transition-colors focus:outline-none focus:ring-1 focus:ring-offset-2 cursor-pointer",
  {
    variants: {
      variant: {
        default: "bg-blue-500 text-white hover:bg-blue-600",
        destructive: "bg-red-500 text-white hover:bg-red-600",
        green: "bg-green-500 text-white hover:bg-green-600",
        outline: "border border-gray-300 text-gray-700 hover:bg-gray-100",
        ghost: "bg-transparent hover:bg-gray-100",
      },
      size: {
        sm: "px-2 py-1 text-sm",
        md: "px-4 py-2 text-base",
        lg: "px-6 py-3 text-lg",
      },
    },
    defaultVariants: {
      variant: "default",
      size: "md",
    },
  }
);
