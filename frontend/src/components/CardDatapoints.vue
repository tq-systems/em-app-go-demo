<!--
    Copyright (c) 2026 TQ-Systems GmbH <license@tq-group.com>, D-82229 Seefeld, Germany. All rights reserved.
    Author: Jonathan Backes and the Energy Manager development team

    This software is licensed under the TQ-Systems Product Software License Agreement Version 1.0.3 or any later version.
    You can obtain a copy of the License Agreement in the TQS (TQ-Systems Software Licenses) folder on the following website:
    https://www.tq-group.com/en/support/downloads/tq-software-license-conditions/
    In case of any license issues please contact license@tq-group.com.
-->

<template>
    <div class="card">
        <div class="card-header">
            <div class="row">
                <div class="col-11">
                    <h2>{{ t("go-demo.headlinetable") }}</h2>
                </div>
                <div class="col-1">
                    <span class="badge bg-primary">{{ $d(time, "timeLTS") }}</span>
                </div>
            </div>
        </div>
        <div class="body">
            <div class="row">
                <div class="col-sm-12">
                    <div v-if="gdr" class="table-responsive">
                        <table class="table">
                            <thead>
                                <tr>
                                    <th style="width: 20%" scope="col">
                                        <div class="row">
                                            <h6 class="col-12 col-md-8 pt-2"> {{ t('dataPoints.enhanced-view') }}</h6>
                                            <div class="col-12 col-md-4 pt-2 d-flex justify-content-md-end">
                                                <div class="form-check form-switch m-0">
                                                    <input id="enhancedViewSwitch" v-model="enhanced"
                                                        class="form-check-input" type="checkbox" role="switch">
                                                </div>
                                            </div>
                                        </div> <!-- row -->
                                    </th>
                                    <th style="width: 20%" class="text-center" scope="col">
                                        <h5>{{ t('dataPoints.phaseL1') }}</h5>
                                    </th>
                                    <th style="width: 20%" class="text-center" scope="col">
                                        <h5>{{ t('dataPoints.phaseL2') }}</h5>
                                    </th>
                                    <th style="width: 20%" class="text-center" scope="col">
                                        <h5>{{ t('dataPoints.phaseL3') }}</h5>
                                    </th>
                                    <th style="width: 20%" class="text-center" scope="col">
                                        <h5>{{ t('dataPoints.total') }}</h5>
                                    </th>
                                </tr>
                            </thead>
                            <tbody>
                                <tr>
                                    <td>{{ t('dataPoints.current') }}</td>
                                    <td class="text-center">{{ $n(currentL1, 'fixed1') }} A</td>
                                    <td class="text-center">{{ $n(currentL2, 'fixed1') }} A</td>
                                    <td class="text-center">{{ $n(currentL3, 'fixed1') }} A</td>
                                    <td class="text-center">{{ $n(currentTotal, 'fixed1') }} A</td>
                                </tr>
                                <tr>
                                    <td>{{ t('dataPoints.voltage') }}</td>
                                    <td class="text-center">{{ $n(voltageL1, 'fixed1') }} V</td>
                                    <td class="text-center">{{ $n(voltageL2, 'fixed1') }} V</td>
                                    <td class="text-center">{{ $n(voltageL3, 'fixed1') }} V</td>
                                    <td class="text-center" />
                                </tr>
                                <tr>
                                    <td>{{ t('dataPoints.powerfactor') }}</td>
                                    <td class="text-center">{{ $n(powerfactorL1 / 1000, 'fixed2') }}</td>
                                    <td class="text-center">{{ $n(powerfactorL2 / 1000, 'fixed2') }}</td>
                                    <td class="text-center">{{ $n(powerfactorL3 / 1000, 'fixed2') }}</td>
                                    <td class="text-center">{{ $n(powerfactorTotal / 1000, 'fixed2') }}</td>
                                </tr>
                                <tr>
                                    <th scope="row"> {{ t('dataPoints.active-power') }}</th>
                                    <th scope="row" class="text-center">
                                        {{ powerL1 > 0 ? "+" : "" }}
                                        {{ $n(powerL1, 'fixed1') }} W
                                    </th>
                                    <th scope="row" class="text-center">
                                        {{ powerL2 > 0 ? "+" : "" }}
                                        {{ $n(powerL2, 'fixed1') }} W
                                    </th>
                                    <th scope="row" class="text-center">
                                        {{ powerL3 > 0 ? "+" : "" }}
                                        {{ $n(powerL3, 'fixed1') }} W
                                    </th>
                                    <th scope="row" class="text-center">
                                        {{ powerTotal > 0 ? "+" : "" }}
                                        {{ $n(powerTotal, 'fixed1') }} W
                                    </th>
                                </tr>
                                <Transition name="slide-y-up">
                                    <tr v-show="enhanced">
                                        <td>{{ t('dataPoints.apparent-power') }}</td>
                                        <td class="text-center">
                                            {{ apparentPowerL1 > 0 ? "+" : "" }}
                                            {{ $n(apparentPowerL1, 'fixed1') }} VA
                                        </td>
                                        <td class="text-center">
                                            {{ apparentPowerL2 > 0 ? "+" : "" }}
                                            {{ $n(apparentPowerL2, 'fixed1') }} VA
                                        </td>
                                        <td class="text-center">
                                            {{ apparentPowerL3 > 0 ? "+" : "" }}
                                            {{ $n(apparentPowerL3, 'fixed1') }} VA
                                        </td>
                                        <td class="text-center">
                                            {{ apparentPowerTotal > 0 ? "+" : "" }}
                                            {{ $n(apparentPowerTotal, 'fixed1') }} VA
                                        </td>
                                    </tr>
                                </Transition>
                                <Transition name="slide-y-up">
                                    <tr v-show="enhanced">
                                        <td>{{ t('dataPoints.reactive-power') }}</td>
                                        <td class="text-center">
                                            {{ reactivePowerL1 > 0 ? "+" : "" }}
                                            {{ $n(reactivePowerL1, 'fixed1') }} var
                                        </td>
                                        <td class="text-center">
                                            {{ reactivePowerL2 > 0 ? "+" : "" }}
                                            {{ $n(reactivePowerL2, 'fixed1') }} var
                                        </td>
                                        <td class="text-center">
                                            {{ reactivePowerL3 > 0 ? "+" : "" }}
                                            {{ $n(reactivePowerL3, 'fixed1') }} var
                                        </td>
                                        <td class="text-center">
                                            {{ reactivePowerTotal > 0 ? "+" : "" }}
                                            {{ $n(reactivePowerTotal, 'fixed1') }} var
                                        </td>
                                    </tr>
                                </Transition>
                                <tr>
                                    <th scope="row" rowspan="2">{{ t('dataPoints.active-energy') }}</th>
                                    <th scope="row" class="text-center">+{{ $n(activeenergyL1pos, 'fixed1') }} kWh</th>
                                    <th scope="row" class="text-center">+{{ $n(activeenergyL2pos, 'fixed1') }} kWh</th>
                                    <th scope="row" class="text-center">+{{ $n(activeenergyL3pos, 'fixed1') }} kWh</th>
                                    <th scope="row" class="text-center">+{{ $n(activeenergyTotalpos, 'fixed1') }} kWh
                                    </th>
                                </tr>
                                <tr>
                                    <th scope="row" class="text-center">{{ $n(activeenergyL1neg * -1, 'fixed1') }} kWh
                                    </th>
                                    <th scope="row" class="text-center">{{ $n(activeenergyL2neg * -1, 'fixed1') }} kWh
                                    </th>
                                    <th scope="row" class="text-center">{{ $n(activeenergyL3neg * -1, 'fixed1') }} kWh
                                    </th>
                                    <th scope="row" class="text-center">{{ $n(activeenergyTotalneg * -1, 'fixed1') }}
                                        kWh</th>
                                </tr>
                                <Transition name="slide-y-up">
                                    <tr v-show="enhanced">
                                        <td rowspan="2">{{ t('dataPoints.apparent-energy') }}</td>
                                        <td class="text-center">+{{ $n(apparentenergyL1pos, 'fixed1') }} kVAh</td>
                                        <td class="text-center">+{{ $n(apparentenergyL2pos, 'fixed1') }} kVAh</td>
                                        <td class="text-center">+{{ $n(apparentenergyL3pos, 'fixed1') }} kVAh</td>
                                        <td class="text-center">+{{ $n(apparentenergyTotalpos, 'fixed1') }} kVAh</td>
                                    </tr>
                                </Transition>
                                <Transition name="slide-y-up">
                                    <tr v-show="enhanced">
                                        <td class="text-center">{{ $n(apparentenergyL1neg * -1, 'fixed1') }} kVAh</td>
                                        <td class="text-center">{{ $n(apparentenergyL2neg * -1, 'fixed1') }} kVAh</td>
                                        <td class="text-center">{{ $n(apparentenergyL3neg * -1, 'fixed1') }} kVAh</td>
                                        <td class="text-center">{{ $n(apparentenergyTotalneg * -1, 'fixed1') }} kVAh
                                        </td>
                                    </tr>
                                </Transition>
                                <Transition name="slide-y-up">
                                    <tr v-show="enhanced">
                                        <td rowspan="2">{{ t('dataPoints.reactive-energy') }}</td>
                                        <td class="text-center">+{{ $n(reactiveenergyL1pos, 'fixed1') }} kvarh</td>
                                        <td class="text-center">+{{ $n(reactiveenergyL2pos, 'fixed1') }} kvarh</td>
                                        <td class="text-center">+{{ $n(reactiveenergyL3pos, 'fixed1') }} kvarh</td>
                                        <td class="text-center">+{{ $n(reactiveenergyTotalpos, 'fixed1') }} kvarh</td>
                                    </tr>
                                </Transition>
                                <Transition name="slide-y-up">
                                    <tr v-show="enhanced">
                                        <td class="text-center">{{ $n(reactiveenergyL1neg * -1, 'fixed1') }} kvarh</td>
                                        <td class="text-center">{{ $n(reactiveenergyL2neg * -1, 'fixed1') }} kvarh</td>
                                        <td class="text-center">{{ $n(reactiveenergyL3neg * -1, 'fixed1') }} kvarh</td>
                                        <td class="text-center">{{ $n(reactiveenergyTotalneg * -1, 'fixed1') }} kvarh
                                        </td>
                                    </tr>
                                </Transition>
                            </tbody>
                        </table>
                    </div>
                    <span v-else class="spinner-border spinner-border-lg d-block mx-auto" role="status" />
                </div>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import http from '@tq-systems/em-http-service'
