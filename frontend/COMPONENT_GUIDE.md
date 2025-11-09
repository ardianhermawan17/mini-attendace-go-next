# Component Guide

This guide explains how to use and extend components in the Mini Attendance Frontend.

## Table of Contents

1. [Atomic Design Pattern](#atomic-design-pattern)
2. [Component Structure](#component-structure)
3. [Atoms](#atoms)
4. [Molecules](#molecules)
5. [Organisms](#organisms)
6. [Creating New Components](#creating-new-components)
7. [Best Practices](#best-practices)

## Atomic Design Pattern

Components are organized following the Atomic Design methodology:

```
Atoms → Molecules → Organisms → Templates → Pages
```

### Hierarchy

- **Atoms**: Basic building blocks (Button, Input, Card)
- **Molecules**: Combinations of atoms (LoginForm, CheckInCard)
- **Organisms**: Complex components (Header, Sidebar, Layout)
- **Templates**: Page layouts (not yet implemented)
- **Pages**: Full pages (in `src/app/`)

## Component Structure

Each component follows a consistent structure:

### File Organization

```
components/
└── atoms/
    └── Button/
        ├── button-ui.tsx      # Presentational component
        ├── button-ui.test.tsx # Tests
        └── index.ts           # Export
```

### Component Files

1. **`*-ui.tsx`**: Presentational component
   - Only accepts props
   - No data fetching
   - No Redux/hooks
   - Pure UI logic

2. **`use-*.ts`**: Logic hook (in features/)
   - State management
   - Data fetching
   - Business logic
   - No UI code

3. **`index.tsx`**: Container component (in features/)
   - Combines UI and logic
   - Passes props to UI component
   - Handles side effects

## Atoms

Basic building blocks that are reusable across the application.

### Button

```typescript
import { ButtonUI } from '@/components/atoms/Button/button-ui';

// Basic usage
<ButtonUI>Click me</ButtonUI>

// With props
<ButtonUI
  colorScheme="blue"
  size="lg"
  isLoading={isLoading}
  onClick={handleClick}
>
  Submit
</ButtonUI>

// Variants
<ButtonUI variant="solid">Solid</ButtonUI>
<ButtonUI variant="outline">Outline</ButtonUI>
<ButtonUI variant="ghost">Ghost</ButtonUI>
<ButtonUI variant="link">Link</ButtonUI>
```

### Input

```typescript
import { InputUI } from '@/components/atoms/Input/input-ui';

// Basic usage
<InputUI
  label="Username"
  placeholder="Enter username"
  value={username}
  onChange={handleChange}
/>

// With error
<InputUI
  label="Email"
  type="email"
  error="Invalid email format"
  value={email}
  onChange={handleChange}
/>

// Different types
<InputUI type="text" />
<InputUI type="password" />
<InputUI type="email" />
<InputUI type="number" />
```

### Card

```typescript
import { CardUI } from '@/components/atoms/Card/card-ui';

// Basic usage
<CardUI>
  <Text>Card content</Text>
</CardUI>

// With variant
<CardUI variant="elevated">Elevated card</CardUI>
<CardUI variant="outline">Outline card</CardUI>
<CardUI variant="filled">Filled card</CardUI>

// With custom padding
<CardUI padding={8}>
  <Text>Custom padding</Text>
</CardUI>
```

### Text

```typescript
import { TextUI, HeadingUI } from '@/components/atoms/Text/text-ui';

// Text variants
<TextUI variant="body">Body text</TextUI>
<TextUI variant="caption">Caption text</TextUI>
<TextUI variant="label">Label text</TextUI>

// Text sizes
<TextUI size="xs">Extra small</TextUI>
<TextUI size="sm">Small</TextUI>
<TextUI size="md">Medium</TextUI>
<TextUI size="lg">Large</TextUI>

// Headings
<HeadingUI level={1}>Heading 1</HeadingUI>
<HeadingUI level={2}>Heading 2</HeadingUI>
<HeadingUI level={3}>Heading 3</HeadingUI>
```

## Molecules

Combinations of atoms that form more complex components.

### LoginForm

```typescript
import { LoginForm } from '@/features/auth';

// Usage in page
export default function LoginPage() {
  return (
    <Center minH="100vh">
      <LoginForm />
    </Center>
  );
}

// Component structure
// features/auth/
// ├── use-login.ts (hook with logic)
// ├── index.tsx (container)
// └── molecules/LoginForm/login-form-ui.tsx (UI)
```

### CheckInCard

```typescript
import { CheckInCard } from '@/features/attendance';

// Usage in page
export default function DashboardPage() {
  return (
    <Container>
      <CheckInCard />
    </Container>
  );
}

// Component structure
// features/attendance/
// ├── use-check-in.ts (check-in logic)
// ├── use-check-out.ts (check-out logic)
// ├── index.tsx (container)
// └── molecules/CheckInCard/check-in-card-ui.tsx (UI)
```

## Organisms

Complex components that combine multiple molecules.

### Layout (To be implemented)

```typescript
// Future implementation
import { Layout } from '@/components/organisms/Layout';

export default function DashboardPage() {
  return (
    <Layout>
      <CheckInCard />
      <AttendanceTable />
    </Layout>
  );
}
```

## Creating New Components

### Step 1: Create Atom Component

Create a new atom in `src/components/atoms/`:

```typescript
// src/components/atoms/Badge/badge-ui.tsx
import React from 'react';
import { Badge as ChakraBadge, BadgeProps as ChakraBadgeProps } from '@chakra-ui/react';

export interface BadgeProps extends ChakraBadgeProps {
  children: React.ReactNode;
  variant?: 'solid' | 'subtle' | 'outline';
}

export const BadgeUI: React.FC<BadgeProps> = ({
  children,
  variant = 'subtle',
  ...props
}) => {
  return (
    <ChakraBadge variant={variant} {...props}>
      {children}
    </ChakraBadge>
  );
};

export default BadgeUI;
```

### Step 2: Create Molecule Component

Create a new molecule in `src/components/molecules/`:

```typescript
// src/components/molecules/UserCard/user-card-ui.tsx
import React from 'react';
import { VStack, HStack } from '@chakra-ui/react';
import { CardUI } from '@/components/atoms/Card/card-ui';
import { TextUI, HeadingUI } from '@/components/atoms/Text/text-ui';
import { User } from '@/types';

export interface UserCardUIProps {
  user: User;
  onEdit?: () => void;
  onDelete?: () => void;
}

export const UserCardUI: React.FC<UserCardUIProps> = ({
  user,
  onEdit,
  onDelete,
}) => {
  return (
    <CardUI>
      <VStack align="start" spacing={3}>
        <HeadingUI level={3}>{user.full_name}</HeadingUI>
        <TextUI variant="caption">{user.email}</TextUI>
        <TextUI variant="label">{user.department}</TextUI>
        <HStack spacing={2}>
          {onEdit && <ButtonUI onClick={onEdit}>Edit</ButtonUI>}
          {onDelete && <ButtonUI onClick={onDelete}>Delete</ButtonUI>}
        </HStack>
      </VStack>
    </CardUI>
  );
};

export default UserCardUI;
```

### Step 3: Create Feature Hook

Create a hook in `src/features/`:

```typescript
// src/features/user/use-user-card.ts
import { useCallback } from 'react';
import { useRouter } from 'next/navigation';
import { User } from '@/types';

export const useUserCard = (user: User) => {
  const router = useRouter();

  const handleEdit = useCallback(() => {
    router.push(`/users/${user.id}/edit`);
  }, [user.id, router]);

  const handleDelete = useCallback(async () => {
    if (confirm('Are you sure?')) {
      // Delete logic
    }
  }, []);

  return {
    handleEdit,
    handleDelete,
  };
};
```

### Step 4: Create Container Component

Create a container in `src/features/`:

```typescript
// src/features/user/index.tsx
import React from 'react';
import { UserCardUI } from '@/components/molecules/UserCard/user-card-ui';
import { useUserCard } from './use-user-card';
import { User } from '@/types';

export interface UserCardProps {
  user: User;
}

export const UserCard: React.FC<UserCardProps> = ({ user }) => {
  const { handleEdit, handleDelete } = useUserCard(user);

  return (
    <UserCardUI
      user={user}
      onEdit={handleEdit}
      onDelete={handleDelete}
    />
  );
};

export default UserCard;
```

### Step 5: Use in Page

```typescript
// src/app/users/page.tsx
'use client';

import { UserCard } from '@/features/user';
import { useGetUsersQuery } from '@/store/api/usersApi';

export default function UsersPage() {
  const { data: users } = useGetUsersQuery();

  return (
    <Container>
      <VStack spacing={4}>
        {users?.map((user) => (
          <UserCard key={user.id} user={user} />
        ))}
      </VStack>
    </Container>
  );
}
```

## Best Practices

### 1. Separation of Concerns

```typescript
// ❌ Bad: UI and logic mixed
export const LoginForm = () => {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');

  const handleSubmit = async (e) => {
    // ... login logic
  };

  return (
    <form onSubmit={handleSubmit}>
      {/* UI */}
    </form>
  );
};

// ✅ Good: Separated concerns
// use-login.ts
export const useLogin = () => {
  const [formData, setFormData] = useState({ username: '', password: '' });
  const { login } = useAuth();
  // ... logic
  return { formData, handleSubmit };
};

// login-form-ui.tsx
export const LoginFormUI = (props) => (
  <form onSubmit={props.onSubmit}>
    {/* UI only */}
  </form>
);

// index.tsx
export const LoginForm = () => {
  const { formData, handleSubmit } = useLogin();
  return <LoginFormUI {...props} />;
};
```

### 2. Props Interface

```typescript
// ✅ Good: Clear, typed props
export interface ButtonProps extends ChakraButtonProps {
  children: React.ReactNode;
  variant?: 'solid' | 'outline' | 'ghost';
  size?: 'sm' | 'md' | 'lg';
  isLoading?: boolean;
}

export const ButtonUI: React.FC<ButtonProps> = (props) => {
  // ...
};
```

### 3. Reusability

```typescript
// ✅ Good: Reusable component
export const CardUI: React.FC<CardProps> = ({
  children,
  variant = 'elevated',
  padding = 6,
  ...props
}) => {
  // Can be used anywhere
};

// Usage
<CardUI variant="outline">Content</CardUI>
<CardUI variant="filled" padding={8}>Content</CardUI>
```

### 4. Type Safety

```typescript
// ✅ Good: Full TypeScript support
interface User {
  id: string;
  name: string;
  email: string;
}

interface UserCardProps {
  user: User;
  onEdit: (user: User) => void;
}

export const UserCard: React.FC<UserCardProps> = ({ user, onEdit }) => {
  // TypeScript ensures correct usage
};
```

### 5. Testing

```typescript
// ✅ Good: Testable components
describe('ButtonUI', () => {
  it('renders with children', () => {
    render(<ButtonUI>Click</ButtonUI>);
    expect(screen.getByText('Click')).toBeInTheDocument();
  });

  it('calls onClick handler', () => {
    const onClick = jest.fn();
    render(<ButtonUI onClick={onClick}>Click</ButtonUI>);
    fireEvent.click(screen.getByText('Click'));
    expect(onClick).toHaveBeenCalled();
  });
});
```

### 6. Documentation

```typescript
/**
 * Button UI Component
 * 
 * A reusable button component with multiple variants and sizes.
 * 
 * @example
 * <ButtonUI variant="solid" size="lg">
 *   Click me
 * </ButtonUI>
 */
export const ButtonUI: React.FC<ButtonProps> = (props) => {
  // ...
};
```

## Component Naming Conventions

- **UI Components**: `*-ui.tsx` (e.g., `button-ui.tsx`)
- **Hooks**: `use-*.ts` (e.g., `use-login.ts`)
- **Containers**: `index.tsx` in feature folder
- **Types**: `*.types.ts` (e.g., `auth.types.ts`)
- **Tests**: `*.test.tsx` (e.g., `button-ui.test.tsx`)

## Folder Structure

```
components/
├── atoms/
│   ├── Button/
│   │   ├── button-ui.tsx
│   │   ├── button-ui.test.tsx
│   │   └── index.ts
│   ├── Input/
│   ��── Card/
│   └── Text/
├── molecules/
│   ├── LoginForm/
│   │   ├── login-form-ui.tsx
│   │   └── index.ts
│   ├── CheckInCard/
│   └── AttendanceTable/
└── organisms/
    ├── Header/
    ├── Sidebar/
    └── Layout/

features/
├── auth/
│   ├── use-login.ts
│   ├── auth.types.ts
│   └── index.tsx
└── attendance/
    ├── use-check-in.ts
    ├── use-check-out.ts
    └── index.tsx
```

## Common Patterns

### Pattern 1: Form Component

```typescript
// Hook
export const useForm = (initialValues) => {
  const [values, setValues] = useState(initialValues);
  const [errors, setErrors] = useState({});

  const handleChange = (e) => {
    const { name, value } = e.target;
    setValues(prev => ({ ...prev, [name]: value }));
  };

  return { values, errors, handleChange };
};

// UI
export const FormUI = ({ values, errors, onSubmit, onInputChange }) => (
  <form onSubmit={onSubmit}>
    <InputUI value={values.name} onChange={onInputChange} error={errors.name} />
  </form>
);

// Container
export const Form = () => {
  const { values, errors, handleChange } = useForm({});
  return <FormUI values={values} errors={errors} onInputChange={handleChange} />;
};
```

### Pattern 2: List Component

```typescript
// Hook
export const useList = (items) => {
  const [filteredItems, setFilteredItems] = useState(items);
  const [search, setSearch] = useState('');

  useEffect(() => {
    setFilteredItems(items.filter(item => item.name.includes(search)));
  }, [search, items]);

  return { filteredItems, search, setSearch };
};

// UI
export const ListUI = ({ items, onSearch }) => (
  <VStack>
    <InputUI onChange={(e) => onSearch(e.target.value)} />
    {items.map(item => <ItemCard key={item.id} item={item} />)}
  </VStack>
);

// Container
export const List = ({ items }) => {
  const { filteredItems, search, setSearch } = useList(items);
  return <ListUI items={filteredItems} onSearch={setSearch} />;
};
```

## Resources

- [Chakra UI Documentation](https://chakra-ui.com)
- [React Documentation](https://react.dev)
- [Next.js Documentation](https://nextjs.org)
- [Atomic Design](https://atomicdesign.bradfrost.com)

## Support

For questions or issues with components, refer to:
- Component examples in `src/components/`
- Feature examples in `src/features/`
- Test examples in `*.test.tsx` files
