// Copyright 2026 ALTESSA SOLUTIONS INC.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file.

package snnote

/*
#include "snnote/snnote.h"
*/
import "C"

import "unsafe"

func (c *converter) point(v C.sn_point) Point {
	out := Point{}
	out.X = float64(v.x)
	out.Y = float64(v.y)
	return out
}

func (c *converter) size(v C.sn_size) Size {
	out := Size{}
	out.Width = float64(v.width)
	out.Height = float64(v.height)
	return out
}

func (c *converter) rect(v C.sn_rect) *Rect {
	if bool(v.is_null) || !c.charge(1, unsafe.Sizeof(Rect{})) {
		return nil
	}
	out := &Rect{}
	out.X = float64(v.x)
	out.Y = float64(v.y)
	out.Width = float64(v.width)
	out.Height = float64(v.height)
	return out
}

func (c *converter) metadataEntry(v C.sn_metadata_entry) MetadataEntry {
	out := MetadataEntry{}
	out.Key = c.requiredString(v.key, "MetadataEntry.Key")
	out.Value = c.requiredString(v.value, "MetadataEntry.Value")
	return out
}

func (c *converter) metadata(v C.sn_metadata) Metadata {
	out := Metadata{}
	if n := c.allocationLength(unsafe.Pointer(v.entries), v.count, "count", unsafe.Sizeof(out.Entries[0])); n > 0 {
		out.Entries = make([]MetadataEntry, n)
		for i, x := range unsafe.Slice(v.entries, n) {
			if c.err != nil {
				return out
			}
			out.Entries[i] = c.metadataEntry(x)
		}
	}
	return out
}

func (c *converter) bitmap(v C.sn_bitmap) Bitmap {
	out := Bitmap{}
	out.Encoding = Encoding(v.encoding)
	out.Name = c.string(v.name)
	out.Bytes = c.bytes(v.bytes)
	return out
}

func (c *converter) ink(v C.sn_ink) Ink {
	out := Ink{}
	out.Gray = uint8(v.gray)
	out.Alpha = uint8(v.alpha)
	return out
}

func (c *converter) palette(v C.sn_palette) Palette {
	out := Palette{}
	for i := range out.Inks {
		out.Inks[i] = c.ink(v.inks[i])
	}
	return out
}

func (c *converter) layerInfo(v C.sn_layer_info) LayerInfo {
	out := LayerInfo{}
	out.LayerID = int64(v.layer_id)
	out.Name = c.requiredString(v.name, "LayerInfo.Name")
	out.IsBackground = bool(v.is_background)
	out.IsVisible = bool(v.is_visible)
	out.IsDeleted = bool(v.is_deleted)
	out.IsCurrent = bool(v.is_current)
	return out
}

func (c *converter) layer(v C.sn_layer) Layer {
	out := Layer{}
	out.Name = LayerName(v.name)
	if v.bitmap != nil && c.charge(1, unsafe.Sizeof(Bitmap{})) {
		value := c.bitmap(*v.bitmap)
		out.Bitmap = &value
	}
	if v.info != nil && c.charge(1, unsafe.Sizeof(LayerInfo{})) {
		value := c.layerInfo(*v.info)
		out.Info = &value
	}
	out.Metadata = c.metadata(v.metadata)
	out.IsStacked = bool(v.is_stacked)
	out.IsVisible = bool(v.is_visible)
	return out
}

func (c *converter) strokePoint(v C.sn_stroke_point) StrokePoint {
	out := StrokePoint{}
	out.Position = c.point(v.position)
	out.Pressure = uint16(v.pressure)
	return out
}

func (c *converter) contour(v C.sn_contour) Contour {
	out := Contour{}
	if n := c.allocationLength(unsafe.Pointer(v.points), v.count, "count", unsafe.Sizeof(out.Points[0])); n > 0 {
		out.Points = make([]Point, n)
		for i, x := range unsafe.Slice(v.points, n) {
			if c.err != nil {
				return out
			}
			out.Points[i] = c.point(x)
		}
	}
	return out
}

