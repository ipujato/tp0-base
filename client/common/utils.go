package common

import (
	"io"
	"runtime/debug"
)

func validateAction(action string, condition bool, err error, id string) error {
	if condition {
		log.Errorf("action: %s | result: fail | client_id: %v | error: %v",
			action, id, err)
		return err
	}
	return nil
}

func validateSend(action string, err error, id string) (int, error) {
	if err != nil {
		log.Errorf("action: %s | result: fail | client_id: %v | error: %v",
			action, id, err)
		return 0, err
	}
	return 0, nil
}

func validateRecv(action string, err error, id string) (string, error) {
	if err != nil && err != io.EOF{
		log.Errorf("action: %s | result: fail | client_id: %v | error: %v",
			action, id, err)
		log.Errorf("stack trace: %v", debug.Stack())
		return "", err
	}
	return "", nil
}