CREATE TABLE BlockRange ( BlockRangeId    INTEGER NOT NULL PRIMARY KEY, BlockType       INTEGER NOT NULL, Identifier      INTEGER NOT NULL, StartToken      INTEGER, EndToken        INTEGER, UserMarkId      INTEGER NOT NULL, CHECK (BlockType BETWEEN 1 AND 2), FOREIGN KEY(UserMarkId) REFERENCES UserMark(UserMarkId) );
CREATE TABLE IF NOT EXISTS "Bookmark" (
    BookmarkId              INTEGER NOT NULL PRIMARY KEY,
    LocationId              INTEGER NOT NULL,
    PublicationLocationId   INTEGER NOT NULL,
    Slot                    INTEGER NOT NULL,
    Title                   TEXT NOT NULL,
    Snippet                 TEXT,
    BlockType               INTEGER NOT NULL DEFAULT 0,
    BlockIdentifier         INTEGER,
    FOREIGN KEY(LocationId) REFERENCES Location(LocationId),
    FOREIGN KEY(PublicationLocationId) REFERENCES Location(LocationId),
    CONSTRAINT PublicationLocationId_Slot UNIQUE (PublicationLocationId, Slot),
    CHECK((BlockType = 0 AND BlockIdentifier IS NULL) OR ((BlockType BETWEEN 1 AND 2) AND BlockIdentifier IS NOT NULL))
);
CREATE TABLE IF NOT EXISTS "TagMap" (
    TagMapId          INTEGER NOT NULL PRIMARY KEY,
    PlaylistItemId    INTEGER,
    LocationId        INTEGER,
    NoteId            INTEGER,
    TagId             INTEGER NOT NULL,
    Position          INTEGER NOT NULL,
    FOREIGN KEY(TagId) REFERENCES Tag(TagId),
    FOREIGN KEY(PlaylistItemId) REFERENCES PlaylistItem(PlaylistItemId),
    FOREIGN KEY(LocationId) REFERENCES Location(LocationId),
    FOREIGN KEY(NoteId) REFERENCES Note(NoteId),
    CONSTRAINT TagId_Position UNIQUE(TagId, Position),
    CONSTRAINT TagId_NoteId UNIQUE(TagId, NoteId),
    CONSTRAINT TagId_LocationId UNIQUE(TagId, LocationId),
    CHECK(
        (NoteId IS NULL AND LocationId IS NULL AND PlaylistItemId IS NOT NULL) OR
        (LocationId IS NULL AND PlaylistItemId IS NULL AND NoteId IS NOT NULL) OR
        (PlaylistItemId IS NULL AND NoteId IS NULL AND LocationId IS NOT NULL))
);
CREATE TABLE IF NOT EXISTS "UserMark" (
    UserMarkId      INTEGER NOT NULL PRIMARY KEY,
    ColorIndex      INTEGER NOT NULL,
    LocationId      INTEGER NOT NULL,
    StyleIndex      INTEGER NOT NULL,
    UserMarkGuid    TEXT NOT NULL UNIQUE,
    Version         INTEGER NOT NULL,
    FOREIGN KEY(LocationId) REFERENCES Location(LocationId)
);
CREATE TABLE InputField (
    LocationId  INTEGER NOT NULL,
    TextTag     TEXT NOT NULL,
    Value       TEXT NOT NULL,
    FOREIGN KEY (LocationId) REFERENCES Location (LocationId),
    CONSTRAINT LocationId_TextTag PRIMARY KEY (LocationId, TextTag)
);
CREATE TABLE grdb_migrations (identifier TEXT NOT NULL PRIMARY KEY);
CREATE TABLE PlaylistItemAccuracy(
    PlaylistItemAccuracyId  INTEGER NOT NULL PRIMARY KEY,
    Description             TEXT NOT NULL UNIQUE);
