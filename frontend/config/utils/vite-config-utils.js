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

import vue from '@vitejs/plugin-vue'
import license from 'rollup-plugin-license'
import sbom from 'rollup-plugin-sbom'
import cssInjectedByJsPlugin from 'vite-plugin-css-injected-by-js'
import { resolve } from 'path'

export function createCommonViteConfig ({ rootDir, basePublicPath, mode, variant, features, outBase }) {
    return {
        root: rootDir,
        base: basePublicPath,
        mode,
        plugins: [
            vue(),
            cssInjectedByJsPlugin(),
            license({
                cwd: rootDir,
                thirdParty: {
                    allow: '(Apache-2.0 OR BSD-2-Clause OR BSD-3-Clause OR CC0-1.0 OR ISC OR MIT OR Unlicense OR WTFPL OR Zlib OR AFL-2.1 OR LicenseRef-TQSSLA-1.0.2 OR LicenseRef-TQSSLA-1.0.3 OR LicenseRef-TQSPSLA-1.0.1 OR OFL-1.1 OR 0BSD)',
                    includePrivate: true,
                    failOnUnlicensed: true,
                    failOnViolation: true,
                    output: {
                        file: resolve(outBase, 'ThirdPartyNotice.txt'),
                        encoding: 'utf-8',
                    },
                },
            }),
            sbom({
                autodetect: false,
                outFormats: ['json'],
                output: resolve(outBase, 'sbom.json'),
            }),
        ],
        define: {
            'process.env.NODE_ENV': JSON.stringify(mode),
            ...Object.fromEntries(
                Object.entries({}).map(([k, v]) => [k, JSON.stringify(v)]),
            ),
        },
        logLevel: 'info',
    }
}
