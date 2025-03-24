package common

import (
	"io"
	"runtime"
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

		// Obtener el backtrace manualmente
		stackBuf := make([]byte, 1024*8) // 8 KB de buffer para la pila
		n := runtime.Stack(stackBuf, false)
		stackTrace := string(stackBuf[:n])

		log.Infof(" %s | result: fail | client_id: %v | error: %v", action, id, err)
		log.Infof("stack trace: %s", stackTrace) // Captura más precisa del stack trace

		return "",err
	}
	return "",nil
}