CREATE TABLE IndependentMedia(
    IndependentMediaId  INTEGER NOT NULL PRIMARY KEY,
    OriginalFilename    TEXT NOT NULL,
    FilePath            TEXT NOT NULL UNIQUE,
    MimeType            TEXT NOT NULL,
    Hash                TEXT NOT NULL,
    CHECK(length(OriginalFilename) > 0),
    CHECK(length(FilePath) > 0),
    CHECK(length(MimeType) > 0),
    CHECK(length(Hash) > 0));
CREATE TABLE PlaylistItemIndependentMediaMap(
    PlaylistItemId      INTEGER NOT NULL,
    IndependentMediaId  INTEGER NOT NULL,
    DurationTicks       INTEGER NOT NULL,
    PRIMARY KEY(PlaylistItemId, IndependentMediaId),
    FOREIGN KEY(PlaylistItemId) REFERENCES PlaylistItem(PlaylistItemId),
    FOREIGN KEY(IndependentMediaId) REFERENCES IndependentMedia(IndependentMediaId))
WITHOUT ROWID;
CREATE TABLE PlaylistItemLocationMap(
    PlaylistItemId      INTEGER NOT NULL,
    LocationId          INTEGER NOT NULL,
    MajorMultimediaType INTEGER NOT NULL,
    BaseDurationTicks   INTEGER,
    PRIMARY KEY(PlaylistItemId, LocationId),
    FOREIGN KEY(PlaylistItemId) REFERENCES PlaylistItem(PlaylistItemId),
    FOREIGN KEY(LocationId) REFERENCES Location(LocationId))
WITHOUT ROWID;
CREATE TABLE PlaylistItemMarker(
    PlaylistItemMarkerId        INTEGER NOT NULL PRIMARY KEY,
    PlaylistItemId              INTEGER NOT NULL,
    Label                       TEXT NOT NULL,
    StartTimeTicks              INTEGER NOT NULL,
    DurationTicks               INTEGER NOT NULL,
    EndTransitionDurationTicks  INTEGER NOT NULL,
    UNIQUE(PlaylistItemId, StartTimeTicks),
    FOREIGN KEY(PlaylistItemId) REFERENCES PlaylistItem(PlaylistItemId));
CREATE TABLE PlaylistItemMarkerBibleVerseMap(
    PlaylistItemMarkerId        INTEGER NOT NULL,
    VerseId                     INTEGER NOT NULL,
    PRIMARY KEY(PlaylistItemMarkerId, VerseId),
    FOREIGN KEY(PlaylistItemMarkerId) REFERENCES PlaylistItemMarker(PlaylistItemMarkerId))
WITHOUT ROWID;
CREATE TABLE PlaylistItemMarkerParagraphMap(
    PlaylistItemMarkerId        INTEGER NOT NULL,
    MepsDocumentId              INTEGER NOT NULL,
    ParagraphIndex              INTEGER NOT NULL,
    MarkerIndexWithinParagraph  INTEGER NOT NULL,
    PRIMARY KEY(PlaylistItemMarkerId, MepsDocumentId, ParagraphIndex, MarkerIndexWithinParagraph),
    FOREIGN KEY(PlaylistItemMarkerId) REFERENCES PlaylistItemMarker(PlaylistItemMarkerId))
WITHOUT ROWID;
CREATE TABLE Tag(
    TagId          INTEGER NOT NULL PRIMARY KEY,
    Type           INTEGER NOT NULL,
    Name           TEXT NOT NULL,
    UNIQUE(Type, Name),
    CHECK(length(Name) > 0),
    CHECK(Type IN (0, 1, 2)));
CREATE TABLE IF NOT EXISTS "Note"(
    NoteId           INTEGER NOT NULL PRIMARY KEY,
    Guid             TEXT NOT NULL UNIQUE,
    UserMarkId       INTEGER,
    LocationId       INTEGER,
    Title            TEXT,
    Content          TEXT,
    LastModified     TEXT NOT NULL DEFAULT(strftime('%Y-%m-%dT%H:%M:%SZ', 'now')),
    Created          TEXT NOT NULL DEFAULT(strftime('%Y-%m-%dT%H:%M:%SZ', 'now')),
    BlockType        INTEGER NOT NULL DEFAULT 0,
    BlockIdentifier  INTEGER,
    CHECK((BlockType = 0 AND BlockIdentifier IS NULL) OR((BlockType BETWEEN 1 AND 2) AND BlockIdentifier IS NOT NULL)),
    FOREIGN KEY(UserMarkId) REFERENCES UserMark(UserMarkId),
    FOREIGN KEY(LocationId) REFERENCES Location(LocationId));
