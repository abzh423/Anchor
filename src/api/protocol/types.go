package protocol

import (
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"io"
	"math"
	"reflect"
	"regexp"

	"github.com/Tnze/go-mc/nbt"
)

var (
	colorFormattingCodes = map[string]string{
		"black":         "0",
		"dark_blue":     "1",
		"dark_green":    "2",
		"dark_aqua":     "3",
		"dark_red":      "4",
		"dark_purple":   "5",
		"gold":          "6",
		"gray":          "7",
		"dark_gray":     "8",
		"blue":          "9",
		"green":         "a",
		"aqua":          "b",
		"red":           "c",
		"light_purple":  "d",
		"yellow":        "e",
		"white":         "f",
		"minecoin_gold": "g",
	}
)

type DataType interface {
	DataTypeReader
	DataTypeWriter
}

type DataTypeReader interface {
	Decode(r io.Reader) (int64, error)
}

type DataTypeWriter interface {
	Encode(w io.Writer) (int64, error)
}

type (
	Boolean       bool
	Byte          int8
	UnsignedByte  uint8
	Short         int16
	UnsignedShort uint16
	ShortLE       int16
	Int           int32
	UnsignedInt   uint32
	Long          int64
	UnsignedLong  uint64
	Float         float32
	Double        float64
	String        string
	StringNT      string
	Chat          struct {
		Text          string `json:"text"`
		Color         string `json:"color"`
		Bold          bool   `json:"bold,omitempty"`
		Italic        bool   `json:"italic,omitempty"`
		Underlined    bool   `json:"underlined,omitempty"`
		Strikethrough bool   `json:"strikethrough,omitempty"`
		Obfuscated    bool   `json:"obfuscated,omitempty"`
		Extra         []Chat `json:"extra,omitempty"`
	}
	Identifier       string
	VarInt           int32
	VarLong          int64
	RelativePosition struct {
		X, Y, Z int64
	}
	AbsolutePosition struct {
		X, Y, Z float64
	}
	Angle     byte
	UUID      string
	ByteArray []byte
	Array     []DataType
	NBT       struct {
		Value   interface{}
		RootTag string
	}
	PackedArray struct {
		Value        []int64
		BitsPerEntry int
	}
)

var (
	uuidReplaceRegExp = regexp.MustCompile("[^A-Fa-f0-9]")
)

func (v Boolean) Encode(w io.Writer) (int64, error) {
	var value byte = 0

	if v {
		value = 1
	}

	n, err := w.Write([]byte{value})

	return int64(n), err
}

func (v *Boolean) Decode(r io.Reader) (int64, error) {
	data := make([]byte, 1)

	n, err := r.Read(data)

	*v = data[0] == 1

	return int64(n), err
}

func (v Byte) Encode(w io.Writer) (int64, error) {
	n, err := w.Write([]byte{byte(v)})

	return int64(n), err
}

func (v *Byte) Decode(r io.Reader) (int64, error) {
	data := make([]byte, 1)

	n, err := r.Read(data)

	*v = Byte(data[0])

	return int64(n), err
}

func (v UnsignedByte) Encode(w io.Writer) (int64, error) {
	n, err := w.Write([]byte{byte(v)})

	return int64(n), err
}

func (v *UnsignedByte) Decode(r io.Reader) (int64, error) {
	data := make([]byte, 1)

	n, err := r.Read(data)

	*v = UnsignedByte(data[0])

	return int64(n), err
}

func (v Short) Encode(w io.Writer) (int64, error) {
	return 2, binary.Write(w, binary.BigEndian, int16(v))
}

func (v *Short) Decode(r io.Reader) (int64, error) {
	var value int16

	err := binary.Read(r, binary.BigEndian, &value)

	*v = Short(value)

	return 2, err
}

func (v UnsignedShort) Encode(w io.Writer) (int64, error) {
	return 2, binary.Write(w, binary.BigEndian, uint16(v))
}

func (v *UnsignedShort) Decode(r io.Reader) (int64, error) {
	var value uint16

	err := binary.Read(r, binary.BigEndian, &value)

	*v = UnsignedShort(value)

	return 2, err
}

func (v ShortLE) Encode(w io.Writer) (int64, error) {
	return 2, binary.Write(w, binary.LittleEndian, uint16(v))
}

