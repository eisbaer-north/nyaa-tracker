# nyaa-tracker
All this does at the moment is check a directory for new json files and if a new file is added it will create a tracker which will then grep the rss feed specified in the json file and check for new files in the item list. as soon as the item list contains a new torrent file, it will be downloaded into a folder specified in the json file.


TODO:
I want this to run either in a proper daemon or have it run in a ncurses app to have some output for the tracking
being able to manually add 

