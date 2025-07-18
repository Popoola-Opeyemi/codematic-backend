package utils

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"regexp"
	"time"

	"crypto/sha256"
	"encoding/hex"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/oapi-codegen/runtime/types"
	"github.com/shopspring/decimal"
)

func ToDBString(s *string) pgtype.Text {
	if s == nil {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{String: *s, Valid: true}
}

func ToDBDate(d *types.Date) pgtype.Timestamp {
	if d == nil {
		return pgtype.Timestamp{Valid: false}
	}
	return pgtype.Timestamp{Time: d.Time, Valid: true}
}

func ToPgxTstamp(t time.Time) pgtype.Timestamp {
	return pgtype.Timestamp{Time: t, Valid: true}
}

func ToPgxText(s string) pgtype.Text {
	return pgtype.Text{String: s, Valid: true}
}

func ToPgxBool(b bool) pgtype.Bool {
	return pgtype.Bool{Bool: b, Valid: true}
}

func ToPgxInt2(val int16) pgtype.Int2 {
	return pgtype.Int2{Int16: val, Valid: true}
}

// OpenAPIToGo converts OpenAPI JSON data to a Go struct
func OpenAPIToGo(data []byte, target interface{}) error {
	// First unmarshal into a map to handle OpenAPI specific formats
	var rawData map[string]interface{}
	if err := json.Unmarshal(data, &rawData); err != nil {
		return fmt.Errorf("failed to unmarshal OpenAPI data: %w", err)
	}

	// Convert OpenAPI specific formats
	convertedData := convertOpenAPIFormats(rawData)

	// Marshal back to JSON with converted formats
	jsonData, err := json.Marshal(convertedData)
	if err != nil {
		return fmt.Errorf("failed to marshal converted data: %w", err)
	}

	// Unmarshal into the target Go struct
	if err := json.Unmarshal(jsonData, target); err != nil {
		return fmt.Errorf("failed to unmarshal into Go struct: %w", err)
	}

	return nil
}

// GoToOpenAPI converts a Go struct to OpenAPI JSON format
func GoToOpenAPI(source interface{}) ([]byte, error) {
	// First marshal the Go struct to JSON
	jsonData, err := json.Marshal(source)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Go struct: %w", err)
	}

	// Unmarshal into a map to handle format conversions
	var rawData map[string]interface{}
	if err := json.Unmarshal(jsonData, &rawData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal Go data: %w", err)
	}

	// Convert to OpenAPI formats
	convertedData := convertGoFormats(rawData)

	// Marshal back to JSON with OpenAPI formats
	result, err := json.Marshal(convertedData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal OpenAPI data: %w", err)
	}

	return result, nil
}

// convertOpenAPIFormats converts OpenAPI format values to Go format
func convertOpenAPIFormats(data map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	for key, value := range data {
		switch v := value.(type) {
		case string:
			// Handle UUID format
			if isUUID(v) {
				if uuid, err := uuid.Parse(v); err == nil {
					result[key] = uuid
					continue
				}
			}
			// Handle date-time format
			if t, err := time.Parse(time.RFC3339, v); err == nil {
				result[key] = t
				continue
			}
			// Handle date format
			if t, err := time.Parse("2006-01-02", v); err == nil {
				result[key] = t
				continue
			}
			// Handle decimal format
			if d, err := decimal.NewFromString(v); err == nil {
				result[key] = d
				continue
			}
			result[key] = v
		case map[string]interface{}:
			result[key] = convertOpenAPIFormats(v)
		case []interface{}:
			converted := make([]interface{}, len(v))
			for i, item := range v {
				if m, ok := item.(map[string]interface{}); ok {
					converted[i] = convertOpenAPIFormats(m)
				} else {
					converted[i] = item
				}
			}
			result[key] = converted
		default:
			result[key] = v
		}
	}

	return result
}

// convertGoFormats converts Go format values to OpenAPI format
func convertGoFormats(data map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	for key, value := range data {
		switch v := value.(type) {
		case uuid.UUID:
			result[key] = v.String()
		case time.Time:
			result[key] = v.Format(time.RFC3339)
		case decimal.Decimal:
			result[key] = v.String()
		case map[string]interface{}:
			result[key] = convertGoFormats(v)
		case []interface{}:
			converted := make([]interface{}, len(v))
			for i, item := range v {
				if m, ok := item.(map[string]interface{}); ok {
					converted[i] = convertGoFormats(m)
				} else {
					converted[i] = item
				}
			}
			result[key] = converted
		default:
			result[key] = v
		}
	}

	return result
}

// isUUID checks if a string is in UUID format
func isUUID(s string) bool {
	_, err := uuid.Parse(s)
	return err == nil
}

// UUID conversion functions
func ToUUID(id uuid.UUID) pgtype.UUID {
	return pgtype.UUID{Bytes: [16]byte(id), Valid: true}
}

func FromPgUUID(id pgtype.UUID) uuid.UUID {
	return uuid.UUID(id.Bytes)
}

func ParseUUID(id string) (uuid.UUID, error) {
	return uuid.Parse(id)
}

// StringToPgUUID converts a string to pgtype.UUID
func StringToPgUUID(id string) (pgtype.UUID, error) {
	u, err := uuid.Parse(id)
	if err != nil {
		return pgtype.UUID{}, err
	}
	return pgtype.UUID{Bytes: [16]byte(u), Valid: true}, nil
}

// Timestamp conversion functions
func ToPgTimestamp(t *time.Time) pgtype.Timestamp {
	if t == nil {
		return pgtype.Timestamp{Valid: false}
	}
	return pgtype.Timestamp{Time: *t, Valid: true}
}

func ToPgTimestamptz(t time.Time) pgtype.Timestamptz {
	return pgtype.Timestamptz{Time: t, Valid: true}
}

func FromPgTimestamp(t pgtype.Timestamp) time.Time {
	if !t.Valid {
		return time.Time{}
	}
	return t.Time
}

// Decimal conversion functions
func ParseDecimal(s string) (decimal.Decimal, error) {
	return decimal.NewFromString(s)
}

func FormatDecimal(d decimal.Decimal) string {
	return d.StringFixedBank(2)
}

// String conversion functions
func ParseString(b [16]byte) string {
	u, err := uuid.FromBytes(b[:])
	if err != nil {
		return ""
	}
	return u.String()
}

func FromPgText(t pgtype.Text) *string {
	if !t.Valid {
		return nil
	}
	return &t.String
}

// Convert row types to Customer type

// Customer conversion function

func NullFloat64(f *float64) sql.NullFloat64 {
	if f == nil {
		return sql.NullFloat64{}
	}
	return sql.NullFloat64{Float64: *f, Valid: true}
}

func NullTimePtrString(s string) sql.NullTime {
	if s == "" {
		return sql.NullTime{}
	}
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return sql.NullTime{}
	}
	return sql.NullTime{Time: t, Valid: true}
}

func NullInt64Ptr(i *int) sql.NullInt64 {
	if i == nil {
		return sql.NullInt64{}
	}
	return sql.NullInt64{Int64: int64(*i), Valid: true}
}

func StringOrEmpty(ptr *string) string {
	if ptr != nil {
		return *ptr
	}
	return ""
}

// HashString returns the SHA256 hash of the input string as a hex string
func HashString(s string) string {
	h := sha256.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

// isValidEmail uses a simple regex for email validation
func IsValidEmail(email string) bool {
	// Simple RFC 5322 regex
	var re = regexp.MustCompile(`^[a-zA-Z0-9._%%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}
