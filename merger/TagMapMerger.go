package merger

import (
	"sort"
	"strconv"
	"strings"

	"github.com/AndreasSko/go-jwlm/model"
)

// MergeTagMaps merges a left and right slice of TagMap. It automatically
// removes redundant entries and also makes sure that the position-order
// stays similar.
func MergeTagMaps(left []*model.TagMap, right []*model.TagMap, conflictSolution map[string]mergeSolution) ([]*model.TagMap, IDChanges, error) {
	tags := make(map[string]map[int]*model.TagMap, len(left)+len(right))

	// Add all entries to tags-map, combine by "PlaylistItemID_LocationID_NoteID"
	for _, side := range [][]*model.TagMap{left, right} {
		for _, tm := range side {
			if tm == nil {
				continue
			}
			var key strings.Builder
			key.Grow(10)
			key.WriteString(strconv.FormatInt(int64(tm.PlaylistItemID.Int32), 10))
			key.WriteString("_")
			key.WriteString(strconv.FormatInt(int64(tm.LocationID.Int32), 10))
			key.WriteString("_")
			key.WriteString(strconv.FormatInt(int64(tm.NoteID.Int32), 10))

			// Add by TagID to sub-map & at the same time filter duplicates
			if _, ok := tags[key.String()]; !ok {
				tags[key.String()] = map[int]*model.TagMap{}
			}

			tags[key.String()][tm.TagID] = tm
		}
	}

	result := make([]*model.TagMap, len(left)+len(right))

	// Go through map in sorted order so we have deterministic results
	sortedTagKeys := make([]string, len(tags))
	i := 0
	for k := range tags {
		sortedTagKeys[i] = k
		i++
	}
	sort.Strings(sortedTagKeys)

	i = 1
	for _, key := range sortedTagKeys {
		tagSet := tags[key]
		sortedTags := make([]*model.TagMap, len(tagSet))

		// Sort by position and update if necessary, so we don't have
		// collision of positions
		j := 0
		for _, tag := range tagSet {
			sortedTags[j] = tag
			j++
		}
		sort.SliceStable(sortedTags, func(i, j int) bool {
			// If equal position, sort by TagID
			if sortedTags[i].Position == sortedTags[j].Position {
				return sortedTags[i].TagID < sortedTags[j].TagID
			}
			return sortedTags[i].Position < sortedTags[j].Position
		})
		for j, tag := range sortedTags {
			tag.Position = j
		}

		// Add to result & update IDs
		for _, tag := range sortedTags {
			result[i] = tag
			result[i].TagMapID = i
			i++
		}
	}

	return result[:i], IDChanges{}, nil
}
