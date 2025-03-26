package common

import (
	"io"
)

func validateAction(action string, condition bool, err error, id string) error {
	if condition {
		log.Errorf(" %s | result: fail | client_id: %v | error: %v",
			action, id, err)
		return err
	}
	return nil
}

func validateSend(action string, err error, id string) (int, error) {
	if err != nil {
		log.Errorf(" %s | result: fail | client_id: %v | error: %v",
			action, id, err)
		return 0, err
	}
	return 0, nil
}

func validateRecv(action string, err error, id string) (string, error) {
	if err != nil {
		if err == io.EOF {
			log.Infof(" %s | result: EOF | client_id: %v", action, id)
			return "", nil
		}
		log.Infof(" %s | result: fail | client_id: %v | error: %v", action, id, err)
		return "",err
	}
	return "",nil
}

func validateScanner(action string, err error, id string) (string, error) {
	if err != nil {
		log.Errorf("%s | result: fail | error: %v", action, err)
		return "", err
	}
	return "", nil
}
