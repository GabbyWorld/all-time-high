import type { FC } from "react"
import React, {useState,  } from "react"
import { Box, Text, Image } from "@chakra-ui/react"
import { ArrowLeftImg, Kline3Img } from '@/assets/images'
import { GeneralButton, Notification, BackButton} from "@/components"
import { cteateAgentApi } from "@/api"
import { useAppDispatch, useAppSelector } from "@/redux/hooks"
import { notificationInfoAction, selectWalletInfo } from "@/redux/reducer"
import { useNavigate } from "react-router-dom"

interface iInput {
    value: string
    maxLen: number
    msg: string
    disable: boolean
}

export const CreateAgent: FC = () => {  
    const [loading, setLoading] = useState(false)

    const [name, setName] = useState<iInput>({
        value: '',
        maxLen: 15,
        msg: '',
        disable: true
    })
    const [ticker, setTicker] = useState<iInput>({
        value: '',
        maxLen: 10,
        msg: '',
        disable: true
    })
    const [prompt, setPrompt] = useState<iInput>({
        value: '',
        maxLen: 50,
        msg: '',
        disable: true
    })

    const [modal, setModal] = useState<{
        open: boolean,
        title: string
    }>({
        open: false,
 
        title: ''
    })

    const dispatch = useAppDispatch()
    const { isConnected, address } = useAppSelector(selectWalletInfo)
    const navigate = useNavigate()

    const onGenerate = async() => {     
       
        if(isConnected) {
            setLoading(true)
            const a = await cteateAgentApi(name.value, ticker.value, prompt.value)
            setLoading(false)
            if(a && a.id) {
                setModal({
                    open: true,
                    title: "agent coin created"
                })                
            }else {
                dispatch(notificationInfoAction({
                    open: true,
                    title: 'creation failed, please try again later'
                }))
            }
        } else {
            dispatch(notificationInfoAction({
                open: true,
                title: 'please connect wallet first'
            }))
        }
    }
    
    const onChangeName = (e: React.ChangeEvent<HTMLInputElement>) => {
        const val = e.target.value

        // Spaces, special characters restrictions

        if(val.length === 0) {
            return setName({...name, value: '', msg: '', disable: true})
        }
        if(val.length > name.maxLen) {
           return setName({...name, msg: `${name.maxLen} characters max`, value: val, disable: true})
        }
        setName({...name, value: val, msg: '', disable: false})
    }
    const onChangeTicker = (e: React.ChangeEvent<HTMLInputElement>) => {
        const val = e.target.value
        if(val.length === 0) {
            return setTicker({...ticker, value: '', msg: '', disable: true})
        }
        if(val.length > ticker.maxLen) {
           return setTicker({...ticker, msg: `${ticker.maxLen} characters max`, value: val, disable: true})
        }
        setTicker({...ticker, value: val.toLocaleUpperCase(), msg: '', disable: false})

    }
    const onChangePrompt = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
        const val = e.target.value
        if(val.length === 0) {
            return setPrompt({...prompt, value: '', msg: '', disable: true})
        }
        if(val.length > prompt.maxLen) {
           return setPrompt({...prompt, msg: `${prompt.maxLen} characters max`, value: val, disable: true})
        }
        setPrompt({...prompt, value: val, msg: '', disable: false})

    }   
    
   
    return (
        <Box pos='relative' className="center"  w="1510px"  overflowY="scroll">
            <BackButton
                onClick={() => navigate(-1)}
                style={{
                    position: 'absolute',
                    top: 0,
                    left: `calc(50% - 528px)`
                }}
            />

            <Box className="fx-col ai-ct " w='528px'> 
                <CreateInput 
                    title="name"
                    maxLen={name.maxLen}
                    currentLen={name.value.length}
                >
                    <input value={name.value} className="agent_input" onChange={onChangeName}/>
                </CreateInput>
                
                <CreateInput title="ticker"  maxLen={ticker.maxLen} currentLen={ticker.value.length}>
                    <input className="agent_input"  value={ticker.value.toLocaleUpperCase()}  onChange={onChangeTicker}/>
                </CreateInput>
                <CreateInput title="prompt"  maxLen={prompt.maxLen} currentLen={prompt.value.length}>
                    <textarea className="agent_input" value={prompt.value} style={{ minHeight: '100px', }}  onChange={onChangePrompt}/>
                </CreateInput>
                <GeneralButton 
                    disable={name.disable || ticker.disable || prompt.disable}
                    loading={loading}
                    style={{ width: '528px', height: '50px', marginTop: '86px' }}
                    title="generate agent & create coin"  // generate agent & create coin for 0.05 SOL
                    onClick={onGenerate}/>
            </Box>

            <Notification 
                visible={modal.open}
                onClose={() => {
                    setModal({ open: false, title: '' })
                    navigate('/')
                }}
                title={modal.title}
            />
        </Box>
    )
}

interface iCreateInput {
    title: string
    children: React.ReactNode
    maxLen: number
    currentLen: number
}
const CreateInput:FC<iCreateInput> = ({
    title,
    children,
    maxLen,
    currentLen
}) => {
    return (
        <Box mb="26px" className=" fx-col ai-start"> 
            <Box className="fx-row ai-ct">
                <Image src={Kline3Img} h="49px" w="17px"/>
                <Text className="fz32 main fw700" ml="10px">{title}</Text>
            </Box>
            { children }
            <Box className="fx-row ai-ct jc-sb w100" h="14px" >
                <Box />
                <Box  >
                    { currentLen > maxLen && <span className="fz12 red mr10">{maxLen} characters max</span>}
                    <span className="fz12" style={{ color: currentLen < maxLen ? '#999' : (currentLen === maxLen ? '#01FDB2' : '#F45B5B') }}>{currentLen}</span>
                    <span className="fz12 gray9">/{maxLen}</span>
                </Box>
            </Box>
        </Box>
    )
}