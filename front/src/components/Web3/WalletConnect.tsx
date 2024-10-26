import React from 'react';
import { Button, message } from 'antd';
import { WalletOutlined, DisconnectOutlined } from '@ant-design/icons';
import { useAccount, useConnect, useDisconnect, useNetwork } from 'wagmi';
import { InjectedConnector } from 'wagmi/connectors/injected';
import { CONFIG } from '../../config';
import {connect} from "@wagmi/core";

const WalletConnect:React.FC = () => {
    const {address, isConnected} = useAccount();
    const {} = useConnect({
        connector: new InjectedConnector(),
        onSuccess(data){
            message.success('钱包连接成功');
        },
        onError(error){
            message.error('钱包连接失败');
        }
    });
    const {disconnect} = useDisconnect();
    const {chain} = useNetwork();
    const isCorrectNetwork = chain?.id === CONFIG.CHAINS.SEPOLIA;

    return (
        <div className="wallet-connect">
            {!isConnected ? (
                <Button
                    type="primary"
                    icon={<WalletOutlined />}
                    onClick={() => connect()}
                >
                    连接钱包
                </Button>
            ) : (
                <div className="wallet-info">
          <span className={`network-status ${isCorrectNetwork ? 'correct' : 'wrong'}`}>
            {isCorrectNetwork ? 'Sepolia' : '请切换到Sepolia测试网'}
          </span>
                    <span className="address">{address}</span>
                    <Button
                        icon={<DisconnectOutlined />}
                        onClick={() => disconnect()}
                        size="small"
                    >
                        断开连接
                    </Button>
                </div>
            )}
        </div>
    );
}


export default WalletConnect;
