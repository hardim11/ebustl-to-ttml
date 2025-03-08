package ttmlgenerate

type TtmlConvertConfiguration struct {
	PreserveSpaces bool
	AddId          bool
	ShuffleTimes   bool
	Debug          bool
}

func TtmlConvertConfigurationDefault() TtmlConvertConfiguration {
	return TtmlConvertConfiguration{
		PreserveSpaces: false,
		AddId:          false,
		ShuffleTimes:   false,
		Debug:          false,
	}
}
