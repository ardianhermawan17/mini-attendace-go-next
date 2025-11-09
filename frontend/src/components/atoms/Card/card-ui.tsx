import React from 'react';
import { Box, BoxProps } from '@chakra-ui/react';

export interface CardProps extends BoxProps {
  children: React.ReactNode;
  variant?: 'elevated' | 'outline' | 'filled';
  padding?: string | number;
}

/**
 * Card UI Component
 * Presentational component that only accepts props
 */
export const CardUI: React.FC<CardProps> = ({
  children,
  variant = 'elevated',
  padding = 6,
  ...props
}) => {
  const variantStyles = {
    elevated: {
      boxShadow: 'md',
      borderRadius: 'lg',
      bg: 'white',
    },
    outline: {
      border: '1px solid',
      borderColor: 'gray.200',
      borderRadius: 'lg',
      bg: 'white',
    },
    filled: {
      bg: 'gray.50',
      borderRadius: 'lg',
    },
  };

  return (
    <Box p={padding} {...variantStyles[variant]} {...props}>
      {children}
    </Box>
  );
};

export default CardUI;
