package merger

// IDChanges represents the changed ids of two slices of a model type
// after a merge has happened, so dependent objects can be updated
// accordingly. So if the ID of an object of the left slice
// changed from id 5 to 20, it will be represented as: {5: 20}.
type IDChanges struct {
	Left  map[int]int
	Right map[int]int
}
