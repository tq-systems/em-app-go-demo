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
} from 'gdr'

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

export function openGDRSocket(topic: string, handler: (gdrs: Record<string, GDR>) => void): Closable {
    return openGDRSocketToURL(`/api/data-transfer/ws/protobuf/gdr/local/values/${topic}`, handler)
}

export function openGDRSocketToURL(url: string, handler: (gdrs: Record<string, GDR>) => void): Closable {
    const scheme = getSocketScheme()
    return openSocket(
        `${scheme}://${window.location.host}${url}`,
        (event) => {
            const gdrs = GDRsDecoder.decode(new Uint8Array(event.data))
            handler(gdrs.GDRs)
        },
    )
}
