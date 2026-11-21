// Copyright 2026 ALTESSA SOLUTIONS INC.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file.

package snnote

import "strconv"

// Encoding is the interpreted encoding. Unknown numeric values remain intact.
type Encoding int

const (
	// EncodingRattaRLE denotes ratta rle.
	EncodingRattaRLE Encoding = 0
	// EncodingPNG denotes png.
	EncodingPNG Encoding = 1
	// EncodingJPEG denotes jpeg.
	EncodingJPEG Encoding = 2
	// EncodingOther denotes other.
	EncodingOther Encoding = 3
)

// String names the value and retains unknown numbers in its fallback.
func (v Encoding) String() string {
	switch v {
	case EncodingRattaRLE:
		return "ratta-rle"
	case EncodingPNG:
		return "png"
	case EncodingJPEG:
		return "jpeg"
	case EncodingOther:
		return "other"
	default:
		return "encoding(" + strconv.Itoa(int(v)) + ")"
	}
}

// Equipment is the interpreted equipment. Unknown numeric values remain intact.
type Equipment int

const (
	// EquipmentOther denotes other.
	EquipmentOther Equipment = 0
	// EquipmentA5X denotes a5x.
	EquipmentA5X Equipment = 1
	// EquipmentA6X denotes a6x.
	EquipmentA6X Equipment = 2
	// EquipmentNomad denotes nomad.
	EquipmentNomad Equipment = 3
	// EquipmentManta denotes manta.
	EquipmentManta Equipment = 4
)

// String names the value and retains unknown numbers in its fallback.
func (v Equipment) String() string {
	switch v {
	case EquipmentOther:
		return "other"
	case EquipmentA5X:
		return "a5x"
	case EquipmentA6X:
		return "a6x"
	case EquipmentNomad:
		return "nomad"
	case EquipmentManta:
		return "manta"
	default:
		return "equipment(" + strconv.Itoa(int(v)) + ")"
	}
}

// Kind is the interpreted kind. Unknown numeric values remain intact.
type Kind int

const (
	// KindNote denotes note.
	KindNote Kind = 0
	// KindMark denotes mark.
	KindMark Kind = 1
)

// String names the value and retains unknown numbers in its fallback.
func (v Kind) String() string {
	switch v {
	case KindNote:
		return "note"
	case KindMark:
		return "mark"
	default:
		return "kind(" + strconv.Itoa(int(v)) + ")"
	}
}

// Orientation is the interpreted orientation. Unknown numeric values remain intact.
type Orientation int

const (
	// OrientationPortrait denotes portrait.
	OrientationPortrait Orientation = 0
	// OrientationLandscape denotes landscape.
	OrientationLandscape Orientation = 1
	// OrientationOther denotes other.
	OrientationOther Orientation = 2
)

// String names the value and retains unknown numbers in its fallback.
func (v Orientation) String() string {
	switch v {
	case OrientationPortrait:
		return "portrait"
	case OrientationLandscape:
		return "landscape"
	case OrientationOther:
		return "other"
	default:
		return "orientation(" + strconv.Itoa(int(v)) + ")"
	}
}

// LayerName is the interpreted layer name. Unknown numeric values remain intact.
type LayerName int

const (
	// LayerMain denotes main.
	LayerMain LayerName = 0
	// Layer1 denotes 1.
	Layer1 LayerName = 1
	// Layer2 denotes 2.
	Layer2 LayerName = 2
	// Layer3 denotes 3.
	Layer3 LayerName = 3
	// LayerBackground denotes background.
	LayerBackground LayerName = 4
)

// String names the value and retains unknown numbers in its fallback.
func (v LayerName) String() string {
	switch v {
	case LayerMain:
		return "main"
	case Layer1:
		return "1"
	case Layer2:
		return "2"
	case Layer3:
		return "3"
	case LayerBackground:
		return "background"
	default:
		return "layer_name(" + strconv.Itoa(int(v)) + ")"
	}
}

// StrokeKind is the interpreted stroke kind. Unknown numeric values remain intact.
type StrokeKind int

