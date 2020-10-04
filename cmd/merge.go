package cmd

import (
	"fmt"
	"os"
	"reflect"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/AndreasSko/go-jwlm/merger"
	"github.com/AndreasSko/go-jwlm/model"
	"github.com/buger/goterm"
	"github.com/jedib0t/go-pretty/table"
	"github.com/spf13/cobra"

	log "github.com/sirupsen/logrus"
)

// mergeCmd represents the merge command
var mergeCmd = &cobra.Command{
	Use:   "merge <left-backup> <right-backup> <dest-filename>",
	Short: "Merge two JW Library backup files",
	Long: `merge imports the left and right .jwlibrary backup file, merges them and 
exports it to the destination file. If a collision between the left and 
the right backup is detected, the user is asked to choose which side should
be included in the merged backup. You are able to let the merger 
automatically solve conflicts using the 'chooseLeft', 'chooseRight', and 
'chooseNewest' resolvers (see Flags).`,
	Example: `go-jwlm merge left.jwlibrary right.jwlibrary merged.jwlibrary
go-jwlm merge left.jwlibrary right.jwlibrary merged.jwlibrary --bookmarks chooseLeft --markings chooseRight --notes chooseNewest`,
	Run: func(cmd *cobra.Command, args []string) {
		leftFilename := args[0]
		rightFilename := args[1]
		mergedFilename := args[2]
		merge(leftFilename, rightFilename, mergedFilename, terminal.Stdio{In: os.Stdin, Out: os.Stdout, Err: os.Stderr})
	},
	Args: cobra.ExactArgs(3),
}

// BookmarkResolver represents a resolver that should be used for conflicting Bookmarks
var BookmarkResolver string

// MarkingResolver represents a resolver that should be used for conflicting UserMarkBlockRanges
var MarkingResolver string

// NoteResolver represents a resolver that should be used for conflicting Notes
var NoteResolver string

