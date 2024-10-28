import React, {useState} from 'react';
import {Card, Row, Col} from 'antd';
import PairSelector from '../../components/Trade/PairSelector';
import PriceDisplay from '../../components/Trade/PriceDisplay';
import TradeForm from '../../components/Trade/TradeForm';
import './trade.css'

const TradePage: React.FC = () => {
    // 模拟数据
    const tradingPairs = [
        {
            symbol: 'ETH/USDT',
            baseToken: 'ETH',
            quoteToken: 'USDT',
            lastPrice: '2145.67',
            priceChange: 2.34,
        },
        // 可以添加更多交易对
    ];

    const [selectedPair, setSelectedPair] = useState<string>('ETH/USDT');

    return (
        <div className="trade-page">
            <Row gutter={[24, 24]}>
                <Col span={24}>
                    <PairSelector
                        pairs={tradingPairs}
                        selectedPair={selectedPair}
                        onPairChange={setSelectedPair}
                    />
                </Col>
                <Col span={24}>
                    <PriceDisplay
                        lastPrice="2145.67"
                        priceChange={2.34}
                        high24h="2200.00"
                        low24h="2100.00"
                        volume24h="1234567.89"
                    />
                </Col>
                <Col span={12}>
                    <Card title="买入">
                        <TradeForm
                            type="buy"
                            maxAmount="1000 USDT"
                            onSubmit={console.log}
                        />
                    </Card>
                </Col>
                <Col span={12}>
                    <Card title="卖出">
                        <TradeForm
                            type="sell"
                            maxAmount="0.5 ETH"
                            onSubmit={console.log}
                        />
                    </Card>
                </Col>
            </Row>
        </div>
    );
};

export default TradePage;