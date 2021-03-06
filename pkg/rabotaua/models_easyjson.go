// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package rabotaua

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

func easyjsonD2b7633eDecodeGithubComIskiyRabotauaTelegramBotPkgRabotaua(in *jlexer.Lexer, out *schedulesStruct) {
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
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "id":
			out.ID = int(in.Int())
		case "ua":
			out.Ua = string(in.String())
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
func easyjsonD2b7633eEncodeGithubComIskiyRabotauaTelegramBotPkgRabotaua(out *jwriter.Writer, in schedulesStruct) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix[1:])
		out.Int(int(in.ID))
	}
	{
		const prefix string = ",\"ua\":"
		out.RawString(prefix)
		out.String(string(in.Ua))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v schedulesStruct) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonD2b7633eEncodeGithubComIskiyRabotauaTelegramBotPkgRabotaua(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v schedulesStruct) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD2b7633eEncodeGithubComIskiyRabotauaTelegramBotPkgRabotaua(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *schedulesStruct) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD2b7633eDecodeGithubComIskiyRabotauaTelegramBotPkgRabotaua(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *schedulesStruct) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD2b7633eDecodeGithubComIskiyRabotauaTelegramBotPkgRabotaua(l, v)
}
func easyjsonD2b7633eDecodeGithubComIskiyRabotauaTelegramBotPkgRabotaua1(in *jlexer.Lexer, out *schedules) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		in.Skip()
		*out = nil
	} else {
		in.Delim('[')
		if *out == nil {
			if !in.IsDelim(']') {
				*out = make(schedules, 0, 2)
			} else {
				*out = schedules{}
			}
		} else {
			*out = (*out)[:0]
		}
		for !in.IsDelim(']') {
			var v1 schedulesStruct
			(v1).UnmarshalEasyJSON(in)
			*out = append(*out, v1)
			in.WantComma()
		}
		in.Delim(']')
	}
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonD2b7633eEncodeGithubComIskiyRabotauaTelegramBotPkgRabotaua1(out *jwriter.Writer, in schedules) {
	if in == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
		out.RawString("null")
	} else {
		out.RawByte('[')
		for v2, v3 := range in {
			if v2 > 0 {
				out.RawByte(',')
			}
			(v3).MarshalEasyJSON(out)
		}
		out.RawByte(']')
	}
}

