# snip-IT
A Code Snippet Manager CLI Tool - Credits to my amazing girlfriend for the name!

# Installation
`bash install.sh`

# Basic Usage
```bash
$ snip-it
> get <language> <filename>
> save <language> <filename>
```

# General
- Create an array that will hold the different languages
    - Each language will hold a hashmap of file created containing the code snippets
- Data (code snippets) will be serialized into JSON to achieve persistence
- When the CLI tool starts up, data will be loaded into the array and hashmap for faster retrieval. When a snippet is added, it is loaded into the DS and then save the DS into JSON.

# TO-DO - 2023-08-16T13:40:03:
1. Create the serialization and flags - DONE
2. Create the Data Structure - DONE
3. Create the saving/loading - DONE
4. Implement Tview into the program (possibly)

# Note
Some errors I haven't caught yet exit the program, so the data.json does not have time repopulate.
I truncate the file to be able to rewrite the new data to it. So if the program exits because of an error, all data inside the json file are gone.

`install.sh` creates a temp.json. I welcome you to copy your data into it for now until I find a good solution (suggestions are welcome... I am beginner)
