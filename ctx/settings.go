package ctx

const (
	SettingsActionMenuText = "Выберите действие: "
	SettingsTimeStartText  = "Введите время начала оповещений в интервале от 0 до 24: "
	SettingsTimeEndText    = "Введите время окончания оповещений в интервале от 0 до 24: "
	SettingsTimeZoneText   = "Введите временную зону в интервале от -12 до 14: "
)

const (
	SettingsTimeStartUpdate = 0
	SettingsTimeEndUpdate   = 1
	SettingsTimeZoneUpdate  = 2
)

type SettingsNotificationInfoStruct struct {
	IsOn      bool
	TimeStart int
	TimeEnd   int
	TimeZone  int
}
