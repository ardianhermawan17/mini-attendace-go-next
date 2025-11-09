import React from 'react';
import { CheckInCardUI } from '@/components/molecules/CheckInCard/check-in-card-ui';
import { useCheckIn } from './use-check-in';
import { useCheckOut } from './use-check-out';

/**
 * Check-In Card Container
 * Combines the UI component with the logic hooks
 */
export const CheckInCard: React.FC = () => {
  const {
    locationError,
    isGettingLocation,
    isCheckingIn,
    checkInError,
    todayRecord,
    getLocation,
    handleCheckIn,
  } = useCheckIn();

  const { isCheckingOut, checkOutError, handleCheckOut } = useCheckOut();

  return (
    <CheckInCardUI
      todayRecord={todayRecord}
      isCheckingIn={isCheckingIn}
      isCheckingOut={isCheckingOut}
      isGettingLocation={isGettingLocation}
      locationError={locationError}
      checkInError={checkInError}
      checkOutError={checkOutError}
      onCheckIn={handleCheckIn}
      onCheckOut={handleCheckOut}
      onGetLocation={getLocation}
    />
  );
};

export default CheckInCard;
