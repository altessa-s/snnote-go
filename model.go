// Copyright 2026 ALTESSA SOLUTIONS INC.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file.

package snnote

// Point preserves the point information supplied by the file.
type Point struct {
	// X holds the horizontal position, in page pixels from the left edge.
	X float64
	// Y holds the vertical position, in page pixels from the top edge.
	Y float64
}

// Size preserves the size information supplied by the file.
type Size struct {
	// Width holds the horizontal extent, in pixels.
	Width float64
	// Height holds the vertical extent, in pixels.
	Height float64
}

// Rect preserves the rect information supplied by the file.
type Rect struct {
	// X holds the horizontal position, in page pixels from the left edge.
	X float64
	// Y holds the vertical position, in page pixels from the top edge.
	Y float64
	// Width holds the horizontal extent, in pixels.
	Width float64
	// Height holds the vertical extent, in pixels.
	Height float64
}

// MetadataEntry preserves the metadata entry information supplied by the file.
type MetadataEntry struct {
	// Key holds the UTF-8 metadata key, compared bytewise.
	Key string
	// Value holds the stored UTF-8 value; colons after the first delimiter are retained.
	Value string
}

// Metadata preserves the metadata information supplied by the file.
type Metadata struct {
	// Entries holds the entries in file order; repeated keys are retained.
	Entries []MetadataEntry
}

// Bitmap preserves the bitmap information supplied by the file.
type Bitmap struct {
	// Encoding holds the detected image signature takes precedence over the declared protocol.
	Encoding Encoding
	// Name holds the unknown protocol spelling; absent for recognized encodings.
	Name *string
	// Bytes holds the encoded payload; not decoded pixels.
	Bytes []byte
}

// Ink preserves the ink information supplied by the file.
type Ink struct {
	// Gray holds the gray intensity, from 0 (black) to 255 (white).
	Gray uint8
	// Alpha holds the coverage, from 0 (transparent) to 255 (opaque).
	Alpha uint8
}

// Palette preserves the palette information supplied by the file.
type Palette struct {
	// Inks maps every possible decoded code to coverage and intensity.
	Inks [256]Ink
}

// LayerInfo preserves the layer info information supplied by the file.
type LayerInfo struct {
	// LayerID holds the stored layer identifier: -1 background, 0 main, 1..3 other ink layers.
	LayerID int64
	// Name holds the stored name; always present, empty when the table gives none.
	Name string
	// IsBackground reports whether the layer-table row identifies the template layer.
	IsBackground bool
	// IsVisible reports whether the layer is shown; deleted or unstacked layers are hidden.
	IsVisible bool
	// IsDeleted reports whether the layer-table row marks the layer deleted.
	IsDeleted bool
	// IsCurrent reports whether this was the active editing layer.
	IsCurrent bool
}

// Layer preserves the layer information supplied by the file.
type Layer struct {
	// Name holds the layer identity derived from its page metadata key.
	Name LayerName
	// Bitmap holds the encoded bitmap, or nil when absent.
	Bitmap *Bitmap
	// Info holds the matching layer-table row, or nil when unavailable.
	Info *LayerInfo
	// Metadata holds the metadata in stored order, including unknown keys.
	Metadata Metadata
	// IsStacked reports whether LAYERSEQ includes this layer; absent sequence includes all layers.
	IsStacked bool
	// IsVisible reports whether the layer is shown; deleted or unstacked layers are hidden.
	IsVisible bool
}

// StrokePoint preserves the stroke point information supplied by the file.
type StrokePoint struct {
	// Position holds the point in page pixels, x rightward and y downward from the top left.
	Position Point
	// Pressure holds the raw digitizer pressure; zero when the pressure array has no matching entry.
	Pressure uint16
}

// Contour preserves the contour information supplied by the file.
type Contour struct {
	// Points holds the points, in page pixel coordinates.
	Points []Point
}

