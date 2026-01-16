import type { Variants } from 'motion/react';

export const blurInVariants = (delay: number = 0, translateY: number = 5): Variants => {
  return {
    hidden: { opacity: 0, y: translateY, filter: 'blur(10px)' },
    visible: {
      opacity: 1,
      y: 0,
      filter: 'blur(0px)',
      transition: {
        duration: 0.4,
        ease: [0.25, 0.1, 0.25, 1],
        delay,
      },
    },
  }
};
