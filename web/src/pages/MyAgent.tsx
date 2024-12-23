import React, { useState, useEffect,useRef } from "react"
import { Spinner, Box } from "@chakra-ui/react"
import { AgentItem, ExpandTable, BackButton } from '@/components'
import { WebSocketManager } from '@/ws/WebSocketManager'
import { CSSTransition, TransitionGroup } from 'react-transition-group'
import "@/components/Agent/BattleLogs.css"
import { useAppDispatch } from "@/redux/hooks"
import {  myAgentsAction } from "@/redux/reducer"
import { userAgentsApi } from '@/api'
import { iAgentReturn, iBattlesReturn } from "@/types"
import { AGENT_ITEM_MAX_WIDTH } from "@/config"
import { useNavigate } from "react-router-dom"



export const MyAgent = () => {
    const [isLoading, setLoading] = useState(false)
    const [selectedItem, setSelectedItem] = useState<number | null>(null) 
    const itemsRef = useRef<iAgentReturn[]>([])
    const battleWSRef = useRef<any>(null)
    const listContainerRef = useRef<HTMLDivElement>(null)
    const [battleDetail, setBattleDetail] = useState<iBattlesReturn | null>(null)

    const dispatch = useAppDispatch()
    const navigate = useNavigate()
    const [agentsList, setAgentsList] = useState<iAgentReturn[]> ([])


    useEffect(() => {
        fetch()
    },[])


    // useEffect(() => {
    //     battleLogsWS()
    //     const wsManager = new WebSocketManager(`${import.meta.env.VITE_API_URL}/ws/agents`, {
    //       heartbeatInterval: 10000,
    //       reconnectInterval: 3000,  
    //       maxReconnectAttempts: 5,  
    //       onOpen: () => {
            
    //       },
    //       onClose: (event: CloseEvent) => {
            
    //       },
    //       onMessage: (message: any) => {
    //         if(message) {
    //             const newList = [JSON.parse(message), ...itemsRef.current]
    //             setItems(newList)           
    //         }
    //       },
    //       onError: (error: ErrorEvent) => {
            
    //       },
    //     })    
    //     return () => {
    //       wsManager.close()
    //       battleWSRef.current.close()
    //     }
    // }, [])

    const fetch = async() => {
        setLoading(true)
        const b = await userAgentsApi()
        setLoading(false)
        if(b && b.agents) {
            setAgentsList(b.agents)
        }
    }
   
    const onExpand = async(id: number, detail: iBattlesReturn) => {
        setSelectedItem(id);     
        setAgentsList((prevItems) => {
            const newItems = prevItems.filter((item) => item.id !== id);
            const selectedItem = prevItems.find((item) => item.id === id);
            return selectedItem ? [selectedItem, ...newItems] : newItems;
        })
        
        setBattleDetail(detail)
    }
    
    const onCollapse = () => {
        setSelectedItem(null)
        setAgentsList(itemsRef.current)     
        setBattleDetail(null)
    }

    // const battleLogsWS = async() => {
    //     const battleWS = new WebSocketManager(`${import.meta.env.VITE_API_URL}/ws/battle`, {
    //         heartbeatInterval: 10000,
    //         reconnectInterval: 3000,  
    //         maxReconnectAttempts: 5,  
    //         onOpen: () => {
    //             console.log('battle ws onOpen')
    //         },
    //         onClose: (event: CloseEvent) => {
    //             console.log('battle ws onClose')
    //         },
    //         onMessage: (message: any) => {
    //             if(message) {
    //                 console.log('battle ws onMessage', JSON.parse(message))
    //                 const { battle, type } = JSON.parse(message)
    //                 if(type === "BATTLE_RESULT") {
    //                     dispatch(lastBattleLogAction(battle))
    //                     dispatch(lastBattleLogTableAction(battle))
    //                 }
    //             }   
    //         },
    //         onError: (error: ErrorEvent) => {
    //             console.log('battle ws onError', error)
    //         },
    //       })
    //       battleWSRef.current = battleWS
    // }

  
    return (    
        <Box>
          
            {/* <BackButton
                onClick={() => navigate('/')}
                style={{
                    marginTop: '-38px'
                }}
            /> */}
            
            <Box h="750px" className="">
                {
                    isLoading ? 
                    <Box h="750px" className="center w100" w={AGENT_ITEM_MAX_WIDTH}>
                        <Spinner size="md" color="white" h="32px" w="32px"/>
                    </Box>
                    :
                    <Box maxH="750px" overflowY="scroll" className="" ref={listContainerRef}>
                        <TransitionGroup className=''>
                            {agentsList.map((item, index) => (
                                <AgentItem 
                                    key={item.id} 
                                    onExpand={onExpand}
                                    onCollapse={onCollapse}
                                    selectedItem={selectedItem}
                                    activeIdx={index} {...item}/>
                            ))}
                        </TransitionGroup>

                        <CSSTransition
                            in={selectedItem !== null}
                            classNames="fade"
                            timeout={300}
                        >
                            <Box style={{ display: selectedItem ? "table" : "none" }} className="w100">
                                <ExpandTable agentId={selectedItem} onCollapse={onCollapse} detail={battleDetail}/>
                            </Box>                
                        </CSSTransition>           
                    </Box>
                }
            </Box>
        </Box>
    );
};
