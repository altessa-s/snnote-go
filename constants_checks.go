//go:build snnote_test

// Copyright 2026 ALTESSA SOLUTIONS INC.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file.

package snnote

/*
#include "snnote/snnote.h"
static const char *sngo_version(void) { return SNNOTE_VERSION_STRING; }
static int64_t sngo_EncodingRattaRLE(void) { return SN_ENCODING_RATTA_RLE; }
static int64_t sngo_EncodingPNG(void) { return SN_ENCODING_PNG; }
static int64_t sngo_EncodingJPEG(void) { return SN_ENCODING_JPEG; }
static int64_t sngo_EncodingOther(void) { return SN_ENCODING_OTHER; }
static int64_t sngo_EquipmentOther(void) { return SN_EQUIPMENT_OTHER; }
static int64_t sngo_EquipmentA5X(void) { return SN_EQUIPMENT_A5X; }
static int64_t sngo_EquipmentA6X(void) { return SN_EQUIPMENT_A6X; }
static int64_t sngo_EquipmentNomad(void) { return SN_EQUIPMENT_NOMAD; }
static int64_t sngo_EquipmentManta(void) { return SN_EQUIPMENT_MANTA; }
static int64_t sngo_KindNote(void) { return SN_KIND_NOTE; }
static int64_t sngo_KindMark(void) { return SN_KIND_MARK; }
static int64_t sngo_OrientationPortrait(void) { return SN_ORIENTATION_PORTRAIT; }
static int64_t sngo_OrientationLandscape(void) { return SN_ORIENTATION_LANDSCAPE; }
static int64_t sngo_OrientationOther(void) { return SN_ORIENTATION_OTHER; }
static int64_t sngo_LayerMain(void) { return SN_LAYER_MAIN; }
static int64_t sngo_Layer1(void) { return SN_LAYER_1; }
static int64_t sngo_Layer2(void) { return SN_LAYER_2; }
static int64_t sngo_Layer3(void) { return SN_LAYER_3; }
static int64_t sngo_LayerBackground(void) { return SN_LAYER_BACKGROUND; }
static int64_t sngo_StrokeKindFreehand(void) { return SN_STROKE_KIND_FREEHAND; }
static int64_t sngo_StrokeKindStraightLine(void) { return SN_STROKE_KIND_STRAIGHT_LINE; }
static int64_t sngo_StrokeKindHeadingBox(void) { return SN_STROKE_KIND_HEADING_BOX; }
static int64_t sngo_StrokeKindFrame(void) { return SN_STROKE_KIND_FRAME; }
static int64_t sngo_StrokeKindStar(void) { return SN_STROKE_KIND_STAR; }
static int64_t sngo_StrokeKindOther(void) { return SN_STROKE_KIND_OTHER; }
static int64_t sngo_RoleInk(void) { return SN_ROLE_INK; }
static int64_t sngo_RoleDerived(void) { return SN_ROLE_DERIVED; }
static int64_t sngo_RoleEraser(void) { return SN_ROLE_ERASER; }
static int64_t sngo_RoleLasso(void) { return SN_ROLE_LASSO; }
static int64_t sngo_RoleOther(void) { return SN_ROLE_OTHER; }
static int64_t sngo_StrokeStatusLive(void) { return SN_STROKE_STATUS_LIVE; }
static int64_t sngo_StrokeStatusCleared(void) { return SN_STROKE_STATUS_CLEARED; }
static int64_t sngo_StrokeStatusMoved(void) { return SN_STROKE_STATUS_MOVED; }
static int64_t sngo_StrokeStatusPartiallyErased(void) { return SN_STROKE_STATUS_PARTIALLY_ERASED; }
static int64_t sngo_StrokeStatusDeleted(void) { return SN_STROKE_STATUS_DELETED; }
static int64_t sngo_StrokeStatusErased(void) { return SN_STROKE_STATUS_ERASED; }
static int64_t sngo_StrokeStatusOther(void) { return SN_STROKE_STATUS_OTHER; }
static int64_t sngo_CategoryOrdinary(void) { return SN_CATEGORY_ORDINARY; }
static int64_t sngo_CategoryHeading(void) { return SN_CATEGORY_HEADING; }
static int64_t sngo_CategoryTextBox(void) { return SN_CATEGORY_TEXT_BOX; }
static int64_t sngo_CategoryLink(void) { return SN_CATEGORY_LINK; }
static int64_t sngo_CategoryPictureLink(void) { return SN_CATEGORY_PICTURE_LINK; }
static int64_t sngo_CategoryTextLink(void) { return SN_CATEGORY_TEXT_LINK; }
static int64_t sngo_CategoryOther(void) { return SN_CATEGORY_OTHER; }
static int64_t sngo_PenInkLegacy(void) { return SN_PEN_INK_LEGACY; }
static int64_t sngo_PenEraser(void) { return SN_PEN_ERASER; }
static int64_t sngo_PenLasso(void) { return SN_PEN_LASSO; }
static int64_t sngo_PenMarkerLegacy(void) { return SN_PEN_MARKER_LEGACY; }
static int64_t sngo_PenEraserSelect(void) { return SN_PEN_ERASER_SELECT; }
static int64_t sngo_PenNeedlePoint(void) { return SN_PEN_NEEDLE_POINT; }
static int64_t sngo_PenMarker(void) { return SN_PEN_MARKER; }
static int64_t sngo_PenCalligraphy(void) { return SN_PEN_CALLIGRAPHY; }
static int64_t sngo_PenInk(void) { return SN_PEN_INK; }
static int64_t sngo_ColorBlack(void) { return SN_COLOR_BLACK; }
static int64_t sngo_ColorMarkerBlack(void) { return SN_COLOR_MARKER_BLACK; }
static int64_t sngo_ColorDarkLegacy(void) { return SN_COLOR_DARK_LEGACY; }
static int64_t sngo_ColorMarkerDarkLegacy(void) { return SN_COLOR_MARKER_DARK_LEGACY; }
static int64_t sngo_ColorGrayLegacy(void) { return SN_COLOR_GRAY_LEGACY; }
static int64_t sngo_ColorMarkerGrayLegacy(void) { return SN_COLOR_MARKER_GRAY_LEGACY; }
static int64_t sngo_ColorDark(void) { return SN_COLOR_DARK; }
static int64_t sngo_ColorMarkerDark(void) { return SN_COLOR_MARKER_DARK; }
static int64_t sngo_ColorGray(void) { return SN_COLOR_GRAY; }
static int64_t sngo_ColorMarkerGray(void) { return SN_COLOR_MARKER_GRAY; }
static int64_t sngo_ColorWhite(void) { return SN_COLOR_WHITE; }
static int64_t sngo_ColorEraser(void) { return SN_COLOR_ERASER; }
static int64_t sngo_ParseBudgetFloor(void) { return SN_PARSE_BUDGET_FLOOR; }
static int64_t sngo_ParseBudgetMultiple(void) { return SN_PARSE_BUDGET_MULTIPLE; }
static int64_t sngo_ParseWorkFloor(void) { return SN_PARSE_WORK_FLOOR; }
static int64_t sngo_ParseWorkMultiple(void) { return SN_PARSE_WORK_MULTIPLE; }
static int64_t sngo_ProtocolMaxLength(void) { return SN_PROTOCOL_MAX_LENGTH; }
static int64_t sngo_JSONMaxDepth(void) { return SN_JSON_MAX_DEPTH; }
static int64_t sngo_MaxRasterPixels(void) { return SN_MAX_RASTER_PIXELS; }
static int64_t sngo_StatusOK(void) { return SN_OK; }
static int64_t sngo_StatusInvalidArgument(void) { return SN_ERR_INVALID_ARGUMENT; }
static int64_t sngo_StatusOutOfMemory(void) { return SN_ERR_OUT_OF_MEMORY; }
static int64_t sngo_StatusUnsupportedFormat(void) { return SN_ERR_UNSUPPORTED_FORMAT; }
static int64_t sngo_StatusUnexpectedEOF(void) { return SN_ERR_UNEXPECTED_EOF; }
static int64_t sngo_StatusMalformed(void) { return SN_ERR_MALFORMED; }
static int64_t sngo_StatusImageEncoding(void) { return SN_ERR_IMAGE_ENCODING; }
static int64_t sngo_StatusUnsupportedEncoding(void) { return SN_ERR_UNSUPPORTED_ENCODING; }
static int64_t sngo_StatusLimit(void) { return SN_ERR_LIMIT; }
static int64_t sngo_LogDebug(void) { return SN_LOG_DEBUG; }
static int64_t sngo_LogError(void) { return SN_LOG_ERROR; }
*/
import "C"

