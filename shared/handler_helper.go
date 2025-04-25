package shared

import (
	"fmt"
	"log"
	"time"
)

func Retry(maxRetries int, delay time.Duration, operation func() error) error {
	var err error
	for i := 0; i < maxRetries; i++ {
		err = operation()
		if err == nil {
			return nil
		}

		log.Printf("Retry %d: %v\n", i+1, err)
		time.Sleep(delay * time.Duration(i+1))
	}
	return fmt.Errorf("retry failed after %d attempts: %w", maxRetries, err)
}

func RetryBasedOnVersionConflict(maxRetries int, delay time.Duration, operation func() error) error {
	var err error
	for i := 0; i < maxRetries; i++ {
		err = operation()
		if err == nil {
			return nil
		}

		if IsVersionConflict(err) {
			log.Printf("RetryBasedOnVersionConflict %d due to version conflict...\n", i+1)
			time.Sleep(delay * time.Duration(i+1))
			continue
		}

		return err
	}
	return fmt.Errorf("retry failed after %d attempts: %w", maxRetries, err)
}

func RetryHandlerBasedOnVersionConflict(operation func() error) error {
	return RetryBasedOnVersionConflict(3, 200*time.Millisecond, operation)
}
