package amf

import (
	"encoding/binary"
	"io"
	"math"
	"time"
)

// amf3 polymorphic router
func (d *Decoder) DecodeAmf3(r io.Reader) (interface{}, error) {
	marker, err := ReadMarker(r)
	if err != nil {
		return nil, err
	}

	switch marker {
	case AMF3_UNDEFINED_MARKER:
		return d.DecodeAmf3Undefined(r, false)
	case AMF3_NULL_MARKER:
		return d.DecodeAmf3Null(r, false)
	case AMF3_FALSE_MARKER:
		return d.DecodeAmf3False(r, false)
	case AMF3_TRUE_MARKER:
		return d.DecodeAmf3True(r, false)
	case AMF3_INTEGER_MARKER:
		return d.DecodeAmf3Integer(r, false)
	case AMF3_DOUBLE_MARKER:
		return d.DecodeAmf3Double(r, false)
	case AMF3_STRING_MARKER:
		return d.DecodeAmf3String(r, false)
	case AMF3_XMLDOC_MARKER:
		return nil, Error("decode amf3: unsupported type xmldoc")
	case AMF3_DATE_MARKER:
		return d.DecodeAmf3Date(r, false)
	case AMF3_ARRAY_MARKER:
		return d.DecodeAmf3Array(r, false)
	case AMF3_OBJECT_MARKER:
		return d.DecodeAmf3Object(r, false)
	case AMF3_XML_MARKER:
		return nil, Error("decode amf3: unsupported type xml")
	case AMF3_BYTEARRAY_MARKER:
		return d.DecodeAmf3ByteArray(r, false)
	}

	return nil, Error("decode amf3: unsupported type %d", marker)
}

// marker: 1 byte 0x00
// no additional data
func (d *Decoder) DecodeAmf3Undefined(r io.Reader, decodeMarker bool) (result interface{}, err error) {
	err = AssertMarker(r, decodeMarker, AMF3_UNDEFINED_MARKER)
	return
}

// marker: 1 byte 0x01
// no additional data
func (d *Decoder) DecodeAmf3Null(r io.Reader, decodeMarker bool) (result interface{}, err error) {
	err = AssertMarker(r, decodeMarker, AMF3_NULL_MARKER)
	return
}

// marker: 1 byte 0x02
// no additional data
func (d *Decoder) DecodeAmf3False(r io.Reader, decodeMarker bool) (result bool, err error) {
	err = AssertMarker(r, decodeMarker, AMF3_FALSE_MARKER)
	result = false
	return
}

// marker: 1 byte 0x03
// no additional data
func (d *Decoder) DecodeAmf3True(r io.Reader, decodeMarker bool) (result bool, err error) {
	err = AssertMarker(r, decodeMarker, AMF3_TRUE_MARKER)
	result = true
	return
}

// marker: 1 byte 0x04
func (d *Decoder) DecodeAmf3Integer(r io.Reader, decodeMarker bool) (result uint32, err error) {
	if err = AssertMarker(r, decodeMarker, AMF3_INTEGER_MARKER); err != nil {
		return
	}

	var b byte

	for i := 0; i < 3; i++ {
		b, err = ReadByte(r)
		if err != nil {
			return
		}
		result = (result << 7) + uint32(b&0x7F)
		if (b & 0x80) == 0 {
			return
		}
	}
	b, err = ReadByte(r)
	if err != nil {
		return
	}

	return ((result << 8) + uint32(b)), nil
}

// marker: 1 byte 0x05
func (d *Decoder) DecodeAmf3Double(r io.Reader, decodeMarker bool) (result float64, err error) {
	if err = AssertMarker(r, decodeMarker, AMF3_DOUBLE_MARKER); err != nil {
		return
	}

	err = binary.Read(r, binary.BigEndian, &result)
	if err != nil {
		return float64(0), Error("amf3 decode: unable to read double: %s", err)
	}

	return
}

// marker: 1 byte 0x06
// format:
// - u29 reference int. if reference, no more data. if not reference,
//   length value of bytes to read to complete string.
func (d *Decoder) DecodeAmf3String(r io.Reader, decodeMarker bool) (result string, err error) {
	if err = AssertMarker(r, decodeMarker, AMF3_STRING_MARKER); err != nil {
		return
	}

	var isRef bool
	var refVal uint32
	isRef, refVal, err = d.decodeReferenceInt(r)
	if err != nil {
		return "", Error("amf3 decode: unable to decode string reference and length: %s", err)
	}

	if isRef {
		if refVal > uint32(len(d.stringRefs)) {
			return "", Error("amf3 decode: bad string reference")
		}

		result = d.stringRefs[refVal]
		return
	}

	buf := make([]byte, refVal)
	_, err = r.Read(buf)
	if err != nil {
		return "", Error("amf3 decode: unable to read string: %s", err)
	}

	result = string(buf)
	if result != "" {
		d.stringRefs = append(d.stringRefs, result)
	}

	return
}

