
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


