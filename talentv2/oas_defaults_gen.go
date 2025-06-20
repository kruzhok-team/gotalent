// Code generated by ogen, DO NOT EDIT.

package talentv2

// setDefaults set default value of fields.
func (s *EventDiplomaSettings) setDefaults() {
	{
		val := DiplomaTemplate("diploma")
		s.Template = val
	}
	{
		val := DiplomaIssueMode("none")
		s.DiplomaIssueMode = val
	}
}

// setDefaults set default value of fields.
func (s *EventDiplomaSettingsCreateReq) setDefaults() {
	{
		val := bool(false)
		s.DiplomasDarkTheme.SetTo(val)
	}
	{
		val := DiplomaTemplate("diploma")
		s.Template.SetTo(val)
	}
	{
		val := DiplomaIssueMode("none")
		s.DiplomaIssueMode.SetTo(val)
	}
}

// setDefaults set default value of fields.
func (s *EventDiplomaSettingsUpdateReq) setDefaults() {
	{
		val := DiplomaTemplate("diploma")
		s.Template.SetTo(val)
	}
	{
		val := DiplomaIssueMode("none")
		s.DiplomaIssueMode.SetTo(val)
	}
}

// setDefaults set default value of fields.
func (s *FileMetaCreateReq) setDefaults() {
	{
		val := bool(false)
		s.IsPublic.SetTo(val)
	}
}