func (v *ShortLE) Decode(r io.Reader) (int64, error) {
	var value uint16

	err := binary.Read(r, binary.LittleEndian, &value)

	*v = ShortLE(value)

	return 2, err
}

func (v Int) Encode(w io.Writer) (int64, error) {
	return 4, binary.Write(w, binary.BigEndian, int32(v))
}

func (v *Int) Decode(r io.Reader) (int64, error) {
	var value int32

	err := binary.Read(r, binary.BigEndian, &value)

	*v = Int(value)

	return 4, err
}

func (v UnsignedInt) Encode(w io.Writer) (int64, error) {
	return 4, binary.Write(w, binary.BigEndian, uint32(v))
}

func (v *UnsignedInt) Decode(r io.Reader) (int64, error) {
	var value uint32

	err := binary.Read(r, binary.BigEndian, &value)

	*v = UnsignedInt(value)

	return 4, err
}

func (v Long) Encode(w io.Writer) (int64, error) {
	return 8, binary.Write(w, binary.BigEndian, int64(v))
}

func (v *Long) Decode(r io.Reader) (int64, error) {
	var value int64

	err := binary.Read(r, binary.BigEndian, &value)

	*v = Long(value)

	return 8, err
}

func (v UnsignedLong) Encode(w io.Writer) (int64, error) {
	return 8, binary.Write(w, binary.BigEndian, uint64(v))
}

func (v *UnsignedLong) Decode(r io.Reader) (int64, error) {
	var value uint64

	err := binary.Read(r, binary.BigEndian, &value)

	*v = UnsignedLong(value)

	return 8, err
}

func (v Float) Encode(w io.Writer) (int64, error) {
	value := math.Float32bits(float32(v))

	return 4, binary.Write(w, binary.BigEndian, value)
}

func (v *Float) Decode(r io.Reader) (int64, error) {
	var value uint32

	err := binary.Read(r, binary.BigEndian, &value)

	*v = Float(math.Float32frombits(value))

	return 4, err
}

func (v Double) Encode(w io.Writer) (int64, error) {
	value := math.Float64bits(float64(v))

	return 8, binary.Write(w, binary.BigEndian, value)
}

func (v *Double) Decode(r io.Reader) (int64, error) {
	var value uint64

	err := binary.Read(r, binary.BigEndian, &value)

	*v = Double(math.Float64frombits(value))

	return 8, err
}

func (v String) Encode(w io.Writer) (int64, error) {
	data := []byte(v)

	n, err := WriteVarInt(int32(len(data)), w)

	if err != nil {
		return n, err
	}

	n2, err := w.Write(data)

	return n + int64(n2), err
}

func (v *String) Decode(r io.Reader) (int64, error) {
	length, n, err := ReadVarInt(r)

	if err != nil {
		return n, err
	}

	data := make([]byte, length)

	n2, err := r.Read(data)

	if err != nil {
		return int64(n2), err
	}

	*v = String(data)

	return n + int64(n2), nil
}

func (v StringNT) Encode(w io.Writer) (int64, error) {
	n, err := w.Write([]byte(v))

	if err != nil {
		return int64(n), err
	}

	n2, err := w.Write([]byte{0x00})

	return int64(n + n2), err
}

func (v *StringNT) Decode(r io.Reader) (int64, error) {
	value := ""

	var bytesRead int64 = 0

	for {
		data := make([]byte, 1)

		n, err := r.Read(data)

		if err != nil {
			return bytesRead, err
		}

		bytesRead += int64(n)

		if data[0] == 0x00 {
			break
		}

		value += string(data[0])
	}

	*v = StringNT(value)

	return bytesRead, nil
}

func (v Chat) Encode(w io.Writer) (int64, error) {
	data, err := json.Marshal(v)

	if err != nil {
		return 0, err
	}

	n, err := String(data).Encode(w)

	return int64(n), err
}

func (v *Chat) Decode(r io.Reader) (int64, error) {
	value := String("")

	n, err := value.Decode(r)

	if err != nil {
		return int64(n), err
	}

	return int64(n), json.Unmarshal([]byte(value), v)
}

func (v Chat) Format() string {
	value := ""

	if colorCode, ok := colorFormattingCodes[v.Color]; ok {
		value += "\u00A7" + colorCode
	}

	if v.Bold {
		value += "\u00A7l"
	}

	if v.Italic {
		value += "\u00A7o"
	}

	if v.Underlined {
		value += "\u00A7n"
	}

	if v.Strikethrough {
		value += "\u00A7m"
	}

	if v.Obfuscated {
		value += "\u00A7k"
	}

	value += v.Text

	for _, extra := range v.Extra {
		value += extra.Format()
	}

	return value
}

