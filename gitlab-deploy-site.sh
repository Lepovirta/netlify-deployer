#!/usr/bin/env bash
set -euo pipefail

GITLAB_BASE_URL="https://gitlab.com/api/v4"
GITLAB_PROJECT_URL="${GITLAB_BASE_URL}/projects/${CI_PROJECT_ID}"
GITLAB_MR_URL="${GITLAB_PROJECT_URL}/merge_requests/${CI_MERGE_REQUEST_IID}"

# Disable drafts for main branch only
NETLIFY_MAIN_BRANCH="${NETLIFY_MAIN_BRANCH:-"master"}"
NETLIFY_DRAFT="true"
if [ "${CI_COMMIT_BRANCH:-}" == "${NETLIFY_MAIN_BRANCH}" ]; then
    NETLIFY_DRAFT="false"
fi

# Get current branch name
BRANCH_NAME="${CI_COMMIT_REF_NAME:-}"
if [ -z "${BRANCH_NAME:-}" ]; then
    BRANCH_NAME="$(git rev-parse --abbrev-ref HEAD)"
fi

# Get current Git commit hash
GIT_HASH=${CI_COMMIT_SHA:-}
if [ -z "${GIT_HASH}" ]; then
    GIT_HASH=$(git rev-parse HEAD)
fi

# Message for the deployment
NETLIFY_DEPLOYMESSAGE="${BRANCH_NAME}/${GIT_HASH}"

# Make the settings available for netlify-deployer
export NETLIFY_DRAFT NETLIFY_DEPLOYMESSAGE

should_post_comment() {
    [ -n "${CI_MERGE_REQUEST_IID:-}" ] && \
    [ -n "${GITLAB_ACCESS_TOKEN}" ]
}

mr_message_template() {
    echo "[Preview site](${1}) for commit ${GIT_HASH}"
}

post_comment() {
    curl \
        -sfSL \
        --retry 3 \
        --retry-connrefused \
        --retry-delay 2 \
        --request POST \
        --header "Private-Token: $GITLAB_ACCESS_TOKEN" \
        --data-urlencode body@- \
        "${GITLAB_MR_URL}/notes" >/dev/null
}

main() {
    local preview_url

    # Deploy
    preview_url=$(netlify-deployer)

    if [ -z "${preview_url:-}" ]; then
        echo "No preview URL available" >&2
        return
    fi
    echo "Preview URL: ${preview_url}"

    # Post comment to the merge request, if it's available.
    if should_post_comment; then
        if ! mr_message_template "${preview_url}" | post_comment; then
            echo "Failed to post comment on merge request!" >&2
        fi
    fi
}

main
