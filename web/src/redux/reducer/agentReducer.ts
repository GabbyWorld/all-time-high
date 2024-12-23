import { iAgentReturn, iBattleItemReturn } from '@/types'
import { PayloadAction, createSlice } from '@reduxjs/toolkit'
import { RootState } from '../store'


interface iAgentState {
  
  expandBattleLogs: iBattleItemReturn | null
  notificationInfo: {
    open: boolean
    title: string
  } 
  walletInfo: {
    isConnected: boolean,
    address: string
  },
  myAgents: iAgentReturn[]
  lastBattleLog: iBattleItemReturn | null 
  lastBattleLogTable: iBattleItemReturn | null 
}

const initialState: iAgentState = {
  
  expandBattleLogs: null,
  notificationInfo: {
    open: false,
    title: ''
  },
  walletInfo: {
    isConnected: false,
    address: ''
  },
  myAgents: [],
  lastBattleLog: null,
  lastBattleLogTable: null
}

const marketSlice = createSlice({
  name: 'agentReducer',
  initialState,
  reducers: {
   
    expandBattleLogsAction: (state, action: PayloadAction<iAgentState['expandBattleLogs']>) => {
      state.expandBattleLogs = action.payload
    },
    notificationInfoAction: (state, action: PayloadAction<iAgentState['notificationInfo']>) => {
      state.notificationInfo = action.payload
    },
    walletInfoAction: (state, action: PayloadAction<iAgentState['walletInfo']>) => {
      state.walletInfo = action.payload
    },
    myAgentsAction: (state, action: PayloadAction<iAgentState['myAgents']>) => {
      state.myAgents = action.payload
    },
    lastBattleLogAction: (state, action: PayloadAction<iAgentState['lastBattleLog']>) => {
      state.lastBattleLog = action.payload
    },
    lastBattleLogTableAction: (state, action: PayloadAction<iAgentState['lastBattleLogTable']>) => {
      state.lastBattleLogTable = action.payload
    },
    
  }
})

export const { 
  expandBattleLogsAction, 
  notificationInfoAction, 
  walletInfoAction, 
  myAgentsAction, 
  lastBattleLogAction,
  lastBattleLogTableAction
 } = marketSlice.actions


export const selectExpandBattleLogs = (state: RootState) => state.agentReducer.expandBattleLogs
export const selectNotificationInfo = (state: RootState) => state.agentReducer.notificationInfo
export const selectWalletInfo = (state: RootState) => state.agentReducer.walletInfo
export const selectMyAgents = (state: RootState) => state.agentReducer.myAgents
export const selectLastBattleLog = (state: RootState) => state.agentReducer.lastBattleLog
export const selectLastBattleLogTable = (state: RootState) => state.agentReducer.lastBattleLogTable

export const agentReducer = marketSlice.reducer
