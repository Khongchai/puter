package box

// For stuff that requires external data to know the value conversion like currencies.
//
// Static units like measurements don't use this. Those are hard-coded.
type ValueConverter = func(fromValue float64, fromUnit string, toUnit string) (float64, error)
