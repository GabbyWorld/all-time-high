import React from "react"
import { Box, Image, Text } from "@chakra-ui/react"
import { ClickButtonWrapper } from './ClickButtonWrapper'
import { ArrowLeftImg } from '@/assets/images'

interface iBackButton {
  onClick: () => void
  style?: React.CSSProperties 
}

export const BackButton: React.FC<iBackButton> = ({  
  onClick,
  style
}) => {

  return (
    <ClickButtonWrapper onClick={onClick} disable={false} clickableDisabled={true}> 
        <Box
            className="fx-row ai-ct click" 
            style={style}
        >
            <Image src={ArrowLeftImg} h="12px" w="7px" />
            <Text className="ml8 main fz24 fw700">back</Text>
        </Box>  
    </ClickButtonWrapper>
  );
};
