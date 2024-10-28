import React from 'react';
import { Card, Table, Tag, Tooltip } from 'antd';
import { CheckCircleOutlined, SyncOutlined, ClockCircleOutlined } from '@ant-design/icons';
import type { ColumnsType } from 'antd/es/table';
import type { AntdIconProps } from '@ant-design/icons/lib/components/AntdIcon';
//交易历史组件
interface Transaction {
    id: string;
    type: 'buy' | 'sell';
    pair: string;
    amount: string;
    price: string;
    total: string;
    status: 'completed' | 'pending' | 'processing';
    time: string;
    hash: string;
}

const getStatusTag = (status: Transaction['status']) => {
    switch (status) {
        case 'completed':
            return (
                <Tag icon={<CheckCircleOutlined onPointerEnterCapture={undefined}
                                                onPointerLeaveCapture={undefined} {...({} as AntdIconProps)} />} color="success">
                    已完成
                </Tag>
            );
        case 'processing':
            return (
                <Tag icon={<SyncOutlined onPointerEnterCapture={undefined}
                                         onPointerLeaveCapture={undefined} {...({} as AntdIconProps)} spin />} color="processing">
                    处理中
                </Tag>
            );
        case 'pending':
            return (
                <Tag icon={<ClockCircleOutlined onPointerEnterCapture={undefined}
                                                onPointerLeaveCapture={undefined} {...({} as AntdIconProps)} />} color="warning">
                    待确认
                </Tag>
            );
        default:
            return null;
    }
};

const columns: ColumnsType<Transaction> = [
    {
        title: '时间',
        dataIndex: 'time',
        key: 'time',
    },
    {
        title: '类型',
        dataIndex: 'type',
        key: 'type',
        render: (type: string) => (
            <Tag color={type === 'buy' ? 'green' : 'red'}>
                {type === 'buy' ? '买入' : '卖出'}
            </Tag>
        ),
    },
    {
        title: '交易对',
        dataIndex: 'pair',
        key: 'pair',
    },
    {
        title: '数量',
        dataIndex: 'amount',
        key: 'amount',
        align: 'right',
    },
    {
        title: '价格',
        dataIndex: 'price',
        key: 'price',
        align: 'right',
    },
    {
        title: '总额',
        dataIndex: 'total',
        key: 'total',
        align: 'right',
    },
    {
        title: '状态',
        dataIndex: 'status',
        key: 'status',
        render: getStatusTag,
    },
    {
        title: '交易哈希',
        dataIndex: 'hash',
        key: 'hash',
        render: (hash: string) => (
            <Tooltip title="点击查看交易详情">
                <a href={`https://sepolia.etherscan.io/tx/${hash}`} target="_blank" rel="noopener noreferrer">
                    {`${hash.substring(0, 6)}...${hash.substring(hash.length - 4)}`}
                </a>
            </Tooltip>
        ),
    },
];

const TransactionHistory: React.FC<{ transactions: Transaction[] }> = ({ transactions }) => {
    return (
        <Card title="交易历史" className="transaction-history">
            <Table <Transaction>
                columns={columns}
                dataSource={transactions}
                rowKey="id"
                pagination={{ pageSize: 10 }}
            />
        </Card>
    );
};

export default TransactionHistory;