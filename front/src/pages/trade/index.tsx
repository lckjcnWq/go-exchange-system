import React from 'react';
import {Card, Row, Col, Typography} from 'antd';
import {useWallet} from '../../hooks/web3/useWallet';

const {Title} = Typography;

const TradePage: React.FC = () => {
    const {isConnected} = useWallet();

    return (
        <div className="trade-page">
            <Title level={2}>交易中心</Title>
            <Row gutter={[24,24]}>
                <Col span={16}>
                    <Card title={'交易区'}>
                        {isConnected ? (
                            <p>交易功能开发中...</p>
                        ) : (
                            <p>请先连接钱包</p>
                        )}
                    </Card>
                </Col>

                <Col span={8}>
                    <Card title="交易信息">
                        <p>账户信息展示区</p>
                    </Card>
                </Col>
            </Row>
        </div>
    );
};

export default TradePage;