


import React, { useState, useEffect,useRef } from "react";
import { Spinner, Image, Text, Box, Grid, GridItem } from "@chakra-ui/react";
import { ClickButtonWrapper, GeneralButton, ExpandTable, BackButton } from '@/components'
import { Kline1Img, MaskImg, ArrowImg, ArrowWhiteImg } from "@/assets/images"
import { WebSocketManager } from '@/ws/WebSocketManager'
import { formatNumber, formatTimeAgo } from '@/utils/tool'
import { CSSTransition, TransitionGroup } from 'react-transition-group';
import "./BattleLogs.css"; 
import { useAppDispatch, useAppSelector } from "@/redux/hooks";
import { expandBattleLogsAction, myAgentsAction, selectMyAgents, selectWalletInfo } from "@/redux/reducer";
import { allAgentsApi } from '@/api'
import { iAgentReturn } from "@/types";
import { AGENT_ITEM_MAX_WIDTH } from "@/config";
import { throttle } from 'lodash';

interface iAgent {
  
}

export const Agent: React.FC<iAgent> = ({

}) => {

    const [allDataLoaded, setAllDataLoaded] = useState(false)
    const [isLoading, setLoading] = useState(false)
    const [isLoadingMore, setLoadingMore] = useState(false)
    const [isOver, setOver] = useState(false)
    const [battleLogs, setBattleLogs] = useState([
        { time: 1734074073245, isWin: true, vs: 'AA' },
        { time: 1734074073245, isWin: false, vs: 'BB' },
        { time: 1734074073245, isWin: true, vs: 'CC' },
    ])

    const [items, setItems] = useState<iAgentReturn[]>([]);
    const [selectedItem, setSelectedItem] = useState<number | null>(null) 
    const [pageIndex, setPageIndex] = useState<number>(1)
    const [pageSize, setPageSize] = useState<number>(8)

    const { isConnected, address } = useAppSelector(selectWalletInfo)
    const itemsRef = useRef<iAgentReturn[]>([])
    const listContainerRef = useRef<HTMLDivElement>(null)


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
    }, 300); 

   

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

    useEffect(() => {
        if(myAgents && !!myAgents.length) {
            setItems(myAgents)
        }else {
            setItems(itemsRef.current)
        }
    },[myAgents])


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
        const interval = setInterval(() => {
          const t = Date.now();
          const newLog = {
            time: t,
            isWin: Math.random() > 0.5,
            vs: `${Math.random()}`.substring(0,6),
          };    
          setBattleLogs((prevLogs) => [newLog, ...prevLogs.slice(0, 2)]);
        }, 5000);
    
        return () => clearInterval(interval);
    }, [])

    useEffect(() => {
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
                const newList = [message, ...itemsRef.current]
                setItems(newList)           
            }
          },
          onError: (error: ErrorEvent) => {
            
          },
        });
    
      
        return () => {
          wsManager.close();
        };
    }, [])

    const onBuy = (address: string) => {
        if(address) {
            window.open(`https://pump.fun/coin/${address}`, '_blank')
        }
    }

   
    const onExpand = (idx: number) => {
        
        setSelectedItem(idx);     
        setItems((prevItems) => {
            const newItems = prevItems.filter((item) => item.id !== idx);
            const selectedItem = prevItems.find((item) => item.id === idx);
            return selectedItem ? [selectedItem, ...newItems] : newItems;
        })
        
    }

    const onCollapse = () => {
        setSelectedItem(null)
        setItems(itemsRef.current)
    }

    const createTimeAgo = (t: string) => {
        if(t) {
            return formatTimeAgo(new Date(t).getTime())
        }
        return '--'
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
                    isLoading ? 
                    <Box h="750px" className="center w100" w={AGENT_ITEM_MAX_WIDTH}>
                        <Spinner size="md" color="white" h="32px" w="32px"/>
                    </Box>
                    :
                    <Box maxH="750px" overflowY="scroll" className="" ref={listContainerRef}>
                        <TransitionGroup className=''>
                            {items.map((item, index) => (
                                <CSSTransition key={item.id} classNames="fade" timeout={300}>
                                    <Box
                                        className=""
                                        mb="10px"
                                        style={{
                                            display: selectedItem === null ? 'flex' : (index === 0 ? "flex" : "none"),
                                            flexDirection: 'row'
                                        }}
                                    >
                                        <Image src={Kline1Img} w="7px" h="59px"/>
                                        <Box 
                                            className="fx-row ai-ct " 
                                            ml="10px" 
                                            p="10px"
                                            h="180px" 
                                            // minW="960px"
                                            w={AGENT_ITEM_MAX_WIDTH} 
                                            borderRadius="5px" 
                                            border="1px solid #01FDB2">
                                            <Image src={item.image_url} h="160px" w="160px" borderRadius="5px"/>
                                            <Box className="fx-row ai-ct jc-sb w100" ml="15px">
                                                {/* row1 */}
                                                <Box className="fx-col ai-start w100 " maxW="350px" >
                                                    <Text className="main fz24 fw700">{item.name}</Text>
                                                    <Text className="gray9 fz12" mt="5px">created by {item.token_address ? item.token_address.substring(0,6) : '--'} · { createTimeAgo(item.created_at) }</Text>
                                                    <Text className="white fz14" h="108px" overflowY="scroll" mt="10px">{item.description}</Text>
                                                </Box>
                                                {/* row2  maxW='200px' */}
                                                <Box className="fx-col " h="180px"> 
                                                    <Text className="main fz20 fw700" mb="5px">{item.ticker}</Text>
                                                    {
                                                        [
                                                            { title: 'market cap:', value: `${formatNumber(item.market_cap)}`},
                                                            { title: 'bonding curve progress:', value: 'coming soon'},
                                                            { title: 'last all-time-high:', value: '1h ago'}
                                                        ].map(item => (
                                                            <Box key={item.title} className='fx-row ai-ct'>
                                                                <Text className="fz14 white">{item.title}&nbsp;</Text>
                                                                <Text className="fz14 main">{item.value}</Text>
                                                            </Box>
                                                        ))
                                                    }
                                                    <GeneralButton onClick={() => onBuy(item.token_address)} title="buy" style={{ height: "35px", width: '100px', marginTop: '10px'}}/>
                                                </Box>
                                                {/* row3 maxW="238px" */}
                                                <Box className="fx-row ">
                                                    <Box className="fx-col " h="180px">
                                                        <Text className="main fz24 fw700" mb="5px">battle stats</Text>
                                                        <Box className='fx-row ai-ct'>
                                                            <Text className="fz14 white">total battles:&nbsp;</Text>
                                                            <Text className="fz14 main">20</Text>
                                                        </Box>

                                                        <Box className='fx-row ai-ct'>
                                                            <Text className="fz14 white">win/loss:&nbsp;</Text>
                                                            <Text className="fz14 main">20
                                                                <span className="white">/</span>
                                                                <span className="red">5</span>
                                                                <span className="white">(</span>
                                                                    75.00%
                                                                <span className="white">)</span>
                                                            </Text>
                                                        </Box>  

                                                        <Text className="fz14 white">battle log:</Text>
                                                        <TransitionGroup component={null}>
                                                            {battleLogs.map((item) => (
                                                                <CSSTransition
                                                                    key={item.time}
                                                                    timeout={300}
                                                                    classNames="log"
                                                                >
                                                                    <Box className="fx-row ai-ct fz14" mb="6px">
                                                                        <Image src={Kline1Img} w="3px" h="14px" mr="5px" />
                                                                        <span className="gray9">{formatTimeAgo(item.time)}: </span>
                                                                        <span
                                                                            className="ml4 mr4"
                                                                            style={{ color: item.isWin ? "#01FDB2" : "#F45B5B" }}
                                                                            >
                                                                            {item.isWin ? "win" : "loss"}
                                                                        </span>
                                                                        <span className="white fw700 underline">vs {item.vs} </span>
                                                                        <Image
                                                                            src={MaskImg}
                                                                            w="17px"
                                                                            h="17px"
                                                                            borderRadius="2px"
                                                                            ml="5px"
                                                                        />
                                                                    </Box>
                                                                </CSSTransition>
                                                            ))}
                                                        </TransitionGroup>


                                                    </Box>

                                                    <Box 
                                                        mt="130px"
                                                        w="26px"
                                                        ml="25px" 
                                                        h="26px" 
                                                        className="click center" 
                                                        border="1px solid #01FDB2" 
                                                        borderRadius="6px"
                                                        _hover={{
                                                            backgroundColor: '#01553C',
                                                            border: '1px solid transparent',
                                                        }}
                                                        style={{
                                                            transform: selectedItem === null ? 'rotate(0deg)' : 'rotate(180deg)',
                                                            transition: 'transform 0.3s'
                                                        }}
                                                        onMouseOver={() => setOver(true)}
                                                        onMouseLeave={() => setOver(false)}
                                                        onClick={selectedItem ? () => onCollapse() : () => onExpand(item.id)}
                                                    >
                                                        <Image src={isOver ? ArrowWhiteImg : ArrowImg} h="7px" w='12px'/>
                                                    </Box>
                                                </Box>
                                            </Box>
                                        </Box>
                                    </Box>
                                </CSSTransition>
                            ))}
                        </TransitionGroup>

                        <CSSTransition
                            in={selectedItem !== null}
                            classNames="fade"
                            timeout={300}
                        >
                            <Box style={{ display: selectedItem ? "table" : "none" }} className="w100">
                                <ExpandTable onCollapse={onCollapse}/>
                            </Box>                
                        </CSSTransition>           
                    </Box>
                }
            </Box>
        </Box>
    );
};
