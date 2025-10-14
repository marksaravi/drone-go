package plotter

const PLOTTER_DUR_DATA_LEN = 4
const PLOTTER_FLOAT_DATA_LEN = 4
const PLOTTER_INT_DATA_LEN = 2

// 2 byte packet size + 2 byte data/packet
const PLOTER_PACKET_HEADER_LEN = PLOTTER_INT_DATA_LEN * 2

// 4 bytes time, 18 bytes rotations
const PLOTTER_DATA_LEN = PLOTTER_DUR_DATA_LEN + 9*PLOTTER_FLOAT_DATA_LEN
