## [1.0.4] - 2026-11-05
### Added
- sandboxing configuration

### Changed
- Updated copyright license and file headers
- frontend: Updates for Yarn update in the toolchain
- frontend: Set Yarn security option enableScripts=false which disables preinstall, install and postinstall scripts
- Normalize GDR values uint64 map keys to decimal strings so OBIS lookups remain stable across runtimes
  Note: This addresses behaviour observed with protobufjs 7.5.9+ where `map<uint64, ...>` keys can be represented as 8-byte hash strings when Long support is active

## [1.0.3] - 2026-07-03
### Changed
- update vite to 6.4.3 (CVE fix)

## [1.0.2] - 2026-06-15

## [1.0.1] - 2026-06-04
### Fixed
- setup of serial modbus connection
- missing language keys
