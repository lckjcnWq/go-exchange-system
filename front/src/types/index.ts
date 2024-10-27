export interface WalletInfo{
    address:string,
    chainId:number,
    isConnected:boolean,
}

export interface Web3Error{
    code:number,
    message:string
}