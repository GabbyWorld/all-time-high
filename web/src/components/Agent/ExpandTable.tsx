import React, { useState, useEffect } from "react";
import { Flex, Image, Text, Box } from "@chakra-ui/react";
import { ClickButtonWrapper, GeneralButton } from '@/components'
import { Kline1Img, Kline3Img, CloseImg } from "@/assets/images"
import { formatTimeAgo } from "@/utils/tool";
import { CSSTransition, TransitionGroup } from 'react-transition-group';
import "./BattleLogs.css"; 

interface iExpandTable {
    onCollapse: () => void
}

export const ExpandTable: React.FC<iExpandTable> = ({
    onCollapse
}) => {
    const [list, setList] = useState([
        {time: Date.now(), UTC: '2024/12/07·12:34:56', attacker: 'attacker1', defender: 'defender1', winner: 'winner1', loser: 'loser1', result: 'result1', battleDescription: '111'},
    ])


    useEffect(() => {
        const interval = setInterval(() => {
          const newEntry = {
            time: Date.now(),
            UTC: new Date().toISOString().replace("T", "·").slice(0, -5),
            attacker: `attacker${Math.floor(Math.random() * 100)}`,
            defender: `defender${Math.floor(Math.random() * 100)}`,
            winner: `winner${Math.floor(Math.random() * 100)}`,
            loser: `loser${Math.floor(Math.random() * 100)}`,
            result: `result${Math.floor(Math.random() * 100)}`,
            battleDescription: `${Math.floor(Math.random() * 1000)}`,
          };
          setList((prevList) => [newEntry, ...prevList]);
        }, 1000);
        return () => clearInterval(interval);
    }, []);



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
                maxH="480px"
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
                        {list.map((item) => (
                        <CSSTransition key={item.time} timeout={300} classNames="row">
                            <tr className="tb_td_h">
                            <td className="tb_bd1 fx-row center tb_td_h">
                                <Image src={Kline3Img} h="20px" w="7px" />
                                <Text className="gray9 fz14" ml="10px">
                                {formatTimeAgo(item.time)}
                                </Text>
                            </td>
                            <td className="tb_bd1">
                                <Text className="gray9 fz14 center">{item.UTC}</Text>
                            </td>
                            <td className="tb_bd1">
                                <Box className="fx-row ai-ct jc-sb">
                                <Text className="white fz14 underline center" w="174px">
                                    {item.attacker}
                                </Text>
                                <Text className="white fz14 underline center" w="174px">
                                    {item.defender}
                                </Text>
                                </Box>
                            </td>
                            <td className="tb_bd1 w174" >
                                <Text className="white fz14 underline center">{item.winner}</Text>
                            </td>
                            <td className="tb_bd1 w174">
                                <Text className="white fz14 underline center">{item.loser}</Text>
                            </td>
                            <td className="tb_bd1 w174">
                                <Text className="main fz14 center">{item.result}</Text>
                            </td>
                            <td className="tb_bd1">
                                <Text className="white fz14 center">
                                {item.battleDescription}
                                </Text>
                            </td>
                            </tr>
                        </CSSTransition>
                        ))}
                    </TransitionGroup>
                </table>
            </Box>

        </Box>
    );
};


