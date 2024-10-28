import React from 'react';
import { Row, Col } from 'antd';
import WalletStatus from '../../components/Account/WalletStatus';
import AssetsOverview from '../../components/Account/AssetsOverview';
import TransactionHistory from '../../components/Account/TransactionHistory';
import { useWallet } from '../../hooks/web3/useWallet';
// 引入 Asset 类型
const AccountPage: React.FC = () => {
    const { isConnected } = useWallet();

    // 模拟数据
    const assets = [
        {
            token: 'ETH',
            balance: '1.2345',
            value: '2654.32',
            price: '2145.67',
            change24h: 2.34,
        },
        {
            token: 'USDT',
            balance: '1000.00',
            value: '1000.00',
            price: '1.00',
            change24h: 0,
        },
    ];

    const transactions = [
        {
            id: '1',
            type: 'buy',
            pair: 'ETH/USDT',
            amount: '0.5 ETH',
            price: '2145.67',
            total: '1072.84',
            status: 'completed',
            time: '2024-03-15 14:30:00',
            hash: '0x1234567890abcdef1234567890abcdef12345678',
        },
        // 可以添加更多交易记录
    ];

    return (
        <div className="account-page">
            <Row gutter={[16, 16]}>
                <Col span={24}>
                    <WalletStatus />
                </Col>
                {isConnected && (
                    <>
                        <Col span={24}>
                            <AssetsOverview
                                totalValue="3654.32"
                                assets={assets}
                            />
                        </Col>
                        <Col span={24}>
                            <TransactionHistory
                                transactions={transactions as any}
                            />
                        </Col>
                    </>
                )}
            </Row>
        </div>
    );
};

export default AccountPage;