func (c *converter) stroke(v C.sn_stroke) Stroke {
	out := Stroke{}
	out.ID = int32(v.id)
	out.Pen = Pen(v.pen)
	out.RoleRaw = int32(v.role_raw)
	out.StatusRaw = int32(v.status_raw)
	out.Operation = int32(v.operation)
	out.Color = InkColor(v.color)
	out.Thickness = uint32(v.thickness)
	out.Layer = uint32(v.layer)
	out.PageNumber = uint32(v.page_number)
	out.KindName = c.requiredString(v.kind_name, "Stroke.KindName")
	out.Kind = StrokeKind(v.kind)
	out.Role = Role(v.role)
	out.Status = StrokeStatus(v.status)
	out.Category = Category(v.category)
	out.IsDrawn = bool(v.is_drawn)
	if bool(v.has_style) && c.charge(1, unsafe.Sizeof(uint32(0))) {
		value := uint32(v.style)
		out.Style = &value
	}
	if n := c.allocationLength(unsafe.Pointer(v.points), v.point_count, "point_count", unsafe.Sizeof(out.Points[0])); n > 0 {
		out.Points = make([]StrokePoint, n)
		for i, x := range unsafe.Slice(v.points, n) {
			if c.err != nil {
				return out
			}
			out.Points[i] = c.strokePoint(x)
		}
	}
	if n := c.allocationLength(unsafe.Pointer(v.pressures), v.pressure_count, "pressure_count", unsafe.Sizeof(out.Pressures[0])); n > 0 {
		out.Pressures = make([]uint16, n)
		for i, x := range unsafe.Slice(v.pressures, n) {
			if c.err != nil {
				return out
			}
			out.Pressures[i] = uint16(x)
		}
	}
	if n := c.allocationLength(unsafe.Pointer(v.contours), v.contour_count, "contour_count", unsafe.Sizeof(out.Contours[0])); n > 0 {
		out.Contours = make([]Contour, n)
		for i, x := range unsafe.Slice(v.contours, n) {
			if c.err != nil {
				return out
			}
			out.Contours[i] = c.contour(x)
		}
	}
	out.Bounds = c.rect(v.bounds)
	if n := c.allocationLength(unsafe.Pointer(v.attributes), v.attribute_count, "attribute_count", unsafe.Sizeof(out.Attributes[0])); n > 0 {
		out.Attributes = make([]string, n)
		for i, x := range unsafe.Slice(v.attributes, n) {
			if c.err != nil {
				return out
			}
			out.Attributes[i] = c.requiredString(x, "Stroke.Attributes")
		}
	}
	return out
}

func (c *converter) textBox(v C.sn_text_box) TextBox {
	out := TextBox{}
	out.ID = int32(v.id)
	out.Rect = c.rect(v.rect)
	out.Text = c.requiredString(v.text, "TextBox.Text")
	out.FontPath = c.string(v.font_path)
	out.Timestamp = c.string(v.timestamp)
	if bool(v.has_font_size) && c.charge(1, unsafe.Sizeof(float64(0))) {
		value := float64(v.font_size)
		out.FontSize = &value
	}
	if n := c.allocationLength(unsafe.Pointer(v.attributes), v.attribute_count, "attribute_count", unsafe.Sizeof(out.Attributes[0])); n > 0 {
		out.Attributes = make([]string, n)
		for i, x := range unsafe.Slice(v.attributes, n) {
			if c.err != nil {
				return out
			}
			out.Attributes[i] = c.requiredString(x, "TextBox.Attributes")
		}
	}
	return out
}

func (c *converter) word(v C.sn_word) Word {
	out := Word{}
	out.Label = c.requiredString(v.label, "Word.Label")
	out.Box = c.rect(v.box)
	return out
}

func (c *converter) element(v C.sn_element) Element {
	out := Element{}
	out.Kind = c.requiredString(v.kind, "Element.Kind")
	out.Label = c.string(v.label)
	if n := c.allocationLength(unsafe.Pointer(v.words), v.word_count, "word_count", unsafe.Sizeof(out.Words[0])); n > 0 {
		out.Words = make([]Word, n)
		for i, x := range unsafe.Slice(v.words, n) {
			if c.err != nil {
				return out
			}
			out.Words[i] = c.word(x)
		}
	}
	return out
}

func (c *converter) recognition(v C.sn_recognition) Recognition {
	out := Recognition{}
	if n := c.allocationLength(unsafe.Pointer(v.elements), v.element_count, "element_count", unsafe.Sizeof(out.Elements[0])); n > 0 {
		out.Elements = make([]Element, n)
		for i, x := range unsafe.Slice(v.elements, n) {
			if c.err != nil {
				return out
			}
			out.Elements[i] = c.element(x)
		}
	}
	out.Text = c.requiredString(v.text, "Recognition.Text")
	return out
}

