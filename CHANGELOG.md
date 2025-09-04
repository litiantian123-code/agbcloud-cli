# Changelog

## [Unreleased]

### Changed
- **BREAKING**: Updated OAuth API endpoint from `/api/oauth/google/login` to `/api/oauth/login_provider`
- Added new required parameters `loginClient` and `oauthProvider` to OAuth login API
  - `loginClient`: Default value "CLI" 
  - `oauthProvider`: Default value "GOOGLE_LOCALHOST"

### Added
- New `GetLoginProviderURL()` method in OAuthAPI interface supporting the updated endpoint
- New response types `OAuthLoginProviderResponse` and `OAuthLoginProviderData`
- Comprehensive test coverage for new API endpoint and parameters
- Backward compatibility support for existing `GetGoogleLoginURL()` method
- API verification demo script in `examples/api_verification_demo.go`

### Maintained
- Full backward compatibility with existing `GetGoogleLoginURL()` method
- All existing functionality preserved while using new endpoint internally
- Legacy response types `OAuthGoogleLoginResponse` and `OAuthGoogleLoginData` still supported

### Technical Details
- The legacy `GetGoogleLoginURL()` method now internally calls `GetLoginProviderURL()` with default parameters
- All tests updated to verify new endpoint and parameters
- Integration tests added for both new and legacy API methods
- Default parameter values ensure seamless migration for existing users

### Migration Guide
- **No action required** for existing code using `GetGoogleLoginURL()`
- To use new functionality, call `GetLoginProviderURL()` with custom `loginClient` and `oauthProvider` values
- New endpoint supports all previous functionality plus additional provider options 