func merge(leftFilename string, rightFilename string, mergedFilename string, stdio terminal.Stdio) {
	fmt.Println("Importing left backup")
	left := model.Database{}
	err := left.ImportJWLBackup(leftFilename)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprintln(stdio.Out, "Importing right backup")
	right := model.Database{}
	err = right.ImportJWLBackup(rightFilename)
	if err != nil {
		log.Fatal(err)
	}

	merged := model.Database{}

	fmt.Fprintln(stdio.Out, "🧭 Merging Locations")
	mergedLocations, locationIDChanges, err := merger.MergeLocations(left.Location, right.Location)
	merged.Location = mergedLocations
	merger.UpdateLRIDs(left.Bookmark, right.Bookmark, "LocationID", locationIDChanges)
	merger.UpdateLRIDs(left.Bookmark, right.Bookmark, "PublicationLocationID", locationIDChanges)
	merger.UpdateLRIDs(left.Note, right.Note, "LocationID", locationIDChanges)
	merger.UpdateLRIDs(left.TagMap, right.TagMap, "LocationID", locationIDChanges)
	merger.UpdateLRIDs(left.UserMark, right.UserMark, "LocationID", locationIDChanges)
	fmt.Fprintln(stdio.Out, "Done.")

	fmt.Fprintln(stdio.Out, "📑 Merging Bookmarks")
	var bookmarksConflictSolution map[string]merger.MergeSolution
	for {
		mergedBookmarks, _, err := merger.MergeBookmarks(left.Bookmark, right.Bookmark, bookmarksConflictSolution)
		if err == nil {
			merged.Bookmark = mergedBookmarks
			break
		}
		switch err := err.(type) {
		case merger.MergeConflictError:
			if BookmarkResolver != "" {
				var resErr error
				bookmarksConflictSolution, resErr = autoResolveConflicts(err.Conflicts, BookmarkResolver)
				if resErr != nil {
					log.Fatal(resErr)
				}
			} else {
				bookmarksConflictSolution = handleMergeConflict(err.Conflicts, &merged, stdio)
			}
		default:
			log.Fatal(err)
		}
	}
	fmt.Fprintln(stdio.Out, "Done.")

	fmt.Fprintln(stdio.Out, "🏷  Merging Tags")
	var tagsConflictSolution map[string]merger.MergeSolution
	for {
		mergedTags, tagIDChanges, err := merger.MergeTags(left.Tag, right.Tag, tagsConflictSolution)
		if err == nil {
			merged.Tag = mergedTags
			merger.UpdateLRIDs(left.TagMap, right.TagMap, "TagID", tagIDChanges)
			break
		}
		switch err := err.(type) {
		case merger.MergeConflictError:
			tagsConflictSolution = handleMergeConflict(err.Conflicts, nil, stdio) // TODO
		default:
			log.Fatal(err)
		}
	}
	fmt.Fprintln(stdio.Out, "Done.")

	fmt.Fprintln(stdio.Out, "🖍  Merging Markings")
	var UMBRConflictSolution map[string]merger.MergeSolution
	for {
		mergedUserMarks, mergedBlockRanges, userMarkIDChanges, err := merger.MergeUserMarkAndBlockRange(left.UserMark, left.BlockRange, right.UserMark, right.BlockRange, UMBRConflictSolution)
		if err == nil {
			merged.UserMark = mergedUserMarks
			merged.BlockRange = mergedBlockRanges
			merger.UpdateLRIDs(left.Note, right.Note, "UserMarkID", userMarkIDChanges)
			break
		}
		switch err := err.(type) {
		case merger.MergeConflictError:
			if MarkingResolver != "" {
				var resErr error
				UMBRConflictSolution, resErr = autoResolveConflicts(err.Conflicts, MarkingResolver)
				if resErr != nil {
					log.Fatal(resErr)
				}
			} else {
				UMBRConflictSolution = handleMergeConflict(err.Conflicts, &merged, stdio)
			}
		default:
			log.Fatal(err)
		}
	}
	fmt.Fprintln(stdio.Out, "Done.")

	fmt.Fprintln(stdio.Out, "📝 Merging Notes")
	var notesConflictSolution map[string]merger.MergeSolution
	for {
		mergedNotes, notesIDChanges, err := merger.MergeNotes(left.Note, right.Note, notesConflictSolution)
		if err == nil {
			merged.Note = mergedNotes
			merger.UpdateLRIDs(left.TagMap, right.TagMap, "NoteID", notesIDChanges)
			break
		}
		switch err := err.(type) {
		case merger.MergeConflictError:
			if NoteResolver != "" {
				var resErr error
				notesConflictSolution, resErr = autoResolveConflicts(err.Conflicts, NoteResolver)
				if resErr != nil {
					log.Fatal(resErr)
				}
			} else {
				notesConflictSolution = handleMergeConflict(err.Conflicts, &merged, stdio)
			}
		default:
			log.Fatal(err)
		}
	}
	fmt.Fprintln(stdio.Out, "Done.")

	fmt.Fprintln(stdio.Out, "🏷  Merging TagMaps")
	var tagMapsConflictSolution map[string]merger.MergeSolution
	for {
		mergedTagMaps, _, err := merger.MergeTagMaps(left.TagMap, right.TagMap, tagMapsConflictSolution)
		if err == nil {
			merged.TagMap = mergedTagMaps
			break
		}
		switch err := err.(type) {
		case merger.MergeConflictError:
			tagMapsConflictSolution = handleMergeConflict(err.Conflicts, nil, stdio)
		default:
			log.Fatal(err)
		}
	}
	fmt.Fprintln(stdio.Out, "Done.")

	fmt.Fprintln(stdio.Out, "🎉 Finished merging!")

	fmt.Fprintln(stdio.Out, "Exporting merged database")
	if err = merged.ExportJWLBackup(mergedFilename); err != nil {
		log.Fatal(err)
	}

}

