import type { Metadata } from 'next';
import { ReduxProvider } from '@/providers/ReduxProvider';
import { ChakraUIProvider } from '@/providers/ChakraProvider';

export const metadata: Metadata = {
  title: 'Mini Attendance System',
  description: 'Employee Attendance Management System',
};

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="en">
      <body>
        <ReduxProvider>
          <ChakraUIProvider>{children}</ChakraUIProvider>
        </ReduxProvider>
      </body>
    </html>
  );
}
