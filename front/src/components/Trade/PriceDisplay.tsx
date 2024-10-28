import React from 'react';
import {Card, Statistic, Row, Col} from 'antd';
import {ArrowUpOutlined, ArrowDownOutlined} from '@ant-design/icons';
import type {AntdIconProps} from '@ant-design/icons/lib/components/AntdIcon';

interface PriceDisplayProps {
    lastPrice: string;
    priceChange: number;
    high24h: string;
    low24h: string;
    volume24h: string;
}
// 创建价格展示组件
const PriceDisplay: React.FC<PriceDisplayProps> = ({
                                                       lastPrice,
                                                       priceChange,
                                                       high24h,
                                                       low24h,
                                                       volume24h,
                                                   }) => {
    return (
        <Card className="price-display">
            <Row gutter={[16, 16]}>
                <Col span={6}>
                    <Statistic
                        title="最新价格"
                        value={lastPrice}
                        precision={2}
                        valueStyle={{color: priceChange >= 0 ? '#3f8600' : '#cf1322'}}
                        prefix={priceChange >= 0 ?
                            <ArrowUpOutlined onPointerEnterCapture={undefined}
                                             onPointerLeaveCapture={undefined} {...({} as AntdIconProps)} /> :
                            <ArrowDownOutlined onPointerEnterCapture={undefined}
                                               onPointerLeaveCapture={undefined} {...({} as AntdIconProps)} />
                        }
                        suffix="USDT"
                    />
                </Col>
                <Col span={6}>
                    <Statistic
                        title="24h涨跌"
                        value={priceChange}
                        precision={2}
                        valueStyle={{color: priceChange >= 0 ? '#3f8600' : '#cf1322'}}
                        suffix="%"
                    />
                </Col>
                <Col span={6}>
                    <Statistic
                        title="24h最高"
                        value={high24h}
                        precision={2}
                        suffix="USDT"
                    />
                </Col>
                <Col span={6}>
                    <Statistic
                        title="24h最低"
                        value={low24h}
                        precision={2}
                        suffix="USDT"
                    />
                </Col>
            </Row>
            <Row style={{marginTop: 16}}>
                <Col span={24}>
                    <Statistic
                        title="24h成交量"
                        value={volume24h}
                        precision={2}
                        suffix="USDT"
                    />
                </Col>
            </Row>
        </Card>
    );
};

export default PriceDisplay;