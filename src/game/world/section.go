package world

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"reflect"

	proto "github.com/anchormc/anchor/src/api/protocol"
	"github.com/anchormc/anchor/src/api/world"
	"github.com/anchormc/anchor/src/util"
)

type ChunkSection struct {
	Index           int64                  `nbt:"Index"`
	Y               int64                  `nbt:"Y"`
	Data            []int16                `nbt:"Data"`
	Palette         []world.Block          `nbt:"Palette"`
	blockUpdateChan chan world.BlockUpdate `nbt:"-"`
}

func NewChunkSection(index int64) world.ChunkSection {
	return &ChunkSection{
		Index: index,
		Y:     index * world.SectionHeight,
		Data:  make([]int16, 16*16*16),
		Palette: []world.Block{
			{
				ID:         0,
				Name:       "minecraft:air",
				Properties: make(map[string]interface{}),
			},
		},
		blockUpdateChan: make(chan world.BlockUpdate),
	}
}

func (c ChunkSection) OnBlockUpdate() world.BlockUpdate {
	return <-c.blockUpdateChan
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

func (c ChunkSection) GetBlock(x, y, z int64) *world.Block {
	index := x + z*16 + y*16*16

	paletteIndex := c.Data[index]

	return &c.Palette[paletteIndex]
}

func (c *ChunkSection) SetBlock(x, y, z int64, item world.Block) {
	index := x + z*16 + y*16*16

	var paletteIndex int16 = -1

	for k, v := range c.Palette {
		if v.Name != item.Name || !reflect.DeepEqual(v.Properties, item.Properties) {
			continue
		}

		paletteIndex = int16(k)

		break
	}

	// oldValue := c.Palette[c.Data[index]]

	if paletteIndex >= 0 {
		c.Data[index] = paletteIndex
	} else {
		c.Data[index] = int16(len(c.Palette))

		c.Palette = append(c.Palette, item)
	}

	/* c.blockUpdateChan <- world.BlockUpdate{
		X:        x,
		Y:        y,
		Z:        z,
		OldValue: oldValue,
		NewValue: c.Palette[c.Data[index]],
	} */
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

var _ world.ChunkSection = &ChunkSection{}
