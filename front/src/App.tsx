import {WagmiConfig, createConfig, configureChains} from 'wagmi';
import {sepolia} from 'wagmi/chains';
import {publicProvider} from 'wagmi/providers/public';
import WalletConnect from './components/Web3/WalletConnect';
import { MetaMaskConnector } from 'wagmi/connectors/metaMask';

// 配置链和提供者
const { chains, publicClient } = configureChains(
    [sepolia],
    [
        publicProvider()
    ]
);

// 创建 wagmi 配置
const config = createConfig({
    autoConnect: true,
    connectors: [
        new MetaMaskConnector({ chains })
    ],
    publicClient,
});


function App() {
    return (
        <WagmiConfig config={config}>
            <div className="app" style={{
                minHeight: '100vh',
                display: 'flex',
                flexDirection: 'column',
                alignItems: 'center',
                padding: '20px',
                background: '#f0f2f5'  // 添加背景色
            }}>
                <h1 style={{
                    marginBottom: '30px',
                    color: '#1890ff',
                    fontSize: '24px',
                    fontWeight: 'bold'
                }}>
                    Web3 交易系统
                </h1>
                <WalletConnect />
            </div>
        </WagmiConfig>
    );
}

export default App
