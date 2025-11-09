import { useState, useCallback } from 'react';
import { useAttendance } from '@/hooks/useAttendance';
import { CheckOutRequest } from '@/types';

export const useCheckOut = () => {
  const { checkOut, isCheckingOut, checkOutError } = useAttendance();
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

  const handleCheckOut = useCallback(async () => {
    if (!location) {
      setLocationError('Location is required for check-out');
      return;
    }

    try {
      const request: CheckOutRequest = {
        latitude: location.latitude,
        longitude: location.longitude,
        device: 'web',
      };

      await checkOut(request);
    } catch (error) {
      console.error('Check-out error:', error);
    }
  }, [location, checkOut]);

  return {
    location,
    locationError,
    isGettingLocation,
    isCheckingOut,
    checkOutError,
    getLocation,
    handleCheckOut,
  };
};
