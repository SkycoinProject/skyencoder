// Code generated by github.com/skycoin/skyencoder. DO NOT EDIT.
package tests

import (
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

func newEmptyMaxLenNestedMapValueStruct2ForEncodeTest() *MaxLenNestedMapValueStruct2 {
	var obj MaxLenNestedMapValueStruct2
	return &obj
}

func newRandomMaxLenNestedMapValueStruct2ForEncodeTest(t *testing.T, rand *mathrand.Rand) *MaxLenNestedMapValueStruct2 {
	var obj MaxLenNestedMapValueStruct2
	err := encodertest.PopulateRandom(&obj, rand, encodertest.PopulateRandomOptions{
		MaxRandLen: 4,
		MinRandLen: 1,
	})
	if err != nil {
		t.Fatalf("encodertest.PopulateRandom failed: %v", err)
	}
	return &obj
}

func newRandomZeroLenMaxLenNestedMapValueStruct2ForEncodeTest(t *testing.T, rand *mathrand.Rand) *MaxLenNestedMapValueStruct2 {
	var obj MaxLenNestedMapValueStruct2
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

func newRandomZeroLenNilMaxLenNestedMapValueStruct2ForEncodeTest(t *testing.T, rand *mathrand.Rand) *MaxLenNestedMapValueStruct2 {
	var obj MaxLenNestedMapValueStruct2
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

func testSkyencoderMaxLenNestedMapValueStruct2(t *testing.T, obj *MaxLenNestedMapValueStruct2) {
	// EncodeSize

	n1 := encoder.Size(obj)
	n2 := EncodeSizeMaxLenNestedMapValueStruct2(obj)

	if n1 != n2 {
		t.Fatalf("encoder.Size() != EncodeSizeMaxLenNestedMapValueStruct2() (%d != %d)", n1, n2)
	}

	// Encode

	data1 := encoder.Serialize(obj)

	data2 := make([]byte, n2)
	err := EncodeMaxLenNestedMapValueStruct2(data2, obj)
	if err != nil {
		t.Fatalf("EncodeMaxLenNestedMapValueStruct2 failed: %v", err)
	}

	if len(data1) != len(data2) {
		t.Fatalf("len(encoder.Serialize()) != len(EncodeMaxLenNestedMapValueStruct2()) (%d != %d)", len(data1), len(data2))
	}

	// Decode

	var obj2 MaxLenNestedMapValueStruct2
	err = encoder.DeserializeRaw(data1, &obj2)
	if err != nil {
		t.Fatalf("encoder.DeserializeRaw failed: %v", err)
	}

	if !cmp.Equal(*obj, obj2, cmpopts.EquateEmpty(), encodertest.IgnoreAllUnexported()) {
		t.Fatal("encoder.DeserializeRaw result wrong")
	}

	var obj3 MaxLenNestedMapValueStruct2
	err = DecodeMaxLenNestedMapValueStruct2(data2, &obj3)
	if err != nil {
		t.Fatalf("DecodeMaxLenNestedMapValueStruct2 failed: %v", err)
	}

	if !cmp.Equal(obj2, obj3, cmpopts.EquateEmpty(), encodertest.IgnoreAllUnexported()) {
		t.Fatal("encoder.DeserializeRaw() != DecodeMaxLenNestedMapValueStruct2()")
	}
}

func TestSkyencoderMaxLenNestedMapValueStruct2(t *testing.T) {
	rand := mathrand.New(mathrand.NewSource(time.Now().Unix()))

	type testCase struct {
		name string
		obj  *MaxLenNestedMapValueStruct2
	}

	cases := []testCase{
		{
			name: "empty object",
			obj:  newEmptyMaxLenNestedMapValueStruct2ForEncodeTest(),
		},
	}

	nRandom := 10

	for i := 0; i < nRandom; i++ {
		cases = append(cases, testCase{
			name: fmt.Sprintf("randomly populated object %d", i),
			obj:  newRandomMaxLenNestedMapValueStruct2ForEncodeTest(t, rand),
		})
		cases = append(cases, testCase{
			name: fmt.Sprintf("randomly populated object %d with zero length variable length contents", i),
			obj:  newRandomZeroLenMaxLenNestedMapValueStruct2ForEncodeTest(t, rand),
		})
		cases = append(cases, testCase{
			name: fmt.Sprintf("randomly populated object %d with zero length variable length contents set to nil", i),
			obj:  newRandomZeroLenNilMaxLenNestedMapValueStruct2ForEncodeTest(t, rand),
		})
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			testSkyencoderMaxLenNestedMapValueStruct2(t, tc.obj)
		})
	}
}

func decodeMaxLenNestedMapValueStruct2ExpectError(t *testing.T, buf []byte, expectedErr error) {
	var obj MaxLenNestedMapValueStruct2
	err := DecodeMaxLenNestedMapValueStruct2(buf, &obj)

	if err == nil {
		t.Fatal("DecodeMaxLenNestedMapValueStruct2: expected error, got nil")
	}

	if err != expectedErr {
		t.Fatalf("DecodeMaxLenNestedMapValueStruct2: expected error %q, got %q", expectedErr, err)
	}
}

func testSkyencoderMaxLenNestedMapValueStruct2DecodeErrors(t *testing.T, k int, obj *MaxLenNestedMapValueStruct2) {
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

	n := EncodeSizeMaxLenNestedMapValueStruct2(obj)
	buf := make([]byte, n)
	err := EncodeMaxLenNestedMapValueStruct2(buf, obj)
	if err != nil {
		t.Fatalf("EncodeMaxLenNestedMapValueStruct2 failed: %v", err)
	}

	// A nil buffer cannot decode, unless the object is a struct with a single omitempty field
	if hasOmitEmptyField(obj) && numEncodableFields(obj) > 1 {
		t.Run(fmt.Sprintf("%d buffer underflow nil", k), func(t *testing.T) {
			decodeMaxLenNestedMapValueStruct2ExpectError(t, nil, encoder.ErrBufferUnderflow)
		})
	}

	// Test all possible truncations of the encoded byte array, but skip
	// a truncation that would be valid where omitempty is removed
	skipN := n - omitEmptyLen(obj)
	for i := 0; i < n; i++ {
		if i == skipN {
			continue
		}
		t.Run(fmt.Sprintf("%d buffer underflow bytes=%d", k, i), func(t *testing.T) {
			decodeMaxLenNestedMapValueStruct2ExpectError(t, buf[:i], encoder.ErrBufferUnderflow)
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

	// Buffer too long
	buf = append(buf, 0)
	t.Run(fmt.Sprintf("%d remaining bytes", k), func(t *testing.T) {
		decodeMaxLenNestedMapValueStruct2ExpectError(t, buf[:], encoder.ErrRemainingBytes)
	})
}

func TestSkyencoderMaxLenNestedMapValueStruct2DecodeErrors(t *testing.T) {
	rand := mathrand.New(mathrand.NewSource(time.Now().Unix()))
	n := 10

	for i := 0; i < n; i++ {
		emptyObj := newEmptyMaxLenNestedMapValueStruct2ForEncodeTest()
		fullObj := newRandomMaxLenNestedMapValueStruct2ForEncodeTest(t, rand)
		testSkyencoderMaxLenNestedMapValueStruct2DecodeErrors(t, i, emptyObj)
		testSkyencoderMaxLenNestedMapValueStruct2DecodeErrors(t, i, fullObj)
	}
}