func (v Identifier) Encode(w io.Writer) (int64, error) {
	return String(v).Encode(w)
}

func (v *Identifier) Decode(r io.Reader) (int64, error) {
	value := String("")

	n, err := value.Decode(r)

	*v = Identifier(value)

	return n, err
}

func (v VarInt) Encode(w io.Writer) (int64, error) {
	return WriteVarInt(int32(v), w)
}

func (v *VarInt) Decode(r io.Reader) (int64, error) {
	value, n, err := ReadVarInt(r)

	*v = VarInt(value)

	return n, err
}

func (v VarLong) Encode(w io.Writer) (int64, error) {
	return WriteVarLong(int64(v), w)
}

func (v *VarLong) Decode(r io.Reader) (int64, error) {
	value, n, err := ReadVarLong(r)

	*v = VarLong(value)

	return n, err
}

func (v RelativePosition) Encode(w io.Writer) (int64, error) {
	return 8, binary.Write(w, binary.BigEndian, ((uint64(v.X)&0x3FFFFFF)<<38)|((uint64(v.Z)&0x3FFFFFF)<<12)|(uint64(v.Y)&0xFFF))
}

func (v *RelativePosition) Decode(r io.Reader) (int64, error) {
	var value uint64

	err := binary.Read(r, binary.BigEndian, &value)

	v.X = int64(value >> 38)
	v.Y = int64(value << 52 >> 52)
	v.Z = int64(value << 26 >> 38)

	return 8, err
}

func (v Angle) Encode(w io.Writer) (int64, error) {
	n, err := w.Write([]byte{byte(v)})

	return int64(n), err
}

func (v *Angle) Decode(r io.Reader) (int64, error) {
	data := make([]byte, 1)

	n, err := r.Read(data)

	*v = Angle(data[0])

	return int64(n), err
}

func (v UUID) Encode(w io.Writer) (int64, error) {
	value, err := hex.DecodeString(uuidReplaceRegExp.ReplaceAllString(string(v), ""))

	if err != nil {
		return 0, err
	}

	n, err := w.Write(value)

	return int64(n), err
}

func (v *UUID) Decode(r io.Reader) (int64, error) {
	data := make([]byte, 16)

	n, err := r.Read(data)

	if err != nil {
		return int64(n), err
	}

	*v = UUID(hex.EncodeToString(data))

	return int64(n), nil
}

func (v ByteArray) Encode(w io.Writer) (int64, error) {
	n, err := WriteVarInt(int32(len(v)), w)

	if err != nil {
		return n, err
	}

	n2, err := w.Write(v)

	return n + int64(n2), err
}

func (v *ByteArray) Decode(r io.Reader) (int64, error) {
	length, n, err := ReadVarInt(r)

	if err != nil {
		return n, err
	}

	data := make([]byte, length)

	n2, err := r.Read(data)

	*v = ByteArray(data)

	return n + int64(n2), err
}

func (v Array) Encode(w io.Writer) (int64, error) {
	n, err := WriteVarInt(int32(len(v)), w)

	if err != nil {
		return n, err
	}

	for _, value := range v {
		n2, err := value.Encode(w)

		if err != nil {
			return n, err
		}

		n += n2
	}

	return n, nil
}

func (v *Array) Decode(r io.Reader) (int64, error) {
	length, n, err := ReadVarInt(r)

	if err != nil {
		return n, err
	}

	s := reflect.TypeOf(*v).Elem()
	slice := reflect.MakeSlice(s, int(length), int(length))

	for k, value := range *v {
		n2, err := value.Decode(r)

		if err != nil {
			return n, err
		}

		n += n2

		slice.Index(k).Set(reflect.ValueOf(value))
	}

	return n, nil
}

func (v NBT) Encode(w io.Writer) (int64, error) {
	if v.Value == nil {
		n, err := w.Write([]byte{0})

		return int64(n), err
	}

	// TODO calculate bytes written for NBT tag

	return -1, nbt.NewEncoder(w).Encode(v.Value, v.RootTag)
}

func (v *NBT) Decode(r io.Reader) (int64, error) {
	// TODO NBT decoder

	return 0, nil
}
