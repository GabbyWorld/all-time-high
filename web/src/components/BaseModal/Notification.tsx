
import type { FC } from "react"
import React from "react"
import { Box, Image, Text, Modal, ModalOverlay, ModalContent, ModalBody } from "@chakra-ui/react"
import { GeneralButton }  from '@/components'
import { Kline3Img, } from '@/assets/images'

interface iNotification {
    visible: boolean
    onClose: () => void
    title: string
}
export const Notification: FC<iNotification> = ({ visible, onClose, title }) => {
  return (
        <Modal isOpen={visible} onClose={onClose} isCentered>
            <ModalOverlay />
            <ModalContent w="730px" pt="100px" pb="30px" maxWidth="none" borderRadius="5px" border='1px solid #01FDB2' bgColor='#16141F'>
                <ModalBody className="fx-col ai-ct">
                    <Box className="fx-row ai-ct">
                        <Image src={Kline3Img } w="17px" h="49px"/>
                        <Text className="main fz32 fw700" ml="15px" whiteSpace="nowrap">{title}</Text>
                    </Box>                
                    <GeneralButton title="ok" onClick={onClose} style={{ height: '50px', width: '528px', marginTop: '100px' }}/>
                </ModalBody>                
            </ModalContent>
        </Modal>
    )
}