func cConstants() (string, map[string][2]int64) {
	return C.GoString(C.sngo_version()), map[string][2]int64{
		"StatusOK":                  {int64(StatusOK), int64(C.sngo_StatusOK())},
		"StatusInvalidArgument":     {int64(StatusInvalidArgument), int64(C.sngo_StatusInvalidArgument())},
		"StatusOutOfMemory":         {int64(StatusOutOfMemory), int64(C.sngo_StatusOutOfMemory())},
		"StatusUnsupportedFormat":   {int64(StatusUnsupportedFormat), int64(C.sngo_StatusUnsupportedFormat())},
		"StatusUnexpectedEOF":       {int64(StatusUnexpectedEOF), int64(C.sngo_StatusUnexpectedEOF())},
		"StatusMalformed":           {int64(StatusMalformed), int64(C.sngo_StatusMalformed())},
		"StatusImageEncoding":       {int64(StatusImageEncoding), int64(C.sngo_StatusImageEncoding())},
		"StatusUnsupportedEncoding": {int64(StatusUnsupportedEncoding), int64(C.sngo_StatusUnsupportedEncoding())},
		"StatusLimit":               {int64(StatusLimit), int64(C.sngo_StatusLimit())},
		"LogDebug":                  {int64(LogDebug), int64(C.sngo_LogDebug())},
		"LogError":                  {int64(LogError), int64(C.sngo_LogError())},

		"EncodingRattaRLE":            {int64(EncodingRattaRLE), int64(C.sngo_EncodingRattaRLE())},
		"EncodingPNG":                 {int64(EncodingPNG), int64(C.sngo_EncodingPNG())},
		"EncodingJPEG":                {int64(EncodingJPEG), int64(C.sngo_EncodingJPEG())},
		"EncodingOther":               {int64(EncodingOther), int64(C.sngo_EncodingOther())},
		"EquipmentOther":              {int64(EquipmentOther), int64(C.sngo_EquipmentOther())},
		"EquipmentA5X":                {int64(EquipmentA5X), int64(C.sngo_EquipmentA5X())},
		"EquipmentA6X":                {int64(EquipmentA6X), int64(C.sngo_EquipmentA6X())},
		"EquipmentNomad":              {int64(EquipmentNomad), int64(C.sngo_EquipmentNomad())},
		"EquipmentManta":              {int64(EquipmentManta), int64(C.sngo_EquipmentManta())},
		"KindNote":                    {int64(KindNote), int64(C.sngo_KindNote())},
		"KindMark":                    {int64(KindMark), int64(C.sngo_KindMark())},
		"OrientationPortrait":         {int64(OrientationPortrait), int64(C.sngo_OrientationPortrait())},
		"OrientationLandscape":        {int64(OrientationLandscape), int64(C.sngo_OrientationLandscape())},
		"OrientationOther":            {int64(OrientationOther), int64(C.sngo_OrientationOther())},
		"LayerMain":                   {int64(LayerMain), int64(C.sngo_LayerMain())},
		"Layer1":                      {int64(Layer1), int64(C.sngo_Layer1())},
		"Layer2":                      {int64(Layer2), int64(C.sngo_Layer2())},
		"Layer3":                      {int64(Layer3), int64(C.sngo_Layer3())},
		"LayerBackground":             {int64(LayerBackground), int64(C.sngo_LayerBackground())},
		"StrokeKindFreehand":          {int64(StrokeKindFreehand), int64(C.sngo_StrokeKindFreehand())},
		"StrokeKindStraightLine":      {int64(StrokeKindStraightLine), int64(C.sngo_StrokeKindStraightLine())},
		"StrokeKindHeadingBox":        {int64(StrokeKindHeadingBox), int64(C.sngo_StrokeKindHeadingBox())},
		"StrokeKindFrame":             {int64(StrokeKindFrame), int64(C.sngo_StrokeKindFrame())},
		"StrokeKindStar":              {int64(StrokeKindStar), int64(C.sngo_StrokeKindStar())},
		"StrokeKindOther":             {int64(StrokeKindOther), int64(C.sngo_StrokeKindOther())},
		"RoleInk":                     {int64(RoleInk), int64(C.sngo_RoleInk())},
		"RoleDerived":                 {int64(RoleDerived), int64(C.sngo_RoleDerived())},
		"RoleEraser":                  {int64(RoleEraser), int64(C.sngo_RoleEraser())},
		"RoleLasso":                   {int64(RoleLasso), int64(C.sngo_RoleLasso())},
		"RoleOther":                   {int64(RoleOther), int64(C.sngo_RoleOther())},
		"StrokeStatusLive":            {int64(StrokeStatusLive), int64(C.sngo_StrokeStatusLive())},
		"StrokeStatusCleared":         {int64(StrokeStatusCleared), int64(C.sngo_StrokeStatusCleared())},
		"StrokeStatusMoved":           {int64(StrokeStatusMoved), int64(C.sngo_StrokeStatusMoved())},
		"StrokeStatusPartiallyErased": {int64(StrokeStatusPartiallyErased), int64(C.sngo_StrokeStatusPartiallyErased())},
		"StrokeStatusDeleted":         {int64(StrokeStatusDeleted), int64(C.sngo_StrokeStatusDeleted())},
		"StrokeStatusErased":          {int64(StrokeStatusErased), int64(C.sngo_StrokeStatusErased())},
		"StrokeStatusOther":           {int64(StrokeStatusOther), int64(C.sngo_StrokeStatusOther())},
		"CategoryOrdinary":            {int64(CategoryOrdinary), int64(C.sngo_CategoryOrdinary())},
		"CategoryHeading":             {int64(CategoryHeading), int64(C.sngo_CategoryHeading())},
		"CategoryTextBox":             {int64(CategoryTextBox), int64(C.sngo_CategoryTextBox())},
		"CategoryLink":                {int64(CategoryLink), int64(C.sngo_CategoryLink())},
		"CategoryPictureLink":         {int64(CategoryPictureLink), int64(C.sngo_CategoryPictureLink())},
		"CategoryTextLink":            {int64(CategoryTextLink), int64(C.sngo_CategoryTextLink())},
		"CategoryOther":               {int64(CategoryOther), int64(C.sngo_CategoryOther())},
		"PenInkLegacy":                {int64(PenInkLegacy), int64(C.sngo_PenInkLegacy())},
		"PenEraser":                   {int64(PenEraser), int64(C.sngo_PenEraser())},
		"PenLasso":                    {int64(PenLasso), int64(C.sngo_PenLasso())},
		"PenMarkerLegacy":             {int64(PenMarkerLegacy), int64(C.sngo_PenMarkerLegacy())},
		"PenEraserSelect":             {int64(PenEraserSelect), int64(C.sngo_PenEraserSelect())},
		"PenNeedlePoint":              {int64(PenNeedlePoint), int64(C.sngo_PenNeedlePoint())},
		"PenMarker":                   {int64(PenMarker), int64(C.sngo_PenMarker())},
		"PenCalligraphy":              {int64(PenCalligraphy), int64(C.sngo_PenCalligraphy())},
		"PenInk":                      {int64(PenInk), int64(C.sngo_PenInk())},
		"ColorBlack":                  {int64(ColorBlack), int64(C.sngo_ColorBlack())},
		"ColorMarkerBlack":            {int64(ColorMarkerBlack), int64(C.sngo_ColorMarkerBlack())},
		"ColorDarkLegacy":             {int64(ColorDarkLegacy), int64(C.sngo_ColorDarkLegacy())},
		"ColorMarkerDarkLegacy":       {int64(ColorMarkerDarkLegacy), int64(C.sngo_ColorMarkerDarkLegacy())},
		"ColorGrayLegacy":             {int64(ColorGrayLegacy), int64(C.sngo_ColorGrayLegacy())},
		"ColorMarkerGrayLegacy":       {int64(ColorMarkerGrayLegacy), int64(C.sngo_ColorMarkerGrayLegacy())},
		"ColorDark":                   {int64(ColorDark), int64(C.sngo_ColorDark())},
		"ColorMarkerDark":             {int64(ColorMarkerDark), int64(C.sngo_ColorMarkerDark())},
		"ColorGray":                   {int64(ColorGray), int64(C.sngo_ColorGray())},
		"ColorMarkerGray":             {int64(ColorMarkerGray), int64(C.sngo_ColorMarkerGray())},
		"ColorWhite":                  {int64(ColorWhite), int64(C.sngo_ColorWhite())},
		"ColorEraser":                 {int64(ColorEraser), int64(C.sngo_ColorEraser())},
		"ParseBudgetFloor":            {int64(ParseBudgetFloor), int64(C.sngo_ParseBudgetFloor())},
		"ParseBudgetMultiple":         {int64(ParseBudgetMultiple), int64(C.sngo_ParseBudgetMultiple())},
		"ParseWorkFloor":              {int64(ParseWorkFloor), int64(C.sngo_ParseWorkFloor())},
		"ParseWorkMultiple":           {int64(ParseWorkMultiple), int64(C.sngo_ParseWorkMultiple())},
		"ProtocolMaxLength":           {int64(ProtocolMaxLength), int64(C.sngo_ProtocolMaxLength())},
		"JSONMaxDepth":                {int64(JSONMaxDepth), int64(C.sngo_JSONMaxDepth())},
		"MaxRasterPixels":             {int64(MaxRasterPixels), int64(C.sngo_MaxRasterPixels())},
	}
}
