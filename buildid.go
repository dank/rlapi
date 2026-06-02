package rlapi

func decodeBuildID(s string) int32 {
	buf := make([]byte, 0, len(s)*2)
	for _, r := range s {
		buf = append(buf, byte(r), byte(r>>8))
	}
	return crc32(buf, 0)
}

// crc32 computes the non-reflected (big-endian) CRC-32
func crc32(data []byte, seed uint32) int32 {
	const poly = uint32(0x04C11DB7)
	crc := seed ^ 0xFFFF_FFFF
	for _, b := range data {
		crc ^= uint32(b) << 24
		for range 8 {
			if crc&0x8000_0000 != 0 {
				crc = (crc << 1) ^ poly
			} else {
				crc <<= 1
			}
		}
	}
	return int32(crc ^ 0xFFFF_FFFF)
}
