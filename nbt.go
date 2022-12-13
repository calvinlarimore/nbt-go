package nbt

import (
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"io"
	"math"
	"os"
)

type TagID byte

const (
	EndTagID = TagID(iota)
	ByteTagID
	ShortTagID
	IntTagID
	LongTagID
	FloatTagID
	DoubleTagID
	ByteArrayTagID
	StringTagID
	ListTagID
	CompoundTagID
	IntArrayTagID
	LongArrayTagID
)

type Tag interface {
	ID() TagID
	Read(r io.Reader)
	Write(w io.Writer)
}

// TAG_End
type EndTag struct{}

func (EndTag) ID() TagID { return EndTagID }

func (t *EndTag) Read(r io.Reader) {}
func (t EndTag) Write(w io.Writer) {}

// TAG_Byte
type ByteTag struct{ val int8 }

func (ByteTag) ID() TagID { return ByteTagID }

func (t *ByteTag) Read(r io.Reader) {
	b := make([]byte, 1)
	r.Read(b)
	t.Set(int8(b[0]))
}
func (t ByteTag) Write(w io.Writer) {
	b := make([]byte, 1)
	b[0] = byte(t.Get())
	w.Write(b)
}

func (t ByteTag) Get() int8     { return t.val }
func (t *ByteTag) Set(val int8) { t.val = val }

func CreateByteTag(val int8) *ByteTag {
	return &ByteTag{val: val}
}

// TAG_Short
type ShortTag struct{ val int16 }

func (ShortTag) ID() TagID { return ShortTagID }

func (t *ShortTag) Read(r io.Reader) {
	b := make([]byte, 2)
	r.Read(b)
	t.Set(int16(binary.BigEndian.Uint16(b)))
}
func (t ShortTag) Write(w io.Writer) {
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, uint16(t.Get()))
	w.Write(b)
}

func (t ShortTag) Get() int16     { return t.val }
func (t *ShortTag) Set(val int16) { t.val = val }

func CreateShortTag(val int16) *ShortTag {
	return &ShortTag{val: val}
}

// TAG_Int
type IntTag struct{ val int32 }

func (IntTag) ID() TagID { return IntTagID }

func (t *IntTag) Read(r io.Reader) {
	b := make([]byte, 4)
	r.Read(b)
	t.Set(int32(binary.BigEndian.Uint32(b)))
}
func (t IntTag) Write(w io.Writer) {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, uint32(t.Get()))
	w.Write(b)
}

func (t IntTag) Get() int32     { return t.val }
func (t *IntTag) Set(val int32) { t.val = val }

func CreateIntTag(val int32) *IntTag {
	return &IntTag{val: val}
}

// TAG_Long
type LongTag struct{ val int64 }

func (LongTag) ID() TagID { return LongTagID }

func (t *LongTag) Read(r io.Reader) {
	b := make([]byte, 8)
	r.Read(b)
	t.Set(int64(binary.BigEndian.Uint64(b)))
}
func (t LongTag) Write(w io.Writer) {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(t.Get()))
	w.Write(b)
}

func (t LongTag) Get() int64     { return t.val }
func (t *LongTag) Set(val int64) { t.val = val }

func CreateLongTag(val int64) *LongTag {
	return &LongTag{val: val}
}

// TAG_Float
type FloatTag struct{ val float32 }

func (FloatTag) ID() TagID { return FloatTagID }

func (t *FloatTag) Read(r io.Reader) {
	b := make([]byte, 4)
	r.Read(b)
	t.Set(math.Float32frombits(binary.BigEndian.Uint32(b)))
}
func (t FloatTag) Write(w io.Writer) {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, math.Float32bits(t.Get()))
	w.Write(b)
}

func (t FloatTag) Get() float32     { return t.val }
func (t *FloatTag) Set(val float32) { t.val = val }

func CreateFloatTag(val float32) *FloatTag {
	return &FloatTag{val: val}
}

// TAG_Double
type DoubleTag struct{ val float64 }

func (DoubleTag) ID() TagID { return DoubleTagID }

func (t *DoubleTag) Read(r io.Reader) {
	b := make([]byte, 8)
	r.Read(b)
	t.Set(math.Float64frombits(binary.BigEndian.Uint64(b)))
}
func (t DoubleTag) Write(w io.Writer) {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, math.Float64bits(t.Get()))
	w.Write(b)
}

func (t DoubleTag) Get() float64     { return t.val }
func (t *DoubleTag) Set(val float64) { t.val = val }

func CreateDoubleTag(val float64) *DoubleTag {
	return &DoubleTag{val: val}
}

