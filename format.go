package skyencoder

import (
	"fmt"
)

func cast(typ, name string) string {
	return fmt.Sprintf("%s(%s)", typ, name)
}

// Options is parsed encoder struct tag options
type Options struct {
	OmitEmpty bool
	MaxLength uint64
}

/* Encode size */

func WrapEncodeSizeFunc(structName, structPackageName, counterName, funcBody string) []byte {
	structType := structName
	if structPackageName != "" {
		structType = fmt.Sprintf("%s.%s", structPackageName, structName)
	}
	return []byte(fmt.Sprintf(`
// EncodeSize%[1]s computes the size of an encoded object of type %[1]s
func EncodeSize%[1]s(obj *%[4]s) int {
	%[2]s := 0

	%[3]s

	return %[2]s
}
`, structName, counterName, funcBody, structType))
}

func BuildEncodeSizeBool(name, counterName string, options *Options) string {
	return fmt.Sprintf(`
		// %[1]s
		%[2]s++
	`, name, counterName)
}

func BuildEncodeSizeUint8(name, counterName string, options *Options) string {
	return fmt.Sprintf(`
		// %[1]s
		%[2]s++
	`, name, counterName)
}

func BuildEncodeSizeUint16(name, counterName string, options *Options) string {
	return fmt.Sprintf(`
		// %[1]s
		%[2]s += 2
	`, name, counterName)
}

func BuildEncodeSizeUint32(name, counterName string, options *Options) string {
	return fmt.Sprintf(`
		// %[1]s
		%[2]s += 4
	`, name, counterName)
}

func BuildEncodeSizeUint64(name, counterName string, options *Options) string {
	return fmt.Sprintf(`
		// %[1]s
		%[2]s += 8
	`, name, counterName)
}

func BuildEncodeSizeInt8(name, counterName string, options *Options) string {
	return fmt.Sprintf(`
		// %[1]s
		%[2]s++
	`, name, counterName)
}

func BuildEncodeSizeInt16(name, counterName string, options *Options) string {
	return fmt.Sprintf(`
		// %[1]s
		%[2]s += 2
	`, name, counterName)
}

func BuildEncodeSizeInt32(name, counterName string, options *Options) string {
	return fmt.Sprintf(`
		// %[1]s
		%[2]s += 4
	`, name, counterName)
}

func BuildEncodeSizeInt64(name, counterName string, options *Options) string {
	return fmt.Sprintf(`
		// %[1]s
		%[2]s += 8
	`, name, counterName)
}

func BuildEncodeSizeString(name, counterName string, options *Options) string {
	body := fmt.Sprintf(`
	// %[1]s
	%[2]s += 4 + len(%[1]s)
	`, name, counterName)

	if options != nil && options.OmitEmpty {
		return fmt.Sprintf(`
		// omitempty
		if len(%[1]s) != 0 {
			%[2]s
		}
		`, name, body)
	}

	return body
}

func BuildEncodeSizeByteArray(name, counterName string, length int64, options *Options) string {
	return fmt.Sprintf(`
	// %[1]s
	%[2]s += %[3]d
	`, name, counterName, length)
}

func BuildEncodeSizeArray(name, counterName, nextCounterName, elemVarName, elemSection string, length int64, isDynamic bool, options *Options) string {
	if isDynamic {
		return fmt.Sprintf(`
		// %[1]s
		for _, %[2]s := range %[1]s {
			%[4]s := 0

			%[3]s

			%[5]s += %[4]s
		}
		`, name, elemVarName, elemSection, nextCounterName, counterName)
	}

	return fmt.Sprintf(`
	// %[1]s
	{
		%[5]s := 0

		%[4]s

		%[2]s += %[3]d * %[5]s
	}
	`, name, counterName, length, elemSection, nextCounterName)

	// return fmt.Sprintf(`
	// // %[1]s
	// i += %[3]d * func() int {
	// 	i := 0

	// 	%[2]s

	// 	return i
	// }()
	// `, name, elemSection, length)
}

func BuildEncodeSizeByteSlice(name, counterName string, options *Options) string {
	body := fmt.Sprintf(`
	// %[1]s
	%[2]s += 4 + len(%[1]s)
	`, name, counterName)

	if options != nil && options.OmitEmpty {
		return fmt.Sprintf(`
		// omitempty
		if len(%[1]s) != 0 {
			%[2]s
		}
		`, name, body)
	}

	return body
}

