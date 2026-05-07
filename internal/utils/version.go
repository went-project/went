package utils

import (
	"errors"
	"net"
	"strconv"
	"strings"
)

var CurrentVersion = "v1.0.0"

func GetCurrentVersion() string {
	return CurrentVersion
}

func NormalizeVersion(version string) string {
	version = strings.TrimSpace(version)
	version = strings.TrimPrefix(version, "v")
	version = strings.TrimPrefix(version, "V")
	return version
}

func CompareVersions(left, right string) int {
	left = NormalizeVersion(left)
	right = NormalizeVersion(right)

	if left == right {
		return 0
	}

	partsA := strings.Split(left, ".")
	partsB := strings.Split(right, ".")
	maxLen := len(partsA)
	if len(partsB) > maxLen {
		maxLen = len(partsB)
	}

	for i := 0; i < maxLen; i++ {
		a := ""
		b := ""
		if i < len(partsA) {
			a = partsA[i]
		}
		if i < len(partsB) {
			b = partsB[i]
		}

		if a == b {
			continue
		}

		aInt, errA := strconv.Atoi(a)
		bInt, errB := strconv.Atoi(b)

		if errA == nil && errB == nil {
			if aInt < bInt {
				return -1
			}
			if aInt > bInt {
				return 1
			}
			continue
		}

		if errA == nil && errB != nil {
			return 1
		}
		if errA != nil && errB == nil {
			return -1
		}

		return strings.Compare(a, b)
	}

	return strings.Compare(left, right)
}

func IsNetworkError(err error) bool {
	var netErr net.Error
	if errors.As(err, &netErr) {
		return true
	}

	lower := strings.ToLower(err.Error())
	if strings.Contains(lower, "dial tcp") || strings.Contains(lower, "no such host") || strings.Contains(lower, "connection refused") || strings.Contains(lower, "tls handshake timeout") || strings.Contains(lower, "request canceled") || strings.Contains(lower, "failed to connect") {
		return true
	}

	return false
}
