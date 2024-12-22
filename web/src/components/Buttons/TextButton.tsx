
import React from "react"
import { Text } from "@chakra-ui/react"

interface iTextButton {
  onClick: () => void
  title: string
}

export const TextButton: React.FC<iTextButton> = ({
  onClick,
  title
}) => {
  return (
    <Text   
        onClick={onClick}
        className="click fz24 fw700 fm1" 
        color="#01FDB2"
        whiteSpace='nowrap'
        _hover={{
            color: '#fff',
            textDecoration: 'underline'
        }}
    >{title}</Text>
  );
};


