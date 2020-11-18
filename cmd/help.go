package cmd

import "github.com/MakeNowJust/heredoc"

func mergeConflictHelp(name string) string {
	helpTexts := map[string]string{
		"*model.Bookmark": `Bookmarks are set for each publication (i.e. a Watchtower issue or a Bible translation)
		and can be placed at ten different „slots“ (the colors you see in the app). A collision 
		happens, if two bookmarks are placed at the same slot of the same publication.  You are 
		able to choose if the bookmark on the left or right should be added to the merged backup. 
		The title and snippet may help you identify the bookmark. Also have a look at the 
		„Related Location“ for a more detailed insight: „BookNumber“ and „ChapterNumber“ 
		generally relate to a Bible book, while „DocumentID“ might be a Watchtower issue or a 
		different publication.`,

		"*model.UserMarkBlockRange": `Markings collide if they overlap at at least one point. To figure out where a marking 
		is located, look at the Identifier, Start- and EndToken. The Identifier generally 
		represents a paragraph in a publication or the verse in a Bible chapter. Start-, and 
		EndToken represent the beginning and the end of a marking, where words and punctuation 
		marks are counted as tokens (e.g. the sentence “You are my witnesses,” contains five 
		tokens, as four words and one comma are counted). Note that a 
		marking can span multiple Identifiers. ColorIndex represents the color of the marking.
		After you located the collision, you can choose between the left and the right marking. 
		In future versions you will be able to merge them into one big marking.`,

		"*model.Note": `A note collides if it exists on both sides (so they must have been synced at least once) 
		and it differers in the title or content. It generally makes sense to choose the note 
		with the newest date.`,
	}

	if text, ok := helpTexts[name]; ok {
		return heredoc.Doc(text)
	}

	return ""
}
