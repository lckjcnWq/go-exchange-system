import React from 'react';
import { Card, Descriptions, Button, Space, message } from 'antd';
import { CopyOutlined, LinkOutlined } from '@ant-design/icons';
import type { AntdIconProps } from '@ant-design/icons/lib/components/AntdIcon';
import { useWallet } from '../../hooks/web3/useWallet';

const WalletStatus: React.FC = () => {
    const { address, chainId, isConnected } = useWallet();

    const copyAddress = () => {
        if (address) {
            navigator.clipboard.writeText(address);
            message.success('地址已复制到剪贴板');
        }
    };

    const viewOnExplorer = () => {
        if (address) {
            window.open(`https://sepolia.etherscan.io/address/${address}`, '_blank');
        }
    };

    if (!isConnected) {
        return (
            <Card title="钱包状态">
                <div className="wallet-not-connected">
                    请先连接钱包
                </div>
            </Card>
        );
    }

    return (
        <Card
            title="钱包状态"
            extra={
                <Space>
                    <Button
                        icon={<CopyOutlined onPointerEnterCapture={undefined}
                                            onPointerLeaveCapture={undefined} {...({} as AntdIconProps)} />}
                        onClick={copyAddress}
                    >
                        复制地址
                    </Button>
                    <Button
                        icon={<LinkOutlined onPointerEnterCapture={undefined}
                                            onPointerLeaveCapture={undefined} {...({} as AntdIconProps)} />}
                        onClick={viewOnExplorer}
                    >
                        在区块浏览器中查看
                    </Button>
                </Space>
            }
        >
            <Descriptions column={1}>
                <Descriptions.Item label="钱包地址">
                    {address}
                </Descriptions.Item>
                <Descriptions.Item label="连接网络">
                    {chainId === 11155111 ? 'Sepolia测试网' : '未知网络'}
                </Descriptions.Item>
                <Descriptions.Item label="连接状态">
                    <span style={{ color: '#52c41a' }}>已连接</span>
                </Descriptions.Item>
            </Descriptions>
        </Card>
    );
};

export default WalletStatus;