# Get the unique dep names in WORKSPACE
awk '/^go_repository/ {gsub("\"",""); gsub(",",""); print}' FS="\n" RS="" WORKSPACE | \
awk '/name/ {print $3}' | \
sort | uniq > /tmp/workspace.deps

# Get the deps in use
bazel query 'kind(go_library, deps(//...)) -//...' |
awk '{split($0, part, "//");gsub("@","", part[1]); print part[1]}' | \
sort | uniq > /tmp/in-use.deps

# Find go_repositories to delete 
grep -f /tmp/in-use.deps -v /tmp/workspace.deps
