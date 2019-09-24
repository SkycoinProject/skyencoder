// Code generated by github.com/SkycoinProject/skyencoder. DO NOT EDIT.
package tests

import (
	"bytes"
	"fmt"
	mathrand "math/rand"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/SkycoinProject/skycoin/src/cipher/encoder"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/skycoin/encodertest"
)

func newEmptyOnlyOmitEmptyStructForEncodeTest() *OnlyOmitEmptyStruct {
	var obj OnlyOmitEmptyStruct
	return &obj
}

func newRandomOnlyOmitEmptyStructForEncodeTest(t *testing.T, rand *mathrand.Rand) *OnlyOmitEmptyStruct {
	var obj OnlyOmitEmptyStruct
	err := encodertest.PopulateRandom(&obj, rand, encodertest.PopulateRandomOptions{
		MaxRandLen: 4,
		MinRandLen: 1,
	})
	if err != nil {
		t.Fatalf("encodertest.PopulateRandom failed: %v", err)
	}
	return &obj
}

func newRandomZeroLenOnlyOmitEmptyStructForEncodeTest(t *testing.T, rand *mathrand.Rand) *OnlyOmitEmptyStruct {
	var obj OnlyOmitEmptyStruct
	err := encodertest.PopulateRandom(&obj, rand, encodertest.PopulateRandomOptions{
		MaxRandLen:    0,
		MinRandLen:    0,
		EmptySliceNil: false,
		EmptyMapNil:   false,
	})
	if err != nil {
		t.Fatalf("encodertest.PopulateRandom failed: %v", err)
	}
	return &obj
}

func newRandomZeroLenNilOnlyOmitEmptyStructForEncodeTest(t *testing.T, rand *mathrand.Rand) *OnlyOmitEmptyStruct {
	var obj OnlyOmitEmptyStruct
	err := encodertest.PopulateRandom(&obj, rand, encodertest.PopulateRandomOptions{
		MaxRandLen:    0,
		MinRandLen:    0,
		EmptySliceNil: true,
		EmptyMapNil:   true,
	})
	if err != nil {
		t.Fatalf("encodertest.PopulateRandom failed: %v", err)
	}
	return &obj
}

