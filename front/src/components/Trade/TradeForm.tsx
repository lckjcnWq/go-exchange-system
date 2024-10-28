import React from 'react';
import { Form, Input, Button, Radio, Typography, Space } from 'antd';
import { useWallet } from '../../hooks/web3/useWallet';

const { Text } = Typography;
interface TradeFormProps {
    type: 'buy' | 'sell';
    maxAmount?: string;
    onSubmit: (values: any) => void;
}
// 创建交易表单组件
const TradeForm: React.FC<TradeFormProps> = ({
                                                 type,
                                                 maxAmount,
                                                 onSubmit,
                                             }) => {
    const { isConnected } = useWallet();
    const [form] = Form.useForm();

    return (
        <Form
            form={form}
            layout="vertical"
            onFinish={onSubmit}
            className="trade-form"
        >
            <Form.Item label="交易类型">
                <Radio.Group value={type} buttonStyle="solid">
                    <Radio.Button value="buy">买入</Radio.Button>
                    <Radio.Button value="sell">卖出</Radio.Button>
                </Radio.Group>
            </Form.Item>

            <Form.Item
                label="价格"
                name="price"
                rules={[{ required: true, message: '请输入价格' }]}
            >
                <Input suffix="USDT" placeholder="输入价格" />
            </Form.Item>

            <Form.Item
                label="数量"
                name="amount"
                rules={[{ required: true, message: '请输入数量' }]}
            >
                <Input suffix="ETH" placeholder="输入数量" />
            </Form.Item>

            <Form.Item label="总额" name="total">
                <Input disabled suffix="USDT" placeholder="交易总额" />
            </Form.Item>

            {maxAmount && (
                <Form.Item>
                    <Space>
                        <Text type="secondary">可用余额：</Text>
                        <Text>{maxAmount}</Text>
                    </Space>
                </Form.Item>
            )}

            <Form.Item>
                <Button
                    type="primary"
                    htmlType="submit"
                    disabled={!isConnected}
                    block
                >
                    {!isConnected ? '请先连接钱包' : `确认${type === 'buy' ? '买入' : '卖出'}`}
                </Button>
            </Form.Item>
        </Form>
    );
};

export default TradeForm;