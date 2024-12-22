import { useRef, useEffect } from 'react';

export const useDebounce = (callback: () => void, delay: number) => {
  const timer = useRef<number>();

  const debouncedCallback = () => {
    if (timer.current) {
      clearTimeout(timer.current);
    }
    timer.current = window.setTimeout(() => {
      callback();
    }, delay);
  };

  useEffect(() => {
    return () => {
      if (timer.current) {
        clearTimeout(timer.current);
      }
    };
  }, [delay]);

  return debouncedCallback;
};

