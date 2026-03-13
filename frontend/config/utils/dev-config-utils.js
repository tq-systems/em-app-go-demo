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

import { loadEnv } from 'vite'
import { resolve } from 'path'

export function loadDevEnv (mode, rootDir, prefix = '') {
    const env = loadEnv(mode, rootDir, prefix)
    return { env }
}

export function createDevDefine (env) {
    return {
        'process.env.NODE_ENV': JSON.stringify(env.NODE_ENV),
    }
}

export function createDevInput (rootDir) {
    const defaultEntry = resolve(rootDir, 'index.js')
    const appEntry = defaultEntry

    return {
        app: appEntry,
        lang: `${rootDir}/src/lang/default/lang.js`,
    }
}

export function createDevEntryNaming () {
    return chunk => chunk.name === 'lang' ? 'i18n/lang.js' : 'index.js'
}
