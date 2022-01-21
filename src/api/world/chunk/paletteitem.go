package chunk

import (
	"encoding/binary"
	"encoding/json"
	"io"
)

type PaletteItem struct {
	ID         uint32                 `nbt:"id"`
	Name       string                 `nbt:"name"`
	Properties map[string]interface{} `nbt:"properties"`
}

func (p PaletteItem) Encode(w io.Writer) error {
	if err := binary.Write(w, binary.BigEndian, uint16(len(p.Name))); err != nil {
		return err
	}

	if _, err := w.Write([]byte(p.Name)); err != nil {
		return err
	}

	properties, err := json.Marshal(p.Properties)

	if err != nil {
		return err
	}

	if err := binary.Write(w, binary.BigEndian, uint32(len(properties))); err != nil {
		return err
	}

	_, err = w.Write(properties)

	return err
}

func (p *PaletteItem) Decode(r io.Reader) error {
	var nameLength uint16
	var propertiesLength uint32

	if err := binary.Read(r, binary.BigEndian, &nameLength); err != nil {
		return err
	}

	nameData := make([]byte, nameLength)

	if _, err := r.Read(nameData); err != nil {
		return err
	}

	p.Name = string(nameData)

	if err := binary.Read(r, binary.BigEndian, &propertiesLength); err != nil {
		return err
	}

	propertiesData := make([]byte, propertiesLength)

	if _, err := r.Read(propertiesData); err != nil {
		return err
	}

	properties := make(map[string]interface{})

	if err := json.Unmarshal(propertiesData, &properties); err != nil {
		return err
	}

	p.Properties = properties

	return nil
}
