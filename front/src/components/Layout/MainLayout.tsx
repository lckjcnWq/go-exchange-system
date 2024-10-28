import React from 'react';
import { Layout, Menu } from 'antd';
import { Link, Outlet, useLocation } from 'react-router-dom';
import { TrademarkOutlined, UserOutlined } from '@ant-design/icons';
import WalletConnect from '../Web3/WalletConnect';
import './MainLayout.css'

const {Header,Content,Sider,Footer} = Layout


const MainLayout:React.FC = () =>{
    const location = useLocation();
    const menuItems = [
        {
            key: '/trade',
            icon: <TrademarkOutlined onPointerEnterCapture={undefined}
                                     onPointerLeaveCapture={undefined} />,
            label: <Link to="/trade">交易</Link>,
        },
        {
            key: '/account',
            icon: <UserOutlined onPointerEnterCapture={undefined}
                                onPointerLeaveCapture={undefined}/>,
            label: <Link to="/account">账户</Link>,
        },
    ]

    return (
        <Layout className="site-layout">
            <Header className="site-header">
                <div className="site-logo">
                    Web3 交易系统
                </div>
                <WalletConnect />
            </Header>

            <Layout className="site-main-layout">
                <Sider className="site-sider" width={200}>
                    <Menu
                        mode="inline"
                        selectedKeys={[location.pathname]}
                        items={menuItems}
                        className="site-menu"
                    />
                </Sider>

                <Layout className="site-inner-layout">
                    <Content className="site-content">
                        <div className="site-container">
                            <Outlet />
                        </div>
                    </Content>
                </Layout>
            </Layout>

            <Footer className="site-footer">
                Web3 Trading System ©2024 Created by Your Team
            </Footer>
        </Layout>
    );
}

export default MainLayout;