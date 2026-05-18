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

export class OBISCode {

    public Media: number
    public Channel: number
    public Indicator: number
    public Mode: number
    public Quantities: number
    public Storage: number

    constructor(media: number, channel: number, indicator: number, mode: number, quantities: number, storage: number) {
        this.Media = media
        this.Channel = channel
        this.Indicator = indicator
        this.Mode = mode
        this.Quantities = quantities
        this.Storage = storage
    }


    public static DecodeStringToOBIS(obisSTR: string): OBISCode {
        // e.g. 1-0:21.4.0*256
        const first = obisSTR.split('-')
        const media = first[0]
        const second = first[1].split(':')
        const channel = second[0]
        const third = second[1].split('.')
        const indicator = third[0]
        const mode = third[1]
        const fourth = third[2].split('*')
        const quantities = fourth[0]
        const storage = fourth[1]

        return new OBISCode(Number(media),
            Number(channel),
            Number(indicator),
            Number(mode),
            Number(quantities),
            Number(storage))
    }



    public static DecodeString(value: string): number {
        const obis = this.DecodeStringToOBIS(value)
        return obis.Encode()
    }

        public Encode (): number {
        let value: number = 0

        const byteArray = [this.Storage, this.Quantities, this.Mode, this.Indicator, this.Channel, this.Media, 0, 0]

        for (let i = byteArray.length - 1; i >= 0; i--) {
            value = (value * 256) + byteArray[i]
        }

        return value
    }

}