func BuildEncodeSizeSlice(name, counterName, nextCounterName, elemVarName, elemSection string, isDynamic bool, options *Options) string {
	var body string

	debugPrintf("BuildEncodeSizeSlice: counterName=%s\n", counterName)

	if isDynamic {
		body = fmt.Sprintf(`
		// %[1]s
		%[2]s += 4
		for _, %[3]s := range %[1]s {
			%[5]s := 0

			%[4]s

			%[2]s += %[5]s
		}
		`, name, counterName, elemVarName, elemSection, nextCounterName)
	} else {
		body = fmt.Sprintf(`
		// %[1]s
		%[2]s += 4
		{
			%[4]s := 0

			%[3]s

			%[2]s += len(%[1]s) * %[4]s
		}
		`, name, counterName, elemSection, nextCounterName)

		// body = fmt.Sprintf(`
		// // %[1]s
		// i += 4
		// i += len(%[1]s) * func() int {
		// 	i := 0

		// 	%[2]s

		// 	return i
		// }()`, name, elemSection)
	}

	if options != nil && options.OmitEmpty {
		return fmt.Sprintf(`
		// omitempty
		if len(%[1]s) != 0 {
			%[2]s
		}
		`, name, body)
	}

	return body
}

func BuildEncodeSizeMap(name, counterName, nextCounterName, keyVarName, elemVarName, keySection, elemSection string, isDynamicKey, isDynamicElem bool, options *Options) string {
	var body string

	if isDynamicKey || isDynamicElem {
		if !isDynamicKey {
			keyVarName = "_"
		}

		if !isDynamicElem {
			elemVarName = "_"
		}

		body = fmt.Sprintf(`
		// %[1]s
		%[2]s += 4
		for %[3]s, %[4]s := range %[1]s {
			%[7]s := 0

			%[5]s

			%[6]s

			%[2]s += %[7]s
		}
		`, name, counterName, keyVarName, elemVarName, keySection, elemSection, nextCounterName)
	} else {
		body = fmt.Sprintf(`
		// %[1]s
		%[2]s += 4
		{
			%[5]s := 0

			%[3]s

			%[4]s

			%[2]s += len(%[1]s) * %[5]s
		}
		`, name, counterName, keySection, elemSection, nextCounterName)

		// body = fmt.Sprintf(`
		// // %[1]s
		// i += 4
		// i += len(%[1]s) * func() int {
		// 	i := 0

		// 	%[2]s

		// 	%[3]s

		// 	return i
		// }()
		// `, name, keySection, elemSection)
	}

	if options != nil && options.OmitEmpty {
		return fmt.Sprintf(`
		// omitempty
		if len(%[1]s) != 0 {
			%[2]s
		}
		`, name, body)
	}

	return body
}

/* Encode */

func WrapEncodeFunc(structName, structPackageName, funcBody string) []byte {
	structType := structName
	if structPackageName != "" {
		structType = fmt.Sprintf("%s.%s", structPackageName, structType)
	}
	return []byte(fmt.Sprintf(`
// Encode%[1]s encodes an object of type %[1]s to the buffer in encoder.Encoder
func Encode%[1]s(e *encoder.Encoder, obj *%[3]s) error {
	%[2]s

	return nil
}
`, structName, funcBody, structType))
}

func BuildEncodeBool(name string, castType bool, options *Options) string {
	castName := name
	if castType {
		castName = cast("bool", name)
	}
	return fmt.Sprintf(`
	// %[1]s
	e.Bool(%[2]s)
	`, name, castName)
}

func BuildEncodeUint8(name string, castType bool, options *Options) string {
	castName := name
	if castType {
		castName = cast("uint8", name)
	}
	return fmt.Sprintf(`
	// %[1]s
	e.Uint8(%[2]s)
	`, name, castName)
}

func BuildEncodeUint16(name string, castType bool, options *Options) string {
	castName := name
	if castType {
		castName = cast("uint16", name)
	}
	return fmt.Sprintf(`
	// %[1]s
	e.Uint16(%[2]s)
	`, name, castName)
}

func BuildEncodeUint32(name string, castType bool, options *Options) string {
	castName := name
	if castType {
		castName = cast("uint32", name)
	}
	return fmt.Sprintf(`
	// %[1]s
	e.Uint32(%[2]s)
	`, name, castName)
}

