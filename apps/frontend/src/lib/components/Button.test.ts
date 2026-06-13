import { describe, it, expect, vi } from 'vitest';
import { render, fireEvent } from '@testing-library/svelte';
import Button from './Button.svelte';
import { createRawSnippet } from 'svelte';

describe('Button Component', () => {
  it('should render button text correctly', () => {
    // Svelte 5 snippet mock for testing children
    const childrenSnippet = createRawSnippet(() => ({
      render: () => '<span data-testid="child">Test Button</span>',
      setup: () => {}
    }));

    const { getByTestId } = render(Button, {
      props: {
        children: childrenSnippet
      }
    });

    const child = getByTestId('child');
    expect(child).toBeInTheDocument();
    expect(child).toHaveTextContent('Test Button');
  });

  it('should fire onclick handler when clicked', async () => {
    const childrenSnippet = createRawSnippet(() => ({
      render: () => '<span>Click Me</span>',
      setup: () => {}
    }));

    const handleClick = vi.fn();
    const { getByRole } = render(Button, {
      props: {
        onclick: handleClick,
        children: childrenSnippet
      }
    });

    const button = getByRole('button');
    await fireEvent.click(button);
    expect(handleClick).toHaveBeenCalledTimes(1);
  });

  it('should apply disabled attribute and disabled styling', () => {
    const childrenSnippet = createRawSnippet(() => ({
      render: () => '<span>Disabled Button</span>',
      setup: () => {}
    }));

    const { getByRole } = render(Button, {
      props: {
        disabled: true,
        children: childrenSnippet
      }
    });

    const button = getByRole('button');
    expect(button).toBeDisabled();
    expect(button.className).toContain('opacity-50');
    expect(button.className).toContain('shadow-none');
  });

  it('should apply custom background color prop', () => {
    const childrenSnippet = createRawSnippet(() => ({
      render: () => '<span>Colored Button</span>',
      setup: () => {}
    }));

    const { getByRole } = render(Button, {
      props: {
        bgColor: 'bg-neubrutalism-pink',
        children: childrenSnippet
      }
    });

    const button = getByRole('button');
    expect(button.className).toContain('bg-neubrutalism-pink');
  });
});
