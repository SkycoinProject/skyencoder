// Code generated by github.com/skycoin/skyencoder. DO NOT EDIT.
package tests

import (
	"bytes"
	"fmt"
	mathrand "math/rand"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/skycoin/skycoin/src/cipher/encoder"
	"github.com/skycoin/skycoin/src/cipher/encoder/encodertest"
)

func newEmptyOmitEmptyStructForEncodeTest() *OmitEmptyStruct {
	var obj OmitEmptyStruct
	return &obj
}

func newRandomOmitEmptyStructForEncodeTest(t *testing.T, rand *mathrand.Rand) *OmitEmptyStruct {
	var obj OmitEmptyStruct
	err := encodertest.PopulateRandom(&obj, rand, encodertest.PopulateRandomOptions{
		MaxRandLen: 4,
		MinRandLen: 1,
	})
	if err != nil {
		t.Fatalf("encodertest.PopulateRandom failed: %v", err)
	}
	return &obj
}

func newRandomZeroLenOmitEmptyStructForEncodeTest(t *testing.T, rand *mathrand.Rand) *OmitEmptyStruct {
	var obj OmitEmptyStruct
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

func newRandomZeroLenNilOmitEmptyStructForEncodeTest(t *testing.T, rand *mathrand.Rand) *OmitEmptyStruct {
	var obj OmitEmptyStruct
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

func testSkyencoderOmitEmptyStruct(t *testing.T, obj *OmitEmptyStruct) {
	// EncodeSize

	n1 := encoder.Size(obj)
	n2 := EncodeSizeOmitEmptyStruct(obj)

	if n1 != n2 {
		t.Fatalf("encoder.Size() != EncodeSizeOmitEmptyStruct() (%d != %d)", n1, n2)
	}

	// Encode

	data1 := encoder.Serialize(obj)

	data2 := make([]byte, n2)
	err := EncodeOmitEmptyStruct(data2, obj)
	if err != nil {
		t.Fatalf("EncodeOmitEmptyStruct failed: %v", err)
	}

	if len(data1) != len(data2) {
		t.Fatalf("len(encoder.Serialize()) != len(EncodeOmitEmptyStruct()) (%d != %d)", len(data1), len(data2))
	}

	if !bytes.Equal(data1, data2) {
		t.Fatal("encoder.Serialize() != EncodeOmitEmptyStruct()")
	}

	// Decode

	var obj2 OmitEmptyStruct
	err = encoder.DeserializeRaw(data1, &obj2)
	if err != nil {
		t.Fatalf("encoder.DeserializeRaw failed: %v", err)
	}

	if !cmp.Equal(*obj, obj2, cmpopts.EquateEmpty(), encodertest.IgnoreAllUnexported()) {
		t.Fatal("encoder.DeserializeRaw result wrong")
	}

	var obj3 OmitEmptyStruct
	n, err := DecodeOmitEmptyStruct(data2, &obj3)
	if err != nil {
		t.Fatalf("DecodeOmitEmptyStruct failed: %v", err)
	}
	if n != len(data2) {
		t.Fatalf("DecodeOmitEmptyStruct bytes read length should be %d, is %d", len(data2), n)
	}

	if !cmp.Equal(obj2, obj3, cmpopts.EquateEmpty(), encodertest.IgnoreAllUnexported()) {
		t.Fatal("encoder.DeserializeRaw() != DecodeOmitEmptyStruct()")
	}

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
	omitEmptyLen := func(obj interface{}) int {
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
			return 4 + f.Len()

		default:
			return 0
		}
	}

	// Check that the bytes read value is correct when providing an extended buffer
	if !hasOmitEmptyField(&obj3) || omitEmptyLen(&obj3) > 0 {
		padding := []byte{0xFF, 0xFE, 0xFD, 0xFC}
		data3 := append(data2[:], padding...)
		n, err = DecodeOmitEmptyStruct(data3, &obj3)
		if err != nil {
			t.Fatalf("DecodeOmitEmptyStruct failed: %v", err)
		}
		if n != len(data2) {
			t.Fatalf("DecodeOmitEmptyStruct bytes read length should be %d, is %d", len(data2), n)
		}
	}
}

func TestSkyencoderOmitEmptyStruct(t *testing.T) {
	rand := mathrand.New(mathrand.NewSource(time.Now().Unix()))

	type testCase struct {
		name string
		obj  *OmitEmptyStruct
	}

	cases := []testCase{
		{
			name: "empty object",
			obj:  newEmptyOmitEmptyStructForEncodeTest(),
		},
	}

	nRandom := 10

	for i := 0; i < nRandom; i++ {
		cases = append(cases, testCase{
			name: fmt.Sprintf("randomly populated object %d", i),
			obj:  newRandomOmitEmptyStructForEncodeTest(t, rand),
		})
		cases = append(cases, testCase{
			name: fmt.Sprintf("randomly populated object %d with zero length variable length contents", i),
			obj:  newRandomZeroLenOmitEmptyStructForEncodeTest(t, rand),
		})
		cases = append(cases, testCase{
			name: fmt.Sprintf("randomly populated object %d with zero length variable length contents set to nil", i),
			obj:  newRandomZeroLenNilOmitEmptyStructForEncodeTest(t, rand),
		})
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			testSkyencoderOmitEmptyStruct(t, tc.obj)
		})
	}
}

