import React, { useState } from 'react';
import { Form, Input, Button, message } from 'antd';
import { useAccount } from 'wagmi';
import BigNumber from 'bignumber.js';
import { tradeApi } from '../../services/api';
import { WebSocketService } from '../../services/WebSocketService';

interface TradeFormProps {
    selectedPair: string;
    onSuccess?: (txHash: string) => void;
}

const TradeForm: React.FC<TradeFormProps> = ({ selectedPair, onSuccess }) => {
    const [form] = Form.useForm();
    const { address, isConnected } = useAccount();
    const [loading, setLoading] = useState(false);
    const [quotePrice, setQuotePrice] = useState<string>('0');

    // 监听输入金额变化，获取报价
    const handleAmountChange = async (value: string) => {
        if (!value || !selectedPair) return;

        try {
            const [tokenIn, tokenOut] = selectedPair.split('/');
            const { amountOut } = await tradeApi.getPrice(tokenIn, tokenOut, value);
            setQuotePrice(amountOut);

            // 设置最小获得量（考虑1%滑点）
            const minAmount = new BigNumber(amountOut)
                .multipliedBy(0.99)
                .toFixed(0);
            form.setFieldsValue({ amountOutMin: minAmount });
        } catch (error) {
            console.error('Failed to get quote:', error);
        }
    };

    // 提交交易
    const handleSubmit = async (values: any) => {
        if (!isConnected || !address) {
            message.error('Please connect wallet first');
            return;
        }

        setLoading(true);
        try {
            const [tokenIn, tokenOut] = selectedPair.split('/');
            const deadline = Math.floor(Date.now() / 1000) + 1200; // 20分钟后过期

            const response = await tradeApi.createTrade({
                tokenIn,
                tokenOut,
                amountIn: values.amountIn,
                amountOutMin: values.amountOutMin,
                deadline,
            });

            message.success('Trade submitted successfully');
            onSuccess?.(response.txHash);

            // 监听交易状态
            WebSocketService.getInstance().subscribe('trade_update', (update: any) => {
                if (update.txHash === response.txHash) {
                    if (update.status === 'confirmed') {
                        message.success('Trade confirmed!');
                    } else if (update.status === 'failed') {
                        message.error('Trade failed');
                    }
                }
            });
        } catch (error: any) {
            message.error(error.message || 'Failed to submit trade');
        } finally {
            setLoading(false);
        }
    };

    return (
        <Form
            form={form}
            layout="vertical"
            onFinish={handleSubmit}
        >
            <Form.Item
                label="Amount In"
                name="amountIn"
                rules={[{ required: true, message: 'Please input amount' }]}
            >
                <Input
                    suffix={selectedPair?.split('/')[0]}
                    onChange={(e) => handleAmountChange(e.target.value)}
                />
            </Form.Item>

            <Form.Item label="Expected Output">
                <Input
                    value={quotePrice}
                    suffix={selectedPair?.split('/')[1]}
                    disabled
                />
            </Form.Item>

            <Form.Item
                label="Minimum Output"
                name="amountOutMin"
                rules={[{ required: true, message: 'Please input minimum amount' }]}
            >
                <Input suffix={selectedPair?.split('/')[1]} />
            </Form.Item>

            <Form.Item>
                <Button
                    type="primary"
                    htmlType="submit"
                    loading={loading}
                    disabled={!isConnected}
                    block
                >
                    {!isConnected ? 'Connect Wallet First' : 'Swap'}
                </Button>
            </Form.Item>
        </Form>
    );
};

export default TradeForm;