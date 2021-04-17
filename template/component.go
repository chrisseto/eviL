package template

type Component struct {
}

type Components struct {
	UUIDs          int
	CIDToComponent map[string]string
	// map[componentName]map[ID]ComponentID
	IDToCID map[string]string
}