CREATE TABLE IF NOT EXISTS "LastModified"(LastModified TEXT NOT NULL);
CREATE TABLE IF NOT EXISTS "PlaylistItem"(
    PlaylistItemId           INTEGER NOT NULL PRIMARY KEY,
    Label                    TEXT NOT NULL,
    StartTrimOffsetTicks     INTEGER,
    EndTrimOffsetTicks       INTEGER,
    Accuracy                 INTEGER NOT NULL,
    EndAction                INTEGER NOT NULL,
    ThumbnailFilePath        TEXT,
    FOREIGN KEY(Accuracy) REFERENCES PlaylistItemAccuracy(PlaylistItemAccuracyId),
    FOREIGN KEY(ThumbnailFilePath) REFERENCES IndependentMedia(FilePath),
    CHECK(length(Label) > 0),
    CHECK(EndAction IN(0, 1, 2, 3)));
CREATE TABLE IF NOT EXISTS "Location"(
    LocationId      INTEGER NOT NULL PRIMARY KEY,
    BookNumber      INTEGER,
    ChapterNumber   INTEGER,
    DocumentId      INTEGER,
    Track           INTEGER,
    IssueTagNumber  INTEGER NOT NULL DEFAULT 0,
    KeySymbol       TEXT,
    MepsLanguage    INTEGER,
    Type            INTEGER NOT NULL,
    Title           TEXT,
    UNIQUE(BookNumber, ChapterNumber, KeySymbol, MepsLanguage, Type),
    UNIQUE(KeySymbol, IssueTagNumber, MepsLanguage, DocumentId, Track, Type),
    CHECK ((
        Type = 0 AND ( -- Document or Bible chapter
            (DocumentId IS NOT NULL AND DocumentId != 0)
                OR ( -- Track based. Requires DocumentId or KeySymbol.
                    Track IS NOT NULL AND (
                        (KeySymbol IS NOT NULL AND (length(KeySymbol) > 0))
                        OR (DocumentId IS NOT NULL AND DocumentId != 0)))
                OR ( -- Bible book. Requires KeySymbol.
                    BookNumber IS NOT NULL AND BookNumber != 0
                    AND KeySymbol IS NOT NULL AND (length(KeySymbol) > 0)
                    AND (ChapterNumber IS NULL OR ChapterNumber = 0))
                OR ( -- Bible chapter. Requires KeySymbol and BookNumber.
                    ChapterNumber IS NOT NULL AND ChapterNumber != 0
                    AND BookNumber IS NOT NULL AND BookNumber != 0
                    AND KeySymbol IS NOT NULL AND (length(KeySymbol) > 0))))
        OR Type != 0),
    CHECK((
        Type = 1 -- Bible
            AND (BookNumber IS NULL OR BookNumber = 0)
            AND (ChapterNumber IS NULL OR ChapterNumber = 0)
            AND (DocumentId IS NULL OR DocumentId = 0)
            AND KeySymbol IS NOT NULL AND (length(KeySymbol) > 0)
            AND Track IS NULL)
        OR Type != 1),
    CHECK((
        Type IN (2, 3) -- Mediator audio/video
            AND (BookNumber IS NULL OR BookNumber = 0)
            AND (ChapterNumber IS NULL OR ChapterNumber = 0))
        OR Type NOT IN (2, 3)));
