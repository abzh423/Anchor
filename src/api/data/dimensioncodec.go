package data

type DimensionCodec struct {
	DimensionType DimensionCodecType          `nbt:"minecraft:dimension_type"`
	WorldGenBiome DimensionCodecWorldGenBiome `nbt:"minecraft:worldgen/biome"`
}

type DimensionCodecType struct {
	Type  string                    `nbt:"minecraft:dimension_type"`
	Value []DimensionCodecTypeValue `nbt:"value"`
}

type DimensionCodecTypeValue struct {
	Name    string                         `nbt:"name"`
	ID      int32                          `nbt:"id"`
	Element DimensionCodecTypeValueElement `nbt:"element"`
}

type DimensionCodecTypeValueElement struct {
	PiglinSafe         byte    `nbt:"piglin_safe"`
	Natural            byte    `nbt:"natural"`
	AmbientLight       float32 `nbt:"ambient_light"`
	FixedTime          int64   `nbt:"fixed_time"`
	InfiniteBurn       string  `nbt:"infiniburn"`
	RespawnAnchorWorks byte    `nbt:"respawn_anchor_works"`
	HasSkylight        byte    `nbt:"has_skylight"`
	BedWorks           byte    `nbt:"bed_works"`
	Effects            string  `nbt:"effects"`
	HasRaids           byte    `nbt:"has_raids"`
	MinY               int32   `nbt:"min_y"`
	Height             int32   `nbt:"height"`
	LogicalHeight      int32   `nbt:"logical_height"`
	CoordinateScale    float32 `nbt:"coordinate_scale"`
	Ultrawarm          byte    `nbt:"ultrawarm"`
	HasCeiling         byte    `nbt:"has_ceiling"`
}

type DimensionCodecWorldGenBiome struct {
	Type  string `nbt:"minecraft:worldgen/biome"`
	Value []map[string]interface{}
}
