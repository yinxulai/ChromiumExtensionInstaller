package browser

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/yinxulai/chromium-extension-installer/internal/utils"
)

// UpdateProfile updates browser profile preferences to add an extension
func UpdateProfile(profile, extensionID, extensionPath string, key []byte, sid string) error {
	prefsPath := filepath.Join(profile, "Preferences")
	securePrefsPath := filepath.Join(profile, "Secure Preferences")

	// Read Preferences as raw map to preserve all existing data
	prefs := make(map[string]interface{})
	if data, err := os.ReadFile(prefsPath); err == nil {
		json.Unmarshal(data, &prefs)
	}

	// Navigate/create extensions structure in prefs
	if prefs["extensions"] == nil {
		prefs["extensions"] = make(map[string]interface{})
	}
	extensions := prefs["extensions"].(map[string]interface{})

	// Handle install_signature
	if extensions["install_signature"] == nil {
		extensions["install_signature"] = make(map[string]interface{})
	}
	installSig := extensions["install_signature"].(map[string]interface{})
	if installSig["ids"] == nil {
		installSig["ids"] = []interface{}{}
	}
	
	// Add extension ID if not exists
	ids := installSig["ids"].([]interface{})
	found := false
	for _, id := range ids {
		if id.(string) == extensionID {
			found = true
			break
		}
	}
	if !found {
		installSig["ids"] = append(ids, extensionID)
	}

	// Handle toolbar
	if extensions["toolbar"] == nil {
		extensions["toolbar"] = []interface{}{}
	}
	toolbar := extensions["toolbar"].([]interface{})
	found = false
	for _, id := range toolbar {
		if id.(string) == extensionID {
			found = true
			break
		}
	}
	if !found {
		extensions["toolbar"] = append(toolbar, extensionID)
	}

	// Read Secure Preferences as raw map to preserve all existing data
	securePrefs := make(map[string]interface{})
	if data, err := os.ReadFile(securePrefsPath); err == nil {
		json.Unmarshal(data, &securePrefs)
	}

	// Navigate/create extensions structure in secure prefs
	if securePrefs["extensions"] == nil {
		securePrefs["extensions"] = make(map[string]interface{})
	}
	secureExtensions := securePrefs["extensions"].(map[string]interface{})
	
	if secureExtensions["settings"] == nil {
		secureExtensions["settings"] = make(map[string]interface{})
	}
	settings := secureExtensions["settings"].(map[string]interface{})

	// Create extension data
	escapedPath := strings.ReplaceAll(extensionPath, "\\", "\\\\")
	extensionData := fmt.Sprintf(`{"active_permissions":{"api":["browsingData","contentSettings","tabs","webRequest","webRequestBlocking"],"explicit_host":["*://*/*","\u003Call_urls>","chrome://favicon/*","http://*/*","https://*/*"],"scriptable_host":["\u003Call_urls>"]},"creation_flags":38,"from_bookmark":false,"from_webstore":false,"granted_permissions":{"api":["browsingData","contentSettings","tabs","webRequest","webRequestBlocking"],"explicit_host":["*://*/*","\u003Call_urls>","chrome://favicon/*","http://*/*","https://*/*"],"scriptable_host":["\u003Call_urls>"]},"install_time":"13188169127141243","location":4,"never_activated_since_loaded":true,"newAllowFileAccess":true,"path":"%s","state":1,"was_installed_by_default":false,"was_installed_by_oem":false}`, escapedPath)

	var extDataMap map[string]interface{}
	json.Unmarshal([]byte(extensionData), &extDataMap)
	settings[extensionID] = extDataMap

	// Calculate HMAC
	message := fmt.Sprintf("%sextensions.settings.%s%s", sid, extensionID, extensionData)
	hash := strings.ToUpper(utils.GetHMACSHA256(key, message))

	// Navigate/create protection structure
	if securePrefs["protection"] == nil {
		securePrefs["protection"] = make(map[string]interface{})
	}
	protection := securePrefs["protection"].(map[string]interface{})
	
	if protection["macs"] == nil {
		protection["macs"] = make(map[string]interface{})
	}
	macs := protection["macs"].(map[string]interface{})
	
	if macs["extensions"] == nil {
		macs["extensions"] = make(map[string]interface{})
	}
	macsExtensions := macs["extensions"].(map[string]interface{})
	
	if macsExtensions["settings"] == nil {
		macsExtensions["settings"] = make(map[string]interface{})
	}
	macsSettings := macsExtensions["settings"].(map[string]interface{})
	
	macsSettings[extensionID] = hash

	// Calculate super_mac
	macsJSON, _ := json.Marshal(macs)
	superMacMessage := fmt.Sprintf("%s%s", sid, string(macsJSON))
	protection["super_mac"] = strings.ToUpper(utils.GetHMACSHA256(key, superMacMessage))

	// Write files
	prefsData, _ := json.MarshalIndent(prefs, "", "  ")
	if err := os.WriteFile(prefsPath, prefsData, 0644); err != nil {
		return err
	}

	securePrefsData, _ := json.MarshalIndent(securePrefs, "", "  ")
	if err := os.WriteFile(securePrefsPath, securePrefsData, 0644); err != nil {
		return err
	}

	return nil
}

