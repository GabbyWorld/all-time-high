import { useState, useEffect, useRef, useCallback } from 'react'

interface CountdownResult {
  secondsLeft: number
  isEnd: boolean
  clearTimer: () => void
}

export const useCountdown = (endTime: number): CountdownResult => {
  const [secondsLeft, setSecondsLeft] = useState<number>(0)
  const [isEnd, setIsEnd] = useState<boolean>(false)
  const timerRef = useRef<any>(null)

  const calculateSecondsLeft = useCallback(() => {
    if (!endTime) {
      setSecondsLeft(0)
      setIsEnd(false)
      return
    }

    const now = new Date().getTime()
    const difference = endTime - now
    const seconds = Math.floor(difference / 1000)

    if (seconds <= 0) {
      setIsEnd(true)
      setSecondsLeft(0)
    } else {
      setSecondsLeft(seconds)
    }
  }, [endTime])

  const clearTimer = useCallback(() => {
    if (timerRef.current) {
      clearInterval(timerRef.current)
      timerRef.current = null
    }
  }, [])

  useEffect(() => {
    if (endTime === undefined) {
      setSecondsLeft(0)
      setIsEnd(false)
      return
    }

    timerRef.current = setInterval(calculateSecondsLeft, 1000)

    return () => clearTimer()
  }, [endTime, calculateSecondsLeft, clearTimer])

  return { secondsLeft, isEnd, clearTimer }
}
