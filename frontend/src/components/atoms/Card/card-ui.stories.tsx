import type { Meta } from '@storybook/react';
import { CardUI } from './card-ui';
import { Box, Text } from '@chakra-ui/react';

const meta: Meta<typeof CardUI> = {
  title: 'Atoms/Card',
  component: CardUI,
  parameters: {
    layout: 'centered',
  },
  tags: ['autodocs'],
};

export default meta;

export const Default = {
  render: () => (
    <CardUI>
      <Box p={4}>
        <Text fontSize="lg" fontWeight="bold" mb={2}>
          Card Title
        </Text>
        <Text>This is a card component with some content inside.</Text>
      </Box>
    </CardUI>
  ),
};

export const Elevated = {
  render: () => (
    <CardUI variant="elevated">
      <Box p={4}>
        <Text fontSize="lg" fontWeight="bold" mb={2}>
          Elevated Card
        </Text>
        <Text>This card has an elevated variant with a shadow.</Text>
      </Box>
    </CardUI>
  ),
};

export const Outlined = {
  render: () => (
    <CardUI variant="outline">
      <Box p={4}>
        <Text fontSize="lg" fontWeight="bold" mb={2}>
          Outlined Card
        </Text>
        <Text>This card has an outline variant.</Text>
      </Box>
    </CardUI>
  ),
};

export const Filled = {
  render: () => (
    <CardUI variant="filled">
      <Box p={4}>
        <Text fontSize="lg" fontWeight="bold" mb={2}>
          Filled Card
        </Text>
        <Text>This card has a filled variant with background color.</Text>
      </Box>
    </CardUI>
  ),
};

export const Clickable = {
  render: () => (
    <CardUI cursor="pointer" _hover={{ boxShadow: 'lg' }}>
      <Box p={4}>
        <Text fontSize="lg" fontWeight="bold" mb={2}>
          Clickable Card
        </Text>
        <Text>This card can be clicked and has a hover effect.</Text>
      </Box>
    </CardUI>
  ),
};
