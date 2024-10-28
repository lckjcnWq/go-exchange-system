import React from 'react';
import {Select, Space, Typography} from 'antd';
import {SyncOutlined} from '@ant-design/icons';
import type {AntdIconProps} from '@ant-design/icons/lib/components/AntdIcon';

const {Option} = Select;
const {Text} = Typography;

interface TradingPair {
    symbol: string;
    baseToken: string;
    quoteToken: string;
    lastPrice: string;
    priceChange: number;
}

interface PairSelectorProps {
    pairs: TradingPair[];
    selectedPair?: string;
    onPairChange: (pair: string) => void;
}

//交易对选择组件
const PairSelector: React.FC<PairSelectorProps> = ({pairs, selectedPair, onPairChange}) => {
    return (
        <div className="pair-selector">
            <Space direction="vertical" size={'small'} style={{width: '100%'}}>
                <Select
                    value={selectedPair}
                    onChange={onPairChange}
                    className={'pair-select'}
                    placeholder={'选择交易对'}
                    optionLabelProp={'label'}
                >
                    {
                        pairs.map((pair) => (
                            <Option
                                key={pair.symbol}
                                value={pair.symbol}
                                label={pair.symbol}>
                                <div className="pair-option">
                                    <Space direction="vertical" size={0}>
                                        <Text strong>{pair.symbol}</Text>
                                        <Space size="small">
                                            <Text type="secondary">{pair.lastPrice}</Text>
                                            <Text
                                                type={pair.priceChange >= 0 ? "success" : "danger"}
                                            >
                                                {pair.priceChange >= 0 ? '+' : ''}{pair.priceChange}%
                                            </Text>
                                        </Space>
                                    </Space>
                                </div>
                            </Option>
                        ))
                    }
                </Select>
                <div className="pair-info">
                    <SyncOutlined onPointerEnterCapture={undefined}
                                  onPointerLeaveCapture={undefined} {...({} as AntdIconProps)} spin/> 实时更新中
                </div>
            </Space>
        </div>
    );
};

export default PairSelector;
