import React from 'react';
import {Card} from 'antd';
import {AreaChart, Area, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer} from 'recharts';
//K线图组件
interface KLineData {
    time: string;
    price: number;
    volume: number;
}

interface KLineChartProps {
    data: KLineData[];
    title?: string;
}

const KLineChrat: React.FC<KLineChartProps> = ({data, title}) => {
    return (
        <Card title={title} className="chart-card">
            <div style={{width: '100%', height: 400}}>
                <ResponsiveContainer>
                    <AreaChart
                        data={data}
                        margin={{top: 10, right: 30, left: 0, bottom: 0}}
                    >
                        <defs>
                            <linearGradient id="colorPrice" x1="0" y1="0" x2="0" y2="1">
                                <stop offset="5%" stopColor="#1890ff" stopOpacity={0.8}/>
                                <stop offset="95%" stopColor="#1890ff" stopOpacity={0}/>
                            </linearGradient>
                        </defs>
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
                        <CartesianGrid strokeDasharray="3 3"/>
                        <Tooltip/>
                        <Area
                            type="monotone"
                            dataKey="price"
                            stroke="#1890ff"
                            fillOpacity={1}
                            fill="url(#colorPrice)"
                        />
                    </AreaChart>
                </ResponsiveContainer>
            </div>
        </Card>
    );
};

export default KLineChrat;