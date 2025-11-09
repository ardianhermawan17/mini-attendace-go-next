# Quick Start Guide

Get the Mini Attendance Frontend running in 5 minutes.

## Prerequisites

- Node.js 18+
- npm or yarn
- Backend API running on `http://localhost:8080`

## Installation

### 1. Install Dependencies

```bash
npm install
```

### 2. Configure Environment

```bash
cp .env.example .env.local
```

The default `.env.local` is already configured for local development.

### 3. Start Development Server

```bash
npm run dev
```

### 4. Open in Browser

Navigate to [http://localhost:3000](http://localhost:3000)

## First Steps

### Login

1. Use credentials from backend seeder:
   - Username: `admin`
   - Password: `admin123`

2. Or create a new user via backend API

### Check In

1. Click "Get Location" to enable geolocation
2. Click "Check In" to record your arrival
3. Location is automatically captured

### Check Out

1. Click "Get Location" again
2. Click "Check Out" to record your departure

### View Reports

Navigate to Reports section to view:
- Monthly attendance summary
- Department statistics
- Individual attendance records

## Available Scripts

```bash
# Development
npm run dev          # Start dev server

# Production
npm run build        # Build for production
npm start            # Start production server

# Testing
npm test             # Run tests
npm run test:watch   # Run tests in watch mode
npm run test:coverage # Generate coverage report

# Code Quality
npm run lint         # Run ESLint
npm run format       # Format code with Prettier
npm run type-check   # Check TypeScript types
```

## Project Structure

```
src/
├── app/              # Next.js pages
├── components/       # Atomic Design components
├── features/         # Feature logic
├── store/            # Redux store
├── hooks/            # Custom hooks
├── lib/              # Utilities
├── providers/        # Context providers
├── types/            # TypeScript types
└── utils/            # Helper functions
```

## Key Features

- ✅ JWT Authentication
- ✅ Check-in/Check-out with geolocation
- ✅ Optimistic updates
- ✅ Attendance reports
- ✅ Responsive design
- ✅ Type-safe with TypeScript

## Common Tasks

### Add a New Component

1. Create component in `src/components/atoms/` (or molecules/organisms)
2. Create `*-ui.tsx` for presentation
3. Create `use-*.ts` for logic (if needed)
4. Create `index.tsx` to combine them

### Add a New Feature

1. Create folder in `src/features/`
2. Create `use-*.ts` hooks
3. Create container components
4. Add Redux slices if needed
5. Add RTK Query endpoints if needed

### Add Tests

1. Create `*.test.tsx` next to component
2. Use React Testing Library
3. Mock Redux store if needed
4. Run `npm test`

## Troubleshooting

### Port 3000 Already in Use

```bash
# Kill process on port 3000
lsof -i :3000 | grep LISTEN | awk '{print $2}' | xargs kill -9
```

### API Connection Error

1. Check backend is running: `http://localhost:8080`
2. Verify `NEXT_PUBLIC_API_URL` in `.env.local`
3. Check browser console for CORS errors

### Build Fails

```bash
# Clear cache and reinstall
rm -rf .next node_modules
npm install
npm run build
```

### Tests Fail

```bash
# Clear Jest cache
npm test -- --clearCache
npm test
```

## Next Steps

1. Read [README.md](./README.md) for detailed documentation
2. Check [DEPLOYMENT.md](./DEPLOYMENT.md) for production setup
3. Review component examples in `src/components/`
4. Explore Redux store in `src/store/`

## Support

- Backend API docs: `../backend/API_DOCUMENTATION.md`
- Architecture: `../docs/ARCHITECTURE.md`
- Issues: Check GitHub issues or create new one

## Tips

- Use Redux DevTools browser extension for debugging
- Use Chakra UI documentation: https://chakra-ui.com
- Use Next.js documentation: https://nextjs.org
- Use Redux Toolkit documentation: https://redux-toolkit.js.org

Happy coding! 🚀