// MarshalJSON supports json.Marshaler interface
func (v schedules) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonD2b7633eEncodeGithubComIskiyRabotauaTelegramBotPkgRabotaua1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v schedules) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD2b7633eEncodeGithubComIskiyRabotauaTelegramBotPkgRabotaua1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *schedules) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD2b7633eDecodeGithubComIskiyRabotauaTelegramBotPkgRabotaua1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *schedules) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD2b7633eDecodeGithubComIskiyRabotauaTelegramBotPkgRabotaua1(l, v)
}
func easyjsonD2b7633eDecodeGithubComIskiyRabotauaTelegramBotPkgRabotaua2(in *jlexer.Lexer, out *VacancyParametersPage) {
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
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "vacancy_parameters":
			(out.VacancyParameters).UnmarshalEasyJSON(in)
		case "page":
			out.Page = int(in.Int())
		case "count":
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
func easyjsonD2b7633eEncodeGithubComIskiyRabotauaTelegramBotPkgRabotaua2(out *jwriter.Writer, in VacancyParametersPage) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"vacancy_parameters\":"
		out.RawString(prefix[1:])
		(in.VacancyParameters).MarshalEasyJSON(out)
	}
	{
		const prefix string = ",\"page\":"
		out.RawString(prefix)
		out.Int(int(in.Page))
	}
	{
		const prefix string = ",\"count\":"
		out.RawString(prefix)
		out.Int(int(in.Count))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v VacancyParametersPage) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonD2b7633eEncodeGithubComIskiyRabotauaTelegramBotPkgRabotaua2(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v VacancyParametersPage) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD2b7633eEncodeGithubComIskiyRabotauaTelegramBotPkgRabotaua2(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *VacancyParametersPage) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD2b7633eDecodeGithubComIskiyRabotauaTelegramBotPkgRabotaua2(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *VacancyParametersPage) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD2b7633eDecodeGithubComIskiyRabotauaTelegramBotPkgRabotaua2(l, v)
}
func easyjsonD2b7633eDecodeGithubComIskiyRabotauaTelegramBotPkgRabotaua3(in *jlexer.Lexer, out *VacancyParameters) {
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
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "keywords":
			out.Keywords = string(in.String())
		case "city_id":
			out.CityID = int(in.Int())
		case "schedule_id":
			out.ScheduleID = int(in.Int())
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
func easyjsonD2b7633eEncodeGithubComIskiyRabotauaTelegramBotPkgRabotaua3(out *jwriter.Writer, in VacancyParameters) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"keywords\":"
		out.RawString(prefix[1:])
		out.String(string(in.Keywords))
	}
	{
		const prefix string = ",\"city_id\":"
		out.RawString(prefix)
		out.Int(int(in.CityID))
	}
	{
		const prefix string = ",\"schedule_id\":"
		out.RawString(prefix)
		out.Int(int(in.ScheduleID))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v VacancyParameters) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonD2b7633eEncodeGithubComIskiyRabotauaTelegramBotPkgRabotaua3(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v VacancyParameters) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD2b7633eEncodeGithubComIskiyRabotauaTelegramBotPkgRabotaua3(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *VacancyParameters) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD2b7633eDecodeGithubComIskiyRabotauaTelegramBotPkgRabotaua3(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *VacancyParameters) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD2b7633eDecodeGithubComIskiyRabotauaTelegramBotPkgRabotaua3(l, v)
}
func easyjsonD2b7633eDecodeGithubComIskiyRabotauaTelegramBotPkgRabotaua4(in *jlexer.Lexer, out *Vacancy) {
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
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "id":
			out.ID = int(in.Int())
		case "name":
			out.Name = string(in.String())
		case "date":
			out.Date = string(in.String())
		case "dateTxt":
			out.DateTxt = string(in.String())
		case "cityName":
			out.CityName = string(in.String())
		case "notebookId":
			out.NotebookID = int(in.Int())
		case "companyName":
			out.CompanyName = string(in.String())
		case "shortDescription":
			out.ShortDescription = string(in.String())
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
func easyjsonD2b7633eEncodeGithubComIskiyRabotauaTelegramBotPkgRabotaua4(out *jwriter.Writer, in Vacancy) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix[1:])
		out.Int(int(in.ID))
	}
	{
		const prefix string = ",\"name\":"
		out.RawString(prefix)
		out.String(string(in.Name))
	}
	{
		const prefix string = ",\"date\":"
		out.RawString(prefix)
		out.String(string(in.Date))
	}
	{
		const prefix string = ",\"dateTxt\":"
		out.RawString(prefix)
		out.String(string(in.DateTxt))
	}
	{
		const prefix string = ",\"cityName\":"
		out.RawString(prefix)
		out.String(string(in.CityName))
	}
	{
		const prefix string = ",\"notebookId\":"
		out.RawString(prefix)
		out.Int(int(in.NotebookID))
	}
	{
		const prefix string = ",\"companyName\":"
		out.RawString(prefix)
		out.String(string(in.CompanyName))
	}
	{
		const prefix string = ",\"shortDescription\":"
		out.RawString(prefix)
		out.String(string(in.ShortDescription))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Vacancy) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonD2b7633eEncodeGithubComIskiyRabotauaTelegramBotPkgRabotaua4(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Vacancy) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD2b7633eEncodeGithubComIskiyRabotauaTelegramBotPkgRabotaua4(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Vacancy) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD2b7633eDecodeGithubComIskiyRabotauaTelegramBotPkgRabotaua4(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Vacancy) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD2b7633eDecodeGithubComIskiyRabotauaTelegramBotPkgRabotaua4(l, v)
}
func easyjsonD2b7633eDecodeGithubComIskiyRabotauaTelegramBotPkgRabotaua5(in *jlexer.Lexer, out *SearchResult) {
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
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "took":
			out.Took = int(in.Int())
		case "start":
			out.Start = int(in.Int())
		case "count":
			out.Count = int(in.Int())
		case "total":
			out.Total = int(in.Int())
		case "errorMessage":
			out.ErrorMessage = string(in.String())
		case "documents":
			if in.IsNull() {
				in.Skip()
				out.Vacancy = nil
			} else {
				in.Delim('[')
				if out.Vacancy == nil {
					if !in.IsDelim(']') {
						out.Vacancy = make([]Vacancy, 0, 0)
					} else {
						out.Vacancy = []Vacancy{}
					}
				} else {
					out.Vacancy = (out.Vacancy)[:0]
				}
				for !in.IsDelim(']') {
					var v4 Vacancy
					(v4).UnmarshalEasyJSON(in)
					out.Vacancy = append(out.Vacancy, v4)
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
func easyjsonD2b7633eEncodeGithubComIskiyRabotauaTelegramBotPkgRabotaua5(out *jwriter.Writer, in SearchResult) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"took\":"
		out.RawString(prefix[1:])
		out.Int(int(in.Took))
	}
	{
		const prefix string = ",\"start\":"
		out.RawString(prefix)
		out.Int(int(in.Start))
	}
	{
		const prefix string = ",\"count\":"
		out.RawString(prefix)
		out.Int(int(in.Count))
	}
	{
		const prefix string = ",\"total\":"
		out.RawString(prefix)
		out.Int(int(in.Total))
	}
	{
		const prefix string = ",\"errorMessage\":"
		out.RawString(prefix)
		out.String(string(in.ErrorMessage))
	}
	{
		const prefix string = ",\"documents\":"
		out.RawString(prefix)
		if in.Vacancy == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v5, v6 := range in.Vacancy {
				if v5 > 0 {
					out.RawByte(',')
				}
				(v6).MarshalEasyJSON(out)
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v SearchResult) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonD2b7633eEncodeGithubComIskiyRabotauaTelegramBotPkgRabotaua5(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v SearchResult) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD2b7633eEncodeGithubComIskiyRabotauaTelegramBotPkgRabotaua5(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *SearchResult) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD2b7633eDecodeGithubComIskiyRabotauaTelegramBotPkgRabotaua5(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *SearchResult) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD2b7633eDecodeGithubComIskiyRabotauaTelegramBotPkgRabotaua5(l, v)
}
func easyjsonD2b7633eDecodeGithubComIskiyRabotauaTelegramBotPkgRabotaua6(in *jlexer.Lexer, out *City) {
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
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "id":
			out.ID = int(in.Int())
		case "ua":
			out.City = string(in.String())
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
func easyjsonD2b7633eEncodeGithubComIskiyRabotauaTelegramBotPkgRabotaua6(out *jwriter.Writer, in City) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix[1:])
		out.Int(int(in.ID))
	}
	{
		const prefix string = ",\"ua\":"
		out.RawString(prefix)
		out.String(string(in.City))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v City) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonD2b7633eEncodeGithubComIskiyRabotauaTelegramBotPkgRabotaua6(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v City) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD2b7633eEncodeGithubComIskiyRabotauaTelegramBotPkgRabotaua6(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *City) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD2b7633eDecodeGithubComIskiyRabotauaTelegramBotPkgRabotaua6(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *City) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD2b7633eDecodeGithubComIskiyRabotauaTelegramBotPkgRabotaua6(l, v)
}
func easyjsonD2b7633eDecodeGithubComIskiyRabotauaTelegramBotPkgRabotaua7(in *jlexer.Lexer, out *Cities) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		in.Skip()
		*out = nil
	} else {
		in.Delim('[')
		if *out == nil {
			if !in.IsDelim(']') {
				*out = make(Cities, 0, 2)
			} else {
				*out = Cities{}
			}
		} else {
			*out = (*out)[:0]
		}
		for !in.IsDelim(']') {
			var v7 City
			(v7).UnmarshalEasyJSON(in)
			*out = append(*out, v7)
			in.WantComma()
		}
		in.Delim(']')
	}
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonD2b7633eEncodeGithubComIskiyRabotauaTelegramBotPkgRabotaua7(out *jwriter.Writer, in Cities) {
	if in == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
		out.RawString("null")
	} else {
		out.RawByte('[')
		for v8, v9 := range in {
			if v8 > 0 {
				out.RawByte(',')
			}
			(v9).MarshalEasyJSON(out)
		}
		out.RawByte(']')
	}
}

// MarshalJSON supports json.Marshaler interface
func (v Cities) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonD2b7633eEncodeGithubComIskiyRabotauaTelegramBotPkgRabotaua7(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Cities) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD2b7633eEncodeGithubComIskiyRabotauaTelegramBotPkgRabotaua7(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Cities) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD2b7633eDecodeGithubComIskiyRabotauaTelegramBotPkgRabotaua7(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Cities) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD2b7633eDecodeGithubComIskiyRabotauaTelegramBotPkgRabotaua7(l, v)
}
