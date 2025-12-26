import { type FC, useMemo } from 'react'
import { StyleSheet, View, Dimensions } from 'react-native'
import { WebView } from 'react-native-webview'
import { dimeTheme } from '@/constants/theme'

interface KLineDataPoint {
    time: number | string  // Unix timestamp or date string
    open: number
    high: number
    low: number
    close: number
    volume?: number
}

interface KLineChartProps {
    data: KLineDataPoint[]
    height?: number
    showVolume?: boolean
    period?: '1D' | '1W' | '1M' | '3M' | '1Y' | 'ALL'
}

const { width: screenWidth } = Dimensions.get('window')

export const KLineChart: FC<KLineChartProps> = ({
    data,
    height = 300,
    showVolume = true,
}) => {
    // Generate HTML with TradingView Lightweight Charts
    const chartHtml = useMemo(() => {
        const candleData = data.map(d => ({
            time: typeof d.time === 'number' ? d.time / 1000 : d.time, // Convert to seconds if ms
            open: d.open,
            high: d.high,
            low: d.low,
            close: d.close,
        }))

        const volumeData = showVolume ? data.map(d => ({
            time: typeof d.time === 'number' ? d.time / 1000 : d.time,
            value: d.volume ?? 0,
            color: d.close >= d.open ? 'rgba(0, 230, 118, 0.5)' : 'rgba(255, 59, 48, 0.5)',
        })) : []

        return `
<!DOCTYPE html>
<html>
<head>
    <meta name="viewport" content="width=device-width, initial-scale=1.0, maximum-scale=1.0, user-scalable=no">
    <script src="https://unpkg.com/lightweight-charts@4.1.0/dist/lightweight-charts.standalone.production.js"></script>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body { 
            background-color: ${dimeTheme.colors.background}; 
            overflow: hidden;
            -webkit-user-select: none;
            -webkit-touch-callout: none;
        }
        #chart { 
            width: 100%; 
            height: ${height}px; 
        }
    </style>
</head>
<body>
    <div id="chart"></div>
    <script>
        const chart = LightweightCharts.createChart(document.getElementById('chart'), {
            width: ${screenWidth},
            height: ${height},
            layout: {
                background: { type: 'solid', color: '${dimeTheme.colors.background}' },
                textColor: '${dimeTheme.colors.textSecondary}',
            },
            grid: {
                vertLines: { color: '${dimeTheme.colors.border}' },
                horzLines: { color: '${dimeTheme.colors.border}' },
            },
            crosshair: {
                mode: LightweightCharts.CrosshairMode.Normal,
                vertLine: {
                    color: '${dimeTheme.colors.primary}',
                    width: 1,
                    style: LightweightCharts.LineStyle.Dashed,
                },
                horzLine: {
                    color: '${dimeTheme.colors.primary}',
                    width: 1,
                    style: LightweightCharts.LineStyle.Dashed,
                },
            },
            rightPriceScale: {
                borderColor: '${dimeTheme.colors.border}',
            },
            timeScale: {
                borderColor: '${dimeTheme.colors.border}',
                timeVisible: true,
                secondsVisible: false,
            },
            handleScroll: {
                mouseWheel: true,
                pressedMouseMove: true,
                horzTouchDrag: true,
                vertTouchDrag: false,
            },
            handleScale: {
                axisPressedMouseMove: true,
                mouseWheel: true,
                pinch: true,
            },
        });

        // Candlestick series
        const candleSeries = chart.addCandlestickSeries({
            upColor: '${dimeTheme.colors.profit}',
            downColor: '${dimeTheme.colors.loss}',
            borderUpColor: '${dimeTheme.colors.profit}',
            borderDownColor: '${dimeTheme.colors.loss}',
            wickUpColor: '${dimeTheme.colors.profit}',
            wickDownColor: '${dimeTheme.colors.loss}',
        });

        const candleData = ${JSON.stringify(candleData)};
        candleSeries.setData(candleData);

        ${showVolume ? `
        // Volume series
        const volumeSeries = chart.addHistogramSeries({
            color: 'rgba(0, 230, 118, 0.5)',
            priceFormat: { type: 'volume' },
            priceScaleId: '',
        });
        volumeSeries.priceScale().applyOptions({
            scaleMargins: { top: 0.8, bottom: 0 },
        });

        const volumeData = ${JSON.stringify(volumeData)};
        volumeSeries.setData(volumeData);
        ` : ''}

        // Fit content
        chart.timeScale().fitContent();
    </script>
</body>
</html>
        `
    }, [data, height, showVolume])

    if (data.length === 0) {
        return (
            <View style={[styles.container, { height }]}>
                <View style={styles.placeholder} />
            </View>
        )
    }

    return (
        <View style={[styles.container, { height }]}>
            <WebView
                source={{ html: chartHtml }}
                style={styles.webview}
                scrollEnabled={false}
                javaScriptEnabled={true}
                domStorageEnabled={true}
                originWhitelist={['*']}
                onError={(error) => console.error('WebView error:', error)}
            />
        </View>
    )
}

const styles = StyleSheet.create({
    container: {
        width: screenWidth,
        backgroundColor: dimeTheme.colors.background,
    },
    webview: {
        flex: 1,
        backgroundColor: 'transparent',
    },
    placeholder: {
        flex: 1,
        backgroundColor: dimeTheme.colors.surface,
        borderRadius: dimeTheme.radius.md,
    },
})

export type { KLineDataPoint }
