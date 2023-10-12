package plotter

const DUR_DATA_LEN = 4
const FLOAT_DATA_LEN = 2
const INT_DATA_LEN = 2
const THROTTLE_DATA_LEN = 1
//2 byte packet size + 2 byte data len + 2 byte data/packet
const PLOTER_PACKET_HEADER_SIZE = INT_DATA_LEN * 3
//4 bytes duration, 18 bytes rotations, 1 byte throttle
const PLOTTER_DATA_LEN = DUR_DATA_LEN + 9 * FLOAT_DATA_LEN + THROTTLE_DATA_LEN
const PLOTTER_DATA_PER_PACKET = 100
const PLOTTER_PACKET_SIZE = PLOTER_PACKET_HEADER_SIZE + PLOTTER_DATA_PER_PACKET*PLOTTER_DATA_LEN
