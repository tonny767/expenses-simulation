<script setup>
import { Primitive } from "reka-ui";
import { cn } from "@/lib/utils";
import { buttonVariants } from "./index";

const props = defineProps({
  variant: { type: String, required: false, default: "default" },
  size: { type: String, required: false, default: "md" },
  class: { type: [String, Array, Object], required: false, default: "" },
  asChild: { type: Boolean, required: false, default: false },
  as: { type: String, required: false, default: "button" },
});
</script>

<template>
  <Primitive
    data-slot="button"
    :as="props.as"
    :as-child="props.asChild"
    :class="cn(buttonVariants({ variant: props.variant, size: props.size }), props.class)"
  >
    <slot />
  </Primitive>
</template>