// TAG_String
type StringTag struct{ val string }

func (StringTag) ID() TagID { return StringTagID }

func (t *StringTag) Read(r io.Reader) {
	lenBuf := make([]byte, 2)
	r.Read(lenBuf)
	len := binary.BigEndian.Uint16(lenBuf)

	b := make([]byte, len)
	r.Read(b)

	t.Set(string(b))
}

func (t StringTag) Write(w io.Writer) {
	b := make([]byte, 2)
	len := uint16(len(t.Get()))
	binary.BigEndian.PutUint16(b, len)
	w.Write(b)

	w.Write([]byte(t.Get()))
}

func (t StringTag) Get() string     { return t.val }
func (t *StringTag) Set(val string) { t.val = val }

func CreateStringTag(val string) *StringTag {
	return &StringTag{val: val}
}

// TAG_Byte_Array
type ByteArrayTag struct{ val []byte }

func (ByteArrayTag) ID() TagID { return ByteArrayTagID }

func (t *ByteArrayTag) Read(r io.Reader) {
	lenBuf := make([]byte, 4)
	r.Read(lenBuf)
	len := binary.BigEndian.Uint32(lenBuf)

	b := make([]byte, len)
	r.Read(b)

	t.SetAll(b)
}
func (t ByteArrayTag) Write(w io.Writer) {
	b := make([]byte, 4)
	len := uint32(len(t.GetAll()))
	binary.BigEndian.PutUint32(b, len)
	w.Write(b)

	w.Write(t.GetAll())
}

func (t *ByteArrayTag) GetAll() []byte      { return t.val }
func (t *ByteArrayTag) SetAll(val []byte)   { t.val = val }
func (t *ByteArrayTag) Get(i int) byte      { return t.val[i] }
func (t *ByteArrayTag) Set(i int, val byte) { t.val[i] = val }
func (t *ByteArrayTag) Append(val byte)     { t.val = append(t.val, val) }
func (t *ByteArrayTag) Remove(i int)        { t.val = append(t.val[:i], t.val[i+1:]...) }

func CreateByteArrayTag(val []byte) *ByteArrayTag {
	return &ByteArrayTag{val: val}
}

// TAG_Int_Array
type IntArrayTag struct{ val []int32 }

func (IntArrayTag) ID() TagID { return IntArrayTagID }

func (t *IntArrayTag) Read(r io.Reader) {
	lenBuf := make([]byte, 4)
	r.Read(lenBuf)
	len := binary.BigEndian.Uint32(lenBuf)

	for i := 0; i < int(len); i++ {
		b := make([]byte, 4)
		r.Read(b)
		t.Append(int32(binary.BigEndian.Uint32(b)))
	}
}
func (t IntArrayTag) Write(w io.Writer) {
	b := make([]byte, 4)
	len := uint32(len(t.GetAll()))
	binary.BigEndian.PutUint32(b, len)
	w.Write(b)

	for i := 0; i < int(len); i++ {
		b := make([]byte, 4)
		binary.BigEndian.PutUint32(b, uint32(t.Get(i)))
		w.Write(b)
	}
}

func (t *IntArrayTag) GetAll() []int32      { return t.val }
func (t *IntArrayTag) SetAll(val []int32)   { t.val = val }
func (t *IntArrayTag) Get(i int) int32      { return t.val[i] }
func (t *IntArrayTag) Set(i int, val int32) { t.val[i] = val }
func (t *IntArrayTag) Append(val int32)     { t.val = append(t.val, val) }
func (t *IntArrayTag) Remove(i int)         { t.val = append(t.val[:i], t.val[i+1:]...) }

func CreateIntArraytag(val []int32) *IntArrayTag {
	return &IntArrayTag{val: val}
}

// TAG_Long_Array
type LongArrayTag struct{ val []int64 }

func (LongArrayTag) ID() TagID { return LongArrayTagID }

func (t *LongArrayTag) Read(r io.Reader) {
	lenBuf := make([]byte, 4)
	r.Read(lenBuf)
	len := binary.BigEndian.Uint32(lenBuf)

	for i := 0; i < int(len); i++ {
		b := make([]byte, 8)
		r.Read(b)
		t.Append(int64(binary.BigEndian.Uint64(b)))
	}
}
func (t LongArrayTag) Write(w io.Writer) {
	b := make([]byte, 4)
	len := uint32(len(t.GetAll()))
	binary.BigEndian.PutUint32(b, len)
	w.Write(b)

	for i := 0; i < int(len); i++ {
		b := make([]byte, 8)
		binary.BigEndian.PutUint64(b, uint64(t.Get(i)))
		w.Write(b)
	}
}

