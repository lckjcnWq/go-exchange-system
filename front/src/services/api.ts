import axios, { AxiosResponse } from 'axios';
import {message} from 'antd';
import {CONFIG} from '../config';

const api = axios.create({
    baseURL: CONFIG.API_BASE_URL,
    timeout: 10000,
});

// 请求拦截器
api.interceptors.request.use(
    (config) => {
        // 可以在这里添加token等认证信息
        return config;
    },
    (error) => {
        return Promise.reject(error);
    }
);

// 响应拦截器
api.interceptors.response.use(
    (response: { data: AxiosResponse<any, any> | Promise<AxiosResponse<any, any>>; }) => {
        return response.data;
    },
    (error) => {
        message.error(error.response?.data?.message || 'Network error');
        return Promise.reject(error);
    }
);

export interface CreateTradeParams {
    tokenIn: string;
    tokenOut: string;
    amountIn: string;
    amountOutMin: string;
    deadline: number;
}

export interface TradeInfo {
    txHash: string;
    tokenIn: string;
    tokenOut: string;
    amountIn: string;
    amountOut: string;
    status: string;
    confirmations: number;
    createdAt: string;
}

export const tradeApi = {
    // 创建交易
    createTrade: (params: CreateTradeParams) => {
        return api.post<any, { txHash: string }>('/trade', params);
    },

    // 获取交易列表
    getTrades: (page: number, size: number) => {
        return api.get<any, { list: TradeInfo[], total: number }>('/trades', {
            params: { page, size }
        });
    },

    // 获取代币价格
    getPrice: (tokenIn: string, tokenOut: string, amountIn: string) => {
        return api.get<any, { amountOut: string, rate: string }>('/price', {
            params: { tokenIn, tokenOut, amountIn }
        });
    }
};