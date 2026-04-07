package output

import (
	"fmt"
	"time"

	"gitcode.com/gitcode-cli/cli/api"
)

// FormatFlexibleTime renders a FlexibleTime using the requested display style.
func FormatFlexibleTime(t api.FlexibleTime, format TimeFormat) string {
	if t.IsZero() {
		return "unknown"
	}
	return FormatTime(t.Time, format)
}

// FormatTime renders a time using the requested display style.
func FormatTime(t time.Time, format TimeFormat) string {
	if t.IsZero() {
		return "unknown"
	}
	switch format {
	case TimeFormatRelative:
		return formatRelativeTime(t, time.Now())
	default:
		return t.Format("2006-01-02 15:04")
	}
}

func formatRelativeTime(t, now time.Time) string {
	diff := now.Sub(t)
	if diff < 0 {
		diff = -diff
	}

	switch {
	case diff < time.Minute:
		return "just now"
	case diff < time.Hour:
		mins := int(diff / time.Minute)
		if mins == 1 {
			return "1 minute ago"
		}
		return fmt.Sprintf("%d minutes ago", mins)
	case diff < 24*time.Hour:
		hours := int(diff / time.Hour)
		if hours == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%d hours ago", hours)
	case diff < 7*24*time.Hour:
		days := int(diff / (24 * time.Hour))
		if days == 1 {
			return "1 day ago"
		}
		return fmt.Sprintf("%d days ago", days)
	case diff < 30*24*time.Hour:
		weeks := int(diff / (7 * 24 * time.Hour))
		if weeks == 1 {
			return "1 week ago"
		}
		return fmt.Sprintf("%d weeks ago", weeks)
	case diff < 365*24*time.Hour:
		months := int(diff / (30 * 24 * time.Hour))
		if months == 1 {
			return "1 month ago"
		}
		return fmt.Sprintf("%d months ago", months)
	default:
		years := int(diff / (365 * 24 * time.Hour))
		if years == 1 {
			return "1 year ago"
		}
		return fmt.Sprintf("%d years ago", years)
	}
}
