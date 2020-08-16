#!/usr/bin/env zsh -e

repo=`git remote get-url origin | sed 's/.*:\(.*\).git/\1/'`
def_ref=`git rev-parse --abbrev-ref HEAD`

env_vars ()
{
  printf "Please set your Github token as an env var\n"
  printf "  export GH_TOKEN=<your-github-token>\n"
  printf "The token requires only the 'repo_deployment' claim\n"
  printf "See https://help.github.com/en/github/authenticating-to-github/creating-a-personal-access-token-for-the-command-line\n"
}

usage ()
{
  printf "Usage: deploy_to_dev [FLAGS]\n\n"
  printf "Flags:\n"
  printf "  %s %-25s %-35s %s\n" "-r" "<github branch or sha>" "Set branch or commit sha to deploy" "Default: [$def_ref]"
  printf "  %s %-25s %-35s %s\n" "-f" "" "Skip deployment confirmation" "Default: false"
  printf "\n"
  env_vars
}

deploy ()
{
  printf "Deploying %s to https://github.com/%s/actions\n\n" $1 $repo

  curl -XPOST \
  -H "Authorization: token $GH_TOKEN" \
  -H "Accept: application/vnd.github.ant-man-preview+json"  \
  -H "Content-Type: application/json" \
  "https://api.github.com/repos/${repo}/deployments" \
  --data "{\"ref\": \"${1}\"}"
}

confirm_deploy ()
{
  if read -q "answer?Deploy ref $1 of $repo (y/n)? "; then
    printf "\n"
    deploy $1
  else
    printf "\nCancelling\n"
  fi
}


if [ -z $GH_TOKEN ]; then
  env_vars >&2
  exit 1
fi

while getopts ":hr:e:f" opt; do
  case $opt in
    h)
      usage
      exit 0
      ;;
    r)
      user_ref=$OPTARG
      ;;
    f)
      force=1
      ;;
    :)
      printf "Option -%s requires an argument.\n" $OPTARG >&2
      exit 1
      ;;
    \?)
      printf "Invalid option: -%s. Please use -h for usage info.\n" $OPTARG >&2
      exit 1
      ;;
  esac
done
shift $((OPTIND-1))

if [[ $# -gt 0 ]]; then
  printf "Invalid arguments: %s. Please use -h for usage info." $@ >&2
  exit 1
fi

ref=${user_ref:-$def_ref}

if [[ -z $force ]]; then
  confirm_deploy $ref
else
  deploy $ref
fi