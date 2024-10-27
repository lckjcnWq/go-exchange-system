import {WagmiConfig, createConfig, configureChains} from 'wagmi';
import {sepolia} from 'wagmi/chains';
import {publicProvider} from 'wagmi/providers/public';
import {MetaMaskConnector} from 'wagmi/connectors/metaMask';
import TradePage from './pages/trade';
import AccountPage from './pages/account';
import './styles/layout.css';
import {BrowserRouter, Navigate, Route, Routes} from "react-router-dom";
import MainLayout from "./components/Layout/MainLayout.tsx";
// 配置链和提供者
const {chains, publicClient} = configureChains(
    [sepolia],
    [
        publicProvider()
    ]
);

// 创建 wagmi 配置
const config = createConfig({
    autoConnect: true,
    connectors: [
        new MetaMaskConnector({chains})
    ],
    publicClient,
});


function App() {
    return (
        <WagmiConfig config={config}>
            <BrowserRouter>
                <Routes>
                    <Route path="/" element={<MainLayout/>}>
                        <Route index element={<Navigate to="/trade" replace/>}/>
                        <Route path="trade" element={<TradePage/>}/>
                        <Route path="account" element={<AccountPage/>}/>
                    </Route>
                </Routes>
            </BrowserRouter>
        </WagmiConfig>
    );
}

export default App
