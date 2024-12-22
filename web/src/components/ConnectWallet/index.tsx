import React, { useEffect, useState } from "react";
import { Flex, Image, Text, Box } from "@chakra-ui/react";
import { ClickButtonWrapper, GeneralButton, showToast } from '@/components'
import { CopyImg, WalletImg } from "@/assets/images"
import { connectWalletApi, profileApi } from '@/api'
import { useWallet, } from '@solana/wallet-adapter-react';
import { CopyToClipboard } from "react-copy-to-clipboard"
import { Connection, PublicKey } from "@solana/web3.js";
import { keepDecimals } from '@/utils/math'
import { useAppDispatch } from "@/redux/hooks";
import { notificationInfoAction, walletInfoAction } from "@/redux/reducer";

interface iConnectWallet {
  
}

export const ConnectWallet: React.FC<iConnectWallet> = ({

}) => {
    const [balance, setBalance] = useState<number>(0)
    const { publicKey, connect, disconnect, connected, select, signTransaction, signMessage,} = useWallet()  

    const dispatch = useAppDispatch()

    useEffect(() => {       
        if(publicKey && publicKey.toBase58()) {
            const address = publicKey.toBase58()
            fetchSOLBalance(address)
            login(address)
        }else {
            const token = localStorage.getItem('Authorization')
            if(!!!token) {
                dispatch(notificationInfoAction({
                    open: true,
                    title: 'please connect wallet first'
                }))
            }
        }
    },[publicKey])

    
    const login = async(address: string) => {
        const { message, user, token } = await connectWalletApi(address)
        if(message === "User connected") {
            dispatch(walletInfoAction({ isConnected: true, address }))
            // showToast('user connected','success')
            localStorage.setItem('Authorization', token)
        }
    }
    const fetchSOLBalance = async(address: string) => {
        const connection = new Connection("https://mainnet.helius-rpc.com/?api-key=f77fbc1f-282a-4bd7-99e7-cad253f17a77");
        try {
            const _publicKey = new PublicKey(address);
            const balance = await connection.getBalance(_publicKey);
            setBalance(balance / 1e9)
        } catch (error) {
            // console.error("Error fetching balance:", error);
        }
    }
    
    const truncateAddress = (address: string) => {
        
        if (!address || address.length < 8) {
            return address
        }
        return address.substring(0,6)
       
    }
    
    const onConnect = async() => {

        // @ts-ignore
        select('Phantom')
        await connect()

    }

    const onCopy = () => {
       showToast('copied','success')
    }
    const onDisconnect = async() => {
        await disconnect()
        localStorage.removeItem('Authorization')
        dispatch(walletInfoAction({ isConnected: false, address: '' }))
    }

    return (
        <Box>
            {
                connected ? 
                <Box h="40px" bgColor='#434343' className="fx-row ai-ct jc-ct" px="5px" borderRadius="5px" w="565px">
                    <Text className="fz24 white" mr="10px" whiteSpace='nowrap'>({keepDecimals(balance)} SOL)</Text>
                    <Image src={WalletImg}  w="30px" h="30px" borderRadius="50%"/>
                    <Text className="fz24 white fw700" ml="10px">{ publicKey && truncateAddress(publicKey.toBase58()) }</Text>

                    <CopyToClipboard text={publicKey ? publicKey.toBase58() : ''} onCopy={onCopy}>
                       <ClickButtonWrapper onClick={() => null}>
                            <Box className="fx-row ai-ct click" mx="20px">
                                <Image src={CopyImg}  w="16.36px" h="15px"/>
                                <Text className="main fz24 fw700" ml="5px">copy</Text>
                            </Box>
                       </ClickButtonWrapper>

                    </CopyToClipboard>                   
                    <GeneralButton style={{ height: "30px" }} title="disconnect" onClick={onDisconnect}/>
                </Box> : 
                <GeneralButton title="connect wallet" onClick={onConnect}/>
            }
        </Box>
    );
};
