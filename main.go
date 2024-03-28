package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/choria-io/fisk"
	"github.com/ppmathis/cloudns-go"
)

var (
	authId   int
	authPass string
	target   string
)

func main() {
	app := fisk.New("cloudns-backup", "Backs up ClouDNS.net hosted zones")

	backup := app.Commandf("backup", "Performs a Backup").Action(backupAction)
	backup.Flag("target", "The directory to write backups to").Envar("CLOUDNS_TARGET").Required().StringVar(&target)
	backup.Flag("auth-id", "Authentication ID").Envar("CLOUDNS_AUTH_ID").Required().IntVar(&authId)
	backup.Flag("password", "Authentication Password").Envar("CLOUDNS_AUTH_PASSWORD").Required().StringVar(&authPass)

	app.MustParseWithUsage(os.Args[1:])
}

func backupAction(_ *fisk.ParseContext) error {
	if !filepath.IsAbs(target) {
		return fmt.Errorf("target must be absolute path")
	}

	client, err := cloudns.New(cloudns.AuthUserID(authId, authPass))
	if err != nil {
		return err
	}

	zones, err := client.Zones.List(context.Background())
	if err != nil {
		return err
	}

	outDir, err := filepath.Abs(target)
	if err != nil {
		return err
	}

	err = os.MkdirAll(outDir, 0700)
	if err != nil {
		return err
	}

	for _, zone := range zones {
		if zone.Type != cloudns.ZoneTypeMaster {
			log.Printf("Skipping non master zone %s", zone.Name)
			continue
		}

		zoneFile := filepath.Join(outDir, zone.Name)
		log.Printf("Fetching export for %s info %s", zone.Name, zoneFile)

		export, err := client.Records.Export(context.Background(), zone.Name)
		if err != nil {
			return err
		}

		if export.Status != "Success" {
			return fmt.Errorf("fetching zone %s failed: %v", zone.Name, export.StatusDescription)
		}

		err = os.WriteFile(zoneFile, []byte(export.Zone), 0644)
		if err != nil {
			return err
		}
	}

	return nil
}
