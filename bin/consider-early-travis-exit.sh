function fooDbg {
	echo "----- begin: $@"
  $@
	echo "------- end: $@"
  echo " "
}

# Exits with status 0 if it can be determined that the
# current PR should not trigger all travis checks.
#
# This could be done with a "git ...|grep -vqE" oneliner
# but as travis triggering is refined it's useful to check
# travis logs to see how branch files were considered.
function consider-early-travis-exit {
  if [ "$TRAVIS_PULL_REQUEST" == "false" ]; then
    echo "Unknown pull request."
    return
  fi

	fooDbg pwd
	fooDbg git diff --name-only HEAD origin/master
	fooDbg cat bin/consider-early-travis-exit.sh
	fooDbg ls -C1
	fooDbg git status
	fooDbg git branch
	fooDbg git remote -v

  echo "TRAVIS_COMMIT_RANGE=$TRAVIS_COMMIT_RANGE"
  echo "---"
  git diff --name-only $TRAVIS_COMMIT_RANGE
  echo "---"
  echo "Branch Files (X==invisible to travis):"
  echo "---"
  local triggers=0
  local invisibles=0
  for fn in $(git diff --name-only HEAD origin/master); do
    if [[ "$fn" =~ (\.md$)|(^docs/) ]]; then
      echo "  X  $fn"
      let invisibles+=1
    else
      echo "     $fn"
      let triggers+=1
    fi
  done
  echo "---"
  printf >&2 "%6d files invisible to travis.\n" $invisibles
  printf >&2 "%6d files trigger travis.\n" $triggers
  if [ $triggers -eq 0 ]; then
    echo "Exiting travis early."
    # see https://github.com/travis-ci/travis-build/blob/master/lib/travis/build/templates/header.sh
    travis_terminate 0
  fi
}
consider-early-travis-exit
unset -f consider-early-travis-exit