// marker: 1 byte 0x08
// format:
// - u29 reference int, if reference, no more data
// - timestamp double
func (d *Decoder) DecodeAmf3Date(r io.Reader, decodeMarker bool) (result time.Time, err error) {
	if err = AssertMarker(r, decodeMarker, AMF3_DATE_MARKER); err != nil {
		return
	}

	var isRef bool
	var refVal uint32
	isRef, refVal, err = d.decodeReferenceInt(r)
	if err != nil {
		return result, Error("amf3 decode: unable to decode date reference and length: %s", err)
	}

	if isRef {
		if refVal > uint32(len(d.objectRefs)) {
			return result, Error("amf3 decode: bad object reference for date")
		}

		log.Debug("is a reference: %v %v", isRef, refVal)
		res, ok := d.objectRefs[refVal].(time.Time)
		if ok != true {
			return result, Error("amf3 decode: unable to extract time from date object references")
		}

		return res, err
	}

	var u64 int64
	err = binary.Read(r, binary.BigEndian, &u64)
	if err != nil {
		return result, Error("amf3 decode: unable to read double: %s", err)
	}

	result = time.Unix(u64, 0)

	d.objectRefs = append(d.objectRefs, result)

	return
}

// marker: 1 byte 0x09
// format:
// - u29 reference int. if reference, no more data.
// - string representing associative array if present
// - n values (length of u29)
func (d *Decoder) DecodeAmf3Array(r io.Reader, decodeMarker bool) (result Array, err error) {
	if err = AssertMarker(r, decodeMarker, AMF3_ARRAY_MARKER); err != nil {
		return
	}

	var isRef bool
	var refVal uint32
	isRef, refVal, err = d.decodeReferenceInt(r)
	if err != nil {
		return result, Error("amf3 decode: unable to decode array reference and length: %s", err)
	}

	if isRef {
		if refVal > uint32(len(d.objectRefs)) {
			return result, Error("amf3 decode: bad object reference for array")
		}

		res, ok := d.objectRefs[refVal].(Array)
		if ok != true {
			return result, Error("amf3 decode: unable to extract array from object references")
		}

		return res, err
	}

	var key string
	key, err = d.DecodeAmf3String(r, false)
	if err != nil {
		return result, Error("amf3 decode: unable to read key for array: %s", err)
	}

	if key != "" {
		return result, Error("amf3 decode: array key is not empty, can't handle associative array")
	}

	for i := uint32(0); i < refVal; i++ {
		tmp, err := d.DecodeAmf3(r)
		if err != nil {
			return result, Error("amf3 decode: array element could not be decoded: %s", err)
		}
		result = append(result, tmp)
	}

	return
}

