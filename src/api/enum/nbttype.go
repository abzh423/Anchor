package enum

type NBTType byte

const (
	NBTTagEnd NBTType = iota
	NBTTagByte
	NBTTagShort
	NBTTagInt
	NBTTagLong
	NBTTagFloat
	NBTTagDouble
	NBTTagByteArray
	NBTTagString
	NBTTagList
	NBTTagCompound
	NBTTagIntArray
	NBTTagLongArray
)
