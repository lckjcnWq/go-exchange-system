import React from 'react';
import { Button, message } from 'antd';
import { WalletOutlined, DisconnectOutlined } from '@ant-design/icons';
import { useAccount, useConnect, useDisconnect, useNetwork } from 'wagmi';
import { InjectedConnector } from 'wagmi/connectors/injected';
import { CONFIG } from '../../config';


const WalletConnect:React.FC = () => {
    const {address, isConnected} = useAccount();
    const {connect} = useConnect({
        connector: new InjectedConnector(),
        onSuccess(_data){
            message.success('钱包连接成功');
        },
        onError(_error){
            message.error('钱包连接失败');
        }
    });
    const {disconnect} = useDisconnect();
    const {chain} = useNetwork();
    const isCorrectNetwork = chain?.id === CONFIG.CHAINS.SEPOLIA;

    const formatAddress = (addr:string)=>{
        return `${addr.slice(0,6)}...${addr.slice(-4)}`
    }

    return (
        <div className="wallet-connect" style={{
            display: 'flex',
            justifyContent: 'center',
            alignItems: 'center',
            padding: '20px',
            gap: '10px'
        }}>
            {!isConnected ? (
                <Button
                    type="primary"
                    icon={<WalletOutlined onPointerEnterCapture={undefined} onPointerLeaveCapture={undefined}  />}
                    onClick={() => connect()}
                    size="large"  // 调整按钮大小
                >
                    连接钱包
                </Button>
            ) : (
                <div className="wallet-info" style={{
                    display: 'flex',
                    alignItems: 'center',
                    gap: '10px',
                    padding: '8px 16px',
                    background: '#f5f5f5',
                    borderRadius: '8px',
                }}>
          <span className={`network-status ${isCorrectNetwork ? 'correct' : 'wrong'}`} style={{
              padding: '4px 8px',
              borderRadius: '4px',
              fontSize: '14px',
              backgroundColor: isCorrectNetwork ? '#f6ffed' : '#fff2f0',
              color: isCorrectNetwork ? '#52c41a' : '#ff4d4f',
          }}>
            {isCorrectNetwork ? 'Sepolia' : '请切换到Sepolia测试网'}
          </span>
                    <span className="address" style={{
                        fontFamily: 'monospace',
                        padding: '4px 8px',
                        background: '#fe4512',
                        borderRadius: '8px',
                        border: '1px solid #d9d9d9',
                    }}>
            {formatAddress(address!)}
          </span>
                    <Button
                        icon={<DisconnectOutlined onPointerEnterCapture={undefined} onPointerLeaveCapture={undefined}  />}
                        onClick={() => disconnect()}
                        size="middle"
                        danger  // 添加危险样式
                    >
                        断开连接
                    </Button>
                </div>
            )}
        </div>
    );
}


export default WalletConnect;