import { GDR } from '../utils/gdr'
import { OBISCode } from '../utils/obis'
import { useI18n } from 'vue-i18n'
import { onBeforeUnmount, onMounted, ref } from 'vue'

defineOptions({ name: 'CardDatapoints' })

const props = defineProps({ gdr: { type: Object as () => GDR, required: true } })

const { t } = useI18n()

const time = ref('')
const enhanced = ref(false)
const currentL1 = ref(0)
const currentL2 = ref(0)
const currentL3 = ref(0)
const currentTotal = ref(0)
const voltageL1 = ref(0)
const voltageL2 = ref(0)
const voltageL3 = ref(0)
const powerfactorL1 = ref(0)
const powerfactorL2 = ref(0)
const powerfactorL3 = ref(0)
const powerfactorTotal = ref(0)
const powerL1 = ref(0)
const powerL2 = ref(0)
const powerL3 = ref(0)
const powerTotal = ref(0)
const apparentPowerL1 = ref(0)
const apparentPowerL2 = ref(0)
const apparentPowerL3 = ref(0)
const apparentPowerTotal = ref(0)
const reactivePowerL1 = ref(0)
const reactivePowerL2 = ref(0)
const reactivePowerL3 = ref(0)
const reactivePowerTotal = ref(0)
const activeenergyL1pos = ref(0)
const activeenergyL2pos = ref(0)
const activeenergyL3pos = ref(0)
const activeenergyTotalpos = ref(0)
const activeenergyL1neg = ref(0)
const activeenergyL2neg = ref(0)
const activeenergyL3neg = ref(0)
const activeenergyTotalneg = ref(0)
const apparentenergyL1pos = ref(0)
const apparentenergyL2pos = ref(0)
const apparentenergyL3pos = ref(0)
const apparentenergyTotalpos = ref(0)
const apparentenergyL1neg = ref(0)
const apparentenergyL2neg = ref(0)
const apparentenergyL3neg = ref(0)
const apparentenergyTotalneg = ref(0)
const reactiveenergyL1pos = ref(0)
const reactiveenergyL2pos = ref(0)
const reactiveenergyL3pos = ref(0)
const reactiveenergyTotalpos = ref(0)
const reactiveenergyL1neg = ref(0)
const reactiveenergyL2neg = ref(0)
const reactiveenergyL3neg = ref(0)
const reactiveenergyTotalneg = ref(0)
const timeUpdateInterval = ref<ReturnType<typeof setInterval> | null>(null)


