package db

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

type contextKey string

const (
	sqlContextKey       = contextKey("pgx_sql")
	argsContextKey      = contextKey("pgx_args")
	startTimeContextKey = contextKey("pgx_start_time")
)

// QueryLogLevel controls when queries are logged based on their duration
type QueryLogLevel struct {
	Debug time.Duration
	Info  time.Duration
	Warn  time.Duration
}

// PgxZapTracer provides human-friendly SQL query logging
type PgxZapTracer struct {
	Logger    *zap.Logger
	Enabled   bool
	LogLevel  QueryLogLevel
	MaxArgLen int
	Colorized bool
}

// NewPgxZapTracer creates a new tracer with sensible defaults
func NewPgxZapTracer(logger *zap.Logger, enabled bool) *PgxZapTracer {
	defaultMaxArgLen := 100
	if defaultMaxArgLen <= 3 {
		logger.Warn("MaxArgLen too small, resetting to default 100")
		defaultMaxArgLen = 100
	}

	return &PgxZapTracer{
		Logger:    logger,
		Enabled:   enabled,
		MaxArgLen: defaultMaxArgLen,
		LogLevel: QueryLogLevel{
			Debug: 100 * time.Millisecond,
			Info:  500 * time.Millisecond,
			Warn:  2 * time.Second,
		},
		Colorized: true,
	}
}

func (t *PgxZapTracer) TraceQueryStart(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	if !t.Enabled {
		return ctx
	}

	ctx = context.WithValue(ctx, sqlContextKey, formatSQL(data.SQL))
	ctx = context.WithValue(ctx, argsContextKey, data.Args)
	ctx = context.WithValue(ctx, startTimeContextKey, time.Now())

	return ctx
}

func (t *PgxZapTracer) TraceQueryEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryEndData) {
	if !t.Enabled {
		return
	}

	sql, _ := ctx.Value(sqlContextKey).(string)
	args, _ := ctx.Value(argsContextKey).([]any)
	startTime, _ := ctx.Value(startTimeContextKey).(time.Time)

	duration := time.Since(startTime)
	finalSQL := interpolateSQL(sql, args, t.MaxArgLen)

	if t.Colorized {
		finalSQL = colorizeSQL(finalSQL)
	}

	fields := []zap.Field{
		zap.String("query", finalSQL),
		zap.Duration("duration", duration),
		zap.String("duration_human", humanizeDuration(duration)),
	}

	if data.Err == nil {
		fields = append(fields, zap.String("result", data.CommandTag.String()))
		if data.CommandTag.RowsAffected() > 0 {
			fields = append(fields, zap.Int64("rows_affected", data.CommandTag.RowsAffected()))
		}
	}

	if data.Err != nil {
		fields = append(fields, zap.Error(data.Err))
		t.Logger.Error("SQL query failed", fields...)
	} else {
		switch {
		case duration <= t.LogLevel.Debug:
			t.Logger.Debug("SQL query completed", fields...)
		case duration <= t.LogLevel.Info:
			t.Logger.Info("SQL query completed", fields...)
		case duration <= t.LogLevel.Warn:
			t.Logger.Warn("SQL query slow", fields...)
		default:
			t.Logger.Info("SQL query completed (slow)", fields...)
		}
	}
}

// --- Helper functions ---

func formatSQL(rawSQL string) string {
	lines := strings.Split(rawSQL, "\n")
	var cleaned []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "--") || line == "" {
			continue
		}
		cleaned = append(cleaned, line)
	}

	return strings.Join(cleaned, " ")
}

func interpolateSQL(sql string, args []any, maxLen int) string {
	if maxLen < 4 {
		maxLen = 100
	}

	result := sql

	for i, arg := range args {
		placeholder := fmt.Sprintf("$%d", i+1)

		var value string
		switch v := arg.(type) {
		case nil:
			value = "NULL"
		case string:
			value = fmt.Sprintf("'%s'", escapeString(truncateString(v, maxLen)))
		case []byte:
			value = fmt.Sprintf("'%s'", escapeString(truncateString(string(v), maxLen)))
		case time.Time:
			value = fmt.Sprintf("'%s'", v.Format(time.RFC3339))
		default:
			value = fmt.Sprintf("%v", v)
		}

		result = strings.Replace(result, placeholder, value, 1)
	}

	return result
}

func escapeString(s string) string {
	return strings.ReplaceAll(s, "'", "''")
}

func truncateString(s string, maxLen int) string {
	if maxLen <= 3 {
		if maxLen <= 0 {
			return ""
		}
		return strings.Repeat(".", maxLen)
	}

	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

func humanizeDuration(d time.Duration) string {
	if d < time.Microsecond {
		return fmt.Sprintf("%d ns", d.Nanoseconds())
	}
	if d < time.Millisecond {
		return fmt.Sprintf("%.2f µs", float64(d.Nanoseconds())/float64(time.Microsecond))
	}
	if d < time.Second {
		return fmt.Sprintf("%.2f ms", float64(d.Nanoseconds())/float64(time.Millisecond))
	}
	return fmt.Sprintf("%.2f s", float64(d.Nanoseconds())/float64(time.Second))
}

// ANSI color codes
const (
	colorReset   = "\033[0m"
	colorRed     = "\033[31m"
	colorGreen   = "\033[32m"
	colorYellow  = "\033[33m"
	colorBlue    = "\033[34m"
	colorMagenta = "\033[35m"
	colorCyan    = "\033[36m"
)

func colorizeSQL(sql string) string {
	keywords := []string{
		"SELECT", "INSERT", "UPDATE", "DELETE", "FROM", "WHERE",
		"JOIN", "LEFT JOIN", "RIGHT JOIN", "INNER JOIN", "OUTER JOIN",
		"GROUP BY", "ORDER BY", "HAVING", "LIMIT", "OFFSET",
		"AND", "OR", "NOT", "IN", "BETWEEN", "LIKE", "IS NULL", "IS NOT NULL",
		"CREATE", "ALTER", "DROP", "TABLE", "INDEX", "VIEW", "FUNCTION",
		"BEGIN", "COMMIT", "ROLLBACK", "TRANSACTION",
	}

	result := sql

	for _, keyword := range keywords {
		upperKeyword := strings.ToUpper(keyword)
		replacement := colorBlue + upperKeyword + colorReset
		result = strings.ReplaceAll(result, upperKeyword, replacement)
		result = strings.ReplaceAll(result, strings.ToLower(keyword), replacement)
	}

	inString := false
	var highlighted strings.Builder

	for i := 0; i < len(result); i++ {
		if result[i] == '\'' {
			if inString {
				highlighted.WriteString("'" + colorReset)
				inString = false
			} else {
				highlighted.WriteString(colorGreen + "'")
				inString = true
			}
		} else {
			highlighted.WriteByte(result[i])
		}
	}

	return highlighted.String()
}
