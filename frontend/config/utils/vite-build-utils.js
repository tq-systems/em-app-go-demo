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

import { build } from 'vite'

export async function buildBundle (
    inputPath,
    outDir,
    entryFileName,
    external = [],
    sourcemap = false,
    commonConfig,
    logMessage,
) {
    console.log(logMessage)
    try {
        await build({
            ...commonConfig,
            build: {
                assetsInlineLimit: 0,
                outDir,
                sourcemap,
                rollupOptions: {
                    input: { [entryFileName.replace(/\..+$/, '')]: inputPath },
                    output: { entryFileNames: entryFileName },
                    external,
                },
            },
        })
    } catch (error) {
        console.error(`   ❌ Failed to build bundle "${entryFileName}" in "${outDir}".`)
        throw error
    }
}
