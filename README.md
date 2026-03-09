![alt text](/assets/icons/icon.png "Logo")

# SmartDiary
SmartDiary is an app runnig on Anroid, Linux and Windows. SmartDiary is a little diary application.  
For every day you can add a text entry.  
For highlighting entries you can choose between colors and symbols, every color can be labeled with a text.  
You can set a password and protect entries so that for viewing / modifying the password must be given.  
You can search in the diary for expressions. You can export the whole diary to JSON format.  
It is also possible to import an previously exported json file.  
SmartDiary is written in [Go](https://go.dev/) and uses [Fyne](https://fyne.io/) as graphical toolkit.

Author: Reiner Pröls  
Licence: MIT  

## Usage of Smartdiary
Tap on a date or entering data.

## Screenshots
![alt text](/screenshots/login.jpg "Login screen")
![alt text](/screenshots/categoryview.jpg "Category view")
![alt text](/screenshots/entryview.jpg "Entry view")
![alt text](/screenshots/entryedit.jpg "Edit view")
![alt text](/screenshots/fieldedit.jpg "Field edit")
![alt text](/screenshots/iconselect.jpg "Icon select")
![alt text](/screenshots/searchresults.jpg "Search results")
![alt text](/screenshots/changepassword.jpg "Change password")
![alt text](/screenshots/settings.jpg "Settings view")

### Precompiled binaries
#### Linux (64 Bit)
[Tar file](https://github.com/bytemystery-com/smartdiary/releases/download/v0.3.0/SmartDiary.tar.xz)  
[Standalone binary](https://github.com/bytemystery-com/smartdiary/releases/download/v0.3.3/smartdiary)  
#### Windows (64 Bit)
[Standalone exe](https://github.com/bytemystery-com/smartdiary/releases/download/v0.3.3/SmartDiary.exe)  
#### Mac
Not available - it could be build but requires Mac + SDK.
#### Android 
[APK all in one](https://github.com/bytemystery-com/smartdiary/releases/download/v0.3.3/SmartDiary.apk)  
[APK only 64 bit](https://github.com/bytemystery-com/smartdiary/releases/download/v0.3.3/SmartDiary_64.apk)  

## Q & A
Q: Where is the database stored ?  
>A: Use the Info dialog :-)  
On Linux it will be located at  
~/.config/fyne/com.bytemystery.smartdiary2/smartdiary.db  
On Windows they are under  
C:\Users\<USERNAME>>\AppData\Roaming\fyne\com.bytemystery.smartdiary2\smartdiary.db  
On Android open the Info dialog at the bottom the path is shown.  

## Statistics
The project consists of round about 3600 lines of go code.  

## Links
[Readme](https://bytemystery-com.github.io/smartdiary/)  
[Repository](https://github.com/bytemystery-com/smartdiary/)  
[Issues](https://github.com/bytemystery-com/smartdiary/issues)  
[Discussions](https://github.com/bytemystery-com/smartdiary/discussions/new)  

© Copyright Reiner Pröls, 2026

