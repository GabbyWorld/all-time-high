


import React from "react";
import { Flex, Box, Spinner, Text } from "@chakra-ui/react";
import { ClickButtonWrapper } from './ClickButtonWrapper'


interface iGeneralButton {
  loading?: boolean;
  disable?: boolean;
  onClick: () => any;
  title: string;
  style?: React.CSSProperties;
}

export const GeneralButton: React.FC<iGeneralButton> = ({

  loading = false,
  disable = false,
  onClick,
  title,
  style

}) => {
  const handleClick = () => {
    if(loading || disable) {
      return false
    }else {
      onClick()
    }
  }



  return (
    <ClickButtonWrapper onClick={handleClick} disable={disable} clickableDisabled={true}> 
      <Box
        className="center"
        cursor={(loading || disable) ? 'not-allowed' : 'pointer'}
        bgColor={ (loading || disable)  ? '#666' : '#01FDB2'}
        color="#16141F"
        _hover={{
          bgColor: "#01553C",
          color: '#CCCCCC'
        }}
        transition="background-color 0.3s ease, color 0.3s ease"
        h='40px'
        w="170px"
        borderRadius="6px"
        style={{
          ...style,
        }}
      >
        {loading ? 
          <Spinner size="md" color="white" h="32px" w="32px"/> : 
          <Text className="fz22 fw700 fm1" >{title}</Text>
        }
      </Box>
    </ClickButtonWrapper>
  );
};