// marker: 1 byte 0x09
// format: oh dear god
func (d *Decoder) DecodeAmf3Object(r io.Reader, decodeMarker bool) (result TypedObject, err error) {
	if err = AssertMarker(r, decodeMarker, AMF3_OBJECT_MARKER); err != nil {
		return
	}

	// decode the initial u29
	isRef, refVal, err := d.decodeReferenceInt(r)
	if err != nil {
		return result, Error("amf3 decode: unable to decode object reference and length: %s", err)
	}

	// if this is a object reference only, grab it and return it
	if isRef {
		objRefId := refVal
		if objRefId > uint32(len(d.objectRefs)) {
			return result, Error("amf3 decode: bad object reference for array")
		}

		res, ok := d.objectRefs[objRefId].(TypedObject)
		if ok != true {
			return result, Error("amf3 decode: unable to extract typed object from object references")
		}

		return res, err
	}

	// each type has traits that are cached, if the peer sent a reference
	// then we'll need to look it up and use it.
	var trait Trait
	result = *NewTypedObject()
	traitIsRef := (refVal & 0x01) == 0

	if traitIsRef {
		traitRef := refVal >> 1
		if traitRef >= uint32(len(d.traitRefs)) {
			return result, Error("amf3 decode: bad trait reference for object")
		}

		trait = d.traitRefs[traitRef]

	} else {
		// build a new trait from what's left of the given u29
		trait = *NewTrait()
		trait.Externalizable = (refVal & 0x02) != 0
		trait.Dynamic = (refVal & 0x04) != 0

		var cls string
		cls, err = d.DecodeAmf3String(r, false)
		if err != nil {
			return result, Error("amf3 decode: unable to read trait type for object: %s", err)
		}
		result.Type = cls
		trait.Type = cls

		// traits have property keys, encoded as amf3 strings
		propLength := refVal >> 3
		for i := uint32(0); i < propLength; i++ {
			tmp, err := d.DecodeAmf3String(r, false)
			if err != nil {
				return result, Error("amf3 decode: unable to read trait property for object: %s", err)
			}
			trait.Properties = append(trait.Properties, tmp)
		}

		d.traitRefs = append(d.traitRefs, trait)
	}

	d.objectRefs = append(d.objectRefs, result)

	// objects can be externalizable, meaning that the system has no concrete understanding of
	// their properties or how they are encoded. in that case, we need to find and delegate behavior
	// to the right object.
	if trait.Externalizable {
		switch trait.Type {
		case "DSK":
			return d.decodeAmf3ExternalizableDSK(r)
		case "DSA":
			return d.decodeAmf3ExternalizableDSA(r)
		case "flex.messaging.io.ArrayCollection":
			return d.decodeAmf3ExternalizableArrayCollection(r)
		default:
			return d.decodeAmf3ExternalizableOther(r, trait.Type)
		}
	}

	var key string
	var val interface{}

	// non-externalizable objects have property keys in traits, iterate through them
	// and add the read values to the object
	for _, key = range trait.Properties {
		val, err = d.DecodeAmf3(r)
		if err != nil {
			return result, Error("amf3 decode: unable to decode object property: %s", err)
		}

		result.Object[key] = val
	}

	// if an object is dynamic, it can have extra key/value data at the end. in this case,
	// read keys until we get an empty one.
	if trait.Dynamic {
		for {
			key, err = d.DecodeAmf3String(r, false)
			if err != nil {
				return result, Error("amf3 decode: unable to decode dynamic key: %s", err)
			}
			if key == "" {
				break
			}
			val, err = d.DecodeAmf3(r)
			if err != nil {
				return result, Error("amf3 decode: unable to decode dynamic value: %s", err)
			}

			result.Object[key] = val
		}
	}

	return
}

// marker: 1 byte 0x0c
// format:
// - u29 reference int. if reference, no more data. if not reference,
//   length value of bytes to read.
func (d *Decoder) DecodeAmf3ByteArray(r io.Reader, decodeMarker bool) (result []byte, err error) {
	if err = AssertMarker(r, decodeMarker, AMF3_BYTEARRAY_MARKER); err != nil {
		return
	}

	var isRef bool
	var refVal uint32
	isRef, refVal, err = d.decodeReferenceInt(r)
	if err != nil {
		return result, Error("amf3 decode: unable to decode byte array reference and length: %s", err)
	}

	if isRef {
		if refVal > uint32(len(d.objectRefs)) {
			return result, Error("amf3 decode: bad object reference for byte array")
		}

		var ok bool
		result, ok = d.objectRefs[refVal].([]byte)
		if ok != true {
			return result, Error("amf3 decode: unable to convert object ref to bytes")
		}

		return
	}

	result = make([]byte, refVal)
	_, err = r.Read(result)
	if err != nil {
		return result, Error("amf3 decode: unable to read bytearray: %s", err)
	}

	return
}

func (d *Decoder) decodeAmf3ExternalizableDSA(r io.Reader) (result TypedObject, err error) {
	result = *NewTypedObject()
	result.Type = "DSA"

	if err = d.decodeAmf3ExternalFields(r, &result.Object,
		[]string{"body", "clientId", "destination", "headers", "messageId", "timeStamp", "timeToLive"},
		[]string{"clientIdBytes", "messageIdBytes"}); err != nil {
		return result, Error("unable to decode dsa: %s", err)
	}

	if err = d.decodeAmf3ExternalFields(r, &result.Object,
		[]string{"correlationId", "correlationIdBytes"}); err != nil {
		return result, Error("unable to decode dsa: %s", err)
	}

	return
}

