import { useEffect, useRef, useCallback } from 'react'
import { useMarketStore } from '@/stores/market'
import { API_CONFIG } from '@/constants'
import type { Quote } from '@/types/market'

export function useWebSocket() {
    const wsRef = useRef<WebSocket | null>(null)
    const reconnectRef = useRef<ReturnType<typeof setTimeout> | null>(null)
    const { watchedSymbols, updateQuote, setWsConnected } = useMarketStore()

    const connect = useCallback(() => {
        if (wsRef.current?.readyState === WebSocket.OPEN) return

        try {
            const ws = new WebSocket(API_CONFIG.wsUrl)

            ws.onopen = () => {
                console.log('[WS] Connected')
                setWsConnected(true)

                if (watchedSymbols.length > 0) {
                    ws.send(JSON.stringify({ type: 'SUBSCRIBE', symbols: watchedSymbols }))
                }
            }

            ws.onmessage = (event) => {
                try {
                    const data = JSON.parse(event.data)
                    if (data.type === 'PRICE_UPDATE') {
                        updateQuote(data as Quote)
                    }
                } catch (err) {
                    console.error('[WS] Parse error:', err)
                }
            }

            ws.onerror = (err) => console.error('[WS] Error:', err)

            ws.onclose = () => {
                console.log('[WS] Disconnected')
                setWsConnected(false)
                reconnectRef.current = setTimeout(connect, 3000)
            }

            wsRef.current = ws
        } catch (err) {
            console.error('[WS] Connection error:', err)
            reconnectRef.current = setTimeout(connect, 3000)
        }
    }, [watchedSymbols, updateQuote, setWsConnected])

    const disconnect = useCallback(() => {
        if (reconnectRef.current) clearTimeout(reconnectRef.current)
        wsRef.current?.close()
        wsRef.current = null
        setWsConnected(false)
    }, [setWsConnected])

    const subscribe = useCallback((symbols: string[]) => {
        if (wsRef.current?.readyState === WebSocket.OPEN) {
            wsRef.current.send(JSON.stringify({ type: 'SUBSCRIBE', symbols }))
        }
    }, [])

    const unsubscribe = useCallback((symbols: string[]) => {
        if (wsRef.current?.readyState === WebSocket.OPEN) {
            wsRef.current.send(JSON.stringify({ type: 'UNSUBSCRIBE', symbols }))
        }
    }, [])

    useEffect(() => {
        connect()
        return disconnect
    }, [connect, disconnect])

    useEffect(() => {
        if (wsRef.current?.readyState === WebSocket.OPEN && watchedSymbols.length > 0) {
            subscribe(watchedSymbols)
        }
    }, [watchedSymbols, subscribe])

    return { connect, disconnect, subscribe, unsubscribe }
}
