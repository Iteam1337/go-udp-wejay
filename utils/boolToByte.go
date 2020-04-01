package utils

func BoolToByte(b bool) byte {
	if b {
		return byte(1)
	}

	return byte(0)
}