func (c *converter) page(v C.sn_page) Page {
	out := Page{}
	out.Number = int64(v.number)
	out.ID = c.string(v.id)
	out.TemplateName = c.string(v.template_name)
	out.Orientation = Orientation(v.orientation)
	out.OrientationRaw = int64(v.orientation_raw)
	out.Size = c.size(v.size)
	if n := c.allocationLength(unsafe.Pointer(v.layers), v.layer_count, "layer_count", unsafe.Sizeof(out.Layers[0])); n > 0 {
		out.Layers = make([]Layer, n)
		for i, x := range unsafe.Slice(v.layers, n) {
			if c.err != nil {
				return out
			}
			out.Layers[i] = c.layer(x)
		}
	}
	if n := c.allocationLength(unsafe.Pointer(v.layer_info), v.layer_info_count, "layer_info_count", unsafe.Sizeof(out.LayerInfo[0])); n > 0 {
		out.LayerInfo = make([]LayerInfo, n)
		for i, x := range unsafe.Slice(v.layer_info, n) {
			if c.err != nil {
				return out
			}
			out.LayerInfo[i] = c.layerInfo(x)
		}
	}
	if n := c.allocationLength(unsafe.Pointer(v.strokes), v.stroke_count, "stroke_count", unsafe.Sizeof(out.Strokes[0])); n > 0 {
		out.Strokes = make([]Stroke, n)
		for i, x := range unsafe.Slice(v.strokes, n) {
			if c.err != nil {
				return out
			}
			out.Strokes[i] = c.stroke(x)
		}
	}
	out.SkippedStrokes = c.length(v.skipped_strokes, "skipped_strokes")
	if n := c.allocationLength(unsafe.Pointer(v.drawn_strokes), v.drawn_stroke_count, "drawn_stroke_count", unsafe.Sizeof(out.DrawnStrokes[0])); n > 0 {
		out.DrawnStrokes = make([]*Stroke, n)
		for i, x := range unsafe.Slice(v.drawn_strokes, n) {
			if c.err != nil {
				return out
			}
			base, address := uintptr(unsafe.Pointer(v.strokes)), uintptr(unsafe.Pointer(x))
			stride := unsafe.Sizeof(*x)
			if x == nil || address < base || (address-base)%stride != 0 || (address-base)/stride >= uintptr(len(out.Strokes)) {
				c.err = &Error{Status: StatusMalformed, Message: "drawn stroke outside strokes"}
				return out
			}
			index := (address - base) / stride
			out.DrawnStrokes[i] = &out.Strokes[index]
		}
	}
	if n := c.allocationLength(unsafe.Pointer(v.text_boxes), v.text_box_count, "text_box_count", unsafe.Sizeof(out.TextBoxes[0])); n > 0 {
		out.TextBoxes = make([]TextBox, n)
		for i, x := range unsafe.Slice(v.text_boxes, n) {
			if c.err != nil {
				return out
			}
			out.TextBoxes[i] = c.textBox(x)
		}
	}
	if v.recognition != nil && c.charge(1, unsafe.Sizeof(Recognition{})) {
		value := c.recognition(*v.recognition)
		out.Recognition = &value
	}
	out.RecognitionArchive = c.bytes(v.recognition_archive)
	if n := c.allocationLength(unsafe.Pointer(v.disabled_areas), v.disabled_area_count, "disabled_area_count", unsafe.Sizeof(out.DisabledAreas[0])); n > 0 {
		out.DisabledAreas = make([]*Rect, n)
		for i, x := range unsafe.Slice(v.disabled_areas, n) {
			if c.err != nil {
				return out
			}
			out.DisabledAreas[i] = c.rect(x)
		}
	}
	out.Metadata = c.metadata(v.metadata)
	return out
}

func (c *converter) title(v C.sn_title) Title {
	out := Title{}
	out.PageNumber = int64(v.page_number)
	out.Level = int64(v.level)
	out.Style = TitleStyle(v.style)
	out.Rect = c.rect(v.rect)
	out.OriginalRect = c.rect(v.original_rect)
	if v.bitmap != nil && c.charge(1, unsafe.Sizeof(Bitmap{})) {
		value := c.bitmap(*v.bitmap)
		out.Bitmap = &value
	}
	out.Metadata = c.metadata(v.metadata)
	return out
}

func (c *converter) keyword(v C.sn_keyword) Keyword {
	out := Keyword{}
	out.PageNumber = int64(v.page_number)
	out.Rect = c.rect(v.rect)
	out.OriginalRect = c.rect(v.original_rect)
	out.Text = c.requiredString(v.text, "Keyword.Text")
	out.Metadata = c.metadata(v.metadata)
	return out
}

func (c *converter) link(v C.sn_link) Link {
	out := Link{}
	out.PageNumber = int64(v.page_number)
	out.Kind = int64(v.kind)
	out.Direction = int64(v.direction)
	if bool(v.has_target_page_number) && c.charge(1, unsafe.Sizeof(int64(0))) {
		value := int64(v.target_page_number)
		out.TargetPageNumber = &value
	}
	out.Style = int64(v.style)
	out.Rect = c.rect(v.rect)
	out.TargetPath = c.string(v.target_path)
	out.TargetFileID = c.string(v.target_file_id)
	out.TargetPageID = c.string(v.target_page_id)
	out.Timestamp = c.string(v.timestamp)
	out.Text = c.string(v.text)
	out.URL = c.string(v.url)
	if v.bitmap != nil && c.charge(1, unsafe.Sizeof(Bitmap{})) {
		value := c.bitmap(*v.bitmap)
		out.Bitmap = &value
	}
	out.Metadata = c.metadata(v.metadata)
	return out
}