func testSkyencoderOnlyOmitEmptyStruct(t *testing.T, obj *OnlyOmitEmptyStruct) {
	isEncodableField := func(f reflect.StructField) bool {
		// Skip unexported fields
		if f.PkgPath != "" {
			return false
		}

		// Skip fields disabled with and enc:"- struct tag
		tag := f.Tag.Get("enc")
		return !strings.HasPrefix(tag, "-,") && tag != "-"
	}

	hasOmitEmptyField := func(obj interface{}) bool {
		v := reflect.ValueOf(obj)
		switch v.Kind() {
		case reflect.Ptr:
			v = v.Elem()
		}

		switch v.Kind() {
		case reflect.Struct:
			t := v.Type()
			n := v.NumField()
			f := t.Field(n - 1)
			tag := f.Tag.Get("enc")
			return isEncodableField(f) && strings.Contains(tag, ",omitempty")
		default:
			return false
		}
	}

	// returns the number of bytes encoded by an omitempty field on a given object
	omitEmptyLen := func(obj interface{}) uint64 {
		if !hasOmitEmptyField(obj) {
			return 0
		}

		v := reflect.ValueOf(obj)
		switch v.Kind() {
		case reflect.Ptr:
			v = v.Elem()
		}

		switch v.Kind() {
		case reflect.Struct:
			n := v.NumField()
			f := v.Field(n - 1)
			if f.Len() == 0 {
				return 0
			}
			return uint64(4 + f.Len())

		default:
			return 0
		}
	}

	// EncodeSize

	n1 := encoder.Size(obj)
	n2 := EncodeSizeOnlyOmitEmptyStruct(obj)

	if uint64(n1) != n2 {
		t.Fatalf("encoder.Size() != EncodeSizeOnlyOmitEmptyStruct() (%d != %d)", n1, n2)
	}

	// Encode

	// encoder.Serialize
	data1 := encoder.Serialize(obj)

	// Encode
	data2, err := EncodeOnlyOmitEmptyStruct(obj)
	if err != nil {
		t.Fatalf("EncodeOnlyOmitEmptyStruct failed: %v", err)
	}
	if uint64(len(data2)) != n2 {
		t.Fatal("EncodeOnlyOmitEmptyStruct produced bytes of unexpected length")
	}
	if len(data1) != len(data2) {
		t.Fatalf("len(encoder.Serialize()) != len(EncodeOnlyOmitEmptyStruct()) (%d != %d)", len(data1), len(data2))
	}

	// EncodeToBuffer
	data3 := make([]byte, n2+5)
	if err := EncodeOnlyOmitEmptyStructToBuffer(data3, obj); err != nil {
		t.Fatalf("EncodeOnlyOmitEmptyStructToBuffer failed: %v", err)
	}

	if !bytes.Equal(data1, data2) {
		t.Fatal("encoder.Serialize() != Encode[1]s()")
	}

	// Decode

	// encoder.DeserializeRaw
	var obj2 OnlyOmitEmptyStruct
	if n, err := encoder.DeserializeRaw(data1, &obj2); err != nil {
		t.Fatalf("encoder.DeserializeRaw failed: %v", err)
	} else if n != uint64(len(data1)) {
		t.Fatalf("encoder.DeserializeRaw failed: %v", encoder.ErrRemainingBytes)
	}
	if !cmp.Equal(*obj, obj2, cmpopts.EquateEmpty(), encodertest.IgnoreAllUnexported()) {
		t.Fatal("encoder.DeserializeRaw result wrong")
	}

	// Decode
	var obj3 OnlyOmitEmptyStruct
	if n, err := DecodeOnlyOmitEmptyStruct(data2, &obj3); err != nil {
		t.Fatalf("DecodeOnlyOmitEmptyStruct failed: %v", err)
	} else if n != uint64(len(data2)) {
		t.Fatalf("DecodeOnlyOmitEmptyStruct bytes read length should be %d, is %d", len(data2), n)
	}
	if !cmp.Equal(obj2, obj3, cmpopts.EquateEmpty(), encodertest.IgnoreAllUnexported()) {
		t.Fatal("encoder.DeserializeRaw() != DecodeOnlyOmitEmptyStruct()")
	}

	// Decode, excess buffer
	var obj4 OnlyOmitEmptyStruct
	n, err := DecodeOnlyOmitEmptyStruct(data3, &obj4)
	if err != nil {
		t.Fatalf("DecodeOnlyOmitEmptyStruct failed: %v", err)
	}

	if hasOmitEmptyField(&obj4) && omitEmptyLen(&obj4) == 0 {
		// 4 bytes read for the omitEmpty length, which should be zero (see the 5 bytes added above)
		if n != n2+4 {
			t.Fatalf("DecodeOnlyOmitEmptyStruct bytes read length should be %d, is %d", n2+4, n)
		}
	} else {
		if n != n2 {
			t.Fatalf("DecodeOnlyOmitEmptyStruct bytes read length should be %d, is %d", n2, n)
		}
	}
	if !cmp.Equal(obj2, obj4, cmpopts.EquateEmpty(), encodertest.IgnoreAllUnexported()) {
		t.Fatal("encoder.DeserializeRaw() != DecodeOnlyOmitEmptyStruct()")
	}

	// DecodeExact
	var obj5 OnlyOmitEmptyStruct
	if err := DecodeOnlyOmitEmptyStructExact(data2, &obj5); err != nil {
		t.Fatalf("DecodeOnlyOmitEmptyStruct failed: %v", err)
	}
	if !cmp.Equal(obj2, obj5, cmpopts.EquateEmpty(), encodertest.IgnoreAllUnexported()) {
		t.Fatal("encoder.DeserializeRaw() != DecodeOnlyOmitEmptyStruct()")
	}

	// Check that the bytes read value is correct when providing an extended buffer
	if !hasOmitEmptyField(&obj3) || omitEmptyLen(&obj3) > 0 {
		padding := []byte{0xFF, 0xFE, 0xFD, 0xFC}
		data4 := append(data2[:], padding...)
		if n, err := DecodeOnlyOmitEmptyStruct(data4, &obj3); err != nil {
			t.Fatalf("DecodeOnlyOmitEmptyStruct failed: %v", err)
		} else if n != uint64(len(data2)) {
			t.Fatalf("DecodeOnlyOmitEmptyStruct bytes read length should be %d, is %d", len(data2), n)
		}
	}
}

