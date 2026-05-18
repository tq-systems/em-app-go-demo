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
import devConfig from './utils/vite-dev-utils.js'
import { loadDevEnv, createDevDefine, createDevInput, createDevEntryNaming } from './utils/dev-config-utils.js'

;(async () => {
    try {
        await runDevBuild({
            rootDir: process.cwd(),
            mode: process.env.NODE_ENV,
            appName: 'go-demo',
        })
    } catch (err) {
        console.error(err)
        process.exit(1)
    }
})()

// Determine the current file path
async function runDevBuild ({
    rootDir = process.cwd(),
    mode = process.env.NODE_ENV || 'development',
    appName,
}) {
    // Load environment variables and determine branding
    const { branding, env } = loadDevEnv(mode, rootDir, '')

    // Load and select variant configuration

    // Load base vite configuration
    const config = devConfig({ mode, appName })

    // Inject define variables: NODE_ENV, BRANDING, and feature flags
    const baseDefine = createDevDefine(env)
    config.define = {
        ...baseDefine,
        ...Object.fromEntries(
            Object.entries({}).map(([key, val]) =>
                [key, JSON.stringify(val)],
            ),
        ),
    }

    // Enable sourcemaps for easier debugging
    config.build.sourcemap = true

    // Enable watch mode
    config.build = config.build || {}
    config.build.watch = {}
    config.build.rollupOptions = config.build.rollupOptions || {}
    config.build.rollupOptions.input = createDevInput(rootDir)
    config.build.rollupOptions.output = {
        entryFileNames: createDevEntryNaming(),
    }

    console.log(`🔁 Dev build watch: mode="${mode}", branding="${branding}"`)
    await build(config)
}
