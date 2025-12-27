import React from 'react'
import { StyleSheet, Dimensions } from 'react-native'
import { WebView } from 'react-native-webview'
import { dimeTheme } from '@/constants/theme'

interface CandleData {
    time: number
    open: number
    high: number
    low: number
    close: number
    volume: number
}

interface TradingViewChartProps {
    candles: CandleData[]
    symbol: string
    height?: number
    chartType?: 'candlestick' | 'line' | 'area'
}

export const TradingViewChart: React.FC<TradingViewChartProps> = ({
    candles,
    symbol,
    height = 300,
    chartType = 'candlestick'
}) => {
    // Convert Unix timestamps to TradingView format (YYYY-MM-DD)
    const formatData = () => {
        return candles.map(c => ({
            time: new Date(c.time * 1000).toISOString().split('T')[0],
            open: c.open,
            high: c.high,
            low: c.low,
            close: c.close,
            value: c.close, // For line/area charts
        }))
    }

    const chartData = JSON.stringify(formatData())

    // Determine colors based on price change
    const priceUp = candles.length > 1 && candles[candles.length - 1].close >= candles[0].close
    const upColor = dimeTheme.colors.profit
    const downColor = dimeTheme.colors.loss

    const html = `
<!DOCTYPE html>
<html>
<head>
    <meta name="viewport" content="width=device-width, initial-scale=1.0, maximum-scale=1.0">
    <script src="https://unpkg.com/lightweight-charts@4.1.0/dist/lightweight-charts.standalone.production.js"></script>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body { background: ${dimeTheme.colors.background}; overflow: hidden; }
        #chart { width: 100%; height: 100vh; }
    </style>
</head>
<body>
    <div id="chart"></div>
    <script>
        const chartOptions = {
            layout: {
                background: { type: 'solid', color: '${dimeTheme.colors.background}' },
                textColor: '${dimeTheme.colors.textSecondary}',
            },
            grid: {
                vertLines: { color: '${dimeTheme.colors.border}22' },
                horzLines: { color: '${dimeTheme.colors.border}22' },
            },
            crosshair: {
                mode: LightweightCharts.CrosshairMode.Magnet,
            },
            rightPriceScale: {
                borderColor: '${dimeTheme.colors.border}',
            },
            timeScale: {
                borderColor: '${dimeTheme.colors.border}',
                timeVisible: true,
                secondsVisible: false,
            },
            handleScale: { axisPressedMouseMove: { time: true, price: true } },
            handleScroll: { vertTouchDrag: false },
        };

        const chart = LightweightCharts.createChart(document.getElementById('chart'), chartOptions);
        
        const data = ${chartData};
        
        // Create chart based on type
        ${chartType === 'candlestick' ? `
        const series = chart.addCandlestickSeries({
            upColor: '${upColor}',
            downColor: '${downColor}',
            borderDownColor: '${downColor}',
            borderUpColor: '${upColor}',
            wickDownColor: '${downColor}',
            wickUpColor: '${upColor}',
        });
        series.setData(data);
        ` : chartType === 'line' ? `
        const series = chart.addLineSeries({
            color: '${priceUp ? upColor : downColor}',
            lineWidth: 2,
        });
        series.setData(data.map(d => ({ time: d.time, value: d.close })));
        ` : `
        const series = chart.addAreaSeries({
            lineColor: '${priceUp ? upColor : downColor}',
            topColor: '${priceUp ? upColor : downColor}44',
            bottomColor: '${priceUp ? upColor : downColor}00',
            lineWidth: 2,
        });
        series.setData(data.map(d => ({ time: d.time, value: d.close })));
        `}

        // Fit content
        chart.timeScale().fitContent();

        // Handle resize
        window.addEventListener('resize', () => {
            chart.applyOptions({ width: window.innerWidth });
        });

        // Add price line at current price
        const lastPrice = data[data.length - 1].close;
        series.createPriceLine({
            price: lastPrice,
            color: '${dimeTheme.colors.primary}',
            lineWidth: 1,
            lineStyle: 2,
            axisLabelVisible: true,
            title: 'Current',
        });
    </script>
    
    <!-- TradingView Attribution (required by Apache 2.0 license) -->
    <div style="position:absolute; bottom:4px; right:8px; opacity:0.5;">
        <a href="https://www.tradingview.com/" target="_blank" style="color:${dimeTheme.colors.textTertiary}; font-size:10px; text-decoration:none;">
            TradingView Lightweight Charts™
        </a>
    </div>
</body>
</html>
`

    return (
        <WebView
            source={{ html }}
            style={[styles.webview, { height }]}
            scrollEnabled={false}
            javaScriptEnabled={true}
            domStorageEnabled={true}
            startInLoadingState={false}
            scalesPageToFit={true}
            originWhitelist={['*']}
            mixedContentMode="always"
        />
    )
}

const styles = StyleSheet.create({
    webview: {
        backgroundColor: dimeTheme.colors.background,
        width: '100%',
    },
})