func handleMergeConflict(conflicts map[string]merger.MergeConflict, mergedDB *model.Database, stdio terminal.Stdio) map[string]merger.MergeSolution {
	helpText := ""
	for _, val := range conflicts {
		helpText = mergeConflictHelp(reflect.TypeOf(val.Left).String())
		break
	}

	prompt := &survey.Select{
		Message: "Select which side should be chosen:",
		Options: []string{"Left", "Right"},
		Help:    helpText,
	}

	result := make(map[string]merger.MergeSolution, len(conflicts))
	for key, conflict := range conflicts {
		t := table.NewWriter()
		t.SetStyle(table.StyleRounded)
		t.Style().Options = table.Options{
			DrawBorder:      true,
			SeparateColumns: true,
			SeparateFooter:  true,
			SeparateHeader:  true,
			SeparateRows:    true,
		}

		t.SetOutputMirror(os.Stdout)
		if goterm.Width() >= 190 {
			t.AppendHeader(table.Row{"Left", "Right"})
			t.AppendRow([]interface{}{conflict.Left.PrettyPrint(mergedDB), conflict.Right.PrettyPrint(mergedDB)})
		} else {
			t.AppendRows([]table.Row{{"Left"}, {conflict.Left.PrettyPrint(mergedDB)}, {"Right"}, {conflict.Right.PrettyPrint(mergedDB)}})
		}

		t.Render()

		fmt.Fprint(stdio.Out, "\n\n")

		var selected string
		err := survey.AskOne(prompt, &selected, survey.WithStdio(stdio.In, stdio.Out, stdio.Err))
		if err == terminal.InterruptErr {
			fmt.Fprintln(stdio.Out, "interrupted")
			os.Exit(0)
		} else if err != nil {
			panic(err)
		}

		if selected == "Left" {
			result[key] = merger.MergeSolution{
				Side:      merger.LeftSide,
				Solution:  conflict.Left,
				Discarded: conflict.Right,
			}
		} else {
			result[key] = merger.MergeSolution{
				Side:      merger.RightSide,
				Solution:  conflict.Right,
				Discarded: conflict.Left,
			}
		}
	}

	return result
}

// autoResolveConflicts resolves mergeConflicts using the resolver
// indicated by resolverName.
func autoResolveConflicts(conflicts map[string]merger.MergeConflict, resolverName string) (map[string]merger.MergeSolution, error) {
	resolver, err := getResolver(resolverName)
	if err != nil {
		return nil, err
	}
	if resolver == nil {
		return nil, nil
	}
	return resolver(conflicts)
}

// getResolver parses the name of the resolver and returns its function.
// If the name is empty, it returns nil.
func getResolver(name string) (merger.MergeConflictSolver, error) {
	if name == "" {
		return nil, nil
	}

	switch name {
	case "chooseLeft":
		return merger.SolveConflictByChoosingLeft, nil
	case "chooseRight":
		return merger.SolveConflictByChoosingRight, nil
	case "chooseNewest":
		return merger.SolveConflictByChoosingNewest, nil
	}

	return nil, fmt.Errorf("%s is not a valid conflict resolver. Can be 'chooseNewest', 'chooseLeft', or 'chooseRight'", name)
}

func init() {
	rootCmd.AddCommand(mergeCmd)
	mergeCmd.Flags().StringVar(&BookmarkResolver, "bookmarks", "", "Resolve conflicting bookmarks with resolver (can be 'chooseLeft' or 'chooseRight')")
	mergeCmd.Flags().StringVar(&MarkingResolver, "markings", "", "Resolve conflicting markings with resolver (can be 'chooseLeft' or 'chooseRight')")
	mergeCmd.Flags().StringVar(&NoteResolver, "notes", "", "Resolve conflicting notes with resolver (can be 'chooseNewest', 'chooseLeft', or 'chooseRight')")
}
