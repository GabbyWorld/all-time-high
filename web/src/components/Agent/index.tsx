import React, { useState, useEffect,useRef } from "react"
import { Spinner, Box } from "@chakra-ui/react"
import { ClickButtonWrapper, GeneralButton, ExpandTable, BackButton } from '@/components'
import { WebSocketManager } from '@/ws/WebSocketManager'
import { CSSTransition, TransitionGroup } from 'react-transition-group'
import "./BattleLogs.css"
import { useAppDispatch, useAppSelector } from "@/redux/hooks"
import { expandBattleLogsAction, lastBattleLogAction, lastBattleLogTableAction, myAgentsAction, selectMyAgents, selectWalletInfo } from "@/redux/reducer"
import { allAgentsApi, battlesApi } from '@/api'
import { iAgentReturn, iBattleItemReturn, iBattlesReturn } from "@/types"
import { AGENT_ITEM_MAX_WIDTH } from "@/config"
import { throttle } from 'lodash'
import { AgentItem } from './AgentItem'

interface iAgent {
  
}

export const Agent: React.FC<iAgent> = ({

}) => {
    const [allDataLoaded, setAllDataLoaded] = useState(false)
    const [isLoading, setLoading] = useState(false)
    const [isLoadingMore, setLoadingMore] = useState(false)
   
    const [items, setItems] = useState<iAgentReturn[]>([]);
    const [selectedItem, setSelectedItem] = useState<number | null>(null) 
    const [pageIndex, setPageIndex] = useState<number>(1)
    // const [pageSize, setPageSize] = useState<number>(2)
    const [pageSize, setPageSize] = useState<number>(8)

    const { isConnected, address } = useAppSelector(selectWalletInfo)
    const itemsRef = useRef<iAgentReturn[]>([])
    const battleWSRef = useRef<any>(null)
    const listContainerRef = useRef<HTMLDivElement>(null)
    const [battleDetail, setBattleDetail] = useState<iBattlesReturn | null>(null)

    const dispatch = useAppDispatch()
    const myAgents = useAppSelector(selectMyAgents)

    const handleScroll = throttle(() => {
        const container = listContainerRef.current;
        if (!isLoadingMore && container) {
            const isBottom = container.scrollHeight - container.scrollTop <= container.clientHeight + 50;
            if (isBottom) {
                setPageIndex((prevPageIndex) => prevPageIndex + 1);
            }
        }
    }, 300)

    useEffect(() => {
        const container = listContainerRef.current;
        if (container) {
            container.addEventListener('scroll', handleScroll)
        }
        return () => {
            if (container) {
                container.removeEventListener('scroll', handleScroll)
            }
        }
    }, [isLoadingMore])


    useEffect(() => {
        if(!allDataLoaded) {
            fetchAllAgents()
        }
    },[isConnected, pageIndex, allDataLoaded])

    // useEffect(() => {
    //     if(myAgents && !!myAgents.length) {
       
    //         setItems(myAgents)
    //     }else {
    //         setItems(itemsRef.current)
    //     }
    // },[myAgents])

   

    const fetchAllAgents = async() => {
        pageIndex === 1 ? setLoading(true) : setLoadingMore(true)
        const { agents, page, page_size, total } = await allAgentsApi(pageIndex, pageSize)
        pageIndex === 1 ? setLoading(false) : setLoadingMore(false)

        if( (page - 1) * page_size > total) {
            setAllDataLoaded(true)
        }else {
            setAllDataLoaded(false)
          
            const newAgents = [...items, ...agents]

            setItems(newAgents)
            itemsRef.current = newAgents
        }
    }

    useEffect(() => {
        battleLogsWS()
        const wsManager = new WebSocketManager(`${import.meta.env.VITE_API_URL}/ws/agents`, {
          heartbeatInterval: 10000,
          reconnectInterval: 3000,  
          maxReconnectAttempts: 5,  
          onOpen: () => {
            
          },
          onClose: (event: CloseEvent) => {
            
          },
          onMessage: (message: any) => {
            if(message) {
                const newList = [JSON.parse(message), ...itemsRef.current]
                setItems(newList)           
            }
          },
          onError: (error: ErrorEvent) => {
            
          },
        })    
        return () => {
          wsManager.close()
          battleWSRef.current.close()
        }
    }, [])
   
    const onExpand = async(id: number, detail: iBattlesReturn) => {
        setSelectedItem(id);     
        setItems((prevItems) => {
            const newItems = prevItems.filter((item) => item.id !== id);
            const selectedItem = prevItems.find((item) => item.id === id);
            return selectedItem ? [selectedItem, ...newItems] : newItems;
        })
        
        setBattleDetail(detail)
    }
    
    const onCollapse = () => {
        setSelectedItem(null)
        setItems(itemsRef.current)
     
        setBattleDetail(null)
    }

    const battleLogsWS = async() => {
        const battleWS = new WebSocketManager(`${import.meta.env.VITE_API_URL}/ws/battle`, {
            heartbeatInterval: 10000,
            reconnectInterval: 3000,  
            maxReconnectAttempts: 5,  
            onOpen: () => {
                // console.log('battle ws onOpen')
            },
            onClose: (event: CloseEvent) => {
                // console.log('battle ws onClose')
            },
            onMessage: (message: any) => {
                if(message) {
                    // console.log('battle ws onMessage', JSON.parse(message))
                    const { battle, type } = JSON.parse(message)
                    if(type === "BATTLE_RESULT") {
                        dispatch(lastBattleLogAction(battle))
                        dispatch(lastBattleLogTableAction(battle))
                    }
                }   
            },
            onError: (error: ErrorEvent) => {
                // console.log('battle ws onError', error)
            },
          })
          battleWSRef.current = battleWS
    }

    return (    
        <Box>
            {
                myAgents && !!myAgents.length && 
                <BackButton
                    onClick={() => dispatch(myAgentsAction([]))}
                    style={{
                        marginTop: '-38px'
                    }}
                />
            }
            <Box h="750px" className="">
                {
                    false ? 
                    <Box h="750px" className="center w100" w={AGENT_ITEM_MAX_WIDTH}>
                        <Spinner size="md" color="white" h="32px" w="32px"/>
                    </Box>
                    :
                    <Box maxH="750px" overflowY="scroll" className="" ref={listContainerRef}>
                        <TransitionGroup className=''>
                            {(!!myAgents.length ? myAgents : items).map((item, index) => (
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
