package types

// Manifest represents the extension manifest.json structure
type Manifest struct {
	Name string `json:"name"`
}

// Preferences represents Chrome's Preferences file structure
type Preferences struct {
	Extensions *ExtensionsPrefs `json:"extensions,omitempty"`
}

// ExtensionsPrefs represents the extensions section in Preferences
type ExtensionsPrefs struct {
	InstallSignature *InstallSignature `json:"install_signature,omitempty"`
	Toolbar          []string          `json:"toolbar,omitempty"`
}

// InstallSignature represents the install signature section
type InstallSignature struct {
	IDs []string `json:"ids,omitempty"`
}

// SecurePreferences represents Chrome's Secure Preferences file structure
type SecurePreferences struct {
	Extensions *SecureExtensions `json:"extensions,omitempty"`
	Protection *Protection       `json:"protection,omitempty"`
}

// SecureExtensions represents the extensions section in Secure Preferences
type SecureExtensions struct {
	Settings map[string]interface{} `json:"settings,omitempty"`
}

// Protection represents the protection section with security signatures
type Protection struct {
	Macs     *Macs  `json:"macs,omitempty"`
	SuperMac string `json:"super_mac,omitempty"`
}

// Macs represents the MAC signatures
type Macs struct {
	Extensions *MacsExtensions `json:"extensions,omitempty"`
}

// MacsExtensions represents the extensions MAC signatures
type MacsExtensions struct {
	Settings map[string]string `json:"settings,omitempty"`
}
