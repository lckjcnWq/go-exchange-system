import React from 'react';
import { Card, Typography } from 'antd';
import { useWallet } from '../../hooks/web3/useWallet';

const { Title } = Typography;

const AccountPage: React.FC = () => {
    const { isConnected, address } = useWallet();

    return (
        <div className="account-page">
            <Title level={2}>账户中心</Title>
            <Card>
                {isConnected ? (
                    <p>当前账户: {address}</p>
                ) : (
                    <p>请先连接钱包</p>
                )}
            </Card>
        </div>
    );
};

export default AccountPage;