func TestSkyencoderOnlyOmitEmptyStruct(t *testing.T) {
	rand := mathrand.New(mathrand.NewSource(time.Now().Unix()))

	type testCase struct {
		name string
		obj  *OnlyOmitEmptyStruct
	}

	cases := []testCase{
		{
			name: "empty object",
			obj:  newEmptyOnlyOmitEmptyStructForEncodeTest(),
		},
	}

	nRandom := 10

	for i := 0; i < nRandom; i++ {
		cases = append(cases, testCase{
			name: fmt.Sprintf("randomly populated object %d", i),
			obj:  newRandomOnlyOmitEmptyStructForEncodeTest(t, rand),
		})
		cases = append(cases, testCase{
			name: fmt.Sprintf("randomly populated object %d with zero length variable length contents", i),
			obj:  newRandomZeroLenOnlyOmitEmptyStructForEncodeTest(t, rand),
		})
		cases = append(cases, testCase{
			name: fmt.Sprintf("randomly populated object %d with zero length variable length contents set to nil", i),
			obj:  newRandomZeroLenNilOnlyOmitEmptyStructForEncodeTest(t, rand),
		})
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			testSkyencoderOnlyOmitEmptyStruct(t, tc.obj)
		})
	}
}

func decodeOnlyOmitEmptyStructExpectError(t *testing.T, buf []byte, expectedErr error) {
	var obj OnlyOmitEmptyStruct
	if _, err := DecodeOnlyOmitEmptyStruct(buf, &obj); err == nil {
		t.Fatal("DecodeOnlyOmitEmptyStruct: expected error, got nil")
	} else if err != expectedErr {
		t.Fatalf("DecodeOnlyOmitEmptyStruct: expected error %q, got %q", expectedErr, err)
	}
}

func decodeOnlyOmitEmptyStructExactExpectError(t *testing.T, buf []byte, expectedErr error) {
	var obj OnlyOmitEmptyStruct
	if err := DecodeOnlyOmitEmptyStructExact(buf, &obj); err == nil {
		t.Fatal("DecodeOnlyOmitEmptyStructExact: expected error, got nil")
	} else if err != expectedErr {
		t.Fatalf("DecodeOnlyOmitEmptyStructExact: expected error %q, got %q", expectedErr, err)
	}
}

