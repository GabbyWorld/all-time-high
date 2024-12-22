


import React, { useState } from "react";
import { Image, Text, Box } from "@chakra-ui/react";
import { Kline3Img } from "@/assets/images"


interface iLeaderboard {
  
}

export const Leaderboard: React.FC<iLeaderboard> = ({

}) => {
    return (    
        <Box 
            maxW="330px" 
            minW="188px"
            w="100%"
            h="750px"
            border="1px solid #999999" 
            borderRadius="10px" 
            px="16px" 
            py="11px"
        >
            <Box className="fx-row ai-ct">
                <Image src={Kline3Img} w="17px" h="49px"/>
                <Text className="fz32 fw700 main" ml="15px">leaderboard</Text>
            </Box>
            <Box className="w100 h100 center">
                <Text className="fz16 gray9">comming soon</Text>
            </Box>
        </Box>
    );
};