CREATE INDEX IX_BlockRange_UserMarkId ON BlockRange(UserMarkId);
CREATE INDEX IX_TagMap_TagId ON TagMap(TagId);
CREATE INDEX IX_TagMap_PlaylistItemId_TagId_Position ON TagMap(PlaylistItemId, TagId, Position);
CREATE INDEX IX_TagMap_LocationId_TagId_Position ON TagMap(LocationId, TagId, Position);
CREATE INDEX IX_TagMap_NoteId_TagId_Position ON TagMap(NoteId, TagId, Position);
CREATE INDEX IX_UserMark_LocationId ON UserMark(LocationId);
CREATE INDEX IX_Note_LastModified_LocationId ON Note(LastModified, LocationId);
CREATE INDEX IX_Note_LocationId_BlockIdentifier ON Note(LocationId, BlockIdentifier);
CREATE INDEX IX_PlaylistItemIndependentMediaMap_IndependentMediaId ON PlaylistItemIndependentMediaMap(IndependentMediaId);
CREATE INDEX IX_PlaylistItemLocationMap_LocationId ON PlaylistItemLocationMap(LocationId);
CREATE INDEX IX_Location_KeySymbol_MepsLanguage_BookNumber_ChapterNumber ON
    Location(KeySymbol, MepsLanguage, BookNumber, ChapterNumber);
CREATE INDEX IX_Location_MepsLanguage_DocumentId ON Location(MepsLanguage, DocumentId);
CREATE INDEX IX_PlaylistItem_ThumbnailFilePath ON PlaylistItem(ThumbnailFilePath);
CREATE TRIGGER TR_Update_LastModified_Insert_Tag
    INSERT ON Tag
    BEGIN
        UPDATE LastModified SET LastModified = strftime('%Y-%m-%dT%H:%M:%SZ', 'now');
    END;
CREATE TRIGGER TR_Update_LastModified_Update_Tag
    UPDATE ON Tag
    BEGIN
        UPDATE LastModified SET LastModified = strftime('%Y-%m-%dT%H:%M:%SZ', 'now');
    END;
CREATE TRIGGER TR_Update_LastModified_Delete_Tag
    DELETE ON Tag
    BEGIN
        UPDATE LastModified SET LastModified = strftime('%Y-%m-%dT%H:%M:%SZ', 'now');
    END;
CREATE TRIGGER TR_Update_LastModified_Insert_TagMap
    INSERT ON TagMap
    BEGIN
        UPDATE LastModified SET LastModified = strftime('%Y-%m-%dT%H:%M:%SZ', 'now');
    END;
CREATE TRIGGER TR_Update_LastModified_Update_TagMap
    UPDATE ON TagMap
    BEGIN
        UPDATE LastModified SET LastModified = strftime('%Y-%m-%dT%H:%M:%SZ', 'now');
    END;
CREATE TRIGGER TR_Update_LastModified_Delete_TagMap
    DELETE ON TagMap
    BEGIN
        UPDATE LastModified SET LastModified = strftime('%Y-%m-%dT%H:%M:%SZ', 'now');
    END;
CREATE TRIGGER TR_Update_LastModified_Insert_Note
    INSERT ON Note
    BEGIN
        UPDATE LastModified SET LastModified = strftime('%Y-%m-%dT%H:%M:%SZ', 'now');
    END;
CREATE TRIGGER TR_Update_LastModified_Update_Note
    UPDATE ON Note
    BEGIN
        UPDATE LastModified SET LastModified = strftime('%Y-%m-%dT%H:%M:%SZ', 'now');
    END;
CREATE TRIGGER TR_Update_LastModified_Delete_Note
    DELETE ON Note
    BEGIN
        UPDATE LastModified SET LastModified = strftime('%Y-%m-%dT%H:%M:%SZ', 'now');
    END;
CREATE TRIGGER TR_Update_LastModified_Insert_Bookmark
    INSERT ON Bookmark
    BEGIN
        UPDATE LastModified SET LastModified = strftime('%Y-%m-%dT%H:%M:%SZ', 'now');
    END;
