package cmd

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"text/tabwriter"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/AndreasSko/go-jwlm/merger"
	"github.com/AndreasSko/go-jwlm/model"
	"github.com/AndreasSko/go-jwlm/publication"
	"github.com/MakeNowJust/heredoc"
	"github.com/buger/goterm"
	"github.com/jedib0t/go-pretty/table"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

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
		catalogExists = checkCatalog()
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

var catalogExists bool = false
var catalogPath string = filepath.Join(viper.GetString("appData"), "catalog.db")

func merge(leftFilename string, rightFilename string, mergedFilename string, stdio terminal.Stdio) {
	fmt.Fprintln(stdio.Out, "Importing left backup")
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
	bookmarksConflictSolution := map[string]merger.MergeSolution{}
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
				newSolutions, resErr := merger.AutoResolveConflicts(err.Conflicts, BookmarkResolver)
				if resErr != nil {
					log.Fatal(resErr)
				}
				addToSolutions(bookmarksConflictSolution, newSolutions)
			} else {
				newSolutions := handleMergeConflict(err.Conflicts, &merged, stdio)
				addToSolutions(bookmarksConflictSolution, newSolutions)
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
	UMBRConflictSolution := map[string]merger.MergeSolution{}
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
				newSolutions, resErr := merger.AutoResolveConflicts(err.Conflicts, MarkingResolver)
				if resErr != nil {
					log.Fatal(resErr)
				}
				addToSolutions(UMBRConflictSolution, newSolutions)
			} else {
				newSolutions := handleMergeConflict(err.Conflicts, &merged, stdio)
				addToSolutions(UMBRConflictSolution, newSolutions)
			}
		default:
			log.Fatal(err)
		}
	}
	fmt.Fprintln(stdio.Out, "Done.")

	fmt.Fprintln(stdio.Out, "📝 Merging Notes")
	notesConflictSolution := map[string]merger.MergeSolution{}
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
				newSolutions, resErr := merger.AutoResolveConflicts(err.Conflicts, NoteResolver)
				if resErr != nil {
					log.Fatal(resErr)
				}
				addToSolutions(notesConflictSolution, newSolutions)
			} else {
				newSolutions := handleMergeConflict(err.Conflicts, &merged, stdio)
				addToSolutions(notesConflictSolution, newSolutions)
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

// addToSolutions adds new mergeSolutions to the existing map of mergeSolutions
func addToSolutions(solutions map[string]merger.MergeSolution, new map[string]merger.MergeSolution) {
	for key, value := range new {
		solutions[key] = value
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
		switch conflict.Left.(type) {
		case *model.Bookmark:
			printBookmarkConflict(conflict, mergedDB, stdio)
		default:
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
		}

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

func printBookmarkConflict(conflict merger.MergeConflict, mergedDB *model.Database, stdio terminal.Stdio) {
	fmt.Fprint(stdio.Out, "\n")
	buf := new(bytes.Buffer)
	w := tabwriter.NewWriter(buf, 0, 0, 1, ' ', 0)

	location := conflict.Left.RelatedEntries(mergedDB).Location
	if location != nil {
		if catalogExists {
			query := publication.Lookup{
				DocumentID:     int(location.DocumentID.Int32),
				KeySymbol:      location.KeySymbol.String,
				IssueTagNumber: location.IssueTagNumber,
				MepsLanguage:   location.MepsLanguage,
			}
			publ, err := publication.LookupPublication(catalogPath, query)
			if err == nil {
				fmt.Fprintf(w, "Publication:\t%s\n", publ.Title)
			}
		} else {
			if location.KeySymbol.Valid {
				fmt.Fprintf(w, "Key Symbol:\t%s\n", location.KeySymbol.String)
			}
			if location.DocumentID.Valid {
				fmt.Fprintf(w, "DocumentID:\t%d\n", location.DocumentID.Int32)
			}
			if location.IssueTagNumber != 0 {
				fmt.Fprintf(w, "IssueTagNumber:\t%d\n", location.IssueTagNumber)
			}
		}
	}
	fmt.Fprintf(w, "Slot:\t%d\n", conflict.Left.(*model.Bookmark).Slot)

	w.Flush()
	fmt.Fprint(stdio.Out, buf)

	printFields := []string{"Title", "Snippet"}
	prettyPrintConflictTable(conflict, printFields, stdio)
}

func printMarkingConflict(conflict merger.MergeConflict, mergedDB *model.Database, stdio terminal.Stdio) {
	fmt.Fprint(stdio.Out, "\n")
	buf := new(bytes.Buffer)
	w := tabwriter.NewWriter(buf, 0, 0, 1, ' ', 0)

	location := conflict.Left.RelatedEntries(mergedDB).Location
	if location != nil {
		if catalogExists {
			query := publication.Lookup{
				DocumentID:     int(location.DocumentID.Int32),
				KeySymbol:      location.KeySymbol.String,
				IssueTagNumber: location.IssueTagNumber,
				MepsLanguage:   location.MepsLanguage,
			}
			publ, err := publication.LookupPublication(catalogPath, query)
			if err == nil {
				fmt.Fprintf(w, "Publication:\t%s\n", publ.Title)
			}
		} else {
			if location.KeySymbol.Valid {
				fmt.Fprintf(w, "Key Symbol:\t%s\n", location.KeySymbol.String)
			}
			if location.DocumentID.Valid {
				fmt.Fprintf(w, "DocumentID:\t%d\n", location.DocumentID.Int32)
			}
			if location.IssueTagNumber != 0 {
				fmt.Fprintf(w, "IssueTagNumber:\t%d\n", location.IssueTagNumber)
			}
		}
		if location.Title.Valid {
			fmt.Fprintf(w, "Title:\t%s\n", location.Title.String)
		}
	}

	w.Flush()
	fmt.Fprint(stdio.Out, buf)

	printFields := []string{"Title", "Snippet", "Slot"}
	prettyPrintConflictTable(conflict, printFields, stdio)
}

func prettyPrintConflictTable(conflict merger.MergeConflict, fields []string, stdio terminal.Stdio) {
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
		t.AppendRow([]interface{}{model.PrettyPrint(conflict.Left, fields), model.PrettyPrint(conflict.Right, fields)})
	} else {
		t.AppendRows([]table.Row{{"Left"}, {model.PrettyPrint(conflict.Left, fields)}, {"Right"}, {model.PrettyPrint(conflict.Right, fields)}})
	}

	t.Render()
	fmt.Fprint(stdio.Out, "\n\n")
}

// checkCatalog makes sure that catalog.db exists and is not older
// than a month, otherwise it will ask the user if it should download
// the latest version. Its returned bool indicates, if catalog.db
// exists.
func checkCatalog() bool {
	download := false

	stat, err := os.Stat(catalogPath)
	if err == nil {
		old := time.Now().Add(-time.Hour * 24 * 30)

		if stat.ModTime().Before(old) {
			prompt := &survey.Confirm{
				Message: heredoc.Doc(`The catalogDB is older than a month and new publications might have been added.
				Should the newest catalog be downloaded?`),
			}
			survey.AskOne(prompt, &download)
		}
	}

	if _, err := os.Stat(catalogPath); os.IsNotExist(err) {
		prompt := &survey.Confirm{
			Message: heredoc.Doc(`The catalogDB doesn't exist yet. Should it be downloaded now?
			It might help you to more easily identify conflicts, but it's not necessary`),
		}
		survey.AskOne(prompt, &download)
	}

	if download {
		prgrsChan := make(chan publication.Progress)
		done := make(chan struct{})
		go func() {
			err := publication.DownloadCatalog(context.Background(), prgrsChan, catalogPath)
			if err != nil {
				log.Fatal(err)
			}
			done <- struct{}{}
		}()

		fmt.Println("Starting download of catalogDB")
		for progress := range prgrsChan {
			fmt.Printf("  transferred %v / %v bytes (%.2f%%)\n", progress.BytesComplete, progress.Size, 100*progress.Progress)
		}
		<-done
		fmt.Println("Done!")
	}
	return download || err == nil
}

func init() {
	rootCmd.AddCommand(mergeCmd)
	mergeCmd.Flags().StringVar(&BookmarkResolver, "bookmarks", "", "Resolve conflicting bookmarks with resolver (can be 'chooseLeft' or 'chooseRight')")
	mergeCmd.Flags().StringVar(&MarkingResolver, "markings", "", "Resolve conflicting markings with resolver (can be 'chooseLeft' or 'chooseRight')")
	mergeCmd.Flags().StringVar(&NoteResolver, "notes", "", "Resolve conflicting notes with resolver (can be 'chooseNewest', 'chooseLeft', or 'chooseRight')")
}
