/*
 * Copyright (c) 2025-2026 TQ-Systems GmbH <license@tq-group.com>, D-82229
 * Seefeld, Germany. All rights reserved.
 * Author: Maximilian Eschenbacher and the Energy Manager development team
 *
 * This software is licensed under the TQ-Systems Product Software License
 * Agreement Version 1.0.3 or any later version.
 * You can obtain a copy of the License Agreement in the TQS (TQ-Systems
 * Software Licenses) folder on the following website:
 * https://www.tq-group.com/en/support/downloads/tq-software-license-conditions/
 * In case of any license issues please contact license@tq-group.com.
 */

package modbus

func bytesToUint16(b []byte) []uint16 {
	out := make([]uint16, 0, len(b)/2)
	for i := 0; i+1 < len(b); i += 2 {
		out = append(out, uint16(b[i])<<8|uint16(b[i+1]))
	}
	return out
}

func bytesToBools(b []byte, n int) []bool {
	out := make([]bool, n)
	for i := range out {
		byteIdx := i / 8
		bitIdx := i % 8
		if byteIdx < len(b) {
			out[i] = (b[byteIdx] & (1 << bitIdx)) != 0
		}
	}
	return out
}

func uint16sToBytes(vals []uint16) []byte {
	out := make([]byte, 0, len(vals)*2)
	for _, v := range vals {
		out = append(out, byte(v>>8), byte(v&0xFF))
	}
	return out
}

func boolsToBytes(bits []bool) []byte {
	if len(bits) == 0 {
		return nil
	}
	out := make([]byte, (len(bits)+7)/8)
	for i, b := range bits {
		if b {
			out[i/8] |= 1 << (i % 8)
		}
	}
	return out
}

func chooseInt(v int, def int) int {
	if v == 0 {
		return def
	}
	return v
}

func chooseString(v, def string) string {
	if v == "" {
		return def
	}
	return v
}