func decodeOmitEmptyStructExpectError(t *testing.T, buf []byte, expectedErr error) {
	var obj OmitEmptyStruct
	_, err := DecodeOmitEmptyStruct(buf, &obj)

	if err == nil {
		t.Fatal("DecodeOmitEmptyStruct: expected error, got nil")
	}

	if err != expectedErr {
		t.Fatalf("DecodeOmitEmptyStruct: expected error %q, got %q", expectedErr, err)
	}
}

func testSkyencoderOmitEmptyStructDecodeErrors(t *testing.T, k int, tag string, obj *OmitEmptyStruct) {
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
	omitEmptyLen := func(obj interface{}) int {
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
			return 4 + f.Len()

		default:
			return 0
		}
	}

	n := EncodeSizeOmitEmptyStruct(obj)
	buf := make([]byte, n)
	err := EncodeOmitEmptyStruct(buf, obj)
	if err != nil {
		t.Fatalf("EncodeOmitEmptyStruct failed: %v", err)
	}

	// A nil buffer cannot decode, unless the object is a struct with a single omitempty field
	if hasOmitEmptyField(obj) && numEncodableFields(obj) > 1 {
		t.Run(fmt.Sprintf("%d %s buffer underflow nil", k, tag), func(t *testing.T) {
			decodeOmitEmptyStructExpectError(t, nil, encoder.ErrBufferUnderflow)
		})
	}

	// Test all possible truncations of the encoded byte array, but skip
	// a truncation that would be valid where omitempty is removed
	skipN := n - omitEmptyLen(obj)
	for i := 0; i < n; i++ {
		if i == skipN {
			continue
		}
		t.Run(fmt.Sprintf("%d %s buffer underflow bytes=%d", k, tag, i), func(t *testing.T) {
			decodeOmitEmptyStructExpectError(t, buf[:i], encoder.ErrBufferUnderflow)
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
}

func TestSkyencoderOmitEmptyStructDecodeErrors(t *testing.T) {
	rand := mathrand.New(mathrand.NewSource(time.Now().Unix()))
	n := 10

	for i := 0; i < n; i++ {
		emptyObj := newEmptyOmitEmptyStructForEncodeTest()
		fullObj := newRandomOmitEmptyStructForEncodeTest(t, rand)
		testSkyencoderOmitEmptyStructDecodeErrors(t, i, "empty", emptyObj)
		testSkyencoderOmitEmptyStructDecodeErrors(t, i, "full", fullObj)
	}
}
