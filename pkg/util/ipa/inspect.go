package ipa

import (
	"archive/zip"
	"fmt"
	"io"
	"strings"

	"howett.net/plist"
)

func PrintURLSchemesFromIPA(ipaPath string) error {
	r, err := zip.OpenReader(ipaPath)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		if strings.HasPrefix(f.Name, "Payload/") &&
			strings.HasSuffix(f.Name, ".app/Info.plist") {

			rc, err := f.Open()
			if err != nil {
				return err
			}
			defer rc.Close()

			data, err := io.ReadAll(rc)
			if err != nil {
				return err
			}

			var plistData map[string]interface{}
			_, err = plist.Unmarshal(data, &plistData)
			if err != nil {
				return fmt.Errorf("failed to parse Info.plist: %w", err)
			}

			urlTypes, ok := plistData["CFBundleURLTypes"]
			if !ok {
				fmt.Println("CFBundleURLSchemes: (none)")
				return nil
			}

			fmt.Println("CFBundleURLSchemes:")

			for _, item := range urlTypes.([]interface{}) {
				dict := item.(map[string]interface{})

				schemes, ok := dict["CFBundleURLSchemes"]
				if !ok {
					continue
				}

				for _, scheme := range schemes.([]interface{}) {
					fmt.Println(scheme.(string))
				}
			}

			return nil
		}
	}

	return fmt.Errorf("Info.plist not found in IPA")
}