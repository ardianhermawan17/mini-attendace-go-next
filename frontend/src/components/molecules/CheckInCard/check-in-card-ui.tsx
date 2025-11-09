import React from 'react';
import { VStack, HStack, Box, Alert, AlertIcon, Badge } from '@chakra-ui/react';
import { CardUI } from '@/components/atoms/Card/card-ui';
import { ButtonUI } from '@/components/atoms/Button/button-ui';
import { HeadingUI, TextUI } from '@/components/atoms/Text/text-ui';
import { AttendanceRecord } from '@/types';
import { format } from 'date-fns';

export interface CheckInCardUIProps {
  todayRecord: AttendanceRecord | null;
  isCheckingIn?: boolean;
  isCheckingOut?: boolean;
  isGettingLocation?: boolean;
  locationError?: string | null;
  checkInError?: string | null;
  checkOutError?: string | null;
  onCheckIn: () => void;
  onCheckOut: () => void;
  onGetLocation: () => void;
}

/**
 * Check-In Card UI Component
 * Presentational component that only accepts props
 */
export const CheckInCardUI: React.FC<CheckInCardUIProps> = ({
  todayRecord,
  isCheckingIn = false,
  isCheckingOut = false,
  isGettingLocation = false,
  locationError,
  checkInError,
  checkOutError,
  onCheckIn,
  onCheckOut,
  onGetLocation,
}) => {
  const hasCheckedIn = !!todayRecord?.check_in_time;
  const hasCheckedOut = !!todayRecord?.check_out_time;

  return (
    <CardUI variant="elevated">
      <VStack spacing={6} align="stretch">
        <HeadingUI level={2} size="lg">
          Attendance
        </HeadingUI>

        {(checkInError || checkOutError || locationError) && (
          <Alert status="error" borderRadius="md">
            <AlertIcon />
            {checkInError || checkOutError || locationError}
          </Alert>
        )}

        {/* Status Display */}
        <Box>
          <TextUI variant="label" mb={2}>
            Status
          </TextUI>
          <HStack spacing={4}>
            <Badge colorScheme={hasCheckedIn ? 'green' : 'gray'}>
              {hasCheckedIn ? '✓ Checked In' : 'Not Checked In'}
            </Badge>
            {hasCheckedIn && (
              <Badge colorScheme={hasCheckedOut ? 'blue' : 'yellow'}>
                {hasCheckedOut ? '✓ Checked Out' : 'Checked In'}
              </Badge>
            )}
          </HStack>
        </Box>

        {/* Time Display */}
        {todayRecord && (
          <VStack align="start" spacing={2}>
            {todayRecord.check_in_time && (
              <Box>
                <TextUI variant="label">Check-In Time</TextUI>
                <TextUI size="md">{format(new Date(todayRecord.check_in_time), 'HH:mm:ss')}</TextUI>
              </Box>
            )}
            {todayRecord.check_out_time && (
              <Box>
                <TextUI variant="label">Check-Out Time</TextUI>
                <TextUI size="md">
                  {format(new Date(todayRecord.check_out_time), 'HH:mm:ss')}
                </TextUI>
              </Box>
            )}
          </VStack>
        )}

        {/* Action Buttons */}
        <VStack spacing={3} align="stretch">
          {!hasCheckedIn && (
            <>
              <ButtonUI
                colorScheme="blue"
                onClick={onGetLocation}
                isLoading={isGettingLocation}
                isDisabled={isGettingLocation || isCheckingIn}
              >
                {isGettingLocation ? 'Getting Location...' : 'Get Location'}
              </ButtonUI>
              <ButtonUI
                colorScheme="green"
                onClick={onCheckIn}
                isLoading={isCheckingIn}
                isDisabled={isCheckingIn || isGettingLocation || locationError !== null}
              >
                {isCheckingIn ? 'Checking In...' : 'Check In'}
              </ButtonUI>
            </>
          )}

          {hasCheckedIn && !hasCheckedOut && (
            <>
              <ButtonUI
                colorScheme="blue"
                onClick={onGetLocation}
                isLoading={isGettingLocation}
                isDisabled={isGettingLocation || isCheckingOut}
              >
                {isGettingLocation ? 'Getting Location...' : 'Get Location'}
              </ButtonUI>
              <ButtonUI
                colorScheme="orange"
                onClick={onCheckOut}
                isLoading={isCheckingOut}
                isDisabled={isCheckingOut || isGettingLocation || locationError !== null}
              >
                {isCheckingOut ? 'Checking Out...' : 'Check Out'}
              </ButtonUI>
            </>
          )}

          {hasCheckedOut && (
            <TextUI variant="caption" textAlign="center" color="green.600">
              You have completed your attendance for today
            </TextUI>
          )}
        </VStack>
      </VStack>
    </CardUI>
  );
};

export default CheckInCardUI;