// RemoveFromProfile removes extension from browser profile preferences
func RemoveFromProfile(profile, extensionID string, key []byte, sid string) error {
	prefsPath := filepath.Join(profile, "Preferences")
	securePrefsPath := filepath.Join(profile, "Secure Preferences")

	// Read Preferences as raw map to preserve all existing data
	prefs := make(map[string]interface{})
	if data, err := os.ReadFile(prefsPath); err == nil {
		json.Unmarshal(data, &prefs)
	}

	// Remove from Preferences
	if prefs["extensions"] != nil {
		extensions := prefs["extensions"].(map[string]interface{})
		
		if extensions["install_signature"] != nil {
			installSig := extensions["install_signature"].(map[string]interface{})
			if installSig["ids"] != nil {
				ids := installSig["ids"].([]interface{})
				newIds := []interface{}{}
				for _, id := range ids {
					if id.(string) != extensionID {
						newIds = append(newIds, id)
					}
				}
				installSig["ids"] = newIds
			}
		}
		
		if extensions["toolbar"] != nil {
			toolbar := extensions["toolbar"].([]interface{})
			newToolbar := []interface{}{}
			for _, id := range toolbar {
				if id.(string) != extensionID {
					newToolbar = append(newToolbar, id)
				}
			}
			extensions["toolbar"] = newToolbar
		}
	}

	// Read Secure Preferences as raw map to preserve all existing data
	securePrefs := make(map[string]interface{})
	if data, err := os.ReadFile(securePrefsPath); err == nil {
		json.Unmarshal(data, &securePrefs)
	}

	// Remove from Secure Preferences
	if securePrefs["extensions"] != nil {
		extensions := securePrefs["extensions"].(map[string]interface{})
		if extensions["settings"] != nil {
			settings := extensions["settings"].(map[string]interface{})
			delete(settings, extensionID)
		}
	}

	if securePrefs["protection"] != nil {
		protection := securePrefs["protection"].(map[string]interface{})
		if protection["macs"] != nil {
			macs := protection["macs"].(map[string]interface{})
			if macs["extensions"] != nil {
				macsExtensions := macs["extensions"].(map[string]interface{})
				if macsExtensions["settings"] != nil {
					macsSettings := macsExtensions["settings"].(map[string]interface{})
					delete(macsSettings, extensionID)
				}
			}
			
			// Recalculate super_mac
			macsJSON, _ := json.Marshal(macs)
			superMacMessage := fmt.Sprintf("%s%s", sid, string(macsJSON))
			protection["super_mac"] = strings.ToUpper(utils.GetHMACSHA256(key, superMacMessage))
		}
	}

	// Write files
	prefsData, _ := json.MarshalIndent(prefs, "", "  ")
	if err := os.WriteFile(prefsPath, prefsData, 0644); err != nil {
		return err
	}

	securePrefsData, _ := json.MarshalIndent(securePrefs, "", "  ")
	if err := os.WriteFile(securePrefsPath, securePrefsData, 0644); err != nil {
		return err
	}

	return nil
}
