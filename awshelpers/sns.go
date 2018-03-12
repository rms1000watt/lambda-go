package awshelpers

import "encoding/json"

func ConvertStringToAlarmMessage(message string) (response *AlarmMessage) {
	response = &AlarmMessage{}

	if err := json.Unmarshal([]byte(message), response); err != nil {
	}

	return
}
