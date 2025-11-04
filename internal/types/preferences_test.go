package types

import (
"encoding/json"
"testing"
)

func TestManifestMarshalUnmarshal(t *testing.T) {
	manifest := Manifest{Name: "Test Extension"}
	data, err := json.Marshal(manifest)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}

	var result Manifest
	err = json.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	if result.Name != manifest.Name {
		t.Errorf("Name = %s, want %s", result.Name, manifest.Name)
	}
}

func TestPreferencesMarshalUnmarshal(t *testing.T) {
	prefs := Preferences{
		Extensions: &ExtensionsPrefs{
			InstallSignature: &InstallSignature{
				IDs: []string{"ext1", "ext2", "ext3"},
			},
			Toolbar: []string{"ext1", "ext2"},
		},
	}

	data, err := json.Marshal(prefs)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}

	var result Preferences
	err = json.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	if result.Extensions == nil {
		t.Fatal("Extensions is nil")
	}
	if len(result.Extensions.InstallSignature.IDs) != 3 {
		t.Errorf("IDs length = %d, want 3", len(result.Extensions.InstallSignature.IDs))
	}
}

func TestSecurePreferencesMarshalUnmarshal(t *testing.T) {
	securePrefs := SecurePreferences{
		Extensions: &SecureExtensions{
			Settings: map[string]interface{}{
				"ext1": map[string]interface{}{"path": "C:\\ext"},
			},
		},
		Protection: &Protection{
			Macs: &Macs{
				Extensions: &MacsExtensions{
					Settings: map[string]string{"ext1": "ABC123"},
				},
			},
			SuperMac: "DEF456",
		},
	}

	data, err := json.Marshal(securePrefs)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}

	var result SecurePreferences
	err = json.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	if result.Protection.SuperMac != "DEF456" {
		t.Errorf("SuperMac = %s, want DEF456", result.Protection.SuperMac)
	}
}
