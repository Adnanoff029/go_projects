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

func NewRetentionMap(ctx context.Context, retentionPeriod time.Duration) RetentionMap {
	rn := make(RetentionMap)
	return rn
}
