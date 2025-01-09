#!/usr/bin/env bash

MOUNT_POINT="/Volumes/Backups"
ERROR_LOG="${HOME}/savetabs-backup-error-$(date +%Y-%m-%d-%H).log"
DIR="/Volumes/Backups/Projects/savetabs"
PREFIX="savetabs"

function ensure_backups_volume_exists {
  mount | grep -q "^${MOUNT_POINT} " || {
     printf "Mount %s does not exist at %s\n" \
      "${MOUNT_POINT}" \
      "$(date)" >> "$ERROR_LOG"
     exit
  }
}
function hash {
  local file="$1"
  /usr/local/bin/sha384sum "${file}"  | awk '{print$1}'
}
function latestFile {
  #ls -lt "${DIR}" | grep "^${PREFIX}" | head -n 1 | awk '{print $NF}'
  find "${DIR}" | grep "/${PREFIX}"| sort -nr | head -n 1
}
function main {
  local srcFile
  local dstFile
  local logFile
  srcFile="${HOME}/.config/savetabs/savetabs.db"
  dstFile="${DIR}/${PREFIX}-$(date "+%Y-%m-%d_%H_%M").db"
  logFile="${DIR}log.txt"

  ensure_backups_volume_exists

  latestFile="$(latestFile)"


  srcHash="$(hash "${srcFile}")"
  latestHash="$(hash "${latestFile}")"

  if [ "${srcHash}" == "${latestHash}" ] ; then
    # No change in source file, no need to backup again
    # echo "Skipping ${srcFile} backup."
    return
  fi
  # echo "Backing up ${srcFile}"
  cp "${srcFile}" "${dstFile}" >> "${logFile}" 2>&1

}
main
