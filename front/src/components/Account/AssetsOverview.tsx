import React from 'react';
import { Card, Row, Col, Statistic, Table } from 'antd';
import { AreaChartOutlined, DollarOutlined, SwapOutlined } from '@ant-design/icons';
import type { AntdIconProps } from '@ant-design/icons/lib/components/AntdIcon';
import type { ColumnsType } from 'antd/es/table';
//账户资产组件
export interface Asset{
    token:string;
    balance:string;
    value:string;
    price:string;
    change24h:number;
}

interface AssetsOverviewProps{
    totalValue:string;
    assets:Asset[];
}

const columns: ColumnsType<Asset> = [
    {
        title: '代币',
        dataIndex: 'token',
        key: 'token',
    },
    {
        title: '余额',
        dataIndex: 'balance',
        key: 'balance',
        align: 'right',
    },
    {
        title: '价值(USDT)',
        dataIndex: 'value',
        key: 'value',
        align: 'right',
    },
    {
        title: '价格',
        dataIndex: 'price',
        key: 'price',
        align: 'right',
    },
    {
        title: '24h涨跌',
        dataIndex: 'change24h',
        key: 'change24h',
        align: 'right',
        render: (change: number) => (
            <span style={{ color: change >= 0 ? '#52c41a' : '#f5222d' }}>
        {change >= 0 ? '+' : ''}{change}%
      </span>
        ),
    },
]


const AssetsOverview: React.FC<AssetsOverviewProps> = ({
                                                           totalValue,
                                                           assets,
                                                       }) => {
    return (
        <div className="assets-overview">
            <Row gutter={[16, 16]}>
                <Col span={8}>
                    <Card>
                        <Statistic
                            title="总资产估值(USDT)"
                            value={totalValue}
                            prefix={<DollarOutlined onPointerEnterCapture={undefined}
                                                    onPointerLeaveCapture={undefined} {...({} as AntdIconProps)} />}
                            precision={2}
                        />
                    </Card>
                </Col>
                <Col span={8}>
                    <Card>
                        <Statistic
                            title="今日收益"
                            value={214.50}
                            precision={2}
                            valueStyle={{color: '#3f8600'}}
                            prefix={<AreaChartOutlined onPointerEnterCapture={undefined}
                                                       onPointerLeaveCapture={undefined} {...({} as AntdIconProps)} />}
                            suffix="USDT"
                        />
                    </Card>
                </Col>
                <Col span={8}>
                    <Card>
                        <Statistic
                            title="待结算交易"
                            value={3}
                            prefix={<SwapOutlined onPointerEnterCapture={undefined}
                                                  onPointerLeaveCapture={undefined} {...({} as AntdIconProps)} />}
                        />
                    </Card>
                </Col>
            </Row>
            <Card title="资产列表" style={{marginTop: '16px'}}>
                <Table<Asset>
                    columns={columns}
                    dataSource={assets}
                    rowKey="token"
                    pagination={false}
                />
            </Card>
        </div>
    );
};
export default AssetsOverview;

