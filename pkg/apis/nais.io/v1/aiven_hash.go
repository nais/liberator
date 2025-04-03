package nais_io_v1

func ApplyAivenGeneration(in AivenInterface, generation uint64, hash uint64) uint64 {
	if usesAiven(in) {
		return hash + generation%10
	}
	return hash
}

func usesAiven(in AivenInterface) bool {
	return in.GetKafka() != nil || in.GetInflux() != nil || in.GetOpenSearch() != nil
}
