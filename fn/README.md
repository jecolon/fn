# fn
Command line tool to fix filenames for shell scripting and overall happines. :)

## Usage
```
fn [options]
-d, --dir    Directory to process. Default is current directory.
-m, --move   Move (rename) instead of copy files.
-o, --out    Output directory (relative to --dir) to save copies.
-r, --report Only report the changes that would be made (dry run).
```

## Fixing Behavior
fn does the following to filenames:
- Removes control characters, so even filenames with newlines are fixed.
- Replace spaces with underscores.
- Remove leading or trailng whitespace.
- Remove leading or trailing dashes.
- Remove special characters such as *{},;:/\<>!$ and many more.
- Limit filename length to 255 characters.
- Identify tricky filenames consisting of only whitespace by renaming to FN_NO_NAME
