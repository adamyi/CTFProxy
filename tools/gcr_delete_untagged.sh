#!/bin/bash

# Copyright [2018] [lahsivjar]
# Copyright [2019] [adamyi]
# 
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
# 
#     http://www.apache.org/licenses/LICENSE-2.0
# 
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -eou pipefail

USAGE="$(basename "$0") [-h] [-r] -- Delete all untagged images in a project for GCR repository
where:
         -h get help (quite a paradox :D)
         -r GCR repository, for example: gcr.io/<project_id>"

delete_untagged() {
    echo "   |-Deleting untagged images for $1"
    while read digest; do
       gcloud container images delete $1@$digest --quiet 2>&1 | sed 's/^/        /'
    done < <(gcloud container images list-tags $1 --filter='-tags:*' --format='get(digest)' --limit=unlimited)
}

delete_for_each_repo() {
    echo "|-Will delete all untagged images in $1"
    while read repo; do
        delete_untagged $repo
	delete_for_each_repo $repo
    done < <(gcloud container images list --repository $1 --format="value(name)")
}

while getopts ':hr:' option; do
    case "$option" in
        h) echo "$USAGE"
           exit
           ;;
        r) REPOSITORY=$OPTARG
           ;;
        *) printf "invalid usage please provide the repository name with -r flag\n" >&2
           echo "$USAGE"
           exit 1
           ;;
    esac
done
shift $((OPTIND-1))

if [ -z ${REPOSITORY-} ]; then
    printf "invalid usage please provide the repository name with -r flag\n" >&2
    echo "$USAGE";
    exit 1;
fi

delete_for_each_repo $REPOSITORY
