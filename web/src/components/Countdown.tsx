import React, { useState, useEffect } from 'react'
import { Text } from '@chakra-ui/react'

interface CountdownProps {
  endTime: number
}

interface TimeLeft {
  H?: number
  M?: number
  S?: number
}

export const Countdown: React.FC<CountdownProps> = ({ endTime }) => {
  const [isEnd, setEnd] = useState<boolean>(false)
  const [timeLeft, setTimeLeft] = useState<TimeLeft>({})

  const calculateTimeLeft = (): TimeLeft => {
    const now = new Date().getTime()
    const difference = endTime - now

    if (difference <= 0) {
      setEnd(true)
      return {}
    }

    return {
      H: Math.floor((difference / (1000 * 60 * 60)) % 24),
      M: Math.floor((difference / 1000 / 60) % 60),
      S: Math.floor((difference / 1000) % 60)
    }
  }

  useEffect(() => {
    const timer = setInterval(() => {
      setTimeLeft(calculateTimeLeft())
    }, 1000)

    return () => clearInterval(timer)
  }, [endTime])

  const timerComponents = Object.keys(timeLeft).map(interval => {
    if (timeLeft[interval as keyof TimeLeft] === undefined) {
      return null
    }

    return (
      <span key={interval}>
        {timeLeft[interval as keyof TimeLeft]} {interval !== 'S' ? ':' : ''}
      </span>
    )
  })

  return (
    <Text className="fz12" color="#BD422C">
      {isEnd ? (
        <span>Ended</span>
      ) : (
        <>
          {timerComponents.length ? timerComponents : <span>Not Started</span>}
        </>
      )}
    </Text>
  )
}


