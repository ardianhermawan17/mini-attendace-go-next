import type { Meta } from '@storybook/react';
import { TextUI, HeadingUI } from './text-ui';
import { VStack } from '@chakra-ui/react';

const meta: Meta<typeof TextUI> = {
  title: 'Atoms/Text',
  component: TextUI,
  parameters: {
    layout: 'centered',
  },
  tags: ['autodocs'],
};

export default meta;

export const Default = {
  args: {
    children: 'This is a text component',
    size: 'md',
  },
};

export const Small = {
  args: {
    children: 'This is small text',
    size: 'sm',
  },
};

export const Large = {
  args: {
    children: 'This is large text',
    size: 'lg',
  },
};

export const Label = {
  args: {
    children: 'Form Label',
    variant: 'label',
  },
};

export const Caption = {
  args: {
    children: 'This is a caption or helper text',
    variant: 'caption',
  },
};

export const Heading1 = {
  render: () => <HeadingUI level={1}>Heading Level 1</HeadingUI>,
};

export const Heading2 = {
  render: () => <HeadingUI level={2}>Heading Level 2</HeadingUI>,
};

export const Heading3 = {
  render: () => <HeadingUI level={3}>Heading Level 3</HeadingUI>,
};

export const AllTextVariants = {
  render: () => (
    <VStack align="start" spacing={4}>
      <HeadingUI level={1}>Heading 1</HeadingUI>
      <HeadingUI level={2}>Heading 2</HeadingUI>
      <HeadingUI level={3}>Heading 3</HeadingUI>
      <TextUI size="lg">Large Text</TextUI>
      <TextUI size="md">Medium Text (Default)</TextUI>
      <TextUI size="sm">Small Text</TextUI>
      <TextUI variant="label">Label Text</TextUI>
      <TextUI variant="caption">Caption or Helper Text</TextUI>
    </VStack>
  ),
};
