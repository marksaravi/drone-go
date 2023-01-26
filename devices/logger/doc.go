/*
Data Package format
Packet Format:

bytes: 0          1            2           3            4          5           6           ...n
┌────────────┬────────────┬───────────┬───────────┬───────────┬───────────┬─────────────┬─────────┐
│ packet len │ packet len │ packet id │ packet id │ packet id │ packet id │ format code │ data... │
│ high byte  │ low byte   │ byte #3   │ byte #2   │ byte #1   │ byte #0   │             │         │
└────────────┴────────────┴───────────┴───────────┴───────────┴───────────┴─────────────┴─────────┘


Type 16, Simple Roll, Pitch, Yaw data serialisation
packet len: 10 + 6 * (number of data)
Packet ID is in ms from start time
Format Code: 16
Roll, Pitch, Yaw range: -360..360
Roll, Pitch, Yaw data type: int16
Decimal Precision: 1 digit
Roll, Pitch, Yaw math conversion: round(original value * 10)
Roll, Pitch, Yaw to byte conversion: LittleEndian
Time Interval range: 1..255ms
Time Interval data type: byte

Packet Information (bytes 0..9):
bytes: 0          1            2           3            4          5        6        7             8                   9
┌────────────┬────────────┬───────────┬───────────┬───────────┬───────────┬────┬───────────────┬────────────────┬────────────────┐
│ packet len │ packet len │ packet id │ packet id │ packet id │ packet id │ 16 │ time interval │ number of data │ number of data │
│ high byte  │ low byte   │ byte #3   │ byte #2   │ byte #1   │ byte #0   │    │ (ms)          │ high byte      │ low byte       │
└────────────┴────────────┴───────────┴───────────┴───────────┴───────────┴────┴───────────────┴────────────────┴────────────────┘

Packet Data (bytes 10..9+6*(number of data)):
bytes: 10       11          12           13         14         15          ...                                                     10..9+6*(number of data)
┌───────────┬──────────┬───────────┬───────────┬───────────┬───────────┬───────────┬──────────┬───────────┬──────────┬───────────┬──────────┐
│ Roll      │ Roll     │ Pitch     │ Pitch     │ Yaw       │ Yaw       │ Roll      │ Roll     │ Pitch     │ Pitch    │ Yaw       │ Yaw      │
│ high byte │ low byte │ high byte │ low byte  │ high byte │ low byte  │ high byte │ low byte │ high byte │ low byte │ high byte │ low byte │
└───────────┴──────────┴───────────┴───────────┴───────────┴───────────┴───────────┴──────────┴───────────┴──────────┴───────────┴──────────┘

*/

package logger
