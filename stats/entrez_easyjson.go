// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package stats

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjson93cb6946DecodeGithubComHscellsGrooveStats(in *jlexer.Lexer, out *term) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson93cb6946EncodeGithubComHscellsGrooveStats(out *jwriter.Writer, in term) {
	out.RawByte('{')
	first := true
	_ = first
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v term) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson93cb6946EncodeGithubComHscellsGrooveStats(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v term) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson93cb6946EncodeGithubComHscellsGrooveStats(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *term) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson93cb6946DecodeGithubComHscellsGrooveStats(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *term) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson93cb6946DecodeGithubComHscellsGrooveStats(l, v)
}
func easyjson93cb6946DecodeGithubComHscellsGrooveStats1(in *jlexer.Lexer, out *esearchresult) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "retstart":
			out.RetStart = string(in.String())
		case "count":
			out.Count = string(in.String())
		case "idlist":
			if in.IsNull() {
				in.Skip()
				out.Idlist = nil
			} else {
				in.Delim('[')
				if out.Idlist == nil {
					if !in.IsDelim(']') {
						out.Idlist = make([]string, 0, 4)
					} else {
						out.Idlist = []string{}
					}
				} else {
					out.Idlist = (out.Idlist)[:0]
				}
				for !in.IsDelim(']') {
					var v1 string
					v1 = string(in.String())
					out.Idlist = append(out.Idlist, v1)
					in.WantComma()
				}
				in.Delim(']')
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson93cb6946EncodeGithubComHscellsGrooveStats1(out *jwriter.Writer, in esearchresult) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"retstart\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.RetStart))
	}
	{
		const prefix string = ",\"count\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Count))
	}
	{
		const prefix string = ",\"idlist\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		if in.Idlist == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v2, v3 := range in.Idlist {
				if v2 > 0 {
					out.RawByte(',')
				}
				out.String(string(v3))
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v esearchresult) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson93cb6946EncodeGithubComHscellsGrooveStats1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v esearchresult) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson93cb6946EncodeGithubComHscellsGrooveStats1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *esearchresult) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson93cb6946DecodeGithubComHscellsGrooveStats1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *esearchresult) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson93cb6946DecodeGithubComHscellsGrooveStats1(l, v)
}
func easyjson93cb6946DecodeGithubComHscellsGrooveStats2(in *jlexer.Lexer, out *esearch) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "esearchresult":
			(out.EsearchResult).UnmarshalEasyJSON(in)
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson93cb6946EncodeGithubComHscellsGrooveStats2(out *jwriter.Writer, in esearch) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"esearchresult\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		(in.EsearchResult).MarshalEasyJSON(out)
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v esearch) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson93cb6946EncodeGithubComHscellsGrooveStats2(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v esearch) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson93cb6946EncodeGithubComHscellsGrooveStats2(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *esearch) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson93cb6946DecodeGithubComHscellsGrooveStats2(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *esearch) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson93cb6946DecodeGithubComHscellsGrooveStats2(l, v)
}
func easyjson93cb6946DecodeGithubComHscellsGrooveStats3(in *jlexer.Lexer, out *Search) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "Count":
			out.Count = int(in.Int())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson93cb6946EncodeGithubComHscellsGrooveStats3(out *jwriter.Writer, in Search) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"Count\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int(int(in.Count))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Search) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson93cb6946EncodeGithubComHscellsGrooveStats3(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Search) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson93cb6946EncodeGithubComHscellsGrooveStats3(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Search) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson93cb6946DecodeGithubComHscellsGrooveStats3(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Search) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson93cb6946DecodeGithubComHscellsGrooveStats3(l, v)
}
func easyjson93cb6946DecodeGithubComHscellsGrooveStats4(in *jlexer.Lexer, out *EntrezStatisticsSource) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson93cb6946EncodeGithubComHscellsGrooveStats4(out *jwriter.Writer, in EntrezStatisticsSource) {
	out.RawByte('{')
	first := true
	_ = first
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (e EntrezStatisticsSource) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson93cb6946EncodeGithubComHscellsGrooveStats4(&w, e)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (e EntrezStatisticsSource) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson93cb6946EncodeGithubComHscellsGrooveStats4(w, e)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (e *EntrezStatisticsSource) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson93cb6946DecodeGithubComHscellsGrooveStats4(&r, e)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (e *EntrezStatisticsSource) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson93cb6946DecodeGithubComHscellsGrooveStats4(l, e)
}
func easyjson93cb6946DecodeGithubComHscellsGrooveStats5(in *jlexer.Lexer, out *EntrezDocument) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "ID":
			out.ID = string(in.String())
		case "Title":
			out.Title = string(in.String())
		case "Text":
			out.Text = string(in.String())
		case "MeSHHeadings":
			if in.IsNull() {
				in.Skip()
				out.MeSHHeadings = nil
			} else {
				in.Delim('[')
				if out.MeSHHeadings == nil {
					if !in.IsDelim(']') {
						out.MeSHHeadings = make([]string, 0, 4)
					} else {
						out.MeSHHeadings = []string{}
					}
				} else {
					out.MeSHHeadings = (out.MeSHHeadings)[:0]
				}
				for !in.IsDelim(']') {
					var v4 string
					v4 = string(in.String())
					out.MeSHHeadings = append(out.MeSHHeadings, v4)
					in.WantComma()
				}
				in.Delim(']')
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson93cb6946EncodeGithubComHscellsGrooveStats5(out *jwriter.Writer, in EntrezDocument) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"ID\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.ID))
	}
	{
		const prefix string = ",\"Title\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Title))
	}
	{
		const prefix string = ",\"Text\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Text))
	}
	{
		const prefix string = ",\"MeSHHeadings\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		if in.MeSHHeadings == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v5, v6 := range in.MeSHHeadings {
				if v5 > 0 {
					out.RawByte(',')
				}
				out.String(string(v6))
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v EntrezDocument) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson93cb6946EncodeGithubComHscellsGrooveStats5(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v EntrezDocument) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson93cb6946EncodeGithubComHscellsGrooveStats5(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *EntrezDocument) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson93cb6946DecodeGithubComHscellsGrooveStats5(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *EntrezDocument) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson93cb6946DecodeGithubComHscellsGrooveStats5(l, v)
}