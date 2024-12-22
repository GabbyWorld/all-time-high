
import { differenceInSeconds, differenceInMinutes, differenceInHours, differenceInDays } from 'date-fns';

export const sleep = (ms: number): Promise<boolean> => {
    return new Promise((resolve) => {
      setTimeout(() => resolve(true), ms)
    })
}


export const formatTimeAgo = (timestamp: number): string => {
  const now = new Date();
  const time = new Date(timestamp);

  const seconds = differenceInSeconds(now, time);
  if (seconds < 60) return `${seconds}s ago`;

  const minutes = differenceInMinutes(now, time);
  if (minutes < 60) return `${minutes}m ago`;

  const hours = differenceInHours(now, time);
  if (hours < 24) return `${hours}h ago`;

  const days = differenceInDays(now, time);
  return `${days}d ago`;
}


export const formatNumber = (value: number): string | number => {
  if(!!!value) return '--'
  if (value < 1000) {
      return `$${Math.floor(value * 100) / 100}`
  }
  const truncatedValue = Math.floor((value / 1000) * 100) / 100
  return `$${truncatedValue}k`
}