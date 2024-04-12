
# **remd**
**remd** recursively searches on a path and it's files/directories against a regex pattern. Works on Linux, macOS and Windows.

## Difference from grep
**grep** regex matches within a directory or recursively matches as it discovers directories ([`--recursive` flag grep](https://www.gnu.org/software/grep/manual/grep.html)).

**remd** will only search recursivly on the path itself,
regex matching against filed/directories as it steps through the path provided.

# Basic Operation 

Given a path *~/some/path/to/dir*,
**remd** will by default perform a regex match from right to left such as 
1. `~/some/path/to/dir/`
2. `~/some/path/to/`
3. `~/some/path/`
4. `~/some/`
5. etc

At each step above, regex match against all files and directories at the right-most directory. If a match is found it will return the first matched file as a full path by default.

Imagine a .git folder at `~/some/path/` from above.
typing `echo "~/some/path/to/dir/" | **remd** -e ".git"` will make
**remd** step through 1 and 2 and on step 3, it will exactly match against .git and print out `~/some/path/.git`

By passing in the `-b` flag as such: `echo "~/some/path/to/dir/" | **remd** -e ".git" -b` only the directoy path will be shown without the matched file/directory itself  `~/some/path`


# Examples
1. `cd $(fzf | remd -b -X ".git")`  
*Explanation:* fzf returns the selected file, remd recursively searches for the matching term, *.git*. If *.git* is found, cd to that directory (`-b` shows the directory of the matched file).  
*Note: In reality this command may need to be improved in the event no match is found.*
2. `tmux new -c $(fzf | remd -b ".git" "go.mod")`  
*Explanation:* fzf returns a selected file, remd recursively searches for a matching term (either .git or go.mod). If a term is found, create a tmux session with that directory (`-b` shows the directory of the matched file). If no match is found, give an error (return 1).
3. `pwd | remd --color "go"`  
*Explanation:* Take the current path and check if there is a file on the path with the provided regex *go* pattern. If a match is found, colour the matched term(s). Return nothing if no matches have been found. 
4. `pwd | remd --color -l "go"`  
*Explanation:* Same as above but begin the search from the left instead of from right.