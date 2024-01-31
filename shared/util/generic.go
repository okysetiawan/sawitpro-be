package util

func In[T comparable](obj T, compared ...T) bool {
	for _, equal := range compared {
		if equal == obj {
			return true
		}
	}
	return false
}
