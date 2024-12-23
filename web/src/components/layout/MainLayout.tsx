import React from "react"
import { Box } from "@chakra-ui/react"
import { Notification } from '@/components'
import { useAppDispatch, useAppSelector } from "@/redux/hooks"
import { notificationInfoAction, selectNotificationInfo } from "@/redux/reducer"
import { PAGE_MAX_WIDTH, PAGE_MIN_WIDTH} from "@/config"

export const MainLayout: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => {
  const { open, title } = useAppSelector(selectNotificationInfo)
  const dispatch = useAppDispatch()

  return (
    <Box className="w100 h100">
      <Box className="w100 h100" border="1px solid transparent" maxW={PAGE_MAX_WIDTH} minW={PAGE_MIN_WIDTH} >
        { children }
      </Box>
      <Notification 
        visible={open}
        onClose={() => dispatch(notificationInfoAction({ open: false, title: '' }))}
        title={title}
      />
    </Box>
  )
}