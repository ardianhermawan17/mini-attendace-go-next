import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import { ChakraProvider } from '@chakra-ui/react';
import { ButtonUI } from './button-ui';

const renderWithChakra = (component: React.ReactElement) => {
  return render(<ChakraProvider>{component}</ChakraProvider>);
};

describe('ButtonUI', () => {
  it('renders button with children', () => {
    renderWithChakra(<ButtonUI>Click me</ButtonUI>);
    expect(screen.getByText('Click me')).toBeInTheDocument();
  });

  it('calls onClick handler when clicked', () => {
    const handleClick = jest.fn();
    renderWithChakra(<ButtonUI onClick={handleClick}>Click me</ButtonUI>);

    fireEvent.click(screen.getByText('Click me'));
    expect(handleClick).toHaveBeenCalledTimes(1);
  });

  it('disables button when isDisabled is true', () => {
    renderWithChakra(<ButtonUI isDisabled>Click me</ButtonUI>);
    expect(screen.getByText('Click me')).toBeDisabled();
  });

  it('shows loading state', () => {
    renderWithChakra(<ButtonUI isLoading>Click me</ButtonUI>);
    const button = screen.getByText('Click me');
    expect(button).toHaveAttribute('data-loading');
  });

  it('applies correct variant', () => {
    const { container } = renderWithChakra(<ButtonUI variant="outline">Click me</ButtonUI>);
    const button = container.querySelector('button');
    expect(button).toHaveClass('chakra-button');
  });

  it('applies correct size', () => {
    renderWithChakra(<ButtonUI size="lg">Click me</ButtonUI>);
    const button = screen.getByText('Click me');
    expect(button).toBeInTheDocument();
  });
});
