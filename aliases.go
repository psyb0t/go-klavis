package klavis

// The API serves the Google Calendar event schema under different names from
// one response to the next, and earlier releases generated the names below
// from those layouts. scripts/patch-spec.sh now merges the identical copies
// into GoogleCalendarEvent. These aliases keep code written against earlier
// releases compiling.

// Deprecated: use GoogleCalendarEvent.
type SandboxSchemasGoogleCalendarGoogleCalendarEventInput = GoogleCalendarEvent

// Deprecated: use GoogleCalendarEvent.
type SandboxSchemasGoogleCalendarGoogleCalendarEventOutput = GoogleCalendarEvent

// Deprecated: use GoogleCalendarEvent.
type SandboxSchemasGoogleWorkspaceAtlasGoogleCalendarEventInput = GoogleCalendarEvent //nolint:lll // generated name

// Deprecated: use GoogleCalendarEvent.
type SandboxSchemasGoogleWorkspaceAtlasGoogleCalendarEventOutput = GoogleCalendarEvent //nolint:lll // generated name
