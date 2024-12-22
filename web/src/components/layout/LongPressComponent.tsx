import React, { useRef } from 'react';

interface LongPressComponentProps {
  onLongPress: () => void;
  children: React.ReactNode;
}

export const LongPressComponent: React.FC<LongPressComponentProps> = ({ onLongPress, children }) => {
  const timeoutRef = useRef<any>(null);

  const handleMouseDown = () => {
    timeoutRef.current = setTimeout(() => {
      onLongPress();
    }, 500);
  };

  const handleMouseUp = () => {

    if (timeoutRef.current) {
      clearTimeout(timeoutRef.current);
    }
  };

  return (
    <div
      onMouseDown={handleMouseDown}
      onMouseUp={handleMouseUp}
      onTouchStart={handleMouseDown}
      onTouchEnd={handleMouseUp}
    >
      {children}
    </div>
  );
};


