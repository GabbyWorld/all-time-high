
import { useEffect } from 'react'
import { Home } from '@/pages/'
import { CreateAgent } from '@/pages/CreateAgent'
import { MyAgent } from '@/pages/MyAgent'
import { AgentDetail } from '@/pages/AgentDetail'
import { MainLayout } from '@/components/layout'
import {BrowserRouter, Route, Routes } from "react-router-dom"
import { ChakraProvider, extendTheme } from '@chakra-ui/react'

import { ReduxProvider } from '@/redux/ReduxProvider'
// import { SOLContextProvider } from '@/lib/solwallet/SOLContextProvider'
import { Navbar, Leaderboard } from '@/components'
import { Box } from "@chakra-ui/react"
// import { WalletsProvider } from '@/lib/solwallet/Wallet'
import { SOLProvider } from '@/lib/solwallet/SOLProvider'

import './types/window.d.ts'

const theme = extendTheme({
  fonts: {
    // body: `'Body Font Name', Salsa`,
  }
})

function App() {
  return (
    <ReduxProvider>
      <ChakraProvider resetCSS theme={theme}>
        <SOLProvider>
          <MainLayout>
            <BrowserRouter>
              <Navbar/>
          
                <Box mt="150px"  className='fx-row ai-ct jc-sb' px={['0px','0px','0px','0px','30px']}>
                  <Routes>          
                    <Route path="/" element={<Home/>}/>               
                    <Route path="/create-agent" element={<CreateAgent/>}/>               
                    <Route path="/my-agent" element={<MyAgent/>}/>               
                    <Route path="/agent-detail" element={<AgentDetail/>}/>               
                  </Routes>
                  <Leaderboard/>
                </Box>
                         
            </BrowserRouter>
          </MainLayout>
        </SOLProvider>
      </ChakraProvider>
    </ReduxProvider>
  )
}


export default App;
