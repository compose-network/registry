// Package chainlist provides read-only access to the embedded
// chain list for Compose networks. It mirrors the public schema
// (including RPC and Explorers arrays) and adds a derived Slug.
//
// Primary entry points:
//   - List: returns all entries
//   - Get: returns an entry by slug (suffix of identifier)
//   - GetByIdentifier: returns an entry by full identifier (e.g. "hoodi/rollup-a")
//   - ListByNetwork: returns entries by parent network (e.g. "hoodi")
//   - Version: synthetic version based on entry count
package chainlist
