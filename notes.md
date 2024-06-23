# CHIP-8 interpreter

- 4KB (4,096 bytes) of RAM, from location 0x000 (0) to 0xFFF (4095). The first 512 bytes, from 0x000 to 0x1FF, are where the original interpreter was located, and should not be used by programs.
- 16 registers (`V0` - `VF`) each of 8-bits
- `VF` is a special register for storing flags
- `DT` and `ST` registers for delay and sound timers 
- `I` register for storing memory addresses (16-bit, lower 12 bits used)
- `PC` is 16-bit
- `SP` is 8-bit
- For graphics, programs may also refer to a group of sprites representing the hexadecimal digits 0 through F. These sprites are 5 bytes long, or 8x5 pixels. The data should be stored in the interpreter area of Chip-8 memory (0x000 to 0x1FF).
- CHIP-8 sprites are always eight pixels wide and between one to fifteen pixels high.

Memory Map:
```text
+---------------+= 0xFFF (4095) End of Chip-8 RAM
|               |
|               |
|               |
|               |
|               |
| 0x200 to 0xFFF|
|     Chip-8    |
| Program / Data|
|     Space     |
|               |
|               |
|               |
+- - - - - - - -+= 0x600 (1536) Start of ETI 660 Chip-8 programs
|               |
|               |
|               |
+---------------+= 0x200 (512) Start of most Chip-8 programs
| 0x000 to 0x1FF|
| Reserved for  |
|  interpreter  |
+---------------+= 0x000 (0) Start of Chip-8 RAM
```

Reference - http://devernay.free.fr/hacks/chip8/C8TECH10.HTM