async function getTime() {
    try {
        const response = await http.get('/api/go-demo/time')
        time.value = response.data
    } catch (error) {
        console.log(error);
    }
}

function mA2A(mA: number): number {
    return mA / 1000
}

function mW2W(val: number): number {
    return val / 1000
}

function mWh2kWh(val: number): number {
    return val / 1000000
}

function mV2V(mV: number): number {
    return mV / 1000
}

function updateDataPoints() {
    getTime()
}

function get(obis: string): number {
    const gdr: GDR = props.gdr
    if (!gdr.values) {
        return 0
    }

    return gdr.values[OBISCode.DecodeString(obis)] ?? 0
}

function getPower(posIdx: string, negIdx: string): number {
    const gdr: GDR = props.gdr
    if (!gdr.values) {
        return 0
    }

    return (gdr.values[OBISCode.DecodeString(posIdx)] ?? 0) - (gdr.values[OBISCode.DecodeString(negIdx)] ?? 0)
}

function getData() {
    if (!props.gdr || !props.gdr.values) {
        console.warn('getData called but gdr or gdr.values is not available. Aborting.')
        return
    }

    currentL1.value = mA2A(get('1-0:31.4.0*255'))
    currentL2.value = mA2A(get('1-0:51.4.0*255'))
    currentL3.value = mA2A(get('1-0:71.4.0*255'))
    currentTotal.value = mA2A((get('1-0:31.4.0*255') + get('1-0:51.4.0*255') + get('1-0:71.4.0*255')))

    voltageL1.value = mV2V(get('1-0:32.4.0*255'))
    voltageL2.value = mV2V(get('1-0:52.4.0*255'))
    voltageL3.value = mV2V(get('1-0:72.4.0*255'))

    powerfactorL1.value = get('1-0:33.4.0*255')
    powerfactorL2.value = get('1-0:53.4.0*255')
    powerfactorL3.value = get('1-0:73.4.0*255')
    powerfactorTotal.value = get('1-0:13.4.0*255')

    powerL1.value = mW2W(getPower('1-0:21.4.0*255', '1-0:22.4.0*255'))
    powerL2.value = mW2W(getPower('1-0:41.4.0*255', '1-0:42.4.0*255'))
    powerL3.value = mW2W(getPower('1-0:61.4.0*255', '1-0:62.4.0*255'))
    powerTotal.value = mW2W(getPower('1-0:1.4.0*255', '1-0:2.4.0*255'))

    apparentPowerL1.value = mW2W(getPower('1-0:29.4.0*255', '1-0:30.4.0*255'))
    apparentPowerL2.value = mW2W(getPower('1-0:49.4.0*255', '1-0:50.4.0*255'))
    apparentPowerL3.value = mW2W(getPower('1-0:69.4.0*255', '1-0:70.4.0*255'))
    apparentPowerTotal.value = mW2W(getPower('1-0:9.4.0*255', '1-0:10.4.0*255'))

    reactivePowerL1.value = mW2W(getPower('1-0:23.4.0*255', '1-0:24.4.0*255'))
    reactivePowerL2.value = mW2W(getPower('1-0:43.4.0*255', '1-0:44.4.0*255'))
    reactivePowerL3.value = mW2W(getPower('1-0:63.4.0*255', '1-0:64.4.0*255'))
    reactivePowerTotal.value = mW2W(getPower('1-0:3.4.0*255', '1-0:4.4.0*255'))

    activeenergyL1pos.value = mWh2kWh(get('1-0:21.8.0*255'))
    activeenergyL2pos.value = mWh2kWh(get('1-0:41.8.0*255'))
    activeenergyL3pos.value = mWh2kWh(get('1-0:61.8.0*255'))
    activeenergyTotalpos.value = mWh2kWh(get('1-0:1.8.0*255'))
    activeenergyL1neg.value = mWh2kWh(get('1-0:22.8.0*255'))
    activeenergyL2neg.value = mWh2kWh(get('1-0:42.8.0*255'))
    activeenergyL3neg.value = mWh2kWh(get('1-0:62.8.0*255'))
    activeenergyTotalneg.value = mWh2kWh(get('1-0:2.8.0*255'))

    apparentenergyL1pos.value = mWh2kWh(get('1-0:29.8.0*255'))
    apparentenergyL2pos.value = mWh2kWh(get('1-0:49.8.0*255'))
    apparentenergyL3pos.value = mWh2kWh(get('1-0:69.8.0*255'))
    apparentenergyTotalpos.value = mWh2kWh(get('1-0:9.8.0*255'))
    apparentenergyL1neg.value = mWh2kWh(get('1-0:30.8.0*255'))
    apparentenergyL2neg.value = mWh2kWh(get('1-0:50.8.0*255'))
    apparentenergyL3neg.value = mWh2kWh(get('1-0:70.8.0*255'))
    apparentenergyTotalneg.value = mWh2kWh(get('1-0:10.8.0*255'))

    reactiveenergyL1pos.value = mWh2kWh(get('1-0:23.8.0*255'))
    reactiveenergyL2pos.value = mWh2kWh(get('1-0:43.8.0*255'))
    reactiveenergyL3pos.value = mWh2kWh(get('1-0:63.8.0*255'))
    reactiveenergyTotalpos.value = mWh2kWh(get('1-0:3.8.0*255'))
    reactiveenergyL1neg.value = mWh2kWh(get('1-0:24.8.0*255'))
    reactiveenergyL2neg.value = mWh2kWh(get('1-0:44.8.0*255'))
    reactiveenergyL3neg.value = mWh2kWh(get('1-0:64.8.0*255'))
    reactiveenergyTotalneg.value = mWh2kWh(get('1-0:4.8.0*255'))
}

onMounted(() => {
    updateDataPoints()
    getData()
    timeUpdateInterval.value = setInterval(() => {
        updateDataPoints()
        getData()
    }, 1000)
})

onBeforeUnmount(() => {
    if (timeUpdateInterval.value) {
        clearInterval(timeUpdateInterval.value)
        timeUpdateInterval.value = null
    }
})

</script>