func BuildEncodeUint64(name string, castType bool, options *Options) string {
	castName := name
	if castType {
		castName = cast("uint64", name)
	}
	return fmt.Sprintf(`
	// %[1]s
	e.Uint64(%[2]s)
	`, name, castName)
}

func BuildEncodeInt8(name string, castType bool, options *Options) string {
	castName := name
	if castType {
		castName = cast("int8", name)
	}
	return fmt.Sprintf(`
	// %[1]s
	e.Int8(%[2]s)
	`, name, castName)
}

func BuildEncodeInt16(name string, castType bool, options *Options) string {
	castName := name
	if castType {
		castName = cast("int16", name)
	}
	return fmt.Sprintf(`
	// %[1]s
	e.Int16(%[2]s)
	`, name, castName)
}

func BuildEncodeInt32(name string, castType bool, options *Options) string {
	castName := name
	if castType {
		castName = cast("int32", name)
	}
	return fmt.Sprintf(`
	// %[1]s
	e.Int32(%[2]s)
	`, name, castName)
}

func BuildEncodeInt64(name string, castType bool, options *Options) string {
	castName := name
	if castType {
		castName = cast("int64", name)
	}
	return fmt.Sprintf(`
	// %[1]s
	e.Int64(%[2]s)
	`, name, castName)
}

func BuildEncodeString(name string, options *Options) string {
	body := fmt.Sprintf(`
	%[2]s

	// %[1]s
	e.ByteSlice([]byte(%[1]s))
	`, name, encodeMaxLengthCheck(name, options))

	if options != nil && options.OmitEmpty {
		return fmt.Sprintf(`
			// omitempty
			if len(%[1]s) != 0 {
				%[2]s
			}
		`, name, body)
	}

	return body
}

func BuildEncodeByteArray(name string, options *Options) string {
	return fmt.Sprintf(`
	// %[1]s
	e.CopyBytes(%[1]s[:])
	`, name)
}

func BuildEncodeArray(name, elemVarName, elemSection string, options *Options) string {
	return fmt.Sprintf(`
	// %[1]s
	for _, %[2]s := range %[1]s {
		%[3]s
	}
	`, name, elemVarName, elemSection)
}

func BuildEncodeByteSlice(name string, options *Options) string {
	body := fmt.Sprintf(`
	%[2]s

	// %[1]s length check
	if len(%[1]s) > math.MaxUint32 {
		return errors.New("%[1]s length exceeds math.MaxUint32")
	}

	// %[1]s length
	e.Uint32(uint32(len(%[1]s)))

	// %[1]s copy
	e.CopyBytes(%[1]s)
	`, name, encodeMaxLengthCheck(name, options))

	if options != nil && options.OmitEmpty {
		return fmt.Sprintf(`
			// omitempty
			if len(%[1]s) != 0 {
				%[2]s
			}
		`, name, body)
	}

	return body
}

func BuildEncodeSlice(name, elemVarName, elemSection string, options *Options) string {
	body := fmt.Sprintf(`
	%[4]s

	// %[1]s length check
	if len(%[1]s) > math.MaxUint32 {
		return errors.New("%[1]s length exceeds math.MaxUint32")
	}

	// %[1]s length
	e.Uint32(uint32(len(%[1]s)))

	// %[1]s
	for _, %[2]s := range %[1]s {
		%[3]s
	}
	`, name, elemVarName, elemSection, encodeMaxLengthCheck(name, options))

	if options != nil && options.OmitEmpty {
		return fmt.Sprintf(`
			// omitempty
			if len(%[1]s) != 0 {
				%[2]s
			}
		`, name, body)
	}

	return body
}

func BuildEncodeMap(name, keyVarName, elemVarName, keySection, elemSection string, options *Options) string {
	if keySection == "" {
		keyVarName = "_"
	}
	if elemSection == "" {
		elemVarName = "_"
	}

	body := fmt.Sprintf(`
	// %[1]s

	%[6]s

	// %[1]s length check
	if len(%[1]s) > math.MaxUint32 {
		return errors.New("%[1]s length exceeds math.MaxUint32")
	}

	// %[1]s length
	e.Uint32(uint32(len(%[1]s)))

	for %[2]s, %[3]s := range %[1]s {
		%[4]s

		%[5]s
	}
	`, name, keyVarName, elemVarName, keySection, elemSection, encodeMaxLengthCheck(name, options))

	if options != nil && options.OmitEmpty {
		return fmt.Sprintf(`
			// omitempty
			if len(%[1]s) != 0 {
				%[2]s
			}
		`, name, body)
	}

	return body
}

