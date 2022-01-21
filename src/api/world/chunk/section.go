package chunk

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"reflect"

	proto "github.com/golangminecraft/minecraft-server/src/api/protocol"
	"github.com/golangminecraft/minecraft-server/src/util"
)

type ChunkSection struct {
	Data    []uint32      `nbt:"data"`
	Palette []PaletteItem `nbt:"palette"`
}

func NewChunkSection() *ChunkSection {
	return &ChunkSection{
		Data: make([]uint32, 16*16*16),
		Palette: []PaletteItem{
			{
				ID:         0,
				Name:       "minecraft:air",
				Properties: make(map[string]interface{}),
			},
		},
	}
}

func (c ChunkSection) BlockCount() (count int16) {
	for _, data := range c.Data {
		if data != 0 {
			continue
		}

		count++
	}

	return
}

func (c ChunkSection) GetBlock(x, y, z int64) *PaletteItem {
	index := x + z*16 + y*16*16

	paletteIndex := c.Data[index]

	return &c.Palette[paletteIndex]
}

func (c *ChunkSection) SetBlock(x, y, z int64, item PaletteItem) {
	index := x + z*16 + y*16*16

	var paletteIndex uint32
	var paletteFound bool = false

	for k, v := range c.Palette {
		if v.Name != item.Name || !reflect.DeepEqual(v.Properties, item.Properties) {
			continue
		}

		paletteIndex = uint32(k)
		paletteFound = true

		break
	}

	if paletteFound {
		c.Data[index] = paletteIndex
	} else {
		c.Data[index] = uint32(len(c.Palette))

		c.Palette = append(c.Palette, item)
	}
}

func (s ChunkSection) Encode(w io.Writer) error {
	if err := binary.Write(w, binary.BigEndian, s.Data); err != nil {
		return err
	}

	if err := binary.Write(w, binary.BigEndian, uint32(len(s.Palette))); err != nil {
		return err
	}

	for _, item := range s.Palette {
		if err := item.Encode(w); err != nil {
			return err
		}
	}

	return nil
}

func (s *ChunkSection) Decode(r io.Reader) error {
	if err := binary.Read(r, binary.BigEndian, s.Data); err != nil {
		return err
	}

	var paletteLength uint32

	if err := binary.Read(r, binary.BigEndian, &paletteLength); err != nil {
		return err
	}

	var i uint32

	s.Palette = make([]PaletteItem, paletteLength)

	for i = 0; i < paletteLength; i++ {
		if err := s.Palette[i].Decode(r); err != nil {
			return err
		}
	}

	return nil
}

func (s ChunkSection) RawSectionData() ([]byte, error) {
	buf := &bytes.Buffer{}

	// 1. Block Count (non-air)
	if _, err := proto.Short(s.BlockCount()).Encode(buf); err != nil {
		return nil, err
	}

	// 2. Block States
	bitsPerEntry := int(math.Ceil(math.Log2(float64(len(s.Palette)))))

	// 2.1 Bits Per Entry
	if _, err := proto.UnsignedByte(bitsPerEntry).Encode(buf); err != nil {
		return nil, err
	}

	// 2.2 Palette
	if bitsPerEntry <= 4 {
		// 2.2.1 Palette Length
		if _, err := proto.VarInt(len(s.Palette)).Encode(buf); err != nil {
			return nil, err
		}

		// 2.2.2 Palette Data
		for _, item := range s.Palette {
			if _, err := proto.VarInt(item.ID).Encode(buf); err != nil {
				return nil, err
			}
		}
	} else if bitsPerEntry >= 5 && bitsPerEntry <= 8 {
		return nil, fmt.Errorf("not implemented") // TODO
	} else {
		return nil, fmt.Errorf("unable to encode section data, bits per entry is outside a valid range: %d", bitsPerEntry)
	}

	blockData := make([]int64, len(s.Data))

	for k, data := range s.Data {
		blockData[k] = int64(data)
	}

	compactedData := util.CompactInt64Array(blockData, bitsPerEntry)

	// 2.3 Data Array Length
	if _, err := proto.VarInt(len(compactedData)).Encode(buf); err != nil {
		return nil, err
	}

	// 2.4 Data Array
	if err := binary.Write(buf, binary.BigEndian, compactedData); err != nil {
		return nil, err
	}

	// 3. Biomes

	// 3.1 Bits Per Entry
	if _, err := proto.UnsignedByte(1).Encode(buf); err != nil {
		return nil, err
	}

	// 3.2.1 Palette Length
	if _, err := proto.VarInt(1).Encode(buf); err != nil {
		return nil, err
	}

	// 3.2.2 Palette
	// TODO implement chunk biomes
	if _, err := proto.VarInt(5).Encode(buf); err != nil { // Desert
		return nil, err
	}

	// 3.3 Data Array Length
	if _, err := proto.VarInt(64).Encode(buf); err != nil {
		return nil, err
	}

	// 3.4 Data Array
	if err := binary.Write(buf, binary.BigEndian, make([]int64, 1)); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
