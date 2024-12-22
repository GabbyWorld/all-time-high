import type { FC } from "react"
import { Box, Text } from "@chakra-ui/react"
import { createStandaloneToast } from '@chakra-ui/react'

const { toast } = createStandaloneToast()

interface ToastProps {
    props: any
}

export const ToastRender:FC<ToastProps> = ({
    props
}) => {
    return (
        <>
            {
                props.status === 'success' && 
                <Box bg='#00230E' border="1px solid #388167" h="32px" className="center">
                    <Text className="white fz14">
                        {props.title}
                    </Text>
                </Box>
            }
            {
                props.status === 'warning' && 
                <Box  bg='#110F28' border="1px solid #332E67" h="32px" className="center">
                    <Box >
                        <Text className="white fz14">
                            {props.title}
                        </Text>
                    </Box>
                </Box>
            }
        </>
    )
}


export const showToast = (title: string, status: 'success' | 'warning') => {
    toast({
        title,
        status,
        render: (props) => <ToastRender props={props} />,
    });
};
