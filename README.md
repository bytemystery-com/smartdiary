![alt text](/assets/icons/icon.png "Logo")

# SmartDiary
SmartDiary is a little diary application runnig on Anroid, Linux and Windows. 
For every day you can add a text entry. For highlighting entries you can choose between colors and symbols, every color can be labeled with a custom text. There are 13 colors available.  
You can set a password and protect entries so that for viewing / modifying these entries the password must be given.  
You can search in the diary for expressions. You can export the whole diary to JSON format.  
It is also possible to import a previously exported json file and replace the existing datavase.  
SmartDiary is written in [Go](https://go.dev/) and uses [Fyne](https://fyne.io/) as graphical toolkit.

Author: Reiner Pröls  
Licence: MIT  

## Usage of SmartDiary
Open the settings. Set some properties fitting your needs. Then you have to give text for categories.  
Choose a color and in the text field right give a new text.
Press ENTER to take over the category. Do not forget to press the Save button to store the changed value.
A category without a name will be hidden in the entry edit dialog.  
You can also change the order of the categories and set a default category for new entries.  
After this you can tap on a day in the calendar and add text, choose categories, display mode and overlay icon.
Also yo can mark an entry as protected so a password must be given to see such entries.

## Screenshots
![alt text](/screenshots/main_d.jpg "Main screen (dark mode)")
![alt text](/screenshots/main_l.jpg "Main screen (light mode)")
![alt text](/screenshots/entry_d.jpg "Edit entry screen (dark mode)")
![alt text](/screenshots/entry_l.jpg "Edit entry screen (light mode)")
![alt text](/screenshots/search_d.jpg "Search screen (dark mode)")
![alt text](/screenshots/search_l.jpg "Search screen (light mode)")
![alt text](/screenshots/settings_d.jpg "Settings screen (dark mode)")
![alt text](/screenshots/settings_l.jpg "Settings screen (light mode)")
![alt text](/screenshots/info.jpg "Infos")

### Precompiled binaries
#### Linux (64 Bit)
[Tar file](https://github.com/bytemystery-com/smartdiary/releases/download/v0.3.3/SmartDiary.tar.xz)  
[Standalone binary](https://github.com/bytemystery-com/smartdiary/releases/download/v0.3.3/smartdiary)  
#### Windows (64 Bit)
[Standalone exe](https://github.com/bytemystery-com/smartdiary/releases/download/v0.3.3/SmartDiary.exe)  
#### Mac
Not available - it could be build but requires Mac + SDK.
#### Android 
[APK all in one](https://github.com/bytemystery-com/smartdiary/releases/download/v0.3.3/SmartDiary.apk)  
[APK only 64 bit](https://github.com/bytemystery-com/smartdiary/releases/download/v0.3.3/SmartDiary_64.apk)  

## Q & A
Q: What is the default password ?  
>A: It's an empty string  
Q: Where is the database stored ?  
>A: Use the Info dialog :-)  
On Linux it will be located at  
~/.config/fyne/com.bytemystery.smartdiary2/smartdiary.db  
On Windows they are under  
C:\Users\<USERNAME>>\AppData\Roaming\fyne\com.bytemystery.smartdiary2\smartdiary.db  
On Android open the Info dialog at the bottom the path is shown.  
Q: Can I change the colors of the categories ?  
>A: No, colors are choosen very carefully - so text is always visible.  
Q: Can I add a category ?  
>A: No, but in the settings you can edit the label of a category.  
Q: Can I delete a category ?  
>A: No, but in the settings you can give it an empty label. And so it is no longer displayed in the entry edit dialog.  
Q: Can I change the order of categories ?  
>A: Yes, in the settings you can move categories up and down.  
Q: What does import do ?  
>A: You can export the caegory settings and all entries in a JSON file. Them you can import it on an other
machine / mobile. Import replaces the existing database !  


## Statistics
The project consists of round about 3600 lines of go code.  

## Links
[Readme](https://bytemystery-com.github.io/smartdiary/)  
[Repository](https://github.com/bytemystery-com/smartdiary/)  
[Issues](https://github.com/bytemystery-com/smartdiary/issues)  
[Discussions](https://github.com/bytemystery-com/smartdiary/discussions/new)  

© Copyright Reiner Pröls, 2026