func (d *Decoder) decodeAmf3ExternalizableDSK(r io.Reader) (result TypedObject, err error) {
	result = *NewTypedObject()
	result.Type = "DSK"

	if err = d.decodeAmf3ExternalFields(r, &result.Object,
		[]string{"body", "clientId", "destination", "headers", "messageId", "timeStamp", "timeToLive"},
		[]string{"clientIdBytes", "messageIdBytes"}); err != nil {
		return result, Error("unable to decode dsa: %s", err)
	}

	if err = d.decodeAmf3ExternalFields(r, &result.Object,
		[]string{"correlationId", "correlationIdBytes"}); err != nil {
		return result, Error("unable to decode dsa: %s", err)
	}

	var flag uint8
	var flags []uint8
	flags, err = d.decodeAmf3ExternalFlags(r)
	if err != nil {
		return result, Error("unable to decode remaining DSK flags: %s", err)
	}
	for i := uint8(0); i < uint8(len(flags)); i++ {
		flag = flags[i]
		err = d.decodeAmf3ExternalRemains(r, flag, 0)
		if err != nil {
			return result, Error("unable to decode remaining DSK: %s", err)
		}
	}

	return
}

func (d *Decoder) decodeAmf3ExternalizableArrayCollection(r io.Reader) (result TypedObject, err error) {
	var obj interface{}
	obj, err = d.DecodeAmf3(r)
	if err != nil {
		return result, Error("amf3 decode: unable to decode array collection typed object element: %s", err)
	}

	result = *NewTypedObject()
	result.Type = "flex.messaging.io.ArrayCollection"
	result.Object["array"] = obj

	return
}

func (d *Decoder) decodeAmf3ExternalizableOther(r io.Reader, name string) (result TypedObject, err error) {
	n := 0
	buf := make([]byte, 4096)
	n, err = r.Read(buf)
	if err != nil {
		return result, Error("unable to read from buffer")
	}

	DumpBytes(name, buf, n)

	return result, Error("amf3 decode: unable to decode externalizable of type %s", name)
}

func (d *Decoder) decodeAmf3ExternalFlags(r io.Reader) (result []uint8, err error) {
	var flag uint8
	for {
		err = binary.Read(r, binary.BigEndian, &flag)
		if err != nil {
			return result, Error("Unable to read flags")
		}
		result = append(result, flag)
		if flag < 128 {
			break
		}
	}

	return
}

func (d *Decoder) decodeAmf3ExternalFields(r io.Reader, obj *Object, fieldSets ...[]string) (err error) {
	var bit uint8
	var length uint8
	var flag uint8
	var flags []uint8

	var field string
	var fields []string

	flags, err = d.decodeAmf3ExternalFlags(r)
	if err != nil {
		return Error("unable to decode flags for fieldsets %s (%s)", fieldSets, err)
	}

	if len(flags) > len(fieldSets) {
		return Error("too many flags for fieldsets %+v (%d flags, %d sets)", fieldSets, len(flags), len(fieldSets))
	}

	for i := 0; i < len(flags); i++ {
		flag = flags[i]
		fields = fieldSets[i]
		length = uint8(len(fields))

		for j := uint8(0); j < length; j++ {
			field = fields[j]
			bit = uint8(math.Exp2(float64(j)))

			if (flag & bit) != 0 {
				tmp, err := d.DecodeAmf3(r)
				if err != nil {
					return Error("unable to decode external field %s (%d %d %d): %s", field, i, j, bit, err)
				}
				(*obj)[field] = tmp
			}

			err = d.decodeAmf3ExternalRemains(r, flag, length)
			if err != nil {
				return Error("unable to decode external field remains (%d, %d)", flag, length)
			}
		}
	}

	return err
}

func (d *Decoder) decodeAmf3ExternalRemains(r io.Reader, flag uint8, bits uint8) (err error) {
	if (flag >> bits) != 0 {
		for i := bits; i < 6; i++ {
			if ((flag >> i) & 1) != 0 {
				_, err := d.DecodeAmf3(r)
				if err != nil {
					return Error("unable to decode remaining field %d: %s", i, err)
				}
			}
		}
	}

	return nil
}

func (d *Decoder) decodeReferenceInt(r io.Reader) (isRef bool, refVal uint32, err error) {
	u29, err := d.DecodeAmf3Integer(r, false)
	if err != nil {
		return false, 0, Error("amf3 decode: unable to decode reference int: %s", err)
	}

	isRef = u29&0x01 == 0
	refVal = u29 >> 1

	return
}
