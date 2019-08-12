package generator



var (
	varTypeToProto3Map = map[string]string{
		"float64": "double",
		"float32": "float32",
		"int":     "int32",
		"int64":   "int64",
		"int32":   "int32",
		"int16":   "int32",
		"int8":    "int32",
		"uint64":  "uint64",
		"uint32":  "uint32",
		"uint16":  "uint32",
		"uint8":   "uint32",
		"bool":    "bool",
		"string":  "string",
	}
)