func testSkyencoderOnlyOmitEmptyStructDecodeErrors(t *testing.T, k int, tag string, obj *OnlyOmitEmptyStruct) {
	isEncodableField := func(f reflect.StructField) bool {
		// Skip unexported fields
		if f.PkgPath != "" {
			return false
		}

		// Skip fields disabled with and enc:"- struct tag
		tag := f.Tag.Get("enc")
		return !strings.HasPrefix(tag, "-,") && tag != "-"
	}

	numEncodableFields := func(obj interface{}) int {
		v := reflect.ValueOf(obj)
		switch v.Kind() {
		case reflect.Ptr:
			v = v.Elem()
		}

		switch v.Kind() {
		case reflect.Struct:
			t := v.Type()

			n := 0
			for i := 0; i < v.NumField(); i++ {
				f := t.Field(i)
				if !isEncodableField(f) {
					continue
				}
				n++
			}
			return n
		default:
			return 0
		}
	}

	hasOmitEmptyField := func(obj interface{}) bool {
		v := reflect.ValueOf(obj)
		switch v.Kind() {
		case reflect.Ptr:
			v = v.Elem()
		}

		switch v.Kind() {
		case reflect.Struct:
			t := v.Type()
			n := v.NumField()
			f := t.Field(n - 1)
			tag := f.Tag.Get("enc")
			return isEncodableField(f) && strings.Contains(tag, ",omitempty")
		default:
			return false
		}
	}

	// returns the number of bytes encoded by an omitempty field on a given object
	omitEmptyLen := func(obj interface{}) uint64 {
		if !hasOmitEmptyField(obj) {
			return 0
		}

		v := reflect.ValueOf(obj)
		switch v.Kind() {
		case reflect.Ptr:
			v = v.Elem()
		}

		switch v.Kind() {
		case reflect.Struct:
			n := v.NumField()
			f := v.Field(n - 1)
			if f.Len() == 0 {
				return 0
			}
			return uint64(4 + f.Len())

		default:
			return 0
		}
	}

	n := EncodeSizeOnlyOmitEmptyStruct(obj)
	buf, err := EncodeOnlyOmitEmptyStruct(obj)
	if err != nil {
		t.Fatalf("EncodeOnlyOmitEmptyStruct failed: %v", err)
	}

	// A nil buffer cannot decode, unless the object is a struct with a single omitempty field
	if hasOmitEmptyField(obj) && numEncodableFields(obj) > 1 {
		t.Run(fmt.Sprintf("%d %s buffer underflow nil", k, tag), func(t *testing.T) {
			decodeOnlyOmitEmptyStructExpectError(t, nil, encoder.ErrBufferUnderflow)
		})

		t.Run(fmt.Sprintf("%d %s exact buffer underflow nil", k, tag), func(t *testing.T) {
			decodeOnlyOmitEmptyStructExactExpectError(t, nil, encoder.ErrBufferUnderflow)
		})
	}

	// Test all possible truncations of the encoded byte array, but skip
	// a truncation that would be valid where omitempty is removed
	skipN := n - omitEmptyLen(obj)
	for i := uint64(0); i < n; i++ {
		if i == skipN {
			continue
		}

		t.Run(fmt.Sprintf("%d %s buffer underflow bytes=%d", k, tag, i), func(t *testing.T) {
			decodeOnlyOmitEmptyStructExpectError(t, buf[:i], encoder.ErrBufferUnderflow)
		})

		t.Run(fmt.Sprintf("%d %s exact buffer underflow bytes=%d", k, tag, i), func(t *testing.T) {
			decodeOnlyOmitEmptyStructExactExpectError(t, buf[:i], encoder.ErrBufferUnderflow)
		})
	}

	// Append 5 bytes for omit empty with a 0 length prefix, to cause an ErrRemainingBytes.
	// If only 1 byte is appended, the decoder will try to read the 4-byte length prefix,
	// and return an ErrBufferUnderflow instead
	if hasOmitEmptyField(obj) {
		buf = append(buf, []byte{0, 0, 0, 0, 0}...)
	} else {
		buf = append(buf, 0)
	}

	t.Run(fmt.Sprintf("%d %s exact buffer remaining bytes", k, tag), func(t *testing.T) {
		decodeOnlyOmitEmptyStructExactExpectError(t, buf, encoder.ErrRemainingBytes)
	})
}

func TestSkyencoderOnlyOmitEmptyStructDecodeErrors(t *testing.T) {
	rand := mathrand.New(mathrand.NewSource(time.Now().Unix()))
	n := 10

	for i := 0; i < n; i++ {
		emptyObj := newEmptyOnlyOmitEmptyStructForEncodeTest()
		fullObj := newRandomOnlyOmitEmptyStructForEncodeTest(t, rand)
		testSkyencoderOnlyOmitEmptyStructDecodeErrors(t, i, "empty", emptyObj)
		testSkyencoderOnlyOmitEmptyStructDecodeErrors(t, i, "full", fullObj)
	}
}