CREATE TRIGGER TR_Update_LastModified_Update_Bookmark
    UPDATE ON Bookmark
    BEGIN
        UPDATE LastModified SET LastModified = strftime('%Y-%m-%dT%H:%M:%SZ', 'now');
    END;
CREATE TRIGGER TR_Update_LastModified_Delete_Bookmark
    DELETE ON Bookmark
    BEGIN
        UPDATE LastModified SET LastModified = strftime('%Y-%m-%dT%H:%M:%SZ', 'now');
    END;
CREATE TRIGGER TR_Update_LastModified_Insert_UserMark
    INSERT ON UserMark
    BEGIN
        UPDATE LastModified SET LastModified = strftime('%Y-%m-%dT%H:%M:%SZ', 'now');
    END;
CREATE TRIGGER TR_Update_LastModified_Update_UserMark
    UPDATE ON UserMark
    BEGIN
        UPDATE LastModified SET LastModified = strftime('%Y-%m-%dT%H:%M:%SZ', 'now');
    END;
CREATE TRIGGER TR_Update_LastModified_Delete_UserMark
    DELETE ON UserMark
    BEGIN
        UPDATE LastModified SET LastModified = strftime('%Y-%m-%dT%H:%M:%SZ', 'now');
    END;
CREATE TRIGGER TR_Update_LastModified_Insert_BlockRange
    INSERT ON BlockRange
    BEGIN
        UPDATE LastModified SET LastModified = strftime('%Y-%m-%dT%H:%M:%SZ', 'now');
    END;
CREATE TRIGGER TR_Update_LastModified_Update_BlockRange
    UPDATE ON BlockRange
    BEGIN
        UPDATE LastModified SET LastModified = strftime('%Y-%m-%dT%H:%M:%SZ', 'now');
    END;
CREATE TRIGGER TR_Update_LastModified_Delete_BlockRange
    DELETE ON BlockRange
    BEGIN
        UPDATE LastModified SET LastModified = strftime('%Y-%m-%dT%H:%M:%SZ', 'now');
    END;
CREATE TRIGGER TR_Update_LastModified_Insert_InputField
    INSERT ON InputField
    BEGIN
        UPDATE LastModified SET LastModified = strftime('%Y-%m-%dT%H:%M:%SZ', 'now');
    END;
CREATE TRIGGER TR_Update_LastModified_Update_InputField
    UPDATE ON InputField
    BEGIN
        UPDATE LastModified SET LastModified = strftime('%Y-%m-%dT%H:%M:%SZ', 'now');
    END;
CREATE TRIGGER TR_Update_LastModified_Delete_InputField
    DELETE ON InputField
    BEGIN
        UPDATE LastModified SET LastModified = strftime('%Y-%m-%dT%H:%M:%SZ', 'now');
    END;
CREATE TRIGGER TR_Update_LastModified_Insert_IndependentMedia
    INSERT ON IndependentMedia
    BEGIN
        UPDATE LastModified SET LastModified = strftime('%Y-%m-%dT%H:%M:%SZ', 'now');
    END;
CREATE TRIGGER TR_Update_LastModified_Update_IndependentMedia
    UPDATE ON IndependentMedia
    BEGIN
        UPDATE LastModified SET LastModified = strftime('%Y-%m-%dT%H:%M:%SZ', 'now');
    END;
CREATE TRIGGER TR_Update_LastModified_Delete_IndependentMedia
    DELETE ON IndependentMedia
    BEGIN
        UPDATE LastModified SET LastModified = strftime('%Y-%m-%dT%H:%M:%SZ', 'now');
    END;
CREATE TRIGGER TR_Raise_Error_Before_Insert_LastModified
    BEFORE INSERT ON LastModified
    BEGIN
        SELECT RAISE (FAIL, 'INSERT INTO LastModified not allowed');
    END;
CREATE TRIGGER TR_Raise_Error_Before_Delete_LastModified
    BEFORE DELETE ON LastModified
    BEGIN
        SELECT RAISE (FAIL, 'DELETE FROM LastModified not allowed');
    END;