const (
	// StrokeKindFreehand denotes freehand.
	StrokeKindFreehand StrokeKind = 0
	// StrokeKindStraightLine denotes straight line.
	StrokeKindStraightLine StrokeKind = 1
	// StrokeKindHeadingBox denotes heading box.
	StrokeKindHeadingBox StrokeKind = 2
	// StrokeKindFrame denotes frame.
	StrokeKindFrame StrokeKind = 3
	// StrokeKindStar denotes star.
	StrokeKindStar StrokeKind = 4
	// StrokeKindOther denotes other.
	StrokeKindOther StrokeKind = 5
)

// String names the value and retains unknown numbers in its fallback.
func (v StrokeKind) String() string {
	switch v {
	case StrokeKindFreehand:
		return "freehand"
	case StrokeKindStraightLine:
		return "straight-line"
	case StrokeKindHeadingBox:
		return "heading-box"
	case StrokeKindFrame:
		return "frame"
	case StrokeKindStar:
		return "star"
	case StrokeKindOther:
		return "other"
	default:
		return "stroke_kind(" + strconv.Itoa(int(v)) + ")"
	}
}

// Role is the interpreted role. Unknown numeric values remain intact.
type Role int

const (
	// RoleInk denotes ink.
	RoleInk Role = 0
	// RoleDerived denotes derived.
	RoleDerived Role = 1
	// RoleEraser denotes eraser.
	RoleEraser Role = 2
	// RoleLasso denotes lasso.
	RoleLasso Role = 3
	// RoleOther denotes other.
	RoleOther Role = 4
)

// String names the value and retains unknown numbers in its fallback.
func (v Role) String() string {
	switch v {
	case RoleInk:
		return "ink"
	case RoleDerived:
		return "derived"
	case RoleEraser:
		return "eraser"
	case RoleLasso:
		return "lasso"
	case RoleOther:
		return "other"
	default:
		return "role(" + strconv.Itoa(int(v)) + ")"
	}
}

// StrokeStatus is the interpreted stroke status. Unknown numeric values remain intact.
type StrokeStatus int

const (
	// StrokeStatusLive denotes live.
	StrokeStatusLive StrokeStatus = 0
	// StrokeStatusCleared denotes cleared.
	StrokeStatusCleared StrokeStatus = 1
	// StrokeStatusMoved denotes moved.
	StrokeStatusMoved StrokeStatus = 2
	// StrokeStatusPartiallyErased denotes partially erased.
	StrokeStatusPartiallyErased StrokeStatus = 3
	// StrokeStatusDeleted denotes deleted.
	StrokeStatusDeleted StrokeStatus = 4
	// StrokeStatusErased denotes erased.
	StrokeStatusErased StrokeStatus = 5
	// StrokeStatusOther denotes other.
	StrokeStatusOther StrokeStatus = 6
)

// String names the value and retains unknown numbers in its fallback.
func (v StrokeStatus) String() string {
	switch v {
	case StrokeStatusLive:
		return "live"
	case StrokeStatusCleared:
		return "cleared"
	case StrokeStatusMoved:
		return "moved"
	case StrokeStatusPartiallyErased:
		return "partially-erased"
	case StrokeStatusDeleted:
		return "deleted"
	case StrokeStatusErased:
		return "erased"
	case StrokeStatusOther:
		return "other"
	default:
		return "stroke_status(" + strconv.Itoa(int(v)) + ")"
	}
}

// Category is the interpreted category. Unknown numeric values remain intact.
type Category int

const (
	// CategoryOrdinary denotes ordinary.
	CategoryOrdinary Category = 0
	// CategoryHeading denotes heading.
	CategoryHeading Category = 1
	// CategoryTextBox denotes text box.
	CategoryTextBox Category = 2
	// CategoryLink denotes link.
	CategoryLink Category = 3
	// CategoryPictureLink denotes picture link.
	CategoryPictureLink Category = 4
	// CategoryTextLink denotes text link.
	CategoryTextLink Category = 5
	// CategoryOther denotes other.
	CategoryOther Category = 6
)