func (c *converter) highlight(v C.sn_highlight) Highlight {
	out := Highlight{}
	out.PageNumber = int64(v.page_number)
	out.Kind = int64(v.kind)
	out.ColorType = int64(v.color_type)
	out.Comment = c.string(v.comment)
	if n := c.allocationLength(unsafe.Pointer(v.rects), v.rect_count, "rect_count", unsafe.Sizeof(out.Rects[0])); n > 0 {
		out.Rects = make([]*Rect, n)
		for i, x := range unsafe.Slice(v.rects, n) {
			if c.err != nil {
				return out
			}
			out.Rects[i] = c.rect(x)
		}
	}
	if n := c.allocationLength(unsafe.Pointer(v.document_rects), v.document_rect_count, "document_rect_count", unsafe.Sizeof(out.DocumentRects[0])); n > 0 {
		out.DocumentRects = make([]*Rect, n)
		for i, x := range unsafe.Slice(v.document_rects, n) {
			if c.err != nil {
				return out
			}
			out.DocumentRects[i] = c.rect(x)
		}
	}
	if bool(v.has_date) && c.charge(1, unsafe.Sizeof(float64(0))) {
		value := float64(v.date)
		out.Date = &value
	}
	return out
}

func (c *converter) template(v C.sn_template) Template {
	out := Template{}
	out.Name = c.requiredString(v.name, "Template.Name")
	out.Bitmap = c.bitmap(v.bitmap)
	return out
}

func (c *converter) notebook(v C.sn_notebook) Notebook {
	out := Notebook{}
	out.Palette = c.palette(v.palette)
	out.Kind = Kind(v.kind)
	out.FormatVersion = int64(v.format_version)
	out.Equipment = Equipment(v.equipment)
	out.EquipmentRaw = c.requiredString(v.equipment_raw, "Notebook.EquipmentRaw")
	out.FileID = c.string(v.file_id)
	out.Sheet = c.size(v.sheet)
	if n := c.allocationLength(unsafe.Pointer(v.pages), v.page_count, "page_count", unsafe.Sizeof(out.Pages[0])); n > 0 {
		out.Pages = make([]Page, n)
		for i, x := range unsafe.Slice(v.pages, n) {
			if c.err != nil {
				return out
			}
			out.Pages[i] = c.page(x)
		}
	}
	out.SkippedPages = c.length(v.skipped_pages, "skipped_pages")
	if n := c.allocationLength(unsafe.Pointer(v.titles), v.title_count, "title_count", unsafe.Sizeof(out.Titles[0])); n > 0 {
		out.Titles = make([]Title, n)
		for i, x := range unsafe.Slice(v.titles, n) {
			if c.err != nil {
				return out
			}
			out.Titles[i] = c.title(x)
		}
	}
	if n := c.allocationLength(unsafe.Pointer(v.keywords), v.keyword_count, "keyword_count", unsafe.Sizeof(out.Keywords[0])); n > 0 {
		out.Keywords = make([]Keyword, n)
		for i, x := range unsafe.Slice(v.keywords, n) {
			if c.err != nil {
				return out
			}
			out.Keywords[i] = c.keyword(x)
		}
	}
	if n := c.allocationLength(unsafe.Pointer(v.links), v.link_count, "link_count", unsafe.Sizeof(out.Links[0])); n > 0 {
		out.Links = make([]Link, n)
		for i, x := range unsafe.Slice(v.links, n) {
			if c.err != nil {
				return out
			}
			out.Links[i] = c.link(x)
		}
	}
	if n := c.allocationLength(unsafe.Pointer(v.highlights), v.highlight_count, "highlight_count", unsafe.Sizeof(out.Highlights[0])); n > 0 {
		out.Highlights = make([]Highlight, n)
		for i, x := range unsafe.Slice(v.highlights, n) {
			if c.err != nil {
				return out
			}
			out.Highlights[i] = c.highlight(x)
		}
	}
	if n := c.allocationLength(unsafe.Pointer(v.templates), v.template_count, "template_count", unsafe.Sizeof(out.Templates[0])); n > 0 {
		out.Templates = make([]Template, n)
		for i, x := range unsafe.Slice(v.templates, n) {
			if c.err != nil {
				return out
			}
			out.Templates[i] = c.template(x)
		}
	}
	if v.cover != nil && c.charge(1, unsafe.Sizeof(Bitmap{})) {
		value := c.bitmap(*v.cover)
		out.Cover = &value
	}
	out.CoverType = int64(v.cover_type)
	out.Header = c.metadata(v.header)
	out.Footer = c.metadata(v.footer)
	return out
}
