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

import { resolve } from 'path'
import { createCommonViteConfig } from './utils/vite-config-utils.js'
import { buildBundle } from './utils/vite-build-utils.js'

;(async () => {
    try {
        await buildVariants({
            rootDir: process.cwd(),
            mode: process.env.NODE_ENV,
            appName: 'go-demo',
        })
    } catch (err) {
        console.error(err)
        process.exit(1)
    }
})()

export async function buildVariants ({
    rootDir = process.cwd(),
    mode = process.env.NODE_ENV || 'production',
    appName,
}) {
    // Determine base path
    const basePublicPath = `/apps/${appName}/`
    let allSucceeded = true
    const entryPoint = resolve(rootDir, 'index.js')
    const outBase = resolve(rootDir, 'dist', 'default')
    const commonConfig = createCommonViteConfig({
        rootDir,
        basePublicPath,
        mode,
        outBase,
    })

    try {
        await buildBundle(
            entryPoint,
            outBase,
            'index.js',
            ['shared/store.js', 'vue'],
            mode !== 'production',
            commonConfig,
            ' -> Building main bundle…',
        )

        // Build i18n bundle without license & SBOM plugins
        const i18nOut = resolve(outBase, 'i18n')
        const i18nPlugins = commonConfig.plugins.filter(p => !['rollup-plugin-license', 'rollup-plugin-sbom'].includes(p.name))
        await buildBundle(
            resolve(rootDir, 'src/lang/default/lang.js'),
            i18nOut,
            'lang.js',
            [],
            mode !== 'production',
            { ...commonConfig, plugins: i18nPlugins },
            ' -> Building i18n bundle (Language: default)…',
        )

        console.log('\n✅ Variant built successfully!')
    } catch (e) {
        console.error('\n❌ Build failed!', e)
        allSucceeded = false
    }
    // Exit with failure code if any variant build failed
    if (!allSucceeded) process.exit(1)
    console.log('\n✅✅✅  Built successfully! ✅✅✅')
}
