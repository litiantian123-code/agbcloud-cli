# Changelog

## [Unreleased]

### Changed
- **BREAKING**: Updated OAuth API endpoint from `/api/oauth/google/login` to `/api/oauth/login_provider`
- **BREAKING**: Removed legacy `GetGoogleLoginURL()` method from OAuthAPI interface
- Added new required parameters `loginClient` and `oauthProvider` to OAuth login API
  - `loginClient`: Default value "CLI" 
  - `oauthProvider`: Default value "GOOGLE_LOCALHOST"

### Added
- New `GetLoginProviderURL()` method in OAuthAPI interface supporting the updated endpoint
- New response types `OAuthLoginProviderResponse` and `OAuthLoginProviderData`
- Comprehensive test coverage for new API endpoint and parameters
- API verification demo script in `examples/api_verification_demo.go`

### Removed
- **BREAKING**: Legacy `GetGoogleLoginURL()` method and associated response types
- Legacy response types `OAuthGoogleLoginResponse` and `OAuthGoogleLoginData`

### Technical Details
- All OAuth requests now use the `/api/oauth/login_provider` endpoint
- All tests updated to use new API method and endpoint
- Integration tests include comprehensive coverage of new endpoint with various parameter combinations
- Login command updated to use new API directly

### Migration Guide
- **Action required**: Replace all calls to `GetGoogleLoginURL()` with `GetLoginProviderURL()`
- Update method signature: `GetLoginProviderURL(ctx, fromUrlPath, "CLI", "GOOGLE_LOCALHOST")`
- Update response type references from `OAuthGoogleLoginResponse` to `OAuthLoginProviderResponse`
- Update data type references from `OAuthGoogleLoginData` to `OAuthLoginProviderData` 