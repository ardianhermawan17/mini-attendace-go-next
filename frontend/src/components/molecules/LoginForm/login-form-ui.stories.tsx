import type { Meta } from '@storybook/react';
import { LoginFormUI } from './login-form-ui';
import { useState } from 'react';

const meta: Meta<typeof LoginFormUI> = {
  title: 'Molecules/LoginForm',
  component: LoginFormUI,
  parameters: {
    layout: 'centered',
  },
  tags: ['autodocs'],
};

export default meta;

function LoginFormWrapper(args: any) {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');

  return (
    <LoginFormUI
      {...args}
      username={username}
      password={password}
      onUsernameChange={(e) => setUsername(e.target.value)}
      onPasswordChange={(e) => setPassword(e.target.value)}
      onSubmit={(e) => {
        e.preventDefault();
        console.log('Form submitted:', { username, password });
      }}
    />
  );
}

export const Default = {
  render: (args: any) => <LoginFormWrapper {...args} />,
};

export const WithError = {
  render: (args: any) => <LoginFormWrapper {...args} error="Invalid email or password" />,
};

export const WithValidationErrors = {
  render: (args: any) => (
    <LoginFormWrapper
      {...args}
      usernameError="Username is required"
      passwordError="Password must be at least 8 characters"
    />
  ),
};

export const Loading = {
  render: (args: any) => <LoginFormWrapper {...args} isLoading={true} />,
};

function PrefilledFormWrapper() {
  const [username] = useState('john.doe@trustmedis.com');
  const [password] = useState('password123');

  return (
    <LoginFormUI
      username={username}
      password={password}
      onUsernameChange={() => {}}
      onPasswordChange={() => {}}
      onSubmit={(e) => {
        e.preventDefault();
      }}
    />
  );
}

export const WithPrefilledData = {
  render: () => <PrefilledFormWrapper />,
};

export const AllErrors = {
  render: (args: any) => (
    <LoginFormWrapper
      {...args}
      error="Authentication failed"
      usernameError="Username is required"
      passwordError="Password is required"
    />
  ),
};