func encodeMaxLengthCheck(name string, options *Options) string {
	if options != nil && options.MaxLength > 0 {
		return fmt.Sprintf(`
		// %[1]s maxlen check
		if len(%[1]s) > %d {
			return encoder.ErrMaxLenExceeded
		}
		`, name, options.MaxLength)
	}

	return ""
}

/* Decode */

func WrapDecodeFunc(structName, structPackageName, funcBody string) []byte {
	structType := structName
	if structPackageName != "" {
		structType = fmt.Sprintf("%s.%s", structPackageName, structName)
	}
	return []byte(fmt.Sprintf(`
// Decode%[1]s decodes an object of type %[1]s from the buffer in encoder.Decoder
func Decode%[1]s(d *encoder.Decoder, obj *%[3]s) error {
	%[2]s

	if len(d.Buffer) != 0 {
		return encoder.ErrRemainingBytes
	}

	return nil
}
`, structName, funcBody, structType))
}

func BuildDecodeBool(name string, castType bool, typeName string, options *Options) string {
	assign := "i"
	if castType {
		assign = cast(typeName, assign)
	}
	return fmt.Sprintf(`{
	// %[1]s
	i, err := d.Bool()
	if err != nil {
		return err
	}
	%[1]s = %[2]s
	}
	`, name, assign)
}

func BuildDecodeUint8(name string, castType bool, typeName string, options *Options) string {
	assign := "i"
	if castType {
		assign = cast(typeName, assign)
	}
	return fmt.Sprintf(`{
	// %[1]s
	i, err := d.Uint8()
	if err != nil {
		return err
	}
	%[1]s = %[2]s
	}`, name, assign)
}

func BuildDecodeUint16(name string, castType bool, typeName string, options *Options) string {
	assign := "i"
	if castType {
		assign = cast(typeName, assign)
	}
	return fmt.Sprintf(`{
	// %[1]s
	i, err := d.Uint16()
	if err != nil {
		return err
	}
	%[1]s = %[2]s
	}
	`, name, assign)
}

func BuildDecodeUint32(name string, castType bool, typeName string, options *Options) string {
	assign := "i"
	if castType {
		assign = cast(typeName, assign)
	}
	return fmt.Sprintf(`{
	// %[1]s
	i, err := d.Uint32()
	if err != nil {
		return err
	}
	%[1]s = %[2]s
	}
	`, name, assign)
}

func BuildDecodeUint64(name string, castType bool, typeName string, options *Options) string {
	assign := "i"
	if castType {
		assign = cast(typeName, assign)
	}
	return fmt.Sprintf(`{
	// %[1]s
	i, err := d.Uint64()
	if err != nil {
		return err
	}
	%[1]s = %[2]s
	}
	`, name, assign)
}

func BuildDecodeInt8(name string, castType bool, typeName string, options *Options) string {
	assign := "i"
	if castType {
		assign = cast(typeName, assign)
	}
	return fmt.Sprintf(`{
	// %[1]s
	i, err := d.Int8()
	if err != nil {
		return err
	}
	%[1]s = %[2]s
	}
	`, name, assign)
}

func BuildDecodeInt16(name string, castType bool, typeName string, options *Options) string {
	assign := "i"
	if castType {
		assign = cast(typeName, assign)
	}
	return fmt.Sprintf(`{
	// %[1]s
	i, err := d.Int16()
	if err != nil {
		return err
	}
	%[1]s = %[2]s
	}
	`, name, assign)
}

func BuildDecodeInt32(name string, castType bool, typeName string, options *Options) string {
	assign := "i"
	if castType {
		assign = cast(typeName, assign)
	}
	return fmt.Sprintf(`{
	// %[1]s
	i, err := d.Int32()
	if err != nil {
		return err
	}
	%[1]s = %[2]s
	}
	`, name, assign)
}

