# SaveTabs

SaveTabs is a Chrome extension and local daemon used for keeping track of your tabs and tab 
groups in Chrome and Edge to make sure you never loose the URLs you have visited and grouped into 
tabs. 

## Feedback Welcome
Although this is a work-in-progress, I would like to see it become a useful solution for 
myself and others where we can own all our own data. 

Any feedback or questions in the form of [issues](https://github.com/mikeschinkel/savetabs/issues) or 
[discussions](https://github.com/mikeschinkel/savetabs/discussions) is welcome. 


## Prerequisites

AFAIK, the only prerequisites are:

1. macOS _(tested)_, Linux _(untested)_, or Windows WSL _(untested)_. 
2. The latest version of Go in your path
2. A recent version of [tailwindcss CLI](https://tailwindcss.com/blog/standalone-cli) in your path
3. Chrome and/or Edge browsers installed

However, as I have only ever run this on my machine I may be forgetting something. Please file 
an issue, start a discussion, and/or submit a PR if you discover something I missed.

## Running the Daemon

To build and run the daemon:

```shell
git clone https://github.com/mikeschinkel/savetabs
cd savetabs
make
```

## Installing the Browser Extension

Assuming you've cloned this Git repo as in the instructions for running the Daemon:

1. Open **Google Chrome** or **Microsoft Edge**
2. Open `chrome://extensions` or `edge://extensions`, as applicable
3. Enable Developer mode by toggling the switch:
   1. In the top-right corner _(Chrome)_, or  
   2. On the left sidebar _(Edge)_ 
4. Click the _"Load unpacked"_ button that appears
5. Navigate to the `extension` directory inside your cloned repository
6. Select the directory (don't select any files inside it, just the directory itself)
7. Click _"Select Folder" (Windows)_ or _"Open" (Mac)_
8. The extension should now appear in your extensions list
9. You can pin it to your toolbar by clicking the puzzle piece icon, or in the 
   toolbar and clicking the pin icon next to the extension

### Troubleshooting

If you see any errors during installation, check that:
* You're selecting the `./extensions` directory of the repo, aka the one that contains the 
  `manifest.json` file
* The `manifest.json` file is still properly formatted
* All referenced files in the manifest exist in the directory

If the extension still does not appear to work:
* Check the extension's error console in the extensions page
* Try disabling and re-enabling the extension
* Ensure all permissions requested by the extension are granted

### Updating the Extension

When you want to update the extension with new changes:
1. Pull the latest changes from the repository
2. Go to the extensions page in your browser
3. Find the extension and click the refresh icon
4. The extension will reload with the latest changes

## To Do

### Bugs to Fix
1. Fix partial saving of Tab Group name.
2. Fix regression for filter query on Linkset postback.
3. ~~Fix non-scrolling of longer lists of links~~
4. Disable triggering of events while editing label names
5. Update ContentURL when edited label is returned by HTML API
6. Fix "Invalid group filter foramt" for slashes ('/') in group name
7. Smooth drop outline so it doesn't flash on and off
8. Fix clickability of `<summary>` in Menu w/o disabling expand. 
9. Maintain menu option highlight after loading links

### Features To Add
1. API
   1. ~~Implement Group Type list API~~
2. Design UI
   1. Popup
      1. ~~Check API status~~  
      2. ~~Implement basic UI~~  
   2. Browse 
      1. ~~Group Types as trees~~
         1. ~~Tag Groups~~
         2. ~~Tags~~
         3. ~~Categories~~
         4. ~~Keywords~~
      2. ~~Groups per Group Types as tree branches~~
      3. ~~Links per Group as leaves~~
      4. ~~Capture from Browser~~
         1. ~~Title~~
      5. ~~On-hover for on-click modal menu~~
         1. ~~For Menu Items/Groups~~
         2. ~~For Links~~
      6. ~~Make Fixed layout w/sticky elements~~
         1. ~~Make `<nav>`, `<thead>`, `<tbody>` sticky~~ 
      7. Drag & Drop 
         1. Move links to different groups
         2. Copy links to different groups
         3. Merge group to other groups
         4. Add visual feedback to illustrate successful drop
         5. ~~Delete moved links from DOM~~
         6. Disallow dropping on same group
      8. Make more Like a web app
         1. Push back status
         2. Change URL on navigation 
         3. Restore state from URL to allow bookmarking
      8. Make Content Nav more generic
         1. Use Go types
      9. Menu Items for Links
         1. ~~Rename~~
         2. Delete
         3. Archive
         4. Merge To
      10. Capture/handle in Caretaker
          1. Screenshot
          2. Meta
          3. Content
          4. Other?
      11. Capture from browser
          1. Browser Type
          2. Other?
      12. Expand/Collapse Item for links
          1. Tabbed panel
             1. Title, Screenshot, Content
             2. URL exploded into component parts
      13. Column resizing
      14. Table Row Sorting
          1. By clicking table Headers
          2. Toolbar control to add a sort
          3. Allow multiple sort levels
      15. Table Row Filters
          1. Use Flex to allow many filters
          2. Filters
             1. Archived — Default=0
             2. Deleted — Default=0
             3. Group/Value
             4. URL parts/Values
      16. Capture from more browser events
          1. TBD
   3. Settings
3. Configuration
   1. Add daemon config file, 
   4. Allow configuring
      1. Submit frequency
      2. Sqlite DB location 
   5. Add CLI command and switches
      1. Specify DB location
4. Tests
   5. Add daemon tests
   6. Add Chrome extension tests

### Project To Do
1.Rename to Stash or Stockpile, or...?
