import type { Meta } from '@storybook/react';
import { CheckInCardUI } from './check-in-card-ui';
import { AttendanceRecord } from '@/types';

const meta: Meta<typeof CheckInCardUI> = {
  title: 'Molecules/CheckInCard',
  component: CheckInCardUI,
  parameters: {
    layout: 'centered',
  },
  tags: ['autodocs'],
};

export default meta;

const mockAttendanceRecord: AttendanceRecord = {
  id: '1',
  user_id: 'user-1',
  attendance_date: new Date().toISOString().split('T')[0],
  check_in_time: new Date().toISOString(),
  check_out_time: null,
  check_in_location: 'Jakarta Office',
  check_out_location: undefined,
  status: 'hadir',
  created_at: new Date().toISOString(),
  updated_at: new Date().toISOString(),
};

const mockCheckOutRecord: AttendanceRecord = {
  ...mockAttendanceRecord,
  check_out_time: new Date(Date.now() + 8 * 60 * 60 * 1000).toISOString(),
  check_out_location: 'Jakarta Office',
};

export const NotCheckedIn = {
  args: {
    todayRecord: null,
    isCheckingIn: false,
    isCheckingOut: false,
    isGettingLocation: false,
    onCheckIn: () => console.log('Check in'),
    onCheckOut: () => console.log('Check out'),
    onGetLocation: () => console.log('Get location'),
  },
};

export const CheckedIn = {
  args: {
    todayRecord: mockAttendanceRecord,
    isCheckingIn: false,
    isCheckingOut: false,
    isGettingLocation: false,
    onCheckIn: () => console.log('Check in'),
    onCheckOut: () => console.log('Check out'),
    onGetLocation: () => console.log('Get location'),
  },
};

export const CheckedOut = {
  args: {
    todayRecord: mockCheckOutRecord,
    isCheckingIn: false,
    isCheckingOut: false,
    isGettingLocation: false,
    onCheckIn: () => console.log('Check in'),
    onCheckOut: () => console.log('Check out'),
    onGetLocation: () => console.log('Get location'),
  },
};

export const CheckingIn = {
  args: {
    todayRecord: null,
    isCheckingIn: true,
    isCheckingOut: false,
    isGettingLocation: false,
    onCheckIn: () => console.log('Check in'),
    onCheckOut: () => console.log('Check out'),
    onGetLocation: () => console.log('Get location'),
  },
};

export const CheckingOut = {
  args: {
    todayRecord: mockAttendanceRecord,
    isCheckingIn: false,
    isCheckingOut: true,
    isGettingLocation: false,
    onCheckIn: () => console.log('Check in'),
    onCheckOut: () => console.log('Check out'),
    onGetLocation: () => console.log('Get location'),
  },
};

export const GettingLocation = {
  args: {
    todayRecord: null,
    isCheckingIn: false,
    isCheckingOut: false,
    isGettingLocation: true,
    onCheckIn: () => console.log('Check in'),
    onCheckOut: () => console.log('Check out'),
    onGetLocation: () => console.log('Get location'),
  },
};

export const WithLocationError = {
  args: {
    todayRecord: null,
    isCheckingIn: false,
    isCheckingOut: false,
    isGettingLocation: false,
    locationError: 'Unable to get your location. Please enable location services.',
    onCheckIn: () => console.log('Check in'),
    onCheckOut: () => console.log('Check out'),
    onGetLocation: () => console.log('Get location'),
  },
};

export const WithCheckInError = {
  args: {
    todayRecord: null,
    isCheckingIn: false,
    isCheckingOut: false,
    isGettingLocation: false,
    checkInError: 'Failed to check in. Please try again.',
    onCheckIn: () => console.log('Check in'),
    onCheckOut: () => console.log('Check out'),
    onGetLocation: () => console.log('Get location'),
  },
};

export const WithCheckOutError = {
  args: {
    todayRecord: mockAttendanceRecord,
    isCheckingIn: false,
    isCheckingOut: false,
    isGettingLocation: false,
    checkOutError: 'Failed to check out. Please try again.',
    onCheckIn: () => console.log('Check in'),
    onCheckOut: () => console.log('Check out'),
    onGetLocation: () => console.log('Get location'),
  },
};
