package serial

// ComputeCRC16 calculates the Modbus CRC16 checksum
func ComputeCRC16(data []byte) uint16 {
	var crc uint16 = 0xFFFF
	for _, b := range data {
		crc = (crc >> 8) ^ Crc16ModbusTable[(crc^uint16(b))&0xFF]
	}
	return crc
}

// AppendCRC16 appends the CRC16 checksum to the command
func AppendCRC16(cmd []byte) []byte {
	crc := ComputeCRC16(cmd)
	cmd = append(cmd, byte(crc&0xFF), byte(crc>>8)) // Low byte first, then high byte
	return cmd
}
