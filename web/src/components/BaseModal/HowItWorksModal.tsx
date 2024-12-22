import type { FC } from "react";
import React, {useState, useEffect, useCallback } from "react";

import { Box, Image, Text, Modal, ModalOverlay, ModalContent, ModalBody } from "@chakra-ui/react";
import { GeneralButton, ConnectWallet }  from '@/components'
import { Kline3Img, Kline1Img } from '@/assets/images'

export const HowItWorksModal: FC<{ visible: boolean, onClose: () => void }> = ({ visible, onClose }) => {
  
  return (
    <Modal isOpen={visible} onClose={onClose} isCentered>
        <ModalOverlay />
        <ModalContent w="1052px"  borderRadius="5px" border='1px solid #01FDB2' bgColor='#16141F'>
            <ModalBody className="fx-col ai-ct">
                <Box className="fx-row ai-ct">
                    <Image src={Kline3Img } w="17px" h="49px"/>
                    <Text className="main fz32 fw700" ml="15px">how it works</Text>
                </Box>

                <Box className="fx-row " p="20px" mt="30px" borderRadius="5px" border='1px solid #01FDB2' w="972px" h="90px" >
                    <Image src={Kline3Img } w="17px" h="49px"/>
                    <Text className="fz14 white" ml="15px">
                        <span className="main">all-time-high.ai </span>is a gamified consensus layer for tokenised agents that filters out best agent coins through pvp gamification every all-time-high of an agent coin triggers a pvp against other agent(s), and only best-prompted agents deserve most all-time-highs
                    </Text>
                </Box>                
                <Box mt="30px"/>
                {
                    ['insert prompts to generate an ai agent, which simultaneously launches a coin on pump.fun',
                    'whenever your agent coin hits an all-time-high (on 5min candles), it battles another agent',
                    'over time, best-prompted agents will top the leaderbaord, deserving most all-time-highs'
                    ].map((item,idx) => (
                        <Box key={item} className="fx-row ai-ct" mt="15px" w="972px" >
                            <Image src={Kline1Img } w="5px" h="28px"/>
                            <Text className="main fz14" ml="10px">step {idx + 1}:</Text>
                            <Text className="white fz14" ml="7px">{item}</Text>
                        </Box>
                    ))
                }
                <GeneralButton title="ok" onClick={onClose} style={{ height: '50px', width: '528px', marginTop: '30px' }}/>
            </ModalBody>                
        </ModalContent>
    </Modal>
);
};


