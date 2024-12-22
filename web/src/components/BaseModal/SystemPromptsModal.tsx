import type { FC } from "react";
import React, {useState, useEffect, useCallback } from "react";

import { Box, Image, Text, Modal, ModalOverlay, ModalContent, ModalBody } from "@chakra-ui/react";
import { GeneralButton, ConnectWallet }  from '@/components'
import { Kline3Img, Kline1Img } from '@/assets/images'

export const SystemPromptsModal: FC<{ visible: boolean, onClose: () => void }> = ({ visible, onClose }) => {
  
  return (
    <Modal isOpen={visible} onClose={onClose} isCentered>
        <ModalOverlay />
        <ModalContent w="1052px"  borderRadius="5px" border='1px solid #01FDB2' bgColor='#16141F'>
            <ModalBody className="fx-col ai-ct">
                <Box className="fx-row ai-ct">
                    <Text className="main fz32 fw700" ml="15px">system prompts</Text>
                </Box>


                <Box className="fx-row ai-ct" p="20px" mt="30px" borderRadius="5px" border='1px solid #01FDB2' w="972px" h="90px" >
                    <Image src={Kline3Img } w="17px" h="49px"/>
                    <Text className="fz14 white" ml="15px">
                        <span className="main">agent prompt: </span>used to generate agent description
                    </Text>
                </Box>

                <Box className="fx-row " mt="15px" w="972px" >
                    <Image src={Kline1Img } w="5px" h="28px"/>
                    <Text className="fz14 white" ml="15px">
                        <p>You’re a creative storyteller and game designer with a talent for crafting engaging character descriptions in a vibrant gaming universe. Your task is to write a short, euphemistic description for an Agent in a player-vs-player battle arena.</p>
                        <p>The description should subtly reflect the Agent’s prompt without directly revealing its purpose or abilities, captivating players and sparking their imagination. Avoid mentioning the Agent’s name in the description. Keep it under 160 characters.</p>
                        <p>Agent’s Name is: <span className="main">#NAME#</span></p>
                        <p>Agent’s Prompt is: <span className="main">#PROMPT#</span></p>                                            
                    </Text>                   
                </Box>             

                <Box className="fx-row ai-ct" p="20px" mt="30px" borderRadius="5px" border='1px solid #01FDB2' w="972px" h="90px" >
                    <Image src={Kline3Img } w="17px" h="49px"/>
                    <Text className="fz14 white" ml="15px">
                        <span className="main">battle prompt: </span>used to generate battle result and description
                    </Text>
                </Box>  

                <Box className="fx-row " mt="15px" w="972px" >
                    <Image src={Kline1Img } w="5px" h="28px"/>
                    <Text className="fz14 white" ml="15px">
                        <p>You’re a game system tasked with predicting the outcome of battles in a player-vs-player arena where user-generated AI agents compete. Your role is to evaluate agents fairly and impartially, focusing solely on their prompts and how their described abilities might interact logically in a real encounter.</p>
                        <p>Analyse the following agents:</p>
                        <ul className="ml24">
                            <li>Attacker Agent Name: <span className="main"> #ATT_NAME#</span></li>
                            <li>Attacker Agent Prompt: <span className="main">#ATT_PROMPT#</span></li>
                            <li>Defender Agent Name: <span className="main">#DEF_NAME#</span></li>
                            <li>Defender Agent Prompt: <span className="main">#DEF_PROMPT#</span></li>
                        </ul>
                       
                        <p>Output Instructions:</p>
                        <div className="ml12">
                            <p>1.Begin by stating the Attack Outcome:</p>
                            <ul className="ml36">
                                <li>“Total Victory!” for clear domination by the Attacker.</li>
                                <li>“Narrow Victory!” for a slight edge to the Attacker.</li>
                                <li>“Narrow Defeat!” for a slight edge to the Defender.</li>
                                <li>“Crushing Defeat!” for clear domination by the Defender.</li>
                            </ul>
                            <p>2.Craft a story under 280 characters, reflecting the battle and its outcome:</p>
                            <ul className="ml36">
                                <li>Mention both names but base the narrative entirely on the interaction of abilities.</li>
                                <li>Avoid directly describing their abilities; focus on the imaginative depiction of how the battle unfolded.</li>
                                <li>Avoid assumptions or biases based on names; rely only on logical implications of abilities.</li>
                                <li>Ensure the outcome aligns with how one ability counters, overpowers, or is neutralized by another.</li>
                            </ul>
                        </div>
                    </Text>                   
                </Box>    


                <GeneralButton title="ok" onClick={onClose} style={{ height: '50px', width: '528px', marginTop: '30px' }}/>
            </ModalBody>                
        </ModalContent>
    </Modal>
);
};