func (t *LongArrayTag) GetAll() []int64      { return t.val }
func (t *LongArrayTag) SetAll(val []int64)   { t.val = val }
func (t *LongArrayTag) Get(i int) int64      { return t.val[i] }
func (t *LongArrayTag) Set(i int, val int64) { t.val[i] = val }
func (t *LongArrayTag) Append(val int64)     { t.val = append(t.val, val) }
func (t *LongArrayTag) Remove(i int)         { t.val = append(t.val[:i], t.val[i+1:]...) }

func CreateLongArrayTag(val []int64) *LongArrayTag {
	return &LongArrayTag{val: val}
}

// TAG_List
type ListTag struct {
	tags []Tag
	id   TagID
}

func (ListTag) ID() TagID { return ListTagID }

func (t *ListTag) Read(r io.Reader) {
	idBuf := make([]byte, 1)
	r.Read(idBuf)
	id := TagID(idBuf[0])
	t.id = id

	lenBuf := make([]byte, 4)
	r.Read(lenBuf)
	len := binary.BigEndian.Uint32(lenBuf)

	for i := 0; i < int(len); i++ {
		var tag Tag

		switch id {
		case ByteTagID:
			tag = &ByteTag{}
		case ShortTagID:
			tag = &ShortTag{}
		case IntTagID:
			tag = &IntTag{}
		case LongTagID:
			tag = &LongTag{}
		case FloatTagID:
			tag = &FloatTag{}
		case DoubleTagID:
			tag = &DoubleTag{}
		case StringTagID:
			tag = &StringTag{}
		case ByteArrayTagID:
			tag = &ByteArrayTag{}
		case IntArrayTagID:
			tag = &IntArrayTag{}
		case LongArrayTagID:
			tag = &LongArrayTag{}
		case ListTagID:
			tag = &ListTag{}
		case CompoundTagID:
			tag = CreateCompoundTag()
		}

		tag.Read(r)

		t.Append(tag)
	}
}
func (t ListTag) Write(w io.Writer) {
	idBuf := make([]byte, 1)
	idBuf[0] = byte(t.id)
	w.Write(idBuf)

	lenBuf := make([]byte, 4)
	len := uint32(len(t.tags))
	binary.BigEndian.PutUint32(lenBuf, len)
	w.Write(lenBuf)

	for i := 0; i < int(len); i++ {
		tag := t.Get(i)

		tag.Write(w)
	}
}

func (t *ListTag) GetAll() []Tag      { return t.tags }
func (t *ListTag) SetAll(tag []Tag)   { t.tags = tag }
func (t *ListTag) Get(i int) Tag      { return t.tags[i] }
func (t *ListTag) Set(i int, tag Tag) { t.tags[i] = tag }
func (t *ListTag) Append(tag Tag)     { t.tags = append(t.tags, tag) }
func (t *ListTag) Remove(i int)       { t.tags = append(t.tags[:i], t.tags[i+1:]...) }

func CreateListTag(id TagID, tags []Tag) *ListTag {
	return &ListTag{
		tags: tags,
		id:   id,
	}
}

// TAG_Compound
type CompoundTag struct {
	tags     map[string]Tag
	implicit bool
}

func (CompoundTag) ID() TagID { return CompoundTagID }

func (t *CompoundTag) Read(r io.Reader) {
	readName := func() string {
		tag := CreateStringTag("")
		tag.Read(r)
		return tag.Get()
	}

	for {
		b := make([]byte, 1)
		r.Read(b)
		id := TagID(b[0])

		if id == EndTagID {
			return
		}

		var tag Tag

		switch id {
		case ByteTagID:
			tag = &ByteTag{}
		case ShortTagID:
			tag = &ShortTag{}
		case IntTagID:
			tag = &IntTag{}
		case LongTagID:
			tag = &LongTag{}
		case FloatTagID:
			tag = &FloatTag{}
		case DoubleTagID:
			tag = &DoubleTag{}
		case StringTagID:
			tag = &StringTag{}
		case ByteArrayTagID:
			tag = &ByteArrayTag{}
		case IntArrayTagID:
			tag = &IntArrayTag{}
		case LongArrayTagID:
			tag = &LongArrayTag{}
		case ListTagID:
			tag = &ListTag{}
		case CompoundTagID:
			tag = CreateCompoundTag()
		}

		name := readName()

		tag.Read(r)

		t.Put(name, tag)
	}
}

