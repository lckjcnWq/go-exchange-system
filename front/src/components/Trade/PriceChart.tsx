import React from 'react';
import { Card } from 'antd';
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer } from 'recharts';
//实时价格变化组件
interface PriceData {
    time: string;
    price: number;
}

interface PriceChartProps {
    data: PriceData[];
    title?: string;
}

const PriceChart: React.FC<PriceChartProps> = ({ data, title = 'Real-time Price' }) => {
    return (
        <Card title={title} className="chart-card">
            <div style={{ width: '100%', height: 200 }}>
                <ResponsiveContainer>
                    <LineChart
                        data={data}
                        margin={{ top: 5, right: 30, left: 20, bottom: 5 }}
                    >
                        <CartesianGrid strokeDasharray="3 3" />
                        <XAxis
                            dataKey="time"
                            axisLine={false}
                            tickLine={false}
                        />
                        <YAxis
                            axisLine={false}
                            tickLine={false}
                            domain={['dataMin', 'dataMax']}
                        />
                        <Tooltip />
                        <Line
                            type="monotone"
                            dataKey="price"
                            stroke="#1890ff"
                            dot={false}
                        />
                    </LineChart>
                </ResponsiveContainer>
            </div>
        </Card>
    );
};

export default PriceChart;