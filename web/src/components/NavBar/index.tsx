import type { FC } from "react"
import React, {useState } from "react"
import { Box, Image, } from "@chakra-ui/react"
import { GeneralButton, ConnectWallet, TextButton, HowItWorksModal, SystemPromptsModal }  from '@/components'
import { LogoImg } from '@/assets/images'
import { useNavigate } from "react-router-dom"
import { userAgentsApi } from "@/api"
import { useAppDispatch } from "@/redux/hooks"
import { myAgentsAction } from "@/redux/reducer"
import { PAGE_MAX_WIDTH } from "@/config"

export const Navbar: FC = () => {
  const [isOpen, setOpen] = useState<boolean>(false)
  const [isSystemOpen, setSystemOpen] = useState<boolean>(false)
  
  const navigate = useNavigate()
  const dispatch = useAppDispatch()

  const createAgent = () => {
    navigate('/create-agent')
  }

  const myAgent = async() => {
    // navigate('/my-agent')
    window.open('/my-agent', '_blank')
    // const a = await profileApi()

    
  }

  const onLogo = () => {
    dispatch(myAgentsAction([]))
    navigate('/')
  }

  return (
    <Box 
      className="fx-row ai-ct jc-sb  w100" 
      maxW={PAGE_MAX_WIDTH}
      h="100px"  
      pos='fixed' 
      top={0} 
      px={['0px','0px','0px','0px','30px']}
      // borderColor={['pink','white','yellow','blue','white',]}
      // borderWidth="1px"
      // borderStyle='solid'
     
      bgColor='#16141F'> 

      <Box className="fx-row ai-ct jc-sb w100" maxW='884px' minW='684px'>
        <Image src={LogoImg} h="80px" w="54.28px" className="click" onClick={onLogo}/>
        <GeneralButton title="create agent" onClick={createAgent} style={{  }}/>
        <GeneralButton title="my agents" onClick={myAgent}/>
        <TextButton
          title="how it works"
          onClick={() => setOpen(true)}
        />
        <TextButton
          title="system prompts"
          onClick={() => setSystemOpen(true)}
        />
      </Box>
        
      <ConnectWallet/>
      <HowItWorksModal visible={isOpen} onClose={() => setOpen(false)} />
      <SystemPromptsModal visible={isSystemOpen} onClose={() => setSystemOpen(false)} />
    </Box>
  )
}
