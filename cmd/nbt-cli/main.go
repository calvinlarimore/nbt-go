package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/calvinlarimore/nbt-go"
)

var names = [...]string{
	"End",
	"Byte",
	"Short",
	"Int",
	"Long",
	"Float",
	"Double",
	"Byte_Array",
	"String",
	"List",
	"Compound",
	"Int_Array",
	"Long_Array",
}

var file string

func main() {
	flags := flag.NewFlagSet("flags", flag.ExitOnError)

	file = os.Args[1]

	path := ""
	flags.StringVar(&path, "path", path, "NBT Path to work on")
	flags.StringVar(&path, "p", path, "Alias of -path")

	raw := false
	flags.BoolVar(&raw, "raw", raw, "Whether or not to print tag type, key, etc.")
	flags.BoolVar(&raw, "r", raw, "Alias of -raw")

	val := ""
	flags.StringVar(&val, "set", val, "Whether or not to print tag type, key, etc.")
	flags.StringVar(&val, "s", val, "Alias of -set")

	flags.Parse(os.Args[2:])

	set := false
	flags.Visit(func(f *flag.Flag) {
		if f.Name == "set" || f.Name == "s" {
			set = true
		}
	})

	nbt := nbt.ReadFile(file)

	if set {
		parsePath(nbt, path)
		// set tag value here
	} else {
		printTag(parsePath(nbt, path), path, 0, raw)
	}
}

func parsePath(parent *nbt.CompoundTag, path string) nbt.Tag {
	if path == "" {
		return parent
	}

	splitPath := strings.Split(path, ".")

	segments := strings.Split(splitPath[0], "[")
	tag := parent.Get(segments[0])

	if tag.ID() == nbt.CompoundTagID {
		if len(splitPath) > 1 {
			compound := tag.(*nbt.CompoundTag)

			return parsePath(compound, strings.Join(splitPath[1:], "."))
		}
	}

	return tag
}

/*func parseIndex(segments []string) int {
	s := strings.Trim(segments[0], "]")
	index, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		log.Fatalln(fmt.Sprintf("Error: \"%s\" is not an integer!"))
	}

	return int(index)
}*/

func printTag(t nbt.Tag, key string, indentLevel int, raw bool) {
	if key == "" {
		key = file
	}

	id := t.ID()
	var data string

	switch id {
	case nbt.EndTagID:
		data = ""
	case nbt.ByteTagID:
		tag := t.(*nbt.ByteTag)
		data = fmt.Sprintf("0x%02x", tag.Get())
	case nbt.ShortTagID:
		tag := t.(*nbt.ShortTag)
		data = strconv.FormatInt(int64(tag.Get()), 10)
	case nbt.IntTagID:
		tag := t.(*nbt.IntTag)
		data = strconv.FormatInt(int64(tag.Get()), 10)
	case nbt.LongTagID:
		tag := t.(*nbt.LongTag)
		data = strconv.FormatInt(int64(tag.Get()), 10)
	case nbt.FloatTagID:
		tag := t.(*nbt.FloatTag)
		data = strconv.FormatFloat(float64(tag.Get()), 'e', -1, 32)
	case nbt.DoubleTagID:
		tag := t.(*nbt.DoubleTag)
		data = strconv.FormatFloat(float64(tag.Get()), 'e', -1, 64)
	case nbt.StringTagID:
		tag := t.(*nbt.StringTag)
		data = fmt.Sprintf("'%s'", tag.Get())
	case nbt.ByteArrayTagID:
		tag := t.(*nbt.ByteArrayTag)
		data = fmt.Sprintf("%d elements [", len(tag.GetAll()))

		for i := range tag.GetAll() {
			val := tag.GetAll()[i]

			sep := ""
			if i+1 != len(tag.GetAll()) {
				sep = ", "
			}

			data = string(fmt.Appendf([]byte(data), "0x%02x%s", val, sep))
		}

		data = string(fmt.Append([]byte(data), "]"))
	case nbt.IntArrayTagID:
		tag := t.(*nbt.IntArrayTag)
		data = fmt.Sprintf("%d elements [", len(tag.GetAll()))

		for i := range tag.GetAll() {
			val := tag.GetAll()[i]

			sep := ""
			if i+1 != len(tag.GetAll()) {
				sep = ", "
			}

			data = string(fmt.Appendf([]byte(data), "%d%s", val, sep))
		}

		data = string(fmt.Append([]byte(data), "]"))
	case nbt.LongArrayTagID:
		tag := t.(*nbt.LongArrayTag)

		data = fmt.Sprintf("%d elements [", len(tag.GetAll()))

		for i := range tag.GetAll() {
			val := tag.GetAll()[i]

			sep := ""
			if i+1 != len(tag.GetAll()) {
				sep = ", "
			}

			data = string(fmt.Appendf([]byte(data), "%d%s", val, sep))
		}

		data = string(fmt.Append([]byte(data), "]"))
	case nbt.ListTagID:
		tag := t.(*nbt.ListTag)
		data = fmt.Sprintf("%d entries", len(tag.GetAll()))
		// TODO: Print List
	case nbt.CompoundTagID:
		tag := t.(*nbt.CompoundTag)
		data = fmt.Sprintf("%d entries", len(tag.GetAll()))
	}

	indent := strings.Repeat("\t", indentLevel)

	if !raw {
		fmt.Printf("%sTAG_%s('%s'): %s\n", indent, names[id], key, data)
	} else {
		fmt.Print(data)
	}

	if id == nbt.CompoundTagID {
		tag := t.(*nbt.CompoundTag)

		fmt.Printf("%s{\n", indent)

		for key := range tag.GetAll() {
			child := tag.GetAll()[key]

			printTag(child, key, indentLevel+1, raw)
		}

		fmt.Printf("%s}\n", indent)

	} else if id == nbt.ListTagID {
		tag := t.(*nbt.ListTag)

		fmt.Printf("%s[\n", indent)

		for i := range tag.GetAll() {
			child := tag.GetAll()[i]

			printTag(child, "N/A", indentLevel+1, raw)
		}

		fmt.Printf("%s]\n", indent)

	}
}
