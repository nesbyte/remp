
# **remp**
**remp** recursively searches on a path and it's files/directories against a regex pattern. Works on Linux and macOS.

## Difference from grep
**grep** regex matches within a directory or recursively matches as it discovers directories ([grep behaviour](https://www.gnu.org/software/grep/manual/grep.html)).

**remp** will only search recursively within the path itself,
regex matching against filed/directories as it steps through the path provided.

## If you only need to find git root
If you only ever need to find a singular git root, it might be better to use `git rev-parse --show-toplevel` instead of remp.

# How to install
**remp** is only available on OSX and Linux. It could be made to support other releases if needed.

To install on linux or OSX amd64 replace `remp-darwin-arm64` below with the [release name matching your OS/Architecture from the release list](https://github.com/nesbyte/remp/releases).

Command to install on Linux or OSX ('name=**arg**' below must change according to [your OS/Architecture](https://github.com/nesbyte/remp/releases) as mentioned above):  
`name=remp-darwin-arm64 && curl -fLO https://github.com/nesbyte/remp/releases/latest/download/$name.tar.gz && tar -xzf $name.tar.gz && sudo mv -i $name-* /usr/local/bin/remp`


### Explanation of the download command
1. `name=remp-darwin-arm64` - Sets a bash variable to a name related to a specific release tarball.
2. `curl -fLO https://github.com/nesbyte/remp/releases/latest/download/$name.tar.gz` - Downloads the latest released tarball
3. `tar -xzf $name.tar.gz` - Extracts and untars the downloaded tarball
4. `sudo mv -i $name-* /usr/local/bin` - Move the extracted tarball to the /usr/local/bin area. If an existing remp installation exists, the *-i* flag will ensure that you are asked if you want to overwrite it or not.

# How **remp** works

Given a path *~/some/path/to/dir*,
**remp** will by default perform a regex match from right to left such as 
1. `~/some/path/to/dir/`
2. `~/some/path/to/`
3. `~/some/path/`
4. `~/some/`
5. etc

At each step above, regex match against all files and directories at the right-most directory. If a match is found it will return the first matched file as a full path by default.

Imagine a .git folder at `~/some/path/` from above.
typing `echo "~/some/path/to/dir/" | remp -e ".git"` will make
**remp** step through 1 and 2 and on step 3, it will exactly match against .git and print out `~/some/path/.git`

By passing in the `-b` flag as such: `echo "~/some/path/to/dir/" | remp -e ".git" -b` only the directoy path will be shown without the matched file/directory itself  `~/some/path`


# More intricate examples
1. `cd $(fzf | remp -b -O . -X ".git")`  
*Explanation:* fzf returns the selected file, remp recursively searches for the exactly matching term using the `-X` flag, *.git*. If *.git* is found, cd to that directory (`-b` shows the directory of the matched file). If no directory is found `.` will be printed to stdout using `-O .` which will cd stay in the current directory.   
*Note: fzf must [be installed.](https://github.com/tmux/tmux/wiki/Installing)*  
2. `tmuxpath=$(fzf | remp -b -O "no match" ".git" "go.mod") && tmux new -c $(echo $tmuxpath)`  
*Explanation:* fzf returns a selected file, remp recursively searches for a matching term (either .git or go.mod). If a term is found it is placed in `tmuxpath` variable, a tmux session is then created with it's root at at the path provided by `echo $tmuxpath` (`-b` shows the directory of the matched file). If no match is found, give an error (return 1) and tmux is not run (due to `&&`).  
*Note: tmux must [be installed.](https://github.com/junegunn/fzf?tab=readme-ov-file#installation)*  
3. `pwd | remp --color "go"`  
*Explanation:* Take the current path and check if there is a file on the path with the provided regex *go* pattern. If a match is found, colour the matched term(s). Return nothing if no matches have been found. 
4. `pwd | remp --color -l "go"`  
*Explanation:* Same as above but begin the search from the left instead of from right.