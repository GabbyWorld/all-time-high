import React from "react"
import {
  Box,
  Popover,
  PopoverTrigger,
  PopoverContent,
  PopoverBody,
  PopoverArrow,
} from "@chakra-ui/react"

interface iTextPopover {
    content: React.ReactNode
    children: React.ReactNode
}
export const TextPopover: React.FC<iTextPopover> = ({ content, children }) => {
  return (
    <Box textAlign="center" className="click">
      <Popover trigger='hover'>
        <PopoverTrigger >
            { children }
        </PopoverTrigger>
        <PopoverContent maxW="500px"  borderRadius="5px" border='1px solid #01FDB2' bgColor='#16141F'>
          <PopoverArrow bgColor="#16141F" boxShadow='-1px -1px 0px 0 #01FDB2' />
          <PopoverBody p="8px">
            { content }
          </PopoverBody>
        </PopoverContent>
      </Popover>
    </Box>
  );
};

