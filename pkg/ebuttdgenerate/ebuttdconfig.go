package ebuttdgenerate

type EbuTtDConvertConfiguration struct {
	PreserveSpaces           bool
	AddId                    bool
	ShuffleTimes             bool
	Debug                    bool
	CellRsolution            string
	TruncateOversizedLines   bool
	IgnoreInvalidDiacritical bool
}

func EbuTtDConvertConfigurationDefault() EbuTtDConvertConfiguration {
	return EbuTtDConvertConfiguration{
		PreserveSpaces:           false,
		AddId:                    false,
		ShuffleTimes:             false,
		Debug:                    false,
		CellRsolution:            "40 24",
		TruncateOversizedLines:   false,
		IgnoreInvalidDiacritical: true,
	}
}