func BuildDecodeInt64(name string, castType bool, typeName string, options *Options) string {
	assign := "i"
	if castType {
		assign = cast(typeName, assign)
	}
	return fmt.Sprintf(`{
	// %[1]s
	i, err := d.Int64()
	if err != nil {
		return err
	}
	%[1]s = %[2]s
	}
	`, name, assign)
}

func BuildDecodeByteArray(name string, options *Options) string {
	return fmt.Sprintf(`{
	// %[1]s
	if len(d.Buffer) < len(%[1]s) {
		return encoder.ErrBufferUnderflow
	}
	copy(%[1]s[:], d.Buffer[:len(%[1]s)])
	d.Buffer = d.Buffer[len(%[1]s):]
	}
	`, name)
}

func BuildDecodeArray(name, elemCounterName, elemVarName, elemSection string, options *Options) string {
	return fmt.Sprintf(`{
	// %[1]s
	for %[2]s := range %[1]s {
		%[4]s
	}
	}
	`, name, elemCounterName, elemVarName, elemSection)
}

func BuildDecodeByteSlice(name string, options *Options) string {
	return fmt.Sprintf(`{
	// %[1]s

	%[3]s

	if len(d.Buffer) < 4 {
		return encoder.ErrBufferUnderflow
	}

	ul, err := d.Uint32()
	if err != nil {
		return err
	}

	length := int(ul)
	if length < 0 || length > len(d.Buffer) {
		return encoder.ErrBufferUnderflow
	}

	%[2]s

	%[1]s = make([]byte, length)
	copy(%[1]s[:], d.Buffer[:length])
	d.Buffer = d.Buffer[length:]
	}`, name, decodeMaxLengthCheck(options), decodeOmitEmptyCheck(options))
}

func BuildDecodeSlice(name, elemCounterName, elemVarName, elemSection, typeName string, options *Options) string {
	return fmt.Sprintf(`{
	// %[1]s

	%[7]s

	if len(d.Buffer) < 4 {
		return encoder.ErrBufferUnderflow
	}

	ul, err := d.Uint32()
	if err != nil {
		return err
	}

	length := int(ul)
	if length < 0 || length > len(d.Buffer) {
		return encoder.ErrBufferUnderflow
	}

	%[6]s

	%[1]s = make(%[5]s, length)

	for %[2]s := range %[1]s {
		%[4]s
	}

	}`, name, elemCounterName, elemVarName, elemSection, typeName, decodeMaxLengthCheck(options), decodeOmitEmptyCheck(options))
}

func BuildDecodeString(name string, options *Options) string {
	return fmt.Sprintf(`{
	// %[1]s

	%[3]s

	if len(d.Buffer) < 4 {
		return encoder.ErrBufferUnderflow
	}

	ul, err := d.Uint32()
	if err != nil {
		return err
	}

	length := int(ul)
	if length < 0 || length > len(d.Buffer) {
		return encoder.ErrBufferUnderflow
	}

	%[2]s

	%[1]s = string(d.Buffer[:length])
	d.Buffer = d.Buffer[length:]
	}`, name, decodeMaxLengthCheck(options), decodeOmitEmptyCheck(options))
}

func BuildDecodeMap(name, keyVarName, elemVarName, keySection, elemSection, typeName string, options *Options) string {
	return fmt.Sprintf(`{
	// %[1]s

	%[8]s

	if len(d.Buffer) < 4 {
		return encoder.ErrBufferUnderflow
	}

	ul, err := d.Uint32()
	if err != nil {
		return err
	}

	length := int(ul)
	if length < 0 || length > len(d.Buffer) {
		return encoder.ErrBufferUnderflow
	}

	%[7]s

	%[1]s = make(%[6]s)

	for %[2]s, %[3]s := range %[1]s {
		%[4]s

		%[5]s

		%[1]s[%[2]s] = %[3]s
	}
	}`, name, keyVarName, elemVarName, keySection, elemSection, typeName, decodeMaxLengthCheck(options), decodeOmitEmptyCheck(options))
}

func decodeMaxLengthCheck(options *Options) string {
	if options != nil && options.MaxLength > 0 {
		return fmt.Sprintf(`if length > %d {
			return encoder.ErrMaxLenExceeded
		}`, options.MaxLength)
	}

	return ""
}

func decodeOmitEmptyCheck(options *Options) string {
	if options != nil && options.OmitEmpty {
		return `if len(d.Buffer) == 0 {
			return nil
		}`
	}

	return ""
}
