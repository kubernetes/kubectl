# A lot of code in this file is originally developed by @monopole in:
# https://github.com/monopole/scratch/blob/master/moveResourcePackage.md#the-morning-refresh-script

# remove repos that are in k8s.io/kubernetes's staging dir.
rm -rf $GOPATH/src/k8s.io/kubectl/vendor/k8s.io/api
rm -rf $GOPATH/src/k8s.io/kubectl/vendor/k8s.io/apimachinery
rm -rf $GOPATH/src/k8s.io/kubectl/vendor/k8s.io/client-go
rm -rf $GOPATH/src/k8s.io/kubectl/vendor/k8s.io/common

cd $GOPATH/src/k8s.io/kubectl
git commit -am "remove old repos"

WORK_TARGET=$GOPATH/src/k8s.io/kubectl
mkdir -p $WORK_TARGET
BRANCH_NAME=contentMove

function copyDirectory {
  local DIR_SOURCE=$1
  local DIR_TARGET=$2
  local REMOTE_NAME=k8supstream

  # Place to clone it.
  local WORK_SOURCE=$(mktemp -d)

  git clone \
      --depth=1 \
      https://github.com/kubernetes/kubernetes \
      $WORK_SOURCE

  cd $WORK_SOURCE
  git checkout -b $BRANCH_NAME

  # Delete everything in the source repo
  # except the files to move:
  git filter-branch \
    --subdirectory-filter $DIR_SOURCE \
    -- --all

  # Show what's left
  ls

  # Move retained content to the target directory
  # in the target repo.
  mkdir -p $DIR_TARGET

  # The -k avoids the error from '*' picking
  # up the target directory itself.
  git mv -k * $DIR_TARGET

  # Commit the change locally.
  git commit -m "Isolated content of $DIR_SOURCE"

  # The repo now contains ONLY the code to copy.
  # Do the copy.
  cd $WORK_TARGET
  # echo Should be on branch $BRANCH_NAME
  git status

  git remote add $REMOTE_NAME $WORK_SOURCE
  git fetch $REMOTE_NAME
  git merge --allow-unrelated-histories \
      -m "Copying $DIR_SOURCE" \
      $REMOTE_NAME/$BRANCH_NAME
  git remote rm $REMOTE_NAME

  # Delete the traumatized `$WORK_SOURCE` directory.
  rm -rf $WORK_SOURCE

  # TODO: Use sed to adjust import statements.
}

copyDirectory staging/src/k8s.io/api          vendor/k8s.io/api
copyDirectory staging/src/k8s.io/apimachinery vendor/k8s.io/apimachinery
copyDirectory staging/src/k8s.io/client-go    vendor/k8s.io/client-go

copyDirectory pkg/kubectl/categories vendor/k8s.io/common/categories
copyDirectory pkg/kubectl/resource   vendor/k8s.io/common/resource
copyDirectory pkg/kubectl/validation vendor/k8s.io/common/validation

cd $WORK_TARGET
for ADJUST_IMPORT_TARGET in "vendor/k8s.io/api" "vendor/k8s.io/apimachinery" "vendor/k8s.io/client-go" "vendor/k8s.io/common"
do
find $ADJUST_IMPORT_TARGET -name "*_test.go" | xargs rm
done

ADJUST_IMPORT_TARGET=vendor/k8s.io
function adjustImport {
  local file=$1
  local old=$2
  local new=$3
  local c="s|\\\"k8s.io/kubernetes/$old/|\\\"k8s.io/common/$new/|"
  sed -i $c $file
  local c="s|\\\"k8s.io/kubernetes/$old\\\"|\\\"k8s.io/common/$new\\\"|"
  sed -i $c $file
}
function adjustAllImports {
  for i in $(find $ADJUST_IMPORT_TARGET -name '*.go' );
  do
    adjustImport $i $1 $2
  done
}

adjustAllImports pkg/kubectl/categories categories
adjustAllImports pkg/kubectl/resource   resource
adjustAllImports pkg/kubectl/validation validation

git commit -am "remove test files and adjust import"
