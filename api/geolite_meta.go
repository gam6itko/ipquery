package api

import (
	"time"
)

// DatabaseInfo describes an mmdb file, as read from its metadata section.
type DatabaseInfo struct {
	Type          string    `json:"type"`
	Description   string    `json:"description"`
	BuildTime     time.Time `json:"build_time"`
	IPVersion     uint      `json:"ip_version"`
	Languages     []string  `json:"languages"`
	NodeCount     uint      `json:"node_count"`
	RecordSize    uint      `json:"record_size"`
	BinaryVersion string    `json:"binary_format_version"`
}

// Info returns the metadata of the currently active reader. A closed database
// reports a zero value.
func (d *ReloadableGeoIPDB) Info() DatabaseInfo {
	d.mu.RLock()
	defer d.mu.RUnlock()

	if d.db == nil {
		return DatabaseInfo{}
	}

	m := d.db.Metadata

	// Not every mmdb carries an English description; fall back to the type,
	// which is always present.
	desc := m.Description["en"]
	if desc == "" {
		desc = m.DatabaseType
	}

	langs := m.Languages
	if langs == nil {
		langs = []string{}
	}

	return DatabaseInfo{
		Type:          m.DatabaseType,
		Description:   desc,
		BuildTime:     m.BuildTime().UTC(),
		IPVersion:     m.IPVersion,
		Languages:     langs,
		NodeCount:     m.NodeCount,
		RecordSize:    m.RecordSize,
		BinaryVersion: fmtVersion(m.BinaryFormatMajorVersion, m.BinaryFormatMinorVersion),
	}
}

// MetaResult is the payload of GET /meta.
type MetaResult struct {
	Databases map[string]DatabaseInfo `json:"databases"`
}
