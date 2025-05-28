[![MITRE ATT&CK® v11](https://img.shields.io/badge/MITRE%20ATT%26CK®-v11-red)](https://attack.mitre.org/versions/v11/)


# Micro Emulation Plans

This project aims to expand the impact of the [Adversary Emulation
Library](https://github.com/center-for-threat-informed-defense/adversary_emulation_library)
by developing easy-to-execute emulation content that targets specific behaviors
and challenges facing defenders. Rather than modeling full-scope breach
scenarios, these micro emulation plans focus on smaller, more specific
combinations of techniques. (We call these short sequences "compound behaviors"
because they involve more than one adversary technique.) The plans produce
signals similar to real adversary techniques (but without the harmful effects)
to enable blue teams to safely test and evaluate their detective and
preventative controls.

These micro emulation plans are optimized for ease-of-use (most are
click-to-run) but are also flexible enough to integrate with breach and attack
simulation (BAS) tools. (See the [CALDERA example](/caldera-integration/).)

**Table Of Contents:**

- [Available Plans](#available-plans)
- [What is a Micro Emulation Plan?](#what-is-a-micro-emulation-plan)
- [User Guidance](#user-guidance)
  - [Getting Started](#getting-started)
  - [Making These Plans Your Own](#making-these-plans-your-own)
- [How Do I Contribute?](#how-do-i-contribute)
- [Notice](#notice)

## Available Plans

|                  Micro Emulation Plan                   | Description                                                                                                                                                                                                                                                                                                                                                                                                                                                           |
| :-----------------------------------------------------: | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [File Access and File Modification](./src/file_access/) | Emulates file access and modification behaviors commonly associated with [TA0009 Collection](https://attack.mitre.org/tactics/TA0009/) as well as [T1486 Data Encrypted for Impact](https://attack.mitre.org/techniques/T1486/). (Data source: [DS0022 File](https://attack.mitre.org/datasources/DS0022/))                                                                                                                                                           |
|            [Named Pipes](./src/named_pipes/)            | Emulates the creation and use of named pipes [commonly abused by malware](https://labs.withsecure.com/blog/detecting-cobalt-strike-default-modules-via-named-pipe-analysis/) (Data source: [DS0023 Named Pipe](https://attack.mitre.org/datasources/DS0023/))                                                                                                                                                                                                         |
|       [Windows Registry](./src/windows_registry/)       | Emulates a few common methods that adversaries use to modify the Windows Registry. (Data Source: [DS0024 Windows Registry](https://attack.mitre.org/datasources/DS0024/))                                                                                                                                                                                                                                                                                             |
|              [Web Shells](./src/webshell/)              | Emulates the compound behavior of planting a web shell ([T1505.003 Server Software Component: Web Shell](https://attack.mitre.org/techniques/T1505/003/)) and then executing arbitrary commands through it ([T1059 Command and Scripting Interpreter](https://attack.mitre.org/techniques/T1059/))                                                                                                                                                                    |
|      [Process Injection](./src/process_injection/)      | Emulates the compound behavior of [T1055 Process Injection](https://attack.mitre.org/techniques/T1055/) followed by execution of arbitrary commands ([T1059 Command and Scripting Interpreter](https://attack.mitre.org/techniques/T1059/))                                                                                                                                                                                                                           |
|         [User Execution](./src/user_execution/)         | Emulates the compound behavior of delivering a malicious `.one`, `.doc`, `.lnk`, or `.iso` file (e.g. via [T1566.001 Phishing: Spearphishing Attachment](https://attack.mitre.org/techniques/T1566/001/)) and then executing arbitrary commands after a user invokes the file ([T1204.002 User Execution: Malicious File](https://attack.mitre.org/techniques/T1204/002/) and [T1059 Command and Scripting Interpreter](https://attack.mitre.org/techniques/T1059/)). |
|     [Active Directory Enumeration](./src/ad_enum/)      | Emulate multiple [TA0007 Discovery](https://attack.mitre.org/tactics/TA0007/) behaviors through commonly abused interfaces and services such as Active Directory (AD)                                                                                                                                                                                                                                                                                                 |
|     [Reflective Loading](./src/reflective_loading/)     | Emulates an adversary performing running malicious code within an arbitrary process to perform [T1620 Reflective Code Loading](https://attack.mitre.org/techniques/T1620/)                                                                                                                                                                                                                                                                                            |
|       [Remote Code Execution](./src/apache_rce/)        | Emulates an adversary performing remote code execution against a vulnerable web server as documented in [T1190 Exploit Public-Facing Application](https://attack.mitre.org/techniques/T1190/)                                                                                                                                                                                                                                                                         |
|           [Log Clearing](./src/log_clearing)            | Emulates an adversary clearing Windows Event Log, as described in [T1070.001 Indicator Removal: Clear Windows Event Logs](https://attack.mitre.org/techniques/T1070/001/)                                                                                                                                                                                                                                                                                             |
|         [Data Exfiltration](./src/data_exfil/)          | Emulates the compound behaviors of an adversary finding, staging, archiving, and extracting sensitive files, as described in [TA0010 Exfiltration](https://attack.mitre.org/tactics/TA0010/)                                                                                                                                                                                                                                                                          |
|        [DLL Sideloading](./src/dll_sideloading/)        | Emulates an adversary executing an otherwise legitimate/benign application in order to hijack its modules/libraries to instead inject their malicious payload, as described in [T1574.002 Hijack Execution Flow: DLL Side-Loading](https://attack.mitre.org/techniques/T1574/002/)                                                                                                                                                                                    |

## What is a Micro Emulation Plan?

**Abstraction matters.** We are HUGE fans of using threat-informed adversary
emulation to inspire, prioritize, and direct security testing and tuning. We
celebrate the community adoption and investment in this approach, but also
recognize that not all organizations have the resources to build and execute
large emulation plans.

Community projects such as [Atomic Red
Team™](https://github.com/redcanaryco/atomic-red-team) have **greatly** reduced
the cost of entry, but we recognize an opportunity for a middle ground between
atomic tests and full-scope emulation plans. There is a need for adversary
emulation plans that focus on evaluating collection of specific data sources or
security controls for short sequences of commonly observed adversary behaviors.
These micro emulation plans answer questions such as:

> Are we collecting the right data about [data source]?
>
> How well do we detect or prevent [a particular part of an intrusion]?

We hope that these plans offer both actionable resources as well as inspiration
for future research and micro emulation plan development.

## User Guidance

### Getting Started

Code and documentation for each micro emulation plan are available in
[`/src`](./src). The `README` file for each plan is structured the same and
includes:

- Description of Emulated Behaviors (What are we doing?)
- CTI / Background (Why you should care?)
- Execution Instructions / Resources (How you can use this plan)
- Execution Demo (What successful execution looks like)
- Defensive Lessons Learned (Detection and mitigation guidance)

### Making These Plans Your Own

All micro emulation plans were designed to be simple to execute. Hard-coded
parameters for each module enable double-click default execution, which can also
be tailored via provided command-line arguments or interactive menus provided by
the module.

For each plan, a compiled and ready to execute binary (within a `.zip` archive)
as well as source code and build/customization instructions are available.

We encourage users to review the `README` documentation for each plan then
consider opportunities for testing and related improvements within their
organization.

## How Do I Contribute?

We welcome your feedback and contributions to help advance **Micro Emulation**.
Please see the guidance for contributors if are you interested in [contributing
or simply reporting issues.](/CONTRIBUTING.md)

Please submit
[issues](https://github.com/center-for-threat-informed-defense/micro-emulation-plans/issues)
for any technical questions/concerns or contact ctid@mitre.org
directly for more general inquiries.

## Notice

Copyright 2023 The MITRE Corporation.

Licensed under the Apache License, Version 2.0 (the "License"); you may not use
this file except in compliance with the License. You may obtain a copy of the
License at:

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed
under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR
CONDITIONS OF ANY KIND, either express or implied. See the License for the
specific language governing permissions and limitations under the License.

This project makes use of ATT&CK®

[ATT&CK Terms of Use](https://attack.mitre.org/resources/terms-of-use/)
