package botStructures

// VariableType
/*
 * Abstraction around environment variables to run checks.
 * Name is the title of the variable entry.
 * RawValue is the value as interpreted by the OS, in string form.
 * IsNumeric is whether the value is actually a numeric type when transformed.
 * SignedInteger is whether the value is signed or unsigned (i64 vs u64)
 */
type VariableType struct {
	Name          string
	RawValue      string
	IsNumeric     bool
	SignedInteger bool
}
