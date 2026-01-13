import type { Variants } from 'motion/react';

export const blurInVariants: Variants = {
  hidden: { opacity: 0, y: 5, filter: 'blur(10px)' },
  visible: {
    opacity: 1,
    y: 0,
    filter: 'blur(0px)',
    transition: {
      duration: 0.4,
      ease: [0.25, 0.1, 0.25, 1],
    },
  },
};
