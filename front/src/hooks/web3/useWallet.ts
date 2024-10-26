import { useAccount, useNetwork } from 'wagmi';
import { CONFIG } from '../../config';
import { WalletInfo } from '../../types';

export const UseWallet = ():WalletInfo =>{
    const { address, isConnected } = useAccount();
    const { chain } = useNetwork();

    return {
        address: address || '',
        chainId: chain?.id || 0,
        isConnected: isConnected && chain?.id === CONFIG.CHAINS.SEPOLIA,
    };
}