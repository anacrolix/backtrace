package backtrace

// I don't think using custom flags or verbs to optionally extract backtraces works because there's
// only a handful of blessed flags and they're all accounted for, and fmt calls Error() on the first
// error seen. We could have a custom Format if you are directly printing an errorWith, but I don't
// see why that would happen.

//func (me errorWith) Format(f fmt.State, verb rune) {
//	fmt.Fprintf(f, fmt.FormatString(f, verb)+"\n", me.wrapped)
//	if !f.Flag('+') {
//		panic("didn't get +")
//		return
//	}
//	// Maybe a bytes.Buffer is better here, or writing straight through?
//	var sb strings.Builder
//	me.toBacktrace().build(&sb)
//	f.Write([]byte(sb.String()))
//}

//var _ fmt.Formatter = errorWith{}