// Stroke preserves the stroke information supplied by the file.
type Stroke struct {
	// ID holds the stored record identifier.
	ID int32
	// Pen holds the raw pen code; known values are Pen constants.
	Pen Pen
	// RoleRaw holds the stored role code, retained even when unknown.
	RoleRaw int32
	// StatusRaw holds the stored status code, retained even when unknown.
	StatusRaw int32
	// Operation holds the stored operation code; no interpretation is imposed.
	Operation int32
	// Color holds the raw ink color code; InkColor.Gray maps known grayscale values.
	Color InkColor
	// Thickness holds the stored pen thickness, in hundredths of a page pixel.
	Thickness uint32
	// Layer holds the stored ink-layer identifier, 0 main and 1..3 other layers.
	Layer uint32
	// PageNumber holds the stored one-based page number; zero when absent or undecodable.
	PageNumber uint32
	// KindName holds the stored stroke-kind spelling, including unknown names.
	KindName string
	// Kind holds the interpreted record kind; raw spelling/code is retained separately where available.
	Kind StrokeKind
	// Role holds the interpreted role; unknown codes are RoleOther.
	Role Role
	// Status holds the interpreted stroke status; raw code remains in StatusRaw.
	Status StrokeStatus
	// Category holds the category derived from the optional style range.
	Category Category
	// IsDrawn reports whether this record represents visible ink before page-layer visibility is applied.
	IsDrawn bool
	// Style holds the stored style code; no clamping of unknown values.
	Style *uint32
	// Points holds the points, in page pixel coordinates.
	Points []StrokePoint
	// Pressures holds the complete raw pressure array, even when its length differs from points.
	Pressures []uint16
	// Contours holds the outline rings in page pixels; empty when the outline cannot be read.
	Contours []Contour
	// Bounds holds the bounds in page pixels; frame and heading bounds are derived from points.
	Bounds *Rect
	// Attributes holds the decoded tail fields, in stored order; unknown positions are retained.
	Attributes []string
}

// TextBox preserves the text box information supplied by the file.
type TextBox struct {
	// ID holds the stored record identifier.
	ID int32
	// Rect holds the rectangle in page pixels from the top-left corner.
	Rect *Rect
	// Text holds the UTF-8 content; always present, possibly empty.
	Text string
	// FontPath holds the stored font reference, or absent; the library never opens it.
	FontPath *string
	// Timestamp holds the stored timestamp spelling, or absent; no timezone conversion is applied.
	Timestamp *string
	// FontSize holds the stored text-box font size; nil when absent.
	FontSize *float64
	// Attributes holds the decoded tail fields, in stored order; unknown positions are retained.
	Attributes []string
}

// Word preserves the word information supplied by the file.
type Word struct {
	// Label holds the recognized text; always present, possibly empty.
	Label string
	// Box holds the recognition bounds in physical millimetres from the top left; nil when no box was supplied.
	Box *Rect
}

// Element preserves the element information supplied by the file.
type Element struct {
	// Kind holds the stored recognition type spelling, empty when absent.
	Kind string
	// Label holds the recognized text; absence and empty text are distinct.
	Label *string
	// Words holds the recognized words in document order.
	Words []Word
}

// Recognition preserves the recognition information supplied by the file.
type Recognition struct {
	// Elements holds the recognition elements in document order, including non-text elements.
	Elements []Element
	// Text holds the UTF-8 content; always present, possibly empty.
	Text string
}

// Page preserves the page information supplied by the file.
type Page struct {
	// Number holds the one-based numeric page key from the footer.
	Number int64
	// ID holds the stored PAGEID, or absent.
	ID *string
	// TemplateName holds the stored PAGESTYLE name, or absent.
	TemplateName *string
	// Orientation holds the interpreted page orientation; missing or zero is portrait.
	Orientation Orientation
	// OrientationRaw holds the stored orientation code, or zero when missing.
	OrientationRaw int64
	// Size holds the oriented page dimensions in pixels; landscape swaps the portrait sheet axes.
	Size Size
	// Layers holds the layers, in front-to-back stack order; unstacked layers follow.
	Layers []Layer
	// LayerInfo holds the layer-table rows in stored order.
	LayerInfo []LayerInfo
	// Strokes holds the readable stroke records, including non-drawn and deleted records.
	Strokes []Stroke
	// SkippedStrokes holds the number of unreadable records; one when the complete log cannot be decoded.
	SkippedStrokes int
	// DrawnStrokes holds pointers into the original Strokes backing array, filtered
	// by classification and layer visibility. Appending may allocate a new array;
	// these pointers still refer to the old records. Re-slicing does not retarget
	// or remove pointers, even to records excluded from the new slice. Reordering
	// or overwriting Strokes changes the records at the pointed-to positions.
	DrawnStrokes []*Stroke
	// TextBoxes holds the live text boxes derived from stroke tails.
	TextBoxes []TextBox
	// Recognition holds the recognition model, or nil when absent or unreadable.
	Recognition *Recognition
	// RecognitionArchive holds the opaque recognition archive; no interpretation is imposed.
	RecognitionArchive []byte
	// DisabledAreas holds the text-box exclusion rectangles in page pixels.
	DisabledAreas []*Rect
	// Metadata holds the metadata in stored order, including unknown keys.
	Metadata Metadata
}

// Title preserves the title information supplied by the file.
type Title struct {
	// PageNumber holds the stored one-based page number; zero when absent or undecodable.
	PageNumber int64
	// Level holds the stored title level; defaults to one.
	Level int64
	// Style holds the stored style code; no clamping of unknown values.
	Style TitleStyle
	// Rect holds the rectangle in page pixels from the top-left corner.
	Rect *Rect
	// OriginalRect holds the stored original rectangle in page pixels; falls back to rect when absent.
	OriginalRect *Rect
	// Bitmap holds the encoded bitmap, or nil when absent.
	Bitmap *Bitmap
	// Metadata holds the metadata in stored order, including unknown keys.
	Metadata Metadata
}

