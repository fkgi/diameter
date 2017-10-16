package msg

// Convert from bool list to byte array
func botob(fs ...bool) (b []byte) {
	b = make([]byte, len(fs)/8+1)
	for i, f := range fs {
		if f {
			b[i/8] |= (0x80 >> uint(i%8))
		}
	}
	return
}

// Convert from byte to bool list
func btobo(b []byte) (fs []bool) {
	fs = make([]bool, 8*len(b))
	for i := range fs {
		fs[i] = (b[i/8]>>uint(7-i%8))&0x01 == 0x01
	}
	return
}
