import type { Meta, StoryObj } from '@storybook/react';
import { ButtonUI } from './button-ui';

const meta = {
  title: 'Atoms/Button',
  component: ButtonUI,
  parameters: {
    layout: 'centered',
  },
  tags: ['autodocs'],
} satisfies Meta<typeof ButtonUI>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Solid: Story = {
  args: {
    children: 'Click Me',
    variant: 'solid',
    colorScheme: 'blue',
    size: 'md',
  },
};

export const Outline: Story = {
  args: {
    children: 'Outline Button',
    variant: 'outline',
    colorScheme: 'blue',
    size: 'md',
  },
};

export const Ghost: Story = {
  args: {
    children: 'Ghost Button',
    variant: 'ghost',
    colorScheme: 'blue',
    size: 'md',
  },
};

export const Link: Story = {
  args: {
    children: 'Link Button',
    variant: 'link',
    colorScheme: 'blue',
    size: 'md',
  },
};

export const Small: Story = {
  args: {
    children: 'Small',
    variant: 'solid',
    size: 'sm',
    colorScheme: 'blue',
  },
};

export const Large: Story = {
  args: {
    children: 'Large',
    variant: 'solid',
    size: 'lg',
    colorScheme: 'blue',
  },
};

export const Loading: Story = {
  args: {
    children: 'Loading',
    variant: 'solid',
    isLoading: true,
    colorScheme: 'blue',
  },
};

export const Disabled: Story = {
  args: {
    children: 'Disabled',
    variant: 'solid',
    isDisabled: true,
    colorScheme: 'blue',
  },
};

export const Success: Story = {
  args: {
    children: 'Success',
    variant: 'solid',
    colorScheme: 'green',
  },
};

export const Danger: Story = {
  args: {
    children: 'Danger',
    variant: 'solid',
    colorScheme: 'red',
  },
};

export const Warning: Story = {
  args: {
    children: 'Warning',
    variant: 'solid',
    colorScheme: 'yellow',
  },
};
