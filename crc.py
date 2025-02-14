def generate_crc16_table():
    polynomial = 0xA001
    table = []

    for i in range(256):
        crc = i
        for _ in range(8):
            if crc & 1:
                crc = (crc >> 1) ^ polynomial
            else:
                crc >>= 1
        table.append(crc)

    return table

def generate_go_file():
    table = generate_crc16_table()

    go_code = '''package sensors

// CRC-16 Modbus Table
var crc16ModbusTable = [256]uint16{
'''

    # Adding the table values to the Go file
    for value in table:
        go_code += f"    0x{value:04X},\n"

    go_code += '''
}
'''

    with open('crc16_table.go', 'w') as f:
        f.write(go_code)

    print("Go file 'crc16_table.go' generated successfully.")

if __name__ == "__main__":
    generate_go_file()
