# ---------------------------------------------------------------------------
# berlios-create-fileserver-admin.ps1 - Create file server admin group and add users

# Copyright 2023 The MITRE Corporation. Approved for public release. Document number CT0005.
# Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in
# compliance with the License. You may obtain a copy of the License at
# http://www.apache.org/licenses/LICENSE-2.0
# Unless required by applicable law or agreed to in writing, software distributed under the License
# is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
# or implied. See the License for the specific language governing permissions and limitations under the License.

# This project makes use of ATT&CK®
# ATT&CK Terms of Use - https://attack.mitre.org/resources/terms-of-use/

# Usage: berlios-create-fileserver-admin.ps1

# ---------------------------------------------------------------------------

# Script to add necessarry roles for File Server Admin in Snake scenario

Import-Module ActiveDirectory

$fileServerAdmin = "EgleAdmin"
$groupName = "File Server Admins"

New-ADGroup -Name $groupName -GroupCategory Security -GroupScope Global

Add-ADGroupMember -Identity $groupName -Members $fileServerAdmin
Add-ADGroupMember -Identity "Schema Admins" -Members "evals_domain_admin"
