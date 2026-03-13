<!--
    Copyright (c) 2026 TQ-Systems GmbH <license@tq-group.com>, D-82229 Seefeld, Germany. All rights reserved.
    Author: Jonathan Backes and the Energy Manager development team

    This software is licensed under the TQ-Systems Product Software License Agreement Version 1.0.3 or any later version.
    You can obtain a copy of the License Agreement in the TQS (TQ-Systems Software Licenses) folder on the following website:
    https://www.tq-group.com/en/support/downloads/tq-software-license-conditions/
    In case of any license issues please contact license@tq-group.com.
-->

<template>
  <div class="row m-2">
    <CardDatapoints v-if="gdr" :gdr="gdr" />
  </div>
</template>

<script setup lang="ts">
import { onBeforeUnmount, onMounted, ref } from 'vue';
import CardDatapoints from './components/CardDatapoints.vue';
import { openGDRSocket, GDR, Closable } from './utils/gdr';

defineOptions({ name: 'ViewApp' })

const gdr = ref<GDR|null>(null);
const gdrSocket = ref<Closable|any>(null);

onMounted(() => {
    try {
        gdrSocket.value = openGDRSocket('smart-meter', (gdrs: Record<string, GDR>) => {
            const gdrData = gdrs['smart-meter']
            if (gdrData) {
                gdr.value = gdrData
            }
        })
    } catch (error) {
        console.log(error);
    }
})

onBeforeUnmount(() => {
    if (gdrSocket.value) {
        gdrSocket.value.close()
    }
})

</script>
