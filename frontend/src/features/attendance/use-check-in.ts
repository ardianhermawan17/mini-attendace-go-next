import { useState, useCallback } from 'react';
import { useAttendance } from '@/hooks/useAttendance';
import { CheckInRequest } from '@/types';

export const useCheckIn = () => {
  const { checkIn, todayRecord, isCheckingIn, checkInError } = useAttendance();
  const [location, setLocation] = useState<{ latitude: number; longitude: number } | null>(null);
  const [locationError, setLocationError] = useState<string | null>(null);
  const [isGettingLocation, setIsGettingLocation] = useState(false);

  const getLocation = useCallback(async () => {
    setIsGettingLocation(true);
    setLocationError(null);

    if (!navigator.geolocation) {
      setLocationError('Geolocation is not supported by your browser');
      setIsGettingLocation(false);
      return;
    }

    navigator.geolocation.getCurrentPosition(
      (position) => {
        setLocation({
          latitude: position.coords.latitude,
          longitude: position.coords.longitude,
        });
        setIsGettingLocation(false);
      },
      (error) => {
        setLocationError(`Failed to get location: ${error.message}`);
        setIsGettingLocation(false);
      }
    );
  }, []);

  const handleCheckIn = useCallback(async () => {
    if (!location) {
      setLocationError('Location is required for check-in');
      return;
    }

    try {
      const request: CheckInRequest = {
        latitude: location.latitude,
        longitude: location.longitude,
        device: 'web',
      };

      await checkIn(request);
    } catch (error) {
      console.error('Check-in error:', error);
    }
  }, [location, checkIn]);

  return {
    location,
    locationError,
    isGettingLocation,
    isCheckingIn,
    checkInError,
    todayRecord,
    getLocation,
    handleCheckIn,
  };
};
