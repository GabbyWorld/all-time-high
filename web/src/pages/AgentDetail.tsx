import React, { useState, useEffect,useRef } from "react"
import { Spinner, Box } from "@chakra-ui/react"
import { AgentItem, ExpandTable, BackButton } from '@/components'
import { WebSocketManager } from '@/ws/WebSocketManager'
import { CSSTransition, TransitionGroup } from 'react-transition-group'
import "@/components/Agent/BattleLogs.css"
import { battlesApi, agent1Api } from '@/api'
import { iAgentReturn, iBattlesReturn } from "@/types"
import { AGENT_ITEM_MAX_WIDTH } from "@/config"
import { useLocation} from "react-router-dom"

export const AgentDetail = () => {
    const [isLoading, setLoading] = useState(false)
    const listContainerRef = useRef<HTMLDivElement>(null)
    const [battleDetail, setBattleDetail] = useState<iBattlesReturn | null>(null)    
    const [agentItem, setAgentItem] = useState<iAgentReturn>()

    const location = useLocation()
    const queryParams = new URLSearchParams(location.search)
    const id = queryParams.get('id')

    useEffect(() => {
        id && fetch(Number(id))
    },[id])

    const fetch = async(agentId: number) => {
        setLoading(true)
        const a = await agent1Api(agentId)
        if(a) {
            setAgentItem(a)
        }
        const b = await battlesApi(agentId)
        setLoading(false)
    }
   
    const onExpand = async(id: number, detail: iBattlesReturn) => {
        setBattleDetail(detail)
    }
    
    const onCollapse = () => {
        setBattleDetail(null)
    }
    return (          
        <Box h="750px" className="">
            {
                isLoading ? 
                <Box h="750px" className="center w100" w={AGENT_ITEM_MAX_WIDTH}>
                    <Spinner size="md" color="white" h="32px" w="32px"/>
                </Box>
                :
                <Box maxH="750px" overflowY="scroll" className="" ref={listContainerRef}>
                    <TransitionGroup className=''>
                        {
                            agentItem && 
                            <AgentItem 
                                isAutoExpand={true}
                                onExpand={onExpand}
                                onCollapse={onCollapse}
                                selectedItem={null}
                                activeIdx={0} 
                                {...agentItem}/>
                        }
                      
                    </TransitionGroup>
                    <CSSTransition
                        in={true}
                        classNames="fade"
                        timeout={300}
                    >                       
                        <ExpandTable agentId={Number(id)} onCollapse={onCollapse} detail={battleDetail}/>
                    </CSSTransition>           
                </Box>
            }
        </Box>
    )
}
