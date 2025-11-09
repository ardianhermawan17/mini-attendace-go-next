import React from 'react';
import { Button as ChakraButton, ButtonProps as ChakraButtonProps } from '@chakra-ui/react';

export interface ButtonProps extends ChakraButtonProps {
  children: React.ReactNode;
  variant?: 'solid' | 'outline' | 'ghost' | 'link';
  size?: 'xs' | 'sm' | 'md' | 'lg';
  isLoading?: boolean;
  isDisabled?: boolean;
  onClick?: () => void;
}

/**
 * Button UI Component
 * Presentational component that only accepts props
 */
export const ButtonUI: React.FC<ButtonProps> = ({
  children,
  variant = 'solid',
  size = 'md',
  isLoading = false,
  isDisabled = false,
  onClick,
  ...props
}) => {
  return (
    <ChakraButton
      variant={variant}
      size={size}
      isLoading={isLoading}
      isDisabled={isDisabled}
      onClick={onClick}
      {...props}
    >
      {children}
    </ChakraButton>
  );
};

export default ButtonUI;
