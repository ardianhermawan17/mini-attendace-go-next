import React from 'react';
import { Text as ChakraText, TextProps as ChakraTextProps, Heading } from '@chakra-ui/react';

export interface TextProps extends ChakraTextProps {
  children: React.ReactNode;
  variant?: 'body' | 'caption' | 'label';
  size?: 'xs' | 'sm' | 'md' | 'lg' | 'xl' | '2xl';
}

export interface HeadingProps extends ChakraTextProps {
  children: React.ReactNode;
  level?: 1 | 2 | 3 | 4 | 5 | 6;
  size?: 'xs' | 'sm' | 'md' | 'lg' | 'xl' | '2xl';
}

/**
 * Text UI Component
 * Presentational component that only accepts props
 */
export const TextUI: React.FC<TextProps> = ({
  children,
  variant = 'body',
  size = 'md',
  ...props
}) => {
  const variantStyles = {
    body: { fontSize: size },
    caption: { fontSize: 'xs', color: 'gray.600' },
    label: { fontSize: 'sm', fontWeight: 'medium' },
  };

  return (
    <ChakraText {...variantStyles[variant]} {...props}>
      {children}
    </ChakraText>
  );
};

/**
 * Heading UI Component
 * Presentational component that only accepts props
 */
export const HeadingUI: React.FC<HeadingProps> = ({
  children,
  level = 1,
  size = 'lg',
  ...props
}) => {
  const sizeMap = {
    1: '2xl',
    2: 'xl',
    3: 'lg',
    4: 'md',
    5: 'sm',
    6: 'xs',
  };

  return (
    <Heading as={`h${level}`} size={sizeMap[level]} {...props}>
      {children}
    </Heading>
  );
};

export default TextUI;
