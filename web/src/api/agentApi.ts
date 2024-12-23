import { iConnectWalletReturn, iAgentReturn, iBattleItemReturn, iBattlesReturn } from '@/types'
import { get, post } from './base'

export const connectWalletApi = ( wallet_address: string) => {
    return post<iConnectWalletReturn>(`/connect_wallet`,{ wallet_address })
}

export const profileApi = () => {
    return get(`/profile`,{  })
}

export const cteateAgentApi = (name: string, ticker: string, prompt: string) => {
    return post<iAgentReturn>(`/agent`,{ name,ticker,prompt })
}

export const userAgentsApi = () => {
    return get<{ agents: iAgentReturn[], page: number, page_size: number, total: number }>(`/agents`,{page: 1, page_size: 100 })
}
export const allAgentsApi = (page: number, page_size: number) => {
    return get<{ agents: iAgentReturn[], page: number, page_size: number, total: number }>(`/agents/all`,{ page, page_size })
}

export const agent1Api = (id: number) => {
    return get<iAgentReturn>(`/agent/${id}`)
}

export const battlesApi = (agent_id: number) => {
    return get<iBattlesReturn>(`/battles`,{ agent_id })
}
export const battleApi = (id: number) => {
    return get<iBattleItemReturn>(`/battle`,{ id })
}