// Keyword preserves the keyword information supplied by the file.
type Keyword struct {
	// PageNumber holds the stored one-based page number; zero when absent or undecodable.
	PageNumber int64
	// Rect holds the rectangle in page pixels from the top-left corner.
	Rect *Rect
	// OriginalRect holds the stored original rectangle in page pixels; falls back to rect when absent.
	OriginalRect *Rect
	// Text holds the UTF-8 content; always present, possibly empty.
	Text string
	// Metadata holds the metadata in stored order, including unknown keys.
	Metadata Metadata
}

// Link preserves the link information supplied by the file.
type Link struct {
	// PageNumber holds the stored one-based page number; zero when absent or undecodable.
	PageNumber int64
	// Kind holds the raw LINKTYPE value: 4 identifies a web link.
	Kind int64
	// Direction holds the raw LINKINOUT value; zero when absent.
	Direction int64
	// TargetPageNumber holds the stored target page number; nil when absent.
	TargetPageNumber *int64
	// Style holds the stored style code; no clamping of unknown values.
	Style int64
	// Rect holds the rectangle in page pixels from the top-left corner.
	Rect *Rect
	// TargetPath holds the decoded stored link target, or absent; never opened by this library.
	TargetPath *string
	// TargetFileID holds the stored target file identifier, or absent.
	TargetFileID *string
	// TargetPageID holds the stored target page identifier, or absent.
	TargetPageID *string
	// Timestamp holds the stored timestamp spelling, or absent; no timezone conversion is applied.
	Timestamp *string
	// Text holds the UTF-8 content; a nil pointer means absent, distinct from empty.
	Text *string
	// URL holds the alias of TargetPath for web links (kind 4), absent otherwise.
	URL *string
	// Bitmap holds the encoded bitmap, or nil when absent.
	Bitmap *Bitmap
	// Metadata holds the metadata in stored order, including unknown keys.
	Metadata Metadata
}

// Highlight preserves the highlight information supplied by the file.
type Highlight struct {
	// PageNumber holds the stored one-based page number; zero when absent or undecodable.
	PageNumber int64
	// Kind holds the raw annotation type; zero when absent.
	Kind int64
	// ColorType holds the raw highlight color code; zero when absent.
	ColorType int64
	// Comment holds the comment; absent for a missing or empty comment.
	Comment *string
	// Rects holds the screen rectangles in page pixels.
	Rects []*Rect
	// DocumentRects holds the rectangles in the imported document coordinate space.
	DocumentRects []*Rect
	// Date holds the seconds since Unix epoch, converted from the stored milliseconds.
	Date *float64
}

// Template preserves the template information supplied by the file.
type Template struct {
	// Name holds the stored name; always present.
	Name string
	// Bitmap holds the encoded bitmap owned by the containing notebook.
	Bitmap Bitmap
}

// Notebook preserves the notebook information supplied by the file.
type Notebook struct {
	// Palette holds the shared code-to-gray mapping for every page bitmap, chosen by format version.
	Palette Palette
	// Kind holds the container signature: ordinary note or document annotation mark.
	Kind Kind
	// FormatVersion holds the decimal SN_FILE_VER number, retained without truncating its year prefix.
	FormatVersion int64
	// Equipment holds the recognized device identity, or EquipmentOther.
	Equipment Equipment
	// EquipmentRaw holds the stored APPLY_EQUIPMENT spelling, empty when missing.
	EquipmentRaw string
	// FileID holds the stored FILE_ID, or absent.
	FileID *string
	// Sheet holds the portrait sheet dimensions in pixels; each page may select a different sheet.
	Sheet Size
	// Pages holds the readable pages, sorted by numeric page key; last duplicate wins.
	Pages []Page
	// SkippedPages holds the number of indexed pages that could not be read.
	SkippedPages int
	// Titles holds the titles; duplicates replace by footer key and rectangle.
	Titles []Title
	// Keywords holds the keywords; duplicates replace by footer key and rectangle.
	Keywords []Keyword
	// Links holds the links, in footer order.
	Links []Link
	// Highlights holds the document highlights, sorted by numeric page key.
	Highlights []Highlight
	// Templates holds the templates; last bitmap wins while first insertion position remains.
	Templates []Template
	// Cover holds the encoded cover, or nil when absent.
	Cover *Bitmap
	// CoverType holds the numeric suffix of the selected COVER_ footer key, or zero.
	CoverType int64
	// Header holds the FILE_FEATURE metadata.
	Header Metadata
	// Footer holds the footer metadata.
	Footer Metadata
}
