package plotter

const PLOTTER_DUR_DATA_LEN = 4   // micro-second*100 which is 0.0001 seconds
const PLOTTER_FLOAT_DATA_LEN = 4 //float32
const PLOTTER_INT_DATA_LEN = 2

// 2 byte packet size + 2 byte data/packet
const PLOTER_PACKET_HEADER_LEN = PLOTTER_INT_DATA_LEN * 2

// 4 bytes time (micro-second*100 which is 0.0001 seconds), 18 bytes rotations (gyro + accelerometer + rotation with roll, pitch, yaw in float32), total 22 bytes
const PLOTTER_DATA_LEN = PLOTTER_DUR_DATA_LEN + 9*PLOTTER_FLOAT_DATA_LEN

// [2 bytes packet size | 2 bytes data/packet + N * 22 bytes data | 4 byte time | 2 byte gyro roll | 2 byte gyro pitch | 2 byte gyro yaw | 2 byte acc roll | 2 byte acc pitch | 2 byte acc yaw | 2 byte rot roll | 2 byte rot pitch | 2 byte rot yaw ])