func (t CompoundTag) Write(w io.Writer) {
	for key := range t.tags {
		tag := t.tags[key]

		b := make([]byte, 1)
		b[0] = byte(tag.ID())
		w.Write(b)

		nameTag := CreateStringTag(key)
		nameTag.Write(w)

		tag.Write(w)
	}

	if !t.implicit {
		b := make([]byte, 1)
		b[0] = byte(EndTagID)
		w.Write(b)
	}
}

func (t *CompoundTag) GetAll() map[string]Tag     { return t.tags }
func (t *CompoundTag) SetAll(tags map[string]Tag) { t.tags = tags }
func (t *CompoundTag) Get(key string) Tag         { return t.tags[key] }
func (t *CompoundTag) Put(key string, tag Tag)    { t.tags[key] = tag }
func (t *CompoundTag) Remove(key string)          { delete(t.tags, key) }

func (t *CompoundTag) Contains(key string) bool {
	for k := range t.tags {
		if k == key {
			return true
		}
	}

	return false
}

func CreateCompoundTag() *CompoundTag {
	return &CompoundTag{
		tags:     make(map[string]Tag),
		implicit: false,
	}
}

func CreateImplicitCompoundTag() *CompoundTag {
	return &CompoundTag{
		tags:     make(map[string]Tag),
		implicit: true,
	}
}

func (t *CompoundTag) GetByte(key string) (i int8, ok bool) {
	tag, okay := t.Get(key).(*ByteTag)

	if okay {
		return tag.Get(), true
	}

	return 0, false
}

func (t *CompoundTag) PutByte(key string, val int8) {
	tag := ByteTag{
		val: val,
	}

	t.Put(key, &tag)
}

func (t *CompoundTag) GetShort(key string) (i int16, ok bool) {
	tag, okay := t.Get(key).(*ShortTag)

	if okay {
		return tag.Get(), true
	}

	return 0, false
}

func (t *CompoundTag) PutShort(key string, val int16) {
	tag := ShortTag{
		val: val,
	}

	t.Put(key, &tag)
}

func (t *CompoundTag) GetInt(key string) (i int32, ok bool) {
	tag, okay := t.Get(key).(*IntTag)

	if okay {
		return tag.Get(), true
	}

	return 0, false
}

func (t *CompoundTag) PutInt(key string, val int32) {
	tag := IntTag{
		val: val,
	}

	t.Put(key, &tag)
}

func (t *CompoundTag) GetLong(key string) (i int64, ok bool) {
	tag, okay := t.Get(key).(*LongTag)

	if okay {
		return tag.Get(), true
	}

	return 0, false
}

func (t *CompoundTag) PutLong(key string, val int64) {
	tag := LongTag{
		val: val,
	}

	t.Put(key, &tag)
}

func (t *CompoundTag) GetFloat(key string) (f float32, ok bool) {
	tag, okay := t.Get(key).(*FloatTag)

	if okay {
		return tag.Get(), true
	}

	return 0, false
}

func (t *CompoundTag) PutFloat(key string, val float32) {
	tag := FloatTag{
		val: val,
	}

	t.Put(key, &tag)
}

func (t *CompoundTag) GetDouble(key string) (f float64, ok bool) {
	tag, okay := t.Get(key).(*DoubleTag)

	if okay {
		return tag.Get(), true
	}

	return 0, false
}

func (t *CompoundTag) PutDouble(key string, val float64) {
	tag := DoubleTag{
		val: val,
	}

	t.Put(key, &tag)
}

func (t *CompoundTag) GetString(key string) (s string, ok bool) {
	tag, okay := t.Get(key).(*StringTag)

	if okay {
		return tag.Get(), true
	}

	return "", false
}

func (t *CompoundTag) PutString(key string, val string) {
	tag := CreateStringTag(val)

	t.Put(key, tag)
}

// TODO: Arrays

func (t *CompoundTag) GetList(key string) (l *ListTag, ok bool) {
	tag, okay := t.Get(key).(*ListTag)

	if okay {
		return tag, true
	}

	return nil, false
}

func (t *CompoundTag) GetCompound(key string) (c *CompoundTag, ok bool) {
	tag, okay := t.Get(key).(*CompoundTag)

	if okay {
		return tag, true
	}

	return nil, false
}

func ReadFile(name string) *CompoundTag {
	data, _ := os.ReadFile(name)

	tag := CreateImplicitCompoundTag()

	gz, err := gzip.NewReader(bytes.NewBuffer(data))

	if err == nil {
		tag.Read(gz)
	} else {
		buf := bytes.NewBuffer(data)
		tag.Read(buf)
	}

	return tag
}
