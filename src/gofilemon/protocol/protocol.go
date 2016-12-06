package protocol

import (
	"encoding/json"
)

const (
	PING = `ping`
	PONG = `pong`

	ERROR_COMMON_CODE = 10

	ERROR_ONLY_DAEMON_START     = "ONLY_DAEMON_START_PLEEZE_RETRY"
	ERROR_ONLY_REGISTER         = "ONLY_REGISTER_PLEEZE_RETRY"
	ERROR_COMMAND_TOO_SMALL     = "COMMAND_TOO_SMALL_ERROR"
	ERROR_CREATE_FILE_PROCESSOR = "CREATE_FILE_PROCESSOR_ERROR"
	ERROR_COMMAND               = "COMMAND_ERROR"
)

func StringSliceToByte(data []string) []byte {
	result, _ := json.Marshal(data)
	return result
}

func ByteToStringSlice(data []byte) ([]string, error) {
	var result []string
	err := json.Unmarshal(data, &result)
	return result, err
}