// String names the value and retains unknown numbers in its fallback.
func (v Category) String() string {
	switch v {
	case CategoryOrdinary:
		return "ordinary"
	case CategoryHeading:
		return "heading"
	case CategoryTextBox:
		return "text-box"
	case CategoryLink:
		return "link"
	case CategoryPictureLink:
		return "picture-link"
	case CategoryTextLink:
		return "text-link"
	case CategoryOther:
		return "other"
	default:
		return "category(" + strconv.Itoa(int(v)) + ")"
	}
}

// Pen preserves the stored code, including unknown values.
type Pen int32

const (
	// PenInkLegacy is the ink legacy code.
	PenInkLegacy Pen = 1
	// PenEraser is the eraser code.
	PenEraser Pen = 3
	// PenLasso is the lasso code.
	PenLasso Pen = 4
	// PenMarkerLegacy is the marker legacy code.
	PenMarkerLegacy Pen = 5
	// PenEraserSelect is the eraser select code.
	PenEraserSelect Pen = 9
	// PenNeedlePoint is the needle point code.
	PenNeedlePoint Pen = 10
	// PenMarker is the marker code.
	PenMarker Pen = 11
	// PenCalligraphy is the calligraphy code.
	PenCalligraphy Pen = 15
	// PenInk is the ink code.
	PenInk Pen = 16
)

// String names the code, preserving unknown values in decimal.
func (v Pen) String() string {
	switch v {
	case PenInkLegacy:
		return "ink-legacy"
	case PenEraser:
		return "eraser"
	case PenLasso:
		return "lasso"
	case PenMarkerLegacy:
		return "marker-legacy"
	case PenEraserSelect:
		return "eraser-select"
	case PenNeedlePoint:
		return "needle-point"
	case PenMarker:
		return "marker"
	case PenCalligraphy:
		return "calligraphy"
	case PenInk:
		return "ink"
	default:
		return "pen(" + strconv.FormatInt(int64(v), 10) + ")"
	}
}

// InkColor preserves the stored code, including unknown values.
type InkColor uint32

const (
	// ColorBlack is the black code.
	ColorBlack InkColor = 0
	// ColorMarkerBlack is the marker black code.
	ColorMarkerBlack InkColor = 1
	// ColorDarkLegacy is the dark legacy code.
	ColorDarkLegacy InkColor = 48
	// ColorMarkerDarkLegacy is the marker dark legacy code.
	ColorMarkerDarkLegacy InkColor = 49
	// ColorGrayLegacy is the gray legacy code.
	ColorGrayLegacy InkColor = 80
	// ColorMarkerGrayLegacy is the marker gray legacy code.
	ColorMarkerGrayLegacy InkColor = 81
	// ColorDark is the dark code.
	ColorDark InkColor = 157
	// ColorMarkerDark is the marker dark code.
	ColorMarkerDark InkColor = 158
	// ColorGray is the gray code.
	ColorGray InkColor = 201
	// ColorMarkerGray is the marker gray code.
	ColorMarkerGray InkColor = 202
	// ColorWhite is the white code.
	ColorWhite InkColor = 254
	// ColorEraser is the eraser code.
	ColorEraser InkColor = 255
)

// String names the code, preserving unknown values in decimal.
func (v InkColor) String() string {
	switch v {
	case ColorBlack:
		return "black"
	case ColorMarkerBlack:
		return "marker-black"
	case ColorDarkLegacy:
		return "dark-legacy"
	case ColorMarkerDarkLegacy:
		return "marker-dark-legacy"
	case ColorGrayLegacy:
		return "gray-legacy"
	case ColorMarkerGrayLegacy:
		return "marker-gray-legacy"
	case ColorDark:
		return "dark"
	case ColorMarkerDark:
		return "marker-dark"
	case ColorGray:
		return "gray"
	case ColorMarkerGray:
		return "marker-gray"
	case ColorWhite:
		return "white"
	case ColorEraser:
		return "eraser"
	default:
		return "inkcolor(" + strconv.FormatInt(int64(v), 10) + ")"
	}
}
