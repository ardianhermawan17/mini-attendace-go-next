import React from 'react';
import {
  Input as ChakraInput,
  InputProps as ChakraInputProps,
  FormControl,
  FormLabel,
  FormErrorMessage,
} from '@chakra-ui/react';

export interface InputProps extends Omit<ChakraInputProps, 'size'> {
  label?: string;
  error?: string;
  size?: 'xs' | 'sm' | 'md' | 'lg';
  placeholder?: string;
  type?: string;
  value?: string;
  onChange?: (e: React.ChangeEvent<HTMLInputElement>) => void;
  disabled?: boolean;
}

/**
 * Input UI Component
 * Presentational component that only accepts props
 */
export const InputUI: React.FC<InputProps> = ({
  label,
  error,
  size = 'md',
  placeholder,
  type = 'text',
  value,
  onChange,
  disabled = false,
  ...props
}) => {
  return (
    <FormControl isInvalid={!!error}>
      {label && <FormLabel>{label}</FormLabel>}
      <ChakraInput
        type={type}
        placeholder={placeholder}
        size={size}
        value={value}
        onChange={onChange}
        isDisabled={disabled}
        {...props}
      />
      {error && <FormErrorMessage>{error}</FormErrorMessage>}
    </FormControl>
  );
};

export default InputUI;
