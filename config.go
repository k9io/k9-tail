/*
 * Copyright (C) 2024 Key9 Identity, Inc <k9.io>
 * Copyright (C) 2024 Champ Clark III <cclark@k9.io>
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License Version 2 as
 * published by the Free Software Foundation.  You may not use, modify or
 * distribute this program under any other version of the GNU General
 * Public License.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, write to the Free Software
 * Foundation, Inc., 59 Temple Place - Suite 330, Boston, MA 02111-1307, USA.
 */

package main

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

var Config *Configuration

func LoadConfig(config_file string) *Configuration {

	/* Load config file */

	file, err := os.Open(config_file)

	if err != nil {
		log.Fatalf("Cannot open '%s' YAML file.", config_file)
	}

	defer file.Close()

	/* Init new YAML decode */

	d := yaml.NewDecoder(file)

	err = d.Decode(&Config)

	if err != nil {
		log.Fatalf("Cannot parse '%s'.", config_file)
	}

	/* Sanity Checks */

	if Config.Authentication.Api_Key == "" {
		log.Fatalf("'api_key' key not found in %s\n", config_file)
	}

	if Config.Authentication.Company_UUID == "" {
		log.Fatalf("'company_uuid' key not found in %s\n", config_file)
	}

	if Config.Tail.Tail_File == "" {
		log.Fatalf("'tail_file' key not found in %s\n", config_file)
	}

	if Config.Tail.Waldo_File == "" {
		log.Fatalf("'waldo_file' key not found in %s\n", config_file)
	}

	if Config.Tail.Client_Logging_URL == "" {
		log.Fatalf("'waldo_file' key not found in %s\n", config_file)
	}

	return Config

}
