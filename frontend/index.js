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

import { createI18n } from 'vue-i18n'
import store from 'shared/store.js'
import { defineVueWebComponent } from '@tq-systems/em-vue3-core'
import ViewApp from './src/ViewApp.vue'
import Icon from './assets/dashboard-icon.svg?url'

const config = {
    path: 'go-demo',
    name: 'go-demo',
    i18nkey: 'go-demo.app_name',
    views: {
        app: {
            component: 'go-demo-main',
            i18nkey: 'go-demo.app_name',
        },
    },
    adapter: [],
}

window.EmRoutes.push(config)

const i18n = createI18n({
    ...window.i18nBaseConfig,
    locale: store.state.language,
    messages: window.EMTranslations,
})

store.commit('addApp', {
    identifier: 'go-demo',
    i18nkey: 'go-demo.app_name',
    widget: [],
    app_icon: Icon,
    iconcolor: 'tertiary1',
    order: 5,
})

store.on('setLanguage:done', (newLang) => {
    i18n.global.locale.value = newLang
})

defineVueWebComponent('go-demo-main', ViewApp, [i18n])
