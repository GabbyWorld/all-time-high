
import React, { FC, useState, useEffect,useRef } from "react"
import { Image, Text, Box } from "@chakra-ui/react"
import { GeneralButton } from '@/components'
import { Kline1Img, ArrowImg, ArrowWhiteImg } from "@/assets/images"
import { formatNumber, createTimeAgo } from '@/utils/tool'
import { CSSTransition, TransitionGroup } from 'react-transition-group'
import { iAgentReturn, iBattleItemReturn, iBattlesReturn } from "@/types"
import { AGENT_ITEM_MAX_WIDTH } from "@/config"
import { battlesApi } from "@/api"
import "./BattleLogs.css"
import { useAppDispatch, useAppSelector } from "@/redux/hooks"
import { lastBattleLogAction, selectLastBattleLog } from "@/redux/reducer"


interface iAgentItem  {
    activeIdx: number
    selectedItem: number | null
    onExpand: (id: number, detail: iBattlesReturn) => void
    onCollapse: () => void
    isAutoExpand?: boolean
}
export const AgentItem: FC<iAgentItem & iAgentReturn> = ({ 
    activeIdx,
    selectedItem = null,
    onExpand,
    onCollapse,
    isAutoExpand,
    id,
    name,
    ticker,
    prompt,
    description,
    image_url,
    token_address,
    created_at, 
    market_cap,
    market_cap_updated_at
 }) => {
    const [isOver, setOver] = useState(false)

   
    const [battleStats, setBattleStats] = useState<iBattlesReturn>({
        losses: 0,
        total: 0,
        win_rate: 0,
        wins: 0,
        battles: []
    })
    const [allBattles, setAllBattles] = useState<iBattleItemReturn[]>([])    
    const lastBattleLog = useAppSelector(selectLastBattleLog)

    const dispatch = useAppDispatch()

    useEffect(() => {
        fetchLogs()
    },[])

    useEffect(() => {
        if(isAutoExpand) {
            handleExpand()
        }
    },[isAutoExpand, battleStats])

    useEffect(() => {
        if(lastBattleLog) {
            const { attacker_id, outcome } = lastBattleLog
            if(attacker_id === id) {
                const { losses, total, win_rate, wins, battles} = battleStats
                const isWin = outcome.includes('VICTORY')
                const _total = total + 1
                const _losses = !isWin ? losses + 1 : losses
                const _wins = isWin ? wins + 1 : wins

                setBattleStats({
                    losses: _losses,
                    total: _total,
                    win_rate: (_wins / _total) * 100,
                    wins: _wins,
                    battles: [lastBattleLog,...battleStats.battles.slice(0,2)]
                })
                dispatch(lastBattleLogAction(null))
            }
        }
    },[lastBattleLog, id, battleStats])

    const fetchLogs = async() => {
        const b = await battlesApi(id)
        if(b) {
            
            setBattleStats({
                ...b,
                battles: b.battles.slice(0,3)
            })
            setAllBattles(b.battles)
        }
    }

    const handleExpand = () => {
        onExpand(id, {
            losses: battleStats.losses,
            total: battleStats.total,
            win_rate: battleStats.win_rate,
            wins: battleStats.wins,
            battles: allBattles
        })
    }
    const onBuy = (address: string) => {
        if(address) {
            window.open(`https://pump.fun/coin/${address}`, '_blank')
        }
    }

    const toItemDetail = (id: number) => {
        window.open(`/agent-detail?id=${id}`,'_blank')
    }

    return (
        <CSSTransition classNames="fade" timeout={300}>
            <Box
                className=""
                mb="10px"
                style={{
                    display: selectedItem === null ? 'flex' : (activeIdx === 0 ? "flex" : "none"),
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
                    <Image src={image_url} h="160px" w="160px" borderRadius="5px"/>
                    <Box className="fx-row ai-ct jc-sb w100" ml="15px">
                        {/* row1 */}
                        <Box className="fx-col ai-start w100 " maxW="350px" >
                            <Text className="main fz24 fw700">{name}</Text>
                            <Text className="gray9 fz12" mt="5px">created by {token_address ? token_address.substring(0,6) : '--'} · { createTimeAgo(created_at) }</Text>
                            <Text className="white fz14" h="108px" overflowY="scroll" mt="10px">{description}</Text>
                                
                        </Box>
                        {/* row2  */}
                        <Box className="fx-col " h="180px"> 
                            <Text className="main fz20 fw700" mb="5px">{ticker}</Text>
                            {
                                [
                                    { title: 'market cap:', value: `${formatNumber(market_cap)}`},
                                    { title: 'bonding curve progress:', value: 'coming soon'},
                                    { title: 'last all-time-high:', value: !!battleStats.battles.length ? createTimeAgo(battleStats.battles[0].created_at) : '--'}
                                ].map(item => (
                                    <Box key={item.title} className='fx-row ai-ct'>
                                        <Text className="fz14 white">{item.title}&nbsp;</Text>
                                        <Text className="fz14 main">{item.value}</Text>
                                    </Box>
                                ))
                            }
                            <GeneralButton onClick={() => onBuy(token_address)} title="buy" style={{ height: "35px", width: '100px', marginTop: '10px'}}/>
                        </Box>
                        {/* row3 */}
                        <Box className="fx-row" w="330px" >
                            <Box className="fx-col"  w="256px" h="180px">
                                <Text className="main fz24 fw700" mb="5px">battle stats</Text>
                                <Box className='fx-row ai-ct'>
                                    <Text className="fz14 white">total battles:&nbsp;</Text>
                                    <Text className="fz14 main">{battleStats.total}</Text>
                                </Box>

                                <Box className='fx-row ai-ct'>
                                    <Text className="fz14 white">win/loss:&nbsp;</Text>
                                    <Text className="fz14 main">{battleStats.wins}
                                        <span className="white">/</span>
                                        <span className="red">{battleStats.losses}</span>
                                        <span className="white">(</span>
                                            {battleStats.win_rate.toFixed(2)}%
                                        <span className="white">)</span>
                                    </Text>
                                </Box>  

                                <Text className="fz14 white">battle log:</Text>
                                <TransitionGroup component={null}>
                                    {battleStats.battles.map((item: iBattleItemReturn) => {
                                     
                                        const {attacker_id,
                                                created_at,
                                                defender_id,
                                                attacker,
                                                defender,
                                                description,
                                                id,
                                                outcome} = item
                                        const isWin = outcome.includes('VICTORY') 
                                        
                                        return(
                                            <CSSTransition
                                                key={id}
                                                timeout={300}
                                                classNames="log "
                                            >
                                                <Box className="fx-row ai-ct fz14" whiteSpace="nowrap" mb="6px">
                                                    <Image src={Kline1Img} w="3px" h="14px" mr="5px" />
                                                    <span className="gray9">{createTimeAgo(created_at)}: </span>
                                                    <span
                                                        className="ml4 mr4"
                                                        style={{ color: isWin? "#01FDB2" : "#F45B5B" }}
                                                        >
                                                        {isWin ? "win" : "loss"}
                                                    </span>
                                                    <span className="white fw700">vs</span>&nbsp;
                                                    <span className="white fw700 underline click" onClick={() => toItemDetail(defender.id)}>{defender.name} </span>
                                                    <Image
                                                        src={defender.image_url}
                                                        w="17px"
                                                        h="17px"
                                                        borderRadius="2px"
                                                        ml="5px"
                                                    />
                                                </Box>
                                            </CSSTransition>
                                        )
                                    })}
                                </TransitionGroup>
                            </Box>
                            {
                                !!!isAutoExpand && 
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
                                    onClick={selectedItem ? onCollapse : handleExpand}
                                >
                                    <Image src={isOver ? ArrowWhiteImg : ArrowImg} h="7px" w='12px'/>
                                </Box>
                            }
                        </Box>
                    </Box>
                </Box>
            </Box>
        </CSSTransition>
    )
}
