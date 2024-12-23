import React, { useEffect, useState } from "react"
import { Image, Text, Box } from "@chakra-ui/react"
import { TextPopover } from '@/components'
import {  Kline3Img, CloseImg } from "@/assets/images"
import { createTimeAgo, localDate2UTC } from "@/utils/tool"
import { CSSTransition, TransitionGroup } from 'react-transition-group'
import { iBattleItemReturn, iBattlesReturn } from "@/types"
import "./BattleLogs.css"
import { useAppDispatch, useAppSelector } from "@/redux/hooks"
import { lastBattleLogTableAction, selectLastBattleLogTable } from "@/redux/reducer"

interface iExpandTable {
    onCollapse: () => void
    detail: iBattlesReturn | null
    agentId: number | null
}

export const ExpandTable: React.FC<iExpandTable> = ({
    onCollapse,
    detail,
    agentId
}) => {
    
    if(detail === null) {
        return null 
    }
    const [battlesList, setBattlesList] = useState<iBattleItemReturn[]>([])
    const lastBattleLogTable = useAppSelector(selectLastBattleLogTable)

    const { losses, total, win_rate, wins, battles} = detail
    const dispatch = useAppDispatch()

    useEffect(() => {
        if(battles && !!battles.length) {
            setBattlesList(battles)
        }
    },[battles])

    useEffect(() => {
        if(lastBattleLogTable) {
            if(lastBattleLogTable.attacker_id === agentId) {
                setBattlesList(prev => [lastBattleLogTable, ...prev])
                dispatch(lastBattleLogTableAction(null))
            }
        }
    },[battlesList, lastBattleLogTable, agentId])

    
    const toItemDetail = (id: number) => {
        window.open(`/agent-detail?id=${id}`,'_blank')
    }
    
    return (    
        <Box mt="30px" maxW="1360px" ml="17px">            
            {/* <Box className="fx-row ai-ct jc-sb">
                <div/>
                <ClickButtonWrapper onClick={onCollapse}>
                    <Image src={CloseImg} w="25px" h="25px"/>
                </ClickButtonWrapper>
            </Box> */}
            <Box 
                border="1px solid #01FDB2"  
                borderRadius="5px" 
                mt="5px"
                maxH="530px"
                overflowY="scroll"
                className="scrollable"
            >
                <table className="tb_container ">
                    <thead>
                        <tr style={{ height: "40px", borderRadius: "10px 10px 0 0" }}>
                        <th className="tb_bd1">
                            <Text className="tb_header">time</Text>
                        </th>
                        <th className="tb_bd1">
                            <Text className="tb_header">UTC</Text>
                        </th>
                        <th className="">
                            <Box className="fx-row ai-ct jc-sb">
                            <Text className="tb_header center" w="174px">
                                attacker
                            </Text>
                            <Text className="tb_header center" w="174px">
                                defender
                            </Text>
                            </Box>
                        </th>
                        <th className="tb_bd1">
                            <Text className="tb_header">winner</Text>
                        </th>
                        <th className="tb_bd1">
                            <Text className="tb_header">loser</Text>
                        </th>
                        <th className="tb_bd1">
                            <Text className="tb_header">attack outcome</Text>
                        </th>
                        <th className="tb_bd1">
                            <Text className="tb_header">battle description</Text>
                        </th>
                        </tr>
                    </thead>
                    <TransitionGroup component="tbody">
                        {battlesList.map((item) => {
                           
                            const {attacker_id,
                                created_at,
                                defender_id,
                                attacker,
                                defender,
                                description,
                                id,
                                outcome} = item
                            const isWin = outcome.includes('VICTORY') 
                         
                        
                            const _outcome  = outcome.toLowerCase().replace('_', ' ');
                            return (
                                <CSSTransition key={id} timeout={300} classNames="row">
                                    <tr className="tb_td_h">
                                    <td className="tb_bd1 fx-row center tb_td_h">
                                        <Image src={Kline3Img} h="20px" w="7px" />
                                        <Text className="gray9 fz14" ml="10px">
                                        {createTimeAgo(created_at)}
                                        </Text>
                                    </td>
                                    <td className="tb_bd1">
                                        <Text className="gray9 fz14 center">{localDate2UTC(created_at)}</Text>
                                    </td>
                                    <td className="tb_bd1">
                                        <Box className="fx-row ai-ct jc-sb">
                                        <Text className="white fz14 underline center click" w="174px" onClick={() => toItemDetail(attacker.id)}>
                                            {attacker.name}
                                        </Text>
                                        <Text className="white fz14 underline center click" w="174px" onClick={() => toItemDetail(defender.id)}>
                                            {defender.name}
                                        </Text>
                                        </Box>
                                    </td>
                                    <td className="tb_bd1 w174" >
                                        <Text className="white fz14 underline center click" onClick={() => toItemDetail(isWin ? attacker.id : defender.id )}>{isWin ? attacker.name : defender.name }</Text>
                                    </td>
                                    <td className="tb_bd1 w174">
                                        <Text className="white fz14 underline center click" onClick={() => toItemDetail(isWin ? defender.id : attacker.id )}>{isWin ? defender.name : attacker.name}</Text>
                                    </td>
                                    <td className="tb_bd1 w174">
                                        <Text 
                                            className="fz14 center"
                                            color={ outcome.includes('VICTORY') ? "#01FDB2" : "#F45B5B" }
                                        >{_outcome}</Text>
                                    </td>
                                    <td className="tb_bd1">
                                    <TextPopover 
                                        content={
                                            <Text className="fz14 white" textAlign='start' dangerouslySetInnerHTML={{ __html: description }} />
                                        }
                                    >
                                        <Text className="white fz14 center">
                                            {`${description.substring(0,16)}...`}
                                        </Text>
                                    </TextPopover>
                                    </td>
                                    </tr>
                                </CSSTransition>
                            )
                        })}
                    </TransitionGroup>
                </table>
            </Box>

        </Box>
    );
};


