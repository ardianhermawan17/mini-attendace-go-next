# Mini Attendance System - Frontend

A modern, responsive attendance management system built with Next.js, Redux Toolkit, and Chakra UI.

## Features

- **JWT Authentication**: Secure login with access and refresh tokens
- **Attendance Management**: Check-in and check-out with geolocation
- **Optimistic Updates**: Instant UI feedback with automatic rollback on errors
- **Atomic Design**: Well-organized component structure
- **State Management**: Redux Toolkit with RTK Query for server state
- **Type Safety**: Full TypeScript support
- **Responsive Design**: Mobile-first approach with Chakra UI
- **Testing**: Jest and React Testing Library setup

## Tech Stack

- **Framework**: Next.js 14
- **UI Library**: Chakra UI
- **State Management**: Redux Toolkit + RTK Query
- **Language**: TypeScript
- **Testing**: Jest + React Testing Library
- **Styling**: Emotion (via Chakra UI)
- **HTTP Client**: Axios

## Project Structure

```
src/
├── app/                          # Next.js app directory
│   ├── layout.tsx               # Root layout with providers
│   ├── page.tsx                 # Home page (redirects to login/dashboard)
│   ├── login/                   # Login page
│   ├── dashboard/               # Dashboard page
│   ├── attendance/              # Attendance pages
│   └── reports/                 # Reports pages
├── components/                   # Atomic Design components
│   ├── atoms/                   # Basic building blocks (Button, Input, Card, Text)
│   ├── molecules/               # Combinations of atoms (LoginForm, CheckInCard)
│   └── organisms/               # Complex components (Header, Sidebar, Layout)
├── features/                     # Feature-specific logic
│   ├── auth/                    # Authentication feature
│   │   ├── use-login.ts        # Login hook
│   │   └── index.tsx           # Login form container
│   └── attendance/              # Attendance feature
│       ├── use-check-in.ts     # Check-in hook
│       ├── use-check-out.ts    # Check-out hook
│       └── index.tsx           # Check-in card container
├── store/                        # Redux store
│   ├── index.ts                # Store configuration
│   ├── api/                    # RTK Query APIs
│   │   ├── authApi.ts
│   │   ├── attendanceApi.ts
│   │   └── reportsApi.ts
│   ├── slices/                 # Redux slices
│   │   ├── auth/
│   │   ├── attendance/
│   │   └── ui/
│   ├── middlewares/            # Custom middlewares
│   │   ├── idempotencyMiddleware.ts
│   │   └── logger.ts
│   └── selectors/              # Memoized selectors
├── hooks/                        # Custom React hooks
│   ├── useAuth.ts
│   └── useAttendance.ts
├── lib/                          # Utility libraries
│   └── api/
│       └── client.ts           # Axios client with interceptors
├── providers/                    # React context providers
│   ├── ReduxProvider.tsx
│   └── ChakraProvider.tsx
├── types/                        # TypeScript type definitions
│   └── index.ts
└── utils/                        # Utility functions
```

## Getting Started

### Prerequisites

- Node.js 18+ 
- npm or yarn

### Installation

1. Install dependencies:

```bash
npm install
```

2. Create `.env.local` file:

```bash
cp .env.example .env.local
```

3. Update environment variables in `.env.local`:

```env
NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1
```

### Development

Start the development server:

```bash
npm run dev
```

Open [http://localhost:3000](http://localhost:3000) in your browser.

### Build

Build for production:

```bash
npm run build
```

Start production server:

```bash
npm start
```

## Testing

Run tests:

```bash
npm test
```

Run tests in watch mode:

```bash
npm run test:watch
```

Generate coverage report:

```bash
npm run test:coverage
```

## Code Quality

Format code:

```bash
npm run format
```

Check formatting:

```bash
npm run format:check
```

Type check:

```bash
npm run type-check
```

## Architecture

### Atomic Design Pattern

Components are organized following the Atomic Design methodology:

- **Atoms**: Basic building blocks (Button, Input, Card, Text)
- **Molecules**: Combinations of atoms (LoginForm, CheckInCard)
- **Organisms**: Complex components (Header, Sidebar, Layout)

### Separation of Concerns

Each feature follows a clear separation:

- **`*-ui.tsx`**: Presentational component (UI only, accepts props)
- **`use-*.ts`**: Custom hook (logic and state management)
- **`index.tsx`**: Container component (combines UI and logic)

Example:

```typescript
// use-login.ts - Logic
export const useLogin = () => {
  const { login, isLoading } = useAuth();
  // ... logic here
  return { formData, handleSubmit, ... };
};

// login-form-ui.tsx - UI
export const LoginFormUI: React.FC<LoginFormUIProps> = (props) => {
  return <form>{/* UI here */}</form>;
};

// index.tsx - Container
export const LoginForm = () => {
  const { formData, handleSubmit } = useLogin();
  return <LoginFormUI {...props} />;
};
```

### State Management

- **Redux Slices**: Local UI state (auth, attendance, ui)
- **RTK Query**: Server state (API calls with caching)
- **Custom Hooks**: Feature-specific logic

### Optimistic Updates

Check-in/check-out operations use optimistic updates:

1. Dispatch optimistic action (update UI immediately)
2. Send request to server
3. On success: Replace optimistic data with real data
4. On error: Rollback to previous state

## API Integration

The frontend integrates with the backend API at `http://localhost:8080/api/v1`.

### Authentication Flow

1. User logs in with credentials
2. Backend returns access and refresh tokens
3. Tokens stored in cookies (secure)
4. Axios interceptor adds token to all requests
5. On token expiry, automatically refresh using refresh token

### RTK Query

Server state is managed through RTK Query with:

- Automatic caching
- Background refetching
- Optimistic updates
- Tag-based invalidation

## Deployment

### Docker

Build Docker image:

```bash
docker build -t mini-attendance-frontend .
```

Run container:

```bash
docker run -p 3000:3000 mini-attendance-frontend
```

### Environment Variables

For production, set these environment variables:

- `NEXT_PUBLIC_API_URL`: Backend API URL
- `NEXT_PUBLIC_APP_NAME`: Application name
- `NEXT_PUBLIC_ENABLE_ANALYTICS`: Enable analytics
- `NEXT_PUBLIC_ENABLE_SENTRY`: Enable error tracking

## Browser Support

- Chrome (latest)
- Firefox (latest)
- Safari (latest)
- Edge (latest)

## Performance

- Code splitting with Next.js
- Image optimization
- CSS-in-JS with Emotion
- Memoized selectors with Reselect
- RTK Query caching

## Security

- JWT tokens in secure cookies
- CSRF protection ready
- XSS prevention with React
- Secure API client with interceptors
- Environment variable validation

## Contributing

1. Follow the Atomic Design pattern
2. Maintain separation of concerns
3. Write tests for new features
4. Use TypeScript for type safety
5. Follow the existing code style

## License

MIT

## Support

For issues and questions, please refer to the main project documentation.
