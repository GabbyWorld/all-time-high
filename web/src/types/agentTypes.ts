
export interface iConnectWalletReturn {
    message: string,
    user: {
        id: number,
        wallet_address: string,
        created_at: string,
        updated_at: string
    },
    token: string
}


export interface iAgentReturn {
    id: number
    name:  string
    ticker: string
    prompt: string
    description: string
    image_url: string
    token_address: string
    created_at: string
    market_cap: number
    market_cap_updated_at: string
}

export enum iOUTCOME {
    TOTAL_VICTORY = "TOTAL_VICTORY",
    NARROW_VICTORY = "NARROW_VICTORY",
    CRUSHING_DEFEAT = "CRUSHING_DEFEAT",
    NARROW_DEFEAT = "NARROW_DEFEAT",
}



export interface iAttacker {
    id: number
    name: string
    ticker: string
    prompt: string
    description: string
    image_url: string
    token_address: string
    user_id: number
    created_at: string
    updated_at: string
    market_cap: number
    market_cap_updated_at: string
    highest_price: number
}
export interface iBattleItemReturn {
    id: number
    created_at: string
    outcome: string
    description: string
    
    attacker_id: number
    attacker: iAttacker
    defender_id: number
    defender: iAttacker
}
export interface iBattlesReturn {
    losses: number
    total: number
    win_rate: number
    wins: number
    battles: iBattleItemReturn[]
}


