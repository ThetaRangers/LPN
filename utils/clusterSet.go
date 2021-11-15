package utils

type ClusterSet struct {
	cluster map[string]struct{}
}

func NewClusterSet() ClusterSet {
	return ClusterSet{
		cluster: make(map[string]struct{}),
	}
}

func (s ClusterSet) Contains(key string) bool {
	_, exists := s.cluster[key]
	return exists
}

func (s ClusterSet) Add(key string) {
	s.cluster[key] = struct{}{}
}

func (s ClusterSet) Remove(key string) bool {
	_, exists := s.cluster[key]
	if !exists {
		return false
	}
	delete(s.cluster, key)
	return true
}

func (s ClusterSet) Len() int {
	return len(s.cluster)
}

func (s ClusterSet) GetAll() []string {
	clusterSlice := make([]string, s.Len())
	i := 0
	for k := range s.cluster {
		clusterSlice[i] = k
		i++
	}
	return clusterSlice
}

func (s ClusterSet) Invalidate() {
	for k:= range s.cluster {
		delete(s.cluster, k)
	}
}
