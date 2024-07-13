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

/* Final JSON */

type JSON_F struct {
	Log      string `json:"log,omitempty"`
	Hostname string `json:"host,omitempty"`
}

/* Configuration from YAML */

type Configuration struct {
	Authentication struct {
		Api_Key      string `yaml:"api_key"`
		Company_UUID string `yaml:"company_uuid"`
	}

	Tail struct {
		Tail_File          string `yaml:"tail_file"`
		Waldo_File         string `yaml:"waldo_file"`
		Client_Logging_URL string `yaml:"client_logging_url"`
	}
}
