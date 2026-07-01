package utils

import (
	"errors"
	"net"
	"strconv"
	"strings"
)

var CurrentVersion = "v1.0.0"

type Channel string

const (
	ChannelStable Channel = "stable"
	ChannelBeta   Channel = "beta"
)

func GetCurrentVersion() string {
	return CurrentVersion
}

func NormalizeVersion(version string) string {
	version = strings.TrimSpace(version)
	version = strings.TrimPrefix(version, "v")
	version = strings.TrimPrefix(version, "V")
	return version
}

func ParseChannelFlag(beta bool) Channel {
	if beta {
		return ChannelBeta
	}
	return ChannelStable
}

func CompareVersions(left, right string) int {
	left = NormalizeVersion(left)
	right = NormalizeVersion(right)

	if left == right {
		return 0
	}

	leftCore, leftPre := splitPrerelease(left)
	rightCore, rightPre := splitPrerelease(right)

	if cmp := compareVersionCore(leftCore, rightCore); cmp != 0 {
		return cmp
	}

	if leftPre == "" && rightPre != "" {
		return 1
	}
	if leftPre != "" && rightPre == "" {
		return -1
	}

	return comparePrerelease(leftPre, rightPre)
}

func splitPrerelease(version string) (string, string) {
	parts := strings.SplitN(version, "-", 2)
	if len(parts) == 1 {
		return parts[0], ""
	}
	return parts[0], parts[1]
}

func compareVersionCore(a, b string) int {
	aParts := strings.Split(a, ".")
	bParts := strings.Split(b, ".")
	maxLen := len(aParts)
	if len(bParts) > maxLen {
		maxLen = len(bParts)
	}

	for i := 0; i < maxLen; i++ {
		aPart := ""
		bPart := ""
		if i < len(aParts) {
			aPart = aParts[i]
		}
		if i < len(bParts) {
			bPart = bParts[i]
		}

		if aPart == bPart {
			continue
		}

		aInt, errA := strconv.Atoi(aPart)
		bInt, errB := strconv.Atoi(bPart)

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

		return strings.Compare(aPart, bPart)
	}

	return 0
}

func comparePrerelease(a, b string) int {
	if a == b {
		return 0
	}

	aParts := strings.Split(a, ".")
	bParts := strings.Split(b, ".")
	maxLen := len(aParts)
	if len(bParts) > maxLen {
		maxLen = len(bParts)
	}

	for i := 0; i < maxLen; i++ {
		aPart := ""
		bPart := ""
		if i < len(aParts) {
			aPart = aParts[i]
		}
		if i < len(bParts) {
			bPart = bParts[i]
		}

		if aPart == bPart {
			continue
		}

		aInt, errA := strconv.Atoi(aPart)
		bInt, errB := strconv.Atoi(bPart)

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

		return strings.Compare(aPart, bPart)
	}

	if len(aParts) < len(bParts) {
		return -1
	}
	if len(aParts) > len(bParts) {
		return 1
	}

	return strings.Compare(a, b)
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
