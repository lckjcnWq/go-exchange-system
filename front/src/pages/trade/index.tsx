import React, {useState} from 'react';
import {Row, Col} from 'antd';
import PairSelector from '../../components/Trade/PairSelector';
import PriceDisplay from '../../components/Trade/PriceDisplay';
import TradeForm from '../../components/Trade/TradeForm';
import KLineChart from '../../components/Trade/KLineChart';
import DepthChart from '../../components/Trade/DepthChart';
import PriceChart from '../../components/Trade/PriceChart';
import './trade.css'
import './charts.css'

const TradePage: React.FC = () => {
    // 模拟数据
    const klineData = Array.from({ length: 30 }, (_, i) => ({
        time: `2024-03-${i + 1}`,
        price: 2000 + Math.random() * 200,
        volume: Math.random() * 1000
    }));

    const depthData = Array.from({ length: 20 }, (_, i) => ({
        price: 2000 + i * 10,
        bids: Math.random() * 100,
        asks: Math.random() * 100
    }));

    const priceData = Array.from({ length: 50 }, (_, i) => ({
        time: `${i}:00`,
        price: 2000 + Math.random() * 50
    }));

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
            <Row gutter={[16, 16]}>
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
                <Col span={24}>
                    <KLineChart data={klineData} title="ETH/USDT K线图" />
                </Col>
                <Col span={12}>
                    <DepthChart data={depthData} title="市场深度" />
                </Col>
                <Col span={12}>
                    <PriceChart data={priceData} title="实时价格" />
                </Col>
                <Col span={12}>
                    <TradeForm
                        type="buy"
                        maxAmount="1000 USDT"
                        onSubmit={console.log}
                    />
                </Col>
                <Col span={12}>
                    <TradeForm
                        type="sell"
                        maxAmount="0.5 ETH"
                        onSubmit={console.log}
                    />
                </Col>
            </Row>
        </div>
    );
};

export default TradePage;