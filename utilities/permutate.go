package utilities

func Permutate(optionSets [][]int) [][]int {
	totalPermutations := 0
	for _, options := range optionSets {
		if totalPermutations == 0 {
			totalPermutations = len(options)
		} else {
			totalPermutations *= len(options)
		}
	}

	if totalPermutations == 0 {
		return [][]int{}
	}

	permutations := make([][]int, totalPermutations)

	optionSetsLen := len(optionSets)
	indices := make([]int, optionSetsLen)

	nextPermutationIndex := 0
	for {
		newPermutation := make([]int, optionSetsLen)
		for i, index := range indices {
			newPermutation[i] = optionSets[i][index]
		}

		permutations[nextPermutationIndex] = newPermutation
		nextPermutationIndex++

		next := optionSetsLen - 1
		for next >= 0 && (indices[next]+1 >= len(optionSets[next])) {
			next--
		}

		if next < 0 {
			return permutations
		}

		indices[next]++

		for i := next + 1; i < optionSetsLen; i++ {
			indices[i] = 0
		}
	}
}
