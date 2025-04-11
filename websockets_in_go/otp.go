package main

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type OTP struct {
	Key     string
	Created time.Time
}

// Key value map for OTPS
type RetentionMap map[string]OTP

func (rn *RetentionMap) NewOTP() OTP {
	otp := OTP{
		Key:     uuid.NewString(),
		Created: time.Now(),
	}
	(*rn)[otp.Key] = otp
	return otp
}

func (rn *RetentionMap) VerifyOTP(otp string) bool {
	if _, ok := (*rn)[otp]; !ok {
		return false
	}
	delete(*rn, otp)
	return true
}

func (rn *RetentionMap) Retention(ctx context.Context, retentionPeriod time.Duration) {
	ticker := time.NewTicker(400 * time.Millisecond)
	for {
		select {
		case <-ticker.C:
			for _, otp := range *rn {
				if otp.Created.Add(retentionPeriod).Before(time.Now()) {
					delete(*rn, otp.Key)
				}
			}
		case <-ctx.Done():
			return
		}
	}
}

func NewRetentionMap(ctx context.Context, retentionPeriod time.Duration) RetentionMap {
	rn := make(RetentionMap)
	go rn.Retention(ctx, retentionPeriod)
	return rn
}
