import { iAgentReturn } from '@/types'
import { PayloadAction, createSlice } from '@reduxjs/toolkit'
import { RootState } from '../store'


interface iHomeState {
  
  expandBattleLogs: {
    open: boolean
    index: number
  }
  notificationInfo: {
    open: boolean
    title: string
  } 
  walletInfo: {
    isConnected: boolean,
    address: string
  },
  myAgents: iAgentReturn[]
}

const initialState: iHomeState = {
  
  expandBattleLogs: {
    open: false,
    index: -1
  },
  notificationInfo: {
    open: false,
    title: ''
  },
  walletInfo: {
    isConnected: false,
    address: ''
  },
  myAgents: []
}

const marketSlice = createSlice({
  name: 'homeReducer',
  initialState,
  reducers: {
   
    expandBattleLogsAction: (state, action: PayloadAction<iHomeState['expandBattleLogs']>) => {
      state.expandBattleLogs = action.payload
    },
    notificationInfoAction: (state, action: PayloadAction<iHomeState['notificationInfo']>) => {
      state.notificationInfo = action.payload
    },
    walletInfoAction: (state, action: PayloadAction<iHomeState['walletInfo']>) => {
      state.walletInfo = action.payload
    },
    myAgentsAction: (state, action: PayloadAction<iHomeState['myAgents']>) => {
      state.myAgents = action.payload
    },
    
  }
})

export const { expandBattleLogsAction, notificationInfoAction, walletInfoAction, myAgentsAction } = marketSlice.actions


export const selectExpandBattleLogs = (state: RootState) => state.homeReducer.expandBattleLogs
export const selectNotificationInfo = (state: RootState) => state.homeReducer.notificationInfo
export const selectWalletInfo = (state: RootState) => state.homeReducer.walletInfo
export const selectMyAgents = (state: RootState) => state.homeReducer.myAgents

export const homeReducer = marketSlice.reducer
