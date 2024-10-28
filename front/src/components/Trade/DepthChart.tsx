import React from 'react';
import { Card } from 'antd';
import { AreaChart, Area, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer } from 'recharts';
//交易深度图组件
interface DepthData {
    price: number;
    bids: number;
    asks: number;
}

interface DepthChartProps {
    data: DepthData[];
    title?: string;
}

const DepthChart: React.FC<DepthChartProps> = ({ data, title = 'Market Depth' }) => {
    return (
        <Card title={title} className="chart-card">
            <div style={{ width: '100%', height: 300 }}>
                <ResponsiveContainer>
                    <AreaChart
                        data={data}
                        margin={{ top: 10, right: 30, left: 0, bottom: 0 }}
                    >
                        <defs>
                            <linearGradient id="colorBids" x1="0" y1="0" x2="0" y2="1">
                                <stop offset="5%" stopColor="#52c41a" stopOpacity={0.8}/>
                                <stop offset="95%" stopColor="#52c41a" stopOpacity={0}/>
                            </linearGradient>
                            <linearGradient id="colorAsks" x1="0" y1="0" x2="0" y2="1">
                                <stop offset="5%" stopColor="#f5222d" stopOpacity={0.8}/>
                                <stop offset="95%" stopColor="#f5222d" stopOpacity={0}/>
                            </linearGradient>
                        </defs>
                        <XAxis
                            dataKey="price"
                            axisLine={false}
                            tickLine={false}
                        />
                        <YAxis
                            axisLine={false}
                            tickLine={false}
                        />
                        <CartesianGrid strokeDasharray="3 3" />
                        <Tooltip />
                        <Area
                            type="monotone"
                            dataKey="bids"
                            stroke="#52c41a"
                            fillOpacity={1}
                            fill="url(#colorBids)"
                        />
                        <Area
                            type="monotone"
                            dataKey="asks"
                            stroke="#f5222d"
                            fillOpacity={1}
                            fill="url(#colorAsks)"
                        />
                    </AreaChart>
                </ResponsiveContainer>
            </div>
        </Card>
    );
};

export default DepthChart;