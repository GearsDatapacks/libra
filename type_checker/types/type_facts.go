package types

func IsInt(ty Type) bool {
	num, ok := Unwrap(ty).(Numeric)
	if !ok {
		return false
	}
	return num.Kind == NumInt || num.Kind == NumUint
}

func IsFloat(ty Type) bool {
	num, ok := Unwrap(ty).(Numeric)
	if !ok {
		return false
	}
	return num.Kind == NumFloat
}

func IsStruct(ty Type) bool {
	_, ok := Unwrap(ty).(*Struct)
	return ok
}

func IsPtr(ty Type) bool {
	_, ok := Unwrap(ty).(*Pointer)
	return ok
}

func IsBool(ty Type) bool {
	pt, ok := Unwrap(ty).(PrimaryType)
	return ok && pt == Bool
}

func IsString(ty Type) bool {
	pt, ok := Unwrap(ty).(PrimaryType)
	return ok && pt == String
}
