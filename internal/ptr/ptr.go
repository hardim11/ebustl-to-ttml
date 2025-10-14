package ptr

// helpers to allow you to assign static values more easily
func String(s string) *string { return &s }
func Int(i int) *int          { return &i }
func Bool(b bool) *bool       { return &b }
