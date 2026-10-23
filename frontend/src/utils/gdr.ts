/*
 * Copyright (c) 2026 TQ-Systems GmbH <license@tq-group.com>, D-82229 Seefeld,
 * Germany. All rights reserved.
 * Author: Jonathan Backes and the Energy Manager development team
 *
 * This software is licensed under the TQ-Systems Product Software License
 * Agreement Version 1.0.3 or any later version.
 * You can obtain a copy of the License Agreement in the TQS (TQ-Systems
 * Software Licenses) folder on the following website:
 * https://www.tq-group.com/en/support/downloads/tq-software-license-conditions/
 * In case of any license issues please contact license@tq-group.com.
 */

import {
    GDRs as GDRsDecoder,
    IGDR as GDR,
} from '@tq-systems/em-gdr'

export {
    GDR,
}

export interface Closable {
    close: () => void;
}

function getSocketScheme(): string {
    return window.location.protocol === 'http:' ? 'ws' : 'wss'
}

function getAuthToken(): string {
    return 'Bearer ' + localStorage.getItem('token')
}

function openSocket(url: string, handler: (event: MessageEvent) => void): Closable {
    let closed = false
    let socket: WebSocket

    const doOpenSocket = () => {
        if (closed) {
            return
        }

        socket = new WebSocket(url)
        socket.binaryType = 'arraybuffer'
        socket.onopen = () => {
            socket.send(getAuthToken())
        }

        socket.onmessage = handler

        socket.onclose = () => setTimeout(doOpenSocket, 5000)
    }

    doOpenSocket()

    return {
        close(): void {
            closed = true
            socket.close()
        },
    }
}

// protobufjs may represent uint64 map keys as 8-byte hash strings.
// We first detect already-decimal keys using /^\d+$/ and keep them unchanged,
// then convert only hashed keys back to decimal strings used by OBIS lookups.
function decodeUInt64HashKey (key: string): string {
    // Fast path for previous/runtime variants where keys are already decimal.
    if (/^\d+$/.test(key)) {
        return key
    }
    if (key.length !== 8) {
        return key
    }

    const bytes = Uint8Array.from(key, (char) => char.charCodeAt(0))
    return new DataView(bytes.buffer).getBigUint64(0, true).toString()
}

// Normalizes GDR values map keys so consumers can always access entries using
// decimal OBIS keys, independent of protobufjs long/hash runtime behavior.
function normalizeGDRValuesKeys (gdrs: Record<string, GDR>): Record<string, GDR> {
    for (const gdrId of Object.keys(gdrs ?? {})) {
        const gdr = gdrs[gdrId]
        if (!gdr || !gdr.values) {
            continue
        }

        const normalizedValues: Record<string, number> = {}
        for (const key of Object.keys(gdr.values)) {
            const rawValue = gdr.values[key]
            const numValue = typeof rawValue === 'object' && rawValue !== null && 'toNumber' in rawValue
                ? (rawValue as any).toNumber()
                : Number(rawValue)
            normalizedValues[decodeUInt64HashKey(key)] = numValue
        }
        gdr.values = normalizedValues
    }

    return gdrs
}

function openGDRSocketToURL(url: string, handler: (gdrs: Record<string, GDR>) => void): Closable {
    const scheme = getSocketScheme()
    return openSocket(
        `${scheme}://${window.location.host}${url}`,
        (event) => {
            const gdrs = GDRsDecoder.decode(new Uint8Array(event.data))
            const decodedGdrs = GDRsDecoder.toObject(gdrs, { longs: Number, defaults: true }).GDRs
            handler(normalizeGDRValuesKeys(decodedGdrs))
        },
    )
}

export function openGDRSocket(topic: string, handler: (gdrs: Record<string, GDR>) => void): Closable {
    return openGDRSocketToURL(`/api/data-transfer/ws/protobuf/gdr/local/values/${topic}`, handler)
}
